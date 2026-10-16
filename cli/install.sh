#!/usr/bin/env bash
# Pasang CipinCode CLI ke ~/.local (tanpa sudo).
#
#   ./install.sh               pasang atau perbarui
#   ./install.sh --uninstall   hapus
#
# Variabel opsional: CIPINCODE_PREFIX (bawaan ~/.local), CIPINCODE_VERSION.
set -euo pipefail

PREFIX="${CIPINCODE_PREFIX:-$HOME/.local}"
VERSION="${CIPINCODE_VERSION:-0.2.2}"
MIN_GO="1.24"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SHARE="$PREFIX/share/cipincode"
BIN="$PREFIX/bin"

say() { printf '\033[1;38;5;166m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[33mPeringatan:\033[0m %s\n' "$*" >&2; }
fail() { printf '\033[31mError:\033[0m %s\n' "$*" >&2; exit 1; }

if [ "${1:-}" = "--uninstall" ]; then
    rm -rf "$SHARE"
    rm -f "$BIN/cipincode"
    say "CipinCode CLI dihapus dari $PREFIX"
    echo "Data login (~/.config/cipincode) tidak dihapus. Hapus manual jika perlu."
    exit 0
fi

# 1. Prasyarat
command -v go >/dev/null 2>&1 || fail "Go belum terpasang. Unduh di https://go.dev/dl/ (butuh $MIN_GO atau lebih baru)."

go_version="$(go env GOVERSION | sed 's/^go//')"
if [ "$(printf '%s\n%s\n' "$MIN_GO" "$go_version" | sort -V | head -n1)" != "$MIN_GO" ]; then
    fail "Go $go_version terlalu lama. Butuh $MIN_GO atau lebih baru."
fi

node_bin="$(command -v node 2>/dev/null || true)"
bun_bin=""
for candidate in "$(command -v bun 2>/dev/null || true)" "$HOME/.bun/bin/bun" "$HOME/.local/bin/bun"; do
    if [ -n "$candidate" ] && [ -x "$candidate" ]; then
        bun_bin="$candidate"
        break
    fi
done
# CLI mengutamakan Node (dist/server.node.js) lalu fallback ke Bun (src/server.ts).
# Jadi cukup salah satu yang ada.
[ -n "$node_bin" ] || [ -n "$bun_bin" ] || fail "Butuh Node.js 18+ dari https://nodejs.org atau Bun dari https://bun.sh"

if [ -n "$node_bin" ]; then
    say "Go $go_version, Node $("$node_bin" --version)"
else
    say "Go $go_version, Bun $("$bun_bin" --version) (Node tidak ada, memakai Bun)"
fi

# 2. Build dan salin
say "Membangun CipinCode CLI $VERSION"
mkdir -p "$SHARE/bin" "$BIN"
(cd "$ROOT/tui" && go build -ldflags "-X main.version=$VERSION" -o "$SHARE/bin/cipincode" .)

# Pastikan dist/server.node.js ada agar mode Node (tanpa Bun) bisa jalan.
# Jika belum ada, coba build dengan tsc via node atau bun.
if [ ! -f "$ROOT/agent/dist/server.node.js" ]; then
    say "Membangun agent (dist/server.node.js)..."
    if [ -n "$node_bin" ] && command -v npx >/dev/null 2>&1; then
        (cd "$ROOT/agent" && npx -y -p typescript tsc -p tsconfig.node.json) || warn "Gagal build dist via npx tsc; fallback ke Bun saat jalan."
    elif [ -n "$bun_bin" ]; then
        (cd "$ROOT/agent" && "$bun_bin" install && "$bun_bin" x tsc -p tsconfig.node.json) || warn "Gagal build dist via bun; fallback ke Bun saat jalan."
    else
        warn "dist/server.node.js tidak ada dan tidak bisa dibangun tanpa npx/bun; CLI akan memakai Bun saat jalan."
    fi
fi

say "Menyalin agent"
rm -rf "$SHARE/agent"
mkdir -p "$SHARE/agent"
# Salin lengkap (src + dist bila ada) agar mode Node dan Bun sama-sama bisa jalan.
cp -R "$ROOT/agent/src" "$SHARE/agent/"
[ -d "$ROOT/agent/dist" ] && cp -R "$ROOT/agent/dist" "$SHARE/agent/"
cp "$ROOT/agent/package.json" "$SHARE/agent/"
[ -f "$ROOT/agent/tsconfig.node.json" ] && cp "$ROOT/agent/tsconfig.node.json" "$SHARE/agent/"
if [ -d "$ROOT/agent/node_modules" ]; then
    cp -R "$ROOT/agent/node_modules" "$SHARE/agent/"
else
    (cd "$SHARE/agent" && (npm install --production 2>/dev/null || true))
fi

ln -sf "$SHARE/bin/cipincode" "$BIN/cipincode"

# 3. Verifikasi
say "Verifikasi: $("$BIN/cipincode" --version)"

case ":$PATH:" in
    *":$BIN:"*) ;;
    *) warn "$BIN belum ada di PATH. Tambahkan ke ~/.bashrc atau ~/.zshrc:
    export PATH=\"$BIN:\$PATH\"" ;;
esac

cat <<'MSG'

Selesai. Langkah berikutnya:
  1. Buat API key di dashboard (menu API Keys)
  2. cipincode login --key <API_KEY>
  3. cipincode
MSG

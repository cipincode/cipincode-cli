#!/usr/bin/env bash
# Pasang CipinCode CLI tanpa npm maupun Go: unduh binary + agent dari GitHub Releases.
#
#   curl -fsSL https://cli.cipincode.com/install.sh | bash
#
set -euo pipefail

DOMAIN="${CIPINCODE_URL:-https://cli.cipincode.com}"
REPO="cipincode/cipincode-cli"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CIPIN_HOME="${CIPINCODE_HOME:-$HOME/.local/share/cipincode}"

info() { printf "\033[1;38;5;166m==>\033[0m %s\n" "$1"; }
error() { printf "\033[31mError:\033[0m %s\n" "$1"; exit 1; }

command -v curl >/dev/null 2>&1 || error "curl belum terpasang."
command -v tar >/dev/null 2>&1 || error "tar belum terpasang."

# Deteksi OS/arch lalu petakan ke nama asset di GitHub Releases.
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$OS" in
  linux) GOOS="linux" ;;
  darwin) GOOS="darwin" ;;
  *) error "OS tidak didukung: $OS (unduh manual di https://github.com/$REPO/releases/latest)" ;;
esac
case "$ARCH" in
  x86_64|amd64) GOARCH="amd64" ;;
  arm64|aarch64) GOARCH="arm64" ;;
  *) error "Arsitektur tidak didukung: $ARCH" ;;
esac
ASSET="cipincode-${GOOS}-${GOARCH}"
info "Mendeteksi sistem: ${GOOS}-${GOARCH}"

# Versi terbaru dari server (sumber kebenaran untuk modal update di CLI).
# Fallback ke tag rilis GitHub terbaru bila server tidak bisa dihubungi.
CLI_VERSION="$(curl -fsSL "$DOMAIN/api/cli/version" 2>/dev/null | grep -o '"version":"[^"]*"' | cut -d'"' -f4 || true)"
if [ -z "$CLI_VERSION" ]; then
  TAG_URL="$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/$REPO/releases/latest" 2>/dev/null || true)"
  CLI_VERSION="${TAG_URL##*/tag/}"
  CLI_VERSION="${CLI_VERSION#v}"
fi
[ -n "$CLI_VERSION" ] || error "Tidak bisa menentukan versi terbaru."
info "Versi terbaru: $CLI_VERSION"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

info "Mengunduh binary $ASSET v$CLI_VERSION..."
curl -fsSL -o "$TMP_DIR/cipincode" "https://github.com/$REPO/releases/download/v$CLI_VERSION/$ASSET" || error "Unduhan binary gagal. Cek https://github.com/$REPO/releases/tag/v$CLI_VERSION"
chmod +x "$TMP_DIR/cipincode"

info "Mengunduh agent..."
curl -fsSL -o "$TMP_DIR/agent.tar.gz" "https://github.com/$REPO/releases/download/v$CLI_VERSION/agent.tar.gz" || error "Unduhan agent gagal."

mkdir -p "$CIPIN_HOME/bin" "$INSTALL_DIR"
mv "$TMP_DIR/cipincode" "$CIPIN_HOME/bin/cipincode"
rm -rf "$CIPIN_HOME/agent"
mkdir -p "$CIPIN_HOME/agent"
tar -xzf "$TMP_DIR/agent.tar.gz" -C "$CIPIN_HOME/agent"

ln -sf "$CIPIN_HOME/bin/cipincode" "$INSTALL_DIR/cipincode"

info "Verifikasi: $("$INSTALL_DIR/cipincode" --version)"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) info "Tambahkan ke PATH di ~/.bashrc atau ~/.zshrc: export PATH=\"$INSTALL_DIR:\$PATH\"" ;;
esac

if ! command -v node >/dev/null 2>&1 && ! command -v bun >/dev/null 2>&1; then
  info "Catatan: Node.js 18+ (atau Bun) hanya dibutuhkan bila memakai agent lokal (CIPINCODE_LOCAL_AGENT=1). Bawaan memakai agent remote sehingga langsung bisa jalan."
fi

info "Selesai! Jalankan: cipincode"

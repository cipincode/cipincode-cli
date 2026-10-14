#!/usr/bin/env bash
set -e

DOMAIN="https://cli.cipincode.com"
INSTALL_DIR="$HOME/.local/bin"
CIPIN_HOME="$HOME/.local/share/cipincode"

info() { printf "\033[1;38;5;166m==>\033[0m %s\n" "$1"; }
error() { printf "\033[31mError:\033[0m %s\n" "$1"; exit 1; }

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux) GOOS="linux" ;;
  darwin) GOOS="darwin" ;;
  *) error "OS tidak didukung: $OS" ;;
esac

case "$ARCH" in
  x86_64|amd64) GOARCH="amd64" ;;
  arm64|aarch64) GOARCH="arm64" ;;
  *) error "Arsitektur tidak didukung: $ARCH" ;;
esac

info "Mendeteksi sistem: $GOOS-$GOARCH"

# Ambil versi terbaru dari server agar binary melaporkan versi yang benar.
# Gagal pun tidak apa-apa: bangun tanpa versi (bawaan 0.1.0).
CLI_VERSION="$(curl -fsSL "$DOMAIN/api/cli/version" 2>/dev/null | grep -o '"version":"[^"]*"' | cut -d'"' -f4)"
if [ -z "$CLI_VERSION" ]; then
  info "Tidak bisa mengecek versi terbaru, memakai bawaan."
  CLI_VERSION="0.1.0"
else
  info "Versi terbaru: $CLI_VERSION"
fi

if ! command -v go &>/dev/null; then
  info "Go belum terpasang. Silakan install Go terlebih dahulu dari https://go.dev/dl/"
  exit 1
fi

# CLI mengutamakan Node (dist/server.node.js), fallback ke Bun (src/server.ts).
# Cukup salah satu; pasang Node bila keduanya belum ada.
if ! command -v node &>/dev/null && ! command -v bun &>/dev/null; then
  if command -v brew &>/dev/null; then
    info "Menginstall Node via brew..."
    brew install node
  else
    info "Silakan install Node.js 18+ dari https://nodejs.org lalu jalankan ulang installer ini."
    exit 1
  fi
fi

if ! command -v bun &>/dev/null; then
  export PATH="$HOME/.bun/bin:$PATH"
fi
command -v node &>/dev/null && info "Node $(node --version)" || true
command -v bun &>/dev/null && info "Bun $(bun --version)" || info "Bun tidak ada (opsional, fallback saja)"

info "Mendownload dan memasang CipinCode CLI..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$CIPIN_HOME/bin"
mkdir -p "$CIPIN_HOME/agent"

TMP_DIR=$(mktemp -d)
git clone https://github.com/cipincode/cipincode-cli.git "$TMP_DIR"
cd "$TMP_DIR/cli/tui"
go build -ldflags "-X main.version=$CLI_VERSION" -o "$CIPIN_HOME/bin/cipincode" .

# Pastikan dist/server.node.js ada agar mode Node jalan tanpa Bun.
if [ ! -f "$TMP_DIR/cli/agent/dist/server.node.js" ]; then
  info "Membangun agent (dist)..."
  (cd "$TMP_DIR/cli/agent" && (npx -y -p typescript tsc -p tsconfig.node.json || bun x tsc -p tsconfig.node.json)) || info "Gagal build dist; CLI akan memakai Bun saat jalan."
fi

# Salin folder agent langsung dari repo yang di-clone agar 100% lengkap (termasuk src & dist).
# node_modules tidak ikut agar ringan.
rm -rf "$CIPIN_HOME/agent"
mkdir -p "$CIPIN_HOME/agent"
cp -R "$TMP_DIR/cli/agent/src" "$CIPIN_HOME/agent/"
[ -d "$TMP_DIR/cli/agent/dist" ] && cp -R "$TMP_DIR/cli/agent/dist" "$CIPIN_HOME/agent/"
cp "$TMP_DIR/cli/agent/package.json" "$CIPIN_HOME/agent/"
[ -f "$TMP_DIR/cli/agent/tsconfig.node.json" ] && cp "$TMP_DIR/cli/agent/tsconfig.node.json" "$CIPIN_HOME/agent/"
(cd "$CIPIN_HOME/agent" && (npm install --production 2>/dev/null || true))

rm -rf "$TMP_DIR"

ln -sf "$CIPIN_HOME/bin/cipincode" "$INSTALL_DIR/cipincode"

# Tambahkan export otomatis ke zshrc / bashrc jika belum ada
PROFILE="$HOME/.zshrc"
[ -f "$HOME/.bashrc" ] && PROFILE="$HOME/.bashrc"

if ! grep -q "CIPINCODE_AGENT_DIR" "$PROFILE" 2>/dev/null; then
  info "Menambahkan CIPINCODE_AGENT_DIR ke $PROFILE..."
  echo "" >> "$PROFILE"
  echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$PROFILE"
  echo "export CIPINCODE_AGENT_DIR=\"$CIPIN_HOME/agent\"" >> "$PROFILE"
fi

info "Selesai! Jalankan: cipincode"

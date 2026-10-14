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

if ! command -v go &>/dev/null; then
  info "Go belum terpasang. Silakan install Go terlebih dahulu dari https://go.dev/dl/"
  exit 1
fi

if ! command -v bun &>/dev/null; then
  info "Menginstall Bun..."
  curl -fsSL https://bun.sh/install | bash
  export PATH="$HOME/.bun/bin:$PATH"
fi

info "Mendownload dan memasang CipinCode CLI..."
mkdir -p "$INSTALL_DIR"
mkdir -p "$CIPIN_HOME/bin"
mkdir -p "$CIPIN_HOME/agent"

TMP_DIR=$(mktemp -d)
git clone https://github.com/cipincode/cipincode-cli.git "$TMP_DIR"
cd "$TMP_DIR/cli/tui"
go build -o "$CIPIN_HOME/bin/cipincode" .
rm -rf "$TMP_DIR"

ln -sf "$CIPIN_HOME/bin/cipincode" "$INSTALL_DIR/cipincode"

info "Mengunduh agent..."
curl -fsSL "$DOMAIN/downloads/agent.tar.gz" -o "/tmp/agent.tar.gz"
tar -xzf "/tmp/agent.tar.gz" -C "$CIPIN_HOME/agent"
rm -f "/tmp/agent.tar.gz"

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

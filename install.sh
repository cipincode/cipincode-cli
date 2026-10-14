#!/usr/bin/env bash
set -e

REPO="cipincode/cipincode-cli"
INSTALL_DIR="$HOME/.local/bin"
CIPIN_HOME="$HOME/.local/share/cipincode"

info() { printf "\033[1;38;5;166m==>\033[0m %s\n" "$1"; }
error() { printf "\033[31mError:\033[0m %s\n" "$1"; exit 1; }

# 1. Detect OS & Arch
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

# 2. Cek / Install Go
if ! command -v go &>/dev/null; then
  info "Go tidak ditemukan. Menginstall Go..."
  GO_VERSION="go1.24.0"
  GO_TAR="${GO_VERSION}.${GOOS}-${GOARCH}.tar.gz"
  curl -fsSL "https://go.dev/dl/$GO_TAR" -o "/tmp/$GO_TAR"
  rm -rf "$HOME/.local/go"
  tar -C "$HOME/.local" -xzf "/tmp/$GO_TAR"
  rm -f "/tmp/$GO_TAR"
  export PATH="$HOME/.local/go/bin:$PATH"
fi

# 3. Cek / Install Bun
if ! command -v bun &>/dev/null; then
  info "Bun tidak ditemukan. Menginstall Bun..."
  curl -fsSL https://bun.sh/install | bash
  export PATH="$HOME/.bun/bin:$PATH"
fi

# 4. Download binary cipincode terbaru dari GitHub Releases
info "Mengunduh CipinCode CLI..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$LATEST_TAG" ]; then
  LATEST_TAG="v0.1.0"
fi

BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/cipincode-$GOOS-$GOARCH"
AGENT_URL="https://github.com/$REPO/releases/download/$LATEST_TAG/agent.tar.gz"

mkdir -p "$INSTALL_DIR"
mkdir -p "$CIPIN_HOME/bin"
mkdir -p "$CIPIN_HOME/agent"

curl -fsSL "$BINARY_URL" -o "$CIPIN_HOME/bin/cipincode"
chmod +x "$CIPIN_HOME/bin/cipincode"

# Buat symlink ke ~/.local/bin/cipincode
ln -sf "$CIPIN_HOME/bin/cipincode" "$INSTALL_DIR/cipincode"

# Download agent
info "Mengunduh agent..."
curl -fsSL "$AGENT_URL" -o "/tmp/agent.tar.gz"
tar -xzf "/tmp/agent.tar.gz" -C "$CIPIN_HOME/agent"
rm -f "/tmp/agent.tar.gz"

# 5. Cek PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "⚠️  Tambahkan baris ini ke ~/.bashrc atau ~/.zshrc:"
    echo "    export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo "    export CIPINCODE_AGENT_DIR=\"$CIPIN_HOME/agent\""
    ;;
esac

info "Selesai! Jalankan: cipincode --version"

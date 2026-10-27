#!/usr/bin/env sh
# Cipin CLI installer (Linux / macOS). Installs straight from GitHub, no npm publish needed.
#   curl -fsSL https://raw.githubusercontent.com/cipincode/cipincode-cli/main/cipin-cli/install.sh | sh
# Optional: CIPIN_REF=branch-or-tag (default: main), CIPIN_DIR=install folder.
set -eu

for cmd in node npm curl tar; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "$cmd tidak ditemukan. Pasang dulu (Node.js 18 atau lebih baru: https://nodejs.org)." >&2
    exit 1
  fi
done

REF="${CIPIN_REF:-main}"
DEST="${CIPIN_DIR:-$HOME/.cipin/cli}"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Node: $(node -v)"
echo "Mengunduh Cipin CLI ($REF) dari GitHub..."
curl -fsSL "https://github.com/cipincode/cipincode-cli/archive/$REF.tar.gz" | tar xz -C "$TMP"

SRC="$(find "$TMP" -mindepth 2 -maxdepth 2 -type d -name cipin-cli | head -n 1)"
if [ -z "$SRC" ]; then
  echo "Folder cipin-cli tidak ditemukan di GitHub (ref: $REF)." >&2
  exit 1
fi

rm -rf "$DEST"
mkdir -p "$(dirname "$DEST")"
cp -R "$SRC" "$DEST"

cd "$DEST"
echo "Memasang dependensi (OpenCode)..."
npm install --omit=dev --no-audit --no-fund
# A previous install leaves a registration behind (and its folder was just replaced): drop it first.
npm uninstall -g cipin-cli >/dev/null 2>&1 || true
npm install -g .

cat <<EOF

Selesai. Terpasang di $DEST
Langkah berikutnya:
  cipin login
  cipin
EOF

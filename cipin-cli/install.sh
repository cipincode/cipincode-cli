#!/usr/bin/env sh
set -eu

for cmd in node npm; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "$cmd tidak ditemukan. Pasang Node.js 18 atau lebih baru dulu: https://nodejs.org" >&2
    exit 1
  fi
done

echo "Node: $(node -v)"
echo "Memasang cipin-cli..."
npm install -g cipin-cli

cat <<'EOF'

Selesai. Langkah berikutnya:
  cipin login
  cipin
EOF

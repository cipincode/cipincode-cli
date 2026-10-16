.PHONY: build dev test install clean

VERSION ?= 0.2.4

# Build binary ke cli/bin/cipincode (plus dist/server.node.js agar mode Node jalan tanpa Bun)
build:
	cd agent && (bun install || npm install --no-audit --no-fund)
	cd agent && (bun x tsc -p tsconfig.node.json || npx -y -p typescript tsc -p tsconfig.node.json)
	cd tui && go build -ldflags "-X main.version=$(VERSION)" -o ../bin/cipincode .

# Jalankan langsung tanpa build
dev:
	cd tui && go run . --agent-dir ../agent

test:
	cd agent && bun test
	cd tui && go test ./...

clean:
	rm -rf bin

# Pasang ke ~/.local (jalankan ./install.sh --uninstall untuk menghapus)
install:
	./install.sh

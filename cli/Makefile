.PHONY: build dev test install clean

VERSION ?= 0.1.0

# Build binary ke cli/bin/cipincode
build:
	cd agent && bun install
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

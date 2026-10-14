# CipinCode CLI

CLI agent berbasis AI untuk developer.

## Instalasi

```bash
npm install -g cipincode
```

## Penggunaan

```bash
cipincode login
cipincode
```

## Development

### Requirement

- Go 1.24+
- Bun

### Build

```bash
cd cli/tui
go build -o cipincode .
```

### Agent

```bash
cd cli/agent
bun install
bun run start
```

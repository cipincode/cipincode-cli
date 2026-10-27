# Cipin CLI

Agen coding AI di terminal untuk [Cipincode](https://cipincode.com). Dibangun di atas [OpenCode](https://opencode.ai) dengan provider Cipin sudah terpasang.

## Instalasi

Butuh Node.js 18 atau lebih baru. Dipasang langsung dari GitHub (folder `cipin-cli/`), tanpa npm publish.

```powershell
# Windows (PowerShell)
irm https://raw.githubusercontent.com/cipincode/cipincode-cli/main/cipin-cli/install.ps1 | iex
```

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/cipincode/cipincode-cli/main/cipin-cli/install.sh | sh
```

Terpasang di `~/.cipin/cli`. Untuk memperbarui, jalankan perintah yang sama lagi.

Opsional: `CIPIN_REF` (branch atau tag, default `main`) dan `CIPIN_DIR` (folder pasang).

## Pemakaian

```bash
cipin login     # masuk lewat browser, tempel kode XXXX-XXXX
cipin           # buka UI di folder proyek
cipin whoami    # akun dan paket
cipin topup     # buka halaman top up
cipin logout
```

Cipin CLI hanya menjawab dengan **paket Unlimited** atau paket Jailbreak yang aktif. Top up di https://cipincode.com/mbr/dashboard/topup.

Variabel opsional: `CIPIN_API_URL` (default `https://cipincode.com`).

## Lisensi

MIT. OpenCode adalah proyek terpisah dengan lisensinya sendiri.

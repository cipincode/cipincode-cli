# cipincode

Pasang CipinCode CLI dengan satu perintah.

```bash
npm install -g cipincode
```

Setelah install, jalankan:

```bash
cipincode --version
cipincode login
cipincode
```

## Yang di-install otomatis

- **Go** ≥ 1.24 (jika belum ada)
- **Bun** (jika belum ada)
- **Binary cipincode** (sesuai OS & arsitektur)
- **Agent** (server TypeScript yang dijalankan Bun)

## Platform yang didukung

| OS      | Arsitektur      |
|---------|-----------------|
| Linux   | x64, arm64      |
| macOS   | x64, arm64      |
| Windows | x64             |

## Variabel lingkungan opsional

| Variabel              | Default   | Keterangan                         |
|-----------------------|-----------|------------------------------------|
| `CIPINCODE_VERSION`   | `0.1.0`   | Paksa versi tertentu saat install  |

## Uninstall

```bash
npm uninstall -g cipincode
```

Go dan Bun yang sudah terpasang sebelumnya tidak dihapus. Hapus manual jika perlu.

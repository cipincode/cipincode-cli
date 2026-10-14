# CipinCode CLI

AI agent untuk terminal, gaya OpenCode. Dua bagian yang saling bicara lewat HTTP:

| Folder   | Bahasa           | Tugas                                             |
| -------- | ---------------- | ------------------------------------------------- |
| `tui/`   | Go (Bubble Tea)  | Tampilan terminal dan perintah `cipincode`        |
| `agent/` | TypeScript (Bun) | Server agent: menerima pesan, mengalirkan balasan |

`cipincode` menjalankan server agent di port acak (dilindungi token acak), lalu menghentikannya saat keluar.

```
CLI (Go) ──▶ agent (Bun) ──▶ dashboard (Laravel) ──▶ OpenRouter
              lokal           /api/cli/chat            model AI
```

Kunci OpenRouter hanya ada di server dashboard. CLI dan agent hanya memegang API key pengguna,
sehingga jatah pemakaian dihitung dan ditegakkan server.

## Kebutuhan

- Go 1.24 atau lebih baru
- Bun 1.x

## Pasang

```bash
cd cli
./install.sh               # pasang atau perbarui, tanpa sudo
./install.sh --uninstall   # hapus
```

Skrip memeriksa Go dan Bun, membangun binary, lalu memasangnya ke `~/.local`
(`bin/cipincode` dan `share/cipincode/agent`). Hasilnya mandiri: folder proyek boleh
dipindah atau dihapus. Pastikan `~/.local/bin` ada di PATH. Ganti lokasi dengan
`CIPINCODE_PREFIX`, ganti versi dengan `CIPINCODE_VERSION`.

## Perintah

```bash
make build   # hasil: cli/bin/cipincode
make dev     # jalankan tanpa build
make test    # tes agent (bun) dan tui (go)

cipincode --version
cipincode                      # buka sesi interaktif

cipincode login                # login interaktif: Email atau Google, tempel kode
cipincode login --key <API_KEY> [--url <alamat dashboard>]
cipincode whoami               # cek siapa yang sedang login
cipincode logout
```

### Login dengan URL dan kode

Saat pertama dijalankan (belum login), CLI menampilkan animasi sambutan lalu layar **Pilih cara masuk**:
Email atau Google. Setelah itu:

1. CLI menampilkan URL (dan mencoba membukanya di browser; `Ctrl+O` untuk membuka lagi).
2. Masuk di browser dengan metode pilihan, lalu klik **Izinkan**.
3. Dashboard menampilkan kode `XXXX-XXXX`. Tempel ke CLI.

Kode berlaku 5 menit dan hanya bisa dipakai sekali. Pertukarannya memakai PKCE: CLI membuat
rahasia acak di komputer Anda dan hanya hash-nya yang dikirim lewat URL, sehingga kode yang
bocor tidak berguna tanpa proses CLI yang meminta. Hasilnya API key baru bernama `CLI · <hostname>`
di menu API Keys. Setel `CIPINCODE_NO_BROWSER=1` untuk menonaktifkan pembukaan browser otomatis.

**Google** hanya muncul di halaman login jika kredensialnya diisi di `.env`:

```
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
```

Buat OAuth client (jenis Web application) di Google Cloud Console, lalu daftarkan
`<APP_URL>/auth/google/callback` sebagai Authorized redirect URI. Akun dengan verifikasi dua
langkah (2FA) tetap harus masuk lewat email dan password.

### Menu perintah

Ketik `/` di kolom input untuk membuka menu. Panah atas/bawah memilih, Tab melengkapi, Enter
menjalankan, Esc menutup. Mengetik melanjutkan penyaringan (`/co` menyisakan `/coding`).

| Perintah                      | Fungsi                             |
| ----------------------------- | ---------------------------------- |
| `/help`                       | bantuan                            |
| `/chat` `/coding` `/debug`    | pindah mode                        |
| `/new` `/clear`               | percakapan baru / bersihkan layar  |
| `/status` `/whoami`           | status akun, agent, mode           |
| `/usage`                      | sisa jatah sesi 5 jam dan mingguan |
| `/login` `/logout`            | masuk atau keluar                  |
| `/dashboard` `/keys` `/setup` | buka halaman di browser            |
| `/version` `/exit`            | versi / keluar                     |

### Membaca, menggulir, dan menyalin

- Jawaban AI dirender sebagai Markdown: **teks penting** tebal dan berwarna, `kode inline` berlatar,
  judul, daftar, kutipan, tabel, dan blok kode sebagai **kotak tersendiri** dengan label bahasa dan
  pewarna sintaks.
- Gulir percakapan dengan roda mouse, panah atas/bawah, atau PgUp/PgDn. Saat Anda membaca ke atas,
  jawaban yang masih mengalir tidak menarik layar ke bawah.
- Seret mouse di area percakapan untuk memilih teks; begitu tombol dilepas teksnya **langsung
  tersalin** dan muncul notifikasi sebentar di kanan atas kolom input. Seretan yang dimulai di dalam
  kotak kode hanya memilih isi kodenya (bingkai dan label tidak ikut). Jika `xclip`/`xsel`/`wl-copy`
  tidak ada, CLI memakai OSC 52 (didukung banyak terminal, juga lewat SSH).
- Karena mouse dipakai CLI, seleksi bawaan terminal tersedia dengan menahan **Shift** saat menyeret.

### Agent (`/agent`)

Buat agent di dashboard (menu **Agents**): nama, instruksi, serta mode dan effort bawaan. Di CLI,
`/agent` membuka daftar pilihan (↑↓, Enter). Instruksi agent ditambahkan ke prompt sistem setiap
pesan. Pilih "Tanpa agent" untuk kembali ke pengaturan bawaan.

### Riwayat sesi

Setiap percakapan (dimulai ulang dengan `/new` atau `/clear`) tercatat di server dan tampil di menu
**Sessions** pada dashboard, lengkap dengan pencarian, filter mode, dan detail pesan.

### Modal (`/usage`, `/status`, `/help`)

Perintah yang menampilkan informasi panjang membuka **modal** di tengah layar, di atas percakapan.
Selama modal terbuka, ketikan diabaikan. Tutup dengan **Esc** (Enter dan `q` juga bisa); panah atau
PgUp/PgDn menggulir jika isinya panjang; Ctrl+C tetap keluar dari CLI. Hasilnya tidak ditambahkan
ke percakapan. `/whoami` dan `/version` cukup satu baris, jadi tetap muncul di percakapan.

### Pemakaian (`/usage`)

Batas dihitung server dalam satuan pesan dan bergantung pada **paket** pengguna:

| Paket    | Model                  | Batas sesi 5 jam | Batas mingguan | Masa berlaku                    |
| -------- | ---------------------- | ---------------- | -------------- | ------------------------------- |
| Gratis   | model gratis (`:free`) | 20 (bawaan)      | tidak ada      | masa percobaan **14 hari**      |
| Berbayar | model berbayar         | 40 (bawaan)      | 400 (bawaan)   | tanpa batas waktu (atau diatur) |

- **Sesi 5 jam**: mulai saat pesan pertama, lalu **reset otomatis** 5 jam kemudian.
- **Mingguan** (paket berbayar): reset setiap Senin pukul 00.00.
- **Masa percobaan** (paket gratis): akun baru berlaku 14 hari sejak daftar. Setelah itu semua pesan
  ditolak dengan pesan untuk beralih ke paket berbayar.

Akun baru otomatis berpaket gratis. Karena belum ada pembayaran, paket diubah lewat perintah artisan:

```bash
php artisan plan:set email@contoh.com paid            # berbayar, tanpa batas waktu
php artisan plan:set email@contoh.com paid --days=30  # berbayar 30 hari
php artisan plan:set email@contoh.com free            # gratis, masa percobaan 14 hari dihitung ulang
```

Tiap batang terdiri dari 20 kotak bersekat; paket gratis hanya menampilkan batang sesi. Belum dipakai berarti tanpa kotak oranye; kotak
oranye bertambah seiring pemakaian; saat jatah habis seluruh kotak berubah **merah** dan pesan
berikutnya ditolak ("Pesan Anda tidak dikirim") sampai jendelanya reset. Ringkasan sesi juga
tampil di header. Halaman **Usage** di dashboard menampilkan batang yang sama.

Pesan dicatat oleh server saat OpenRouter menerima permintaan, dan ditolak (429) jika jatah habis,
sehingga tidak bisa dilewati dengan memanggil API langsung. Batas diatur lewat `.env`:
`USAGE_FREE_LIMIT_SESSION` (bawaan 20), `USAGE_LIMIT_SESSION` (40), dan `USAGE_LIMIT_WEEKLY` (400).

API key dibuat di menu **API Keys** pada dashboard (tampil sekali, disimpan sebagai hash).
Login tersimpan di `~/.config/cipincode/config.json` (izin 0600). Alamat dashboard bawaan
`http://localhost:8000`, bisa diubah dengan `--url` atau `CIPINCODE_URL`.

Di dalam sesi: `/help`, `/clear`, `/exit`. Enter kirim, Esc batalkan, PgUp/PgDn gulir, Ctrl+C keluar.

**Mode** (Shift+Tab berpindah berurutan, bawaan Chat). Mode aktif tampil di atas kolom input:

| Mode      | Fungsi                                  |
| --------- | --------------------------------------- |
| Chat      | tanya jawab biasa                       |
| Coding    | menulis dan mengubah kode, auto-approve |
| Debugging | mencari dan memperbaiki bug             |

Header menampilkan logo CP, nama dan email akun, serta waktu **sesi terakhir**
(disimpan di `~/.config/cipincode/state.json`).

Opsi: `--agent-dir <folder>` (atau `CIPINCODE_AGENT_DIR`) untuk menunjuk folder agent,
`--server <url>` (atau `CIPINCODE_AGENT_URL`) untuk memakai server agent yang sudah berjalan.

## Model AI (OpenRouter)

Jawaban datang dari OpenRouter lewat server dashboard. Isi di `.env` server:

```
OPENROUTER_API_KEY=sk-or-...          # dari https://openrouter.ai/keys
OPENROUTER_FREE_MODEL=nvidia/nemotron-3-super-120b-a12b:free   # model paket GRATIS (akhiran :free)
OPENROUTER_MODEL=openrouter/auto      # model paket BERBAYAR, atau id model tertentu
OPENROUTER_MODEL_CODING=              # opsional, model berbeda per mode (berbayar)
```

Paket **gratis** selalu memakai `OPENROUTER_FREE_MODEL` (bawaan kode: `openrouter/free`, pemilih otomatis
di antara model gratis; lebih lambat karena sering memilih model penalaran). Model gratis kadang kena
batas laju dari penyedia (HTTP 429); CLI menampilkan "Layanan AI sedang sibuk" dan jatah tidak
terpotong, coba lagi sebentar. Paket **berbayar** memakai `OPENROUTER_MODEL`, atau model per mode lewat
`OPENROUTER_MODEL_CHAT`, `_CODING`, `_DEBUGGING`. Tiap mode punya prompt sistem sendiri. Jika kunci belum diisi, CLI
menampilkan "OPENROUTER_API_KEY belum diisi" dan jatah tidak terpotong.

Percakapan dikirim beserta 20 giliran terakhir sebagai konteks. Pesan butuh login. Jatah
dipotong **setelah** OpenRouter menerima permintaan, jadi galat di layanan AI tidak
menghabiskan jatah. Agent belum bisa membaca atau mengubah file; mode Coding menulis kode
dalam blok markdown untuk disalin.

Untuk pengembangan tanpa jaringan: `CIPINCODE_AGENT_STUB=1` membuat agent membalas teks contoh.

## API agent

- `GET /health` → `{ ok, name, version }`
- `POST /chat` dengan `{ "message", "mode": "chat|coding|debugging", "history": [{role, content}] }` dan header `X-Agent-Token`, `X-Cipincode-Url`, `X-Cipincode-Key` → aliran SSE: `usage`, `token`, `blocked`, `error`, `done`

## API dashboard (Laravel)

- `POST /api/cli/exchange` dengan `{ code, verifier, device }` → `{ key, user }` (kode sekali pakai, 10 permintaan/menit).
- `GET /api/cli/usage` → `{ session, weekly }`, masing-masing `{ used, limit, percent, exhausted, resets_at }`.
- `POST /api/cli/usage` dengan `{ units }` mencatat pemakaian secara manual; membalas 429 (beserta snapshot) jika jatah habis.
- `POST /api/cli/chat` dengan `{ message, mode, history }` → aliran SSE `usage`, `token`, lalu `done` atau `error`. 429 jika jatah habis, 503 jika `OPENROUTER_API_KEY` belum diisi, 502 jika OpenRouter menolak, 403 (`blocked: plan`) jika masa berlaku paket habis.
- `GET /api/cli/me` dengan `Authorization: Bearer <API_KEY>` → `{ "user": { "name", "email" } }`.
  Mengembalikan 401 jika key tidak ada atau sudah dicabut.

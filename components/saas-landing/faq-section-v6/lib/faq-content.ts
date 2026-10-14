export type FaqItem = {
  id: string
  question: string
  answer: string
}

export const faqSectionMeta = {
  title: "Pertanyaan yang Sering Diajukan",
  description: "Jawaban cepat soal instalasi, kuota, paket, dan akun Anda.",
}

export const faqItems: FaqItem[] = [
  {
    id: "instalasi",
    question: "Bagaimana cara menginstal CipinCode CLI?",
    answer:
      "Salin perintah instalasi di bagian atas halaman ini ke terminal Linux atau macOS Anda. Go, Bun, dan binary CLI terpasang otomatis dalam satu langkah.",
  },
  {
    id: "masa-coba",
    question: "Apakah ada masa percobaan gratis?",
    answer:
      "Ya. Paket Gratis memberi masa percobaan 14 hari dengan 20 pesan per 5 jam tanpa batas mingguan. Paket berbayar Go, Pro, dan Max masih Segera Hadir.",
  },
  {
    id: "kuota",
    question: "Bagaimana sistem kuota pesan bekerja?",
    answer:
      "Setiap paket punya batas pesan per 5 jam dan per minggu. Contoh paket Go: 40 pesan per 5 jam dan 400 pesan per minggu. Pantau pemakaian Anda langsung dari dashboard web.",
  },
  {
    id: "pembayaran",
    question: "Kapan pembayaran online tersedia?",
    answer:
      "Pembayaran online belum aktif, jadi tombol paket berbayar masih nonaktif. Akun Anda tetap berjalan di paket Gratis sampai pembayaran dibuka.",
  },
  {
    id: "keamanan",
    question: "Apakah kode saya aman?",
    answer:
      "Ya. Sesi dan token terenkripsi, dan Anda memegang kendali penuh atas file yang diubah agen. API key bisa dicabut kapan saja dari dashboard.",
  },
]

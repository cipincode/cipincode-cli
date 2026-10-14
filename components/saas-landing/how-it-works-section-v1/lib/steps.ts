export type Step = {
  id: string
  number: string
  title: string
  description: string
}

export const steps: Step[] = [
  {
    id: "install",
    number: "01",
    title: "Instal",
    description:
      "Jalankan satu perintah instalasi di terminal Anda. Go, Bun, dan binary CLI terpasang otomatis di Linux dan macOS.",
  },
  {
    id: "connect",
    number: "02",
    title: "Hubungkan",
    description:
      "Masuk dengan akun CipinCode Anda dan sinkronkan token yang aman. Sesi, API key, dan penggunaan tetap sinkron dengan dashboard web.",
  },
  {
    id: "ship",
    number: "03",
    title: "Rilis",
    description:
      "Pilih mode, jelaskan tugasnya, dan biarkan agen menulis kode, debug, lalu verifikasi. Review diff-nya dan rilis langsung ke production.",
  },
]

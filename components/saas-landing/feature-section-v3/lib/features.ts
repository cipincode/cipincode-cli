export type FeatureIllustrationId =
  | "funnel-opportunities"
  | "dropoff-diagnosis"
  | "experiment-impact"

export type Feature = {
  id: FeatureIllustrationId
  step: string
  title: string
  description: string
  illustration: FeatureIllustrationId
}

export const features: Feature[] = [
  {
    id: "funnel-opportunities",
    step: "01",
    title: "Ikuti setiap tahap jalannya agen",
    description:
      "Mengikuti agen melewati scan, plan, dan execute sehingga kemacetan dan perilaku aneh terlihat sebelum merugikan Anda.",
    illustration: "funnel-opportunities",
  },
  {
    id: "dropoff-diagnosis",
    step: "02",
    title: "Pahami apa yang menghambat jalan",
    description:
      "Menunjukkan tepat di mana agen macet dalam alur kerja Anda dan apa yang menahan perbaikan.",
    illustration: "dropoff-diagnosis",
  },
  {
    id: "experiment-impact",
    step: "03",
    title: "Rilis hasil dari temuan hingga dampak",
    description:
      "Mengubah diff menjadi run yang bisa direview dan mengikuti hasilnya sampai terlihat dampaknya.",
    illustration: "experiment-impact",
  },
]

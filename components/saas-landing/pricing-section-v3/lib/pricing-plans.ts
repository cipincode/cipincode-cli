export type FeatureGroup = {
  title: string
  features: string[]
}

// Harga & fitur disamakan persis dengan App\Support\Plans::catalog()
// (config/usage.php) yang dipakai halaman Billing dashboard.
// Semua harga rupiah per bulan, tanpa siklus tahunan.
export type PricingPlan = {
  id: "free" | "go" | "pro" | "max"
  name: string
  description: string
  popular?: boolean
  price: string
  priceSuffix: string
  featureGroups: FeatureGroup[]
  ctaLabel: string
}

export const pricingPlans: PricingPlan[] = [
  {
    id: "free",
    name: "Gratis",
    description: "Untuk mencoba CLI.",
    price: "Gratis",
    priceSuffix: "coba 14 hari",
    ctaLabel: "Mulai Gratis",
    featureGroups: [
      {
        title: "Kuota",
        features: ["20 pesan per 5 jam", "Tanpa batas mingguan"],
      },
      {
        title: "Fitur",
        features: [
          "Semua mode dan effort",
          "Model gratis",
          "Masa percobaan 14 hari",
        ],
      },
    ],
  },
  {
    id: "go",
    name: "Go",
    description: "Untuk developer harian.",
    price: "Rp 49.000",
    priceSuffix: "per bulan",
    ctaLabel: "Pilih Go",
    featureGroups: [
      {
        title: "Kuota",
        features: ["40 pesan per 5 jam", "400 pesan per minggu"],
      },
      {
        title: "Fitur",
        features: ["Semua mode dan effort", "Model berbayar"],
      },
    ],
  },
  {
    id: "pro",
    name: "Pro",
    description: "Untuk profesional produktif.",
    popular: true,
    price: "Rp 149.000",
    priceSuffix: "per bulan",
    ctaLabel: "Pilih Pro",
    featureGroups: [
      {
        title: "Kuota",
        features: ["100 pesan per 5 jam", "1.500 pesan per minggu"],
      },
      {
        title: "Fitur",
        features: ["Semua mode dan effort", "Model berbayar"],
      },
    ],
  },
  {
    id: "max",
    name: "Max",
    description: "Untuk beban kerja berat.",
    price: "Rp 499.000",
    priceSuffix: "per bulan",
    ctaLabel: "Pilih Max",
    featureGroups: [
      {
        title: "Kuota",
        features: ["300 pesan per 5 jam", "5.000 pesan per minggu"],
      },
      {
        title: "Fitur",
        features: ["Semua mode dan effort", "Model berbayar"],
      },
    ],
  },
]

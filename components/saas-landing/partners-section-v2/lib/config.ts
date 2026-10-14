export const partnersSectionV2Content = {
  heading: "Tentang CipinCode untuk developer Indonesia",
  description:
    "CipinCode membangun agen AI terminal yang hidup di ekosistem tools modern — dari editor dan Git sampai infrastruktur tempat kode Anda berjalan di production.",
} as const

/** Number of logo cells in the grid. Each cell rotates through `partners.length / CELL_COUNT` logos. */
export const CELL_COUNT = 6

export const motion = {
  /** Pause between one cell swap and the next (seconds). */
  swapInterval: 1.6,
  /** Duration of a single logo swap (seconds). */
  swapDuration: 0.7,
} as const

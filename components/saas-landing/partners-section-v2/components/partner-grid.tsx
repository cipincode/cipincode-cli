import { cn } from "@/lib/utils"

import PartnerCell from "./partner-cell"
import { CELL_COUNT } from "../lib/config"
import { partners, type Partner } from "../lib/partners"

const cells: Partner[][] = Array.from({ length: CELL_COUNT }, (_, i) =>
  partners.filter((_, j) => j % CELL_COUNT === i)
)

function CornerMark({ className }: { className: string }) {
  return (
    <span
      data-pv2-mark
      className={cn("absolute size-3 -translate-x-1/2 -translate-y-1/2", className)}
    >
      <span className="absolute top-1/2 left-0 h-px w-full bg-foreground" />
      <span className="absolute top-0 left-1/2 h-full w-px bg-foreground" />
    </span>
  )
}

export default function PartnerGrid() {
  return (
    <div aria-hidden data-pv2-grid className="relative">
      {/* Frame */}
      <span data-pv2-line="x" className="absolute inset-x-0 top-0 h-px bg-border" />
      <span data-pv2-line="x" className="absolute inset-x-0 bottom-0 h-px bg-border" />
      <span data-pv2-line="y" className="absolute inset-y-0 left-0 w-px bg-border" />
      <span data-pv2-line="y" className="absolute inset-y-0 right-0 w-px bg-border" />

      {/* Inner rules: 2 × 3 on mobile, 3 × 2 from md */}
      <span data-pv2-line="y" className="absolute inset-y-0 left-1/2 w-px bg-border md:hidden" />
      <span data-pv2-line="x" className="absolute inset-x-0 top-1/3 h-px bg-border md:hidden" />
      <span data-pv2-line="x" className="absolute inset-x-0 top-2/3 h-px bg-border md:hidden" />
      <span data-pv2-line="y" className="absolute inset-y-0 left-1/3 hidden w-px bg-border md:block" />
      <span data-pv2-line="y" className="absolute inset-y-0 left-2/3 hidden w-px bg-border md:block" />
      <span data-pv2-line="x" className="absolute inset-x-0 top-1/2 hidden h-px bg-border md:block" />

      <CornerMark className="top-0 left-0" />
      <CornerMark className="top-0 left-full" />
      <CornerMark className="top-full left-0" />
      <CornerMark className="top-full left-full" />

      <ul className="relative grid grid-cols-2 md:grid-cols-3">
        {cells.map((rotation) => (
          <PartnerCell key={rotation[0].name} rotation={rotation} />
        ))}
      </ul>
    </div>
  )
}

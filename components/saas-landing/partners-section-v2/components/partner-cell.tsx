import { cn } from "@/lib/utils"

import PartnerLogo from "./partner-logo"
import type { Partner } from "../lib/partners"

export default function PartnerCell({ rotation }: { rotation: Partner[] }) {
  return (
    <li data-pv2-cell className="group relative h-28 overflow-hidden sm:h-32 md:h-40">
      {rotation.map((partner, j) => (
        <div
          key={partner.name}
          data-pv2-layer
          className={cn(
            "absolute inset-0 flex items-center justify-center",
            j > 0 && "invisible opacity-0"
          )}
        >
          <span data-pv2-logo>
            <span className="block transition-transform duration-500 ease-out group-hover:scale-[1.06]">
              <PartnerLogo partner={partner} />
            </span>
          </span>
        </div>
      ))}
    </li>
  )
}

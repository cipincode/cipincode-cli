import { cn } from "@/lib/utils"

import type { Partner } from "../lib/partners"

const partnerLogoShellClass =
  "inline-flex h-7 w-32 max-w-full items-center justify-center md:h-8 md:w-36"

export default function PartnerLogo({
  partner,
  className,
}: {
  partner: Partner
  className?: string
}) {
  const Light = partner.light
  const Dark = partner.dark
  const shellClass = cn(partnerLogoShellClass, className)

  if (Dark) {
    return (
      <span className={shellClass}>
        <Light className="block h-full w-full dark:hidden" aria-hidden />
        <Dark className="hidden h-full w-full dark:block" aria-hidden />
      </span>
    )
  }

  return (
    <span className={shellClass}>
      <Light className="h-full w-full" aria-hidden />
    </span>
  )
}

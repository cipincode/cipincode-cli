import { cn } from "@/lib/utils"

import type { Brand } from "./brands"

const brandLogoShellClass =
  "inline-flex h-7 w-32 max-w-full items-center justify-center md:h-8 md:w-36"

export default function BrandLogo({
  brand,
  className,
}: {
  brand: Brand
  className?: string
}) {
  const Light = brand.light
  const Dark = brand.dark
  const shellClass = cn(brandLogoShellClass, className)

  if (Dark) {
    return (
      <span className={shellClass}>
        <Light
          className="block h-full w-full dark:hidden"
          aria-label={brand.name}
        />
        <Dark
          className="hidden h-full w-full dark:block"
          aria-label={brand.name}
        />
      </span>
    )
  }

  return (
    <span className={shellClass}>
      <Light className="h-full w-full" aria-label={brand.name} />
    </span>
  )
}

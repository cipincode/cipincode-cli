import type { ReactNode } from "react"

import { cn } from "@/lib/utils"

type MockScreenProps = {
  children: ReactNode
  className?: string
  label: string
}

export function MockScreen({ children, className, label }: MockScreenProps) {
  return (
    <div
      className={cn(
        "fsv3-mock flex size-full min-h-0 flex-col overflow-hidden rounded-xl border border-border/70 bg-card text-foreground opacity-70 shadow-sm motion-safe:transition-opacity duration-200 ease-out group-hover/card:opacity-100",
        className
      )}
      role="img"
      aria-label={label}
    >
      {children}
    </div>
  )
}

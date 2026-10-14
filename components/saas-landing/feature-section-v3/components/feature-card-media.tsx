import type { ReactNode } from "react"

import { cn } from "@/lib/utils"

import { FEATURE_CARD_MEDIA_HEIGHT } from "./feature-illustration"

type FeatureCardMediaProps = {
  children: ReactNode
  className?: string
}

export function FeatureCardMedia({ children, className }: FeatureCardMediaProps) {
  return (
    <div
      className={cn(
        "flex w-full shrink-0 items-center justify-center border-b border-border/60 bg-muted/20 p-4",
        FEATURE_CARD_MEDIA_HEIGHT,
        className
      )}
      aria-hidden
    >
      <div className="h-full max-h-full w-full max-w-full min-h-0 overflow-hidden">
        {children}
      </div>
    </div>
  )
}

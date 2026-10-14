"use client"

import { cn } from "@/lib/utils"

import { MockScreen } from "./mock-screen"
import { useInViewOnce } from "../lib/use-in-view-once"

const STOPS: ReadonlyArray<{
  label: string
  stuck?: boolean
}> = [
  { label: "Pindai" },
  { label: "Uji" },
  { label: "Perbaiki", stuck: true },
  { label: "Rilis" },
]

export default function DropoffIllustration() {
  const { ref, active } = useInViewOnce()

  return (
    <div ref={ref} className="size-full min-h-0 p-1">
      <MockScreen
        label="Agen macet di Perbaiki karena tes terus gagal"
        className={cn("fsv3-do-mock", active && "fsv3-do-mock--active")}
      >
        <div className="flex min-h-0 flex-1 flex-col justify-between px-3 py-3">
          <div className="fsv3-do-path relative flex items-start justify-between gap-1 pt-1">
            <div
              className="absolute top-[0.7rem] right-4 left-4 h-px bg-border"
              aria-hidden
            />
            {STOPS.map((stop, index) => (
              <div
                key={stop.label}
                className={cn(
                  "fsv3-do-stop relative z-10 flex flex-1 flex-col items-center gap-1.5",
                  `fsv3-do-stop--${index + 1}`
                )}
              >
                <span
                  className={cn(
                    "flex size-5 items-center justify-center rounded-full border border-border bg-card",
                    stop.stuck &&
                      "size-6 border-foreground/50 bg-foreground text-background motion-safe:transition-transform duration-200 ease-out group-hover/card:scale-110"
                  )}
                >
                  {stop.stuck ? (
                    <span className="font-mono text-[0.5rem] font-semibold">!</span>
                  ) : (
                    <span className="size-1.5 rounded-full bg-foreground/35" />
                  )}
                </span>
                <span
                  className={cn(
                    "font-mono text-[0.5rem] text-muted-foreground",
                    stop.stuck && "font-medium text-foreground/85"
                  )}
                >
                  {stop.label}
                </span>
              </div>
            ))}
          </div>

          <div
            className={cn(
              "fsv3-do-pin relative mx-auto w-full max-w-[11.5rem] rounded-lg border border-foreground/30 bg-muted/50 px-2.5 py-2",
              "motion-safe:transition-transform duration-200 ease-out group-hover/card:-translate-y-0.5"
            )}
          >
            <span
              className="absolute -top-1.5 left-1/2 size-2.5 -translate-x-1/2 rotate-45 border-t border-l border-foreground/30 bg-muted/50"
              aria-hidden
            />
            <p className="font-mono text-[0.5rem] tracking-[0.1em] text-muted-foreground uppercase">
              Kenapa macet
            </p>
            <p className="mt-0.5 text-[0.6875rem] leading-snug font-medium tracking-tight text-foreground">
              Tes gagal menghalangi perbaikan
            </p>
            <p className="mt-1 font-mono text-[0.5625rem] text-foreground/70">
              3 dari 5 run
            </p>
          </div>
        </div>
      </MockScreen>
    </div>
  )
}

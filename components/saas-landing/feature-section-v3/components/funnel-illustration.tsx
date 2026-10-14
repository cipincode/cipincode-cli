"use client"

import { cn } from "@/lib/utils"

import { MockScreen } from "./mock-screen"
import { useInViewOnce } from "../lib/use-in-view-once"

const WEDGES: ReadonlyArray<{
  label: string
  users: string
  width: string
  gap?: boolean
}> = [
  { label: "Pindai", users: "48 file", width: "92%" },
  { label: "Rencana", users: "12 tugas", width: "74%" },
  { label: "Eksekusi", users: "5 perubahan", width: "52%", gap: true },
  { label: "Rilis", users: "3 rilis", width: "34%" },
]

export default function FunnelIllustration() {
  const { ref, active } = useInViewOnce()

  return (
    <div ref={ref} className="size-full min-h-0 p-1">
      <MockScreen
        label="Tahapan jalan agen dengan kemacetan di Eksekusi"
        className={cn("fsv3-fn-mock", active && "fsv3-fn-mock--active")}
      >
        <div className="flex min-h-0 flex-1 flex-col items-center justify-center gap-1.5 px-3 py-3">
          {WEDGES.map((wedge, index) => (
            <div
              key={wedge.label}
              className={cn(
                "fsv3-fn-wedge relative flex w-full flex-col items-center",
                `fsv3-fn-wedge--${index + 1}`,
                wedge.gap &&
                  "origin-center motion-safe:transition-transform duration-200 ease-out group-hover/card:scale-[1.03]"
              )}
            >
              <div
                className={cn(
                  "flex h-7 items-center justify-between gap-2 rounded-sm border border-border/65 bg-muted/35 px-2",
                  wedge.gap && "border-foreground/40 bg-muted/65"
                )}
                style={{ width: wedge.width }}
              >
                <span className="font-mono text-[0.5625rem] text-foreground/80">
                  {wedge.label}
                </span>
                <span className="font-mono text-[0.5625rem] text-muted-foreground tabular-nums">
                  {wedge.users}
                </span>
              </div>
              {wedge.gap ? (
                <span className="fsv3-fn-gap-mark mt-0.5 font-mono text-[0.5rem] tracking-[0.12em] text-foreground/70 uppercase">
                  macet · −56%
                </span>
              ) : null}
            </div>
          ))}
        </div>
      </MockScreen>
    </div>
  )
}

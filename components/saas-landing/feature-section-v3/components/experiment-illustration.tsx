"use client"

import { cn } from "@/lib/utils"

import { MockScreen } from "./mock-screen"
import { useInViewOnce } from "../lib/use-in-view-once"

function CheckoutDraft({
  fields,
  winner,
}: {
  fields: number
  winner?: boolean
}) {
  return (
    <div
      className={cn(
        "flex min-h-0 flex-1 flex-col overflow-hidden rounded-md border",
        winner
          ? "border-foreground/45 bg-card"
          : "border-border/50 bg-muted/15 opacity-60"
      )}
    >
      <div className="flex h-3 shrink-0 items-center gap-0.5 border-b border-border/45 px-1.5">
        <span className="size-1 rounded-full bg-foreground/25" />
        <span className="size-1 rounded-full bg-foreground/25" />
        <span className="size-1 rounded-full bg-foreground/25" />
      </div>
      <div className="flex min-h-0 flex-1 flex-col justify-end gap-1 p-1.5">
        {Array.from({ length: fields }, (_, i) => (
          <div
            key={i}
            className="h-1.5 w-full rounded-[2px] border border-border/50 bg-muted/40"
          />
        ))}
        <div
          className={cn(
            "mt-0.5 flex h-3.5 items-center justify-center rounded-[2px]",
            winner ? "bg-foreground" : "bg-foreground/20"
          )}
        >
          <span
            className={cn(
              "font-mono text-[0.4375rem] tracking-wide uppercase",
              winner ? "text-background" : "text-foreground/40"
            )}
          >
            Kirim
          </span>
        </div>
      </div>
    </div>
  )
}

export default function ExperimentIllustration() {
  const { ref, active } = useInViewOnce()

  return (
    <div ref={ref} className="size-full min-h-0 p-1">
      <MockScreen
        label="Run terpandu yang baru mengalahkan run otomatis yang lama"
        className={cn("fsv3-ex-mock", active && "fsv3-ex-mock--active")}
      >
        <div className="flex min-h-0 flex-1 flex-col gap-2 px-2.5 py-2.5">
          <div className="grid min-h-0 flex-1 grid-cols-[1fr_auto_1fr] items-stretch gap-1.5">
            <div className="fsv3-ex-side fsv3-ex-side--a flex min-h-0 flex-col gap-1">
              <p className="shrink-0 font-mono text-[0.5rem] tracking-[0.1em] text-muted-foreground uppercase">
                A · Otomatis
              </p>
              <CheckoutDraft fields={4} />
            </div>

            <div
              className="fsv3-ex-vs flex shrink-0 items-center self-center font-mono text-[0.5rem] text-muted-foreground"
              aria-hidden
            >
              vs
            </div>

            <div className="fsv3-ex-side fsv3-ex-side--b flex min-h-0 flex-col gap-1">
              <div className="flex shrink-0 items-center justify-between gap-1">
                <p className="font-mono text-[0.5rem] tracking-[0.1em] text-foreground/80 uppercase">
                  B · Terpandu
                </p>
                <span className="rounded-[2px] bg-foreground px-1 py-px font-mono text-[0.4375rem] font-semibold tracking-wide text-background uppercase">
                  Menang
                </span>
              </div>
              <CheckoutDraft fields={1} winner />
            </div>
          </div>

          <div
            className={cn(
              "fsv3-ex-result flex shrink-0 items-baseline justify-between gap-2 rounded-md border border-foreground/30 bg-muted/55 px-2.5 py-2",
              "motion-safe:transition-transform duration-200 ease-out group-hover/card:-translate-y-0.5"
            )}
          >
            <div>
              <p className="font-mono text-[0.5rem] tracking-[0.1em] text-muted-foreground uppercase">
                Hasil
              </p>
              <p className="mt-0.5 text-[0.75rem] leading-none font-medium tracking-tight text-foreground">
                +32% lebih sedikit percobaan ulang
              </p>
            </div>
            <p className="font-mono text-[0.5rem] tracking-[0.08em] text-foreground/70 uppercase">
              Rilis
            </p>
          </div>
        </div>
      </MockScreen>
    </div>
  )
}

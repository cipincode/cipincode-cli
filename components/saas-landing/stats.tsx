import { useEffect, useRef, useState } from "react"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"
import { Globe, Users, Star, Clock, ArrowUp } from "lucide-react"

type IconProps = { className?: string; size?: number | string }

type IconRenderer = (props: IconProps) => React.ReactNode

type Metric = {
  id: string
  icon: IconRenderer
  target: number
  decimals: number
  prefix: string
  suffix: string
  label: string
  delta: string
}

const metrics: Metric[] = [
  {
    id: "uptime",
    icon: (p: IconProps) => (
      <Globe {...p} />
    ),
    target: 99.9,
    decimals: 1,
    prefix: "",
    suffix: "%",
    label: "Service uptime",
    delta: "+0.4%",
  },
  {
    id: "users",
    icon: (p: IconProps) => (
      <Users {...p} />
    ),
    target: 12,
    decimals: 0,
    prefix: "",
    suffix: "k+",
    label: "Active teams",
    delta: "+18%",
  },
  {
    id: "rating",
    icon: (p: IconProps) => (
      <Star {...p} />
    ),
    target: 4.9,
    decimals: 1,
    prefix: "",
    suffix: "/5",
    label: "Average rating",
    delta: "+0.2",
  },
  {
    id: "integrations",
    icon: (p: IconProps) => (
      <Clock {...p} />
    ),
    target: 250,
    decimals: 0,
    prefix: "",
    suffix: "+",
    label: "Integrations shipped",
    delta: "+32",
  },
]

function formatValue(value: number, metric: Metric) {
  const fixed = value.toFixed(metric.decimals)
  return `${metric.prefix}${fixed}${metric.suffix}`
}

export default function Stats() {
  const sectionRef = useRef<HTMLElement | null>(null)
  const [progress, setProgress] = useState(1)

  useEffect(() => {
    const node = sectionRef.current
    if (!node) return

    const prefersReduced = window.matchMedia(
      "(prefers-reduced-motion: reduce)"
    ).matches
    if (prefersReduced) return

    let frame = 0
    let started = false

    const animate = (start: number) => {
      const duration = 1600

      const step = (now: number) => {
        const elapsed = now - start
        const linear = Math.min(elapsed / duration, 1)
        const eased = 1 - Math.pow(1 - linear, 3)
        setProgress(eased)
        if (linear < 1) {
          frame = requestAnimationFrame(step)
        }
      }

      frame = requestAnimationFrame(step)
    }

    const box = node.getBoundingClientRect()
    const alreadyVisible = box.top < window.innerHeight && box.bottom > 0
    if (alreadyVisible) return

    setProgress(0)

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting && !started) {
            started = true
            animate(performance.now())
            observer.disconnect()
          }
        }
      },
      { threshold: 0.35 }
    )

    observer.observe(node)

    return () => {
      observer.disconnect()
      if (frame) cancelAnimationFrame(frame)
    }
  }, [])

  return (
    <section
      ref={sectionRef}
      className="flex w-full items-center justify-center px-6 py-16 sm:py-24"
    >
      <div className="mx-auto w-full max-w-6xl">
        <div className="mx-auto max-w-2xl text-center">
          <Badge variant="outline" className="mb-4 tracking-widest uppercase">
            By The Numbers
          </Badge>
          <h2 className="font-heading text-3xl font-bold tracking-tight text-balance sm:text-4xl">
            Numbers teams at Acme trust
          </h2>
          <p className="mt-4 text-base text-balance text-muted-foreground">
            Real performance from a platform built to scale with you, measured
            and reported every single quarter.
          </p>
        </div>

        <ul className="mt-12 grid grid-cols-1 gap-px overflow-hidden rounded-lg border border-border bg-border sm:grid-cols-2 md:grid-cols-4">
          {metrics.map((metric) => {
            const Icon = metric.icon
            const current = metric.target * progress
            return (
              <li
                key={metric.id}
                className="group flex flex-col gap-4 bg-card p-6 transition-colors hover:bg-card/60 sm:p-8"
              >
                <div className="flex items-center justify-between">
                  <span className="flex size-10 items-center justify-center rounded-md bg-muted text-muted-foreground transition-colors group-hover:bg-primary/10 group-hover:text-primary">
                    <Icon className="size-5" aria-hidden="true" />
                  </span>
                  <Badge
                    variant="outline"
                    className="border-transparent bg-primary/10 text-primary"
                  >
                    <ArrowUp data-icon="inline-start" />
                    {metric.delta}
                  </Badge>
                </div>
                <div className="mt-auto">
                  <p
                    className={cn(
                      "text-4xl font-semibold tracking-tight tabular-nums sm:text-5xl"
                    )}
                  >
                    {formatValue(current, metric)}
                  </p>
                  <p className="mt-2 text-sm text-muted-foreground">
                    {metric.label}
                  </p>
                </div>
              </li>
            )
          })}
        </ul>
      </div>
    </section>
  )
}

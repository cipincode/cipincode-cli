import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { Users, BarChart, ShieldCheck, Check, ArrowRight } from "lucide-react"

type IconProps = { className?: string; size?: number | string }

const rows = [
  {
    eyebrow: "Workflow",
    Icon: (p: IconProps) => (
      <Users {...p} />
    ),
    title: "Interactive TUI & Slash Commands",
    body: "A first-class command palette in your terminal. Run /skills, /missions, /mcp, and dozens more without leaving the prompt.",
    bullets: [
      "Bubbletea TUI engine with modal pop-ups",
      "Real-time agent status and step feedback",
      "Persistent session history and fork support",
      "Custom themes and keybinding configurations",
    ],
    cta: "Explore TUI Engine",
    img: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&q=80",
    imgAlt: "Terminal coding interface showing code and syntax",
    avatars: [
      {
        name: "Sarah Kim",
        initials: "SK",
        src: "https://i.pravatar.cc/150?img=47",
      },
      {
        name: "Marcus Webb",
        initials: "MW",
        src: "https://i.pravatar.cc/150?img=12",
      },
      {
        name: "Priya Nair",
        initials: "PN",
        src: "https://i.pravatar.cc/150?img=32",
      },
    ],
    stat: { value: "10x", label: "Faster Coding Workflow" },
  },
  {
    eyebrow: "AI Execution",
    Icon: (p: IconProps) => (
      <BarChart {...p} />
    ),
    title: "Multi-Mode Autonomous Agents",
    body: "Switch instantly between Chat for Q&A, Coding with auto-approve, and Debugging for deep root-cause analysis.",
    bullets: [
      "Adjustable effort levels (low, medium, high)",
      "Autonomous patch generation and verification",
      "Multi-model routing (Grok, Claude, GPT)",
      "Automated test runner and error reporting",
    ],
    cta: "See Multi-Mode in Action",
    img: "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=800&q=80",
    imgAlt: "Matrix code visualization representing AI execution",
    avatars: [
      {
        name: "James Okafor",
        initials: "JO",
        src: "https://i.pravatar.cc/150?img=68",
      },
      {
        name: "Lena Strauss",
        initials: "LS",
        src: "https://i.pravatar.cc/150?img=5",
      },
    ],
    stat: { value: "< 1s", label: "Agent Response Latency" },
  },
  {
    eyebrow: "Security",
    Icon: (p: IconProps) => (
      <ShieldCheck {...p} />
    ),
    title: "Encrypted Credentials & Token Auth",
    body: "Securely manage API keys, monitor 5-hour and weekly message quotas, and revoke sessions instantly from the web dashboard.",
    bullets: [
      "OAuth cryptographic token verification",
      "Real-time usage metering and alerts",
      "Instant API key revocation from web",
      "Local credential encryption in ~/.config",
    ],
    cta: "Review Security Docs",
    img: "https://images.unsplash.com/photo-1563986768609-322da13575f3?w=800&q=80",
    imgAlt: "Cyber security lock and token representation",
    avatars: [
      {
        name: "Diana Reyes",
        initials: "DR",
        src: "https://i.pravatar.cc/150?img=9",
      },
      {
        name: "Tom Eriksen",
        initials: "TE",
        src: "https://i.pravatar.cc/150?img=53",
      },
      {
        name: "Aiko Tanaka",
        initials: "AT",
        src: "https://i.pravatar.cc/150?img=25",
      },
    ],
    stat: { value: "100%", label: "Local Control & Privacy" },
  },
]

export default function Features() {
  return (
    <section className="flex w-full items-center justify-center px-6 py-16 sm:py-24">
      <div className="mx-auto w-full max-w-6xl">
        <div className="mx-auto mb-16 max-w-xl text-center">
          <Badge variant="outline" className="mb-4 tracking-widest uppercase">
            Platform
          </Badge>
          <h2 className="font-heading text-3xl font-bold tracking-tight text-balance sm:text-4xl">
            Built for High-Performance Development
          </h2>
          <p className="mt-4 text-base leading-relaxed text-muted-foreground">
            CipinCode CLI brings TUI speed, multi-mode AI execution, and secure tracking into one cohesive workflow.
          </p>
        </div>

        <div className="flex flex-col">
          {rows.map((row, index) => {
            const isEven = index % 2 === 0
            return (
              <div key={row.eyebrow}>
                <div
                  className={`flex flex-col gap-10 py-16 md:flex-row md:items-center md:gap-16 ${
                    isEven ? "" : "md:flex-row-reverse"
                  }`}
                >
                  <div className="flex flex-1 flex-col gap-6">
                    <div className="flex items-center gap-2">
                      <span className="flex size-6 shrink-0 items-center justify-center rounded-md border border-border bg-muted">
                        <row.Icon
                          className="size-3.5 text-muted-foreground"
                          aria-hidden="true"
                        />
                      </span>
                      <span className="text-xs font-semibold tracking-widest text-muted-foreground uppercase">
                        {row.eyebrow}
                      </span>
                    </div>

                    <h3 className="font-heading text-2xl font-bold tracking-tight sm:text-[1.75rem] sm:leading-snug">
                      {row.title}
                    </h3>

                    <p className="leading-relaxed text-muted-foreground">
                      {row.body}
                    </p>

                    <ul className="flex flex-col gap-3">
                      {row.bullets.map((bullet) => (
                        <li
                          key={bullet}
                          className="flex items-center gap-3 text-sm"
                        >
                          <span className="flex size-[18px] shrink-0 items-center justify-center rounded-md bg-primary text-primary-foreground">
                            <Check className="size-2.5" aria-hidden="true" />
                          </span>
                          <span className="text-foreground">{bullet}</span>
                        </li>
                      ))}
                    </ul>

                    <div className="flex items-center gap-4 border-t border-border pt-5">
                      <div className="flex -space-x-2">
                        {row.avatars.map((av) => (
                          <Avatar
                            key={av.name}
                            className="size-7 border-2 border-background"
                          >
                            <AvatarImage
                              src={av.src}
                              alt={av.name}
                              className="grayscale"
                            />
                            <AvatarFallback className="text-[10px]">
                              {av.initials}
                            </AvatarFallback>
                          </Avatar>
                        ))}
                      </div>
                      <div className="flex items-baseline gap-1.5">
                        <span className="text-lg font-bold tabular-nums">
                          {row.stat.value}
                        </span>
                        <span className="text-xs text-muted-foreground">
                          {row.stat.label}
                        </span>
                      </div>
                    </div>

                    <div>
                      <Button variant="outline" size="sm">
                        {row.cta}
                        <ArrowRight data-icon="inline-end" />
                      </Button>
                    </div>
                  </div>

                  <div className="flex flex-1 items-center justify-center">
                    <div className="relative aspect-[4/3] w-full overflow-hidden rounded-lg border border-border bg-muted">
                      <img
                        src={row.img}
                        alt={row.imgAlt}
                        className="size-full object-cover"
                      />
                      <div className="absolute inset-x-0 bottom-0 flex items-center gap-2 border-t border-border bg-background/95 px-4 py-2.5 backdrop-blur-sm">
                        <span className="flex size-5 shrink-0 items-center justify-center rounded-md border border-border bg-background">
                          <row.Icon
                            className="size-3 text-foreground"
                            aria-hidden="true"
                          />
                        </span>
                        <span className="text-[11px] font-medium tracking-wide text-foreground">
                          {row.eyebrow} Preview
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                {index < rows.length - 1 && <Separator />}
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}

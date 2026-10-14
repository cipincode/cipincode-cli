import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Separator } from "@/components/ui/separator"
import { ArrowRight, PlayCircle, Play, Check } from "lucide-react"

const TRUST_ITEMS = [
  "No credit card required",
  "14-day free trial",
  "Cancel anytime",
]

const SCREENSHOT_ROWS = [
  {
    label: "Revenue",
    value: "$48,290",
    change: "+12.4%",
  },
  {
    label: "Users",
    value: "3,841",
    change: "+8.1%",
  },
  {
    label: "Churn",
    value: "1.7%",
    change: "-0.3%",
  },
  {
    label: "Session",
    value: "4m 22s",
    change: "+0.9%",
  },
]

const BAR_HEIGHTS = [30, 55, 42, 70, 48, 80, 65, 90, 58, 75, 62, 88]

const TOP_CHANNELS = [
  { name: "Organic search", pct: 74, sessions: "2,841" },
  { name: "Direct", pct: 53, sessions: "2,032" },
  { name: "Referral", pct: 31, sessions: "1,190" },
]

export default function Hero() {
  return (
    <section className="flex w-full items-center justify-center px-6 py-16 sm:py-24">
      <div className="mx-auto grid w-full max-w-6xl grid-cols-1 gap-12 md:grid-cols-2 md:items-center md:gap-16">
        <div className="flex flex-col">
          <Badge variant="outline" className="w-fit">
            Next-Gen AI Terminal Agent
          </Badge>

          <h1 className="mt-6 font-heading text-4xl font-bold tracking-tight text-balance sm:text-5xl">
            Interactive AI Agents
            <br className="hidden sm:block" /> in Your Terminal.
          </h1>

          <p className="mt-5 text-lg text-muted-foreground">
            The CipinCode Terminal UI gives you slash commands, persistent sessions, and full agent support for delegating tasks without leaving your keyboard.
          </p>

          <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
            <Button
              render={<a href="/login" />}
              nativeButton={false}
              className="w-full sm:w-auto"
            >
              Get Started
              <ArrowRight data-icon="inline-end" aria-hidden="true" />
            </Button>
            <Dialog>
              <DialogTrigger
                render={<Button variant="ghost" className="w-full sm:w-auto" />}
              >
                <PlayCircle data-icon="inline-start" aria-hidden="true" />
                Watch Demo
              </DialogTrigger>
              <DialogContent className="sm:max-w-2xl">
                <DialogHeader>
                  <DialogTitle>Product demo</DialogTitle>
                  <DialogDescription>
                    A 2-minute tour of the Acme analytics dashboard.
                  </DialogDescription>
                </DialogHeader>
                <div className="relative flex aspect-video w-full items-center justify-center rounded-lg border border-border bg-muted">
                  <span
                    className="flex size-14 items-center justify-center rounded-full bg-background/80 ring-1 ring-border"
                    aria-hidden="true"
                  >
                    <Play className="size-6 text-foreground" />
                  </span>
                  <span className="sr-only">Demo video placeholder</span>
                </div>
              </DialogContent>
            </Dialog>
          </div>

          <Separator className="my-8" />

          <ul className="flex flex-col gap-2.5">
            {TRUST_ITEMS.map((item) => (
              <li
                key={item}
                className="flex items-center gap-2 text-sm text-muted-foreground"
              >
                <Check className="size-4 shrink-0 text-foreground" aria-hidden="true" />
                {item}
              </li>
            ))}
          </ul>
        </div>

        <div className="relative">
          <div className="rounded-xl border border-border bg-card p-1">
            <div className="flex items-center gap-1.5 rounded-t-lg border-b border-border bg-muted px-3 py-2">
              <span className="size-2.5 rounded-full bg-red-500" />
              <span className="size-2.5 rounded-full bg-yellow-500" />
              <span className="size-2.5 rounded-full bg-green-500" />
              <span className="ml-3 font-mono text-xs text-muted-foreground">terminal — zsh</span>
            </div>

            <div className="flex flex-col rounded-b-lg bg-[#0d0d0f] p-5 font-mono text-sm text-zinc-200">
              <div className="flex items-center justify-between mb-4 pb-3 border-b border-zinc-800">
                <span className="text-xs text-orange-400 font-bold">cipincode agent session v1.2</span>
                <span className="text-xs text-emerald-400">Connected 🟢</span>
              </div>
              <div className="space-y-3 text-xs leading-relaxed">
                <p className="text-zinc-400">&gt; cipincode --mode=coding --effort=high</p>
                <p className="text-orange-400">⚡ Initializing autonomous agent environment...</p>
                <p className="text-zinc-300">🔍 Scanning project structure... 14 files indexed.</p>
                <p className="text-emerald-400">✔ Found bug in authentication middleware (token expiry check).</p>
                <p className="text-zinc-300">✍ Applying patch to app/Http/Middleware/VerifyToken.php...</p>
                <p className="text-emerald-400 font-bold">🎉 Successfully fixed! All tests passed.</p>
              </div>

              <div className="mt-6 pt-4 border-t border-zinc-800 flex items-center justify-between text-xs text-zinc-400">
                <span>curl -fsSL https://cli.cipincode.com/install.sh | bash</span>
                <button
                  onClick={() => navigator.clipboard.writeText("curl -fsSL https://cli.cipincode.com/install.sh | bash")}
                  className="px-2.5 py-1 bg-zinc-800 hover:bg-zinc-700 text-zinc-200 rounded text-xs font-sans font-semibold transition-colors"
                >
                  Copy
                </button>
              </div>
            </div>
          </div>

          <div
            className="absolute -right-3 -bottom-3 -z-10 size-full rounded-xl border border-border bg-muted"
            aria-hidden="true"
          />
        </div>
      </div>
    </section>
  )
}

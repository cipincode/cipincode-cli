"use client"

import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { ArrowUpRightIcon, Folder } from "lucide-react"

export default function HeroSectionV1() {
  return (
    <section className="py-12 md:py-20 lg:py-24">
      <div className="mx-auto max-w-6xl px-4">
        <div className="grid gap-10 lg:grid-cols-2 lg:items-center lg:gap-12">
          <div className="flex flex-col gap-6">
            <h1 className="text-4xl font-extrabold tracking-tight sm:text-5xl lg:text-6xl text-foreground">
              Agents built <span className="text-primary">for coding</span>
            </h1>
            <p className="text-lg leading-relaxed text-muted-foreground sm:text-xl">
              Build and deploy AI coding agents with a single SDK. Get composable tools, intelligent model routing, and deep codebase context that works across any repository.
            </p>
            <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap">
              <Button asChild variant="default" size="lg">
                <a href="/login">
                  Start building
                  <ArrowUpRightIcon className="ml-2 size-4" />
                </a>
              </Button>
              <Button asChild variant="outline" size="lg">
                <a href="/login">
                  View docs
                  <ArrowUpRightIcon className="ml-2 size-4" />
                </a>
              </Button>
            </div>
          </div>

          <div className="relative mx-auto w-full max-w-lg rounded-xl border border-border bg-card shadow-2xl overflow-hidden font-mono text-xs">
            <div className="flex h-9 items-center border-b border-border bg-muted/50 px-4 gap-2">
              <div className="flex gap-1.5">
                <div className="size-3 rounded-full bg-red-500" />
                <div className="size-3 rounded-full bg-yellow-500" />
                <div className="size-3 rounded-full bg-green-500" />
              </div>
              <div className="flex flex-1 items-center justify-center gap-1.5 text-muted-foreground">
                <Folder className="size-3.5 text-blue-500" />
                <span>projects — -zsh</span>
              </div>
            </div>
            <div className="p-4 space-y-2 text-foreground/90 bg-card/50">
              <div className="text-muted-foreground">$ npx @7ovr/saas-landing init</div>
              <div className="text-green-500">✔ Initializing CipinCode CLI environment...</div>
              <div className="text-green-500">✔ Loading intelligent agent modules...</div>
              <div className="text-cyan-500">→ Ready to build, test & deploy across repositories.</div>
              <div className="flex items-center gap-1 pt-2">
                <span className="text-primary">➜</span>
                <span className="text-muted-foreground">~/cipincode</span>
                <span className="animate-pulse">_</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

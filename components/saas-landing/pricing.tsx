"use client"

import { useState } from "react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { cn } from "@/lib/utils"
import { Sparkles, Check } from "lucide-react"
const loginUrl = () => {
  return '/login'
}

const tiers = [
  {
    name: "Free Plan",
    monthly: 0,
    annual: 0,
    description: "Cocok untuk mencoba dan eksplorasi fitur dasar CLI.",
    features: ["Akses penuh TUI CLI dasar", "Kuota pesan standar 5 jam", "Manajemen 2 API key aktif"],
    cta: "Daftar Gratis",
    featured: false,
  },
  {
    name: "Go Plan (Pro)",
    monthly: 99,
    annual: 79,
    description: "Untuk developer profesional dengan produktivitas tinggi.",
    features: [
      "Kuota pesan tinggi & prioritas",
      "Multi-mode AI & Effort levels",
      "Agent kustom tanpa batas",
      "API keys & session analytics lengkap",
    ],
    cta: "Mulai Sekarang",
    featured: true,
  },
]

export default function Pricing() {
  const [annual, setAnnual] = useState(false)

  return (
    <section className="flex w-full items-center justify-center px-6 py-16 sm:py-24">
      <div className="mx-auto w-full max-w-6xl">
        <div className="mx-auto max-w-2xl text-center">
          <Badge variant="outline" className="mb-4 tracking-widest uppercase">
            Harga & Plan
          </Badge>
          <h2 className="font-heading text-3xl font-bold tracking-tight text-balance sm:text-4xl">
            Pilih Paket Sesuai Kebutuhan
          </h2>
          <p className="mt-4 text-muted-foreground">
            Mulai secara gratis dan tingkatkan seiring berkembangnya skala proyek Anda. Hemat hingga 17% dengan penagihan tahunan.
          </p>
        </div>

        <div className="mt-8 flex justify-center">
          <div
            role="group"
            aria-label="Billing period"
            className="inline-flex items-center gap-1 rounded-lg border border-border bg-muted/50 p-1"
          >
            <button
              type="button"
              onClick={() => setAnnual(false)}
              aria-pressed={!annual}
              className={cn(
                "px-4 py-1.5 text-sm font-medium transition-colors",
                !annual
                  ? "bg-background text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              Monthly
            </button>
            <button
              type="button"
              onClick={() => setAnnual(true)}
              aria-pressed={annual}
              className={cn(
                "flex items-center gap-2 px-4 py-1.5 text-sm font-medium transition-colors",
                annual
                  ? "bg-background text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              )}
            >
              Annual
              <Badge className="px-1.5 text-[10px] font-semibold">-17%</Badge>
            </button>
          </div>
        </div>

        <div className="mt-12 grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-3">
          {tiers.map((tier) => {
            const price = annual ? tier.annual : tier.monthly
            return (
              <Card
                key={tier.name}
                className={cn(
                  "flex flex-col",
                  tier.featured && "bg-muted/30 ring-2 ring-primary"
                )}
              >
                <CardHeader>
                  <CardTitle className="text-base font-semibold">
                    {tier.name}
                  </CardTitle>
                  {tier.featured && (
                    <CardAction>
                      <Badge className="gap-1.5">
                        <Sparkles data-icon="inline-start" className="size-3.5" aria-hidden="true" />
                        Most popular
                      </Badge>
                    </CardAction>
                  )}
                  <CardDescription className="text-sm sm:min-h-10">
                    {tier.description}
                  </CardDescription>
                </CardHeader>
                <CardContent className="flex flex-1 flex-col gap-6">
                  <div>
                    <p className="flex items-baseline gap-1">
                      <span className="text-4xl font-bold tracking-tight text-foreground tabular-nums">
                        {price === 0 ? "Rp 0" : `Rp ${price}rb`}
                      </span>
                      <span className="text-sm text-muted-foreground">
                        {price === 0 ? "/selamanya" : "/bulan"}
                      </span>
                    </p>
                    <p className="mt-1 h-5 text-xs text-muted-foreground">
                      {annual && price > 0
                        ? `Rp ${price * 12}rb billed yearly`
                        : " "}
                    </p>
                  </div>
                  <ul className="flex flex-1 flex-col gap-3 text-sm text-foreground">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex items-center gap-2">
                        <Check className="size-4 shrink-0 text-primary" aria-hidden="true" />
                        {feature}
                      </li>
                    ))}
                  </ul>
                </CardContent>
                <CardFooter>
                  <Button
                    nativeButton={false}
                    variant={tier.featured ? "default" : "outline"}
                    size="lg"
                    className="w-full"
                    render={<a href={loginUrl()} />}
                  >
                    {tier.cta}
                  </Button>
                </CardFooter>
              </Card>
            )
          })}
        </div>
      </div>
    </section>
  )
}

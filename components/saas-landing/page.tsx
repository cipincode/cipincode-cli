import { Toaster } from "@/components/ui/sonner"
import Bento from "./bento"
import Cta from "./cta"
import Faqs from "./faqs"
import Features from "./features"
import Footer from "./footer"
import Header from "./header"
import HeroSection from "./hero-section"
import LogoCloud from "./logo-cloud"
import Pricing from "./pricing"
import Stats from "./stats"
import Testimonials from "./testimonials"

const pool = (size: string, color: string, mid: string, out: string) =>
  `radial-gradient(${size}, ${color}, ${mid} 38%, ${out} 78%)`

const wash = (base: string, a: number, b: number) => {
  const at = (v: number) => `rgba(${base},${v})`
  return [
    pool("38% 13% at 84% 9%", at(a), at(a / 2), at(0)),
    pool("42% 13% at 12% 28%", at(b), at(b / 2), at(0)),
    pool("40% 12% at 70% 47%", at(a), at(a / 2), at(0)),
    pool("42% 12% at 14% 66%", at(b), at(b / 2), at(0)),
    pool("40% 13% at 74% 87%", at(a), at(a / 2), at(0)),
  ].join(",")
}

const WASH_LIGHT = wash("0,0,0", 0.04, 0.034)
const WASH_DARK = wash("255,255,255", 0.06, 0.051)

export default function SaasLandingTemplate() {
  return (
    <div className="relative isolate flex min-h-screen w-full flex-col bg-background text-foreground">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[720px] bg-[radial-gradient(65%_60%_at_50%_-8%,rgba(0,0,0,0.06),transparent_72%)] dark:bg-[radial-gradient(65%_60%_at_50%_-8%,rgba(255,255,255,0.13),transparent_72%)]"
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 dark:hidden"
        style={{ backgroundImage: WASH_LIGHT }}
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10 hidden dark:block"
        style={{ backgroundImage: WASH_DARK }}
      />
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-x-0 bottom-0 -z-10 h-[560px] bg-[radial-gradient(55%_50%_at_50%_100%,rgba(0,0,0,0.05),transparent_70%)] dark:bg-[radial-gradient(55%_50%_at_50%_100%,rgba(255,255,255,0.08),transparent_70%)]"
      />
      <Toaster />
      <Header />
      <main className="flex w-full flex-col flex-1">
        <HeroSection />
        <LogoCloud />
        <Features />
        <Bento />
        <Stats />
        <Testimonials />
        <Pricing />
        <Faqs />
        <Cta />
      </main>
      <div className="border-t border-border/50">
        <Footer />
      </div>
    </div>
  )
}

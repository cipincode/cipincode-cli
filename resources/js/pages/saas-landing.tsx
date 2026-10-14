import { Toaster } from "@/components/ui/sonner"
import Header from "../../../components/saas-landing/header"
import Hero from "../../../components/saas-landing/hero"
import Footer from "../../../components/saas-landing/footer"

export default function SaasLandingPage() {
  return (
    <div className="relative isolate flex min-h-screen w-full flex-col bg-background text-foreground">
      <Toaster />
      <Header />
      <main className="flex w-full flex-col flex-1">
        <Hero />
      </main>
      <div className="border-t border-border/50">
        <Footer />
      </div>
    </div>
  )
}

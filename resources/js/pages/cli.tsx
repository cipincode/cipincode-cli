import { Toaster } from "@/components/ui/sonner"
import Header from "../../../components/saas-landing/header"
import HeroSectionV1 from "../../../components/saas-landing/hero-section"
import FeatureSectionV3 from "../../../components/saas-landing/feature-section-v3/components/feature-section-v3"
import PartnersSectionV2 from "../../../components/saas-landing/partners-section-v2/components/partners-section-v2"
import HowItWorksSectionV1 from "../../../components/saas-landing/how-it-works-section-v1/components/how-it-works-section-v1"
import PricingSectionV3 from "../../../components/saas-landing/pricing-section-v3/components/pricing-section-v3"
import FaqSectionV6 from "../../../components/saas-landing/faq-section-v6/components/faq-section-v6"
import Footer from "../../../components/saas-landing/footer"

export default function CliLandingPage() {
  return (
    <div className="relative isolate flex min-h-screen w-full flex-col bg-background text-foreground">
      <Toaster />
      <Header />
      <main className="flex w-full flex-col flex-1">
        <HeroSectionV1 />
        <FeatureSectionV3 />
        <PartnersSectionV2 />
        <HowItWorksSectionV1 />
        <PricingSectionV3 />
        <FaqSectionV6 />
      </main>
      <div className="border-t border-border/50">
        <Footer />
      </div>
    </div>
  )
}

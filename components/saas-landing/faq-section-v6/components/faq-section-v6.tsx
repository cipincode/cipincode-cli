import FaqAccordion from "./faq-accordion"

export default function FaqSectionV6() {
  return (
    <section id="faq" className="scroll-mt-20 py-16 md:py-20 lg:py-24">
      <div className="mx-auto max-w-6xl px-4">
        <div className="grid gap-10 md:grid-cols-[minmax(0,1fr)_minmax(0,1.7fr)] md:gap-16 lg:gap-24">
          <div className="md:sticky md:top-24 md:self-start">
            <h2 className="text-6xl leading-none font-semibold tracking-tight text-foreground sm:text-7xl lg:text-8xl">
              FAQ
            </h2>
            <p className="mt-4 max-w-xs text-base leading-relaxed text-muted-foreground">
              Jawaban cepat soal instalasi, kuota, paket, dan akun Anda.
            </p>
          </div>

          <FaqAccordion />
        </div>
      </div>
    </section>
  )
}

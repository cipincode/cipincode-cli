import FeatureShowcase from "./feature-showcase"
import "../styles/feature-section-v3.css"

export default function FeatureSectionV3() {
  return (
    <section className="py-4 md:py-6 lg:py-8">
      <div className="mx-auto max-w-6xl px-4">
        <h2 className="max-w-3xl text-3xl leading-tight font-medium tracking-tighter sm:text-4xl lg:text-5xl">
          Ketahui yang terjadi di terminal Anda
        </h2>
        <p className="mt-4 max-w-2xl text-3.75 leading-relaxed text-muted-foreground">
          Mengikuti setiap jalannya agen dari scan hingga ship, memahami apa
          yang menghambat pekerjaan, dan menghantarkan perbaikan dari temuan
          hingga dampak.
        </p>

        <div className="mt-12 lg:mt-16">
          <FeatureShowcase />
        </div>
      </div>
    </section>
  )
}

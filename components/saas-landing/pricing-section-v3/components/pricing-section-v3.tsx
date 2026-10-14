"use client"

import PricingCard from "./pricing-card"
import { pricingPlans } from "../lib/pricing-plans"

export default function PricingSectionV3() {
  return (
    <section id="harga" className="scroll-mt-20 py-4 md:py-6 lg:py-8">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto mb-10 max-w-2xl text-center">
          <h2 className="text-4xl leading-[1.15] tracking-tight text-balance md:text-5xl">
            <span className="block text-foreground">Pilih Paket</span>
            <span className="mt-1 block text-muted-foreground">
              Sesuai Kebutuhan Anda
            </span>
          </h2>
          <p className="mx-auto mt-5 max-w-xl text-[0.9375rem] leading-relaxed text-muted-foreground">
            Pilih paket yang sesuai dengan alur kerja Anda. Semua harga dalam
            rupiah per bulan — paket dan kuota sama seperti di dashboard
            billing Anda.
          </p>
        </div>
        <ul className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
          {pricingPlans.map((plan) => (
            <PricingCard key={plan.id} plan={plan} />
          ))}
        </ul>
      </div>
    </section>
  )
}

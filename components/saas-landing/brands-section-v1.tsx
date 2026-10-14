import Marquee from "./brands-marquee"
import { brands } from "./brands"

import "./brands-section-v1.css"

export default function BrandsSectionV1() {
  return (
    <section className="py-4 md:py-6 lg:py-8">
      <div className="mx-auto max-w-6xl px-4">
        <div className="grid items-center gap-8 md:grid-cols-[minmax(0,260px)_minmax(0,1fr)] md:gap-10 lg:gap-14">
          <p className="max-w-[16rem] text-base leading-snug font-medium text-muted-foreground md:text-lg">
            Dipercaya brand AI terdepan:
          </p>
          <Marquee brands={brands} />
        </div>
      </div>
    </section>
  )
}

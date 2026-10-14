import type { Brand } from "./brands"
import BrandLogo from "./brand-logo"

export default function Marquee({ brands }: { brands: Brand[] }) {
  return (
    <div className="relative min-w-0 overflow-hidden">
      <div
        className="pointer-events-none absolute top-0 left-0 z-1 h-full w-16 bg-linear-to-r from-background to-transparent motion-reduce:hidden"
        aria-hidden
      />
      <div
        className="pointer-events-none absolute top-0 right-0 z-1 h-full w-16 bg-linear-to-l from-background to-transparent motion-reduce:hidden"
        aria-hidden
      />

      <div className="bs1-track flex w-max flex-nowrap items-center py-1.4 will-change-transform backface-hidden">
        <div
          className="flex shrink-0 items-center gap-10 px-2 md:gap-15 md:px-3 lg:gap-16"
          role="list"
          aria-label="Brand logos"
        >
          {brands.map((brand) => (
            <div
              key={brand.name}
              role="listitem"
              className="flex shrink-0 items-center justify-center px-4 md:px-6"
            >
              <BrandLogo brand={brand} />
            </div>
          ))}
        </div>
        <div
          className="flex shrink-0 items-center gap-10 px-2 md:gap-15 md:px-3 lg:gap-16"
          aria-hidden
        >
          {brands.map((brand) => (
            <div
              key={`dup-${brand.name}`}
              className="flex shrink-0 items-center justify-center px-4 md:px-6"
            >
              <BrandLogo brand={brand} />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

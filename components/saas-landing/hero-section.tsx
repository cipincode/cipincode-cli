import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { ArrowUpRightIcon } from "lucide-react"

import MacOsTerminal from "@/components/mac-os-terminal/mac-os-terminal"
import "@/components/mac-os-terminal/mac-os-terminal.css"

import BrandsSectionV1 from "./brands-section-v1"
import { HeroV1Image } from "./hero-v1-image"

export default function HeroSectionV1() {
  return (
    <>
      <section className="py-4 md:py-6 lg:py-8">
        <div className="mx-auto max-w-6xl px-4">
          <div className="grid gap-10 lg:grid-cols-2 lg:items-start lg:gap-12">
            <h1
              className={cn(
                "text-4xl font-medium tracking-tight sm:text-5xl lg:text-6xl xl:text-7xl"
              )}
            >
              Agen AI{" "}
              <span
                className="font-normal tracking-normal text-transparent bg-clip-text bg-[radial-gradient(#333_1px,transparent_1px)] dark:bg-[radial-gradient(#fff_1px,transparent_1px)] [background-size:4px_4px] px-1"
              >
                untuk coding
              </span>
            </h1>

            <div className="flex flex-col gap-6">
              <p
                className={cn(
                  "text-lg leading-relaxed text-muted-foreground sm:text-xl"
                )}
              >
                Bangun dan jalankan agen coding AI dengan satu SDK. Dapatkan
                tools yang bisa dikombinasikan, routing model yang cerdas,
                dan konteks codebase yang mendalam di semua repository.
              </p>
              <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap">
                <Button
                  variant="default"
                  size="lg"
                >
                  <a href="/login" className="flex items-center gap-1">
                    Mulai Membangun
                    <ArrowUpRightIcon className="size-4" />
                  </a>
                </Button>
                <Button
                  variant="outline"
                  size="lg"
                >
                  <a href="/login" className="flex items-center gap-1">
                    Lihat Dokumentasi
                    <ArrowUpRightIcon className="size-4" />
                  </a>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <BrandsSectionV1 />

      <section className="px-4 py-4 md:py-6 lg:py-8">
        <div className="relative isolate mx-auto max-w-6xl overflow-hidden rounded-2xl px-4 py-10 md:px-8 md:py-14">
          <HeroV1Image />
          <MacOsTerminal className="relative z-1 mx-auto w-full max-w-2xl min-w-0" />
        </div>
      </section>
    </>
  )
}

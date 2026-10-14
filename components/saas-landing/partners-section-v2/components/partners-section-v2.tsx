"use client"

import { useRef } from "react"

import PartnerGrid from "./partner-grid"
import { usePartnersV2Motion } from "../hooks/use-partners-v2-motion"
import { partnersSectionV2Content } from "../lib/config"
import { partners } from "../lib/partners"

import "../styles/partners-section-v2.css"

export default function PartnersSectionV2() {
  const { heading, description } = partnersSectionV2Content
  const sectionRef = useRef<HTMLElement>(null)

  usePartnersV2Motion(sectionRef)

  const words = heading.split(" ")

  return (
    <section ref={sectionRef} data-pv2-root className="py-16 md:py-24">
      <div className="mx-auto grid max-w-6xl gap-12 px-4 lg:grid-cols-12 lg:items-center lg:gap-16">
        <div className="flex flex-col gap-5 lg:col-span-5">
          <h2 className="text-3xl leading-[1.05] font-medium tracking-[-0.035em] text-balance text-foreground sm:text-4xl lg:text-5xl">
            {words.map((word, i) => (
              <span key={`${word}-${i}`}>
                <span className="-mb-[0.12em] inline-block overflow-hidden pb-[0.12em] align-top">
                  <span data-pv2-word className="inline-block">
                    {word}
                  </span>
                </span>
                {i < words.length - 1 ? " " : null}
              </span>
            ))}
          </h2>
          <p
            data-pv2-description
            className="max-w-md text-base leading-relaxed text-pretty text-muted-foreground md:text-lg"
          >
            {description}
          </p>
        </div>

        <div className="lg:col-span-7">
          <ul aria-label="Ekosistem teknologi" className="sr-only">
            {partners.map((partner) => (
              <li key={partner.name}>{partner.name}</li>
            ))}
          </ul>
          <PartnerGrid />
        </div>
      </div>
    </section>
  )
}

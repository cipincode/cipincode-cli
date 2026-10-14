"use client"

import { PlusIcon } from "lucide-react"
import { Accordion as AccordionPrimitive } from "radix-ui"

import { faqItems } from "../lib/faq-content"

export default function FaqAccordion() {
  return (
    <AccordionPrimitive.Root type="single" collapsible className="w-full">
      {faqItems.map((item) => (
        <AccordionPrimitive.Item
          key={item.id}
          value={item.id}
          className="border-b border-dashed border-border"
        >
          <AccordionPrimitive.Header className="flex">
            <AccordionPrimitive.Trigger className="group/faq flex flex-1 items-start justify-between gap-5 py-6 text-left text-xl font-medium tracking-tight text-foreground outline-none transition-colors hover:text-foreground/70 focus-visible:underline sm:text-2xl">
              {item.question}
              <PlusIcon className="mt-1 size-7 shrink-0 text-muted-foreground transition-transform duration-200 group-data-[state=open]/faq:rotate-45" />
            </AccordionPrimitive.Trigger>
          </AccordionPrimitive.Header>
          <AccordionPrimitive.Content className="overflow-hidden data-[state=open]:animate-accordion-down data-[state=closed]:animate-accordion-up">
            <p className="pb-6 pr-12 font-mono text-base leading-relaxed text-muted-foreground sm:text-lg">
              {item.answer}
            </p>
          </AccordionPrimitive.Content>
        </AccordionPrimitive.Item>
      ))}
    </AccordionPrimitive.Root>
  )
}

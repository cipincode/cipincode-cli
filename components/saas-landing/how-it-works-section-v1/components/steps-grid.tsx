"use client"

import { steps } from "../lib/steps"
import StepCard from "./step-card"

export default function StepsGrid() {
  return (
    <div className="hiw1-grid relative grid grid-cols-1 gap-12 sm:grid-cols-2 lg:grid-cols-3 lg:gap-10">
      {steps.map((step, i) => (
        <div key={step.id} className="hiw1-step-item relative isolate">
          <StepCard step={step} />
          {i < steps.length - 1 ? (
            <>

              <span
                aria-hidden
                className={
                  i === 0
                    ? "hiw1-beam-h pointer-events-none absolute top-1/2 hidden -translate-y-1/2 sm:block"
                    : "hiw1-beam-h pointer-events-none absolute top-1/2 hidden -translate-y-1/2 lg:block"
                }
              />

              <span
                aria-hidden
                className={
                  i === 0
                    ? "hiw1-beam-vertical pointer-events-none absolute -bottom-12 left-1/2 h-12 w-px -translate-x-1/2 lg:hidden"
                    : "hiw1-beam-vertical pointer-events-none absolute -bottom-12 left-1/2 h-12 w-px -translate-x-1/2 sm:hidden"
                }
              />
            </>
          ) : null}
        </div>
      ))}
    </div>
  )
}

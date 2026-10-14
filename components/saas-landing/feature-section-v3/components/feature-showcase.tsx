import type { ComponentType } from "react"

import { Card, CardContent, CardTitle } from "@/components/ui/card"

import DropoffIllustration from "./dropoff-illustration"
import ExperimentIllustration from "./experiment-illustration"
import { FeatureCardMedia } from "./feature-card-media"
import FunnelIllustration from "./funnel-illustration"
import { features, type FeatureIllustrationId } from "../lib/features"

const illustrationComponents: Record<FeatureIllustrationId, ComponentType> = {
  "funnel-opportunities": FunnelIllustration,
  "dropoff-diagnosis": DropoffIllustration,
  "experiment-impact": ExperimentIllustration,
}

export default function FeatureShowcase() {
  return (
    <ul className="m-0 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-3">
      {features.map(({ id, step, title, description, illustration }) => {
        const Illustration = illustrationComponents[illustration]

        return (
          <li key={id} className="flex min-w-0 flex-col">
            <Card className="group/card flex h-full min-h-80 flex-col gap-0 py-0">
              <FeatureCardMedia>
                <Illustration />
              </FeatureCardMedia>
              <CardContent className="flex flex-col gap-0 px-3.5 pt-3 pb-2.25">
                <p
                  className="m-0 mb-1 shrink-0 font-mono text-2.75 leading-[1.2] font-medium tracking-[0.14em] text-foreground/85"
                >
                  {step}
                </p>
                <CardTitle className="m-0 mb-1.25 min-h-[calc(0.95rem*1.35*2)] text-3.8 leading-[1.35] font-semibold tracking-normal text-foreground">
                  {title}
                </CardTitle>
                <p className="m-0 shrink-0 text-3.25 leading-normal text-muted-foreground">
                  {description}
                </p>
              </CardContent>
            </Card>
          </li>
        )
      })}
    </ul>
  )
}

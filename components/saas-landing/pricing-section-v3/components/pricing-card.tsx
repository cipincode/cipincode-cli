import { CheckIcon } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Separator } from "@/components/ui/separator"
import type { PricingPlan } from "../lib/pricing-plans"

type PricingCardProps = {
  plan: PricingPlan
}

export default function PricingCard({ plan }: PricingCardProps) {
  return (
    <li className="min-w-0">
      <Card className="relative h-full gap-0 border border-border py-0 shadow-none ring-0">
        <CardHeader className="gap-3 border-b border-border px-5 pt-5 pb-5">
          <div className="flex items-start justify-between gap-3">
            <div className="flex flex-col gap-1">
              <CardTitle className="text-lg font-semibold tracking-tight text-foreground">
                {plan.name}
              </CardTitle>
              <CardDescription className="text-sm leading-relaxed text-muted-foreground">
                {plan.description}
              </CardDescription>
            </div>
            {plan.popular ? (
              <Badge variant="outline" className="shrink-0">
                Populer
              </Badge>
            ) : null}
          </div>
          <div className="flex flex-col gap-1 pt-1">
            <p className="text-4xl leading-[1.15] tracking-normal text-foreground">
              {plan.price}
            </p>
            <p className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
              {plan.priceSuffix}
            </p>
          </div>
        </CardHeader>
        <CardContent className="flex flex-1 flex-col gap-0 px-5 pt-0 pb-5">
          {plan.featureGroups.map((group, groupIndex) => (
            <div key={group.title}>
              {groupIndex > 0 ? <Separator className="my-4" /> : null}
              <div className={groupIndex === 0 ? "pt-5" : undefined}>
                <p className="mb-3 text-sm text-muted-foreground">{group.title}</p>
                <ul className="flex flex-col gap-2.5">
                  {group.features.map((feature) => (
                    <li
                      key={feature}
                      className="flex items-start gap-2.5 text-sm leading-relaxed text-muted-foreground"
                    >
                      <span className="mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border border-border text-muted-foreground">
                        <CheckIcon aria-hidden className="size-2.5" />
                      </span>
                      <span>{feature}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          ))}
        </CardContent>
        <CardFooter className="border-t-0 bg-transparent p-5 pt-0">
          {plan.id === "free" ? (
            <Button
              asChild
              className="w-full rounded-full bg-orange-600 text-white hover:bg-orange-700 hover:text-white"
              size="lg"
            >
              <a href="/login">{plan.ctaLabel}</a>
            </Button>
          ) : (
            <Button
              className="w-full rounded-full bg-orange-600 text-white hover:bg-orange-600 hover:text-white disabled:opacity-50"
              size="lg"
              disabled
            >
              Segera Hadir
            </Button>
          )}
        </CardFooter>
      </Card>
    </li>
  )
}

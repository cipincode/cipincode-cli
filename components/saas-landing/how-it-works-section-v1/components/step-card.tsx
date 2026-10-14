import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import type { Step } from "../lib/steps"

type Props = {
  step: Step
}

export default function StepCard({ step }: Props) {
  return (
    <Card className="flex h-full flex-col rounded-[2rem] py-0">
      <CardHeader className="p-8 pb-0">
        <CardTitle className="text-base text-muted-foreground">
          {step.number}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col p-8 pt-6 pb-12">
        <CardTitle
          className="text-[2.75rem] leading-[1.05] font-medium tracking-tight text-foreground/90"
        >
          {step.title}
        </CardTitle>
        <CardDescription className="mt-auto line-clamp-4 min-h-24 pt-4 text-[15px] leading-relaxed text-muted-foreground/80">
          {step.description}
        </CardDescription>
      </CardContent>
    </Card>
  )
}

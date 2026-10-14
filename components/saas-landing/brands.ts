import type { ComponentType, SVGProps } from "react"

import { NvidiaWordmarkLight } from "@/components/ui/svgs/nvidiaWordmarkLight"
import { NvidiaWordmarkDark } from "@/components/ui/svgs/nvidiaWordmarkDark"
import { OpenaiWordmarkLight } from "@/components/ui/svgs/openaiWordmarkLight"
import { OpenaiWordmarkDark } from "@/components/ui/svgs/openaiWordmarkDark"
import { GrokWordmarkLight } from "@/components/ui/svgs/grokWordmarkLight"
import { GrokWordmarkDark } from "@/components/ui/svgs/grokWordmarkDark"
import { PerplexityWordmarkLight } from "@/components/ui/svgs/perplexityWordmarkLight"
import { PerplexityWordmarkDark } from "@/components/ui/svgs/perplexityWordmarkDark"
import { ClaudeAiWordmarkIconLight } from "@/components/ui/svgs/claudeAiWordmarkIconLight"
import { ClaudeAiWordmarkIconDark } from "@/components/ui/svgs/claudeAiWordmarkIconDark"
import { GoogleWordmark } from "@/components/ui/svgs/googleWordmark"

export type SvgComponent = ComponentType<SVGProps<SVGSVGElement>>

export type Brand = {
  name: string
  light: SvgComponent
  dark?: SvgComponent
}

export const brands: Brand[] = [
  { name: "Nvidia",        light: NvidiaWordmarkLight,       dark: NvidiaWordmarkDark },
  { name: "OpenAI",        light: OpenaiWordmarkLight,       dark: OpenaiWordmarkDark },
  { name: "Grok",          light: GrokWordmarkLight,         dark: GrokWordmarkDark },
  { name: "Perplexity AI", light: PerplexityWordmarkLight,   dark: PerplexityWordmarkDark },
  { name: "Claude AI",     light: ClaudeAiWordmarkIconLight, dark: ClaudeAiWordmarkIconDark },
  { name: "Google",        light: GoogleWordmark },
]

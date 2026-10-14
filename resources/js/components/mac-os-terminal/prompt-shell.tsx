import type { ReactNode } from "react"
import { PROMPT } from "./config"

export default function PromptShell({ children }: { children?: ReactNode }) {
  return (
    <div className="flex min-w-0 flex-wrap items-center gap-x-1">
      <span className="shrink-0 tabular-nums">{PROMPT}</span>
      {children}
    </div>
  )
}

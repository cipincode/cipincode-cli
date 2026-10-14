import type { ReactNode } from "react"
import { CODE_SYNTAX_PARTS } from "./mac-os-terminal-code-segments"

export default function SyntaxCode({ visibleChars }: { visibleChars: number }) {
  let remaining = visibleChars
  const nodes: ReactNode[] = []
  let key = 0
  for (const { t, c } of CODE_SYNTAX_PARTS) {
    if (remaining <= 0) break
    const take = Math.min(t.length, remaining)
    if (take > 0) {
      nodes.push(
        <span key={key++} className={c}>
          {t.slice(0, take)}
        </span>
      )
      remaining -= take
    }
  }
  return <>{nodes}</>
}

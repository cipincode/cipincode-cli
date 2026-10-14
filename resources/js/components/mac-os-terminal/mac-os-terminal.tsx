"use client"

import { Folder } from "lucide-react"
import { cn } from "@/lib/utils"

import { useTerminalAnimation } from "./use-terminal-animation"
import Spinner from "./spinner"
import PromptShell from "./prompt-shell"
import TrafficLights from "./traffic-lights"
import SyntaxCode from "./syntax-code"
import { CODE_DEMO } from "./mac-os-terminal-code-segments"
import type { LineColor } from "./config"

const codeBlockWrap =
  "mot-codeblock my-2 border-l-2 py-2 pr-2 pl-3 font-mono text-[11px] leading-[1.45]"

function colorCls(c?: LineColor) {
  if (!c) return "text-[var(--mot-t-body-fg)]"
  switch (c) {
    case "green":
      return "text-[var(--mot-t-line-green)]"
    case "cyan":
      return "text-[var(--mot-t-line-cyan)]"
    case "dim":
      return "text-[var(--mot-t-line-dim)]"
  }
}

export default function MacOsTerminal({ className }: { className?: string }) {
  const { lines, typed, phase, codeIdx, scrollRef } = useTerminalAnimation()

  return (
    <div
      className={cn(
        "mot-terminal mx-auto w-full max-w-3xl min-w-0 overflow-hidden rounded-md border border-[var(--mot-t-border)] bg-[var(--mot-t-shell-bg)] font-mono shadow-[var(--mot-t-shadow)] outline-none outline-1 -outline-offset-1 outline-[var(--mot-t-ring)]",
        className
      )}
    >
      <div className="relative flex h-7 shrink-0 items-center border-b border-[var(--mot-t-titlebar-border)] bg-[var(--mot-t-titlebar-bg)] select-none">
        <TrafficLights />
        <div
          className={cn(
            "flex w-full items-center justify-center gap-1.5 px-20 text-center text-[10px] font-medium tracking-tight text-[var(--mot-t-title-fg)]"
          )}
        >
          <Folder
            size={12}
            className="size-3 shrink-0"
            fill="#3b8eea"
            stroke="none"
            strokeWidth={0}
            aria-hidden
          />
          <span>cipincode</span>
          <span className="text-[#a1a1aa] dark:text-[#6b6b6b]" aria-hidden>
            —
          </span>
          <span>-zsh</span>
          <span className="text-[#a1a1aa] dark:text-[#6b6b6b]" aria-hidden>
            —
          </span>
          <span className="tabular-nums">80×24</span>
        </div>
      </div>

      <div
        ref={scrollRef}
        className="mot-terminal-body h-128 overflow-y-auto bg-[var(--mot-t-body-bg)] px-3.5 py-2 text-left font-mono text-[11px] leading-[1.45] text-[var(--mot-t-body-fg)] [&::-webkit-scrollbar]:hidden"
        style={{ scrollbarWidth: "none" }}
      >
        {lines.map((line) => (
          <div
            key={line.id}
            className={cn("mot-terminal-line-in", colorCls(line.color))}
          >
            {line.spinning ? (
              <span>
                <Spinner /> {line.text}
              </span>
            ) : (
              line.text || "\u00A0"
            )}
          </div>
        ))}

        {(phase === "typingCode" || phase === "done") && (
          <div className={codeBlockWrap}>
            <pre className="m-0 font-mono wrap-break-word whitespace-pre-wrap">
              <SyntaxCode
                visibleChars={phase === "done" ? CODE_DEMO.length : codeIdx}
              />
              {phase === "typingCode" && (
                <span className="mot-terminal-cursor ml-px inline-block h-[1em] w-2 translate-y-px bg-[var(--mot-terminal-cursor-bg)] align-text-bottom" />
              )}
            </pre>
          </div>
        )}

        {(phase === "idle" || phase === "typing") && (
          <div className="flex min-w-0 flex-wrap items-center gap-x-1 gap-y-0">
            <PromptShell />
            <span className="min-w-0">{typed}</span>
            <span className="mot-terminal-cursor ml-px inline-block h-[1em] w-2 translate-y-px bg-[var(--mot-terminal-cursor-bg)] align-text-bottom" />
          </div>
        )}

        {phase === "done" && (
          <div className="mt-2 flex items-center gap-x-1">
            <PromptShell />
            <span className="mot-terminal-cursor ml-px inline-block h-[1em] w-2 translate-y-px bg-[var(--mot-terminal-cursor-bg)] align-text-bottom" />
          </div>
        )}
      </div>
    </div>
  )
}

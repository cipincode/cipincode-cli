import { useEffect, useRef, useState } from "react"
import {
  COMMAND, INIT_LINES, OUTPUT_SEQUENCE, PROMPT,
  TYPE_CODE_MS, TYPE_CODE_NEWLINE_MS, TYPE_SPEED,
  type Line, type Phase,
} from "./config"
import { CODE_DEMO } from "./mac-os-terminal-code-segments"

export function useTerminalAnimation() {
  const [lines, setLines] = useState<Line[]>(INIT_LINES)
  const [typed, setTyped] = useState("")
  const [phase, setPhase] = useState<Phase>("idle")
  const [outIdx, setOutIdx] = useState(0)
  const [codeIdx, setCodeIdx] = useState(0)
  const idRef = useRef(10)
  const spinTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [lines, typed, codeIdx, phase])

  useEffect(() => {
    if (phase !== "idle") return
    const t = setTimeout(() => setPhase("typing"), 1600)
    return () => clearTimeout(t)
  }, [phase])

  useEffect(() => {
    if (phase !== "typing") return
    if (typed.length < COMMAND.length) {
      const t = setTimeout(() => setTyped(COMMAND.slice(0, typed.length + 1)), TYPE_SPEED)
      return () => clearTimeout(t)
    }
    const t = setTimeout(() => {
      setLines((prev) => [...prev, { id: idRef.current++, text: `${PROMPT}${COMMAND}` }])
      setTyped("")
      setOutIdx(0)
      setPhase("output")
    }, 400)
    return () => clearTimeout(t)
  }, [phase, typed])

  useEffect(() => {
    if (phase !== "output") return
    if (outIdx >= OUTPUT_SEQUENCE.length) {
      const t = setTimeout(() => setPhase("typingCode"), 350)
      return () => clearTimeout(t)
    }
    const item = OUTPUT_SEQUENCE[outIdx]
    if (item.type === "spinner") {
      const spinId = idRef.current++
      setLines((prev) => [...prev, { id: spinId, text: item.text, spinning: true, color: "dim" }])
      const t = setTimeout(() => {
        spinTimerRef.current = setTimeout(() => {
          setLines((prev) =>
            prev.map((l) =>
              l.id === spinId
                ? { ...l, spinning: false, text: `✓ ${item.doneText}`, color: "green" as const }
                : l,
            ),
          )
          setOutIdx((i) => i + 1)
        }, item.duration)
      }, item.delay)
      return () => {
        clearTimeout(t)
        if (spinTimerRef.current) clearTimeout(spinTimerRef.current)
      }
    }
    const t = setTimeout(() => {
      setLines((prev) => [...prev, { id: idRef.current++, text: item.text, color: item.color }])
      setOutIdx((i) => i + 1)
    }, item.delay)
    return () => clearTimeout(t)
  }, [phase, outIdx])

  useEffect(() => {
    if (phase !== "typingCode") return
    if (codeIdx >= CODE_DEMO.length) {
      const t = setTimeout(() => setPhase("done"), 320)
      return () => clearTimeout(t)
    }
    const ch = CODE_DEMO[codeIdx]
    const delay = ch === "\n" ? TYPE_CODE_NEWLINE_MS : TYPE_CODE_MS
    const t = setTimeout(() => setCodeIdx((i) => i + 1), delay)
    return () => clearTimeout(t)
  }, [phase, codeIdx])

  useEffect(() => {
    if (phase !== "done") return
    const t = setTimeout(() => {
      idRef.current = 10
      setTyped("")
      setOutIdx(0)
      setCodeIdx(0)
      setLines(INIT_LINES)
      setPhase("idle")
    }, 8000)
    return () => clearTimeout(t)
  }, [phase])

  return { lines, typed, phase, codeIdx, scrollRef }
}

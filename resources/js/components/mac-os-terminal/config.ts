export type LineColor = "green" | "cyan" | "dim"

export type Phase = "idle" | "typing" | "output" | "typingCode" | "done"

export interface Line {
  id: number
  text: string
  color?: LineColor
  spinning?: boolean
}

export type OutputItem =
  | { type: "line"; text: string; color?: LineColor; delay: number }
  | { type: "spinner"; text: string; doneText: string; delay: number; duration: number }

export const PROMPT = "user@cipincode ~ % "
export const COMMAND = "curl -fsSL https://cli.cipincode.com/install.sh | bash"
export const TYPE_SPEED = 38
export const TYPE_CODE_MS = 11
export const TYPE_CODE_NEWLINE_MS = 48

export const INIT_LINES: Line[] = []

export const OUTPUT_SEQUENCE: OutputItem[] = [
  {
    type: "spinner",
    text: "Downloading installer…",
    doneText: "Installer downloaded",
    delay: 240,
    duration: 900,
  },
  {
    type: "spinner",
    text: "Installing Go, Bun + cipincode…",
    doneText: "CipinCode CLI installed",
    delay: 140,
    duration: 1400,
  },
  {
    type: "line",
    text: "Run `cipincode login` to get started",
    color: "green",
    delay: 200,
  },
]

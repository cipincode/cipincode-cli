export type MotShClass =
  | "mot-sh-kw"
  | "mot-sh-str"
  | "mot-sh-fn"
  | "mot-sh-num"
  | "mot-sh-prop"
  | "mot-sh-var"
  | "mot-sh-punct"

export const CODE_SYNTAX_PARTS: ReadonlyArray<{ t: string; c: MotShClass }> = [
  { t: "package ", c: "mot-sh-kw" },
  { t: "main", c: "mot-sh-var" },
  { t: "\n\n", c: "mot-sh-punct" },
  { t: "import ", c: "mot-sh-kw" },
  { t: `"github.com/cipincode/cli/agent"`, c: "mot-sh-str" },
  { t: "\n\n", c: "mot-sh-punct" },
  { t: "func ", c: "mot-sh-kw" },
  { t: "main", c: "mot-sh-fn" },
  { t: "() {\n", c: "mot-sh-punct" },
  { t: "    ", c: "mot-sh-punct" },
  { t: "bot", c: "mot-sh-var" },
  { t: " := ", c: "mot-sh-punct" },
  { t: "agent", c: "mot-sh-fn" },
  { t: ".", c: "mot-sh-punct" },
  { t: "New", c: "mot-sh-fn" },
  { t: "(", c: "mot-sh-punct" },
  { t: `"coding"`, c: "mot-sh-str" },
  { t: ")\n", c: "mot-sh-punct" },
  { t: "    ", c: "mot-sh-punct" },
  { t: "bot", c: "mot-sh-var" },
  { t: ".", c: "mot-sh-punct" },
  { t: "Effort", c: "mot-sh-fn" },
  { t: "(", c: "mot-sh-punct" },
  { t: `"high"`, c: "mot-sh-str" },
  { t: ")\n", c: "mot-sh-punct" },
  { t: "    ", c: "mot-sh-punct" },
  { t: "if ", c: "mot-sh-kw" },
  { t: "err", c: "mot-sh-var" },
  { t: " := ", c: "mot-sh-punct" },
  { t: "bot", c: "mot-sh-var" },
  { t: ".", c: "mot-sh-punct" },
  { t: "Run", c: "mot-sh-fn" },
  { t: "(); ", c: "mot-sh-punct" },
  { t: "err", c: "mot-sh-var" },
  { t: " != ", c: "mot-sh-punct" },
  { t: "nil", c: "mot-sh-kw" },
  { t: " {\n", c: "mot-sh-punct" },
  { t: "        ", c: "mot-sh-punct" },
  { t: "panic", c: "mot-sh-fn" },
  { t: "(", c: "mot-sh-punct" },
  { t: "err", c: "mot-sh-var" },
  { t: ")\n", c: "mot-sh-punct" },
  { t: "    }\n", c: "mot-sh-punct" },
  { t: "}", c: "mot-sh-punct" },
]

export const CODE_DEMO = CODE_SYNTAX_PARTS.map((p) => p.t).join("")

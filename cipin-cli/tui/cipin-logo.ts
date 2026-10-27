// Cipin CLI: replaces the OpenCode wordmark on the home screen with the Cipin one.
// Loaded only by `cipin` (through OPENCODE_TUI_CONFIG), never by a plain `opencode`.
const LEFT = ["█▀▀▀ █ █▀▀█ █ █▀▀▄ ", "█___ █ █__█ █ █__█ ", "▀▀▀▀ ▀ █▀▀▀ ▀ ▀~~▀ "]
const RIGHT = ["█▀▀▀ █▀▀█ █▀▀█ █▀▀█", "█___ █__█ █__█ █^^^", "▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀"]

const VIOLET = "#8b5cf6"
const LIGHT = "#e5e7eb"

export default {
    id: "cipin-logo",
    tui: async (api: any) => {
        const { jsx } = await import("@opentui/solid/jsx-runtime")

        api.slots.register({
            slots: {
                home_logo: () =>
                    jsx("box", {
                        flexDirection: "column",
                        children: LEFT.map((left, row) =>
                            jsx("box", {
                                flexDirection: "row",
                                children: [
                                    jsx("text", { content: left, fg: VIOLET, selectable: false }),
                                    jsx("text", { content: " ", selectable: false }),
                                    jsx("text", { content: RIGHT[row], fg: LIGHT, selectable: false }),
                                ],
                            }),
                        ),
                    }),
            },
        })
    },
}

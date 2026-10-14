package main

import "github.com/charmbracelet/lipgloss"

type mode int

const (
	modeChat mode = iota
	modeCoding
	modeDebug
)

// next berpindah mode berurutan (Shift+Tab): chat, coding, debugging, lalu kembali.
func (m mode) next() mode { return (m + 1) % 3 }

// key adalah nilai yang dikirim ke server agent.
func (m mode) key() string {
	switch m {
	case modeCoding:
		return "coding"
	case modeDebug:
		return "debugging"
	}
	return "chat"
}

func (m mode) name() string {
	switch m {
	case modeCoding:
		return "Coding"
	case modeDebug:
		return "Debugging"
	}
	return "Chat"
}

func (m mode) icon() string {
	switch m {
	case modeCoding:
		return "⏵⏵"
	case modeDebug:
		return "✱"
	}
	return "●"
}

func (m mode) description() string {
	switch m {
	case modeCoding:
		return "menulis dan mengubah kode · auto-approve"
	case modeDebug:
		return "mencari dan memperbaiki bug"
	}
	return "tanya jawab biasa"
}

func (m mode) placeholder() string {
	switch m {
	case modeCoding:
		return "Jelaskan kode yang ingin dibuat atau diubah…"
	case modeDebug:
		return "Tempel error atau jelaskan bug-nya…"
	}
	return "Tanya apa saja, atau ketik /help"
}

func (m mode) color() lipgloss.Color {
	switch m {
	case modeCoding:
		return lipgloss.Color("#ea580c")
	case modeDebug:
		return lipgloss.Color("#fb7185")
	}
	return lipgloss.Color("#38bdf8")
}

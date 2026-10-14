package main

import "github.com/charmbracelet/lipgloss"

var (
	orange = lipgloss.Color("#ea580c")
	muted  = lipgloss.Color("#71717a")
	red    = lipgloss.Color("#f87171")

	brandStyle  = lipgloss.NewStyle().Bold(true).Foreground(orange)
	mutedStyle  = lipgloss.NewStyle().Foreground(muted)
	userStyle   = lipgloss.NewStyle().Bold(true).Foreground(orange)
	errorStyle  = lipgloss.NewStyle().Foreground(red)
	inputBorder = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(muted)
)

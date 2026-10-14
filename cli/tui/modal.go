package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// modalView adalah kotak informasi di tengah layar (/usage, /status, /help).
// Selama terbuka, semua tombol diabaikan kecuali Esc (tutup), panah (gulir), dan Ctrl+C.
type modalKind int

const (
	modalInfo   modalKind = iota // hanya dibaca: /usage, /status, /help
	modalEffort                  // slider: ← → mengubah level
	modalAgents                  // daftar pilihan agent: ↑↓ dan Enter
	modalProviders               // daftar pilihan provider API: ↑↓ dan Enter
	modalModels                  // daftar pilihan model: ↑↓ dan Enter
	modalConnect                 // konfirmasi provider + kunci API langsung
	modalUpdate                  // pengingat update: Spasi centang, Enter update, Esc nanti
)

type modalView struct {
	title  string
	lines  []string // isi, sudah berwarna
	scroll int
	kind   modalKind
	accent lipgloss.Color // warna bingkai dan judul
	footer string         // petunjuk tombol; kosong memakai bawaan
}

var modalBox = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1)

func newModal(title, body string) *modalView {
	return &modalView{title: title, lines: strings.Split(body, "\n"), accent: orange}
}

// chrome: baris yang dipakai bingkai, judul, garis, jarak, dan petunjuk.
const modalChrome = 8

func (v *modalView) capacity(screenHeight int) int { return max(3, screenHeight-modalChrome) }

func (v *modalView) maxScroll(screenHeight int) int {
	return max(0, len(v.lines)-v.capacity(screenHeight))
}

func (v *modalView) render(screenWidth, screenHeight int) string {
	capacity := v.capacity(screenHeight)
	scrollable := len(v.lines) > capacity
	end := min(len(v.lines), v.scroll+capacity)
	body := v.lines[v.scroll:end]

	footer := "Esc untuk menutup"
	if v.footer != "" {
		footer = v.footer
	}
	if scrollable && v.footer == "" {
		footer = fmt.Sprintf("↑↓ gulir (%d–%d dari %d) · Esc untuk menutup", v.scroll+1, end, len(v.lines))
	}

	width := max(lipgloss.Width(v.title), lipgloss.Width(footer))
	for _, line := range body {
		width = max(width, lipgloss.Width(line))
	}
	width = min(width, max(screenWidth-8, 20))

	clipped := make([]string, len(body))
	for i, line := range body {
		clipped[i] = ansi.Truncate(line, width, "…")
	}

	inner := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(v.accent).Render(v.title),
		mutedStyle.Render(strings.Repeat("─", width)),
		strings.Join(clipped, "\n"),
		"",
		mutedStyle.Render(footer),
	)
	return modalBox.BorderForeground(v.accent).Width(width + 2).Render(inner)
}

// overlay menempelkan box di tengah base tanpa merusak warna di sekelilingnya.
func overlay(base, box string, width int) string {
	baseLines := strings.Split(base, "\n")
	boxLines := strings.Split(box, "\n")

	boxWidth := lipgloss.Width(box)
	x := max(0, (width-boxWidth)/2)
	y := max(0, (len(baseLines)-len(boxLines))/2)

	for i, boxLine := range boxLines {
		row := y + i
		if row >= len(baseLines) {
			break
		}
		line := baseLines[row]

		left := ansi.Truncate(line, x, "")
		if pad := x - ansi.StringWidth(left); pad > 0 {
			left += strings.Repeat(" ", pad)
		}
		right := ansi.TruncateLeft(line, x+boxWidth, "")

		baseLines[row] = left + "\x1b[0m" + boxLine + "\x1b[0m" + right
	}
	return strings.Join(baseLines, "\n")
}

func (m model) openModal(title, body string) model {
	m.modal = newModal(title, body)
	m.layout()
	return m
}

// openEffortModal membuka slider effort. Tiap geseran langsung berlaku dan tersimpan.
func (m model) openEffortModal() model {
	v := newModal("Effort", effortBody(m.effort))
	v.kind = modalEffort
	v.accent = m.effort.color()
	v.footer = "← → ubah level · Esc untuk menutup"
	m.modal = v
	m.layout()
	return m
}

func (m model) keyEffortModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	delta := 0
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter:
		m.modal = nil
		m.layout()
		return m, nil
	case tea.KeyLeft:
		delta = -1
	case tea.KeyRight:
		delta = 1
	default:
		switch msg.String() {
		case "q":
			m.modal = nil
			m.layout()
			return m, nil
		case "h":
			delta = -1
		case "l":
			delta = 1
		}
	}

	if next := stepEffort(m.effort, delta); delta != 0 && next != m.effort {
		m.effort = next
		saveEffort(next)
		m = m.openEffortModal()
	}
	return m, nil
}

func (m model) keyModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.modal.kind {
	case modalEffort:
		return m.keyEffortModal(msg)
	case modalAgents:
		return m.keyAgentsModal(msg)
	case modalProviders:
		return m.keyProvidersModal(msg)
	case modalModels:
		return m.keyModelsModal(msg)
	case modalConnect:
		return m.keyConnectModal(msg)
	case modalUpdate:
		return m.keyUpdateModal(msg)
	}
	v := *m.modal
	limit := v.maxScroll(m.height)
	page := v.capacity(m.height)

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter:
		m.modal = nil
		m.layout()
		return m, nil
	case tea.KeyUp:
		v.scroll = max(0, v.scroll-1)
	case tea.KeyDown:
		v.scroll = min(limit, v.scroll+1)
	case tea.KeyPgUp:
		v.scroll = max(0, v.scroll-page)
	case tea.KeyPgDown:
		v.scroll = min(limit, v.scroll+page)
	default:
		switch msg.String() {
		case "q":
			m.modal = nil
			m.layout()
			return m, nil
		case "k":
			v.scroll = max(0, v.scroll-1)
		case "j":
			v.scroll = min(limit, v.scroll+1)
		}
	}

	m.modal = &v
	return m, nil
}

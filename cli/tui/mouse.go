package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

const (
	wheelStep   = 3
	toastTime   = 2 * time.Second
	toastCopied = "✓ Berhasil disalin ke clipboard"
	toastOSC    = "✓ Disalin ke clipboard terminal"
)

// viewportPlain adalah baris-baris percakapan yang sedang terlihat, tanpa warna.
func (m model) viewportPlain() []string {
	lines := strings.Split(m.viewport.View(), "\n")
	for i, line := range lines {
		lines[i] = ansi.Strip(line)
	}
	return lines
}

func (m model) toastCmd() tea.Cmd {
	id := m.toastID
	return tea.Tick(toastTime, func(time.Time) tea.Msg { return toastClearMsg{id} })
}

// handleMouse: roda mouse menggulir percakapan; menyeret di area percakapan memilih teks
// dan langsung menyalinnya saat tombol dilepas.
func (m model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.screen != screenChat {
		return m, nil
	}

	wheel := 0
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		wheel = -wheelStep
	case tea.MouseButtonWheelDown:
		wheel = wheelStep
	}

	if m.modal != nil { // modal hanya bisa digulir (jika berisi bacaan panjang)
		if wheel != 0 && m.modal.kind == modalInfo {
			v := *m.modal
			v.scroll = min(max(v.scroll+wheel, 0), v.maxScroll(m.height))
			m.modal = &v
		}
		return m, nil
	}

	if wheel != 0 {
		if wheel < 0 {
			m.viewport.ScrollUp(-wheel)
		} else {
			m.viewport.ScrollDown(wheel)
		}
		m.sel = nil // isi bergeser, sorotan lama tidak berlaku lagi
		return m, nil
	}

	if msg.Button != tea.MouseButtonLeft && !(msg.Action == tea.MouseActionRelease && m.sel != nil) {
		return m, nil
	}

	row, col := msg.Y-m.headerRows(), msg.X
	inView := row >= 0 && row < m.viewport.Height && len(m.menuItems()) == 0

	switch msg.Action {
	case tea.MouseActionPress:
		m.sel = nil
		if inView {
			plain := m.viewportPlain()
			kind, start, end := classifyRow(plain[row])
			if kind == rowBorder { // garis atas/bawah kotak kode tidak bisa dipilih
				return m, nil
			}
			if kind == rowContent { // mulai di dalam kotak: anchor selalu di area kode
				col = min(max(col, start), max(end-1, start))
			}
			m.sel = &selection{anchorRow: row, anchorCol: col, curRow: row, curCol: col, dragging: true}
		}

	case tea.MouseActionMotion:
		if m.sel != nil && m.sel.dragging {
			s := *m.sel
			s.curRow = min(max(row, 0), m.viewport.Height-1)
			s.curCol = min(max(col, 0), max(m.width-1, 0))

			// Seretan yang dimulai di dalam kotak kode tidak boleh keluar dari kotak itu.
			plain := m.viewportPlain()
			if first, last, ok := codeBox(plain, s.anchorRow); ok {
				_, start, end := classifyRow(plain[s.anchorRow])
				s.curRow = min(max(s.curRow, first), last)
				s.curCol = min(max(s.curCol, start), max(end-1, start))
			}
			s.moved = s.moved || s.curRow != s.anchorRow || s.curCol != s.anchorCol
			m.sel = &s
		}

	case tea.MouseActionRelease:
		if m.sel == nil || !m.sel.dragging {
			return m, nil
		}
		s := *m.sel
		s.dragging = false
		if !s.moved { // klik biasa, bukan seretan
			m.sel = nil
			return m, nil
		}
		m.sel = &s

		text := selectionText(m.viewportPlain(), s)
		if text == "" {
			return m, nil
		}

		m.toast = toastCopied
		if err := copyToClipboard(text); err != nil {
			m.toast = toastOSC
		}
		m.toastID++
		return m, m.toastCmd()
	}
	return m, nil
}

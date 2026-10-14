package main

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func sized(t *testing.T, w, h int) model {
	t.Helper()
	cfg := &Config{Name: "Naufal Hamka", Email: "customer@gmail.com"}
	last := time.Date(2026, 10, 13, 22, 5, 0, 0, time.Local)
	m := newModel(&Agent{URL: "http://127.0.0.1:4096"}, cfg, last, options{screen: screenChat})
	next, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return next.(model)
}

func plain(m model) string { return ansi.Strip(m.View()) }

func TestViewShowsLogoSessionAndMode(t *testing.T) {
	m := sized(t, 100, 30)
	view := plain(m)
	t.Log("\n" + view)

	for _, want := range []string{
		"Naufal Hamka", "customer@gmail.com", "Sesi terakhir: 13 Okt 2026, 22.05",
		"Chat mode", "shift+tab ganti mode", "Selamat datang, Naufal", "▄█",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("tampilan tidak memuat %q", want)
		}
	}
	if strings.Contains(view, "CipinCode ·") {
		t.Error("teks CipinCode seharusnya diganti logo")
	}
}

func TestShiftTabCyclesModesAndUpdatesLabel(t *testing.T) {
	m := sized(t, 100, 30)

	want := []string{"Coding mode", "Debugging mode", "Chat mode"}
	for _, label := range want {
		next, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
		m = next.(model)
		if !strings.Contains(plain(m), label) {
			t.Fatalf("setelah Shift+Tab tampilan tidak memuat %q", label)
		}
	}
	if !strings.Contains(plain(sizedAfter(t, 1)), "auto-approve") {
		t.Error("mode coding harus menandai auto-approve")
	}
}

func sizedAfter(t *testing.T, presses int) model {
	m := sized(t, 100, 30)
	for i := 0; i < presses; i++ {
		next, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
		m = next.(model)
	}
	return m
}

func TestCompactLayoutOnSmallTerminal(t *testing.T) {
	m := sized(t, 50, 20)
	view := plain(m)

	if strings.Contains(view, "Sesi terakhir") {
		t.Error("terminal kecil seharusnya memakai header ringkas")
	}
	if !strings.Contains(view, "Naufal Hamka") || !strings.Contains(view, "Chat mode") {
		t.Errorf("header ringkas atau label mode hilang:\n%s", view)
	}
}

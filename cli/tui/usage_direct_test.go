package main

import (
	"net/http"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func TestUsageDirectModeShowsProviderAccount(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.provider = "openai"
	m.modelID = "openai/gpt-4o-mini"
	if m.config != nil {
		if m.config.ProviderKeys == nil {
			m.config.ProviderKeys = map[string]string{}
		}
		m.config.ProviderKeys["openai"] = "sk-1234567890abcdef"
	}

	if !m.directMode() {
		t.Fatal("ada kunci sendiri seharusnya mode langsung")
	}

	m, _ = update(m, typed("/usage"))
	m, _ = update(m, key(tea.KeyEnter))
	if m.modal == nil {
		t.Fatal("/usage mode langsung seharusnya membuka modal")
	}
	view := plain(m)
	for _, want := range []string{"Mode langsung", "OpenAI", "cdef", "tanpa jatah server", "GPT-4o mini"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal usage langsung tidak memuat %q:\n%s", want, view)
		}
	}
	if strings.Contains(view, "Sesi 5 jam") {
		t.Errorf("mode langsung tidak boleh menampilkan batang kuota server:\n%s", view)
	}
}

func TestUsageDirectHeaderLine(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.provider = "anthropic"
	if m.config != nil {
		m.config.ClaudeAccount = true
	}
	line := ansi.Strip(m.usageLine())
	if !strings.Contains(line, "Langsung") || !strings.Contains(line, "akun Claude") {
		t.Errorf("header seharusnya menandai mode langsung akun Claude, got %q", line)
	}
	if strings.Contains(line, "Sesi 5 jam") {
		t.Errorf("header mode langsung tidak boleh memakai batang server, got %q", line)
	}
}

func TestUsageServerStillShowsMeters(t *testing.T) {
	reset := time.Now().Add(3 * time.Hour)
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(snapshotJSON(
			usageMeter{Used: 15, Limit: 40, Percent: 38, ResetsAt: &reset},
			usageMeter{Used: 48, Limit: 400, Percent: 12},
		))
	})
	m := newScreen(t, screenChat, loggedIn(srv.URL), "")
	m.provider = "" // Otomatis = via server
	if m.directMode() {
		t.Fatal("tanpa provider seharusnya bukan mode langsung")
	}
	m, _ = update(m, typed("/usage"))
	m, cmd := update(m, key(tea.KeyEnter))
	m, _ = update(m, cmd())
	if m.modal == nil {
		t.Fatal("/usage via server seharusnya membuka modal")
	}
	if view := plain(m); !strings.Contains(view, "Sesi 5 jam") {
		t.Errorf("via server harus tetap menampilkan batang kuota:\n%s", view)
	}
}

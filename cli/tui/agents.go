package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// agentItem adalah agent buatan pengguna di dashboard (menu Agents).
type agentItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
	Effort      string `json:"effort"`
}

type agentsMsg struct {
	items []agentItem
	err   error
}

const agentWindow = 9 // jumlah baris agent yang terlihat sekaligus di modal

// newUUID membuat UUID v4 untuk mengenali satu percakapan di server.
func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func fetchAgents(baseURL, key string) ([]agentItem, error) {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/cli/agents", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa terhubung ke %s", baseURL)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		return nil, errInvalidKey
	default:
		return nil, fmt.Errorf("server membalas status %d", resp.StatusCode)
	}

	var body struct {
		Agents []agentItem `json:"agents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("balasan server tidak valid: %w", err)
	}
	return body.Agents, nil
}

func fetchAgentsCmd(cfg *Config) tea.Cmd {
	return func() tea.Msg {
		items, err := fetchAgents(cfg.URL, cfg.Key)
		return agentsMsg{items: items, err: err}
	}
}

// agentsBody menggambar daftar pilihan: "Tanpa agent" lalu agent milik pengguna.
// cursor 0 berarti "Tanpa agent"; active adalah ID agent yang sedang dipakai (0 = tidak ada).
func agentsBody(items []agentItem, cursor, active int) string {
	total := len(items) + 1
	start := 0
	if cursor >= agentWindow {
		start = cursor - agentWindow + 1
	}
	end := min(start+agentWindow, total)

	selected := lipgloss.NewStyle().Bold(true).Foreground(orange)
	row := func(i int) string {
		name, detail, id := "Tanpa agent", "bawaan, tanpa instruksi tambahan", 0
		if i > 0 {
			it := items[i-1]
			name, id = it.Name, it.ID
			detail = it.Mode + " · " + it.Effort
		}
		mark := ""
		if id == active {
			mark = " ✓"
		}
		line := padRight(name, 24) + detail + mark
		if i == cursor {
			return selected.Render("❯ " + line)
		}
		return "  " + padRight(name, 24) + mutedStyle.Render(detail+mark)
	}

	lines := []string{mutedStyle.Render("Pilih agent untuk percakapan ini. Mode dan effort mengikuti pengaturan agent."), ""}
	for i := start; i < end; i++ {
		lines = append(lines, row(i))
	}
	if total > agentWindow {
		lines = append(lines, mutedStyle.Render(fmt.Sprintf("  ↕ %d dari %d, gunakan ↑↓", cursor+1, total)))
	}
	if cursor > 0 {
		if desc := items[cursor-1].Description; desc != "" {
			lines = append(lines, "", mutedStyle.Render(desc))
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) openAgentsModal() model {
	v := newModal("Agent", agentsBody(m.agentChoices, m.agentCursor, m.activeAgentID()))
	v.kind = modalAgents
	v.footer = "↑↓ pilih · Enter pakai · Esc tutup"
	m.modal = v
	m.layout()
	return m
}

func (m model) activeAgentID() int {
	if m.chosen == nil {
		return 0
	}
	return m.chosen.ID
}

func (m model) keyAgentsModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(m.agentChoices) + 1
	moved := false

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.modal = nil
		m.layout()
		return m, nil
	case tea.KeyUp:
		m.agentCursor, moved = (m.agentCursor-1+total)%total, true
	case tea.KeyDown:
		m.agentCursor, moved = (m.agentCursor+1)%total, true
	case tea.KeyEnter:
		return m.selectAgent(), nil
	default:
		switch msg.String() {
		case "q":
			m.modal = nil
			m.layout()
			return m, nil
		case "k":
			m.agentCursor, moved = (m.agentCursor-1+total)%total, true
		case "j":
			m.agentCursor, moved = (m.agentCursor+1)%total, true
		}
	}

	if moved {
		m = m.openAgentsModal()
	}
	return m, nil
}

// selectAgent memakai agent di bawah kursor dan menerapkan mode dan effort bawaannya.
func (m model) selectAgent() model {
	m.modal = nil

	if m.agentCursor == 0 {
		m.chosen = nil
		m = m.note(roleSystem, "Agent dinonaktifkan. Percakapan memakai pengaturan bawaan.")
	} else {
		it := m.agentChoices[m.agentCursor-1]
		m.chosen = &it
		m.mode = parseModeKey(it.Mode)
		m.effort = parseEffort(it.Effort)
		m.input.Placeholder = m.mode.placeholder()
		m = m.note(roleSystem, fmt.Sprintf("Agent aktif: %s · mode %s · effort %s", it.Name, m.mode.name(), m.effort.key()))
	}

	m.layout()
	m.refresh()
	return m
}

func parseModeKey(key string) mode {
	for _, md := range []mode{modeChat, modeCoding, modeDebug} {
		if md.key() == key {
			return md
		}
	}
	return modeChat
}

func (m model) handleAgents(msg agentsMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m = m.note(roleError, usageErrorText(msg.err))
		m.layout()
		m.refresh()
		return m, nil
	}

	if len(msg.items) == 0 {
		m = m.note(roleSystem, "Belum ada agent. Buat di dashboard, menu Agents (ketik /dashboard).")
		m.layout()
		m.refresh()
		return m, nil
	}

	m.agentChoices = msg.items
	m.agentCursor = 0
	for i, it := range msg.items { // kursor mulai di agent yang sedang aktif
		if m.agent != nil && it.ID == m.chosen.ID {
			m.agentCursor = i + 1
		}
	}
	return m.openAgentsModal(), nil
}

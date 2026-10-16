package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Jalur "Akun Claude": CLI resmi `claude` (Claude Code) dipakai sebagai
// jembatan. Login (OAuth ke browser) dan token langganan Pro/Max sepenuhnya
// ditangani Claude Code; CipinCode tidak pernah melihat atau menyimpan token.

const claudeInstallHint = "Pasang dulu: npm i -g @anthropic-ai/claude-code"

type claudeStatus struct {
	LoggedIn         bool   `json:"loggedIn"`
	AuthMethod       string `json:"authMethod"`
	Email            string `json:"email"`
	SubscriptionType string `json:"subscriptionType"`
}

// claudeBin mencari binary `claude` di PATH.
func claudeBin() (string, error) {
	p, err := exec.LookPath("claude")
	if err != nil {
		return "", errors.New("Claude Code belum terpasang. " + claudeInstallHint)
	}
	return p, nil
}

// claudeAuthStatus membaca status login Claude Code.
func claudeAuthStatus() (*claudeStatus, error) {
	bin, err := claudeBin()
	if err != nil {
		return nil, err
	}
	out, err := exec.Command(bin, "auth", "status", "--json").Output()
	var st claudeStatus
	if jerr := json.Unmarshal(out, &st); jerr != nil {
		if err != nil {
			return nil, fmt.Errorf("gagal membaca status login Claude: %w", err)
		}
		return nil, fmt.Errorf("status login Claude tidak terbaca: %w", jerr)
	}
	return &st, nil
}

// label ringkas akun untuk ditampilkan ke pengguna.
func (s claudeStatus) label() string {
	plan := s.SubscriptionType
	if plan != "" {
		plan = strings.ToUpper(plan[:1]) + plan[1:]
	}
	switch {
	case s.Email != "" && plan != "":
		return fmt.Sprintf("%s (%s)", s.Email, plan)
	case s.Email != "":
		return s.Email
	}
	return "akun Claude"
}

// useClaudeAccount: chat harus lewat Claude Code lokal, bukan server/kunci API.
func (m model) useClaudeAccount() bool {
	return m.provider == "anthropic" && m.config != nil && m.config.ClaudeAccount
}

type claudeLoginMsg struct{ err error }

// claudeLoginCmd menjalankan `claude auth login --claudeai`, yang membuka
// browser untuk login akun Pro/Max. Layar diserahkan ke proses itu sementara.
func claudeLoginCmd() tea.Cmd {
	bin, err := claudeBin()
	if err != nil {
		return func() tea.Msg { return claudeLoginMsg{err} }
	}
	cmd := exec.Command(bin, "auth", "login", "--claudeai")
	return tea.ExecProcess(cmd, func(err error) tea.Msg { return claudeLoginMsg{err} })
}

// claudeCheckCmd memeriksa status login tanpa membuka browser.
func claudeCheckCmd() tea.Cmd {
	return func() tea.Msg {
		st, err := claudeAuthStatus()
		if err != nil {
			return claudeLoginMsg{err}
		}
		if !st.LoggedIn || st.AuthMethod != "claude.ai" {
			return claudeLoginMsg{errNeedLogin}
		}
		return claudeLoginMsg{nil}
	}
}

var errNeedLogin = errors.New("belum login akun Claude")

// ---- modal pilihan metode ----

type claudeChoice struct{ name, desc string }

var claudeChoices = []claudeChoice{
	{"Akun Claude", "login Pro/Max lewat browser (tanpa kunci API)"},
	{"Kunci API", "tempel kunci dari console.anthropic.com"},
	{"Via server", "tetap lewat server CipinCode"},
}

func claudeMethodBody(cursor int) string {
	selected := lipgloss.NewStyle().Bold(true).Foreground(orange)
	lines := []string{mutedStyle.Render("Pilih cara terhubung ke Anthropic."), ""}
	for i, c := range claudeChoices {
		if i == cursor {
			lines = append(lines, selected.Render("❯ "+padRight(c.name, 14)+c.desc))
		} else {
			lines = append(lines, "  "+padRight(c.name, 14)+mutedStyle.Render(c.desc))
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) openClaudeMethodModal() model {
	v := newModal("Connect · Anthropic", claudeMethodBody(m.claudeCursor))
	v.kind = modalClaudeMethod
	v.footer = "↑↓ pilih · Enter pakai · Esc lewati"
	m.modal = v
	m.layout()
	return m
}

func (m model) keyClaudeMethodModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(claudeChoices)
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.modal = nil
		m.connectPending = nil
		m = m.note(roleSystem, "Anthropic via server CipinCode.")
		m.layout()
		m.refresh()
		return m, nil
	case tea.KeyUp:
		m.claudeCursor = (m.claudeCursor - 1 + total) % total
		return m.openClaudeMethodModal(), nil
	case tea.KeyDown:
		m.claudeCursor = (m.claudeCursor + 1) % total
		return m.openClaudeMethodModal(), nil
	case tea.KeyEnter:
		return m.chooseClaudeMethod()
	}
	return m, nil
}

func (m model) chooseClaudeMethod() (tea.Model, tea.Cmd) {
	switch m.claudeCursor {
	case 0:
		m.modal = nil
		m.connectPending = nil
		if _, err := claudeBin(); err != nil {
			m = m.note(roleError, err.Error())
			m.layout()
			m.refresh()
			return m, nil
		}
		m = m.note(roleSystem, "Memeriksa login akun Claude…")
		m.layout()
		m.refresh()
		return m, claudeCheckCmd()
	case 1:
		cmd := m.startKeyEntry()
		return m, cmd
	}
	m.modal = nil
	m.connectPending = nil
	if m.config != nil {
		m.config.ClaudeAccount = false
		delete(m.config.ProviderKeys, "anthropic")
		_ = saveConfig(*m.config)
	}
	m = m.note(roleSystem, "Anthropic via server CipinCode.")
	m.layout()
	m.refresh()
	return m, nil
}

// handleClaudeLogin menanggapi hasil pemeriksaan / login akun Claude.
func (m model) handleClaudeLogin(msg claudeLoginMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.err == nil:
		st, err := claudeAuthStatus()
		if err != nil || !st.LoggedIn {
			m = m.note(roleError, "Login akun Claude belum berhasil. Coba /connect lagi.")
			break
		}
		if m.config != nil {
			m.config.ClaudeAccount = true
			delete(m.config.ProviderKeys, "anthropic")
			_ = saveConfig(*m.config)
		}
		m = m.note(roleSystem, fmt.Sprintf("Terhubung ke %s · API kini memakai %s lewat akun Claude (tanpa kunci API).", st.label(), modelName(m.modelID)))
	case errors.Is(msg.err, errNeedLogin):
		m = m.note(roleSystem, "Membuka browser untuk login akun Claude…")
		m.layout()
		m.refresh()
		return m, claudeLoginCmd()
	default:
		m = m.note(roleError, msg.err.Error())
	}
	m.layout()
	m.refresh()
	return m, nil
}

// ---- chat lewat `claude -p` ----

// claudeModelArg memetakan model katalog ke argumen --model Claude Code.
func claudeModelArg(modelID string) string {
	for _, p := range providers {
		for _, mdl := range p.models {
			if mdl.id == modelID && mdl.direct != "" {
				return mdl.direct
			}
		}
	}
	return strings.TrimPrefix(modelID, "anthropic/")
}

var claudeSystems = map[string]string{
	"chat":      "You are a helpful coding assistant. Answer concisely in Indonesian.",
	"coding":    "You are an expert coding assistant. Write clean, working code with brief explanations in Indonesian.",
	"debugging": "You are a debugging expert. Analyze errors deeply, find root causes, and explain fixes in Indonesian.",
}

// claudePrompt menyusun riwayat + pesan baru menjadi satu prompt teks.
func claudePrompt(history []historyItem, message string) string {
	var b strings.Builder
	if len(history) > 0 {
		b.WriteString("Percakapan sebelumnya:\n")
		for _, h := range history {
			who := "Pengguna"
			if h.Role == "assistant" {
				who = "Asisten"
			}
			fmt.Fprintf(&b, "%s: %s\n\n", who, h.Content)
		}
		b.WriteString("Pesan baru dari pengguna:\n")
	}
	b.WriteString(message)
	return b.String()
}

func claudeArgs(req chatRequest) []string {
	system := claudeSystems[req.Mode]
	if system == "" {
		system = claudeSystems["chat"]
	}
	args := []string{
		"-p",
		"--output-format", "stream-json",
		"--include-partial-messages",
		"--verbose",
		"--tools", "",
		"--strict-mcp-config",
		"--disable-slash-commands",
		"--no-session-persistence",
		"--system-prompt", system,
	}
	if req.Model != "" {
		args = append(args, "--model", claudeModelArg(req.Model))
	}
	return args
}

type claudeEvent struct {
	Type   string `json:"type"`
	Result string `json:"result"`
	IsErr  bool   `json:"is_error"`
	Event  struct {
		Type  string `json:"type"`
		Delta struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"delta"`
	} `json:"event"`
}

// readClaudeStream membaca stream-json dan meneruskan teks ke out.
func readClaudeStream(lines *bufio.Scanner, out chan<- chatEvent) error {
	got := false
	for lines.Scan() {
		var ev claudeEvent
		if err := json.Unmarshal(lines.Bytes(), &ev); err != nil {
			continue
		}
		switch ev.Type {
		case "stream_event":
			if ev.Event.Type == "content_block_delta" && ev.Event.Delta.Type == "text_delta" && ev.Event.Delta.Text != "" {
				got = true
				out <- chatEvent{kind: "token", text: ev.Event.Delta.Text}
			}
		case "result":
			if ev.IsErr {
				return errors.New(claudeErrorText(ev.Result))
			}
			if !got && ev.Result != "" {
				out <- chatEvent{kind: "token", text: ev.Result}
			}
			return nil
		}
	}
	return lines.Err()
}

func claudeErrorText(s string) string {
	low := strings.ToLower(s)
	switch {
	case strings.Contains(low, "login") || strings.Contains(low, "authenticat") || strings.Contains(low, "401"):
		return "Sesi akun Claude habis. Jalankan /connect → Anthropic → Akun Claude untuk login lagi."
	case strings.Contains(low, "limit"):
		return "Batas pemakaian akun Claude tercapai: " + s
	}
	if s == "" {
		return "Claude Code mengembalikan galat."
	}
	return s
}

// chatClaude menjalankan satu giliran lewat Claude Code lokal.
func chatClaude(ctx context.Context, req chatRequest, out chan<- chatEvent) error {
	bin, err := claudeBin()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, bin, claudeArgs(req)...)
	cmd.Stdin = strings.NewReader(claudePrompt(req.History, req.Message))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("gagal menjalankan Claude Code: %w", err)
	}

	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	rerr := readClaudeStream(sc, out)
	werr := cmd.Wait()
	if rerr != nil {
		return rerr
	}
	if werr != nil && ctx.Err() == nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return errors.New(claudeErrorText(lastLines(msg, 5)))
		}
		return fmt.Errorf("Claude Code berhenti: %w", werr)
	}
	return nil
}

package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// provider adalah penyedia API AI (lewat OpenRouter). ID kosong berarti
// Otomatis: model mengikuti paket pengguna seperti sebelumnya.
type provider struct {
	id      string
	name    string
	desc    string
	models  []providerModel
	baseURL string // endpoint langsung (BYOK); kosong = hanya via server
}

type providerModel struct {
	id     string
	name   string
	desc   string
	direct string // ID model untuk API langsung (kosong = samakan id)
}

// providers adalah katalog yang didukung /connect. Daftar model di sini
// harus sama dengan config/ai.php di dashboard agar pilihan benar dipakai.
var providers = []provider{
	{
		id:      "anthropic",
		name:    "Anthropic",
		desc:    "Claude · kuat untuk coding dan instruksi panjang",
		baseURL: "https://api.anthropic.com",
		models: []providerModel{
			{id: "anthropic/claude-sonnet-4-5", direct: "claude-sonnet-4-5-20250929", name: "Claude Sonnet 4.5", desc: "seimbang untuk coding harian"},
			{id: "anthropic/claude-opus-4-1", direct: "claude-opus-4-1-20250805", name: "Claude Opus 4.1", desc: "paling pintar, paling mahal"},
			{id: "anthropic/claude-3-5-haiku", direct: "claude-3-5-haiku-20241022", name: "Claude Haiku 3.5", desc: "cepat dan hemat"},
		},
	},
	{
		id:      "openai",
		name:    "OpenAI",
		desc:    "GPT · serbaguna dan cepat",
		baseURL: "https://api.openai.com/v1",
		models: []providerModel{
			{id: "openai/gpt-4o", name: "GPT-4o", desc: "serbaguna untuk chat dan kode"},
			{id: "openai/gpt-4o-mini", name: "GPT-4o mini", desc: "ringan dan hemat"},
			{id: "openai/o4-mini", name: "o4-mini", desc: "nalar cepat untuk debugging"},
		},
	},
	{
		id:      "google",
		name:    "Google",
		desc:    "Gemini · konteks sangat panjang",
		baseURL: "https://generativelanguage.googleapis.com/v1beta/openai",
		models: []providerModel{
			{id: "google/gemini-2.5-pro", name: "Gemini 2.5 Pro", desc: "konteks besar untuk repo besar"},
			{id: "google/gemini-2.5-flash", name: "Gemini 2.5 Flash", desc: "cepat untuk tugas ringan"},
		},
	},
	{
		id:      "xai",
		name:    "xAI",
		desc:    "Grok · respons lugas dan terkini",
		baseURL: "https://api.x.ai/v1",
		models: []providerModel{
			{id: "x-ai/grok-3", name: "Grok 3", desc: "utama untuk chat dan kode"},
			{id: "x-ai/grok-3-mini", name: "Grok 3 Mini", desc: "cepat dan hemat"},
		},
	},
	{
		id:      "deepseek",
		name:    "DeepSeek",
		desc:    "nalar mendalam, harga hemat",
		baseURL: "https://api.deepseek.com/v1",
		models: []providerModel{
			{id: "deepseek/deepseek-chat", name: "DeepSeek Chat", desc: "chat dan coding harian"},
			{id: "deepseek/deepseek-r1", name: "DeepSeek R1", desc: "nalar untuk masalah sulit"},
		},
	},
	{
		id:   "meta",
		name: "Meta",
		desc: "Llama · terbuka dan cepat",
		models: []providerModel{
			{id: "meta-llama/llama-3.3-70b-instruct", name: "Llama 3.3 70B", desc: "terbuka untuk umum"},
		},
	},
}

const providerWindow = 9 // jumlah baris provider yang terlihat sekaligus di modal

// parseProvider mencari provider dari ID; tidak dikenal berarti Otomatis.
func parseProvider(id string) *provider {
	for i := range providers {
		if providers[i].id == id {
			return &providers[i]
		}
	}
	return nil
}

// providerModels mengembalikan daftar model untuk ID provider ("auto" kosong).
func providerModels(id string) []providerModel {
	if p := parseProvider(id); p != nil {
		return p.models
	}
	return nil
}

// providerName untuk ditampilkan ("Otomatis" jika kosong/tak dikenal).
func providerName(id string) string {
	if p := parseProvider(id); p != nil {
		return p.name
	}
	return "Otomatis"
}

// modelName mencari nama pendek model dari ID-nya di semua provider.
func modelName(id string) string {
	if id == "" {
		return "Otomatis"
	}
	for _, p := range providers {
		for _, mdl := range p.models {
			if mdl.id == id {
				return mdl.name
			}
		}
	}
	return id
}

// providersBody menggambar daftar pilihan: "Otomatis" lalu semua provider.
func providersBody(cursor int, active string) string {
	total := len(providers) + 1
	start := 0
	if cursor >= providerWindow {
		start = cursor - providerWindow + 1
	}
	end := min(start+providerWindow, total)

	selected := lipgloss.NewStyle().Bold(true).Foreground(orange)
	row := func(i int) string {
		name, detail, id := "Otomatis", "ikut paket Anda", ""
		if i > 0 {
			p := providers[i-1]
			name, detail, id = p.name, p.desc, p.id
		}
		mark := ""
		if id == active {
			mark = " ✓"
		}
		line := padRight(name, 14) + detail + mark
		if i == cursor {
			return selected.Render("❯ " + line)
		}
		return "  " + padRight(name, 14) + mutedStyle.Render(detail+mark)
	}

	lines := []string{mutedStyle.Render("Pilih penyedia API. Model mengikuti provider yang dipilih."), ""}
	for i := start; i < end; i++ {
		lines = append(lines, row(i))
	}
	if total > providerWindow {
		lines = append(lines, mutedStyle.Render(fmt.Sprintf("  ↕ %d dari %d, gunakan ↑↓", cursor+1, total)))
	}
	return strings.Join(lines, "\n")
}

// modelsBody menggambar daftar model milik satu provider.
func modelsBody(p *provider, models []providerModel, cursor int, active string) string {
	total := len(models)
	start := 0
	if cursor >= providerWindow {
		start = cursor - providerWindow + 1
	}
	end := min(start+providerWindow, total)

	selected := lipgloss.NewStyle().Bold(true).Foreground(orange)
	lines := []string{mutedStyle.Render("Model " + p.name + " yang dipakai untuk percakapan ini."), ""}
	for i := start; i < end; i++ {
		mdl := models[i]
		mark := ""
		if mdl.id == active {
			mark = " ✓"
		}
		line := padRight(mdl.name, 24) + mdl.desc + mark
		if i == cursor {
			lines = append(lines, selected.Render("❯ "+line))
		} else {
			lines = append(lines, "  "+padRight(mdl.name, 24)+mutedStyle.Render(mdl.desc+mark))
		}
	}
	return strings.Join(lines, "\n")
}

func (m model) openProvidersModal() model {
	v := newModal("Connect", providersBody(m.providerCursor, m.provider))
	v.kind = modalProviders
	v.footer = "↑↓ pilih · Enter pakai · Esc tutup"
	m.modal = v
	m.layout()
	return m
}

func (m model) openModelsModal() model {
	p := parseProvider(m.provider)
	v := newModal("Model · "+p.name, modelsBody(p, providerModels(m.provider), m.modelCursor, m.modelID))
	v.kind = modalModels
	v.footer = "↑↓ pilih · Enter pakai · Esc tutup"
	m.modal = v
	m.layout()
	return m
}

func (m model) keyProvidersModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(providers) + 1
	moved := false

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.modal = nil
		m.layout()
		return m, nil
	case tea.KeyUp:
		m.providerCursor, moved = (m.providerCursor-1+total)%total, true
	case tea.KeyDown:
		m.providerCursor, moved = (m.providerCursor+1)%total, true
	case tea.KeyEnter:
		return m.selectProvider()
	default:
		switch msg.String() {
		case "q":
			m.modal = nil
			m.layout()
			return m, nil
		case "k":
			m.providerCursor, moved = (m.providerCursor-1+total)%total, true
		case "j":
			m.providerCursor, moved = (m.providerCursor+1)%total, true
		}
	}

	if moved {
		m = m.openProvidersModal()
	}
	return m, nil
}

func (m model) keyModelsModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	models := providerModels(m.provider)
	total := max(len(models), 1)
	moved := false

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.modal = nil
		m.layout()
		return m, nil
	case tea.KeyUp:
		m.modelCursor, moved = (m.modelCursor-1+total)%total, true
	case tea.KeyDown:
		m.modelCursor, moved = (m.modelCursor+1)%total, true
	case tea.KeyEnter:
		return m.selectModel(), nil
	default:
		switch msg.String() {
		case "q":
			m.modal = nil
			m.layout()
			return m, nil
		case "k":
			m.modelCursor, moved = (m.modelCursor-1+total)%total, true
		case "j":
			m.modelCursor, moved = (m.modelCursor+1)%total, true
		}
	}

	if moved {
		m = m.openModelsModal()
	}
	return m, nil
}

// selectProvider memakai provider di bawah kursor lalu membuka modal
// konfirmasi: API otomatis berganti mengikuti provider itu dan model
// direset ke bawaan provider.
func (m model) selectProvider() (tea.Model, tea.Cmd) {
	m.modal = nil

	if m.providerCursor == 0 {
		m.provider, m.modelID = "", ""
		m = m.note(roleSystem, "Provider Otomatis. API mengikuti paket Anda.")
		m.layout()
		m.refresh()
		return m, nil
	}

	p := providers[m.providerCursor-1]
	m.provider = p.id
	m.modelID = ""
	if len(p.models) > 0 {
		m.modelID = p.models[0].id
	}
	saveProvider(m.provider)
	saveModel(m.modelID)

	if p.baseURL == "" {
		m = m.note(roleSystem, fmt.Sprintf("Provider: %s · API kini memakai %s (%s) via server CipinCode.", p.name, modelName(m.modelID), p.name))
		m.layout()
		m.refresh()
		return m, nil
	}

	m.connectPending = &providers[m.providerCursor-1]
	if p.id == "anthropic" {
		m.claudeCursor = 0
		return m.openClaudeMethodModal(), nil
	}
	cmd := m.startKeyEntry()
	return m, cmd
}

// startKeyEntry membuka modal tempel kunci API untuk connectPending.
func (m *model) startKeyEntry() tea.Cmd {
	keyIn := textinput.New()
	keyIn.Prompt = "Kunci ❯ "
	keyIn.EchoMode = textinput.EchoPassword
	keyIn.CharLimit = 200
	keyIn.Width = 40
	keyIn.Focus()
	m.connectInput = keyIn
	*m = m.openConnectModal()
	return textinput.Blink
}

// providerKey mengembalikan kunci API milik pengguna untuk provider aktif.
func (m model) providerKey() string {
	if m.config == nil || m.provider == "" {
		return ""
	}
	return m.config.ProviderKeys[m.provider]
}

// connectModalBody adalah isi modal konfirmasi + kunci API langsung.
func connectModalBody(p *provider, modelID, typed string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Provider: %s · API kini memakai %s (%s)\n\n", p.name, modelName(modelID), p.name)
	b.WriteString("Punya kunci API " + p.name + "? Tempel agar\n")
	b.WriteString("terhubung langsung (tanpa lewat server dan kuota).\n")
	b.WriteString("Kosongkan untuk tetap via server CipinCode.\n\n")
	b.WriteString("Kunci ❯ " + typed + "\n\n")
	b.WriteString("Enter: simpan · Esc: lewati")
	return b.String()
}

func (m model) openConnectModal() model {
	p := m.connectPending
	title, body := "Connect", ""
	if p != nil {
		title = "Connect · " + p.name
		body = connectModalBody(p, m.modelID, m.connectInput.View())
	}
	v := newModal(title, body)
	v.kind = modalConnect
	v.footer = "Enter simpan · Esc lewati"
	m.modal = v
	m.layout()
	return m
}

func (m model) keyConnectModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.modal = nil
		m.connectPending = nil
		m = m.note(roleSystem, "Kunci dilewati. API via server CipinCode.")
		m.layout()
		m.refresh()
		return m, nil
	case tea.KeyEnter:
		return m.saveConnectKey(), nil
	}

	var cmd tea.Cmd
	m.connectInput, cmd = m.connectInput.Update(msg)
	m = m.openConnectModal()
	return m, cmd
}

// saveConnectKey menyimpan kunci API langsung (boleh kosong = via server).
func (m model) saveConnectKey() model {
	p := m.connectPending
	key := strings.TrimSpace(m.connectInput.Value())
	m.modal = nil
	m.connectPending = nil

	if p == nil {
		m.layout()
		m.refresh()
		return m
	}
	if m.config != nil {
		if m.config.ProviderKeys == nil {
			m.config.ProviderKeys = map[string]string{}
		}
		if key == "" {
			delete(m.config.ProviderKeys, p.id)
		} else {
			m.config.ProviderKeys[p.id] = key
			m.config.ClaudeAccount = false
		}
		_ = saveConfig(*m.config)
	}

	if key == "" {
		m = m.note(roleSystem, fmt.Sprintf("Provider: %s · API kini memakai %s (%s) via server CipinCode.", p.name, modelName(m.modelID), p.name))
	} else {
		m = m.note(roleSystem, fmt.Sprintf("Provider: %s · API kini memakai %s (%s) langsung dengan kunci Anda.", p.name, modelName(m.modelID), p.name))
	}
	m.layout()
	m.refresh()
	return m
}

// selectModel memakai model di bawah kursor untuk provider yang aktif.
func (m model) selectModel() model {
	m.modal = nil

	models := providerModels(m.provider)
	if m.modelCursor >= 0 && m.modelCursor < len(models) {
		mdl := models[m.modelCursor]
		m.modelID = mdl.id
		saveModel(m.modelID)
		m = m.note(roleSystem, fmt.Sprintf("Model: %s · API kini memakai %s", mdl.name, mdl.id))
	}

	m.layout()
	m.refresh()
	return m
}

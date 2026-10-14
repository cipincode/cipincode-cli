package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// provider adalah penyedia API AI (lewat OpenRouter). ID kosong berarti
// Otomatis: model mengikuti paket pengguna seperti sebelumnya.
type provider struct {
	id     string
	name   string
	desc   string
	models []providerModel
}

type providerModel struct {
	id   string
	name string
	desc string
}

// providers adalah katalog yang didukung /connect. Daftar model di sini
// harus sama dengan config/ai.php di dashboard agar pilihan benar dipakai.
var providers = []provider{
	{
		id:   "anthropic",
		name: "Anthropic",
		desc: "Claude · kuat untuk coding dan instruksi panjang",
		models: []providerModel{
			{id: "anthropic/claude-sonnet-4-5", name: "Claude Sonnet 4.5", desc: "seimbang untuk coding harian"},
			{id: "anthropic/claude-opus-4-1", name: "Claude Opus 4.1", desc: "paling pintar, paling mahal"},
			{id: "anthropic/claude-3-5-haiku", name: "Claude Haiku 3.5", desc: "cepat dan hemat"},
		},
	},
	{
		id:   "openai",
		name: "OpenAI",
		desc: "GPT · serbaguna dan cepat",
		models: []providerModel{
			{id: "openai/gpt-4o", name: "GPT-4o", desc: "serbaguna untuk chat dan kode"},
			{id: "openai/gpt-4o-mini", name: "GPT-4o mini", desc: "ringan dan hemat"},
			{id: "openai/o4-mini", name: "o4-mini", desc: "nalar cepat untuk debugging"},
		},
	},
	{
		id:   "google",
		name: "Google",
		desc: "Gemini · konteks sangat panjang",
		models: []providerModel{
			{id: "google/gemini-2.5-pro", name: "Gemini 2.5 Pro", desc: "konteks besar untuk repo besar"},
			{id: "google/gemini-2.5-flash", name: "Gemini 2.5 Flash", desc: "cepat untuk tugas ringan"},
		},
	},
	{
		id:   "xai",
		name: "xAI",
		desc: "Grok · respons lugas dan terkini",
		models: []providerModel{
			{id: "x-ai/grok-3", name: "Grok 3", desc: "utama untuk chat dan kode"},
			{id: "x-ai/grok-3-mini", name: "Grok 3 Mini", desc: "cepat dan hemat"},
		},
	},
	{
		id:   "deepseek",
		name: "DeepSeek",
		desc: "nalar mendalam, harga hemat",
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
		return m.selectProvider(), nil
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

// selectProvider memakai provider di bawah kursor. API otomatis berganti
// mengikuti provider itu dan model direset ke bawaan provider.
func (m model) selectProvider() model {
	m.modal = nil

	if m.providerCursor == 0 {
		m.provider, m.modelID = "", ""
		m = m.note(roleSystem, "Provider Otomatis. API mengikuti paket Anda.")
	} else {
		p := providers[m.providerCursor-1]
		m.provider = p.id
		m.modelID = ""
		if len(p.models) > 0 {
			m.modelID = p.models[0].id
		}
		saveProvider(m.provider)
		saveModel(m.modelID)
		m = m.note(roleSystem, fmt.Sprintf("Provider: %s · API kini memakai %s (%s)", p.name, modelName(m.modelID), p.name))
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

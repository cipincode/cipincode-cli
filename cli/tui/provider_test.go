package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func openConnect(t *testing.T, m model) model {
	t.Helper()
	m, _ = update(m, typed("/connect"))
	m, _ = update(m, key(tea.KeyEnter))
	return m
}

func TestProviderCatalogValid(t *testing.T) {
	if len(providers) != 6 {
		t.Fatalf("seharusnya 6 provider, got %d", len(providers))
	}
	seen := map[string]bool{}
	for _, p := range providers {
		if p.id == "" || p.name == "" {
			t.Errorf("provider tanpa id/nama: %+v", p)
		}
		if len(p.models) == 0 {
			t.Errorf("provider %s tidak punya model", p.id)
		}
		for _, mdl := range p.models {
			if mdl.id == "" || seen[mdl.id] {
				t.Errorf("model duplikat atau kosong: %q", mdl.id)
			}
			seen[mdl.id] = true
		}
	}
}

func TestParseProviderFallback(t *testing.T) {
	if parseProvider("anthropic") == nil || parseProvider("anthropic").name != "Anthropic" {
		t.Error("parseProvider(anthropic) seharusnya ketemu")
	}
	if parseProvider("") != nil || parseProvider("ngawur") != nil {
		t.Error("ID tak dikenal seharusnya nil (Otomatis)")
	}
	if providerName("") != "Otomatis" || providerName("openai") != "OpenAI" {
		t.Error("providerName salah")
	}
	if modelName("") != "Otomatis" || modelName("openai/gpt-4o") != "GPT-4o" {
		t.Error("modelName salah")
	}
	if modelName("entah/apa") != "entah/apa" {
		t.Error("model tak dikenal seharusnya dikembalikan apa adanya")
	}
}

func TestConnectOpensProviderModal(t *testing.T) {
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))

	if m.modal == nil || m.modal.kind != modalProviders {
		t.Fatal("/connect seharusnya membuka modal provider")
	}
	view := plain(m)
	for _, want := range []string{"Connect", "Otomatis", "Anthropic", "OpenAI", "↑↓ pilih"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal tidak memuat %q:\n%s", want, view)
		}
	}

	m, _ = update(m, key(tea.KeyEsc))
	if m.modal != nil {
		t.Error("Esc seharusnya menutup modal")
	}
}

func TestSelectProviderOpensConnectModal(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))

	m, _ = update(m, key(tea.KeyDown)) // kursor ke Anthropic
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal == nil || m.modal.kind != modalClaudeMethod {
		t.Fatal("Enter pada Anthropic seharusnya membuka pilihan metode")
	}
	for _, want := range []string{"Akun Claude", "Kunci API", "Via server"} {
		if !strings.Contains(plain(m), want) {
			t.Errorf("modal metode tidak memuat %q", want)
		}
	}
	m, _ = update(m, key(tea.KeyDown)) // Kunci API
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal == nil || m.modal.kind != modalConnect {
		t.Fatal("Kunci API seharusnya membuka modal connect")
	}
	view := plain(m)
	for _, want := range []string{"Connect · Anthropic", "API kini memakai Claude Sonnet 4.5 (Anthropic)", "Kunci"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal connect tidak memuat %q:\n%s", want, view)
		}
	}
}

func TestSelectAutoProvider(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))

	m, _ = update(m, key(tea.KeyEnter)) // kursor 0 = Otomatis

	if m.provider != "" || m.modelID != "" {
		t.Fatalf("Otomatis seharusnya mengosongkan provider dan model, got %q/%q", m.provider, m.modelID)
	}
	if !strings.Contains(plain(m), "Otomatis") {
		t.Error("seharusnya ada catatan kembali ke Otomatis")
	}
}

func TestModelNeedsProviderFirst(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, _ = update(m, typed("/model"))
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal != nil {
		t.Fatal("tanpa provider, /model seharusnya tidak membuka modal")
	}
	if !strings.Contains(plain(m), "/connect") {
		t.Error("seharusnya menyarankan /connect")
	}
}

func TestModelModalListsAndSelects(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.provider = "openai"
	m.modelID = ""

	m, _ = update(m, typed("/model"))
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal == nil || m.modal.kind != modalModels {
		t.Fatal("/model seharusnya membuka modal daftar model")
	}
	if view := plain(m); !strings.Contains(view, "GPT-4o") {
		t.Errorf("modal seharusnya memuat model OpenAI:\n%s", view)
	}

	m, _ = update(m, key(tea.KeyDown)) // ke GPT-4o mini
	m, _ = update(m, key(tea.KeyEnter))

	if m.modelID != "openai/gpt-4o-mini" {
		t.Fatalf("modelID = %q, seharusnya openai/gpt-4o-mini", m.modelID)
	}
	if loadModel() != "openai/gpt-4o-mini" {
		t.Error("model seharusnya tersimpan")
	}
}

func TestConnectModalSaveSkipsKey(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	newConnect := func() model {
		m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))
		m, _ = update(m, key(tea.KeyDown))
		m, _ = update(m, key(tea.KeyEnter))
		m, _ = update(m, key(tea.KeyDown)) // metode: Kunci API
		m, _ = update(m, key(tea.KeyEnter))
		return m
	}

	// Esc = lewati kunci, provider tetap terpilih via server.
	m := newConnect()
	m, _ = update(m, key(tea.KeyEsc))
	if m.modal != nil {
		t.Fatal("Esc seharusnya menutup modal connect")
	}
	if m.provider != "anthropic" || m.providerKey() != "" {
		t.Fatalf("provider=%q key=%q, seharusnya anthropic tanpa kunci", m.provider, m.providerKey())
	}
	if !strings.Contains(plain(m), "via server") {
		t.Error("seharusnya ada catatan via server")
	}

	// Ketik kunci + Enter = tersimpan dan langsung.
	m = newConnect()
	for _, r := range "sk-ant-test" {
		m, _ = update(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m, _ = update(m, key(tea.KeyEnter))
	if m.modal != nil {
		t.Fatal("Enter seharusnya menutup modal connect")
	}
	if m.providerKey() != "sk-ant-test" {
		t.Fatalf("kunci tidak tersimpan, got %q", m.providerKey())
	}
	cfg, err := loadConfig()
	if err != nil || cfg == nil || cfg.ProviderKeys["anthropic"] != "sk-ant-test" {
		t.Fatalf("config seharusnya menyimpan kunci: %+v %v", cfg, err)
	}
	if view := plain(m); !strings.Contains(view, "langsung dengan kunci Anda") {
		t.Errorf("seharusnya ada catatan langsung:\n%s", view)
	}
}

func TestMetaProviderHasNoDirectKey(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))
	for i := 0; i < len(providers); i++ {
		m, _ = update(m, key(tea.KeyDown))
	}
	m, _ = update(m, key(tea.KeyEnter))
	if m.modal != nil {
		t.Fatal("Meta seharusnya catatan biasa (tanpa modal kunci)")
	}
	if m.provider != "meta" {
		t.Fatalf("provider = %q, seharusnya meta", m.provider)
	}
}

func TestClaudeAccountSelectedAfterLogin(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))
	m, _ = update(m, key(tea.KeyDown))
	m, _ = update(m, key(tea.KeyEnter)) // Anthropic -> pilihan metode

	if m.useClaudeAccount() {
		t.Fatal("belum login, seharusnya belum memakai akun Claude")
	}
	m.config.ClaudeAccount = true
	if !m.useClaudeAccount() {
		t.Fatal("ClaudeAccount aktif + provider anthropic seharusnya memakai akun Claude")
	}
	m.provider = "openai"
	if m.useClaudeAccount() {
		t.Fatal("provider lain tidak boleh memakai akun Claude")
	}
}

func TestClaudePromptAndArgs(t *testing.T) {
	p := claudePrompt([]historyItem{{"user", "hai"}, {"assistant", "halo"}}, "apa kabar")
	for _, want := range []string{"Pengguna: hai", "Asisten: halo", "apa kabar"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt tidak memuat %q:\n%s", want, p)
		}
	}
	args := strings.Join(claudeArgs(chatRequest{Mode: "coding", Model: "anthropic/claude-sonnet-4-5"}), " ")
	for _, want := range []string{"-p", "stream-json", "--model claude-sonnet-4-5-20250929", "expert coding assistant"} {
		if !strings.Contains(args, want) {
			t.Errorf("args tidak memuat %q: %s", want, args)
		}
	}
}

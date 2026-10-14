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

func TestSelectProviderSwitchesAPI(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := openConnect(t, newScreen(t, screenChat, loggedIn(""), ""))

	m, _ = update(m, key(tea.KeyDown)) // kursor ke Anthropic
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal != nil {
		t.Fatal("modal seharusnya tertutup setelah Enter")
	}
	if m.provider != "anthropic" {
		t.Fatalf("provider = %q, seharusnya anthropic", m.provider)
	}
	if m.modelID == "" || parseProvider(m.provider) == nil {
		t.Fatalf("model bawaan provider seharusnya terisi, got %q", m.modelID)
	}
	view := plain(m)
	if !strings.Contains(view, "Anthropic") {
		t.Errorf("seharusnya ada catatan ganti provider:\n%s", view)
	}

	if loadProvider() != "anthropic" || loadModel() != m.modelID {
		t.Error("pilihan provider dan model seharusnya tersimpan")
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

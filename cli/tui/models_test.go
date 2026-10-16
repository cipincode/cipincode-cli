package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// modelsServer adalah API /models palsu untuk pengujian fetch live.
func modelsServer(t *testing.T, ids []string, status int) string {
	t.Helper()
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models" && r.URL.Path != "/v1/models" {
			t.Errorf("path salah: %s", r.URL.Path)
		}
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		data := make([]map[string]string, 0, len(ids))
		for _, id := range ids {
			data = append(data, map[string]string{"id": id, "display_name": "Nama " + id})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
	})
	return srv.URL
}

func TestFetchLiveFiltersNonChat(t *testing.T) {
	url := modelsServer(t, []string{
		"gpt-4o", "gpt-4o-mini", "whisper-1", "text-embedding-3-small",
		"dall-e-3", "o4-mini", "gpt-4o-realtime-preview",
	}, http.StatusOK)

	models, err := fetchModelsFrom("openai", "sk-test", url)
	if err != nil {
		t.Fatalf("fetchModelsFrom: %v", err)
	}
	var ids []string
	for _, mdl := range models {
		ids = append(ids, mdl.id)
	}
	for _, want := range []string{"gpt-4o", "gpt-4o-mini", "o4-mini"} {
		if !strings.Contains(strings.Join(ids, ","), want) {
			t.Errorf("seharusnya memuat %q, got %v", want, ids)
		}
	}
	for _, junk := range []string{"whisper-1", "text-embedding-3-small", "dall-e-3", "gpt-4o-realtime-preview"} {
		if strings.Contains(strings.Join(ids, ","), junk) {
			t.Errorf("seharusnya menyaring %q, got %v", junk, ids)
		}
	}
	// Nama ramah dari katalog untuk ID yang dikenal.
	for _, mdl := range models {
		if mdl.id == "gpt-4o" && mdl.name != "GPT-4o" {
			t.Errorf("nama ramah gpt-4o = %q", mdl.name)
		}
	}
}

func TestFetchLiveAnthropic(t *testing.T) {
	url := modelsServer(t, []string{"claude-sonnet-4-5-20250929", "buatan-sendiri"}, http.StatusOK)

	models, err := fetchModelsFrom("anthropic", "sk-ant-test", url)
	if err != nil {
		t.Fatalf("fetchModelsFrom: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("seharusnya 2 model, got %v", models)
	}
	// ID langsung yang dikenal memakai nama katalog.
	for _, mdl := range models {
		if mdl.id == "claude-sonnet-4-5-20250929" && mdl.name != "Claude Sonnet 4.5" {
			t.Errorf("nama katalog tidak dipakai, got %q", mdl.name)
		}
	}
}

func TestFetchLiveRejectsBadKey(t *testing.T) {
	url := modelsServer(t, nil, http.StatusUnauthorized)
	if _, err := fetchModelsFrom("openai", "sk-salah", url); err == nil {
		t.Error("kunci salah seharusnya gagal")
	} else if !strings.Contains(err.Error(), "ditolak") {
		t.Errorf("pesan seharusnya menyebut kunci ditolak, got %q", err)
	}
}

func TestModelLiveOpensFromAccount(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	url := modelsServer(t, []string{"gpt-4o", "gpt-4o-mini"}, http.StatusOK)

	m := newScreen(t, screenChat, loggedIn(""), "")
	m.provider = "openai"
	if m.config != nil {
		if m.config.ProviderKeys == nil {
			m.config.ProviderKeys = map[string]string{}
		}
		m.config.ProviderKeys["openai"] = "sk-test"
	}
	if !m.modelsLiveOK() {
		t.Fatal("ada kunci sendiri seharusnya bisa daftar live")
	}

	// Simulasi hasil fetch: ganti baseURL lewat fetchModelsFrom langsung.
	models, err := fetchModelsFrom("openai", "sk-test", url)
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	m, _ = update(m, modelsMsg{models: models})
	if m.modal == nil || m.modal.kind != modalModels {
		t.Fatal("seharusnya membuka modal model")
	}
	if !m.modelsLive || !strings.Contains(plain(m), "langsung dari akun") {
		t.Errorf("modal seharusnya menandai daftar live:\n%s", plain(m))
	}
}

func TestModelFallsBackToCatalog(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.provider = "openai"

	m, _ = update(m, modelsMsg{err: errInvalidKey})
	if m.modal == nil || m.modal.kind != modalModels {
		t.Fatal("gagal fetch seharusnya tetap membuka modal katalog")
	}
	if m.modelsLive {
		t.Error("fallback tidak boleh menandai live")
	}
	if !strings.Contains(plain(m), "GPT-4o") {
		t.Errorf("modal seharusnya memuat katalog bawaan:\n%s", plain(m))
	}

	// Via server (tanpa kunci) langsung memakai katalog tanpa fetch.
	m2 := newScreen(t, screenChat, loggedIn(""), "")
	m2.provider = "openai"
	if m2.modelsLiveOK() {
		t.Error("tanpa kunci seharusnya tidak fetch live")
	}
	m2, _ = update(m2, typed("/model"))
	m2, _ = update(m2, key(tea.KeyEnter))
	if m2.modal == nil || m2.modal.kind != modalModels {
		t.Fatal("/model via server seharusnya membuka katalog")
	}
}

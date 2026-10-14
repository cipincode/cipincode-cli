package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type seenRequest struct {
	Token string
	URL   string
	Key   string
	Body  struct {
		Message string        `json:"message"`
		Mode    string        `json:"mode"`
		Effort  string        `json:"effort"`
		History []historyItem `json:"history"`
	}
}

// fakeAgent meniru server agent: merekam permintaan dan mengalirkan peristiwa yang diberikan.
type fakeAgent struct {
	*httptest.Server
	mu   sync.Mutex
	seen []seenRequest
}

func newFakeAgent(t *testing.T, events ...string) *fakeAgent {
	t.Helper()
	fa := &fakeAgent{}
	fa.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req seenRequest
		req.Token = r.Header.Get("X-Agent-Token")
		req.URL = r.Header.Get("X-Cipincode-Url")
		req.Key = r.Header.Get("X-Cipincode-Key")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &req.Body)

		fa.mu.Lock()
		fa.seen = append(fa.seen, req)
		fa.mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		for _, ev := range events {
			_, _ = w.Write([]byte("data: " + ev + "\n\n"))
		}
	}))
	t.Cleanup(fa.Close)
	return fa
}

func (fa *fakeAgent) requests() []seenRequest {
	fa.mu.Lock()
	defer fa.mu.Unlock()
	return append([]seenRequest(nil), fa.seen...)
}

const (
	evUsage = `{"type":"usage","session":{"used":1,"limit":40,"percent":3,"exhausted":false,"resets_at":null},"weekly":{"used":1,"limit":400,"percent":0,"exhausted":false,"resets_at":null}}`
	evHalo  = `{"type":"token","text":"Halo"}`
	evDunia = `{"type":"token","text":" dunia"}`
	evDone  = `{"type":"done"}`
)

// chatWith membuat model chat yang terhubung ke fakeAgent.
func chatWith(t *testing.T, fa *fakeAgent, cfg *Config) model {
	t.Helper()
	m := newScreen(t, screenChat, cfg, "")
	m.agent = &Agent{URL: fa.URL, Token: "rahasia"}
	return m
}

// drain menjalankan perintah hingga aliran selesai dan mengembalikan model akhir.
func drain(t *testing.T, m model, cmd tea.Cmd) model {
	t.Helper()
	for i := 0; cmd != nil && i < 500; i++ {
		msg := cmd()
		if msg == nil {
			break
		}
		m, cmd = update(m, msg)
	}
	return m
}

func send(t *testing.T, m model, text string) model {
	t.Helper()
	m, _ = update(m, typed(text))
	m, cmd := update(m, key(tea.KeyEnter))
	return drain(t, m, cmd)
}

func loggedInCfg() *Config {
	return &Config{Name: "Budi Santoso", Email: "budi@example.com", Key: "cpc_ok", URL: "http://dash.test"}
}

func TestChatStreamsTheAnswerIntoTheConversation(t *testing.T) {
	fa := newFakeAgent(t, evUsage, evHalo, evDunia, evDone)
	m := send(t, chatWith(t, fa, loggedInCfg()), "halo")

	view := plain(m)
	if !strings.Contains(view, "Halo dunia") || !strings.Contains(view, "❯ halo") {
		t.Fatalf("percakapan tidak tampil:\n%s", view)
	}
	if m.streaming || m.cancel != nil {
		t.Fatal("setelah selesai tidak boleh lagi streaming")
	}
}

func TestChatSendsCredentialsTokenModeAndMessage(t *testing.T) {
	fa := newFakeAgent(t, evHalo, evDone)
	m := chatWith(t, fa, loggedInCfg())
	m, _ = update(m, key(tea.KeyShiftTab)) // mode coding
	send(t, m, "buat fungsi")

	reqs := fa.requests()
	if len(reqs) != 1 {
		t.Fatalf("requests=%d", len(reqs))
	}
	r := reqs[0]
	if r.Token != "rahasia" || r.URL != "http://dash.test" || r.Key != "cpc_ok" {
		t.Fatalf("header salah: %+v", r)
	}
	if r.Body.Message != "buat fungsi" || r.Body.Mode != "coding" {
		t.Fatalf("body salah: %+v", r.Body)
	}
}

func TestChatCarriesHistoryToTheNextTurn(t *testing.T) {
	fa := newFakeAgent(t, evHalo, evDunia, evDone)
	m := chatWith(t, fa, loggedInCfg())
	m = send(t, m, "pertanyaan satu")
	send(t, m, "pertanyaan dua")

	reqs := fa.requests()
	if len(reqs) != 2 {
		t.Fatalf("requests=%d", len(reqs))
	}
	if len(reqs[0].Body.History) != 0 {
		t.Fatalf("giliran pertama tanpa riwayat, got %+v", reqs[0].Body.History)
	}
	want := []historyItem{{"user", "pertanyaan satu"}, {"assistant", "Halo dunia"}}
	got := reqs[1].Body.History
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("riwayat giliran kedua salah: %+v", got)
	}
	if reqs[1].Body.Message != "pertanyaan dua" {
		t.Fatalf("message=%q", reqs[1].Body.Message)
	}
}

func TestHistoryIsCappedAndSkipsNoise(t *testing.T) {
	m := newScreen(t, screenChat, loggedInCfg(), "")
	for i := 0; i < 30; i++ {
		m.messages = append(m.messages, message{role: roleUser, content: "u"}, message{role: roleAssistant, content: "a"})
	}
	m.messages = append(m.messages, message{role: roleError, content: "galat"}, message{role: roleSystem, content: "info"}, message{role: roleRich, content: "blok"}, message{role: roleAssistant, content: ""})

	items := m.history()
	if len(items) != historyLimit {
		t.Fatalf("len=%d, want %d", len(items), historyLimit)
	}
	for _, it := range items {
		if it.Content == "" || (it.Role != "user" && it.Role != "assistant") {
			t.Fatalf("riwayat memuat entri tak wajar: %+v", it)
		}
	}
}

func TestUsageEventUpdatesTheHeaderDuringTheAnswer(t *testing.T) {
	fa := newFakeAgent(t, evUsage, evHalo, evDone)
	m := send(t, chatWith(t, fa, loggedInCfg()), "halo")

	if m.usage == nil || m.usage.Session.Used != 1 {
		t.Fatalf("pemakaian dari aliran seharusnya tersimpan: %+v", m.usage)
	}
	if !strings.Contains(plain(m), "Sesi 5 jam") {
		t.Fatal("header seharusnya menampilkan ringkasan sesi")
	}
}

func TestBlockedAnswerRollsBackAndShowsUsage(t *testing.T) {
	blocked := `{"type":"blocked","message":"Batas pemakaian sesi 5 jam tercapai.","blocked":"session",` +
		`"session":{"used":40,"limit":40,"percent":100,"exhausted":true,"resets_at":"` + time.Now().Add(2*time.Hour).UTC().Format(time.RFC3339) + `"},` +
		`"weekly":{"used":40,"limit":400,"percent":10,"exhausted":false,"resets_at":null}}`
	fa := newFakeAgent(t, blocked)
	m := send(t, chatWith(t, fa, loggedInCfg()), "tolong buatkan fungsi")

	if m.streaming || m.cancel != nil {
		t.Fatal("tidak boleh ada jawaban berjalan")
	}
	if len(m.messages) != 2 { // galat dan blok pemakaian saja
		t.Fatalf("pertanyaan dan jawaban kosong seharusnya dihapus, messages=%d", len(m.messages))
	}
	if m.input.Value() != "tolong buatkan fungsi" {
		t.Fatalf("teks seharusnya kembali ke input, got %q", m.input.Value())
	}

	view := plain(m)
	for _, want := range []string{"Batas pemakaian sesi 5 jam tercapai", "Pesan Anda tidak dikirim", "100%", "Habis · reset dalam"} {
		if !strings.Contains(view, want) {
			t.Errorf("tidak memuat %q:\n%s", want, view)
		}
	}
	if m.usage == nil || !m.usage.Session.Exhausted {
		t.Fatal("snapshot habis seharusnya tersimpan")
	}
}

func TestErrorEventShowsMessageAndDropsEmptyAnswer(t *testing.T) {
	fa := newFakeAgent(t, `{"type":"error","message":"Layanan AI belum dikonfigurasi."}`)
	m := send(t, chatWith(t, fa, loggedInCfg()), "halo")

	if m.streaming {
		t.Fatal("galat seharusnya menghentikan streaming")
	}
	view := plain(m)
	if !strings.Contains(view, "Layanan AI belum dikonfigurasi") {
		t.Fatalf("pesan galat tidak tampil:\n%s", view)
	}
	for _, msg := range m.messages {
		if msg.role == roleAssistant && msg.content == "" {
			t.Fatal("jawaban kosong seharusnya dihapus")
		}
	}
	if !strings.Contains(view, "❯ halo") {
		t.Fatal("pertanyaan pengguna seharusnya tetap terlihat")
	}
}

func TestErrorAfterPartialAnswerKeepsWhatWasReceived(t *testing.T) {
	fa := newFakeAgent(t, evHalo, `{"type":"error","message":"Provider terputus"}`)
	m := send(t, chatWith(t, fa, loggedInCfg()), "halo")

	view := plain(m)
	if !strings.Contains(view, "Halo") || !strings.Contains(view, "Provider terputus") {
		t.Fatalf("jawaban parsial dan galat seharusnya tampil:\n%s", view)
	}
}

func TestChatNeedsLoginAndKeepsTheText(t *testing.T) {
	fa := newFakeAgent(t, evHalo, evDone)
	m := chatWith(t, fa, nil)

	m, _ = update(m, typed("halo"))
	m, cmd := update(m, key(tea.KeyEnter))

	if cmd != nil || m.streaming {
		t.Fatal("tanpa login tidak boleh ada permintaan")
	}
	if !strings.Contains(plain(m), "Masuk dulu dengan /login") {
		t.Fatalf("seharusnya meminta login:\n%s", plain(m))
	}
	if m.input.Value() != "halo" {
		t.Fatalf("teks seharusnya tetap di input, got %q", m.input.Value())
	}
	if len(fa.requests()) != 0 {
		t.Fatal("agent tidak boleh dipanggil tanpa login")
	}
}

func TestUnreachableAgentShowsError(t *testing.T) {
	m := newScreen(t, screenChat, loggedInCfg(), "")
	m.agent = &Agent{URL: "http://127.0.0.1:1"}

	m = send(t, m, "halo")

	if m.streaming || !strings.Contains(plain(m), "Error:") {
		t.Fatalf("seharusnya menampilkan galat:\n%s", plain(m))
	}
}

func TestAgentRejectsWrongTokenSoCLIReportsIt(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Agent-Token") != "benar" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		_, _ = w.Write([]byte("data: " + evDone + "\n\n"))
	}))
	t.Cleanup(srv.Close)

	out := make(chan chatEvent, 4)
	err := (&Agent{URL: srv.URL, Token: "salah"}).Chat(t.Context(), chatRequest{Message: "x"}, out)
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("got %v", err)
	}

	if err := (&Agent{URL: srv.URL, Token: "benar"}).Chat(t.Context(), chatRequest{Message: "x"}, out); err != nil {
		t.Fatalf("token benar seharusnya diterima: %v", err)
	}
}

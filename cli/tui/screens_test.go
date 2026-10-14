package main

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func key(t tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: t} }

func typed(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func update(m model, msg tea.Msg) (model, tea.Cmd) {
	next, cmd := m.Update(msg)
	return next.(model), cmd
}

func newScreen(t *testing.T, s screen, cfg *Config, url string) model {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("CIPINCODE_NO_BROWSER", "1")

	m := newModel(&Agent{URL: "http://127.0.0.1:4096"}, cfg, time.Time{}, options{screen: s, url: url})
	m, _ = update(m, tea.WindowSizeMsg{Width: 100, Height: 30})
	return m
}

func TestOnboardingAnimatesThenWaitsForEnter(t *testing.T) {
	m := newScreen(t, screenOnboarding, nil, "")

	// animasi mulai kosong lalu logo muncul bertahap
	if strings.Contains(plain(m), "Selamat datang") {
		t.Fatal("judul seharusnya belum tampil di frame 0")
	}
	for i := 0; i < 12; i++ {
		m, _ = update(m, animTickMsg{})
	}
	if !strings.Contains(plain(m), "▄█") {
		t.Fatal("logo seharusnya sudah tampil setelah 12 frame")
	}

	// Enter sebelum selesai hanya melewati animasi, tidak pindah layar
	m, _ = update(m, key(tea.KeyEnter))
	if m.screen != screenOnboarding || m.frame != onboardFrame {
		t.Fatalf("screen=%v frame=%d", m.screen, m.frame)
	}
	view := plain(m)
	for _, want := range []string{"Selamat datang di CipinCode", "AI agent untuk terminal Anda", "Tekan Enter"} {
		if !strings.Contains(view, want) {
			t.Errorf("onboarding tidak memuat %q", want)
		}
	}

	m, _ = update(m, key(tea.KeyEnter))
	if m.screen != screenMethod {
		t.Fatalf("setelah selesai Enter seharusnya ke pilihan login, screen=%v", m.screen)
	}
}

func TestOnboardingTitleIsTypedProgressively(t *testing.T) {
	m := newScreen(t, screenOnboarding, nil, "")

	m.frame = 22
	if v := plain(m); !strings.Contains(v, "Selamat") || strings.Contains(v, "CipinCode") {
		t.Fatalf("judul seharusnya baru separuh di frame 20:\n%s", v)
	}
}

func TestOnboardingCanBeSkippedToChat(t *testing.T) {
	m := newScreen(t, screenOnboarding, nil, "")

	m, _ = update(m, key(tea.KeyEsc))
	if m.screen != screenChat {
		t.Fatalf("Esc seharusnya melewati ke chat, screen=%v", m.screen)
	}
}

func TestMethodScreenOffersEmailAndGoogle(t *testing.T) {
	m := newScreen(t, screenMethod, nil, "")
	view := plain(m)
	for _, want := range []string{"Pilih cara masuk", "Email", "Google"} {
		if !strings.Contains(view, want) {
			t.Errorf("layar pilihan tidak memuat %q", want)
		}
	}

	m, _ = update(m, key(tea.KeyDown))
	if m.methodIndex != 1 {
		t.Fatalf("methodIndex=%d", m.methodIndex)
	}
}

func TestCodeScreenShowsURLForChosenMethod(t *testing.T) {
	m := newScreen(t, screenMethod, nil, "http://dash.test")

	m, _ = update(m, key(tea.KeyEnter)) // email
	if m.screen != screenCode || !strings.Contains(m.authURL, "http://dash.test/cli/authorize?challenge=") || strings.Contains(m.authURL, "method=google") {
		t.Fatalf("url email: %q", m.authURL)
	}
	if !strings.Contains(plain(m), "Tempel kode") {
		t.Fatal("layar kode harus meminta kode ditempel")
	}

	m, _ = update(m, key(tea.KeyEsc))
	m, _ = update(m, key(tea.KeyDown))
	m, _ = update(m, key(tea.KeyEnter)) // google
	if !strings.Contains(m.authURL, "method=google") {
		t.Fatalf("url google: %q", m.authURL)
	}
	if !strings.Contains(plain(m), "Masuk dengan Google") {
		t.Fatal("judul layar seharusnya menyebut Google")
	}
}

func TestPastingCodeLogsInAndEntersChat(t *testing.T) {
	srv := exchangeServer(t, 200, map[string]any{
		"key":  "cpc_baru",
		"user": map[string]string{"name": "Budi Santoso", "email": "budi@example.com"},
	})
	m := newScreen(t, screenMethod, nil, srv.URL)
	m, _ = update(m, key(tea.KeyEnter))

	m, _ = update(m, typed("ABCD-EFGH"))
	m, cmd := update(m, key(tea.KeyEnter))
	if !m.exchanging || cmd == nil {
		t.Fatal("Enter seharusnya memulai penukaran kode")
	}

	m, cmd = update(m, cmd()) // hasil penukaran
	if !m.loginOK || m.config == nil || m.config.Name != "Budi Santoso" {
		t.Fatalf("login belum berhasil: %+v", m)
	}
	if !strings.Contains(plain(m), "Berhasil masuk sebagai Budi Santoso") {
		t.Fatal("pesan sukses tidak tampil")
	}
	saved, _ := loadConfig()
	if saved == nil || saved.Key != "cpc_baru" || saved.URL != srv.URL {
		t.Fatalf("config tersimpan salah: %+v", saved)
	}

	m, _ = update(m, loginDoneMsg{})
	if m.screen != screenChat || !strings.Contains(plain(m), "Selamat datang, Budi") {
		t.Fatalf("setelah login seharusnya masuk chat, screen=%v", m.screen)
	}
}

func TestWrongCodeShowsErrorAndAllowsRetry(t *testing.T) {
	srv := exchangeServer(t, 422, map[string]string{"message": "Kode tidak valid atau sudah kedaluwarsa."})
	m := newScreen(t, screenMethod, nil, srv.URL)
	m, _ = update(m, key(tea.KeyEnter))

	m, _ = update(m, typed("SALAH"))
	m, cmd := update(m, key(tea.KeyEnter))
	m, _ = update(m, cmd())

	if m.loginOK || m.exchanging {
		t.Fatal("login tidak boleh berhasil dengan kode salah")
	}
	if !strings.Contains(plain(m), "Kode tidak valid") {
		t.Fatalf("pesan error tidak tampil:\n%s", plain(m))
	}
	if saved, _ := loadConfig(); saved != nil {
		t.Fatal("config tidak boleh tersimpan")
	}
	if m.codeInput.Value() != "SALAH" {
		t.Fatal("kode yang diketik seharusnya tetap ada agar bisa diperbaiki")
	}
}

func TestLoginOnlyQuitsAfterSuccess(t *testing.T) {
	m := newScreen(t, screenCode, &Config{}, "")
	m.loginOnly = true

	_, cmd := update(m, loginDoneMsg{})
	if cmd == nil {
		t.Fatal("seharusnya keluar")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("perintah yang dikembalikan seharusnya tea.Quit")
	}
}

// ---- menu "/" ----

func chatModel(t *testing.T) model {
	t.Helper()
	return newScreen(t, screenChat, &Config{Name: "Budi", Email: "budi@example.com"}, "http://dash.test")
}

func TestSlashOpensMenuWithManyOptions(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/"))

	view := plain(m)
	for _, want := range []string{"/help", "/coding", "/status", "↑↓ pilih", "perintah lagi"} {
		if !strings.Contains(view, want) {
			t.Errorf("menu tidak memuat %q", want)
		}
	}
	if len(m.menuItems()) != len(commands) {
		t.Fatalf("menu memuat %d, want %d", len(m.menuItems()), len(commands))
	}
	if strings.Contains(view, "/exit") {
		t.Error("perintah terakhir seharusnya belum terlihat sebelum menggulir")
	}

	// menggulir ke bawah (Up dari atas melompat ke item terakhir) memperlihatkan sisanya
	m, _ = update(m, key(tea.KeyUp))
	if view := plain(m); !strings.Contains(view, "/exit") || !strings.Contains(view, "/login") {
		t.Errorf("setelah menggulir, menu seharusnya memuat /login dan /exit:\n%s", view)
	}
}

func TestMenuFiltersWhileTypingAndClosesOnSpace(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/cod"))

	items := m.menuItems()
	if len(items) != 1 || items[0].name != "coding" {
		t.Fatalf("items=%v", items)
	}
	if strings.Contains(plain(m), "/status") {
		t.Fatal("menu seharusnya tersaring")
	}

	m, _ = update(m, typed(" "))
	if len(m.menuItems()) != 0 {
		t.Fatal("menu seharusnya tertutup setelah spasi")
	}
}

func TestMenuNavigationRunsSelectedCommand(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/"))
	m, _ = update(m, key(tea.KeyDown)) // /chat
	m, _ = update(m, key(tea.KeyDown)) // /coding
	m, _ = update(m, key(tea.KeyEnter))

	if m.mode != modeCoding {
		t.Fatalf("mode=%v, seharusnya coding", m.mode)
	}
	if m.input.Value() != "" || len(m.menuItems()) != 0 {
		t.Fatal("input dan menu seharusnya bersih setelah perintah dijalankan")
	}
	if !strings.Contains(plain(m), "Mode Coding") {
		t.Fatal("konfirmasi pindah mode tidak tampil")
	}
}

func TestMenuUpWrapsAroundAndEscCloses(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/"))
	m, _ = update(m, key(tea.KeyUp))
	if m.menuIndex != len(commands)-1 {
		t.Fatalf("menuIndex=%d", m.menuIndex)
	}

	m, _ = update(m, key(tea.KeyEsc))
	if m.input.Value() != "" || len(m.menuItems()) != 0 {
		t.Fatal("Esc seharusnya menutup menu dan mengosongkan input")
	}
}

func TestTabCompletesTheSelectedCommand(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/sta"))
	m, _ = update(m, key(tea.KeyTab))

	if m.input.Value() != "/status " {
		t.Fatalf("input=%q", m.input.Value())
	}
}

func TestTypedCommandRunsOnEnter(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/version"))
	m, _ = update(m, key(tea.KeyEnter)) // menu terbuka: menjalankan item terpilih

	if !strings.Contains(plain(m), "cipincode "+version) {
		t.Fatalf("keluaran /version tidak tampil:\n%s", plain(m))
	}
}

func TestStatusAndWhoamiReportAccount(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/whoami"))
	m, _ = update(m, key(tea.KeyEnter))
	if !strings.Contains(plain(m), "Budi (budi@example.com)") {
		t.Fatalf("whoami:\n%s", plain(m))
	}

	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))
	view := plain(m)
	for _, want := range []string{"Akun", "Dashboard   http://dash.test", "Agent", "Mode"} {
		if !strings.Contains(view, want) {
			t.Errorf("status tidak memuat %q", want)
		}
	}
}

func TestLoginAndLogoutCommands(t *testing.T) {
	m := chatModel(t)
	if err := saveConfig(*m.config); err != nil {
		t.Fatal(err)
	}

	m, _ = update(m, typed("/logout"))
	m, _ = update(m, key(tea.KeyEnter))
	if m.config != nil {
		t.Fatal("config seharusnya kosong setelah /logout")
	}
	if saved, _ := loadConfig(); saved != nil {
		t.Fatal("file config seharusnya terhapus")
	}
	if !strings.Contains(plain(m), "Belum login") {
		t.Fatalf("tampilan seharusnya menandai belum login:\n%s", plain(m))
	}

	m, _ = update(m, typed("/login"))
	m, _ = update(m, key(tea.KeyEnter))
	if m.screen != screenMethod {
		t.Fatalf("/login seharusnya membuka pilihan cara masuk, screen=%v", m.screen)
	}
}

func TestNewAndClearCommands(t *testing.T) {
	m := chatModel(t)
	m.messages = []message{{role: roleUser, content: "halo"}, {role: roleAssistant, content: "hai"}}
	m.home = false

	m, _ = update(m, typed("/new"))
	m, _ = update(m, key(tea.KeyEnter))
	if len(m.messages) != 0 || !m.home {
		t.Fatalf("/new: messages=%d home=%v", len(m.messages), m.home)
	}

	m.messages = []message{{role: roleUser, content: "halo"}}
	m, _ = update(m, typed("/clear"))
	m, _ = update(m, key(tea.KeyEnter))
	if len(m.messages) != 0 || m.home {
		t.Fatalf("/clear: messages=%d home=%v", len(m.messages), m.home)
	}
}

func TestUnknownCommandShowsError(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/xyz "))
	m, _ = update(m, key(tea.KeyEnter))

	if !strings.Contains(plain(m), "Perintah tidak dikenal: /xyz") {
		t.Fatalf("pesan error tidak tampil:\n%s", plain(m))
	}
}

func TestMenuShrinksTheViewportAndRestoresIt(t *testing.T) {
	m := chatModel(t)
	base := m.viewport.Height

	m, _ = update(m, typed("/"))
	if m.viewport.Height >= base {
		t.Fatalf("viewport seharusnya mengecil saat menu terbuka: %d >= %d", m.viewport.Height, base)
	}

	m, _ = update(m, key(tea.KeyEsc))
	if m.viewport.Height != base {
		t.Fatalf("viewport seharusnya kembali: %d != %d", m.viewport.Height, base)
	}
}

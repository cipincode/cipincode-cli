package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func usageModal(t *testing.T, cfgURL string) model {
	t.Helper()
	m := newScreen(t, screenChat, loggedIn(cfgURL), "")
	m, _ = update(m, typed("/usage"))
	m, cmd := update(m, key(tea.KeyEnter))
	m, _ = update(m, cmd())
	return m
}

func okUsageServer(t *testing.T) string {
	reset := time.Now().Add(3 * time.Hour)
	return usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(snapshotJSON(
			usageMeter{Used: 15, Limit: 40, Percent: 38, ResetsAt: &reset},
			usageMeter{Used: 48, Limit: 400, Percent: 12},
		))
	}).URL
}

func TestUsageOpensAsAModalThatNeedsEscToClose(t *testing.T) {
	m := usageModal(t, okUsageServer(t))

	if m.modal == nil {
		t.Fatal("/usage seharusnya membuka modal")
	}
	view := plain(m)
	for _, want := range []string{"Pemakaian · Budi Santoso", "Sesi 5 jam", "Mingguan", "Esc untuk menutup"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal tidak memuat %q:\n%s", want, view)
		}
	}
	// modal ditempel di atas chat: header chat tetap terlihat di sekelilingnya
	if !strings.Contains(view, "Naufal") && !strings.Contains(view, "Budi Santoso") {
		t.Error("tampilan chat di belakang modal seharusnya tetap ada")
	}

	m, _ = update(m, key(tea.KeyEsc))
	if m.modal != nil {
		t.Fatal("Esc seharusnya menutup modal")
	}
	if strings.Contains(plain(m), "Esc untuk menutup") || strings.Contains(plain(m), "Mingguan") {
		t.Fatal("setelah ditutup, isi modal tidak boleh tersisa di layar")
	}
	// dan tidak meninggalkan catatan di percakapan
	for _, msg := range m.messages {
		if msg.role == roleRich {
			t.Fatal("hasil /usage tidak boleh ditambahkan ke percakapan")
		}
	}
}

func TestModalIgnoresTypingUntilClosed(t *testing.T) {
	m := usageModal(t, okUsageServer(t))

	m, _ = update(m, typed("halo"))
	m, _ = update(m, typed("/"))
	if m.input.Value() != "" {
		t.Fatalf("input tidak boleh berubah saat modal terbuka, got %q", m.input.Value())
	}
	if len(m.menuItems()) != 0 {
		t.Fatal("menu / tidak boleh terbuka di bawah modal")
	}
	if m.modal == nil {
		t.Fatal("modal harus tetap terbuka sampai Esc")
	}

	m, _ = update(m, key(tea.KeyEsc))
	m, _ = update(m, typed("halo"))
	if m.input.Value() != "halo" {
		t.Fatalf("setelah modal ditutup input harus aktif lagi, got %q", m.input.Value())
	}
}

func TestEnterAndQAlsoCloseTheModal(t *testing.T) {
	for _, closer := range []tea.KeyMsg{key(tea.KeyEnter), typed("q")} {
		m := usageModal(t, okUsageServer(t))
		m, _ = update(m, closer)
		if m.modal != nil {
			t.Errorf("%v seharusnya menutup modal", closer)
		}
	}
}

func TestCtrlCStillQuitsFromAModal(t *testing.T) {
	m := usageModal(t, okUsageServer(t))

	_, cmd := update(m, key(tea.KeyCtrlC))
	if cmd == nil {
		t.Fatal("Ctrl+C seharusnya keluar")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("perintah seharusnya tea.Quit")
	}
}

func TestStatusAndHelpAreModalsToo(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://dash.test"), "")

	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))
	if m.modal == nil || !strings.Contains(plain(m), "Dashboard   http://dash.test") || !strings.Contains(plain(m), "Esc untuk menutup") {
		t.Fatalf("/status seharusnya modal:\n%s", plain(m))
	}
	m, _ = update(m, key(tea.KeyEsc))

	m, _ = update(m, typed("/help"))
	m, _ = update(m, key(tea.KeyEnter))
	view := plain(m)
	if m.modal == nil || !strings.Contains(view, "Bantuan") || !strings.Contains(view, "/usage") {
		t.Fatalf("/help seharusnya modal berisi daftar perintah:\n%s", view)
	}
	m, _ = update(m, key(tea.KeyEsc))
	if m.modal != nil {
		t.Fatal("Esc seharusnya menutup")
	}
}

func TestShortInfoCommandsStayInline(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	m, _ = update(m, typed("/version"))
	m, _ = update(m, key(tea.KeyEnter))

	if m.modal != nil {
		t.Fatal("/version cukup satu baris, bukan modal")
	}
}

func TestUsageErrorsStayInlineAndOpenNoModal(t *testing.T) {
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	m := usageModal(t, srv.URL)

	if m.modal != nil {
		t.Fatal("galat tidak boleh membuka modal kosong")
	}
	if !strings.Contains(plain(m), "/login") {
		t.Fatalf("galat seharusnya tampil di percakapan:\n%s", plain(m))
	}
}

func TestModalDoesNotShiftTheChatBehindIt(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://dash.test"), "")
	before := strings.Split(plain(m), "\n")

	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))
	after := strings.Split(plain(m), "\n")

	if len(after) != len(before) {
		t.Fatalf("jumlah baris layar berubah: %d -> %d", len(before), len(after))
	}
	if after[0] != before[0] || after[len(after)-1] != before[len(before)-1] {
		t.Error("baris paling atas dan bawah (di luar modal) seharusnya tidak berubah")
	}
}

func TestLongModalScrollsAndShowsPosition(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	m, _ = update(m, tea.WindowSizeMsg{Width: 100, Height: 16})

	var body []string
	for i := 1; i <= 30; i++ {
		body = append(body, fmt.Sprintf("baris %02d", i))
	}
	m = m.openModal("Panjang", strings.Join(body, "\n"))

	view := plain(m)
	if !strings.Contains(view, "baris 01") || strings.Contains(view, "baris 30") || !strings.Contains(view, "gulir") {
		t.Fatalf("awal: harus tampil dari atas dengan petunjuk gulir:\n%s", view)
	}

	m, _ = update(m, key(tea.KeyPgDown))
	m, _ = update(m, key(tea.KeyDown))
	if m.modal.scroll == 0 || strings.Contains(plain(m), "baris 01") {
		t.Fatalf("setelah PgDn/Down seharusnya tergulir (scroll=%d)", m.modal.scroll)
	}

	for i := 0; i < 20; i++ {
		m, _ = update(m, key(tea.KeyPgDown))
	}
	if !strings.Contains(plain(m), "baris 30") {
		t.Fatalf("di dasar harus tampil baris terakhir:\n%s", plain(m))
	}
	for i := 0; i < 20; i++ {
		m, _ = update(m, key(tea.KeyPgUp))
	}
	if m.modal.scroll != 0 {
		t.Fatalf("scroll tidak boleh negatif, got %d", m.modal.scroll)
	}
}

func TestOverlayPlacesTheBoxCenteredAndKeepsTheRest(t *testing.T) {
	base := strings.Join([]string{
		strings.Repeat("a", 20),
		strings.Repeat("b", 20),
		strings.Repeat("c", 20),
		strings.Repeat("d", 20),
		strings.Repeat("e", 20),
	}, "\n")

	got := ansi.Strip(overlay(base, "XX\nXX", 20))
	lines := strings.Split(got, "\n")

	if lines[0] != strings.Repeat("a", 20) || lines[4] != strings.Repeat("e", 20) {
		t.Fatal("baris di luar box harus utuh")
	}
	if lines[1] != strings.Repeat("b", 9)+"XX"+strings.Repeat("b", 9) {
		t.Errorf("baris ke-2: %q", lines[1])
	}
	if lines[2] != strings.Repeat("c", 9)+"XX"+strings.Repeat("c", 9) {
		t.Errorf("baris ke-3: %q", lines[2])
	}
	for i, l := range lines {
		if ansi.StringWidth(l) != 20 {
			t.Errorf("lebar baris %d berubah: %d", i, ansi.StringWidth(l))
		}
	}
}

func TestModalClipsWideLinesInsteadOfBreakingLayout(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	m, _ = update(m, tea.WindowSizeMsg{Width: 60, Height: 30})

	widest := func(view string) int {
		w := 0
		for _, line := range strings.Split(view, "\n") {
			w = max(w, ansi.StringWidth(line))
		}
		return w
	}
	before := widest(plain(m))

	m = m.openModal("Lebar", strings.Repeat("x", 200))
	if after := widest(plain(m)); after > before {
		t.Fatalf("modal melebarkan layar: %d > %d", after, before)
	}
	if !strings.Contains(plain(m), "…") {
		t.Error("baris yang terlalu panjang seharusnya dipotong dengan …")
	}
}

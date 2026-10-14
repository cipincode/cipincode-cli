package main

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

func TestFilledSegments(t *testing.T) {
	cases := []struct{ used, limit, want int }{
		{0, 40, 0},   // belum dipakai: tidak ada kotak oranye
		{1, 40, 1},   // begitu ada pemakaian, minimal satu kotak
		{2, 40, 1},   // 5%
		{3, 40, 2},   // dibulatkan ke atas
		{20, 40, 10}, // separuh
		{39, 40, 19}, // hampir habis tetap belum penuh
		{40, 40, 20}, // habis: penuh
		{55, 40, 20}, // melebihi jatah: tetap penuh
		{5, 0, 0},    // tanpa batas valid
		{-3, 40, 0},  // nilai tak wajar
	}
	for _, c := range cases {
		if got := filledSegments(c.used, c.limit, 20); got != c.want {
			t.Errorf("filledSegments(%d,%d)=%d, want %d", c.used, c.limit, got, c.want)
		}
	}
}

// Warna pada glyph "■": hitung berapa yang oranye, merah, dan kosong.
func countColors(bar string) (orangeN, redN, emptyN int) {
	for _, cell := range strings.Split(bar, " ") {
		switch {
		case strings.Contains(cell, "234;88;12"): // #ea580c
			orangeN++
		case strings.Contains(cell, "248;113;113"): // #f87171
			redN++
		case strings.Contains(cell, "■"):
			emptyN++
		}
	}
	return
}

func TestUsageBarColorsFollowUsage(t *testing.T) {
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })

	zero := usageBar(usageMeter{Used: 0, Limit: 40}, 20)
	if o, r, e := countColors(zero); o != 0 || r != 0 || e != 20 {
		t.Errorf("0%%: oranye=%d merah=%d kosong=%d", o, r, e)
	}

	partial := usageBar(usageMeter{Used: 20, Limit: 40}, 20)
	if o, r, e := countColors(partial); o != 10 || r != 0 || e != 10 {
		t.Errorf("50%%: oranye=%d merah=%d kosong=%d", o, r, e)
	}

	full := usageBar(usageMeter{Used: 40, Limit: 40, Exhausted: true}, 20)
	if o, r, e := countColors(full); o != 0 || r != 20 || e != 0 {
		t.Errorf("habis: oranye=%d merah=%d kosong=%d", o, r, e)
	}

	if got := strings.Count(ansi.Strip(partial), "■"); got != 20 {
		t.Errorf("jumlah kotak=%d, want 20", got)
	}
	if !strings.Contains(ansi.Strip(partial), "■ ■") {
		t.Error("kotak seharusnya bersekat (dipisah spasi)")
	}
}

func usageServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv
}

func snapshotJSON(session, weekly usageMeter) []byte {
	b, _ := json.Marshal(usageSnapshot{Session: session, Weekly: &weekly})
	return b
}

func TestFetchUsage(t *testing.T) {
	reset := time.Now().Add(3 * time.Hour).UTC().Truncate(time.Second)
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer cpc_ok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, _ = w.Write(snapshotJSON(
			usageMeter{Used: 15, Limit: 40, Percent: 38, ResetsAt: &reset},
			usageMeter{Used: 48, Limit: 400, Percent: 12},
		))
	})

	snap, err := fetchUsage(srv.URL, "cpc_ok")
	if err != nil || snap.Session.Used != 15 || snap.Weekly.Limit != 400 || snap.Session.ResetsAt == nil || !snap.Session.ResetsAt.Equal(reset) {
		t.Fatalf("got %+v, %v", snap, err)
	}

	if _, err := fetchUsage(srv.URL, "salah"); err != errInvalidKey {
		t.Fatalf("want errInvalidKey, got %v", err)
	}
}

func TestFormatIn(t *testing.T) {
	cases := map[time.Duration]string{
		10 * time.Second:             "kurang dari 1 mnt",
		42 * time.Minute:             "42 mnt",
		3*time.Hour + 12*time.Minute: "3 j 12 m",
		4*24*time.Hour + 6*time.Hour: "4 hari 6 j",
	}
	for d, want := range cases {
		if got := formatIn(d); got != want {
			t.Errorf("formatIn(%v)=%q, want %q", d, got, want)
		}
	}
}

func TestFormatResetAt(t *testing.T) {
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.Local) // Rabu
	if got := formatResetAt(time.Date(2026, 10, 14, 19, 45, 0, 0, time.Local), now); got != "pukul 19.45" {
		t.Errorf("hari ini: %q", got)
	}
	if got := formatResetAt(time.Date(2026, 10, 19, 0, 0, 0, 0, time.Local), now); got != "Sen 19 Okt, 00.00" {
		t.Errorf("hari lain: %q", got)
	}
}

func TestRenderUsage(t *testing.T) {
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.Local)
	sessionReset := now.Add(3*time.Hour + 12*time.Minute)
	weeklyReset := time.Date(2026, 10, 19, 0, 0, 0, 0, time.Local)

	out := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Used: 15, Limit: 40, Percent: 38, ResetsAt: &sessionReset},
		Weekly:  &usageMeter{Used: 48, Limit: 400, Percent: 12, ResetsAt: &weeklyReset},
	}, "Budi", 110, now))

	for _, want := range []string{
		"Pemakaian · Budi", "Sesi 5 jam", "38%", "15 / 40 pesan", "reset dalam 3 j 12 m (pukul 12.12)",
		"Mingguan", "12%", "48 / 400 pesan", "reset dalam 4 hari 15 j (Sen 19 Okt, 00.00)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("tidak memuat %q:\n%s", want, out)
		}
	}
}

func TestRenderUsageIdleAndExhausted(t *testing.T) {
	now := time.Now()
	idle := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Limit: 40},
		Weekly:  &usageMeter{Limit: 400},
	}, "", 110, now))
	if !strings.Contains(idle, "0 / 40 pesan · mulai saat pesan pertama") || !strings.Contains(idle, "0%") {
		t.Errorf("belum dimulai:\n%s", idle)
	}

	reset := now.Add(2 * time.Hour)
	full := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Used: 40, Limit: 40, Percent: 100, Exhausted: true, ResetsAt: &reset},
		Weekly:  &usageMeter{Used: 40, Limit: 400, Percent: 10},
	}, "", 110, now))
	if !strings.Contains(full, "Habis · reset dalam") || !strings.Contains(full, "100%") {
		t.Errorf("habis:\n%s", full)
	}
}

func TestSmallTerminalUsesTenSegments(t *testing.T) {
	if usageSegments(60) != 10 || usageSegments(110) != 20 {
		t.Fatal("jumlah kotak seharusnya menyesuaikan lebar layar")
	}
}

// ---- integrasi dengan model ----

func loggedIn(srvURL string) *Config {
	return &Config{Name: "Budi Santoso", Email: "budi@example.com", Key: "cpc_ok", URL: srvURL}
}

func TestUsageCommandShowsMeters(t *testing.T) {
	reset := time.Now().Add(3 * time.Hour)
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(snapshotJSON(
			usageMeter{Used: 15, Limit: 40, Percent: 38, ResetsAt: &reset},
			usageMeter{Used: 48, Limit: 400, Percent: 12},
		))
	})
	m := newScreen(t, screenChat, loggedIn(srv.URL), "")

	m, _ = update(m, typed("/usage"))
	m, cmd := update(m, key(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("/usage seharusnya mengambil data")
	}
	m, _ = update(m, cmd())

	view := plain(m)
	for _, want := range []string{"Pemakaian · Budi Santoso", "Sesi 5 jam", "38%", "Mingguan", "15 / 40 pesan"} {
		if !strings.Contains(view, want) {
			t.Errorf("tampilan /usage tidak memuat %q", want)
		}
	}
	if m.usage == nil || m.usage.Session.Used != 15 {
		t.Fatal("snapshot terbaru seharusnya tersimpan")
	}
}

func TestUsageCommandRequiresLogin(t *testing.T) {
	m := newScreen(t, screenChat, nil, "")
	m, _ = update(m, typed("/usage"))
	m, cmd := update(m, key(tea.KeyEnter))

	if cmd != nil || !strings.Contains(plain(m), "Masuk dulu") {
		t.Fatalf("seharusnya meminta login:\n%s", plain(m))
	}
}

func TestUsageCommandReportsServerErrors(t *testing.T) {
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusUnauthorized) })
	m := newScreen(t, screenChat, loggedIn(srv.URL), "")

	m, _ = update(m, typed("/usage"))
	m, cmd := update(m, key(tea.KeyEnter))
	m, _ = update(m, cmd())

	if !strings.Contains(plain(m), "/login") {
		t.Fatalf("key tidak valid seharusnya mengarahkan ke /login:\n%s", plain(m))
	}
}

func TestHeaderShowsUsageOnceLoaded(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	if !strings.Contains(plain(m), "Pemakaian: memuat") {
		t.Fatal("sebelum data tiba seharusnya menampilkan 'memuat'")
	}

	m, _ = update(m, usageMsg{snap: &usageSnapshot{Session: usageMeter{Used: 20, Limit: 40, Percent: 50}}})
	view := plain(m)
	if !strings.Contains(view, "Sesi 5 jam") || !strings.Contains(view, "50%") {
		t.Fatalf("header seharusnya memuat ringkasan pemakaian:\n%s", view)
	}
	if strings.Contains(view, "Pemakaian · ") {
		t.Fatal("pembaruan senyap tidak boleh menambah pesan ke percakapan")
	}
}

func TestStatusIncludesUsageSummary(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	m, _ = update(m, usageMsg{snap: &usageSnapshot{Session: usageMeter{Percent: 38}, Weekly: &usageMeter{Percent: 12}}})

	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))
	if !strings.Contains(plain(m), "sesi 38% · mingguan 12%") {
		t.Fatalf("status seharusnya memuat ringkasan pemakaian:\n%s", plain(m))
	}
}

// ---- paket gratis dan berbayar ----

func freePlan(expires *time.Time, expired bool) usagePlan {
	return usagePlan{Name: "free", Label: "Gratis", ExpiresAt: expires, Expired: expired}
}

func TestFreePlanShowsOnlyTheSessionBarAndTrialExpiry(t *testing.T) {
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.Local)
	trialEnd := now.Add(13*24*time.Hour + 5*time.Hour)
	reset := now.Add(2 * time.Hour)

	out := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Used: 5, Limit: 20, Percent: 25, ResetsAt: &reset},
		Plan:    freePlan(&trialEnd, false),
	}, "Budi", 110, now))

	for _, want := range []string{"Paket Gratis", "masa percobaan berakhir dalam 13 hari 5 j", "Sesi 5 jam", "5 / 20 pesan", "tidak memiliki batas mingguan"} {
		if !strings.Contains(out, want) {
			t.Errorf("tidak memuat %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "Mingguan ") || strings.Contains(out, "/ 400") {
		t.Errorf("paket gratis tidak boleh menampilkan batang mingguan:\n%s", out)
	}
}

func TestPaidPlanShowsBothBarsAndNoTrial(t *testing.T) {
	now := time.Now()
	out := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Used: 1, Limit: 40, Percent: 3},
		Weekly:  &usageMeter{Used: 1, Limit: 400, Percent: 0},
		Plan:    usagePlan{Name: "paid", Label: "Berbayar"},
	}, "", 110, now))

	for _, want := range []string{"Paket Berbayar", "aktif tanpa batas waktu", "Sesi 5 jam", "Mingguan", "1 / 400 pesan"} {
		if !strings.Contains(out, want) {
			t.Errorf("tidak memuat %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "percobaan") {
		t.Error("paket berbayar tidak punya masa percobaan")
	}
}

func TestExpiredPlansAreMarkedInRed(t *testing.T) {
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })

	raw := renderUsage(usageSnapshot{
		Session: usageMeter{Limit: 20},
		Plan:    freePlan(nil, true),
	}, "", 110, time.Now())

	if !strings.Contains(ansi.Strip(raw), "masa percobaan sudah berakhir") {
		t.Fatalf("pesan kedaluwarsa tidak tampil:\n%s", ansi.Strip(raw))
	}
	if !strings.Contains(raw, "248;113;113") {
		t.Error("masa berakhir seharusnya berwarna merah")
	}

	paid := ansi.Strip(renderUsage(usageSnapshot{
		Session: usageMeter{Limit: 40},
		Plan:    usagePlan{Name: "paid", Label: "Berbayar", Expired: true},
	}, "", 110, time.Now()))
	if !strings.Contains(paid, "masa aktif sudah berakhir") {
		t.Errorf("paket berbayar kedaluwarsa:\n%s", paid)
	}
}

func TestRenderUsageToleratesServersWithoutPlanData(t *testing.T) {
	out := ansi.Strip(renderUsage(usageSnapshot{Session: usageMeter{Used: 1, Limit: 40, Percent: 3}}, "", 110, time.Now()))

	if strings.Contains(out, "Paket") && !strings.Contains(out, "tidak memiliki batas mingguan") {
		t.Errorf("tanpa data paket tidak boleh ada baris paket:\n%s", out)
	}
	if !strings.Contains(out, "Sesi 5 jam") {
		t.Error("batang sesi tetap harus tampil")
	}
}

func TestStreamedUsageEventWithNullWeeklyIsParsed(t *testing.T) {
	events, err := collectEvents(t,
		`data: {"type":"usage","session":{"used":2,"limit":20,"percent":10},"weekly":null,`+
			`"plan":{"name":"free","label":"Gratis","expires_at":null,"expired":false}}`+"\n\n"+
			`data: {"type":"done"}`+"\n\n")
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%+v err=%v", events, err)
	}
	snap := events[0].snap
	if snap.Weekly != nil || snap.Plan.Name != "free" || snap.Session.Limit != 20 {
		t.Fatalf("snapshot paket gratis salah: %+v", snap)
	}
}

func TestBlockedByExpiredTrialRollsBackAndExplains(t *testing.T) {
	blocked := `{"type":"blocked","message":"Masa percobaan gratis 14 hari telah berakhir. Beralih ke paket berbayar untuk melanjutkan.","blocked":"plan",` +
		`"session":{"used":0,"limit":20,"percent":0,"exhausted":false,"resets_at":null},"weekly":null,` +
		`"plan":{"name":"free","label":"Gratis","expires_at":null,"expired":true}}`
	fa := newFakeAgent(t, blocked)
	m := send(t, chatWith(t, fa, loggedInCfg()), "halo")

	view := strings.Join(strings.Fields(plain(m)), " ") // kalimat panjang terbungkus ke beberapa baris
	for _, want := range []string{"Masa percobaan gratis 14 hari telah berakhir", "Pesan Anda tidak dikirim", "Paket Gratis", "masa percobaan sudah berakhir"} {
		if !strings.Contains(view, want) {
			t.Errorf("tidak memuat %q:\n%s", want, view)
		}
	}
	if m.input.Value() != "halo" || m.streaming {
		t.Fatal("teks harus kembali ke input dan tidak ada jawaban berjalan")
	}
}

func TestStatusShowsThePlanName(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn("http://x"), "")
	m, _ = update(m, usageMsg{snap: &usageSnapshot{
		Session: usageMeter{Percent: 25},
		Plan:    usagePlan{Name: "free", Label: "Gratis"},
	}})
	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))

	view := plain(m)
	if !strings.Contains(view, "Paket       Gratis") || !strings.Contains(view, "sesi 25%") || strings.Contains(view, "mingguan") {
		t.Fatalf("status paket gratis salah:\n%s", view)
	}
}

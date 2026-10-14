package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"0.1.0", "0.2.0", -1},
		{"0.2.0", "0.2.0", 0},
		{"1.0.0", "0.9.9", 1},
		{"v0.2.0", "0.2.0", 0},
		{"1.2", "1.2.0", 0},
		{"1.10.0", "1.9.9", 1},
		{"0.1.0", "0.1.1", -1},
	}
	for _, c := range cases {
		if got := compareVersions(c.a, c.b); got != c.want {
			t.Errorf("compareVersions(%q, %q) = %d, seharusnya %d", c.a, c.b, got, c.want)
		}
	}
	if !isNewer("0.1.0", "0.2.0") || isNewer("0.2.0", "0.2.0") || isNewer("0.3.0", "0.2.0") {
		t.Error("isNewer salah")
	}
}

func versionServer(t *testing.T, ver string) string {
	t.Helper()
	srv := usageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cli/version" {
			t.Errorf("path salah: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"version": ver,
			"added":   []string{"Fitur A", "Fitur B"},
			"fixed":   []string{"Perbaikan C"},
		})
	})
	return srv.URL
}

func TestUpdateModalOpensForNewerVersion(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, _ = update(m, versionMsg{info: &versionInfo{
		Version: "9.9.9",
		Added:   []string{"Fitur A"},
		Fixed:   []string{"Perbaikan C"},
	}})

	if m.modal == nil || m.modal.kind != modalUpdate {
		t.Fatal("versi baru seharusnya membuka modal update")
	}
	view := plain(m)
	for _, want := range []string{"Pembaruan tersedia", "9.9.9", "Fitur A", "Perbaikan C", "Jangan ingatkan", "Enter"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal tidak memuat %q:\n%s", want, view)
		}
	}
}

func TestUpdateModalCheckboxSkip(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, _ = update(m, versionMsg{info: &versionInfo{Version: "9.9.9"}})

	m, _ = update(m, key(tea.KeySpace)) // centang
	m, _ = update(m, key(tea.KeyEsc))   // nanti

	if m.modal != nil {
		t.Fatal("Esc seharusnya menutup modal")
	}
	if loadSkippedVersion() != "9.9.9" {
		t.Fatal("versi seharusnya tercatat sebagai dilewati")
	}

	// tidak muncul lagi untuk versi yang sama
	m, _ = update(m, versionMsg{info: &versionInfo{Version: "9.9.9"}})
	if m.modal != nil {
		t.Error("versi yang dilewati tidak boleh membuka modal lagi")
	}

	// tapi versi lebih baru tetap muncul
	m, _ = update(m, versionMsg{info: &versionInfo{Version: "9.9.10"}})
	if m.modal == nil {
		t.Error("versi yang lebih baru seharusnya tetap membuka modal")
	}
}

func TestNoModalWhenUpToDate(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, _ = update(m, versionMsg{info: &versionInfo{Version: version}})
	if m.modal != nil {
		t.Error("versi sama seharusnya tidak membuka modal")
	}
}

func TestFetchVersionFromServer(t *testing.T) {
	url := versionServer(t, "0.2.0")
	info, err := fetchVersion(url)
	if err != nil {
		t.Fatalf("fetchVersion: %v", err)
	}
	if info.Version != "0.2.0" || len(info.Added) != 2 || len(info.Fixed) != 1 {
		t.Fatalf("isi salah: %+v", info)
	}
}

func TestManualUpdateReportsLatest(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.updateManual = true
	m, _ = update(m, versionMsg{info: &versionInfo{Version: version}})
	if m.modal != nil {
		t.Fatal("sudah terbaru seharusnya tidak membuka modal")
	}
	if !strings.Contains(plain(m), "versi terbaru") {
		t.Error("seharusnya ada catatan sudah terbaru")
	}
}

func TestUpdateCommandFetches(t *testing.T) {
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, cmd := update(m, typed("/update"))
	m, _ = update(m, key(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("/update seharusnya menjalankan pemeriksaan")
	}
	if !m.updateManual {
		t.Error("flag manual seharusnya menyala")
	}
}

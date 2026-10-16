package main

import (
	"errors"
	"strings"
	"testing"
)

func TestUpdateModalUsesCurlWhenNotNpm(t *testing.T) {
	// Binary tes tidak berada di node_modules/.cipincode → metode curl.
	if updateMethod() != "curl" {
		t.Fatalf("updateMethod = %q, seharusnya curl di lingkungan tes", updateMethod())
	}
	body := updateModalBody(&versionInfo{Version: "9.9.9"}, false)
	if !strings.Contains(body, "install.sh") || strings.Contains(body, "npm update") {
		t.Errorf("modal seharusnya menawarkan update via install.sh:\n%s", body)
	}
}

func TestUpdateDoneErrorSuggestsCurlManual(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, _ = update(m, updateDoneMsg{version: "9.9.9", via: "install.sh", err: errors.New("jaringan putus")})
	if view := plain(m); !strings.Contains(view, "curl -fsSL") {
		t.Errorf("pesan gagal seharusnya menyarankan curl manual:\n%s", view)
	}
}

func TestUpdateDoneSuccessSchedulesRestart(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m, cmd := update(m, updateDoneMsg{version: "9.9.9", via: "install.sh"})
	if cmd == nil {
		t.Fatal("sukses update seharusnya menjadwalkan restart")
	}
	if view := plain(m); !strings.Contains(view, "Me-restart") {
		t.Errorf("seharusnya ada catatan restart otomatis:\n%s", view)
	}

	// restartMsg memakai stub agar tes tidak benar-benar exec.
	old := restartSelfFn
	restartSelfFn = func() error { return nil }
	defer func() { restartSelfFn = old }()
	m, _ = update(m, restartMsg{})
	if view := plain(m); strings.Contains(view, "Restart otomatis gagal") {
		t.Errorf("restart stub sukses tidak boleh error:\n%s", view)
	}

	restartSelfFn = func() error { return errors.New("macet") }
	m, _ = update(m, restartMsg{})
	if view := plain(m); !strings.Contains(view, "Restart otomatis gagal") {
		t.Errorf("restart gagal harus menyuruh buka manual:\n%s", view)
	}
}

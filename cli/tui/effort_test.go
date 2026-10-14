package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

func openEffort(t *testing.T, m model) model {
	t.Helper()
	m, _ = update(m, typed("/effort"))
	m, _ = update(m, key(tea.KeyEnter))
	return m
}

func TestEffortKeysAndParsing(t *testing.T) {
	want := []string{"low", "medium", "high", "xhigh", "ultra"}
	for i, e := range allEfforts {
		if e.key() != want[i] || parseEffort(want[i]) != e {
			t.Errorf("level %d: key=%q parse=%v", i, e.key(), parseEffort(want[i]))
		}
	}
	if parseEffort("") != effortMedium || parseEffort("mega") != effortMedium {
		t.Fatal("nilai tak dikenal seharusnya menjadi medium")
	}
	if len(allEfforts) != 5 {
		t.Fatalf("seharusnya 5 level, got %d", len(allEfforts))
	}
}

func TestStepEffortClampsAtTheEnds(t *testing.T) {
	if stepEffort(effortLow, -1) != effortLow || stepEffort(effortUltra, 1) != effortUltra {
		t.Fatal("slider tidak boleh melewati ujung")
	}
	if stepEffort(effortMedium, 1) != effortHigh || stepEffort(effortMedium, -1) != effortLow {
		t.Fatal("geser satu langkah salah")
	}
}

func TestEveryLevelHasADistinctColor(t *testing.T) {
	seen := map[lipgloss.Color]effort{}
	for _, e := range allEfforts {
		if prev, dup := seen[e.color()]; dup {
			t.Fatalf("warna %s dipakai %s dan %s", e.color(), prev.key(), e.key())
		}
		seen[e.color()] = e
	}
}

func TestSliderShowsAllFiveLabelsAndTheThumbPosition(t *testing.T) {
	for i, e := range allEfforts {
		out := ansi.Strip(effortSlider(e))
		lines := strings.Split(out, "\n")
		if len(lines) != 2 {
			t.Fatalf("slider harus 2 baris, got %d", len(lines))
		}
		for _, name := range []string{"low", "medium", "high", "xhigh", "ultra"} {
			if !strings.Contains(lines[1], name) {
				t.Errorf("label %q hilang saat %s dipilih", name, e.key())
			}
		}
		if strings.Count(lines[0], "●") != 1 {
			t.Errorf("harus ada tepat satu penanda ●: %q", lines[0])
		}
		if idx := strings.Index(lines[0], "●"); idx != -1 {
			// rune ● diawali karakter-karakter 1 kolom per titik
			col := len([]rune(lines[0][:idx]))
			if col != i*sliderGap {
				t.Errorf("posisi penanda level %s: kolom %d, want %d", e.key(), col, i*sliderGap)
			}
		}
	}
}

func TestSliderColorFollowsTheSelectedLevel(t *testing.T) {
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })

	rgb := map[effort]string{
		effortLow:   "52;211;153",  // #34d399
		effortHigh:  "250;204;21",  // #facc15
		effortXHigh: "251;146;60",  // #fb923c
		effortUltra: "232;121;249", // #e879f9
	}
	for e, code := range rgb {
		if !strings.Contains(effortSlider(e), code) {
			t.Errorf("slider level %s seharusnya berwarna %s", e.key(), code)
		}
	}
	if strings.Contains(effortSlider(effortLow), rgb[effortUltra]) {
		t.Error("level low tidak boleh memakai warna ultra")
	}
}

func TestSlashEffortOpensTheSliderModal(t *testing.T) {
	m := openEffort(t, chatModel(t))

	if m.modal == nil || m.modal.kind != modalEffort {
		t.Fatal("/effort seharusnya membuka modal slider")
	}
	view := plain(m)
	for _, want := range []string{"Effort", "low", "medium", "high", "xhigh", "ultra", "seimbang", "semua paket", "← → ubah level", "Esc untuk menutup"} {
		if !strings.Contains(view, want) {
			t.Errorf("modal tidak memuat %q:\n%s", want, view)
		}
	}
	if m.effort != effortMedium {
		t.Fatalf("bawaan seharusnya medium, got %s", m.effort.key())
	}
}

func TestArrowKeysSlideTheEffortAndUpdateTheDescription(t *testing.T) {
	m := openEffort(t, chatModel(t))

	m, _ = update(m, key(tea.KeyRight))
	if m.effort != effortHigh || !strings.Contains(plain(m), "teliti dan menjelaskan alasan") {
		t.Fatalf("→ seharusnya ke high:\n%s", plain(m))
	}
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyRight))
	if m.effort != effortUltra || !strings.Contains(plain(m), "paling boros token") {
		t.Fatalf("seharusnya ultra:\n%s", plain(m))
	}
	m, _ = update(m, key(tea.KeyRight)) // sudah di ujung
	if m.effort != effortUltra {
		t.Fatal("tidak boleh melewati ultra")
	}

	for i := 0; i < 6; i++ {
		m, _ = update(m, key(tea.KeyLeft))
	}
	if m.effort != effortLow || !strings.Contains(plain(m), "paling hemat token") {
		t.Fatalf("seharusnya low:\n%s", plain(m))
	}
}

func TestVimKeysAlsoSlide(t *testing.T) {
	m := openEffort(t, chatModel(t))
	m, _ = update(m, typed("l"))
	m, _ = update(m, typed("l"))
	m, _ = update(m, typed("h"))
	if m.effort != effortHigh {
		t.Fatalf("l,l,h dari medium seharusnya high, got %s", m.effort.key())
	}
}

func TestModalBorderColorChangesWithTheLevel(t *testing.T) {
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })

	m := openEffort(t, chatModel(t))
	medium := m.modal.render(m.width, m.height)
	m, _ = update(m, key(tea.KeyRight)) // high
	m, _ = update(m, key(tea.KeyRight)) // xhigh
	m, _ = update(m, key(tea.KeyRight)) // ultra
	ultra := m.modal.render(m.width, m.height)

	if m.modal.accent != effortUltra.color() {
		t.Fatalf("aksen modal seharusnya warna ultra, got %s", m.modal.accent)
	}
	if !strings.Contains(medium, "56;189;248") || !strings.Contains(ultra, "232;121;249") {
		t.Error("bingkai modal seharusnya berwarna sesuai level")
	}
	if strings.Contains(ultra, "56;189;248") {
		t.Error("warna lama (medium) tidak boleh tersisa setelah digeser")
	}
}

func TestEffortSelectionIsSavedAndRestored(t *testing.T) {
	m := openEffort(t, chatModel(t))
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyRight)) // xhigh

	if loadEffort() != effortXHigh {
		t.Fatalf("pilihan seharusnya tersimpan, got %s", loadEffort().key())
	}

	// sesi baru dengan config yang sama membaca pilihan terakhir
	fresh := newModel(&Agent{URL: "http://x"}, loggedInCfg(), time.Time{}, options{screen: screenChat})
	if fresh.effort != effortXHigh {
		t.Fatalf("sesi baru seharusnya xhigh, got %s", fresh.effort.key())
	}
}

func TestSavingEffortKeepsTheLastSessionTime(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	first := time.Date(2026, 10, 13, 22, 5, 0, 0, time.Local)
	touchSession(first)

	saveEffort(effortUltra)

	if got := touchSession(first.Add(time.Hour)); !got.Equal(first) {
		t.Fatalf("waktu sesi terakhir hilang setelah menyimpan effort: %v", got)
	}
	if loadEffort() != effortUltra {
		t.Fatal("effort seharusnya tetap tersimpan setelah touchSession")
	}
}

func TestCorruptStateFallsBackToMediumWithoutCrashing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	path, _ := statePath()
	_ = os.MkdirAll(filepath.Dir(path), 0o700)
	_ = os.WriteFile(path, []byte("{bukan json"), 0o600)

	if loadEffort() != effortMedium {
		t.Fatal("state rusak seharusnya menghasilkan medium")
	}
}

func TestEscClosesTheEffortModalAndKeepsTheChoice(t *testing.T) {
	m := openEffort(t, chatModel(t))
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyEsc))

	if m.modal != nil {
		t.Fatal("Esc seharusnya menutup")
	}
	if m.effort != effortHigh {
		t.Fatal("pilihan seharusnya tetap setelah ditutup")
	}
	m, _ = update(m, typed("halo"))
	if m.input.Value() != "halo" {
		t.Fatal("input seharusnya aktif lagi")
	}
}

func TestEffortModalIgnoresTypingAndUpDown(t *testing.T) {
	m := openEffort(t, chatModel(t))

	m, _ = update(m, typed("x"))
	m, _ = update(m, key(tea.KeyUp))
	m, _ = update(m, key(tea.KeyDown))
	if m.effort != effortMedium || m.input.Value() != "" || m.modal == nil {
		t.Fatalf("effort=%s input=%q modal=%v", m.effort.key(), m.input.Value(), m.modal != nil)
	}
}

func TestChatFooterShowsTheCurrentEffort(t *testing.T) {
	m := chatModel(t)
	if !strings.Contains(plain(m), "effort medium") {
		t.Fatalf("label seharusnya memuat effort medium:\n%s", plain(m))
	}

	m = openEffort(t, m)
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyEsc))
	if !strings.Contains(plain(m), "effort high") {
		t.Fatalf("label seharusnya memuat effort high:\n%s", plain(m))
	}
}

func TestEffortIsSentWithEveryMessage(t *testing.T) {
	fa := newFakeAgent(t, evHalo, evDone)
	m := chatWith(t, fa, loggedInCfg())
	m = openEffort(t, m)
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyRight))
	m, _ = update(m, key(tea.KeyRight)) // ultra
	m, _ = update(m, key(tea.KeyEsc))

	m = send(t, m, "kerjakan dengan teliti")
	send(t, m, "pertanyaan kedua")

	reqs := fa.requests()
	if len(reqs) != 2 || reqs[0].Body.Effort != "ultra" || reqs[1].Body.Effort != "ultra" {
		t.Fatalf("effort tidak terkirim: %+v", reqs)
	}
}

func TestDefaultEffortIsSentAsMedium(t *testing.T) {
	fa := newFakeAgent(t, evHalo, evDone)
	send(t, chatWith(t, fa, loggedInCfg()), "halo")

	if reqs := fa.requests(); len(reqs) != 1 || reqs[0].Body.Effort != "medium" {
		t.Fatalf("bawaan seharusnya medium: %+v", reqs)
	}
}

func TestEffortAppearsInStatusHelpAndMenu(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, typed("/eff"))
	if items := m.menuItems(); len(items) != 1 || items[0].name != "effort" {
		t.Fatalf("menu: %+v", items)
	}
	m, _ = update(m, key(tea.KeyEsc))

	m, _ = update(m, typed("/status"))
	m, _ = update(m, key(tea.KeyEnter))
	if !strings.Contains(plain(m), "Effort      medium") {
		t.Fatalf("status tanpa effort:\n%s", plain(m))
	}
	m, _ = update(m, key(tea.KeyEsc))

	m, _ = update(m, typed("/help"))
	m, _ = update(m, key(tea.KeyEnter))
	if !strings.Contains(plain(m), "/effort") {
		t.Fatalf("bantuan tanpa /effort:\n%s", plain(m))
	}
}

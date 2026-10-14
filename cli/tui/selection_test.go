package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

// ---- unit: seleksi dan pembersihan ----

func TestSelectionBoundsAreNormalisedForReversedDrags(t *testing.T) {
	down := selection{anchorRow: 1, anchorCol: 2, curRow: 3, curCol: 5}
	up := selection{anchorRow: 3, anchorCol: 5, curRow: 1, curCol: 2}

	r1, c1, r2, c2 := down.bounds()
	u1, d1, u2, d2 := up.bounds()
	if r1 != u1 || c1 != d1 || r2 != u2 || c2 != d2 {
		t.Fatalf("seret ke atas harus sama dengan seret ke bawah: %v vs %v", []int{r1, c1, r2, c2}, []int{u1, d1, u2, d2})
	}
	if r1 != 1 || c1 != 2 || r2 != 3 || c2 != 6 {
		t.Fatalf("bounds=%v, want [1 2 3 6]", []int{r1, c1, r2, c2})
	}

	sameRow := selection{anchorRow: 2, anchorCol: 9, curRow: 2, curCol: 4}
	if _, c1, _, c2 := sameRow.bounds(); c1 != 4 || c2 != 10 {
		t.Fatalf("satu baris, seret ke kiri: c1=%d c2=%d", c1, c2)
	}
}

func TestSelectionTextSingleAndMultiLine(t *testing.T) {
	lines := []string{"alfa beta gamma", "delta epsilon", "zeta eta theta"}

	single := selection{anchorRow: 0, anchorCol: 5, curRow: 0, curCol: 8}
	if got := selectionText(lines, single); got != "beta" {
		t.Errorf("satu baris: %q", got)
	}

	multi := selection{anchorRow: 0, anchorCol: 5, curRow: 2, curCol: 3}
	want := "beta gamma\ndelta epsilon\nzeta"
	if got := selectionText(lines, multi); got != want {
		t.Errorf("banyak baris: %q, want %q", got, want)
	}

	reversed := selection{anchorRow: 2, anchorCol: 3, curRow: 0, curCol: 5}
	if got := selectionText(lines, reversed); got != want {
		t.Errorf("terbalik: %q", got)
	}
}

func TestSelectionTextTrimsTrailingPadding(t *testing.T) {
	lines := []string{"abc          ", "def          "}
	got := selectionText(lines, selection{anchorRow: 0, anchorCol: 0, curRow: 1, curCol: 12})

	if got != "abc\ndef" {
		t.Fatalf("got %q", got)
	}
}

func TestCopiedTextDropsCodeBoxBordersAndMessageMarkers(t *testing.T) {
	lines := []string{
		"● Berikut kodenya:",
		"  ┌─ php ──────────┐",
		"  │ echo 1;        │",
		"  │ echo 2;        │",
		"  └────────────────┘",
		"  Selesai.",
	}
	got := selectionText(lines, selection{anchorRow: 0, anchorCol: 0, curRow: 5, curCol: 30})

	want := "Berikut kodenya:\necho 1;\necho 2;\nSelesai."
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestCopyingOnlyCodeLinesGivesPureCode(t *testing.T) {
	lines := []string{"  │ function a() {  │", "  │     return 1;   │", "  │ }               │"}
	got := selectionText(lines, selection{anchorRow: 0, anchorCol: 0, curRow: 2, curCol: 30})

	if got != "function a() {\n    return 1;\n}" {
		t.Fatalf("indentasi kode harus terjaga: %q", got)
	}
}

func TestPartialSelectionKeepsTextAsIs(t *testing.T) {
	lines := []string{"● jawaban penting di sini"}
	got := selectionText(lines, selection{anchorRow: 0, anchorCol: 2, curRow: 0, curCol: 9})

	if got != "jawaban" {
		t.Fatalf("seleksi parsial tidak boleh dipotong lagi: %q", got)
	}
}

func TestUserMarkerIsDroppedWhenSelectingFromTheLineStart(t *testing.T) {
	got := selectionText([]string{"❯ halo dunia"}, selection{anchorRow: 0, anchorCol: 0, curRow: 0, curCol: 20})
	if got != "halo dunia" {
		t.Fatalf("got %q", got)
	}
}

func TestHighlightReversesOnlyTheSelectedCells(t *testing.T) {
	trueColor(t)
	lines := []string{"alfa beta gamma", "delta epsilon"}
	out := highlightLines(lines, selection{anchorRow: 0, anchorCol: 5, curRow: 0, curCol: 8})

	if !strings.Contains(out[0], "\x1b[7m") {
		t.Fatalf("bagian terpilih seharusnya dibalik warnanya: %q", out[0])
	}
	if ansi.Strip(out[0]) != lines[0] {
		t.Fatalf("teks tidak boleh berubah: %q", ansi.Strip(out[0]))
	}
	if out[1] != lines[1] {
		t.Fatal("baris di luar seleksi harus utuh")
	}
	reversed := ansi.Strip(strings.Split(strings.Split(out[0], "\x1b[7m")[1], "\x1b[0m")[0])
	if reversed != "beta" {
		t.Fatalf("yang dibalik seharusnya 'beta', got %q", reversed)
	}
}

func TestHighlightKeepsStylesAroundTheSelection(t *testing.T) {
	trueColor(t)
	styled := mdBold.Render("tebal") + " biasa " + mdItalic.Render("miring")
	width := ansi.StringWidth(styled)

	out := highlightLines([]string{styled}, selection{anchorRow: 0, anchorCol: 6, curRow: 0, curCol: 9})
	if ansi.StringWidth(out[0]) != width {
		t.Fatalf("lebar baris berubah: %d -> %d", width, ansi.StringWidth(out[0]))
	}
	if !strings.Contains(out[0], "253;186;116") {
		t.Error("gaya tebal di sebelah kiri seleksi harus tetap ada")
	}
}

// ---- alur mouse di model ----

func chatWithConversation(t *testing.T, content string) model {
	t.Helper()
	m := chatModel(t)
	m.home = false
	m.messages = []message{{role: roleUser, content: "tolong"}, {role: roleAssistant, content: content}}
	m.layout()
	m.refresh()
	return m
}

func viewportPlain(m model) []string {
	var out []string
	for _, line := range strings.Split(m.viewport.View(), "\n") {
		out = append(out, ansi.Strip(line))
	}
	return out
}

func rowOf(t *testing.T, m model, needle string) int {
	t.Helper()
	for i, line := range viewportPlain(m) {
		if strings.Contains(line, needle) {
			return i
		}
	}
	t.Fatalf("%q tidak ada di viewport:\n%s", needle, strings.Join(viewportPlain(m), "\n"))
	return -1
}

func mouse(m model, action tea.MouseAction, button tea.MouseButton, x, y int) (model, tea.Cmd) {
	return update(m, tea.MouseMsg(tea.MouseEvent{X: x, Y: y, Action: action, Button: button}))
}

// drag menyeret dari (r1,c1) ke (r2,c2) dalam koordinat viewport dan melepas tombol.
func drag(m model, r1, c1, r2, c2 int) (model, tea.Cmd) {
	top := m.headerRows()
	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, c1, top+r1)
	m, _ = mouse(m, tea.MouseActionMotion, tea.MouseButtonLeft, c2, top+r2)
	return mouse(m, tea.MouseActionRelease, tea.MouseButtonLeft, c2, top+r2)
}

func captureClipboard(t *testing.T) *[]string {
	t.Helper()
	var copied []string
	old := copyToClipboard
	copyToClipboard = func(s string) error { copied = append(copied, s); return nil }
	t.Cleanup(func() { copyToClipboard = old })
	return &copied
}

func TestDraggingSelectsAndCopiesImmediatelyWithAToast(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")

	m, cmd := drag(m, row, 2, row, 9) // "Jawaban"

	if len(*copied) != 1 || (*copied)[0] != "Jawaban" {
		t.Fatalf("clipboard=%q, want [Jawaban]", *copied)
	}
	if m.toast != toastCopied || !strings.Contains(plain(m), "Berhasil disalin ke clipboard") {
		t.Fatalf("notifikasi tidak tampil:\n%s", plain(m))
	}
	if cmd == nil {
		t.Fatal("notifikasi seharusnya dijadwalkan hilang")
	}
	if m.sel == nil {
		t.Fatal("sorotan seharusnya tetap terlihat setelah menyalin")
	}
}

func TestToastReplacesTheRightHintThenRestoresIt(t *testing.T) {
	captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	m, _ = drag(m, row, 2, row, 9)

	view := plain(m)
	if !strings.Contains(view, "Berhasil disalin") || strings.Contains(view, "shift+tab ganti mode") {
		t.Fatalf("toast seharusnya menggantikan petunjuk effort/shift+tab:\n%s", view)
	}

	m, _ = update(m, toastClearMsg{id: m.toastID})
	view = plain(m)
	if strings.Contains(view, "Berhasil disalin") || !strings.Contains(view, "shift+tab ganti mode") || !strings.Contains(view, "effort medium") {
		t.Fatalf("setelah beberapa detik petunjuk harus kembali:\n%s", view)
	}
}

func TestAStaleToastTimerDoesNotClearANewerToast(t *testing.T) {
	captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")

	m, _ = drag(m, row, 2, row, 5)
	firstID := m.toastID
	m, _ = drag(m, row, 2, row, 9)

	m, _ = update(m, toastClearMsg{id: firstID})
	if m.toast == "" {
		t.Fatal("timer lama tidak boleh menghapus notifikasi yang lebih baru")
	}
	m, _ = update(m, toastClearMsg{id: m.toastID})
	if m.toast != "" {
		t.Fatal("timer terbaru seharusnya menghapus notifikasi")
	}
}

func TestDraggingUpwardsAndLeftwardsAlsoCopies(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")

	drag(m, row, 9, row, 2)

	if len(*copied) != 1 || (*copied)[0] != "Jawaban" {
		t.Fatalf("clipboard=%q", *copied)
	}
}

func TestPlainClickDoesNotCopyAndClearsTheHighlight(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	top := m.headerRows()

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 4, top+row)
	m, _ = mouse(m, tea.MouseActionRelease, tea.MouseButtonLeft, 4, top+row)

	if len(*copied) != 0 || m.sel != nil || m.toast != "" {
		t.Fatalf("klik biasa tidak boleh menyalin: copied=%v sel=%v toast=%q", *copied, m.sel, m.toast)
	}
}

func TestSelectionAcrossLinesCopiesCleanText(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Berikut kodenya:\n\n```php\necho 1;\necho 2;\n```\n\nSelesai.")
	first := rowOf(t, m, "Berikut kodenya")
	last := rowOf(t, m, "Selesai.")

	drag(m, first, 0, last, 40)

	if len(*copied) != 1 {
		t.Fatalf("clipboard=%q", *copied)
	}
	want := "Berikut kodenya:\n\necho 1;\necho 2;\n\nSelesai."
	if (*copied)[0] != want {
		t.Fatalf("hasil salinan:\n%q\nseharusnya:\n%q", (*copied)[0], want)
	}
}

func TestCopyingJustTheCodeLinesYieldsRunnableCode(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "```go\nfunc main() {\n    fmt.Println(\"halo\")\n}\n```")
	first := rowOf(t, m, "func main")
	last := rowOf(t, m, "}")

	drag(m, first, 0, last, 60)

	want := "func main() {\n    fmt.Println(\"halo\")\n}"
	if len(*copied) != 1 || (*copied)[0] != want {
		t.Fatalf("hasil salinan %q, want %q", *copied, want)
	}
}

func TestClipboardFailureFallsBackToTheTerminalToast(t *testing.T) {
	old := copyToClipboard
	copyToClipboard = func(string) error { return errors.New("xclip tidak ada") }
	t.Cleanup(func() { copyToClipboard = old })

	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	m, _ = drag(m, row, 2, row, 9)

	if m.toast != toastOSC || !strings.Contains(plain(m), "Disalin ke clipboard terminal") {
		t.Fatalf("toast=%q", m.toast)
	}
}

func TestSelectionIsDrawnWhileDragging(t *testing.T) {
	trueColor(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	top := m.headerRows()

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 2, top+row)
	m, _ = mouse(m, tea.MouseActionMotion, tea.MouseButtonLeft, 9, top+row)

	if !strings.Contains(m.View(), "\x1b[7m") {
		t.Fatal("teks yang sedang diseret seharusnya disorot")
	}
	if m.sel == nil || !m.sel.dragging {
		t.Fatal("seleksi seharusnya masih berlangsung sebelum tombol dilepas")
	}
}

func TestPressingOutsideTheConversationAreaIsIgnored(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama")

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 5, 1) // di header
	m, _ = mouse(m, tea.MouseActionMotion, tea.MouseButtonLeft, 20, 1)
	m, _ = mouse(m, tea.MouseActionRelease, tea.MouseButtonLeft, 20, 1)

	if m.sel != nil || len(*copied) != 0 {
		t.Fatal("seretan di luar area percakapan harus diabaikan")
	}

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 5, m.headerRows()+m.viewport.Height+1) // di kolom input
	if m.sel != nil {
		t.Fatal("klik di area input tidak boleh memulai seleksi")
	}
}

func TestOnlyTheLeftButtonSelects(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	top := m.headerRows()

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonRight, 2, top+row)
	m, _ = mouse(m, tea.MouseActionMotion, tea.MouseButtonRight, 9, top+row)
	m, _ = mouse(m, tea.MouseActionRelease, tea.MouseButtonRight, 9, top+row)

	if m.sel != nil || len(*copied) != 0 {
		t.Fatal("tombol kanan tidak boleh memilih teks")
	}
}

func TestDraggingBeyondTheEdgesIsClamped(t *testing.T) {
	copied := captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama")
	row := rowOf(t, m, "Jawaban pertama")
	top := m.headerRows()

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 2, top+row)
	m, _ = mouse(m, tea.MouseActionMotion, tea.MouseButtonLeft, 500, top+row+500)
	m, _ = mouse(m, tea.MouseActionRelease, tea.MouseButtonLeft, 500, top+row+500)

	if m.sel == nil || m.sel.curRow >= m.viewport.Height || m.sel.curCol >= m.width {
		t.Fatalf("seretan harus dibatasi pada area percakapan: %+v", m.sel)
	}
	if len(*copied) != 1 || !strings.HasPrefix((*copied)[0], "Jawaban pertama") {
		t.Fatalf("clipboard=%q", *copied)
	}
}

func TestTypingClearsTheHighlight(t *testing.T) {
	captureClipboard(t)
	m := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, m, "Jawaban pertama")
	m, _ = drag(m, row, 2, row, 9)

	m, _ = update(m, typed("a"))
	if m.sel != nil {
		t.Fatal("mengetik seharusnya menghapus sorotan")
	}
}

func TestSelectionIsDisabledWhileAModalOrMenuIsOpen(t *testing.T) {
	copied := captureClipboard(t)

	modal := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	row := rowOf(t, modal, "Jawaban pertama")
	modal = modal.openModal("Info", "isi")
	modal, _ = drag(modal, row, 2, row, 9)
	if modal.sel != nil || len(*copied) != 0 {
		t.Fatal("modal terbuka: seleksi di belakangnya harus dimatikan")
	}

	menu := chatWithConversation(t, "Jawaban pertama yang cukup panjang")
	menu, _ = update(menu, typed("/"))
	menu, _ = drag(menu, 1, 2, 1, 9)
	if menu.sel != nil || len(*copied) != 0 {
		t.Fatal("menu terbuka: seleksi harus dimatikan")
	}
}

func TestMouseIsIgnoredOutsideTheChatScreen(t *testing.T) {
	m := newScreen(t, screenMethod, nil, "")
	m, cmd := mouse(m, tea.MouseActionPress, tea.MouseButtonLeft, 5, 5)

	if cmd != nil || m.sel != nil || m.screen != screenMethod {
		t.Fatal("layar login tidak boleh bereaksi pada mouse")
	}
}

// ---- scroll ----

func longConversation(t *testing.T) model {
	t.Helper()
	var b strings.Builder
	for i := 1; i <= 80; i++ {
		fmt.Fprintf(&b, "baris nomor %02d\n\n", i)
	}
	return chatWithConversation(t, b.String())
}

func TestMouseWheelScrollsTheConversation(t *testing.T) {
	m := longConversation(t)
	if !m.viewport.AtBottom() || m.viewport.YOffset == 0 {
		t.Fatalf("percakapan panjang seharusnya mulai di dasar (offset=%d)", m.viewport.YOffset)
	}
	bottom := m.viewport.YOffset

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelUp, 10, 12)
	if m.viewport.YOffset != bottom-wheelStep {
		t.Fatalf("roda ke atas: offset=%d, want %d", m.viewport.YOffset, bottom-wheelStep)
	}

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelDown, 10, 12)
	if m.viewport.YOffset != bottom {
		t.Fatalf("roda ke bawah: offset=%d, want %d", m.viewport.YOffset, bottom)
	}

	for i := 0; i < 100; i++ {
		m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelUp, 10, 12)
	}
	if m.viewport.YOffset != 0 {
		t.Fatalf("di puncak offset harus 0, got %d", m.viewport.YOffset)
	}
}

func TestScrollingBackKeepsTheViewWhileANewAnswerStreams(t *testing.T) {
	m := longConversation(t)
	for i := 0; i < 5; i++ {
		m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelUp, 10, 12)
	}
	offset := m.viewport.YOffset

	m.messages[len(m.messages)-1].content += "\n\nteks baru dari streaming"
	m.refresh()

	if m.viewport.YOffset != offset {
		t.Fatalf("pengguna yang sedang membaca ke atas tidak boleh ditarik ke bawah: %d -> %d", offset, m.viewport.YOffset)
	}
}

func TestScrollingClearsAnExistingHighlight(t *testing.T) {
	captureClipboard(t)
	m := longConversation(t)
	m, _ = drag(m, 3, 0, 3, 8)
	if m.sel == nil {
		t.Fatal("prasyarat: seleksi harus ada")
	}

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelUp, 10, 12)
	if m.sel != nil {
		t.Fatal("isi bergeser: sorotan lama harus dihapus")
	}
}

func TestArrowKeysScrollTheConversationLineByLine(t *testing.T) {
	m := longConversation(t)
	bottom := m.viewport.YOffset

	m, _ = update(m, key(tea.KeyUp))
	m, _ = update(m, key(tea.KeyUp))
	if m.viewport.YOffset != bottom-2 {
		t.Fatalf("↑↑: offset=%d, want %d", m.viewport.YOffset, bottom-2)
	}
	m, _ = update(m, key(tea.KeyDown))
	if m.viewport.YOffset != bottom-1 {
		t.Fatalf("↓: offset=%d, want %d", m.viewport.YOffset, bottom-1)
	}
}

func TestPageKeysStillScroll(t *testing.T) {
	m := longConversation(t)
	bottom := m.viewport.YOffset

	m, _ = update(m, key(tea.KeyPgUp))
	if m.viewport.YOffset >= bottom {
		t.Fatalf("PgUp tidak menggulir: %d", m.viewport.YOffset)
	}
}

func TestWheelScrollsALongModalButNotTheEffortSlider(t *testing.T) {
	m := chatModel(t)
	m, _ = update(m, tea.WindowSizeMsg{Width: 100, Height: 16})
	var body []string
	for i := 1; i <= 30; i++ {
		body = append(body, fmt.Sprintf("baris %02d", i))
	}
	m = m.openModal("Panjang", strings.Join(body, "\n"))

	m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelDown, 10, 5)
	if m.modal.scroll != wheelStep {
		t.Fatalf("roda di modal: scroll=%d, want %d", m.modal.scroll, wheelStep)
	}
	for i := 0; i < 50; i++ {
		m, _ = mouse(m, tea.MouseActionPress, tea.MouseButtonWheelUp, 10, 5)
	}
	if m.modal.scroll != 0 {
		t.Fatalf("scroll tidak boleh negatif: %d", m.modal.scroll)
	}

	slider := openEffort(t, chatModel(t))
	slider, _ = mouse(slider, tea.MouseActionPress, tea.MouseButtonWheelDown, 10, 5)
	if slider.effort != effortMedium || slider.modal == nil {
		t.Fatal("roda mouse tidak boleh mengubah slider effort")
	}
}

// ---- Markdown di percakapan ----

func TestAssistantAnswerIsRenderedAsMarkdownInTheChat(t *testing.T) {
	m := chatWithConversation(t, "Ini **sangat penting** dan `kode inline`.\n\n- satu\n- dua\n\n```php\necho 1;\n```")
	view := plain(m)

	for _, want := range []string{"Ini sangat penting dan  kode inline .", "• satu", "• dua", "┌─ php", "echo 1;", "└"} {
		if !strings.Contains(view, want) {
			t.Errorf("tidak memuat %q:\n%s", want, view)
		}
	}
	for _, bad := range []string{"**", "```", "- satu"} {
		if strings.Contains(view, bad) {
			t.Errorf("sintaks Markdown %q seharusnya sudah dirender:\n%s", bad, view)
		}
	}
}

func TestUserMessagesAreNotRenderedAsMarkdown(t *testing.T) {
	m := chatModel(t)
	m.home = false
	m.messages = []message{{role: roleUser, content: "tolong jelaskan **ini** dan `itu`"}}
	m.refresh()

	if !strings.Contains(plain(m), "**ini**") || !strings.Contains(plain(m), "`itu`") {
		t.Fatalf("pesan pengguna harus apa adanya:\n%s", plain(m))
	}
}

func TestPartialCodeBlockWhileStreamingIsAlreadyABox(t *testing.T) {
	m := chatWithConversation(t, "Berikut:\n```js\nconst a = 1;")
	view := plain(m)

	if !strings.Contains(view, "┌─ js") || !strings.Contains(view, "const a = 1;") || !strings.Contains(view, "└") {
		t.Fatalf("kode yang masih mengalir tetap harus berkotak:\n%s", view)
	}
}

func TestMarkdownRenderingIsCachedUntilTheContentOrWidthChanges(t *testing.T) {
	m := chatWithConversation(t, "Halo **dunia**")
	idx := len(m.messages) - 1

	first := m.messages[idx].rendered
	if first == "" || m.messages[idx].renderedContent != "Halo **dunia**" {
		t.Fatal("hasil render seharusnya tersimpan")
	}

	m.messages[idx].rendered = "SENTINEL"
	m.refresh()
	if m.messages[idx].rendered != "SENTINEL" {
		t.Fatal("tanpa perubahan, hasil cache seharusnya dipakai ulang")
	}

	m.messages[idx].content += " lagi"
	m.refresh()
	if strings.Contains(m.messages[idx].rendered, "SENTINEL") || !strings.Contains(ansi.Strip(m.messages[idx].rendered), "dunia lagi") {
		t.Fatal("isi berubah: harus dirender ulang")
	}

	m.messages[idx].rendered = "SENTINEL"
	m, _ = update(m, tea.WindowSizeMsg{Width: 70, Height: 30})
	if strings.Contains(m.messages[idx].rendered, "SENTINEL") {
		t.Fatal("lebar berubah: harus dirender ulang")
	}
}

func TestCodeBoxFitsInsideNarrowTerminals(t *testing.T) {
	m := chatWithConversation(t, "```go\n"+strings.Repeat("x", 200)+"\n```")
	m, _ = update(m, tea.WindowSizeMsg{Width: 50, Height: 30})

	for i, line := range viewportPlain(m) {
		if w := ansi.StringWidth(line); w > 50 {
			t.Fatalf("baris %d melebihi lebar terminal: %d", i, w)
		}
	}
}

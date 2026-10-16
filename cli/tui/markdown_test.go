package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

func trueColor(t *testing.T) {
	t.Helper()
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })
}

func md(src string) string { return ansi.Strip(renderMarkdown(src, 80)) }

func TestBoldIsStyledNotShownAsAsterisks(t *testing.T) {
	trueColor(t)
	raw := renderMarkdown("ini **sangat penting** ya", 80)

	if strings.Contains(ansi.Strip(raw), "*") {
		t.Fatalf("tanda ** seharusnya hilang: %q", ansi.Strip(raw))
	}
	if ansi.Strip(raw) != "ini sangat penting ya" {
		t.Fatalf("teks salah: %q", ansi.Strip(raw))
	}
	if !strings.Contains(raw, "\x1b[1;") && !strings.Contains(raw, "\x1b[1m") {
		t.Error("teks penting seharusnya tebal")
	}
	if !strings.Contains(raw, "253;186;116") { // #fdba74
		t.Error("teks penting seharusnya berwarna oranye muda, bukan hanya tebal")
	}
}

func TestInlineCodeHasABackgroundAndNoBackticks(t *testing.T) {
	trueColor(t)
	raw := renderMarkdown("jalankan `php artisan serve` dulu", 80)

	if ansi.Strip(raw) != "jalankan  php artisan serve  dulu" {
		t.Fatalf("got %q", ansi.Strip(raw))
	}
	if !strings.Contains(raw, "48;2;39;39;42") { // latar #27272a
		t.Error("kode inline seharusnya punya latar")
	}
}

func TestItalicStrikeAndEscapes(t *testing.T) {
	if got := md("ini *miring* dan _juga miring_"); got != "ini miring dan juga miring" {
		t.Errorf("miring: %q", got)
	}
	if got := md("~~salah~~ benar"); got != "salah benar" {
		t.Errorf("coret: %q", got)
	}
	if got := md(`harga \*bukan\* miring`); got != "harga *bukan* miring" {
		t.Errorf("escape: %q", got)
	}
	if got := md("kali 2 * 3 * 4"); got != "kali 2 * 3 * 4" {
		t.Errorf("perkalian tidak boleh jadi miring: %q", got)
	}
}

func TestSnakeCaseIsNotItalicised(t *testing.T) {
	for _, src := range []string{"pakai nama_variabel_ini saja", "file __init__.py", "set MAX_RETRY_COUNT"} {
		raw := renderMarkdown(src, 80)
		if ansi.Strip(raw) != src && src != "file __init__.py" {
			t.Errorf("%q berubah menjadi %q", src, ansi.Strip(raw))
		}
	}
	if got := md("pakai nama_variabel_ini saja"); got != "pakai nama_variabel_ini saja" {
		t.Errorf("snake_case rusak: %q", got)
	}
}

func TestUnclosedMarkersStayLiteralWhileStreaming(t *testing.T) {
	if got := md("ini **belum selesai"); got != "ini **belum selesai" {
		t.Errorf("got %q", got)
	}
	if got := md("kode `belum"); got != "kode `belum" {
		t.Errorf("got %q", got)
	}
}

func TestLinksShowTextAndUrl(t *testing.T) {
	if got := md("lihat [dokumentasi](https://x.test/docs) ya"); got != "lihat dokumentasi (https://x.test/docs) ya" {
		t.Errorf("got %q", got)
	}
	if got := md("[https://x.test](https://x.test)"); got != "https://x.test" {
		t.Errorf("tautan yang teksnya sama dengan url tidak diulang: %q", got)
	}
	if got := md("bukan [tautan] ya"); got != "bukan [tautan] ya" {
		t.Errorf("got %q", got)
	}
}

func TestHeadings(t *testing.T) {
	got := md("# Judul Besar\n\n## Sub\n\n### Kecil\n\nisi")
	for _, want := range []string{"JUDUL BESAR", "━", "Sub", "Kecil", "isi"} {
		if !strings.Contains(got, want) {
			t.Errorf("tidak memuat %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "#") {
		t.Errorf("tanda # seharusnya hilang:\n%s", got)
	}
}

func TestBulletedNumberedAndTaskLists(t *testing.T) {
	got := md("- satu\n- dua\n  - bersarang\n\n1. pertama\n2. kedua\n\n- [ ] belum\n- [x] selesai")

	for _, want := range []string{"• satu", "• dua", "  ◦ bersarang", "1. pertama", "2. kedua", "☐ belum", "☑ selesai"} {
		if !strings.Contains(got, want) {
			t.Errorf("tidak memuat %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "• satu\n\n• dua") {
		t.Error("butir daftar berurutan tidak boleh diberi baris kosong")
	}
}

func TestLongListItemsWrapWithAHangingIndent(t *testing.T) {
	out := ansi.Strip(renderMarkdown("- "+strings.Repeat("kata ", 30), 40))
	lines := strings.Split(out, "\n")

	if len(lines) < 3 {
		t.Fatalf("seharusnya terbungkus ke beberapa baris:\n%s", out)
	}
	if !strings.HasPrefix(lines[0], "• ") {
		t.Errorf("baris pertama: %q", lines[0])
	}
	for _, l := range lines[1:] {
		if !strings.HasPrefix(l, "  ") {
			t.Errorf("baris lanjutan harus menjorok: %q", l)
		}
	}
	for _, l := range lines {
		if ansi.StringWidth(l) > 40 {
			t.Errorf("baris melebihi lebar: %q", l)
		}
	}
}

func TestBlockquoteAndRule(t *testing.T) {
	got := md("> catatan penting\n> baris dua\n\n---\n\nsetelah")
	for _, want := range []string{"▎ catatan penting", "▎ baris dua", "────", "setelah"} {
		if !strings.Contains(got, want) {
			t.Errorf("tidak memuat %q:\n%s", want, got)
		}
	}
}

func TestTablesAreAligned(t *testing.T) {
	got := md("| Paket | Harga |\n|---|---|\n| Go | Rp 49.000 |\n| Maksimal | Rp 499.000 |")
	lines := strings.Split(got, "\n")

	if len(lines) != 4 {
		t.Fatalf("seharusnya 4 baris:\n%s", got)
	}
	if !strings.Contains(lines[0], "Paket") || !strings.Contains(lines[1], "┼") {
		t.Errorf("header atau pemisah salah:\n%s", got)
	}
	col := strings.Index(lines[2], "│")
	for _, l := range []string{lines[0], lines[2], lines[3]} {
		if strings.Index(l, "│") != col {
			t.Errorf("kolom tidak rata:\n%s", got)
		}
	}
}

func TestTooWideTableFallsBackToRawText(t *testing.T) {
	wide := "| " + strings.Repeat("a", 60) + " | " + strings.Repeat("b", 60) + " |\n|---|---|\n| x | y |"
	if got := ansi.Strip(renderMarkdown(wide, 40)); !strings.Contains(got, "|---|---|") {
		t.Errorf("tabel terlalu lebar seharusnya ditampilkan apa adanya:\n%s", got)
	}
}

func TestParagraphsWrapAndAreSeparatedByABlankLine(t *testing.T) {
	out := ansi.Strip(renderMarkdown("alfa beta gamma delta epsilon zeta eta theta\n\nparagraf dua", 20))
	lines := strings.Split(out, "\n")

	for _, l := range lines {
		if ansi.StringWidth(l) > 20 {
			t.Errorf("melebihi lebar: %q", l)
		}
	}
	if !strings.Contains(out, "\n\nparagraf dua") {
		t.Errorf("paragraf harus dipisah baris kosong:\n%s", out)
	}
}

// ---- blok kode ----

func TestCodeFenceBecomesABoxWithTheLanguageLabel(t *testing.T) {
	got := md("Contoh:\n\n```php\n<?php\necho 'halo';\n```\n\nSelesai.")
	lines := strings.Split(got, "\n")

	var top, bottom int = -1, -1
	for i, l := range lines {
		if strings.HasPrefix(l, "┌─ php ") {
			top = i
		}
		if strings.HasPrefix(l, "└") {
			bottom = i
		}
	}
	if top < 0 || bottom < top {
		t.Fatalf("kotak kode tidak ditemukan:\n%s", got)
	}
	for _, l := range lines[top+1 : bottom] {
		if !strings.HasPrefix(l, "│") || !strings.HasSuffix(l, "│") {
			t.Errorf("baris kode tanpa bingkai: %q", l)
		}
	}
	for _, want := range []string{"<?php", "echo 'halo';", "Contoh:", "Selesai."} {
		if !strings.Contains(got, want) {
			t.Errorf("tidak memuat %q", want)
		}
	}
	if strings.Contains(got, "```") {
		t.Error("tanda ``` seharusnya hilang")
	}
}

func TestCodeBoxHasAFixedWidthSoItNeverBreaks(t *testing.T) {
	out := ansi.Strip(renderMarkdown("```go\nfunc main() {}\n"+strings.Repeat("x", 300)+"\n```", 60))

	widths := map[int]bool{}
	for _, l := range strings.Split(out, "\n") {
		widths[ansi.StringWidth(l)] = true
	}
	if len(widths) != 1 {
		t.Fatalf("semua baris kotak harus sama lebar, got %v:\n%s", widths, out)
	}
	for w := range widths {
		if w > 60 {
			t.Fatalf("kotak melebihi lebar tersedia: %d", w)
		}
	}
}

func TestCodeBoxIsCappedForReadabilityOnWideTerminals(t *testing.T) {
	out := ansi.Strip(renderMarkdown("```\nx\n```", 300))
	if w := ansi.StringWidth(strings.Split(out, "\n")[0]); w != maxCodeWidth {
		t.Fatalf("lebar kotak=%d, want %d", w, maxCodeWidth)
	}
}

func TestUnlabelledCodeIsCalledKode(t *testing.T) {
	if got := md("```\nabc\n```"); !strings.Contains(got, "┌─ kode ") {
		t.Errorf("label bawaan seharusnya 'kode':\n%s", got)
	}
}

func TestUnclosedFenceIsStillRenderedAsABoxWhileStreaming(t *testing.T) {
	got := md("Berikut:\n```js\nconst a = 1;")

	if !strings.Contains(got, "┌─ js ") || !strings.Contains(got, "const a = 1;") || !strings.Contains(got, "└") {
		t.Errorf("kode yang belum selesai tetap harus berbentuk kotak utuh:\n%s", got)
	}
}

func TestCodeInsideTheFenceIsNotParsedAsMarkdown(t *testing.T) {
	got := md("```md\n# bukan judul\n- bukan daftar\n**bukan tebal**\n```")

	for _, want := range []string{"# bukan judul", "- bukan daftar", "**bukan tebal**"} {
		if !strings.Contains(got, want) {
			t.Errorf("isi kode harus apa adanya, tidak memuat %q:\n%s", want, got)
		}
	}
}

func TestTabsInCodeAreExpanded(t *testing.T) {
	if got := md("```\n\tindent\n```"); strings.Contains(got, "\t") || !strings.Contains(got, "    indent") {
		t.Errorf("tab seharusnya diubah menjadi spasi:\n%q", got)
	}
}

func TestSyntaxHighlightingAddsColorsAndAPanelBackground(t *testing.T) {
	trueColor(t)
	raw := renderMarkdown("```php\n<?php\nfunction add($a, $b) { return $a + $b; }\n```", 80)

	if !strings.Contains(raw, "48;2;30;30;46") { // latar panel #1e1e2e
		t.Error("kotak kode seharusnya punya latar gelap sendiri")
	}
	if !strings.Contains(raw, "38;2;205;214;24") { // teks bawaan panel #cdd6f4 (dibulatkan terminal)
		t.Error("teks bawaan kotak kode harus terbaca di latar gelap")
	}
	colors := map[string]bool{}
	for _, part := range strings.Split(raw, "\x1b[38;2;") {
		if end := strings.Index(part, "m"); end > 0 && strings.Count(part[:end], ";") == 2 {
			colors[part[:end]] = true
		}
	}
	if len(colors) < 4 {
		t.Errorf("pewarna sintaks seharusnya memakai beberapa warna, got %d", len(colors))
	}
}

func TestPanelBackgroundSurvivesTokenResets(t *testing.T) {
	trueColor(t)
	raw := renderMarkdown("```go\nfunc main() { fmt.Println(\"halo\") }\n```", 80)

	// setiap kali pewarna mereset gaya, latar panel harus dipasang lagi di baris yang sama
	for _, line := range strings.Split(raw, "\n")[1 : len(strings.Split(raw, "\n"))-1] {
		for _, seg := range strings.Split(line, "\x1b[0m")[1:] {
			if seg == "" || strings.HasPrefix(seg, "\x1b[38;2;") && strings.Contains(seg, "48;2;30;30;46") {
				continue
			}
			if strings.Contains(seg, "│") {
				continue // tepi kanan di luar panel
			}
			if !strings.Contains(seg, "48;2;30;30;46") && ansi.Strip(seg) != "" {
				t.Fatalf("latar panel hilang setelah reset: %q", seg)
			}
		}
	}
}

func TestNoColorTerminalsStillGetAReadableBox(t *testing.T) {
	old := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.Ascii)
	t.Cleanup(func() { lipgloss.SetColorProfile(old) })

	raw := renderMarkdown("**tebal** dan\n```go\nx := 1\n```", 60)
	if strings.Contains(raw, "\x1b") {
		t.Errorf("tanpa warna tidak boleh ada kode escape: %q", raw)
	}
	if !strings.Contains(raw, "┌─ go ") || !strings.Contains(raw, "x := 1") {
		t.Errorf("kotak tetap harus tampil:\n%s", raw)
	}
}

func TestUnknownLanguageFallsBackToGuessingOrPlainText(t *testing.T) {
	trueColor(t)
	got := ansi.Strip(renderMarkdown("```bahasa-fiktif\nteks biasa saja\n```", 60))
	if !strings.Contains(got, "teks biasa saja") || !strings.Contains(got, "bahasa-fiktif") {
		t.Errorf("got %q", got)
	}
}

func TestEmptyAndTinyInputsDoNotPanic(t *testing.T) {
	for _, src := range []string{"", "\n\n", "```", "```\n```", "#", "-", ">", "|", "**", "[", "`"} {
		for _, w := range []int{0, 5, 20, 120} {
			_ = renderMarkdown(src, w)
		}
	}
}

func TestToolCallsRenderAsBox(t *testing.T) {
	src := "Jalankan ini:\n<function_calls>\n<invoke name=\"execute_command\">\n<parameter name=\"command\">ls -la</parameter>\n</invoke>\n</function_calls>\nSelesai."
	got := md(src)
	for _, want := range []string{"execute_command", "command: ls -la", "⚙", "Otomatis dijalankan"} {
		if !strings.Contains(got, want) {
			t.Errorf("kotak tool tidak memuat %q:\n%s", want, got)
		}
	}
	for _, raw := range []string{"<invoke", "<parameter", "</invoke>", "<function_calls>"} {
		if strings.Contains(got, raw) {
			t.Errorf("XML mentah %q seharusnya hilang:\n%s", raw, got)
		}
	}
}

func TestBareInvokeWithoutWrapper(t *testing.T) {
	src := "<invoke name=\"execute_command\">\n<parameter name=\"command\">tree -L 2 -a</parameter>\n</invoke>"
	got := md(src)
	if !strings.Contains(got, "execute_command") || !strings.Contains(got, "tree -L 2 -a") {
		t.Errorf("invoke tanpa pembungkus harus tetap rapi:\n%s", got)
	}
}

func TestUnclosedToolCallWhileStreaming(t *testing.T) {
	src := "<function_calls>\n<invoke name=\"execute_command\">\n<parameter name=\"command\">ls -la"
	got := md(src)
	if !strings.Contains(got, "execute_command") || !strings.Contains(got, "ls -la") {
		t.Errorf("blok tool yang belum tertutup harus tetap tampil:\n%s", got)
	}
	if strings.Contains(got, "<invoke") {
		t.Errorf("XML mentah seharusnya hilang:\n%s", got)
	}
}

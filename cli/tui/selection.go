package main

import (
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

// selection adalah teks yang disorot dengan menyeret mouse di area percakapan.
// Koordinat relatif terhadap pojok kiri atas viewport.
type selection struct {
	anchorRow, anchorCol int
	curRow, curCol       int
	moved                bool
	dragging             bool
}

var selectedStyle = lipgloss.NewStyle().Reverse(true)

// bounds mengurutkan awal dan akhir (seret ke atas atau ke kiri tetap valid). Kolom akhir
// bersifat eksklusif dan sel di bawah kursor ikut terpilih.
func (s selection) bounds() (r1, c1, r2, c2 int) {
	r1, c1, r2, c2 = s.anchorRow, s.anchorCol, s.curRow, s.curCol
	if r1 > r2 || (r1 == r2 && c1 > c2) {
		r1, c1, r2, c2 = r2, c2, r1, c1
	}
	return r1, c1, r2, c2 + 1
}

// span adalah rentang kolom yang terpilih pada baris tertentu.
func (s selection) span(row, lineWidth int) (from, to int) {
	r1, c1, r2, c2 := s.bounds()
	from, to = 0, lineWidth
	if row == r1 {
		from = c1
	}
	if row == r2 {
		to = c2
	}
	return from, min(to, max(lineWidth, to))
}

// Jenis baris menurut tampilan polosnya: bagian kotak kode atau teks biasa.
const (
	rowNone    = iota
	rowContent // baris isi kotak kode: "│ kode │"
	rowBorder  // garis atas atau bawah kotak kode
)

// classifyRow mengenali baris kotak kode dan kolom isinya [start, end), tanpa bingkai "│ "
// di kiri dan " │" di kanan. Seleksi dibatasi ke rentang ini agar yang tersorot dan
// tersalin hanya kodenya.
func classifyRow(plain string) (kind, start, end int) {
	trimmed := strings.TrimLeft(plain, " ")
	indent := len(plain) - len(trimmed)
	right := strings.TrimRight(trimmed, " ")

	switch {
	case strings.HasPrefix(trimmed, "┌"), strings.HasPrefix(trimmed, "└"):
		return rowBorder, 0, 0
	case strings.HasPrefix(trimmed, "│") && strings.HasSuffix(right, "│") && ansi.StringWidth(right) >= 4:
		return rowContent, indent + 2, indent + ansi.StringWidth(right) - 2
	}
	return rowNone, 0, 0
}

// codeBox mengembalikan baris isi pertama dan terakhir dari kotak kode yang memuat row.
func codeBox(plain []string, row int) (first, last int, ok bool) {
	if row < 0 || row >= len(plain) {
		return 0, 0, false
	}
	if kind, _, _ := classifyRow(plain[row]); kind != rowContent {
		return 0, 0, false
	}
	first, last = row, row
	for first > 0 {
		if kind, _, _ := classifyRow(plain[first-1]); kind != rowContent {
			break
		}
		first--
	}
	for last+1 < len(plain) {
		if kind, _, _ := classifyRow(plain[last+1]); kind != rowContent {
			break
		}
		last++
	}
	return first, last, true
}

// spanIn adalah span() yang menghormati kotak kode: pada baris kotak hanya isinya yang
// terpilih, dan garis atas/bawah kotak tidak pernah terpilih.
func (s selection) spanIn(row int, plain string) (from, to int) {
	from, to = s.span(row, ansi.StringWidth(plain))

	switch kind, start, end := classifyRow(plain); kind {
	case rowBorder:
		return 0, 0
	case rowContent:
		from, to = max(from, start), min(to, end)
	}
	return from, to
}

// highlightLines membalik warna teks terpilih. Isi di luar seleksi tidak diubah.
func highlightLines(lines []string, s selection) []string {
	r1, _, r2, _ := s.bounds()
	out := make([]string, len(lines))
	copy(out, lines)

	for row := max(r1, 0); row <= r2 && row < len(lines); row++ {
		line := lines[row]
		from, to := s.spanIn(row, ansi.Strip(line))
		if to <= from {
			continue
		}
		left := ansi.Truncate(line, from, "")
		mid := ansi.Strip(ansi.Truncate(ansi.TruncateLeft(line, from, ""), to-from, ""))
		right := ansi.TruncateLeft(line, to, "")
		out[row] = left + "\x1b[0m" + selectedStyle.Render(mid) + "\x1b[0m" + right
	}
	return out
}

// selectionText mengambil teks polos yang terpilih dan membersihkannya agar enak ditempel:
// bingkai kotak kode dan penanda pesan tidak ikut tersalin.
func selectionText(plainLines []string, s selection) string {
	r1, _, r2, _ := s.bounds()
	var out []string

	for row := max(r1, 0); row <= r2 && row < len(plainLines); row++ {
		line := plainLines[row]
		from, to := s.spanIn(row, line)
		if to <= from {
			if kind, _, _ := classifyRow(line); kind != rowBorder { // garis kotak kode dilewati
				out = append(out, "")
			}
			continue
		}
		segment := ansi.Truncate(ansi.TruncateLeft(line, from, ""), to-from, "")
		if cleaned, keep := cleanCopiedLine(segment, from == 0); keep {
			out = append(out, cleaned)
		}
	}
	return strings.Trim(strings.Join(out, "\n"), "\n")
}

func cleanCopiedLine(segment string, fromLineStart bool) (string, bool) {
	line := strings.TrimRight(segment, " ")
	bare := strings.TrimLeft(line, " ")

	switch {
	case strings.HasPrefix(bare, "┌"), strings.HasPrefix(bare, "└"):
		return "", false // garis atas dan bawah kotak kode
	case strings.HasPrefix(bare, "│"):
		code := strings.TrimPrefix(bare, "│")
		code = strings.TrimPrefix(code, " ")
		code = strings.TrimSuffix(strings.TrimRight(code, " "), "│")
		return strings.TrimRight(code, " "), true
	case fromLineStart:
		for _, marker := range []string{"● ", "❯ "} {
			if strings.HasPrefix(line, marker) {
				return strings.TrimPrefix(line, marker), true
			}
		}
		return strings.TrimPrefix(line, "  "), true // jorokan lanjutan pesan
	}
	return line, true
}

// copyToClipboard bisa diganti di test.
var copyToClipboard = copyText

// copyText menyalin lewat dua jalur: OSC 52 (bekerja lewat SSH dan di banyak terminal)
// dan perintah clipboard sistem (xclip, xsel, wl-copy, pbcopy).
func copyText(text string) error {
	termenv.NewOutput(os.Stdout).Copy(text)
	return clipboard.WriteAll(text)
}

package main

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/muesli/termenv"
)

// Renderer Markdown untuk jawaban AI: teks penting ditebalkan dengan warna, kode inline
// diberi latar, dan blok kode dibuat sebagai kotak tersendiri dengan pewarna sintaks.

var (
	mdBold   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#c2410c", Dark: "#fdba74"})
	mdItalic = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.AdaptiveColor{Light: "#52525b", Dark: "#d4d4d8"})
	mdStrike = lipgloss.NewStyle().Strikethrough(true).Foreground(muted)
	mdCode   = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#be123c", Dark: "#fca5a5"}).
			Background(lipgloss.AdaptiveColor{Light: "#f4f4f5", Dark: "#27272a"})
	mdLink   = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.AdaptiveColor{Light: "#0369a1", Dark: "#7dd3fc"})
	mdH1     = lipgloss.NewStyle().Bold(true).Foreground(orange)
	mdH2     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#c2410c", Dark: "#fb923c"})
	mdH3     = lipgloss.NewStyle().Bold(true)
	mdBullet = lipgloss.NewStyle().Foreground(orange)
	mdQuote  = lipgloss.NewStyle().Foreground(muted)
	mdBorder = lipgloss.NewStyle().Foreground(muted)
	mdLabel  = lipgloss.NewStyle().Bold(true).Foreground(orange)
)

const (
	codeBg = lipgloss.Color("#1e1e2e")
	codeFg = lipgloss.Color("#cdd6f4")

	maxCodeWidth = 100
)

var (
	reFence   = regexp.MustCompile("^\\s*(```+|~~~+)\\s*([\\w+#.-]*)\\s*$")
	reHeading = regexp.MustCompile(`^\s{0,3}(#{1,6})\s+(.+?)\s*#*\s*$`)
	reRule    = regexp.MustCompile(`^\s{0,3}([-*_])(\s*[-*_]){2,}\s*$`)
	reList    = regexp.MustCompile(`^(\s*)([-*+]|\d{1,3}[.)])\s+(.*)$`)
	reQuote   = regexp.MustCompile(`^\s{0,3}>\s?(.*)$`)
	reTable   = regexp.MustCompile(`^\s*\|.*\|\s*$`)
	reTableSp = regexp.MustCompile(`^\s*\|?\s*:?-{2,}:?\s*(\|\s*:?-{2,}:?\s*)*\|?\s*$`)
	reTask    = regexp.MustCompile(`^\[( |x|X)\]\s+(.*)$`)
)

func wrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	return ansi.Wrap(s, width, "")
}

// renderMarkdown mengubah Markdown menjadi teks berwarna selebar width kolom.
func renderMarkdown(src string, width int) string {
	width = max(width, 20)
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	var out []string
	prev := "" // jenis blok sebelumnya, untuk mengatur jarak
	add := func(kind string, block string) {
		if prev != "" && !(kind == "list" && prev == "list") {
			out = append(out, "")
		}
		out = append(out, strings.Split(block, "\n")...)
		prev = kind
	}

	for i := 0; i < len(lines); {
		line := lines[i]

		switch {
		case strings.TrimSpace(line) == "":
			i++

		case reFence.MatchString(line):
			m := reFence.FindStringSubmatch(line)
			marker, lang := m[1], m[2]
			var code []string
			i++
			for i < len(lines) && !isFenceEnd(lines[i], marker) {
				code = append(code, lines[i])
				i++
			}
			i++ // lewati penutup (atau akhir teks jika belum tertutup, mis. saat streaming)
			add("code", renderCodeBlock(lang, strings.Join(code, "\n"), width))

		case reHeading.MatchString(line):
			m := reHeading.FindStringSubmatch(line)
			add("heading", renderHeading(len(m[1]), m[2], width))
			i++

		case reRule.MatchString(line) && !reList.MatchString(line):
			add("rule", mdBorder.Render(strings.Repeat("─", min(width, 48))))
			i++

		case reTable.MatchString(line) && i+1 < len(lines) && reTableSp.MatchString(lines[i+1]):
			j := i + 2
			for j < len(lines) && reTable.MatchString(lines[j]) {
				j++
			}
			add("table", renderTable(lines[i:j], width))
			i = j

		case reQuote.MatchString(line):
			var quoted []string
			for i < len(lines) && reQuote.MatchString(lines[i]) {
				quoted = append(quoted, reQuote.FindStringSubmatch(lines[i])[1])
				i++
			}
			body := wrapText(renderInlineLines(quoted), width-2)
			var b []string
			for _, l := range strings.Split(body, "\n") {
				b = append(b, mdQuote.Render("▎ ")+mdItalic.Render(ansi.Strip(l)))
			}
			add("quote", strings.Join(b, "\n"))

		case reList.MatchString(line):
			add("list", renderListItem(line, width))
			i++

		default:
			var para []string
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" && !startsBlock(lines[i], lines, i) {
				para = append(para, lines[i])
				i++
			}
			add("para", wrapText(renderInlineLines(para), width))
		}
	}
	return strings.Join(out, "\n")
}

func isFenceEnd(line, marker string) bool {
	t := strings.TrimSpace(line)
	return strings.HasPrefix(t, marker[:3]) && strings.Trim(t, string(marker[0])) == ""
}

func startsBlock(line string, lines []string, i int) bool {
	switch {
	case reFence.MatchString(line), reHeading.MatchString(line), reQuote.MatchString(line):
		return true
	case reList.MatchString(line):
		return true
	case reRule.MatchString(line):
		return true
	case reTable.MatchString(line) && i+1 < len(lines) && reTableSp.MatchString(lines[i+1]):
		return true
	}
	return false
}

func renderInlineLines(lines []string) string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = renderInline(strings.TrimRight(l, " \t"))
	}
	return strings.Join(out, "\n")
}

func renderHeading(level int, text string, width int) string {
	body := renderInline(text)
	switch level {
	case 1:
		plain := ansi.Strip(body)
		return mdH1.Render(strings.ToUpper(plain)) + "\n" + mdBorder.Render(strings.Repeat("━", min(width, max(ansi.StringWidth(plain), 8))))
	case 2:
		return mdH2.Render(ansi.Strip(body))
	case 3:
		return mdH3.Render(ansi.Strip(body))
	}
	return mdH3.Foreground(muted).Render(ansi.Strip(body))
}

func renderListItem(line string, width int) string {
	m := reList.FindStringSubmatch(line)
	level := len(strings.ReplaceAll(m[1], "\t", "  ")) / 2
	marker, text := m[2], m[3]

	glyph := marker
	switch {
	case marker == "-" || marker == "*" || marker == "+":
		glyph = []string{"•", "◦", "▪"}[min(level, 2)]
	}
	if t := reTask.FindStringSubmatch(text); t != nil {
		glyph = "☐"
		if t[1] != " " {
			glyph = "☑"
		}
		text = t[2]
	}

	prefix := strings.Repeat("  ", level) + glyph + " "
	body := wrapText(renderInline(text), max(width-ansi.StringWidth(prefix), 10))

	rows := strings.Split(body, "\n")
	pad := strings.Repeat(" ", ansi.StringWidth(prefix))
	rows[0] = strings.Repeat("  ", level) + mdBullet.Render(glyph) + " " + rows[0]
	for i := 1; i < len(rows); i++ {
		rows[i] = pad + rows[i]
	}
	return strings.Join(rows, "\n")
}

func splitRow(line string) []string {
	t := strings.TrimSpace(line)
	t = strings.TrimPrefix(t, "|")
	t = strings.TrimSuffix(t, "|")
	cells := strings.Split(t, "|")
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}
	return cells
}

func renderTable(lines []string, width int) string {
	header := splitRow(lines[0])
	var rows [][]string
	for _, l := range lines[2:] {
		rows = append(rows, splitRow(l))
	}

	cols := len(header)
	widths := make([]int, cols)
	cell := func(row []string, c int) string {
		if c < len(row) {
			return renderInline(row[c])
		}
		return ""
	}
	for c := 0; c < cols; c++ {
		widths[c] = ansi.StringWidth(cell(header, c))
		for _, r := range rows {
			widths[c] = max(widths[c], ansi.StringWidth(cell(r, c)))
		}
	}

	total := cols*3 - 3
	for _, w := range widths {
		total += w
	}
	if total > width { // terlalu lebar: tampilkan apa adanya, tetap terbaca
		return mdBorder.Render(strings.Join(lines, "\n"))
	}

	line := func(row []string, style func(string) string) string {
		parts := make([]string, cols)
		for c := 0; c < cols; c++ {
			text := cell(row, c)
			parts[c] = style(text) + strings.Repeat(" ", widths[c]-ansi.StringWidth(text))
		}
		return strings.Join(parts, mdBorder.Render(" │ "))
	}

	sep := make([]string, cols)
	for c := range sep {
		sep[c] = strings.Repeat("─", widths[c])
	}

	out := []string{line(header, func(s string) string { return mdH3.Render(ansi.Strip(s)) }), mdBorder.Render(strings.Join(sep, "─┼─"))}
	for _, r := range rows {
		out = append(out, line(r, func(s string) string { return s }))
	}
	return strings.Join(out, "\n")
}

// ---- inline ----

func renderInline(s string) string {
	r := []rune(s)
	var out strings.Builder

	for i := 0; i < len(r); {
		c := r[i]
		switch {
		case c == '\\' && i+1 < len(r) && strings.ContainsRune("\\\\`*_{}[]()#+-.!|~>", r[i+1]):
			out.WriteRune(r[i+1])
			i += 2

		case c == '`':
			if end := indexOf(r, "`", i+1); end > i+1 {
				out.WriteString(mdCode.Render(" " + string(r[i+1:end]) + " "))
				i = end + 1
				continue
			}
			out.WriteRune(c)
			i++

		case hasAt(r, i, "**"), hasAt(r, i, "__") && boundaryBefore(r, i):
			marker := string(r[i : i+2])
			if end := indexOf(r, marker, i+2); end > i+2 && strings.TrimSpace(string(r[i+2:end])) != "" &&
				(marker == "**" || boundaryAfter(r, end+2)) {
				out.WriteString(mdBold.Render(string(r[i+2 : end])))
				i = end + 2
				continue
			}
			out.WriteString(marker)
			i += 2

		case hasAt(r, i, "~~"):
			if end := indexOf(r, "~~", i+2); end > i+2 {
				out.WriteString(mdStrike.Render(string(r[i+2 : end])))
				i = end + 2
				continue
			}
			out.WriteString("~~")
			i += 2

		case c == '*' || (c == '_' && boundaryBefore(r, i)):
			if end := italicEnd(r, i); end > 0 {
				out.WriteString(mdItalic.Render(string(r[i+1 : end])))
				i = end + 1
				continue
			}
			out.WriteRune(c)
			i++

		case c == '[':
			if text, url, next, ok := parseLink(r, i); ok {
				out.WriteString(mdLink.Render(text))
				if url != text {
					out.WriteString(mdQuote.Render(" (" + url + ")"))
				}
				i = next
				continue
			}
			out.WriteRune(c)
			i++

		default:
			out.WriteRune(c)
			i++
		}
	}
	return out.String()
}

func hasAt(r []rune, i int, s string) bool {
	sr := []rune(s)
	if i+len(sr) > len(r) {
		return false
	}
	for k, x := range sr {
		if r[i+k] != x {
			return false
		}
	}
	return true
}

func indexOf(r []rune, s string, from int) int {
	for i := from; i < len(r); i++ {
		if hasAt(r, i, s) {
			return i
		}
	}
	return -1
}

func isWordRune(c rune) bool { return unicode.IsLetter(c) || unicode.IsDigit(c) }

// boundaryBefore: tanda garis bawah hanya berarti penekanan jika tidak menempel pada kata
// (agar nama_variabel_ini tidak ikut dimiringkan).
func boundaryBefore(r []rune, i int) bool { return i == 0 || !isWordRune(r[i-1]) }

func boundaryAfter(r []rune, i int) bool { return i >= len(r) || !isWordRune(r[i]) }

func italicEnd(r []rune, i int) int {
	marker := r[i]
	if i+1 >= len(r) || unicode.IsSpace(r[i+1]) || r[i+1] == marker {
		return -1
	}
	for j := i + 2; j < len(r); j++ {
		if r[j] != marker || unicode.IsSpace(r[j-1]) {
			continue
		}
		if j+1 < len(r) && r[j+1] == marker { // bagian dari ** atau __
			j++
			continue
		}
		if marker == '_' && !boundaryAfter(r, j+1) {
			continue
		}
		return j
	}
	return -1
}

func parseLink(r []rune, i int) (text, url string, next int, ok bool) {
	close := indexOf(r, "]", i+1)
	if close < 0 || close+1 >= len(r) || r[close+1] != '(' {
		return "", "", 0, false
	}
	end := indexOf(r, ")", close+2)
	if end < 0 {
		return "", "", 0, false
	}
	text, url = string(r[i+1:close]), string(r[close+2:end])
	if text == "" || url == "" || strings.ContainsAny(url, " \t") {
		return "", "", 0, false
	}
	return text, url, end + 1, true
}

// ---- blok kode ----

func colorSeq(fg, bg lipgloss.Color) string {
	probe := lipgloss.NewStyle().Foreground(fg).Background(bg).Render("X")
	if idx := strings.Index(probe, "X"); idx > 0 {
		return probe[:idx]
	}
	return ""
}

func chromaFormatter(p termenv.Profile) chroma.Formatter {
	switch p {
	case termenv.TrueColor:
		return formatters.Get("terminal16m")
	case termenv.ANSI256:
		return formatters.Get("terminal256")
	case termenv.ANSI:
		return formatters.Get("terminal16")
	}
	return nil
}

// highlightCode memberi warna sintaks. Jika bahasa tidak dikenal, tebak dari isinya.
func highlightCode(code, lang string) string {
	formatter := chromaFormatter(lipgloss.ColorProfile())
	if formatter == nil {
		return code
	}
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		return code
	}
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}
	it, err := chroma.Coalesce(lexer).Tokenise(nil, code)
	if err != nil {
		return code
	}
	var b strings.Builder
	if err := formatter.Format(&b, style, it); err != nil {
		return code
	}
	return strings.TrimRight(b.String(), "\n")
}

// renderCodeBlock menggambar kode sebagai kotak tersendiri: judul bahasa di atas,
// latar gelap, dan baris yang terlalu panjang dibungkus agar kotak tidak pecah.
func renderCodeBlock(lang, code string, width int) string {
	width = min(max(width, 24), maxCodeWidth)
	inner := width - 4

	label := lang
	if label == "" {
		label = "kode"
	}

	// Deteksi diff / patch (git diff) untuk memperindah tampilan tambah/hapus baris
	isDiff := lang == "diff" || lang == "patch" || strings.Contains(code, "\n+") || strings.Contains(code, "\n-")

	code = strings.TrimRight(strings.ReplaceAll(code, "\t", "    "), "\n")
	
	var rows []string
	if isDiff {
		for _, line := range strings.Split(code, "\n") {
			prefix := "  "
			colorStyle := lipgloss.NewStyle().Foreground(codeFg)
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				prefix = "+ "
				colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")) // Hijau lembut
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				prefix = "- "
				colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")) // Merah lembut
			} else if strings.HasPrefix(line, "@@") {
				colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")).Bold(true)
			}

			content := line
			if (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, " ")) && (lang == "diff" || lang == "patch") {
				content = line[1:]
			}

			formatted := colorStyle.Render(prefix + content)
			for _, part := range strings.Split(ansi.Hardwrap(formatted, inner, true), "\n") {
				text := " " + part + strings.Repeat(" ", max(0, inner-ansi.StringWidth(ansi.Strip(part)))) + " "
				rows = append(rows, mdBorder.Render("│")+text+mdBorder.Render("│"))
			}
		}
	} else {
		highlighted := highlightCode(code, lang)
		panel := colorSeq(codeFg, codeBg)
		for _, line := range strings.Split(highlighted, "\n") {
			for _, part := range strings.Split(ansi.Hardwrap(line, inner, true), "\n") {
				text := " " + part + strings.Repeat(" ", max(0, inner-ansi.StringWidth(ansi.Strip(part)))) + " "
				if panel != "" {
					text = panel + strings.ReplaceAll(text, "\x1b[0m", "\x1b[0m"+panel) + "\x1b[0m"
				}
				rows = append(rows, mdBorder.Render("│")+text+mdBorder.Render("│"))
			}
		}
	}

	top := mdBorder.Render("┌─ ") + mdLabel.Render(label) + mdBorder.Render(" "+strings.Repeat("─", max(1, width-len([]rune(label))-5))+"┐")
	bottom := mdBorder.Render("└" + strings.Repeat("─", width-2) + "┘")
	return top + "\n" + strings.Join(rows, "\n") + "\n" + bottom
}

package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Logo CP dibuat dari public/images/logoclii.png (setengah blok, 16x6).
// logoClasses menandai warna tiap karakter: o = oranye, g = abu-abu.
const logoWidth = 16

var logoLines = []string{
	"   ▄███████▀    ",
	"  ██▀ ▄▄▄▄▄▄▄▄▄ ",
	"▄██▀ ▀▀▀▀▀▀▀▀███",
	"▀█████▀▄██████▀ ",
	"     ▄██▀       ",
	"    ██▀         ",
}

var logoClasses = []string{
	"   ggggggggg    ",
	"  ggg ooooooooo ",
	"gggg ooooooooooo",
	"gggggggoooooooo ",
	"     oooo       ",
	"    ooo         ",
}

var (
	logoOrange = lipgloss.NewStyle().Foreground(lipgloss.Color("#e06000"))
	logoGray   = lipgloss.NewStyle().Foreground(lipgloss.Color("#a1a1aa"))
)

func renderLogo() string { return renderLogoRows(len(logoLines)) }

// renderLogoRows menggambar n baris pertama; sisanya kosong agar tinggi tetap
// (dipakai animasi onboarding).
func renderLogoRows(n int) string {
	rows := make([]string, len(logoLines))
	for i, line := range logoLines {
		if i >= n {
			rows[i] = strings.Repeat(" ", logoWidth)
			continue
		}
		classes := []rune(logoClasses[i])
		var b strings.Builder
		for j, r := range []rune(line) {
			switch classes[j] {
			case 'o':
				b.WriteString(logoOrange.Render(string(r)))
			case 'g':
				b.WriteString(logoGray.Render(string(r)))
			default:
				b.WriteRune(' ')
			}
		}
		rows[i] = b.String()
	}
	return strings.Join(rows, "\n")
}

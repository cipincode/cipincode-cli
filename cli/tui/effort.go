package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// effort adalah tingkat kedalaman kerja model. Berlaku untuk semua paket.
type effort int

const (
	effortLow effort = iota
	effortMedium
	effortHigh
	effortXHigh
	effortUltra
)

var allEfforts = []effort{effortLow, effortMedium, effortHigh, effortXHigh, effortUltra}

func (e effort) key() string {
	switch e {
	case effortLow:
		return "low"
	case effortHigh:
		return "high"
	case effortXHigh:
		return "xhigh"
	case effortUltra:
		return "ultra"
	}
	return "medium"
}

// parseEffort membaca nilai tersimpan; yang tidak dikenal menjadi medium.
func parseEffort(s string) effort {
	for _, e := range allEfforts {
		if e.key() == s {
			return e
		}
	}
	return effortMedium
}

func (e effort) description() string {
	switch e {
	case effortLow:
		return "cepat dan singkat · paling hemat token"
	case effortHigh:
		return "teliti dan menjelaskan alasan · lebih lambat"
	case effortXHigh:
		return "mendalam, memeriksa kasus tepi · lebih lambat dan lebih boros token"
	case effortUltra:
		return "maksimum: banyak pendekatan dan verifikasi · paling lambat dan paling boros token"
	}
	return "seimbang · bawaan"
}

// color berubah di sepanjang slider: hijau, biru, kuning, oranye, lalu magenta.
func (e effort) color() lipgloss.Color {
	switch e {
	case effortLow:
		return lipgloss.Color("#34d399")
	case effortHigh:
		return lipgloss.Color("#facc15")
	case effortXHigh:
		return lipgloss.Color("#fb923c")
	case effortUltra:
		return lipgloss.Color("#e879f9")
	}
	return lipgloss.Color("#38bdf8")
}

// stepEffort menggeser level dan berhenti di ujung (tidak berputar).
func stepEffort(e effort, delta int) effort {
	next := int(e) + delta
	return allEfforts[min(max(next, 0), len(allEfforts)-1)]
}

const sliderGap = 10 // lebar antar titik slider

// effortSlider menggambar slider dua baris: lintasan dan label. Bagian yang sudah
// "terisi" sampai level terpilih memakai warna level itu.
func effortSlider(selected effort) string {
	on := lipgloss.NewStyle().Foreground(selected.color())
	off := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#d4d4d8", Dark: "#3f3f46"})

	var track strings.Builder
	for i := range allEfforts {
		switch {
		case i == int(selected):
			track.WriteString(on.Bold(true).Render("●"))
		case i < int(selected):
			track.WriteString(on.Render("○"))
		default:
			track.WriteString(off.Render("○"))
		}
		if i == len(allEfforts)-1 {
			break
		}
		seg := strings.Repeat("━", sliderGap-1)
		if i < int(selected) {
			track.WriteString(on.Render(seg))
		} else {
			track.WriteString(off.Render(seg))
		}
	}

	var labels strings.Builder
	col := 0
	for i, e := range allEfforts {
		name := e.key()
		start := max(col, i*sliderGap-len(name)/2)
		labels.WriteString(strings.Repeat(" ", start-col))
		if i == int(selected) {
			labels.WriteString(on.Bold(true).Render(name))
		} else {
			labels.WriteString(mutedStyle.Render(name))
		}
		col = start + len(name)
	}

	return track.String() + "\n" + labels.String()
}

// effortBody adalah isi modal /effort.
func effortBody(selected effort) string {
	on := lipgloss.NewStyle().Foreground(selected.color()).Bold(true)
	return strings.Join([]string{
		mutedStyle.Render("Atur seberapa dalam dan teliti model bekerja. Berlaku untuk semua paket."),
		"",
		effortSlider(selected),
		"",
		on.Render(selected.key()) + mutedStyle.Render(" · "+selected.description()),
	}, "\n")
}

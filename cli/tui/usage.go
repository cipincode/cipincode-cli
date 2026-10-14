package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Pemakaian dihitung server (Laravel) dengan dua jendela: sesi 5 jam yang reset
// otomatis, dan mingguan. CLI hanya menampilkan dan mencatat lewat API.

type usageMeter struct {
	Used      int        `json:"used"`
	Limit     int        `json:"limit"`
	Percent   int        `json:"percent"`
	Exhausted bool       `json:"exhausted"`
	ResetsAt  *time.Time `json:"resets_at"`
}

// usagePlan: paket gratis hanya punya batas sesi dan masa percobaan; berbayar punya sesi dan mingguan.
type usagePlan struct {
	Name      string     `json:"name"`
	Label     string     `json:"label"`
	ExpiresAt *time.Time `json:"expires_at"`
	Expired   bool       `json:"expired"`
}

type usageSnapshot struct {
	Session usageMeter  `json:"session"`
	Weekly  *usageMeter `json:"weekly"` // nil jika paket tanpa batas mingguan
	Plan    usagePlan   `json:"plan"`
	Blocked string      `json:"blocked,omitempty"`
	Message string      `json:"message,omitempty"`
}

func usageRequest(method, baseURL, key string, body any) (*http.Response, error) {
	var reader *bytes.Reader
	if body != nil {
		payload, _ := json.Marshal(body)
		reader = bytes.NewReader(payload)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, strings.TrimRight(baseURL, "/")+"/api/cli/usage", reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa terhubung ke %s", baseURL)
	}
	return resp, nil
}

func decodeUsage(resp *http.Response) (*usageSnapshot, error) {
	var snap usageSnapshot
	if err := json.NewDecoder(resp.Body).Decode(&snap); err != nil {
		return nil, fmt.Errorf("balasan server tidak valid: %w", err)
	}
	return &snap, nil
}

func fetchUsage(baseURL, key string) (*usageSnapshot, error) {
	resp, err := usageRequest(http.MethodGet, baseURL, key, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return decodeUsage(resp)
	case http.StatusUnauthorized:
		return nil, errInvalidKey
	}
	return nil, fmt.Errorf("server membalas status %d", resp.StatusCode)
}

// ---- tampilan ----

var (
	emptyStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#d4d4d8", Dark: "#3f3f46"})
	usedStyle  = lipgloss.NewStyle().Foreground(orange)
	fullStyle  = lipgloss.NewStyle().Foreground(red)
)

// filledSegments: 0 jika belum dipakai, bertambah seiring pemakaian (minimal 1
// begitu ada pemakaian), dan penuh hanya saat jatah benar-benar habis.
func filledSegments(used, limit, segments int) int {
	switch {
	case limit <= 0 || used <= 0:
		return 0
	case used >= limit:
		return segments
	}
	filled := (used*segments + limit - 1) / limit // pembulatan ke atas
	return min(filled, segments-1)
}

// usageBar menggambar kotak-kotak bersekat: kosong, oranye saat terpakai, merah saat habis.
func usageBar(m usageMeter, segments int) string {
	filled := filledSegments(m.Used, m.Limit, segments)
	on := usedStyle
	if m.Exhausted {
		on = fullStyle
	}

	cells := make([]string, segments)
	for i := range cells {
		if i < filled {
			cells[i] = on.Render("■")
		} else {
			cells[i] = emptyStyle.Render("■")
		}
	}
	return strings.Join(cells, " ")
}

func percentText(m usageMeter) string {
	text := fmt.Sprintf("%3d%%", m.Percent)
	if m.Exhausted {
		return fullStyle.Bold(true).Render(text)
	}
	return boldStyle.Render(text)
}

var hari = [...]string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}

// formatIn menulis sisa waktu singkat: "kurang dari 1 mnt", "42 mnt", "3 j 12 m", "4 hari 6 j".
func formatIn(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "kurang dari 1 mnt"
	case d < time.Hour:
		return fmt.Sprintf("%d mnt", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d j %d m", int(d.Hours()), int(d.Minutes())%60)
	}
	return fmt.Sprintf("%d hari %d j", int(d.Hours())/24, int(d.Hours())%24)
}

// formatResetAt: "pukul 19.45" untuk hari ini, selain itu "Sen 20 Okt, 00.00".
func formatResetAt(t, now time.Time) string {
	t = t.Local()
	clock := fmt.Sprintf("%02d.%02d", t.Hour(), t.Minute())
	y1, m1, d1 := t.Date()
	y2, m2, d2 := now.Local().Date()
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return "pukul " + clock
	}
	return fmt.Sprintf("%s %d %s, %s", hari[t.Weekday()], t.Day(), bulan[t.Month()-1], clock)
}

func meterDetail(m usageMeter, idle string, now time.Time) string {
	text := fmt.Sprintf("%d / %d pesan", m.Used, m.Limit)

	if m.ResetsAt == nil {
		return mutedStyle.Render(text + " · " + idle)
	}
	reset := fmt.Sprintf("reset dalam %s (%s)", formatIn(m.ResetsAt.Sub(now)), formatResetAt(*m.ResetsAt, now))
	if m.Exhausted {
		return fullStyle.Render("Habis · " + reset)
	}
	return mutedStyle.Render(text + " · " + reset)
}

func usageSegments(width int) int {
	if width < 70 {
		return 10
	}
	return 20
}

// planLine: nama paket dan masa berlakunya. Kosong jika server tidak mengirim data paket.
func planLine(p usagePlan, now time.Time) string {
	if p.Label == "" {
		return ""
	}
	name := boldStyle.Render("Paket " + p.Label)
	trial := p.Name == "free"

	switch {
	case p.Expired && trial:
		return name + fullStyle.Render(" · masa percobaan sudah berakhir")
	case p.Expired:
		return name + fullStyle.Render(" · masa aktif sudah berakhir")
	case p.ExpiresAt != nil:
		verb := "aktif, berakhir"
		if trial {
			verb = "masa percobaan berakhir"
		}
		return name + mutedStyle.Render(fmt.Sprintf(" · %s dalam %s (%s)", verb, formatIn(p.ExpiresAt.Sub(now)), formatResetAt(*p.ExpiresAt, now)))
	}
	return name + mutedStyle.Render(" · aktif tanpa batas waktu")
}

// usageBody menyusun isi /usage tanpa judul: paket, lalu batang kotak-kotak beserta rinciannya.
// Paket tanpa batas mingguan hanya menampilkan batang sesi.
func usageBody(s usageSnapshot, width int, now time.Time) string {
	segments := usageSegments(width)
	indent := strings.Repeat(" ", 14)

	row := func(label string, m usageMeter, idle string) string {
		return padRight(boldStyle.Render(label), 14) + usageBar(m, segments) + "  " + percentText(m) + "\n" +
			indent + meterDetail(m, idle, now)
	}

	var b strings.Builder
	if line := planLine(s.Plan, now); line != "" {
		b.WriteString(line + "\n\n")
	}
	b.WriteString(row("Sesi 5 jam", s.Session, "mulai saat pesan pertama"))
	if s.Weekly != nil {
		b.WriteString("\n\n" + row("Mingguan", *s.Weekly, "reset setiap Senin 00.00"))
	} else {
		label := s.Plan.Label
		if label == "" {
			label = "ini"
		}
		b.WriteString("\n\n" + mutedStyle.Render("Paket "+label+" tidak memiliki batas mingguan."))
	}
	return b.String()
}

func usageTitle(name string) string {
	if name == "" {
		return "Pemakaian"
	}
	return "Pemakaian · " + name
}

// renderUsage menyusun blok berjudul untuk dipasang di percakapan (mis. saat jatah habis).
func renderUsage(s usageSnapshot, name string, width int, now time.Time) string {
	return brandStyle.Render(usageTitle(name)) + "\n" + usageBody(s, width, now)
}

// usageHeaderLine adalah satu baris ringkas untuk header chat.
func usageHeaderLine(s *usageSnapshot) string {
	if s == nil {
		return mutedStyle.Render("Pemakaian: memuat…")
	}
	return mutedStyle.Render("Sesi 5 jam ") + usageBar(s.Session, 10) + " " + percentText(s.Session)
}

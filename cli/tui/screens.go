package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	greenStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#34d399"))
	linkStyle  = lipgloss.NewStyle().Foreground(orange).Underline(true)
	boldStyle  = lipgloss.NewStyle().Bold(true)
)

// ---- onboarding (animasi sambutan) ----

const onboardTitle = "Selamat datang di CipinCode"

// Garis waktu animasi (1 frame = 70 ms):
//
//	frame 0-12   logo muncul baris demi baris
//	frame 14-41  judul diketik huruf demi huruf
//	frame 44     subjudul muncul
//	frame 52     petunjuk berkedip, siap lanjut
func (m model) onboardingView() string {
	logo := renderLogoRows(m.frame / 2)

	titleRunes := []rune(onboardTitle)
	n := min(max(0, m.frame-14), len(titleRunes))
	title := lipgloss.NewStyle().Width(len(titleRunes)).Align(lipgloss.Center).
		Render(brandStyle.Render(string(titleRunes[:n])))

	subtitle := " "
	if m.frame >= 44 {
		subtitle = mutedStyle.Render("AI agent untuk terminal Anda")
	}

	hint := " "
	if m.frame >= onboardFrame {
		hint = mutedStyle.Render("Tekan Enter untuk masuk · Esc untuk lewati")
		if (m.frame/8)%2 == 1 {
			hint = lipgloss.NewStyle().Foreground(orange).Render("Tekan Enter untuk masuk · Esc untuk lewati")
		}
	}

	block := lipgloss.JoinVertical(lipgloss.Center, logo, "", title, subtitle, "", hint)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

func (m model) keyOnboarding(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		return m.enterChat(), textinput.Blink
	case tea.KeyEnter:
		if m.frame >= onboardFrame {
			m.screen, m.methodIndex, m.loginErr = screenMethod, 0, ""
			return m, nil
		}
	}
	// tombol apa pun melewati animasi
	if m.frame < onboardFrame {
		m.frame = onboardFrame
	}
	return m, nil
}

func (m model) enterChat() model {
	m.screen = screenChat
	m.layout()
	m.refresh()
	return m
}

// ---- pilih cara masuk ----

var methods = []struct{ icon, name, desc, key string }{
	{"✉", "Email", "masuk dengan email dan password akun dashboard", "email"},
	{"G", "Google", "masuk dengan akun Google Anda", "google"},
}

func (m model) methodView() string {
	rows := make([]string, 0, len(methods))
	for i, md := range methods {
		line := padRight(md.icon+"  "+md.name, 12) + md.desc
		if i == m.methodIndex {
			rows = append(rows, brandStyle.Render("❯ "+padRight(md.icon+"  "+md.name, 12))+md.desc)
		} else {
			rows = append(rows, "  "+mutedStyle.Render(line))
		}
	}

	hint := "↑↓ pilih · Enter lanjut · Esc lewati"
	if m.loginOnly {
		hint = "↑↓ pilih · Enter lanjut · Esc batal"
	}

	block := lipgloss.JoinVertical(lipgloss.Center,
		renderLogo(),
		"",
		boldStyle.Render("Pilih cara masuk"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
		"",
		mutedStyle.Render(hint),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

func (m model) keyMethod(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if m.loginOnly {
			return m, tea.Quit
		}
		return m.enterChat(), textinput.Blink
	case tea.KeyUp, tea.KeyDown, tea.KeyTab, tea.KeyShiftTab:
		m.methodIndex = (m.methodIndex + 1) % len(methods)
		return m, nil
	case tea.KeyEnter:
		return m.startCode()
	}

	switch msg.String() {
	case "k", "j":
		m.methodIndex = (m.methodIndex + 1) % len(methods)
	case "1":
		m.methodIndex = 0
		return m.startCode()
	case "2":
		m.methodIndex = 1
		return m.startCode()
	}
	return m, nil
}

// ---- URL dan tempel kode ----

func (m model) startCode() (tea.Model, tea.Cmd) {
	verifier, err := newVerifier()
	if err != nil {
		m.loginErr = "Gagal membuat kunci keamanan: " + err.Error()
		return m, nil
	}

	m.verifier = verifier
	m.baseURL = baseURLFor(m.urlOverride, m.config)
	m.authURL = authorizeURL(m.baseURL, challengeFor(verifier), methods[m.methodIndex].key)
	m.browser, m.loginErr, m.exchanging, m.loginOK = 0, "", false, false
	m.codeInput.SetValue("")
	m.codeInput.Focus()
	m.screen = screenCode

	return m, tea.Batch(textinput.Blink, openBrowserCmd(m.authURL))
}

func openBrowserCmd(target string) tea.Cmd {
	return func() tea.Msg { return browserMsg{err: openBrowser(target)} }
}

func exchangeCmd(baseURL, code, verifier string) tea.Cmd {
	return func() tea.Msg {
		key, user, err := exchangeCode(baseURL, code, verifier)
		return exchangeMsg{key: key, user: user, err: err}
	}
}

func (m model) keyCode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if !m.exchanging && !m.loginOK {
			m.screen, m.loginErr = screenMethod, ""
		}
		return m, nil
	case tea.KeyCtrlO:
		return m, openBrowserCmd(m.authURL)
	case tea.KeyEnter:
		code := strings.TrimSpace(m.codeInput.Value())
		if code == "" || m.exchanging || m.loginOK {
			return m, nil
		}
		m.exchanging, m.loginErr = true, ""
		return m, exchangeCmd(m.baseURL, code, m.verifier)
	}

	if m.exchanging || m.loginOK {
		return m, nil
	}
	var cmd tea.Cmd
	m.codeInput, cmd = m.codeInput.Update(msg)
	return m, cmd
}

func (m model) handleExchange(msg exchangeMsg) (tea.Model, tea.Cmd) {
	m.exchanging = false
	if msg.err != nil {
		m.loginErr = msg.err.Error()
		return m, nil
	}

	cfg, err := saveLogin(m.baseURL, msg.key, msg.user.Name, msg.user.Email)
	if err != nil {
		m.loginErr = "Gagal menyimpan login: " + err.Error()
		return m, nil
	}

	m.config = &cfg
	m.loginOK = true
	return m, tea.Tick(1400*time.Millisecond, func(time.Time) tea.Msg { return loginDoneMsg{} })
}

func (m model) finishLogin() (tea.Model, tea.Cmd) {
	if m.loginOnly {
		return m, tea.Quit
	}
	m.home = true
	cmds := []tea.Cmd{textinput.Blink}
	if m.config != nil {
		cmds = append(cmds, fetchUsageCmd(m.config, false))
	}
	return m.enterChat(), tea.Batch(cmds...)
}

func (m model) codeView() string {
	// Lebar mengikuti panjang URL agar tidak terpotong ke baris baru (mudah disalin).
	width := min(max(m.width-4, 30), max(84, lipgloss.Width(m.authURL)))
	wrap := lipgloss.NewStyle().Width(width)
	method := methods[m.methodIndex]

	var b strings.Builder
	b.WriteString(brandStyle.Render("Masuk dengan " + method.name))
	b.WriteString("\n\n")

	if m.loginOK {
		b.WriteString(greenStyle.Render("✓ Berhasil masuk sebagai " + m.config.Name))
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(m.config.Email))
		return m.framed(b.String(), width)
	}

	b.WriteString(boldStyle.Render("1. Buka URL ini di browser"))
	b.WriteString("\n")
	b.WriteString(wrap.Render(linkStyle.Render(m.authURL)))
	b.WriteString("\n")
	switch m.browser {
	case 1:
		b.WriteString(mutedStyle.Render("Browser dibuka otomatis. Ctrl+O untuk membukanya lagi."))
	case 2:
		b.WriteString(mutedStyle.Render("Tidak bisa membuka browser otomatis. Salin URL di atas secara manual."))
	default:
		b.WriteString(mutedStyle.Render("Membuka browser…"))
	}
	b.WriteString("\n\n")

	b.WriteString(boldStyle.Render("2. Masuk dengan " + method.name + ", lalu klik “Izinkan”"))
	b.WriteString("\n\n")

	b.WriteString(boldStyle.Render("3. Tempel kode yang muncul di sini"))
	b.WriteString("\n")
	b.WriteString(inputBorder.BorderForeground(orange).Render(m.codeInput.View()))
	b.WriteString("\n")

	switch {
	case m.exchanging:
		b.WriteString(m.spinner.View() + mutedStyle.Render(" Memverifikasi kode…"))
	case m.loginErr != "":
		b.WriteString(wrap.Render(errorStyle.Render("✗ " + m.loginErr)))
	default:
		b.WriteString(mutedStyle.Render("Kode berlaku 5 menit · Enter kirim · Esc kembali"))
	}
	return m.framed(b.String(), width)
}

func (m model) framed(content string, width int) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center,
			renderLogo(),
			"",
			lipgloss.NewStyle().Width(width).Render(content),
		))
}

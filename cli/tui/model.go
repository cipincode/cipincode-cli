package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type role int

const (
	roleUser role = iota
	roleAssistant
	roleSystem
	roleError
	roleRich // sudah berwarna, ditampilkan apa adanya
)

type message struct {
	role    role
	content string

	// cache hasil Markdown agar tidak dihitung ulang tiap token
	rendered        string
	renderedContent string
	renderedWidth   int
}

type screen int

const (
	screenOnboarding screen = iota // animasi sambutan
	screenMethod                   // pilih email atau Google
	screenCode                     // buka URL, tempel kode
	screenChat
)

type options struct {
	screen    screen
	loginOnly bool   // keluar setelah login (cipincode login)
	url       string // alamat dashboard yang dipilih lewat --url (opsional)
}

type (
	tokenMsg     string
	doneMsg      struct{}
	errMsg       struct{ err error }
	animTickMsg  struct{}
	loginDoneMsg struct{}
	browserMsg   struct{ err error }
	usageMsg     struct {
		snap     *usageSnapshot
		err      error
		announce bool // tampilkan hasilnya di percakapan (/usage)
	}
	usageStreamMsg struct{ snap *usageSnapshot }
	toastClearMsg  struct{ id int }
	blockedMsg     struct {
		snap    *usageSnapshot
		message string
	}
	exchangeMsg struct {
		key  string
		user *userInfo
		err  error
	}
)

const (
	inputRows    = 3 // kotak input bergaris
	gapRows      = 1 // jarak kosong antara area chat dan label mode
	labelRows    = 1 // label mode di atas input
	statusRows   = 1
	maxMenuRows  = 8
	onboardFrame = 52 // frame saat animasi onboarding selesai
	frameTime    = 70 * time.Millisecond
)

type model struct {
	agent       *Agent
	config      *Config
	lastSession time.Time
	screen      screen
	loginOnly   bool
	urlOverride string

	// chat
	mode           mode
	effort         effort
	sessionID      string
	chosen         *agentItem // agent buatan pengguna yang dipilih lewat /agent
	home           bool
	messages       []message
	input          textinput.Model
	viewport       viewport.Model
	spinner        spinner.Model
	stream         chan tea.Msg
	cancel         context.CancelFunc
	streaming      bool
	menuIndex      int
	usage          *usageSnapshot
	pending        string // pesan yang sedang diproses, dikembalikan ke input jika diblokir
	modal          *modalView
	sel            *selection
	agentChoices   []agentItem
	agentCursor    int
	provider       string // ID provider API ("": Otomatis ikut paket)
	modelID        string // ID model yang dipakai ("": bawaan provider)
	modelChoices   []providerModel // daftar live dari akun provider (nil = katalog bawaan)
	modelsLive     bool             // modelChoices berasal dari fetch akun
	providerCursor int
	modelCursor    int
	claudeCursor   int             // pilihan di modal metode Anthropic
	connectPending *provider       // provider yang menunggu kunci API di modal
	connectInput   textinput.Model // kolom kunci API di modal connect
	updateInfo     *versionInfo    // info update yang sedang ditawarkan
	updateSkip     bool            // checkbox "jangan ingatkan lagi"
	updateManual   bool            // /update manual: laporkan juga bila sudah terbaru
	toast          string
	toastID        int

	// onboarding dan login
	frame       int
	methodIndex int // 0 email, 1 google
	baseURL     string
	verifier    string
	authURL     string
	browser     int // 0 belum diketahui, 1 terbuka, 2 gagal
	codeInput   textinput.Model
	exchanging  bool
	loginErr    string
	loginOK     bool

	width  int
	height int
	ready  bool
}

func newModel(agent *Agent, cfg *Config, lastSession time.Time, opts options) model {
	in := textinput.New()
	in.Prompt = "❯ "
	in.PromptStyle = userStyle
	in.Placeholder = modeChat.placeholder()
	in.Focus()

	code := textinput.New()
	code.Prompt = "Kode ❯ "
	code.PromptStyle = userStyle
	code.Placeholder = "XXXX-XXXX"
	code.CharLimit = 24
	code.Width = 24

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = lipgloss.NewStyle().Foreground(orange)

	return model{
		agent:       agent,
		config:      cfg,
		lastSession: lastSession,
		screen:      opts.screen,
		loginOnly:   opts.loginOnly,
		urlOverride: opts.url,
		mode:        modeChat,
		effort:      loadEffort(),
		provider:    loadProvider(),
		modelID:     loadModel(),
		sessionID:   newUUID(),
		home:        true,
		input:       in,
		codeInput:   code,
		spinner:     sp,
		baseURL:     baseURLFor(opts.url, cfg),
	}
}

// baseURLFor memilih alamat dashboard: pilihan eksplisit, lalu login sebelumnya, env, bawaan.
func baseURLFor(override string, cfg *Config) string {
	if override != "" {
		return override
	}
	return loginBaseURL(cfg)
}

func animTick() tea.Cmd {
	return tea.Tick(frameTime, func(time.Time) tea.Msg { return animTickMsg{} })
}

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}
	switch m.screen {
	case screenOnboarding:
		cmds = append(cmds, animTick())
	case screenChat:
		cmds = append(cmds, textinput.Blink)
		cmds = append(cmds, fetchVersionCmd(m.baseURL))
		if m.config != nil {
			cmds = append(cmds, fetchUsageCmd(m.config, false))
		}
	}
	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()
		m.refresh()
		return m, nil

	case animTickMsg:
		if m.screen != screenOnboarding {
			return m, nil
		}
		m.frame++
		return m, animTick()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.streaming {
			m.refresh()
		}
		return m, cmd

	case claudeLoginMsg:
		return m.handleClaudeLogin(msg)
	case browserMsg:
		m.browser = 1
		if msg.err != nil {
			m.browser = 2
		}
		return m, nil

	case exchangeMsg:
		return m.handleExchange(msg)

	case loginDoneMsg:
		return m.finishLogin()

	case usageMsg:
		return m.handleUsage(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case toastClearMsg:
		if msg.id == m.toastID {
			m.toast = ""
		}
		return m, nil

	case agentsMsg:
		return m.handleAgents(msg)

	case claudeUsageMsg:
		return m.handleClaudeUsage(msg)

	case modelsMsg:
		return m.handleModels(msg)

	case versionMsg:
		return m.handleVersion(msg)

	case updateDoneMsg:
		return m.handleUpdateDone(msg)

	case usageStreamMsg:
		m.usage = msg.snap
		return m, waitForStream(m.stream)

	case blockedMsg:
		return m.handleBlocked(msg)

	case tokenMsg:
		m.appendToAssistant(string(msg))
		m.refresh()
		return m, waitForStream(m.stream)

	case doneMsg:
		m.streaming = false
		m.cancel = nil
		m.refresh()
		return m, nil

	case errMsg:
		m.streaming = false
		m.cancel = nil
		if n := len(m.messages); n > 0 && m.messages[n-1].role == roleAssistant && m.messages[n-1].content == "" {
			m.messages = m.messages[:n-1] // jawaban kosong tidak perlu ditampilkan
		}
		m.messages = append(m.messages, message{role: roleError, content: "Error: " + msg.err.Error()})
		m.refresh()
		return m, nil

	case tea.KeyMsg:
		switch m.screen {
		case screenOnboarding:
			return m.keyOnboarding(msg)
		case screenMethod:
			return m.keyMethod(msg)
		case screenCode:
			return m.keyCode(msg)
		}
		if m.modal != nil {
			return m.keyModal(msg)
		}
		return m.keyChat(msg)
	}

	// Pesan lain (mis. kedip kursor) diteruskan ke kolom input yang aktif.
	var cmd tea.Cmd
	switch m.screen {
	case screenChat:
		m.input, cmd = m.input.Update(msg)
	case screenCode:
		m.codeInput, cmd = m.codeInput.Update(msg)
	}
	return m, cmd
}

// ---- chat ----

func (m model) keyChat(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.sel = nil // mengetik menghapus sorotan
	if items := m.menuItems(); len(items) > 0 {
		switch msg.Type {
		case tea.KeyUp:
			m.menuIndex = (m.menuIndex - 1 + len(items)) % len(items)
			return m, nil
		case tea.KeyDown:
			m.menuIndex = (m.menuIndex + 1) % len(items)
			return m, nil
		case tea.KeyTab:
			m.input.SetValue("/" + items[m.menuIndex].name + " ")
			m.input.CursorEnd()
			m.menuIndex = 0
			m.layout()
			return m, nil
		case tea.KeyEnter:
			name := items[m.menuIndex].name
			m.input.Reset()
			m.menuIndex = 0
			m.layout()
			return m.runCommand(name)
		case tea.KeyEsc:
			m.input.Reset()
			m.menuIndex = 0
			m.layout()
			return m, nil
		}
	}

	switch msg.Type {
	case tea.KeyCtrlC:
		if m.cancel != nil {
			m.cancel()
		}
		return m, tea.Quit
	case tea.KeyShiftTab:
		m.mode = m.mode.next()
		m.input.Placeholder = m.mode.placeholder()
		return m, nil
	case tea.KeyEsc:
		if m.streaming && m.cancel != nil {
			m.cancel()
		}
		return m, nil
	case tea.KeyEnter:
		return m.submit()
	case tea.KeyPgUp, tea.KeyPgDown:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case tea.KeyUp:
		m.viewport.ScrollUp(1)
		return m, nil
	case tea.KeyDown:
		m.viewport.ScrollDown(1)
		return m, nil
	}

	before := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != before {
		m.menuIndex = 0
	}
	m.layout()
	return m, cmd
}

func (m model) submit() (tea.Model, tea.Cmd) {
	text := strings.TrimSpace(m.input.Value())
	if text == "" || m.streaming {
		return m, nil
	}

	if strings.HasPrefix(text, "/") {
		m.input.Reset()
		name := strings.TrimPrefix(strings.Fields(text)[0], "/")
		return m.runCommand(name)
	}

	// AI berjalan lewat server dashboard, jadi butuh login. Teks tetap di kolom input.
	if m.config == nil {
		m = m.note(roleError, "Masuk dulu dengan /login untuk memakai AI.")
		m.layout()
		m.refresh()
		return m, nil
	}
	m.input.Reset()

	cwd, _ := os.Getwd()
	request := chatRequest{
		Message:     text,
		Mode:        m.mode.key(),
		Effort:      m.effort.key(),
		Provider:    m.provider,
		Model:       m.modelID,
		ProviderKey: m.providerKey(),
		Session:     m.sessionID,
		Agent:       m.activeAgentID(),
		History:     m.history(),
		Cwd:         cwd,
		URL:         m.config.URL,
		Key:         m.config.Key,
	}

	m.home = false
	m.pending = text
	m.messages = append(m.messages,
		message{role: roleUser, content: text},
		message{role: roleAssistant, content: ""},
	)
	m.streaming = true
	m.layout()
	m.refresh()

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.stream = make(chan tea.Msg, 64)
	go streamChat(ctx, m.agent, request, m.useClaudeAccount(), m.stream)

	return m, waitForStream(m.stream)
}

const historyLimit = 20

// history adalah percakapan sebelumnya (tanpa pesan sistem dan galat) untuk konteks model.
func (m model) history() []historyItem {
	var items []historyItem
	for _, msg := range m.messages {
		switch {
		case msg.role == roleUser && msg.content != "":
			items = append(items, historyItem{"user", msg.content})
		case msg.role == roleAssistant && msg.content != "":
			items = append(items, historyItem{"assistant", msg.content})
		}
	}
	if len(items) > historyLimit {
		items = items[len(items)-historyLimit:]
	}
	return items
}

func fetchUsageCmd(cfg *Config, announce bool) tea.Cmd {
	return func() tea.Msg {
		snap, err := fetchUsage(cfg.URL, cfg.Key)
		return usageMsg{snap: snap, err: err, announce: announce}
	}
}

func (m model) handleUsage(msg usageMsg) (tea.Model, tea.Cmd) {
	if msg.err == nil && msg.snap != nil {
		m.usage = msg.snap
	}
	if !msg.announce {
		return m, nil
	}

	if msg.err != nil {
		m = m.note(roleError, usageErrorText(msg.err))
	} else {
		m = m.openModal(usageTitle(m.accountName()), usageBody(*msg.snap, m.width, time.Now()))
	}
	m.layout()
	m.refresh()
	return m, nil
}

// handleBlocked: server menolak karena jatah habis. Pesan tidak terkirim, jadi
// pertanyaan dan jawaban kosongnya dihapus dan teks dikembalikan ke kolom input.
func (m model) handleBlocked(msg blockedMsg) (tea.Model, tea.Cmd) {
	if n := len(m.messages); n >= 2 {
		m.messages = m.messages[:n-2]
	}
	m.streaming = false
	m.cancel = nil
	m.input.SetValue(m.pending)
	m.input.CursorEnd()

	reason := msg.message
	if reason == "" {
		reason = "Batas pemakaian tercapai."
	}
	m = m.note(roleError, "✗ "+reason+" Pesan Anda tidak dikirim.")
	if msg.snap != nil {
		m.usage = msg.snap
		m = m.note(roleRich, renderUsage(*msg.snap, m.accountName(), m.width, time.Now()))
	}

	m.layout()
	m.refresh()
	return m, nil
}

func usageErrorText(err error) string {
	if errors.Is(err, errInvalidKey) {
		return "Login Anda tidak valid lagi. Ketik /login untuk masuk kembali."
	}
	return "Tidak bisa memeriksa pemakaian: " + err.Error()
}

func (m model) accountName() string {
	if m.config != nil {
		return m.config.Name
	}
	return ""
}

func (m model) note(r role, text string) model {
	m.home = false
	m.messages = append(m.messages, message{role: r, content: text})
	return m
}

func (m model) runCommand(name string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch name {
	case "exit", "quit":
		return m, tea.Quit
	case "clear":
		m.messages, m.home, m.sessionID = nil, false, newUUID()
	case "new":
		m.messages, m.home, m.sessionID = nil, true, newUUID()
	case "help":
		m = m.openModal("Bantuan", helpText())
	case "chat", "coding", "debug":
		m.mode = map[string]mode{"chat": modeChat, "coding": modeCoding, "debug": modeDebug}[name]
		m.input.Placeholder = m.mode.placeholder()
		m = m.note(roleSystem, fmt.Sprintf("Mode %s · %s", m.mode.name(), m.mode.description()))
	case "effort":
		m = m.openEffortModal()
	case "connect":
		m.providerCursor = 0
		for i, p := range providers {
			if p.id == m.provider {
				m.providerCursor = i + 1
			}
		}
		m = m.openProvidersModal()
	case "model":
		if parseProvider(m.provider) == nil {
			m = m.note(roleSystem, "Belum ada provider. Pilih dulu dengan /connect.")
		} else if m.modelsLiveOK() {
			// Ada kunci sendiri: muat daftar live dari akun provider.
			m = m.note(roleSystem, "Memuat daftar model dari akun "+providerName(m.provider)+"…")
			m.layout()
			m.refresh()
			return m, fetchModelsCmd(m.provider, m.providerKey())
		} else {
			// Via server atau akun Claude: memakai katalog bawaan.
			m.modelChoices, m.modelsLive = nil, false
			models := providerModels(m.provider)
			m.modelCursor = 0
			for i, mdl := range models {
				if mdl.id == m.modelID {
					m.modelCursor = i
				}
			}
			m = m.openModelsModal()
		}
	case "agent":
		if m.config == nil {
			m = m.note(roleSystem, "Masuk dulu dengan /login untuk memakai agent.")
		} else {
			return m, fetchAgentsCmd(m.config)
		}
	case "usage":
		if m.config == nil {
			m = m.note(roleSystem, "Masuk dulu dengan /login untuk melihat pemakaian.")
		} else if m.directMode() && m.useClaudeAccount() {
			return m, claudeUsageCmd()
		} else if m.directMode() {
			m = m.openModal("Pemakaian · "+providerName(m.provider), m.directUsageBody(nil))
		} else {
			return m, fetchUsageCmd(m.config, true)
		}
	case "version":
		m = m.openModal("Versi", "CipinCode CLI v"+version+"\n\n"+mutedStyle.Render("Tool AI terminal interaktif untuk developer."))
	case "update":
		m.updateManual = true
		m = m.note(roleSystem, "Memeriksa pembaruan…")
		m.layout()
		m.refresh()
		return m, fetchVersionCmd(m.baseURL)
	case "whoami":
		if m.config == nil {
			m = m.note(roleSystem, "Belum login. Ketik /login untuk masuk.")
		} else {
			m = m.note(roleSystem, fmt.Sprintf("%s (%s)\n%s", m.config.Name, m.config.Email, m.config.URL))
		}
	case "status":
		m = m.openModal("Status", m.statusText())
	case "login":
		m.screen, m.methodIndex, m.loginErr = screenMethod, 0, ""
		return m, nil
	case "logout":
		if m.config == nil {
			m = m.note(roleSystem, "Anda belum login.")
		} else if err := deleteConfig(); err != nil {
			m = m.note(roleError, "Gagal keluar: "+err.Error())
		} else {
			m.config = nil
			m = m.note(roleSystem, "✓ Berhasil keluar. Ketik /login untuk masuk lagi.")
		}
	case "dashboard", "keys", "setup":
		path := map[string]string{"dashboard": "/dashboard", "keys": "/api-keys", "setup": "/setup"}[name]
		target := strings.TrimRight(baseURLFor(m.urlOverride, m.config), "/") + path
		m = m.note(roleSystem, "Membuka "+target)
		cmd = func() tea.Msg { _ = openBrowser(target); return nil }
	default:
		m = m.note(roleError, "Perintah tidak dikenal: /"+name+" (ketik / untuk melihat daftar)")
	}

	m.layout()
	m.refresh()
	return m, cmd
}

func helpText() string {
	var b strings.Builder
	b.WriteString("Perintah (ketik / untuk membuka menu):\n")
	for _, c := range commands {
		b.WriteString("  " + padRight("/"+c.name, 12) + c.desc + "\n")
	}
	b.WriteString("\nMode (Shift+Tab untuk berganti): Chat · Coding · Debugging\n")
	b.WriteString("Effort (/effort, lalu ← →): low · medium · high · xhigh · ultra\n")
	b.WriteString("Provider (/connect): Otomatis atau Anthropic · OpenAI · Google · xAI · DeepSeek · Meta\n")
	b.WriteString("Model (/model): daftar model dari provider yang aktif\n")
	b.WriteString("Agent (/agent): pilih agent buatan Anda dari menu Agents di dashboard\n")
	b.WriteString("Tombol: Enter kirim · Esc batalkan · PgUp/PgDn gulir · Ctrl+C keluar")
	return b.String()
}

func (m model) statusText() string {
	account := "belum login"
	if m.config != nil {
		account = fmt.Sprintf("%s (%s)", m.config.Name, m.config.Email)
	}
	return strings.Join([]string{
		"Akun        " + account,
		"Dashboard   " + baseURLFor(m.urlOverride, m.config),
		"Agent       " + m.agent.URL,
		"Mode        " + m.mode.name() + " · " + m.mode.description(),
		"Effort      " + m.effort.key() + " · " + m.effort.description(),
		"Provider    " + providerName(m.provider),
		"Model       " + modelName(m.modelID),
		"API         " + m.apiModeText(),
		"Agent       " + m.agentName(),
		"Sesi lalu   " + formatSession(m.lastSession),
		"Versi       " + version,
	}, "\n") + m.usageSummary()
}

func (m model) agentName() string {
	if m.chosen == nil {
		return "tidak ada (/agent untuk memilih)"
	}
	return m.chosen.Name
}

// apiModeText menjelaskan jalur API yang dipakai: langsung dengan kunci
// milik pengguna, atau via server CipinCode.
func (m model) apiModeText() string {
	if m.provider == "" {
		return "via CipinCode (ikut paket)"
	}
	if m.useClaudeAccount() {
		return "akun Claude (Pro/Max) lewat Claude Code"
	}
	if m.providerKey() != "" {
		return "langsung ke " + providerName(m.provider) + " (kunci Anda)"
	}
	return "via server CipinCode"
}

func (m model) usageSummary() string {
	if m.directMode() {
		return "\nPemakaian   langsung via " + providerName(m.provider) + ", jatah server tidak terpotong (ketik /usage)"
	}
	if m.usage == nil {
		return ""
	}
	summary := fmt.Sprintf("sesi %d%%", m.usage.Session.Percent)
	if m.usage.Weekly != nil {
		summary += fmt.Sprintf(" · mingguan %d%%", m.usage.Weekly.Percent)
	}
	plan := ""
	if m.usage.Plan.Label != "" {
		plan = "\nPaket       " + m.usage.Plan.Label
	}
	return plan + "\nPemakaian   " + summary + " (ketik /usage)"
}

func (m *model) appendToAssistant(text string) {
	if n := len(m.messages); n > 0 && m.messages[n-1].role == roleAssistant {
		m.messages[n-1].content += text
	}
}

// menuItems adalah perintah yang cocok dengan teks input (kosong jika menu tertutup).
func (m model) menuItems() []command {
	if m.screen != screenChat || m.streaming || m.modal != nil {
		return nil
	}
	return matchCommands(m.input.Value())
}

func (m model) menuRows() int {
	n := len(m.menuItems())
	if n == 0 {
		return 0
	}
	rows := min(n, maxMenuRows) + 2
	if n > maxMenuRows {
		rows++ // baris petunjuk "n lagi"
	}
	return rows
}

// compact dipakai di terminal kecil: header satu baris tanpa logo.
func (m model) compact() bool { return m.width < 56 || m.height < 24 }

func (m model) headerRows() int {
	if m.compact() {
		return 2
	}
	return len(logoLines) + 1
}

func (m *model) layout() {
	h := m.height - m.headerRows() - gapRows - labelRows - inputRows - statusRows - m.menuRows()
	if h < 3 {
		h = 3
	}
	if !m.ready {
		m.viewport = viewport.New(m.width, h)
		m.ready = true
	} else {
		m.viewport.Width, m.viewport.Height = m.width, h
	}
	m.input.Width = m.width - 6
}

func (m *model) refresh() {
	if !m.ready {
		return
	}
	wrap := lipgloss.NewStyle().Width(max(m.width-2, 10))
	var b strings.Builder

	if m.home {
		b.WriteString(m.homeView())
		b.WriteString("\n\n")
	}

	for i, msg := range m.messages {
		switch msg.role {
		case roleUser:
			b.WriteString(wrap.Render(userStyle.Render("❯ ") + msg.content))
		case roleAssistant:
			if msg.content == "" && m.streaming && i == len(m.messages)-1 {
				b.WriteString(brandStyle.Render("● ") + m.spinner.View() + mutedStyle.Render(" berpikir…"))
				break
			}
			b.WriteString(m.assistantText(i))
		case roleSystem:
			b.WriteString(wrap.Render(mutedStyle.Render(msg.content)))
		case roleError:
			b.WriteString(wrap.Render(errorStyle.Render(msg.content)))
		case roleRich:
			b.WriteString(msg.content)
		}
		b.WriteString("\n\n")
	}

	atBottom := m.viewport.AtBottom() || m.viewport.TotalLineCount() == 0
	m.viewport.SetContent(strings.TrimRight(b.String(), "\n"))
	if atBottom {
		m.viewport.GotoBottom()
	}
}

// assistantText merender jawaban AI sebagai Markdown, dengan penanda "●" di baris pertama dan
// jorokan di baris berikutnya. Hasilnya di-cache selama isi dan lebar tidak berubah.
func (m *model) assistantText(i int) string {
	msg := &m.messages[i]
	width := max(m.width-2, 20)

	if msg.rendered == "" || msg.renderedContent != msg.content || msg.renderedWidth != width {
		lines := strings.Split(renderMarkdown(msg.content, width), "\n")
		for j := range lines {
			if j == 0 {
				lines[j] = brandStyle.Render("● ") + lines[j]
			} else {
				lines[j] = "  " + lines[j]
			}
		}
		msg.rendered = strings.Join(lines, "\n")
		msg.renderedContent, msg.renderedWidth = msg.content, width
	}
	return msg.rendered
}

func (m model) homeView() string {
	greeting := "Selamat datang"
	if m.config != nil && m.config.Name != "" {
		greeting += ", " + firstName(m.config.Name)
	}

	var b strings.Builder
	b.WriteString(brandStyle.Render(greeting + "."))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("Pilih mode dengan Shift+Tab, ketik / untuk menu perintah, lalu tulis pesan di bawah."))
	b.WriteString("\n\n")
	for _, md := range []mode{modeChat, modeCoding, modeDebug} {
		style := lipgloss.NewStyle().Foreground(md.color()).Bold(true)
		b.WriteString(style.Render(padRight(md.icon()+" "+md.name(), 14)))
		b.WriteString(mutedStyle.Render(md.description()))
		b.WriteString("\n")
	}
	if m.config == nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render("Belum login.") + mutedStyle.Render(" Ketik /login untuk masuk."))
	}
	return b.String()
}

func (m model) View() string {
	if !m.ready {
		return "Memuat…"
	}
	switch m.screen {
	case screenOnboarding:
		return m.onboardingView()
	case screenMethod:
		return m.methodView()
	case screenCode:
		return m.codeView()
	}

	parts := []string{m.header(), m.viewportView(), ""} // baris kosong: jarak ke label mode
	if items := m.menuItems(); len(items) > 0 {
		parts = append(parts, m.menuView(items))
	}
	parts = append(parts, m.footer())
	view := strings.Join(parts, "\n")
	if m.modal != nil {
		view = overlay(view, m.modal.render(m.width, m.height), m.width)
	}
	return view
}

// viewportView adalah isi percakapan, dengan teks terpilih disorot.
func (m model) viewportView() string {
	view := m.viewport.View()
	if m.sel == nil || m.modal != nil {
		return view
	}
	return strings.Join(highlightLines(strings.Split(view, "\n"), *m.sel), "\n")
}

func (m model) header() string {
	name, detail := "Belum login", "ketik /login untuk masuk"
	if m.config != nil {
		name = m.config.Name
		detail = m.config.Email + " · API key"
	}

	if m.compact() {
		line := brandStyle.Render("CP") + mutedStyle.Render(" · "+version+" · "+name)
		return line + "\n" + mutedStyle.Render(strings.Repeat("─", m.width))
	}

	info := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render(name),
		mutedStyle.Render(detail),
		mutedStyle.Render("Sesi terakhir: "+formatSession(m.lastSession)),
		mutedStyle.Render("Agent: "+m.agent.URL),
		mutedStyle.Render("Versi "+version),
		m.usageLine(),
	)
	block := lipgloss.JoinHorizontal(lipgloss.Center, renderLogo(), "   ", info)
	return block + "\n" + mutedStyle.Render(strings.Repeat("─", m.width))
}

// usageLine: ringkasan pemakaian di header; kosong jika belum login.
// Mode langsung menampilkan akun provider, bukan batang kuota server.
func (m model) usageLine() string {
	if m.config == nil {
		return ""
	}
	if m.directMode() {
		via := "kunci Anda"
		if m.useClaudeAccount() {
			via = "akun Claude"
		}
		return mutedStyle.Render("Langsung · " + providerName(m.provider) + " (" + via + ") · tanpa jatah server")
	}
	return usageHeaderLine(m.usage)
}

func (m model) menuView(items []command) string {
	start := 0
	if m.menuIndex >= maxMenuRows {
		start = m.menuIndex - maxMenuRows + 1
	}
	end := min(start+maxMenuRows, len(items))

	selected := lipgloss.NewStyle().Bold(true).Foreground(orange)
	lines := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		name := padRight("/"+items[i].name, 14)
		if i == m.menuIndex {
			lines = append(lines, selected.Render("❯ "+name)+items[i].desc)
		} else {
			lines = append(lines, "  "+name+mutedStyle.Render(items[i].desc))
		}
	}
	if more := len(items) - (end - start); more > 0 {
		lines = append(lines, mutedStyle.Render(fmt.Sprintf("  ↕ %d perintah lagi, gunakan ↑↓", more)))
	}
	return inputBorder.Width(m.width - 2).Render(strings.Join(lines, "\n"))
}

func (m model) footer() string {
	md := m.mode
	accent := lipgloss.NewStyle().Foreground(md.color())

	left := accent.Bold(true).Render(md.icon()+" "+md.name()+" mode") + mutedStyle.Render(" · "+md.description())
	if m.chosen != nil {
		left += mutedStyle.Render(" · agent ") + lipgloss.NewStyle().Bold(true).Render(m.chosen.Name)
	}
	right := lipgloss.NewStyle().Foreground(m.effort.color()).Render("effort "+m.effort.key()) +
		mutedStyle.Render(" · shift+tab ganti mode")
	if m.toast != "" { // notifikasi sementara menggantikan petunjuk di sebelah kanan
		right = greenStyle.Render(m.toast)
	}
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	label := left
	if gap > 1 {
		label += strings.Repeat(" ", gap) + right
	}

	box := inputBorder.BorderForeground(md.color()).Width(m.width - 2).Render(m.input.View())

	status := mutedStyle.Render("Enter kirim · / menu perintah · Esc batal")
	switch {
	case m.streaming:
		status = m.spinner.View() + mutedStyle.Render(" menjawab · Esc untuk membatalkan")
	case len(m.menuItems()) > 0:
		status = mutedStyle.Render("↑↓ pilih · Tab lengkapi · Enter jalankan · Esc tutup")
	}
	return label + "\n" + box + "\n" + status
}

func firstName(full string) string {
	if fields := strings.Fields(full); len(fields) > 0 {
		return fields[0]
	}
	return full
}

func padRight(s string, width int) string {
	if w := lipgloss.Width(s); w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}

func streamChat(ctx context.Context, agent *Agent, request chatRequest, claude bool, out chan<- tea.Msg) {
	defer close(out)

	events := make(chan chatEvent, 64)
	errc := make(chan error, 1)
	go func() {
		if claude {
			errc <- chatClaude(ctx, request, events)
		} else {
			errc <- agent.Chat(ctx, request, events)
		}
		close(events)
	}()

	for ev := range events {
		switch ev.kind {
		case "token":
			out <- tokenMsg(ev.text)
		case "usage":
			out <- usageStreamMsg{ev.snap}
		case "blocked":
			out <- blockedMsg{snap: ev.snap, message: ev.message}
			<-errc
			return
		}
	}
	if err := <-errc; err != nil && ctx.Err() == nil {
		out <- errMsg{err}
		return
	}
	out <- doneMsg{}
}

func waitForStream(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return doneMsg{}
		}
		return msg
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// versionInfo adalah versi terbaru + changelog dari GET /api/cli/version.
type versionInfo struct {
	Version string   `json:"version"`
	Added   []string `json:"added"`
	Fixed   []string `json:"fixed"`
}

type versionMsg struct {
	info *versionInfo
	err  error
}

type updateDoneMsg struct {
	version string
	err     error
}

func fetchVersionCmd(baseURL string) tea.Cmd {
	return func() tea.Msg {
		info, err := fetchVersion(baseURL)
		return versionMsg{info: info, err: err}
	}
}

func fetchVersion(baseURL string) (*versionInfo, error) {
	base := strings.TrimRight(baseURL, "/")
	if base == "" {
		base = "https://cli.cipincode.com"
	}
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Get(base + "/api/cli/version")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server membalas status %d", resp.StatusCode)
	}
	var info versionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if strings.TrimSpace(info.Version) == "" {
		return nil, fmt.Errorf("balasan server tidak memuat versi")
	}
	return &info, nil
}

// compareVersions membandingkan "1.2.3" secara numerik per bagian.
func compareVersions(a, b string) int {
	pa := strings.Split(strings.TrimPrefix(strings.TrimSpace(a), "v"), ".")
	pb := strings.Split(strings.TrimPrefix(strings.TrimSpace(b), "v"), ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var na, nb int
		if i < len(pa) {
			na, _ = strconv.Atoi(strings.TrimSpace(pa[i]))
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(strings.TrimSpace(pb[i]))
		}
		if na < nb {
			return -1
		}
		if na > nb {
			return 1
		}
	}
	return 0
}

// isNewer melaporkan apakah latest lebih baru dari have.
func isNewer(have, latest string) bool { return compareVersions(have, latest) < 0 }

func updateModalBody(info *versionInfo, skip bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Versi %s → %s tersedia.\n\n", version, info.Version)
	if len(info.Added) > 0 {
		b.WriteString("Yang baru:\n")
		for _, s := range info.Added {
			b.WriteString("  • " + s + "\n")
		}
		b.WriteString("\n")
	}
	if len(info.Fixed) > 0 {
		b.WriteString("Diperbaiki:\n")
		for _, s := range info.Fixed {
			b.WriteString("  • " + s + "\n")
		}
		b.WriteString("\n")
	}
	box := "[ ]"
	if skip {
		box = "[x]"
	}
	b.WriteString(box + " Jangan ingatkan untuk versi ini lagi\n\n")
	b.WriteString("Enter: update sekarang    Esc: nanti")
	return b.String()
}

func (m model) openUpdateModal() model {
	v := newModal("Pembaruan tersedia", updateModalBody(m.updateInfo, m.updateSkip))
	v.kind = modalUpdate
	v.footer = "Spasi centang · Enter update sekarang · Esc nanti"
	m.modal = v
	m.layout()
	return m
}

func (m model) keyUpdateModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		return m.dismissUpdate(), nil
	case tea.KeyEnter:
		return m.startUpdate()
	case tea.KeySpace:
		m.updateSkip = !m.updateSkip
		m = m.openUpdateModal()
		return m, nil
	default:
		switch msg.String() {
		case "q":
			return m.dismissUpdate(), nil
		case " ":
			m.updateSkip = !m.updateSkip
			m = m.openUpdateModal()
			return m, nil
		}
	}
	return m, nil
}

// dismissUpdate menutup modal; versi dicatat agar tidak muncul lagi bila dicentang.
func (m model) dismissUpdate() model {
	if m.updateSkip && m.updateInfo != nil {
		saveSkippedVersion(m.updateInfo.Version)
	}
	m.modal = nil
	m.updateInfo = nil
	m.layout()
	return m
}

// startUpdate menjalankan installer resmi lalu kembali ke TUI.
func (m model) startUpdate() (tea.Model, tea.Cmd) {
	info := m.updateInfo
	m.modal = nil
	m.updateInfo = nil
	if info == nil {
		m.layout()
		return m, nil
	}
	saveSkippedVersion(info.Version)
	m = m.note(roleSystem, fmt.Sprintf("Mengunduh dan memasang CipinCode CLI v%s…", info.Version))
	m.layout()
	m.refresh()

	base := strings.TrimRight(m.baseURL, "/")
	if base == "" {
		base = "https://cli.cipincode.com"
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("curl -fsSL %s/install.sh | bash", base))
	cmd.Env = append(os.Environ(), "CIPINCODE_VERSION="+info.Version)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return updateDoneMsg{version: info.Version, err: err}
	})
}

func (m model) handleVersion(msg versionMsg) (tea.Model, tea.Cmd) {
	manual := m.updateManual
	m.updateManual = false

	if msg.err != nil || msg.info == nil || msg.info.Version == "" {
		if manual {
			m = m.note(roleError, "Tidak bisa memeriksa pembaruan sekarang. Coba lagi nanti.")
			m.layout()
			m.refresh()
		}
		return m, nil
	}

	if !isNewer(version, msg.info.Version) {
		if manual {
			m = m.note(roleSystem, fmt.Sprintf("Sudah memakai versi terbaru (v%s).", version))
			m.layout()
			m.refresh()
		}
		return m, nil
	}

	if !manual && loadSkippedVersion() == msg.info.Version {
		return m, nil
	}
	if m.modal != nil {
		return m, nil
	}

	m.updateInfo = msg.info
	m.updateSkip = false
	return m.openUpdateModal(), nil
}

func (m model) handleUpdateDone(msg updateDoneMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m = m.note(roleError, "Gagal memperbarui otomatis: "+msg.err.Error()+". Jalankan manual: curl -fsSL https://cli.cipincode.com/install.sh | bash")
	} else {
		m = m.note(roleSystem, fmt.Sprintf("✓ CLI diperbarui ke v%s. Mulai ulang CLI untuk memakainya.", msg.version))
	}
	m.layout()
	m.refresh()
	return m, nil
}

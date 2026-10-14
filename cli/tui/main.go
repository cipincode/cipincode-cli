package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// version diisi saat build: go build -ldflags "-X main.version=1.2.3".
var version = "0.1.0"

func main() {
	if handled, err := runSubcommand(os.Args[1:]); handled {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	var showVersion bool
	var agentDir, serverURL string

	flag.BoolVar(&showVersion, "version", false, "tampilkan versi")
	flag.BoolVar(&showVersion, "v", false, "tampilkan versi (singkat)")
	flag.StringVar(&agentDir, "agent-dir", os.Getenv("CIPINCODE_AGENT_DIR"), "folder agent (berisi src/server.ts)")
	flag.StringVar(&serverURL, "server", os.Getenv("CIPINCODE_AGENT_URL"), "pakai server agent yang sudah berjalan")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "CipinCode CLI %s\n\nPemakaian:\n  cipincode [opsi]        buka sesi interaktif\n  cipincode login --key <API_KEY> [--url <alamat>]\n  cipincode logout\n  cipincode whoami\n\nOpsi:\n", version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Printf("cipincode %s\n", version)
		return
	}

	agent, err := connectAgent(serverURL, agentDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Gagal menyiapkan agent:", err)
		os.Exit(1)
	}
	defer agent.Stop()

	cfg, _ := loadConfig()
	lastSession := touchSession(time.Now())
	opts := options{screen: screenChat}
	if cfg == nil {
		opts.screen = screenOnboarding
	}
	program := tea.NewProgram(newModel(agent, cfg, lastSession, opts), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		agent.Stop()
		os.Exit(1)
	}
}

// runSubcommand menjalankan login, logout, atau whoami.
// handled bernilai false jika args bukan subperintah.
func runSubcommand(args []string) (handled bool, err error) {
	if len(args) == 0 {
		return false, nil
	}
	switch args[0] {
	case "login":
		return true, runLogin(args[1:])
	case "logout":
		return true, runLogout()
	case "whoami":
		return true, runWhoami()
	}
	return false, nil
}

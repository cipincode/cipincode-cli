package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

const defaultURL = "http://localhost:8000"

var errInvalidKey = errors.New("API key tidak valid, nonaktif, atau sudah dihapus")

type userInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// fetchUser memverifikasi API key ke server dan mengembalikan pemiliknya.
func fetchUser(baseURL, key string) (*userInfo, error) {
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/cli/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Accept", "application/json")

	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("tidak bisa terhubung ke %s: %w", baseURL, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		return nil, errInvalidKey
	default:
		return nil, fmt.Errorf("server membalas status %d", resp.StatusCode)
	}

	var body struct {
		User userInfo `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("balasan server tidak valid: %w", err)
	}
	return &body.User, nil
}

func runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	key := fs.String("key", "", "API key dari menu API Keys di dashboard")
	url := fs.String("url", envOr("CIPINCODE_URL", defaultURL), "alamat dashboard")
	if err := fs.Parse(args); err != nil {
		return err
	}
	urlSet := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "url" {
			urlSet = true
		}
	})

	if *key == "" {
		if term.IsTerminal(os.Stdin.Fd()) {
			if !urlSet {
				return runLoginInteractive("")
			}
			return runLoginInteractive(*url)
		}
		return errors.New("API key belum diisi.\nBuat di menu API Keys pada dashboard, lalu jalankan:\n  cipincode login --key <API_KEY>")
	}

	user, err := fetchUser(*url, *key)
	if err != nil {
		return err
	}
	if err := saveConfig(Config{URL: *url, Key: *key, Name: user.Name, Email: user.Email}); err != nil {
		return fmt.Errorf("menyimpan login: %w", err)
	}

	fmt.Printf("✓ Berhasil masuk sebagai %s (%s)\n", user.Name, user.Email)
	return nil
}

func runLogout() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if cfg == nil {
		fmt.Println("Anda belum login.")
		return nil
	}
	if err := deleteConfig(); err != nil {
		return err
	}
	fmt.Println("✓ Berhasil keluar.")
	return nil
}

func runWhoami() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if cfg == nil {
		return errors.New("Anda belum login. Jalankan: cipincode login --key <API_KEY>")
	}
	user, err := fetchUser(cfg.URL, cfg.Key)
	if err != nil {
		return err
	}
	fmt.Printf("%s (%s)\n%s\n", user.Name, user.Email, cfg.URL)
	return nil
}

func envOr(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

// runLoginInteractive membuka layar login (pilih email atau Google, URL, tempel kode).
// url kosong berarti memakai login sebelumnya, CIPINCODE_URL, atau alamat bawaan.
func runLoginInteractive(url string) error {
	final, err := tea.NewProgram(
		newModel(nil, nil, time.Time{}, options{screen: screenMethod, loginOnly: true, url: url}),
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		return err
	}

	if m, ok := final.(model); ok && m.loginOK && m.config != nil {
		fmt.Printf("✓ Berhasil masuk sebagai %s (%s)\n", m.config.Name, m.config.Email)
	}
	return nil
}

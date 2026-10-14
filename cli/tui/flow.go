package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Alur login "tempel kode" dengan PKCE:
//  1. CLI membuat verifier acak, hanya challenge-nya (hash) yang dikirim lewat URL.
//  2. Pengguna login di browser (email atau Google), klik Izinkan, lalu dapat kode.
//  3. CLI menukar kode + verifier dengan API key. Kode yang bocor tidak berguna
//     tanpa verifier yang hanya ada di proses CLI ini.

func newVerifier() (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func challengeFor(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func authorizeURL(baseURL, challenge, method string) string {
	q := url.Values{"challenge": {challenge}}
	if method == "google" {
		q.Set("method", "google")
	}
	return strings.TrimRight(baseURL, "/") + "/cli/authorize?" + q.Encode()
}

func deviceName() string {
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	return ""
}

// exchangeCode menukar kode dari browser dengan API key baru.
func exchangeCode(baseURL, code, verifier string) (string, *userInfo, error) {
	payload, _ := json.Marshal(map[string]string{
		"code":     code,
		"verifier": verifier,
		"device":   deviceName(),
	})

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/api/cli/exchange", bytes.NewReader(payload))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("tidak bisa terhubung ke %s", baseURL)
	}
	defer resp.Body.Close()

	var body struct {
		Key     string   `json:"key"`
		User    userInfo `json:"user"`
		Message string   `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	switch {
	case resp.StatusCode == http.StatusOK && body.Key != "":
		return body.Key, &body.User, nil
	case resp.StatusCode == http.StatusTooManyRequests:
		return "", nil, errors.New("terlalu banyak percobaan, tunggu sebentar lalu coba lagi")
	case body.Message != "":
		return "", nil, errors.New(body.Message)
	}
	return "", nil, fmt.Errorf("server membalas status %d", resp.StatusCode)
}

// openBrowser membuka URL di browser bawaan. Dinonaktifkan dengan CIPINCODE_NO_BROWSER=1
// (berguna lewat SSH atau saat pengujian).
func openBrowser(target string) error {
	if os.Getenv("CIPINCODE_NO_BROWSER") != "" {
		return errors.New("pembukaan browser dinonaktifkan")
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() { _ = cmd.Wait() }()
	return nil
}

// loginBaseURL memilih alamat dashboard: login sebelumnya, lalu env, lalu bawaan.
func loginBaseURL(cfg *Config) string {
	if cfg != nil && cfg.URL != "" {
		return cfg.URL
	}
	return envOr("CIPINCODE_URL", defaultURL)
}

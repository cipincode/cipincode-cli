package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State menyimpan data ringan antar sesi, di sebelah config.json.
type State struct {
	LastSession time.Time `json:"last_session"`
	Effort      string    `json:"effort,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Model       string    `json:"model,omitempty"`
	// SkippedVersion adalah versi yang tidak ingin diingatkan lagi.
	SkippedVersion string `json:"skipped_version,omitempty"`
}

func statePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "cipincode", "state.json"), nil
}

func loadState() (state State, ok bool) {
	path, err := statePath()
	if err != nil {
		return state, false
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return state, true
	}
	if err != nil || json.Unmarshal(data, &state) != nil {
		return State{}, false
	}
	return state, true
}

func saveState(state State) {
	path, err := statePath()
	if err != nil {
		return
	}
	data, err := json.Marshal(state)
	if err != nil || os.MkdirAll(filepath.Dir(path), 0o700) != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o600)
}

// touchSession mencatat sesi ini dan mengembalikan waktu sesi sebelumnya
// (zero value jika ini sesi pertama).
func touchSession(now time.Time) time.Time {
	state, ok := loadState()
	if !ok {
		return time.Time{}
	}
	previous := state.LastSession
	state.LastSession = now
	saveState(state)
	return previous
}

// loadEffort membaca pilihan effort terakhir (medium jika belum ada).
func loadEffort() effort {
	state, _ := loadState()
	return parseEffort(state.Effort)
}

func saveEffort(e effort) {
	state, ok := loadState()
	if !ok {
		return
	}
	state.Effort = e.key()
	saveState(state)
}

// loadProvider membaca ID provider terakhir ("" berarti Otomatis).
func loadProvider() string {
	state, _ := loadState()
	if parseProvider(state.Provider) == nil {
		return ""
	}
	return state.Provider
}

func saveProvider(id string) {
	state, ok := loadState()
	if !ok {
		return
	}
	state.Provider = id
	saveState(state)
}

// loadModel membaca ID model terakhir yang masih dikenal ("" berarti bawaan provider).
func loadModel() string {
	state, _ := loadState()
	if modelName(state.Model) == state.Model && state.Model != "" {
		return ""
	}
	return state.Model
}

func saveModel(id string) {
	state, ok := loadState()
	if !ok {
		return
	}
	state.Model = id
	saveState(state)
}

// loadSkippedVersion membaca versi yang dilewati lewat checkbox update.
func loadSkippedVersion() string {
	state, _ := loadState()
	return state.SkippedVersion
}

func saveSkippedVersion(v string) {
	state, ok := loadState()
	if !ok {
		return
	}
	state.SkippedVersion = v
	saveState(state)
}

var bulan = [...]string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}

// formatSession menulis waktu dalam bahasa Indonesia, mis. "13 Okt 2026, 22.05".
func formatSession(t time.Time) string {
	if t.IsZero() {
		return "belum ada (ini sesi pertama Anda)"
	}
	t = t.Local()
	return fmt.Sprintf("%d %s %d, %02d.%02d", t.Day(), bulan[t.Month()-1], t.Year(), t.Hour(), t.Minute())
}

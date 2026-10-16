package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Config adalah data login yang disimpan di ~/.config/cipincode/config.json.
type Config struct {
	URL  string `json:"url"`
	Key  string `json:"key"`
	Name string `json:"name"`
	// Email hanya untuk ditampilkan.
	Email string `json:"email"`
	// ProviderKeys adalah kunci API milik pengguna per provider
	// (BYOK) agar bisa terhubung langsung tanpa lewat server.
	ProviderKeys map[string]string `json:"provider_keys,omitempty"`
	// ClaudeAccount: Anthropic dipakai lewat akun Claude Pro/Max
	// (Claude Code lokal), bukan kunci API.
	ClaudeAccount bool `json:"claude_account,omitempty"`
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "cipincode", "config.json"), nil
}

// loadConfig mengembalikan nil jika belum login.
func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Key == "" {
		return nil, nil
	}
	return &cfg, nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func deleteConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// saveLogin menyimpan hasil login (dashboard) baru tanpa menghapus pilihan
// auth provider yang sudah ada (kunci API / akun Claude): auth yang dipilih
// tetap dipakai walau login ulang atau terminal ditutup-buka.
func saveLogin(baseURL, key, name, email string) (Config, error) {
	cfg := Config{URL: baseURL, Key: key, Name: name, Email: email}
	if old, _ := loadConfig(); old != nil {
		cfg.ProviderKeys = old.ProviderKeys
		cfg.ClaudeAccount = old.ClaudeAccount
	}
	if err := saveConfig(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

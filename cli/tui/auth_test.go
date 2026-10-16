package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cli/me" || r.Header.Get("Authorization") != "Bearer cpc_valid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":{"name":"Budi","email":"budi@example.com"}}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestFetchUser(t *testing.T) {
	srv := fakeServer(t)

	user, err := fetchUser(srv.URL, "cpc_valid")
	if err != nil || user.Name != "Budi" || user.Email != "budi@example.com" {
		t.Fatalf("got %+v, %v", user, err)
	}

	if _, err := fetchUser(srv.URL, "cpc_salah"); err != errInvalidKey {
		t.Fatalf("want errInvalidKey, got %v", err)
	}
}

func TestLoginLogoutRoundTrip(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	srv := fakeServer(t)

	if err := runLogin([]string{"--key", "cpc_valid", "--url", srv.URL}); err != nil {
		t.Fatalf("login: %v", err)
	}
	cfg, err := loadConfig()
	if err != nil || cfg == nil || cfg.Name != "Budi" || cfg.Key != "cpc_valid" || cfg.URL != srv.URL {
		t.Fatalf("config: %+v, %v", cfg, err)
	}

	if err := runLogout(); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if cfg, _ := loadConfig(); cfg != nil {
		t.Fatalf("config masih ada setelah logout: %+v", cfg)
	}
}

func TestLoginRejectsBadKeyAndMissingKey(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	srv := fakeServer(t)

	if err := runLogin([]string{"--key", "cpc_salah", "--url", srv.URL}); err != errInvalidKey {
		t.Fatalf("want errInvalidKey, got %v", err)
	}
	if cfg, _ := loadConfig(); cfg != nil {
		t.Fatal("config tersimpan padahal key salah")
	}
	if err := runLogin(nil); err == nil {
		t.Fatal("login tanpa --key seharusnya gagal")
	}
}

func TestSaveLoginKeepsProviderAuth(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	old := Config{URL: "https://lama", Key: "k", Name: "A", Email: "a@x",
		ProviderKeys: map[string]string{"openai": "sk-lama"}, ClaudeAccount: true}
	if err := saveConfig(old); err != nil {
		t.Fatal(err)
	}
	cfg, err := saveLogin("https://baru", "k2", "B", "b@x")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.URL != "https://baru" || cfg.Key != "k2" {
		t.Fatalf("login baru tidak tersimpan: %+v", cfg)
	}
	if cfg.ProviderKeys["openai"] != "sk-lama" || !cfg.ClaudeAccount {
		t.Fatalf("auth provider ikut terhapus: %+v", cfg)
	}
}

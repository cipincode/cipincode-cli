package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestChallengeMatchesKnownSha256Vector(t *testing.T) {
	// base64url(sha256("abc")), sama dengan hasil CliAuthCode::challengeFor di Laravel.
	want := "ungWv48Bz-pBQUDeXa4iI7ADYaOWF3qctBD_YfIAFa0"
	if got := challengeFor("abc"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNewVerifierIsLongRandomAndURLSafe(t *testing.T) {
	a, err := newVerifier()
	if err != nil {
		t.Fatal(err)
	}
	b, _ := newVerifier()

	if len(a) < 43 || len(a) > 128 {
		t.Fatalf("panjang verifier %d di luar 43..128", len(a))
	}
	if a == b {
		t.Fatal("verifier seharusnya acak")
	}
	if strings.ContainsAny(a, "+/=") {
		t.Fatalf("verifier tidak URL-safe: %q", a)
	}
}

func TestAuthorizeURL(t *testing.T) {
	challenge := challengeFor("abc")

	email, _ := url.Parse(authorizeURL("http://localhost:8000/", challenge, "email"))
	if email.Path != "/cli/authorize" || email.Query().Get("challenge") != challenge || email.Query().Has("method") {
		t.Fatalf("url email salah: %s", email)
	}

	google, _ := url.Parse(authorizeURL("http://localhost:8000", challenge, "google"))
	if google.Query().Get("method") != "google" {
		t.Fatalf("url google tanpa method: %s", google)
	}
}

func exchangeServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got map[string]string
		_ = json.NewDecoder(r.Body).Decode(&got)
		if r.URL.Path != "/api/cli/exchange" || r.Method != http.MethodPost || got["code"] == "" || got["verifier"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestExchangeCodeSuccess(t *testing.T) {
	srv := exchangeServer(t, 200, map[string]any{
		"key":  "cpc_abc",
		"user": map[string]string{"name": "Budi", "email": "budi@example.com"},
	})

	key, user, err := exchangeCode(srv.URL, "ABCD-EFGH", "verifier")
	if err != nil || key != "cpc_abc" || user.Name != "Budi" {
		t.Fatalf("got %q, %+v, %v", key, user, err)
	}
}

func TestExchangeCodeSurfacesServerMessages(t *testing.T) {
	invalid := exchangeServer(t, 422, map[string]string{"message": "Kode tidak valid atau sudah kedaluwarsa."})
	if _, _, err := exchangeCode(invalid.URL, "X", "v"); err == nil || !strings.Contains(err.Error(), "tidak valid") {
		t.Fatalf("got %v", err)
	}

	limited := exchangeServer(t, 429, map[string]string{"message": "Too Many Attempts."})
	if _, _, err := exchangeCode(limited.URL, "X", "v"); err == nil || !strings.Contains(err.Error(), "terlalu banyak") {
		t.Fatalf("got %v", err)
	}

	if _, _, err := exchangeCode("http://127.0.0.1:1", "X", "v"); err == nil || !strings.Contains(err.Error(), "tidak bisa terhubung") {
		t.Fatalf("got %v", err)
	}
}

func TestMatchCommands(t *testing.T) {
	names := func(cs []command) string {
		var n []string
		for _, c := range cs {
			n = append(n, c.name)
		}
		return strings.Join(n, ",")
	}

	if got := matchCommands("/"); len(got) != len(commands) {
		t.Fatalf("'/' seharusnya menampilkan semua perintah, got %d", len(got))
	}
	if got := names(matchCommands("/co")); got != "coding,connect" {
		t.Fatalf("'/co' => %s", got)
	}
	if got := names(matchCommands("/lo")); got != "login,logout" {
		t.Fatalf("'/lo' => %s", got)
	}
	// yang berawalan didahulukan, baru yang sekadar memuat
	if got := names(matchCommands("/e")); !strings.HasPrefix(got, "effort,exit,") {
		t.Fatalf("'/e' => %s", got)
	}
	if got := names(matchCommands("/eff")); got != "effort" {
		t.Fatalf("'/eff' => %s", got)
	}
	if matchCommands("/xyz") != nil || matchCommands("halo") != nil || matchCommands("/coding ") != nil {
		t.Fatal("input yang bukan awalan perintah tidak boleh membuka menu")
	}
	if len(commands) < 10 {
		t.Fatal("menu seharusnya punya banyak pilihan")
	}
}

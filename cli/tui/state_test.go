package main

import (
	"strings"
	"testing"
	"time"
)

func TestTouchSessionReturnsPreviousSession(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	first := time.Date(2026, 10, 13, 22, 5, 0, 0, time.Local)
	second := first.Add(48 * time.Hour)

	if prev := touchSession(first); !prev.IsZero() {
		t.Fatalf("sesi pertama seharusnya tanpa sesi sebelumnya, got %v", prev)
	}
	if prev := touchSession(second); !prev.Equal(first) {
		t.Fatalf("got %v, want %v", prev, first)
	}
	if prev := touchSession(second.Add(time.Hour)); !prev.Equal(second) {
		t.Fatalf("got %v, want %v", prev, second)
	}
}

func TestFormatSession(t *testing.T) {
	got := formatSession(time.Date(2026, 10, 13, 22, 5, 0, 0, time.Local))
	if got != "13 Okt 2026, 22.05" {
		t.Fatalf("got %q", got)
	}
	if got := formatSession(time.Time{}); !strings.Contains(got, "sesi pertama") {
		t.Fatalf("got %q", got)
	}
}

func TestModeCyclesThroughThreeModes(t *testing.T) {
	m := modeChat
	var seen []string
	for i := 0; i < 4; i++ {
		seen = append(seen, m.key())
		m = m.next()
	}
	want := []string{"chat", "coding", "debugging", "chat"}
	for i := range want {
		if seen[i] != want[i] {
			t.Fatalf("urutan mode %v, want %v", seen, want)
		}
	}
}

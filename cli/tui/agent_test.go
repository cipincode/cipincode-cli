package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func collectEvents(t *testing.T, stream string) ([]chatEvent, error) {
	t.Helper()
	out := make(chan chatEvent, 16)
	err := readEvents(strings.NewReader(stream), out)
	close(out)

	var events []chatEvent
	for ev := range out {
		events = append(events, ev)
	}
	return events, err
}

func TestReadEventsStreamsTokensUntilDone(t *testing.T) {
	events, err := collectEvents(t,
		"data: {\"type\":\"token\",\"text\":\"Hai\"}\n\n"+
			"data: {\"type\":\"token\",\"text\":\" dunia\"}\n\n"+
			"data: {\"type\":\"done\"}\n\n"+
			"data: {\"type\":\"token\",\"text\":\"diabaikan\"}\n\n")
	if err != nil {
		t.Fatalf("readEvents: %v", err)
	}

	var got string
	for _, ev := range events {
		got += ev.text
	}
	if got != "Hai dunia" {
		t.Fatalf("got %q, want %q", got, "Hai dunia")
	}
}

func TestReadEventsReturnsAgentError(t *testing.T) {
	_, err := collectEvents(t, "data: {\"type\":\"error\",\"message\":\"boom\"}\n\n")
	if err == nil || err.Error() != "boom" {
		t.Fatalf("got %v, want boom", err)
	}
}

func TestReadEventsParsesUsageSnapshot(t *testing.T) {
	events, err := collectEvents(t,
		"data: {\"type\":\"usage\",\"session\":{\"used\":1,\"limit\":40,\"percent\":3},\"weekly\":{\"used\":1,\"limit\":400,\"percent\":0}}\n\n"+
			"data: {\"type\":\"done\"}\n\n")
	if err != nil || len(events) != 1 || events[0].kind != "usage" {
		t.Fatalf("events=%+v err=%v", events, err)
	}
	if snap := events[0].snap; snap.Session.Used != 1 || snap.Weekly.Limit != 400 {
		t.Fatalf("snapshot salah: %+v", snap)
	}
}

func TestReadEventsStopsAtBlockedWithSnapshotAndMessage(t *testing.T) {
	events, err := collectEvents(t,
		"data: {\"type\":\"blocked\",\"message\":\"Batas habis.\",\"session\":{\"used\":40,\"limit\":40,\"exhausted\":true},\"weekly\":{\"used\":40,\"limit\":400}}\n\n"+
			"data: {\"type\":\"token\",\"text\":\"tidak boleh\"}\n\n")
	if err != nil || len(events) != 1 {
		t.Fatalf("events=%+v err=%v", events, err)
	}
	if ev := events[0]; ev.kind != "blocked" || ev.message != "Batas habis." || !ev.snap.Session.Exhausted {
		t.Fatalf("event salah: %+v", ev)
	}
}

func TestReadEventsIgnoresMalformedLines(t *testing.T) {
	events, err := collectEvents(t, ": ping\n\ndata: bukan-json\n\ndata: {\"type\":\"token\",\"text\":\"ok\"}\n\n")
	if err != nil || len(events) != 1 || events[0].text != "ok" {
		t.Fatalf("events=%+v err=%v", events, err)
	}
}

func TestLastLinesMengambilEkor(t *testing.T) {
	if got := lastLines("", 20); got != "" {
		t.Fatalf("kosong harus tetap kosong, got %q", got)
	}
	in := "a\nb\nc\nd\ne"
	if got := lastLines(in, 2); got != "d\ne" {
		t.Fatalf("got %q, want %q", got, "d\ne")
	}
	if got := lastLines(in, 10); got != in {
		t.Fatalf("got %q, want %q", got, in)
	}
}

func TestSyncBufferMenampungTulisan(t *testing.T) {
	var b syncBuffer
	if _, err := b.Write([]byte("halo")); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Write([]byte(" dunia")); err != nil {
		t.Fatal(err)
	}
	if got := b.String(); got != "halo dunia" {
		t.Fatalf("got %q", got)
	}
}

func TestFindRuntimesGagalBilaFolderKosong(t *testing.T) {
	dir := t.TempDir()
	if _, err := findRuntimes(dir); err == nil {
		t.Fatal("harus error untuk folder agent kosong")
	}
}

func TestFindRuntimesMenemukanDistBilaAdaNode(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dist", "server.node.js"), []byte("// fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	rts, err := findRuntimes(dir)
	if err != nil {
		t.Skipf("node/bun tidak tersedia di CI: %v", err)
	}
	foundNode := false
	for _, r := range rts {
		if r.mode == "node" {
			foundNode = true
		}
		if r.healthTimeout() <= 0 || r.healthTimeout() > 30*time.Second {
			t.Fatalf("timeout tidak wajar: %v", r.healthTimeout())
		}
		if r.describe() == "" {
			t.Fatal("describe kosong")
		}
	}
	// Bila node terpasang, mode node harus diutamakan (pertama).
	if len(rts) > 0 && rts[0].mode != "node" && foundNode {
		t.Fatalf("node harus diutamakan, got %+v", rts)
	}
}

func TestDumpAgentLogMenulisFile(t *testing.T) {
	p := dumpAgentLog("baris1\nbaris2")
	if p == "" {
		t.Fatal("path kosong")
	}
	defer os.Remove(p)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "baris1\nbaris2" {
		t.Fatalf("got %q", string(data))
	}
	if lastAgentLogPath() != p {
		t.Fatalf("lastAgentLogPath=%q want %q", lastAgentLogPath(), p)
	}
}

package main

import (
	"strings"
	"testing"
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

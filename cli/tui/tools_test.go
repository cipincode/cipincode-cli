package main

import (
	"strings"
	"testing"
)

func TestRunShellEcho(t *testing.T) {
	out, code := runShell("echo halo-dunia")
	if code != 0 || out != "halo-dunia" {
		t.Fatalf("got %q, %d", out, code)
	}
}

func TestRunShellExitCode(t *testing.T) {
	_, code := runShell("exit 3")
	if code != 3 {
		t.Fatalf("exit code = %d, seharusnya 3", code)
	}
}

func TestRunShellOutputCapped(t *testing.T) {
	out, _ := runShell("printf '%9000s' 'x'")
	if len(out) > toolOutputCap+100 {
		t.Fatalf("output tidak dipotong: %d byte", len(out))
	}
	if !strings.Contains(out, "dipotong") {
		t.Error("seharusnya ada penanda pemotongan")
	}
}

func TestShellCallsFiltersExecuteCommand(t *testing.T) {
	src := "<function_calls>\n" +
		"<invoke name=\"execute_command\">\n<parameter name=\"command\">ls -la</parameter>\n</invoke>\n" +
		"<invoke name=\"baca_file\">\n<parameter name=\"path\">x</parameter>\n</invoke>\n" +
		"<invoke name=\"execute_command\">\n<parameter name=\"command\"></parameter>\n</invoke>\n" +
		"</function_calls>"
	calls := shellCalls(src)
	if len(calls) != 1 || calls[0].command != "ls -la" {
		t.Fatalf("got %+v", calls)
	}
}

func TestToolLoopRunsInCodingMode(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := newScreen(t, screenChat, loggedIn(""), "")
	m.mode = modeCoding
	m.messages = append(m.messages,
		message{role: roleUser, content: "lihat isi folder"},
		message{role: roleAssistant, content: "<invoke name=\"execute_command\">\n<parameter name=\"command\">echo ok</parameter>\n</invoke>"},
	)

	m2, cmd := update(m, doneMsg{})
	if cmd == nil {
		t.Fatal("mode Coding seharusnya menjalankan tool otomatis")
	}
	// Jalankan cmd tool secara langsung (tanpa stream lanjutan): hasilnya
	// harus memuat output perintah untuk dikirim balik ke AI.
	got, ok := cmd().(toolResultMsg)
	if !ok {
		t.Fatalf("cmd seharusnya toolResultMsg, got %T", cmd())
	}
	_ = m2
	if !strings.Contains(got.results, "Hasil perintah") || !strings.Contains(got.results, "echo ok") {
		t.Errorf("hasil eksekusi tidak lengkap:\n%s", got.results)
	}
}

func TestToolLoopSkippedInChatAndAtDepthCap(t *testing.T) {
	mk := func(md mode, depth int) model {
		m := newScreen(t, screenChat, loggedIn(""), "")
		m.mode = md
		m.toolDepth = depth
		m.messages = append(m.messages,
			message{role: roleAssistant, content: "<invoke name=\"execute_command\">\n<parameter name=\"command\">echo ok</parameter>\n</invoke>"},
		)
		return m
	}

	m, cmd := update(mk(modeChat, 0), doneMsg{})
	if cmd != nil {
		t.Error("mode Chat tidak boleh auto-jalan (hanya pratinjau)")
	}
	_ = m

	m, cmd = update(mk(modeCoding, maxToolDepth), doneMsg{})
	if cmd != nil {
		t.Error("lewat batas depth seharusnya berhenti")
	}
	if view := plain(m); !strings.Contains(view, "Berhenti otomatis") {
		t.Errorf("seharusnya ada catatan berhenti:\n%s", view)
	}
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Eksekusi tool AI di komputer pengguna (mode Coding/Debugging, auto-approve).
// Satu-satunya tool yang didukung: execute_command (shell lewat bash).

const (
	maxToolCalls  = 5
	toolTimeout   = 30 * time.Second
	toolOutputCap = 8000
	maxToolDepth  = 4 // langkah tool beruntun per satu pesan pengguna
)

type shellCall struct{ command string }

// shellCalls menyaring panggilan execute_command dari jawaban AI.
func shellCalls(src string) []shellCall {
	var out []shellCall
	for _, c := range parseToolCalls(src) {
		if !strings.EqualFold(c.name, "execute_command") {
			continue
		}
		for _, p := range c.params {
			if strings.EqualFold(p[0], "command") && strings.TrimSpace(p[1]) != "" {
				out = append(out, shellCall{command: strings.TrimSpace(p[1])})
			}
		}
		if len(out) >= maxToolCalls {
			break
		}
	}
	return out
}

type toolResultMsg struct{ results string }

// runToolsCmd menjalankan perintah shell lalu mengembalikan hasilnya
// sebagai pesan lanjutan untuk AI.
func runToolsCmd(calls []shellCall) tea.Cmd {
	return func() tea.Msg {
		var b strings.Builder
		b.WriteString("Hasil perintah shell yang kamu minta (dijalankan otomatis):\n")
		for _, c := range calls {
			out, code := runShell(c.command)
			fmt.Fprintf(&b, "\n$ %s\n%s\n(keluar: %d)\n", c.command, out, code)
		}
		return toolResultMsg{results: strings.TrimRight(b.String(), "\n")}
	}
}

// runShell menjalankan satu perintah lewat bash dengan batas waktu dan
// potongan output agar aman dan tidak membanjiri percakapan.
func runShell(command string) (string, int) {
	ctx, cancel := context.WithTimeout(context.Background(), toolTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()

	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			code = 1
		}
	}
	out := strings.TrimRight(buf.String(), "\n")
	if ctx.Err() == context.DeadlineExceeded {
		out += "\n(waktu habis 30 dtk, proses dihentikan)"
		code = 124
	}
	if len(out) > toolOutputCap {
		out = out[:toolOutputCap] + "\n…(dipotong, 8000 karakter pertama)"
	}
	if out == "" {
		out = "(tanpa output)"
	}
	return out, code
}

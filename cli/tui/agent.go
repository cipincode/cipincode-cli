package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Agent adalah koneksi ke server agent (Bun/TypeScript).
type Agent struct {
	URL string
	// Token dibuat acak saat CLI menjalankan agent, agar hanya CLI ini yang bisa memakainya.
	Token string
	cmd   *exec.Cmd
}

type historyItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest adalah satu giliran percakapan. URL dan Key (login pengguna) dikirim
// lewat header agar agent bisa memanggil dashboard atas nama pengguna.
type chatRequest struct {
	Message string        `json:"message"`
	Mode    string        `json:"mode"`
	Effort  string        `json:"effort"`
	Session string        `json:"session,omitempty"`
	Agent   int           `json:"agent,omitempty"`
	History []historyItem `json:"history"`
	URL     string        `json:"-"`
	Key     string        `json:"-"`
}

// chatEvent adalah peristiwa dari aliran jawaban: token, usage, atau blocked.
type chatEvent struct {
	kind    string
	text    string
	snap    *usageSnapshot
	message string
}

type event struct {
	Type    string      `json:"type"`
	Text    string      `json:"text"`
	Message string      `json:"message"`
	Session *usageMeter `json:"session"`
	Weekly  *usageMeter `json:"weekly"` // null untuk paket tanpa batas mingguan
	Plan    *usagePlan  `json:"plan"`
}

func (e event) snapshot() *usageSnapshot {
	if e.Session == nil {
		return nil
	}
	snap := &usageSnapshot{Session: *e.Session, Weekly: e.Weekly}
	if e.Plan != nil {
		snap.Plan = *e.Plan
	}
	return snap
}

// connectAgent memakai server yang sudah ada, atau menjalankan agent lewat Node
// (diutamakan, hasil build npm tanpa Bun) lalu fallback ke Bun untuk dev.
func connectAgent(serverURL, agentDir string) (*Agent, error) {
	if serverURL != "" {
		a := &Agent{URL: strings.TrimRight(serverURL, "/"), Token: os.Getenv("CIPINCODE_AGENT_TOKEN")}
		if err := a.waitHealthy(3 * time.Second); err != nil {
			return nil, err
		}
		return a, nil
	}

	dir, err := findAgentDir(agentDir)
	if err != nil {
		return nil, err
	}
	rt, err := findRuntime(dir)
	if err != nil {
		return nil, err
	}

	port, err := freePort()
	if err != nil {
		return nil, err
	}

	token, err := newVerifier()
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	if rt.mode == "node" {
		cmd = exec.Command(rt.bin, rt.entry)
	} else {
		cmd = exec.Command(rt.bin, "run", filepath.Join("src", "server.ts"))
	}
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		"CIPINCODE_VERSION="+version,
		"CIPINCODE_AGENT_TOKEN="+token,
	)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("menjalankan agent: %w", err)
	}

	a := &Agent{URL: fmt.Sprintf("http://127.0.0.1:%d", port), Token: token, cmd: cmd}
	if err := a.waitHealthy(8 * time.Second); err != nil {
		a.Stop()
		return nil, err
	}
	return a, nil
}

// runtimeInfo adalah interpreter untuk menjalankan agent.
type runtimeInfo struct {
	bin   string // path node atau bun
	mode  string // "node" atau "bun"
	entry string // file agent relatif dari agent dir
}

// findRuntime mengutamakan Node + dist/server.node.js (package npm, tanpa Bun),
// lalu fallback ke Bun + src/server.ts (dev / install.sh lama).
func findRuntime(agentDir string) (runtimeInfo, error) {
	if node, err := exec.LookPath("node"); err == nil {
		for _, e := range []string{
			filepath.Join("dist", "server.node.js"),
			filepath.Join("src", "server.node.js"),
		} {
			if _, err := os.Stat(filepath.Join(agentDir, e)); err == nil {
				return runtimeInfo{bin: node, mode: "node", entry: e}, nil
			}
		}
	}
	if bun, err := findBun(); err == nil {
		if _, err := os.Stat(filepath.Join(agentDir, "src", "server.ts")); err == nil {
			return runtimeInfo{bin: bun, mode: "bun", entry: filepath.Join("src", "server.ts")}, nil
		}
	}
	return runtimeInfo{}, errors.New("Node.js tidak ditemukan. Pasang Node.js 18+ dari https://nodejs.org (atau Bun dari https://bun.sh untuk dev)")
}

func findBun() (string, error) {
	if p, err := exec.LookPath("bun"); err == nil {
		return p, nil
	}
	if home, err := os.UserHomeDir(); err == nil {
		for _, p := range []string{
			filepath.Join(home, ".bun", "bin", "bun"),
			filepath.Join(home, ".local", "bin", "bun"),
		} {
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}
	}
	return "", errors.New("Bun tidak ditemukan. Pasang dari https://bun.sh")
}

func findAgentDir(explicit string) (string, error) {
	var candidates []string
	if explicit != "" {
		candidates = append(candidates, explicit)
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates, filepath.Join(base, "..", "agent"), filepath.Join(base, "agent"))
	}
	candidates = append(candidates, "agent", filepath.Join("cli", "agent"))

	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, "dist", "server.node.js")); err == nil {
			abs, err := filepath.Abs(c)
			if err != nil {
				return "", err
			}
			return abs, nil
		}
		if _, err := os.Stat(filepath.Join(c, "src", "server.ts")); err == nil {
			abs, err := filepath.Abs(c)
			if err != nil {
				return "", err
			}
			return abs, nil
		}
	}
	return "", errors.New("folder agent tidak ditemukan; gunakan --agent-dir atau CIPINCODE_AGENT_DIR")
}

func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func (a *Agent) waitHealthy(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(a.URL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("agent tidak merespons di %s", a.URL)
}

// Stop menghentikan proses agent yang dijalankan oleh CLI ini.
func (a *Agent) Stop() {
	if a == nil || a.cmd == nil || a.cmd.Process == nil {
		return
	}
	_ = a.cmd.Process.Kill()
	_ = a.cmd.Wait()
	a.cmd = nil
}

// Chat mengirim satu giliran percakapan dan mengalirkan peristiwanya ke out.
func (a *Agent) Chat(ctx context.Context, chat chatRequest, out chan<- chatEvent) error {
	body, _ := json.Marshal(chat)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.URL+"/chat", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if a.Token != "" {
		req.Header.Set("X-Agent-Token", a.Token)
	}
	if chat.URL != "" && chat.Key != "" {
		req.Header.Set("X-Cipincode-Url", chat.URL)
		req.Header.Set("X-Cipincode-Key", chat.Key)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent membalas status %d", resp.StatusCode)
	}
	return readEvents(resp.Body, out)
}

// readEvents membaca aliran SSE dari agent. Aliran berakhir pada peristiwa done,
// blocked (jatah habis), atau error.
func readEvents(r io.Reader, out chan<- chatEvent) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var ev event
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &ev); err != nil {
			continue
		}
		switch ev.Type {
		case "token":
			out <- chatEvent{kind: "token", text: ev.Text}
		case "usage":
			if snap := ev.snapshot(); snap != nil {
				out <- chatEvent{kind: "usage", snap: snap}
			}
		case "blocked":
			out <- chatEvent{kind: "blocked", snap: ev.snapshot(), message: ev.Message}
			return nil
		case "error":
			return errors.New(ev.Message)
		case "done":
			return nil
		}
	}
	return scanner.Err()
}

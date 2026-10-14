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
	"sync"
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
	Message  string        `json:"message"`
	Mode     string        `json:"mode"`
	Effort   string        `json:"effort"`
	Provider string        `json:"provider,omitempty"`
	Model    string        `json:"model,omitempty"`
	// ProviderKey hanya untuk agent lokal (tidak diteruskan ke dashboard).
	ProviderKey string        `json:"provider_key,omitempty"`
	Session  string        `json:"session,omitempty"`
	Agent   int           `json:"agent,omitempty"`
	History []historyItem `json:"history"`
	Cwd     string        `json:"cwd,omitempty"`
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

// connectAgent memakai server yang sudah ada, atau menjalankan agent lokal.
// Mencoba semua runtime yang tersedia berurutan (Node dulu, lalu Bun) agar
// satu runtime yang rusak tidak langsung menggagalkan CLI. Log agent ditangkap
// dan ditampilkan saat gagal agar mudah didiagnosis di laptop pengguna.
func connectAgent(serverURL, agentDir string) (*Agent, error) {
	// Gunakan domain remote default jika tidak diset secara lokal agar CLI tidak perlu menjalankan server node lokal.
	if serverURL == "" {
		serverURL = os.Getenv("CIPINCODE_AGENT_URL")
	}
	if serverURL == "" {
		serverURL = "https://cli.cipincode.com"
	}
	a := &Agent{URL: strings.TrimRight(serverURL, "/"), Token: os.Getenv("CIPINCODE_AGENT_TOKEN")}
	if err := a.waitHealthy(5 * time.Second); err == nil {
		return a, nil
	}

	if serverURL != "https://cli.cipincode.com" {
		return nil, fmt.Errorf("agent di %s tidak merespons", a.URL)
	}

	// Fallback ke proses lokal jika domain utama tidak merespons
	dir, err := findAgentDir(agentDir)
	if err != nil {
		return nil, err
	}
	runtimes, err := findRuntimes(dir)
	if err != nil {
		return nil, err
	}

	token, err := newVerifier()
	if err != nil {
		return nil, err
	}

	var failures []string
	for _, rt := range runtimes {
		port, err := freePort()
		if err != nil {
			return nil, err
		}
		a, logs, startErr := startAgent(dir, rt, port, token)
		if startErr != nil {
			failures = append(failures, fmt.Sprintf("%s: gagal dijalankan: %v", rt.describe(), startErr))
			continue
		}
		if err := a.waitHealthy(rt.healthTimeout()); err != nil {
			full := logs.String()
			tail := lastLines(full, 20)
			a.Stop()
			logPath := dumpAgentLog(full)
			detail := fmt.Sprintf("%s: %v", rt.describe(), err)
			if tail != "" {
				detail += "\nLog agent:\n" + tail
			}
			if logPath != "" {
				detail += "\nLog lengkap: " + logPath
			}
			// Jika proses mati cepat, kemungkinan crash (modul hilang, syntax, izin).
			// Log di atas biasanya sudah menjelaskan sebabnya.
			failures = append(failures, detail)
			continue
		}
		return a, nil
	}

	var sb strings.Builder
	for _, f := range failures {
		sb.WriteString("  • " + f + "\n")
	}
	sb.WriteString("Periksa: node -v (butuh 18+), bun --version, folder agent di " + dir +
		" (harus ada dist/server.node.js atau src/server.ts)." +
		" Coba: cipincode --agent-dir <folder-agent>" +
		" atau CIPINCODE_AGENT_DEBUG=1 cipincode untuk log penuh.")
	return nil, errors.New(sb.String())
}

// startAgent menjalankan satu runtime agent di port tertentu dan menangkap
// stdout/stderr-nya ke buffer. Bila CIPINCODE_AGENT_DEBUG=1, log juga
// diteruskan ke stderr agar terlihat langsung di terminal.
func startAgent(dir string, rt runtimeInfo, port int, token string) (*Agent, *syncBuffer, error) {
	var cmd *exec.Cmd
	if rt.mode == "node" {
		cmd = exec.Command(rt.bin, filepath.Join(dir, rt.entry))
	} else {
		cmd = exec.Command(rt.bin, "run", filepath.Join(dir, "src", "server.ts"))
	}
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		"CIPINCODE_VERSION="+version,
		"CIPINCODE_AGENT_TOKEN="+token,
	)
	logs := &syncBuffer{}
	if os.Getenv("CIPINCODE_AGENT_DEBUG") == "1" {
		cmd.Stdout = io.MultiWriter(logs, os.Stderr)
		cmd.Stderr = io.MultiWriter(logs, os.Stderr)
	} else {
		cmd.Stdout = logs
		cmd.Stderr = logs
	}
	if err := cmd.Start(); err != nil {
		return nil, logs, err
	}
	a := &Agent{URL: fmt.Sprintf("http://127.0.0.1:%d", port), Token: token, cmd: cmd}
	return a, logs, nil
}

// syncBuffer adalah bytes.Buffer yang aman untuk ditulis dari goroutine
// pemantau proses agent.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

var (
	agentLogPathMu sync.Mutex
	agentLogPath   string
)

func rememberAgentLogPath(p string) {
	agentLogPathMu.Lock()
	defer agentLogPathMu.Unlock()
	agentLogPath = p
}

func lastAgentLogPath() string {
	agentLogPathMu.Lock()
	defer agentLogPathMu.Unlock()
	return agentLogPath
}

// dumpAgentLog menyimpan log penuh ke file temp untuk diagnosis lanjut.
func dumpAgentLog(logs string) string {
	f, err := os.CreateTemp("", "cipincode-agent-*.log")
	if err != nil {
		return ""
	}
	_, _ = f.WriteString(logs)
	_ = f.Close()
	rememberAgentLogPath(f.Name())
	return f.Name()
}

// lastLines mengambil N baris terakhir agar pesan error tetap ringkas.
func lastLines(s string, n int) string {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	// Batasi total panjang agar tidak membanjiri terminal.
	out := strings.Join(lines, "\n")
	if len(out) > 2000 {
		out = out[len(out)-2000:]
		if i := strings.Index(out, "\n"); i >= 0 {
			out = out[i+1:]
		}
	}
	return strings.TrimSpace(out)
}

// runtimeInfo adalah interpreter untuk menjalankan agent.
type runtimeInfo struct {
	bin   string // path node atau bun
	mode  string // "node" atau "bun"
	entry string // file agent relatif dari agent dir
}

func (r runtimeInfo) describe() string {
	if r.mode == "node" {
		return fmt.Sprintf("node (%s %s)", r.bin, r.entry)
	}
	return fmt.Sprintf("bun (%s)", r.bin)
}

// healthTimeout memberi Bun waktu lebih lama karena kompilasi awal di macOS
// bisa lambat; Node biasanya siap <1 detik.
func (r runtimeInfo) healthTimeout() time.Duration {
	if r.mode == "bun" {
		return 15 * time.Second
	}
	return 10 * time.Second
}

// findRuntimes mengembalikan semua runtime yang bisa dipakai, berurutan:
// Node + dist/server.node.js dulu (tanpa Bun), lalu Bun + src/server.ts.
func findRuntimes(agentDir string) ([]runtimeInfo, error) {
	var out []runtimeInfo
	if node, err := exec.LookPath("node"); err == nil {
		if _, err := os.Stat(filepath.Join(agentDir, "dist", "server.node.js")); err == nil {
			out = append(out, runtimeInfo{bin: node, mode: "node", entry: filepath.Join("dist", "server.node.js")})
		}
	}
	if bun, err := findBun(); err == nil {
		if _, err := os.Stat(filepath.Join(agentDir, "src", "server.ts")); err == nil {
			out = append(out, runtimeInfo{bin: bun, mode: "bun", entry: filepath.Join("src", "server.ts")})
		}
	}
	if len(out) > 0 {
		return out, nil
	}
	// Pesan spesifik agar pengguna tahu apa yang hilang.
	hasDist := false
	hasSrc := false
	if _, err := os.Stat(filepath.Join(agentDir, "dist", "server.node.js")); err == nil {
		hasDist = true
	}
	if _, err := os.Stat(filepath.Join(agentDir, "src", "server.ts")); err == nil {
		hasSrc = true
	}
	switch {
	case !hasDist && !hasSrc:
		return nil, errors.New("folder agent tidak lengkap di " + agentDir + " (tidak ada dist/server.node.js maupun src/server.ts); install ulang: curl -fsSL https://cli.cipincode.com/install.sh | bash")
	case hasDist && !hasSrc:
		return nil, errors.New("Node.js tidak ditemukan. Pasang Node.js 18+ dari https://nodejs.org (agent hanya ada versi dist, butuh Node)")
	default:
		return nil, errors.New("Node.js maupun Bun tidak ditemukan. Pasang Node.js 18+ dari https://nodejs.org (atau Bun dari https://bun.sh)")
	}
}

// findRuntime dipertahankan untuk kompatibilitas test lama.
func findRuntime(agentDir string) (runtimeInfo, error) {
	rts, err := findRuntimes(agentDir)
	if err != nil {
		return runtimeInfo{}, err
	}
	return rts[0], nil
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
	if env := os.Getenv("CIPINCODE_AGENT_DIR"); env != "" {
		candidates = append(candidates, env)
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".local", "share", "cipincode", "agent"))
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

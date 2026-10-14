#!/usr/bin/env node
/**
 * install.js – postinstall script untuk npm package cipincode
 *
 * Yang dilakukan:
 *   1. Detect OS + arch
 *   2. Install Go jika belum ada (atau versinya terlalu lama)
 *   3. Install Bun jika belum ada
 *   4. Download binary cipincode dari GitHub Releases
 *   5. Simpan ke node_modules/.cipincode/bin/ dan agen-nya
 */

import { execSync, spawnSync } from 'node:child_process';
import { createWriteStream, existsSync, mkdirSync, chmodSync, rmSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
import { pipeline } from 'node:stream/promises';
import https from 'node:https';
import os from 'node:os';

const __dirname = dirname(fileURLToPath(import.meta.url));

// ─── Konfigurasi ─────────────────────────────────────────────────────────────

const CIPINCODE_VERSION = process.env.CIPINCODE_VERSION || '0.1.0';
const GITHUB_REPO       = 'cipincode/cipincode-cli';
const RELEASES_BASE     = `https://github.com/${GITHUB_REPO}/releases/download/v${CIPINCODE_VERSION}`;

const MIN_GO_MAJOR = 1;
const MIN_GO_MINOR = 24;

// Direktori install di dalam package (node_modules/cipincode/.cipincode/)
const INSTALL_DIR = join(__dirname, '.cipincode');
const BIN_DIR     = join(INSTALL_DIR, 'bin');
const AGENT_DIR   = join(INSTALL_DIR, 'agent');

// ─── Helper ──────────────────────────────────────────────────────────────────

const isWindows = process.platform === 'win32';

function log(msg)  { process.stdout.write(`\x1b[1;38;5;166m==>\x1b[0m ${msg}\n`); }
function warn(msg) { process.stderr.write(`\x1b[33mPeringatan:\x1b[0m ${msg}\n`); }
function fail(msg) { process.stderr.write(`\x1b[31mError:\x1b[0m ${msg}\n`); process.exit(1); }

function run(cmd, opts = {}) {
  return spawnSync(cmd, { shell: true, stdio: 'inherit', ...opts });
}

function capture(cmd) {
  try {
    return execSync(cmd, { stdio: ['ignore', 'pipe', 'ignore'] }).toString().trim();
  } catch {
    return null;
  }
}

function which(bin) {
  return capture(isWindows ? `where ${bin}` : `command -v ${bin}`);
}

/** Resolusi platform → (goOS, goArch, bunOS, bunArch, ext) */
function platform() {
  const nodePlatform = process.platform; // linux | darwin | win32
  const nodeArch     = process.arch;     // x64 | arm64 | ia32

  const goOS = { linux: 'linux', darwin: 'darwin', win32: 'windows' }[nodePlatform];
  if (!goOS) fail(`Platform tidak didukung: ${nodePlatform}`);

  const goArch = { x64: 'amd64', arm64: 'arm64', ia32: '386' }[nodeArch];
  if (!goArch) fail(`Arsitektur tidak didukung: ${nodeArch}`);

  const bunPlatform = {
    linux:  { x64: 'linux-x64',         arm64: 'linux-aarch64' },
    darwin: { x64: 'darwin-x64',        arm64: 'darwin-aarch64' },
    win32:  { x64: 'windows-x64',       arm64: 'windows-x64' }, // Bun belum arm64 Windows
  }[nodePlatform]?.[nodeArch];
  if (!bunPlatform) fail(`Kombinasi platform/arch tidak didukung: ${nodePlatform}/${nodeArch}`);

  const ext = nodePlatform === 'win32' ? '.exe' : '';

  // Nama asset cipincode binary di GitHub Releases:
  // cipincode-linux-amd64, cipincode-darwin-arm64, cipincode-windows-amd64.exe
  const assetName = `cipincode-${goOS}-${goArch}${ext}`;

  return { goOS, goArch, bunPlatform, ext, assetName };
}

/** Download file dari URL ke destPath dengan redirect support */
async function download(url, destPath) {
  await new Promise((resolve, reject) => {
    function get(u) {
      https.get(u, { headers: { 'User-Agent': 'cipincode-npm-installer' } }, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
          get(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`HTTP ${res.statusCode} dari ${u}`));
          return;
        }
        const dest = createWriteStream(destPath);
        pipeline(res, dest).then(resolve).catch(reject);
      }).on('error', reject);
    }
    get(url);
  });
}

// ─── 1. Cek / Install Go ─────────────────────────────────────────────────────

async function ensureGo() {
  const goPath = which('go');
  if (goPath) {
    const ver = capture('go env GOVERSION');
    if (ver) {
      const m = ver.match(/^go(\d+)\.(\d+)/);
      if (m) {
        const [, major, minor] = m.map(Number);
        if (major > MIN_GO_MAJOR || (major === MIN_GO_MAJOR && minor >= MIN_GO_MINOR)) {
          log(`Go sudah terpasang: ${ver}`);
          return;
        }
        warn(`Go ${ver} terlalu lama (butuh ${MIN_GO_MAJOR}.${MIN_GO_MINOR}+). Menginstall ulang...`);
      }
    }
  }

  log('Menginstall Go...');
  const { goOS, goArch } = platform();

  // Ambil versi Go terbaru
  const goVersion = capture(
    `node -e "const h=require('https');h.get('https://go.dev/VERSION?m=text',r=>{let d='';r.on('data',c=>d+=c);r.on('end',()=>process.stdout.write(d.split('\\n')[0]))})"`
  ) || 'go1.24.0';

  const ext     = goOS === 'windows' ? 'zip' : 'tar.gz';
  const archive = `${goVersion}.${goOS}-${goArch}.${ext}`;
  const dlURL   = `https://go.dev/dl/${archive}`;
  const tmpFile = join(os.tmpdir(), archive);

  log(`Download ${dlURL}`);
  await download(dlURL, tmpFile);

  const goInstallDir = join(os.homedir(), '.local', 'go');
  mkdirSync(goInstallDir, { recursive: true });

  if (goOS === 'windows') {
    run(`tar -xf "${tmpFile}" -C "${join(os.homedir(), '.local')}"`);
  } else {
    run(`tar -xzf "${tmpFile}" -C "${join(os.homedir(), '.local')}"`);
  }

  rmSync(tmpFile, { force: true });

  // Tambahkan ke PATH untuk sesi ini
  const goBin = join(goInstallDir, 'bin');
  process.env.PATH = `${goBin}${isWindows ? ';' : ':'}${process.env.PATH}`;

  const installed = capture('go env GOVERSION');
  log(`Go terpasang: ${installed}`);
  warnPath(goBin, 'Go');
}

// ─── 2. Cek / Install Bun ────────────────────────────────────────────────────

async function ensureBun() {
  // Cek di PATH dan lokasi default
  const candidates = [
    which('bun'),
    join(os.homedir(), '.bun', 'bin', 'bun' + (isWindows ? '.exe' : '')),
    join(os.homedir(), '.local', 'bin', 'bun' + (isWindows ? '.exe' : '')),
  ].filter(Boolean);

  for (const c of candidates) {
    if (existsSync(c)) {
      const ver = capture(`"${c}" --version`);
      if (ver) {
        log(`Bun sudah terpasang: ${ver}`);
        // Pastikan bun ada di PATH untuk sesi ini
        const bunDir = dirname(c);
        if (!process.env.PATH.includes(bunDir)) {
          process.env.PATH = `${bunDir}${isWindows ? ';' : ':'}${process.env.PATH}`;
        }
        return;
      }
    }
  }

  log('Menginstall Bun...');

  if (isWindows) {
    // Windows: PowerShell installer resmi Bun
    const result = run(
      'powershell -Command "irm bun.sh/install.ps1 | iex"',
      { stdio: 'inherit' }
    );
    if (result.status !== 0) fail('Gagal install Bun di Windows.');
  } else {
    // Linux / macOS: curl installer resmi Bun
    const result = run('curl -fsSL https://bun.sh/install | bash', { stdio: 'inherit' });
    if (result.status !== 0) fail('Gagal install Bun.');
  }

  // Tambahkan ke PATH untuk sesi ini
  const bunBin = join(os.homedir(), '.bun', 'bin');
  process.env.PATH = `${bunBin}${isWindows ? ';' : ':'}${process.env.PATH}`;

  const ver = capture('bun --version');
  log(`Bun terpasang: ${ver}`);
  warnPath(bunBin, 'Bun');
}

// ─── 3. Download binary cipincode ────────────────────────────────────────────

async function installCipincode() {
  const { assetName, ext } = platform();

  const binaryURL  = `${RELEASES_BASE}/${assetName}`;
  const agentURL   = `${RELEASES_BASE}/agent.tar.gz`;

  mkdirSync(BIN_DIR,   { recursive: true });
  mkdirSync(AGENT_DIR, { recursive: true });

  // Binary
  const binaryDest = join(BIN_DIR, `cipincode${ext}`);
  log(`Download cipincode binary (${assetName})...`);
  await download(binaryURL, binaryDest);
  if (!isWindows) chmodSync(binaryDest, 0o755);

  // Agent (bundle TypeScript yang dijalankan Bun)
  log('Download agent...');
  const agentTar = join(os.tmpdir(), 'cipincode-agent.tar.gz');
  await download(agentURL, agentTar);
  run(`tar -xzf "${agentTar}" -C "${AGENT_DIR}" --strip-components=1`);
  rmSync(agentTar, { force: true });

  log(`CipinCode ${CIPINCODE_VERSION} terpasang.`);
}

// ─── 4. PATH warning ─────────────────────────────────────────────────────────

function warnPath(dir, label) {
  const shellRcs = ['.bashrc', '.zshrc', '.profile'].map(f => join(os.homedir(), f));
  const exportLine = isWindows ? '' : `export PATH="${dir}:$PATH"`;

  const inPath = process.env.PATH.split(isWindows ? ';' : ':').includes(dir);
  if (!inPath && !isWindows) {
    warn(
      `${label} terpasang di ${dir} tapi belum ada di PATH.\n` +
      `Tambahkan baris berikut ke ~/.bashrc atau ~/.zshrc:\n  ${exportLine}`
    );
  }
}

// ─── Main ─────────────────────────────────────────────────────────────────────

(async () => {
  try {
    log(`Menginstall CipinCode CLI v${CIPINCODE_VERSION}...`);
    // npm path memakai binary prebuilt + agent dist (Node), jadi Go tidak dibutuhkan.
    // Bun hanya fallback bila dist gagal; jangan gagalkan install bila Bun bermasalah.
    try {
      await ensureBun();
    } catch (e) {
      warn(`Bun opsional tidak terpasang (${e.message || e}); CLI tetap jalan via Node.`);
    }
    await installCipincode();

    log('Selesai! Jalankan: cipincode --version');
    process.stdout.write('\nLangkah berikutnya:\n  1. cipincode login\n  2. cipincode\n\n');
  } catch (err) {
    fail(err.message || String(err));
  }
})();

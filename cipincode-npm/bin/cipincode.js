#!/usr/bin/env node
/**
 * bin/cipincode.js – shim yang memanggil binary native cipincode
 *
 * Dicari di urutan:
 *   1. .cipincode/bin/cipincode (di-download saat postinstall)
 *   2. binary cipincode di PATH sistem
 */

import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';
import os from 'node:os';

const __dirname = dirname(fileURLToPath(import.meta.url));
const isWindows = process.platform === 'win32';
const ext = isWindows ? '.exe' : '';

// Lokasi binary yang di-download oleh install.js
const bundledBinary = join(__dirname, '..', '.cipincode', 'bin', `cipincode${ext}`);

// Lokasi agent yang di-download oleh install.js
const agentDir = join(__dirname, '..', '.cipincode', 'agent');

let binary = null;

if (existsSync(bundledBinary)) {
  binary = bundledBinary;
} else {
  // Fallback: coba di PATH
  const inPath = spawnSync(isWindows ? 'where' : 'which', ['cipincode'], {
    stdio: ['ignore', 'pipe', 'ignore'],
  });
  const found = inPath.stdout?.toString().trim().split('\n')[0];
  if (found && existsSync(found) && found !== process.argv[1]) {
    binary = found;
  }
}

if (!binary) {
  process.stderr.write(
    '\x1b[31mError:\x1b[0m Binary cipincode tidak ditemukan.\n' +
    'Coba jalankan ulang: npm install -g cipincode\n'
  );
  process.exit(1);
}

const args = process.argv.slice(2);

// Teruskan env variable lokasi agent agar CLI tau di mana agent berada
const env = {
  ...process.env,
  CIPINCODE_AGENT_DIR: agentDir,
};

const result = spawnSync(binary, args, {
  stdio: 'inherit',
  env,
});

process.exit(result.status ?? 1);

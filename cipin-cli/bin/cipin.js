#!/usr/bin/env node
// Cipin CLI: OpenCode pointed at Cipincode.
//
//   cipin login     sign in through the browser (code exchange)
//   cipin logout    forget the saved key
//   cipin whoami    show the account and its plan
//   cipin topup     open the top-up page
//   cipin ...       anything else goes to OpenCode with the Cipin provider preset
import { spawn, spawnSync } from "node:child_process";
import { createHash, randomBytes } from "node:crypto";
import { chmodSync, existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { createRequire } from "node:module";
import { homedir, hostname } from "node:os";
import { join } from "node:path";
import { createInterface } from "node:readline/promises";

const VERSION = "2.1.4";
const SITE = (process.env.CIPIN_API_URL ?? "https://cipincode.com").replace(/\/+$/, "");
const CLIENT = `cipin-cli/${VERSION}`;
const FALLBACK_MODELS = ["gpt-5", "claude-sonnet-4-5"];

const configDir = join(process.env.XDG_CONFIG_HOME ?? join(homedir(), ".config"), "cipin");
const authFile = join(configDir, "auth.json");
const modelsFile = join(configDir, "models.json");

const fail = (message) => {
    console.error(message);
    process.exit(1);
};

const readAuth = () => {
    try {
        const saved = JSON.parse(readFileSync(authFile, "utf8").replace(/^\uFEFF/, ""));
        return typeof saved.key === "string" && saved.key ? saved : undefined;
    } catch {
        return undefined;
    }
};

const saveAuth = (auth) => {
    mkdirSync(configDir, { recursive: true, mode: 0o700 });
    writeFileSync(authFile, JSON.stringify(auth, null, 2), { mode: 0o600 });

    try {
        chmodSync(authFile, 0o600);
    } catch {
        // Windows ignores POSIX modes.
    }
};

const base64Url = (buffer) => buffer.toString("base64").replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");

const openBrowser = (url) => {
    const [command, args] =
        process.platform === "win32"
            ? ["cmd", ["/c", "start", "", url.replace(/&/g, "^&")]]
            : process.platform === "darwin"
              ? ["open", [url]]
              : ["xdg-open", [url]];

    try {
        spawn(command, args, { stdio: "ignore", detached: true }).on("error", () => {}).unref();
    } catch {
        // The link is printed too.
    }
};

const api = (path, init = {}) =>
    fetch(`${SITE}${path}`, {
        ...init,
        headers: {
            "content-type": "application/json",
            accept: "application/json",
            "x-cipin-client": CLIENT,
            ...init.headers,
        },
        signal: AbortSignal.timeout(15000),
    });

async function login() {
    const verifier = base64Url(randomBytes(48));
    const challenge = base64Url(createHash("sha256").update(verifier).digest());
    const url = `${SITE}/cli/authorize?challenge=${challenge}`;

    console.log("Cipin CLI — masuk ke akun Cipincode.\n");
    console.log("Buka tautan ini di browser, masuk, lalu setujui:\n");
    console.log(`  ${url}\n`);
    openBrowser(url);

    const rl = createInterface({ input: process.stdin, output: process.stdout });
    let code;

    try {
        code = (await rl.question("Tempel kode (XXXX-XXXX): ")).trim();
    } finally {
        rl.close();
    }

    if (!code) {
        return fail("Kode kosong.");
    }

    const response = await api("/api/cli/exchange", {
        method: "POST",
        body: JSON.stringify({ code, verifier, device: hostname() }),
    }).catch((error) => fail(`Tidak bisa menghubungi ${SITE}: ${error.message}`));

    const body = await response.json().catch(() => ({}));

    if (!response.ok || !body.key) {
        return fail(body.message ?? `Login gagal (${response.status}).`);
    }

    saveAuth({ key: body.key, name: body.user?.name, email: body.user?.email });
    console.log(`\nBerhasil masuk sebagai ${body.user?.name ?? "pengguna"} (${body.user?.email ?? "-"}).`);
    console.log(`Cipin CLI hanya menjawab dengan paket Unlimited aktif. Top up: ${SITE}/mbr/dashboard/topup`);
    console.log(`Syarat & ketentuan: ${SITE}/syarat-ketentuan`);
}

async function whoami() {
    const auth = readAuth();

    if (!auth) {
        return fail("Belum masuk. Jalankan: cipin login");
    }

    const headers = { authorization: `Bearer ${auth.key}` };
    const me = await api("/api/cli/me", { headers }).catch(() => undefined);

    if (me?.status === 401) {
        return fail("Sesi tidak berlaku lagi. Jalankan: cipin login");
    }

    console.log(`${auth.name ?? "-"} <${auth.email ?? "-"}>`);

    const usage = await (await api("/api/cli/usage", { headers }).catch(() => undefined))?.json().catch(() => undefined);

    if (usage?.plan) {
        const paid = usage.plan.name === "paid";
        console.log(`Paket: ${usage.plan.label}${usage.plan.expires_at ? ` (sampai ${usage.plan.expires_at})` : ""}`);

        if (!paid) {
            console.log(`Cipin CLI butuh paket Unlimited. Top up: ${SITE}/mbr/dashboard/topup`);
        }
    }
}

function logout() {
    rmSync(authFile, { force: true });
    console.log("Sudah keluar.");
}

function topup() {
    const url = `${SITE}/mbr/dashboard/topup`;
    console.log(url);
    openBrowser(url);
}

async function models(key) {
    try {
        const response = await api("/v1/models", {
            headers: { authorization: `Bearer ${key}` },
            signal: AbortSignal.timeout(4000),
        });
        const ids = ((await response.json()).data ?? []).map((model) => model.id).filter((id) => typeof id === "string");

        if (ids.length) {
            mkdirSync(configDir, { recursive: true, mode: 0o700 });
            writeFileSync(modelsFile, JSON.stringify(ids));
            return ids;
        }
    } catch {
        // Offline or slow: fall back to the last list.
    }

    try {
        if (existsSync(modelsFile)) {
            return JSON.parse(readFileSync(modelsFile, "utf8"));
        }
    } catch {
        // Fall through.
    }

    return FALLBACK_MODELS;
}

/** The OpenCode binary: the platform package that ships with opencode-ai, else `opencode` on PATH. */
function opencodeBinary() {
    const require = createRequire(import.meta.url);
    const os = { win32: "windows", darwin: "darwin", linux: "linux" }[process.platform];
    const arch = { x64: "x64", arm64: "arm64" }[process.arch];
    const exe = process.platform === "win32" ? "opencode.exe" : "opencode";

    if (os && arch) {
        for (const name of [`opencode-${os}-${arch}`, `opencode-${os}-${arch}-baseline`]) {
            try {
                const binary = join(require.resolve(`${name}/package.json`, { paths: [require.resolve("opencode-ai/package.json")] }), "..", "bin", exe);

                if (existsSync(binary)) {
                    return binary;
                }
            } catch {
                // Try the next one.
            }
        }
    }

    return "opencode";
}

async function main() {
    const [command, ...rest] = process.argv.slice(2);

    switch (command) {
        case "login":
            await login();
            return process.exit(0);
        case "logout":
            logout();
            return process.exit(0);
        case "whoami":
            await whoami();
            return process.exit(0);
        case "topup":
            topup();
            return process.exit(0);
    }

    const bypass = command !== undefined && ["--version", "-v", "--help", "-h", "help"].includes(command) && rest.length === 0;
    const env = { ...process.env };
    const binary = opencodeBinary();

    if (!bypass) {
        const auth = readAuth();

        if (!auth) {
            return fail("Belum masuk. Jalankan: cipin login");
        }

        const ids = await models(auth.key);

        env.OPENCODE_CONFIG_CONTENT = JSON.stringify({
            provider: {
                cipin: {
                    npm: "@ai-sdk/openai-compatible",
                    name: "Cipin",
                    options: { baseURL: `${SITE}/v1`, apiKey: auth.key, headers: { "x-cipin-client": CLIENT } },
                    models: Object.fromEntries(ids.map((id) => [id, { name: id }])),
                },
            },
            enabled_providers: ["cipin"],
            model: `cipin/${ids[0]}`,
            autoupdate: false,
            share: "disabled",
        });
        env.OPENCODE_DISABLE_MODELS_FETCH = "1";
        env.OPENCODE_DISABLE_AUTOUPDATE = "1";
        // Own data folder, so it never shares a database with another OpenCode install.
        env.XDG_DATA_HOME ??= join(homedir(), ".local", "share", "cipin");
    }

    if (command === "--version" || command === "-v") {
        const result = spawnSync(binary, ["--version"], { encoding: "utf8", shell: false });
        console.log(`cipin ${VERSION} (opencode ${(result.stdout ?? "").trim() || "?"})`);
        return process.exit(0);
    }

    const child = spawn(binary, process.argv.slice(2), { stdio: "inherit", env });

    child.on("error", (error) => fail(`Tidak bisa menjalankan OpenCode: ${error.message}`));
    child.on("exit", (code, signal) => process.exit(code ?? (signal ? 1 : 0)));
}

try {
    await main();
} catch (error) {
    // Ctrl+C while waiting for input: leave quietly instead of printing a stack trace.
    if (error?.code === "ABORT_ERR") {
        console.log("\nDibatalkan.");
        process.exit(130);
    }

    fail(error?.message ?? String(error));
}

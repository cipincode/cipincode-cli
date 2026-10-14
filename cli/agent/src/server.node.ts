import { createServer, type IncomingMessage, type ServerResponse } from 'node:http';
import { parseEffort, parseMode, run, type Credentials, type Message } from './agent.js';

const port = Number(process.env.PORT ?? 4096);
const version = process.env.CIPINCODE_VERSION ?? '0.1.0';

// Jika diisi, /chat hanya melayani permintaan yang membawa token ini. CLI membuat token
// acak saat menjalankan agent, sehingga halaman web atau proses lain di komputer ini
// tidak bisa memakai agent (dan API key pengguna) lewat port lokal.
const agentToken = process.env.CIPINCODE_AGENT_TOKEN ?? '';

const encoder = new TextEncoder();
const sseLine = (event: unknown) => encoder.encode(`data: ${JSON.stringify(event)}\n\n`);

function credentialsFromHeaders(headers: IncomingMessage['headers']): Credentials | null {
    const url = (headers['x-cipincode-url'] as string | undefined) ?? '';
    const key = (headers['x-cipincode-key'] as string | undefined) ?? '';

    return /^https?:\/\//.test(url) && key !== '' ? { url, key } : null;
}

function readJsonBody(req: IncomingMessage): Promise<Record<string, unknown>> {
    return new Promise((resolve) => {
        let raw = '';
        req.on('data', (chunk: Buffer) => {
            raw += chunk.toString('utf8');
            // Batasi 2MB agar tidak OOM dari request jahat lokal.
            if (raw.length > 2_000_000) req.destroy();
        });
        req.on('end', () => {
            try {
                resolve(raw === '' ? {} : (JSON.parse(raw) as Record<string, unknown>));
            } catch {
                resolve({});
            }
        });
        req.on('error', () => resolve({}));
    });
}

const server = createServer(async (req: IncomingMessage, res: ServerResponse) => {
    const url = new URL(req.url ?? '/', 'http://127.0.0.1');

    if (url.pathname === '/health' && req.method === 'GET') {
        res.writeHead(200, { 'content-type': 'application/json' });
        res.end(JSON.stringify({ ok: true, name: 'cipincode-agent', version }));
        return;
    }

    if (url.pathname === '/chat' && req.method === 'POST') {
        if (agentToken !== '' && req.headers['x-agent-token'] !== agentToken) {
            res.writeHead(403, { 'content-type': 'application/json' });
            res.end(JSON.stringify({ error: 'forbidden' }));
            return;
        }

        const body = await readJsonBody(req);
        const message = typeof body.message === 'string' ? body.message.trim() : '';

        if (message === '') {
            res.writeHead(400, { 'content-type': 'application/json' });
            res.end(JSON.stringify({ error: 'message is required' }));
            return;
        }

        const credentials = credentialsFromHeaders(req.headers);

        res.writeHead(200, {
            'content-type': 'text/event-stream',
            'cache-control': 'no-cache',
            connection: 'keep-alive',
        });

        try {
            for await (const event of run({
                message,
                history: Array.isArray(body.history) ? (body.history as Message[]) : [],
                mode: parseMode(body.mode),
                effort: parseEffort(body.effort),
                session: typeof body.session === 'string' ? body.session : undefined,
                agent: typeof body.agent === 'number' ? body.agent : undefined,
                credentials,
            })) {
                if (!res.write(sseLine(event))) {
                    // Tunggu drain agar tidak menumpuk memori saat CLI lambat membaca.
                    await new Promise<void>((resolve) => res.once('drain', () => resolve()));
                }
            }
        } catch (error) {
            res.write(
                sseLine({
                    type: 'error',
                    message: error instanceof Error ? error.message : String(error),
                }),
            );
        } finally {
            res.end();
        }
        return;
    }

    res.writeHead(404, { 'content-type': 'text/plain' });
    res.end('Not found');
});

server.listen(port, '127.0.0.1', () => {
    console.log(`cipincode-agent listening on http://127.0.0.1:${port}`);
});

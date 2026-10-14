import {
    parseEffort,
    parseMode,
    run,
    type Credentials,
    type Message,
} from './agent';

const port = Number(process.env.PORT ?? 4096);
const version = process.env.CIPINCODE_VERSION ?? '0.1.0';

// Jika diisi, /chat hanya melayani permintaan yang membawa token ini. CLI membuat token
// acak saat menjalankan agent, sehingga halaman web atau proses lain di komputer ini
// tidak bisa memakai agent (dan API key pengguna) lewat port lokal.
const agentToken = process.env.CIPINCODE_AGENT_TOKEN ?? '';

const encoder = new TextEncoder();
const sse = (event: unknown) =>
    encoder.encode(`data: ${JSON.stringify(event)}\n\n`);

function credentialsFrom(req: Request): Credentials | null {
    const url = req.headers.get('x-cipincode-url') ?? '';
    const key = req.headers.get('x-cipincode-key') ?? '';

    return /^https?:\/\//.test(url) && key !== '' ? { url, key } : null;
}

const server = Bun.serve({
    hostname: '127.0.0.1',
    port,
    idleTimeout: 120,
    async fetch(req) {
        const url = new URL(req.url);

        if (url.pathname === '/health') {
            return Response.json({
                ok: true,
                name: 'cipincode-agent',
                version,
            });
        }

        if (url.pathname === '/chat' && req.method === 'POST') {
            if (
                agentToken !== '' &&
                req.headers.get('x-agent-token') !== agentToken
            ) {
                return Response.json({ error: 'forbidden' }, { status: 403 });
            }

            const body = (await req.json().catch(() => ({}))) as {
                message?: unknown;
                history?: Message[];
                mode?: unknown;
                effort?: unknown;
                provider?: unknown;
                model?: unknown;
                provider_key?: unknown;
                session?: unknown;
                agent?: unknown;
            };
            const message =
                typeof body.message === 'string' ? body.message.trim() : '';

            if (message === '') {
                return Response.json(
                    { error: 'message is required' },
                    { status: 400 },
                );
            }

            const credentials = credentialsFrom(req);
            const stream = new ReadableStream({
                async start(controller) {
                    try {
                        for await (const event of run({
                            message,
                            history: Array.isArray(body.history)
                                ? body.history
                                : [],
                            mode: parseMode(body.mode),
                            effort: parseEffort(body.effort),
                            provider:
                                typeof body.provider === 'string' &&
                                body.provider !== ''
                                    ? body.provider
                                    : undefined,
                            model:
                                typeof body.model === 'string' &&
                                body.model !== ''
                                    ? body.model
                                    : undefined,
                            providerKey:
                                typeof body.provider_key === 'string' &&
                                body.provider_key !== ''
                                    ? body.provider_key
                                    : undefined,
                            session:
                                typeof body.session === 'string'
                                    ? body.session
                                    : undefined,
                            agent:
                                typeof body.agent === 'number'
                                    ? body.agent
                                    : undefined,
                            credentials,
                        })) {
                            controller.enqueue(sse(event));
                        }
                    } catch (error) {
                        controller.enqueue(
                            sse({
                                type: 'error',
                                message:
                                    error instanceof Error
                                        ? error.message
                                        : String(error),
                            }),
                        );
                    } finally {
                        controller.close();
                    }
                },
            });

            return new Response(stream, {
                headers: {
                    'content-type': 'text/event-stream',
                    'cache-control': 'no-cache',
                },
            });
        }

        return new Response('Not found', { status: 404 });
    },
});

console.log(`cipincode-agent listening on ${server.url}`);

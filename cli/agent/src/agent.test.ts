import { afterEach, beforeEach, describe, expect, test } from 'bun:test';
import {
    parseEffort,
    parseMode,
    parseStream,
    run,
    type AgentEvent,
    type Credentials,
} from './agent';

const encoder = new TextEncoder();

function sseBody(...events: unknown[]): ReadableStream<Uint8Array> {
    return new ReadableStream({
        start(controller) {
            for (const event of events) {
                controller.enqueue(
                    encoder.encode(`data: ${JSON.stringify(event)}\n\n`),
                );
            }
            controller.close();
        },
    });
}

async function collect(gen: AsyncGenerator<AgentEvent>): Promise<AgentEvent[]> {
    const out: AgentEvent[] = [];
    for await (const event of gen) out.push(event);
    return out;
}

const credentials: Credentials = { url: 'http://dash.test/', key: 'cpc_abc' };

type Handler = (url: string, init: RequestInit) => Response | Promise<Response>;
const fakeFetch = (handler: Handler) =>
    ((url: string, init: RequestInit) =>
        Promise.resolve(handler(url, init))) as unknown as typeof fetch;

beforeEach(() => {
    delete process.env.CIPINCODE_AGENT_STUB;
});
afterEach(() => {
    delete process.env.CIPINCODE_AGENT_STUB;
});

describe('run', () => {
    test('forwards the server stream and sends credentials, mode and history', async () => {
        let seen: { url: string; init: RequestInit } | undefined;
        const events = await collect(
            run(
                {
                    message: 'halo',
                    mode: 'coding',
                    effort: 'ultra',
                    history: [{ role: 'user', content: 'sebelumnya' }],
                    credentials,
                },
                fakeFetch((url, init) => {
                    seen = { url, init };
                    return new Response(
                        sseBody(
                            {
                                type: 'usage',
                                session: { used: 1 },
                                weekly: { used: 1 },
                            },
                            { type: 'token', text: 'Halo' },
                            { type: 'token', text: ' dunia' },
                            { type: 'done' },
                        ),
                        { status: 200 },
                    );
                }),
            ),
        );

        expect(events.map((e) => e.type)).toEqual([
            'usage',
            'token',
            'token',
            'done',
        ]);
        const request = seen as unknown as { url: string; init: RequestInit };
        expect(request.url).toBe('http://dash.test/api/cli/chat');
        expect(
            (request.init.headers as Record<string, string>).Authorization,
        ).toBe('Bearer cpc_abc');
        expect(JSON.parse(request.init.body as string)).toEqual({
            message: 'halo',
            mode: 'coding',
            effort: 'ultra',
            history: [{ role: 'user', content: 'sebelumnya' }],
        });
    });

    test('forwards provider and model when chosen', async () => {
        let seen: { url: string; init: RequestInit } | undefined;
        await collect(
            run(
                {
                    message: 'halo',
                    provider: 'anthropic',
                    model: 'anthropic/claude-sonnet-4-5',
                    credentials,
                },
                fakeFetch((url, init) => {
                    seen = { url, init };
                    return new Response(
                        sseBody({ type: 'done' }),
                        { status: 200 },
                    );
                }),
            ),
        );

        const request = seen as unknown as { url: string; init: RequestInit };
        expect(JSON.parse(request.init.body as string)).toEqual({
            message: 'halo',
            mode: 'chat',
            effort: 'medium',
            history: [],
            provider: 'anthropic',
            model: 'anthropic/claude-sonnet-4-5',
        });
    });

    test('streams directly from Anthropic with the user key', async () => {
        let seen: { url: string; init: RequestInit } | undefined;
        const events = await collect(
            run(
                {
                    message: 'halo',
                    provider: 'anthropic',
                    model: 'anthropic/claude-sonnet-4-5',
                    providerKey: 'sk-ant-x',
                    credentials,
                },
                fakeFetch((url, init) => {
                    seen = { url, init };
                    return new Response(
                        sseBody(
                            {
                                type: 'content_block_delta',
                                delta: { type: 'text_delta', text: 'Hai' },
                            },
                            {
                                type: 'content_block_delta',
                                delta: { type: 'text_delta', text: ' juga' },
                            },
                        ),
                        { status: 200 },
                    );
                }),
            ),
        );

        expect(events).toEqual([
            { type: 'token', text: 'Hai' },
            { type: 'token', text: ' juga' },
            { type: 'done' },
        ]);
        const request = seen as unknown as { url: string; init: RequestInit };
        expect(request.url).toBe('https://api.anthropic.com/v1/messages');
        expect(
            (request.init.headers as Record<string, string>)['x-api-key'],
        ).toBe('sk-ant-x');
        expect(JSON.parse(request.init.body as string).model).toBe(
            'claude-sonnet-4-5-20250929',
        );
    });

    test('streams directly from OpenAI-compatible providers', async () => {
        let seen: { url: string; init: RequestInit } | undefined;
        const events = await collect(
            run(
                {
                    message: 'halo',
                    provider: 'openai',
                    model: 'openai/gpt-4o',
                    providerKey: 'sk-x',
                    credentials,
                },
                fakeFetch((url, init) => {
                    seen = { url, init };
                    return new Response(
                        sseBody({
                            choices: [{ delta: { content: 'Halo' } }],
                        }),
                        { status: 200 },
                    );
                }),
            ),
        );

        expect(events).toEqual([
            { type: 'token', text: 'Halo' },
            { type: 'done' },
        ]);
        const request = seen as unknown as { url: string; init: RequestInit };
        expect(request.url).toBe(
            'https://api.openai.com/v1/chat/completions',
        );
        expect(
            (request.init.headers as Record<string, string>).Authorization,
        ).toBe('Bearer sk-x');
        // Kunci provider tidak boleh ikut ke dashboard (jalur langsung).
        expect(request.url).not.toContain('dash.test');
    });

    test('rejects a wrong direct key with a clear message', async () => {
        const events = await collect(
            run(
                {
                    message: 'halo',
                    provider: 'openai',
                    model: 'openai/gpt-4o',
                    providerKey: 'salah',
                    credentials,
                },
                fakeFetch(
                    () =>
                        new Response(
                            JSON.stringify({
                                error: { message: 'Incorrect API key' },
                            }),
                            { status: 401 },
                        ),
                ),
            ),
        );

        expect(events).toEqual([
            {
                type: 'error',
                message:
                    'Kunci API OpenAI ditolak. Periksa kunci di /connect.',
            },
        ]);
    });

    test('asks the user to log in when there are no credentials', async () => {        const events = await collect(
            run(
                { message: 'halo' },
                fakeFetch(() => new Response('x')),
            ),
        );

        expect(events).toEqual([
            {
                type: 'error',
                message: 'Belum login. Ketik /login untuk masuk.',
            },
        ]);
    });

    test('turns a 429 into a blocked event carrying the usage snapshot', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() =>
                    Response.json(
                        {
                            message: 'Batas pemakaian sesi 5 jam tercapai.',
                            blocked: 'session',
                            session: { used: 40, exhausted: true },
                            weekly: { used: 40 },
                        },
                        { status: 429 },
                    ),
                ),
            ),
        );

        expect(events).toHaveLength(1);
        expect(events[0]).toMatchObject({
            type: 'blocked',
            blocked: 'session',
            message: 'Batas pemakaian sesi 5 jam tercapai.',
            session: { exhausted: true },
        });
    });

    test('turns a 403 with a blocked reason (expired trial) into a blocked event', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() =>
                    Response.json(
                        {
                            message:
                                'Masa percobaan gratis 14 hari telah berakhir.',
                            blocked: 'plan',
                            plan: { expired: true },
                        },
                        { status: 403 },
                    ),
                ),
            ),
        );

        expect(events).toHaveLength(1);
        expect(events[0]).toMatchObject({
            type: 'blocked',
            blocked: 'plan',
            plan: { expired: true },
        });
    });

    test('a plain 403 without a blocked reason stays an error', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() =>
                    Response.json({ message: 'Dilarang.' }, { status: 403 }),
                ),
            ),
        );

        expect(events).toEqual([{ type: 'error', message: 'Dilarang.' }]);
    });

    test('reports an invalid login on 401', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() => new Response('', { status: 401 })),
            ),
        );

        expect(events[0]).toMatchObject({ type: 'error' });
        expect((events[0] as { message: string }).message).toContain('/login');
    });

    test('surfaces the server message for other failures', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() =>
                    Response.json(
                        { message: 'Layanan AI belum dikonfigurasi.' },
                        { status: 503 },
                    ),
                ),
            ),
        );

        expect(events).toEqual([
            { type: 'error', message: 'Layanan AI belum dikonfigurasi.' },
        ]);
    });

    test('falls back to the first validation error', async () => {
        const events = await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch(() =>
                    Response.json(
                        { errors: { message: ['Pesan terlalu panjang.'] } },
                        { status: 422 },
                    ),
                ),
            ),
        );

        expect(events).toEqual([
            { type: 'error', message: 'Pesan terlalu panjang.' },
        ]);
    });

    test('reports a network failure', async () => {
        const failing = (() =>
            Promise.reject(
                new Error('ECONNREFUSED'),
            )) as unknown as typeof fetch;
        const events = await collect(
            run({ message: 'halo', credentials }, failing),
        );

        expect(events).toEqual([
            {
                type: 'error',
                message: 'Tidak bisa terhubung ke http://dash.test/',
            },
        ]);
    });

    test('stub mode answers without any network access', async () => {
        process.env.CIPINCODE_AGENT_STUB = '1';
        const never = (() => {
            throw new Error('tidak boleh memanggil jaringan');
        }) as unknown as typeof fetch;

        const events = await collect(
            run({ message: 'perbaiki', mode: 'coding' }, never),
        );
        const text = events
            .map((e) => (e.type === 'token' ? e.text : ''))
            .join('');

        expect(text).toContain('mode coding');
        expect(text).toContain('perbaiki');
        expect(events.at(-1)).toEqual({ type: 'done' });
    });
});

describe('parseStream', () => {
    test('stops at the first terminal event', async () => {
        const events = await collect(
            parseStream(
                sseBody(
                    { type: 'token', text: 'a' },
                    { type: 'error', message: 'rusak' },
                    { type: 'token', text: 'b' },
                ),
            ),
        );

        expect(events.map((e) => e.type)).toEqual(['token', 'error']);
    });

    test('ends with done when the stream closes without a terminal event', async () => {
        const events = await collect(
            parseStream(sseBody({ type: 'token', text: 'a' })),
        );

        expect(events.map((e) => e.type)).toEqual(['token', 'done']);
    });

    test('handles events split across chunks and ignores noise', async () => {
        const body = new ReadableStream<Uint8Array>({
            start(controller) {
                const text = ': ping\n\ndata: {"type":"token","te';
                controller.enqueue(encoder.encode(text));
                controller.enqueue(
                    encoder.encode(
                        'xt":"Hai"}\n\ndata: tidak-json\n\ndata: {"type":"done"}\n\n',
                    ),
                );
                controller.close();
            },
        });

        const events = await collect(parseStream(body));

        expect(events).toEqual([
            { type: 'token', text: 'Hai' },
            { type: 'done' },
        ]);
    });
});

describe('effort', () => {
    test('defaults to medium when not given', async () => {
        let sent: Record<string, unknown> = {};
        await collect(
            run(
                { message: 'halo', credentials },
                fakeFetch((_url, init) => {
                    sent = JSON.parse(init.body as string);
                    return new Response(sseBody({ type: 'done' }), {
                        status: 200,
                    });
                }),
            ),
        );

        expect(sent.effort).toBe('medium');
    });

    test('parseEffort accepts the five levels and falls back to medium', () => {
        for (const level of ['low', 'medium', 'high', 'xhigh', 'ultra']) {
            expect(parseEffort(level)).toBe(level as never);
        }
        expect(parseEffort('mega')).toBe('medium');
        expect(parseEffort(undefined)).toBe('medium');
    });

    test('stub mode mentions the effort', async () => {
        process.env.CIPINCODE_AGENT_STUB = '1';
        const events = await collect(run({ message: 'x', effort: 'xhigh' }));

        expect(
            events.map((e) => (e.type === 'token' ? e.text : '')).join(''),
        ).toContain('effort xhigh');
    });
});

describe('parseMode', () => {
    test('falls back to chat for unknown values', () => {
        expect(parseMode('debugging')).toBe('debugging');
        expect(parseMode('hack')).toBe('chat');
        expect(parseMode(undefined)).toBe('chat');
    });
});

import type { AgentEvent, Effort, Message, Mode } from './agent';

export type DirectRequest = {
    provider: string;
    providerKey: string;
    model: string;
    message: string;
    history: Message[];
    mode: Mode;
    effort: Effort;
};

const OPENAI_COMPAT_BASES: Record<string, string> = {
    openai: 'https://api.openai.com/v1',
    xai: 'https://api.x.ai/v1',
    deepseek: 'https://api.deepseek.com/v1',
    google: 'https://generativelanguage.googleapis.com/v1beta/openai',
};

const ANTHROPIC_BASE = 'https://api.anthropic.com';

/** OpenRouter ID -> ID API langsung Anthropic. */
const ANTHROPIC_DIRECT: Record<string, string> = {
    'anthropic/claude-sonnet-4-5': 'claude-sonnet-4-5-20250929',
    'anthropic/claude-opus-4-1': 'claude-opus-4-1-20250805',
    'anthropic/claude-3-5-haiku': 'claude-3-5-haiku-20241022',
};

const MAX_TOKENS: Record<Effort, number> = {
    low: 1024,
    medium: 2048,
    high: 4096,
    xhigh: 8192,
    ultra: 16384,
};

const SYSTEMS: Record<Mode, string> = {
    chat: 'You are a helpful coding assistant. Answer concisely in Indonesian.',
    coding: 'You are an expert coding assistant. Write clean, working code with brief explanations in Indonesian.',
    debugging:
        'You are a debugging expert. Analyze errors deeply, find root causes, and explain fixes in Indonesian.',
};

/** Base URL API langsung, atau null bila provider hanya via server. */
export function directBase(provider: string): string | null {
    if (provider === 'anthropic') return ANTHROPIC_BASE;
    return OPENAI_COMPAT_BASES[provider] ?? null;
}

/** Petakan ID model OpenRouter ke ID API langsung provider. */
export function directModel(provider: string, model: string): string {
    if (provider === 'anthropic') {
        return (
            ANTHROPIC_DIRECT[model] ?? model.replace(/^anthropic\//, '')
        );
    }
    const parts = model.split('/');
    return parts.length > 1 ? parts[parts.length - 1] : model;
}

function toMessages(req: DirectRequest): { role: string; content: string }[] {
    return [
        ...req.history.map((h) => ({ role: h.role, content: h.content })),
        { role: 'user', content: req.message },
    ];
}

async function* readSseLines(
    body: ReadableStream<Uint8Array>,
): AsyncGenerator<unknown> {
    const reader = body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    try {
        while (true) {
            const { done, value } = await reader.read();
            if (done) break;
            buffer += decoder.decode(value, { stream: true });
            let newline: number;
            while ((newline = buffer.indexOf('\n')) !== -1) {
                const line = buffer.slice(0, newline).trim();
                buffer = buffer.slice(newline + 1);
                if (!line.startsWith('data:')) continue;
                const payload = line.slice(5).trim();
                if (payload === '[DONE]') return;
                try {
                    yield JSON.parse(payload);
                } catch {
                    continue;
                }
            }
        }
    } finally {
        reader.releaseLock();
    }
}

const DISPLAY_NAMES: Record<string, string> = {
    anthropic: 'Anthropic',
    openai: 'OpenAI',
    google: 'Google',
    xai: 'xAI',
    deepseek: 'DeepSeek',
    meta: 'Meta',
};

function providerError(
    provider: string,
    status: number,
    detail?: string,
): AgentEvent {
    const name = DISPLAY_NAMES[provider] ?? provider;
    if (status === 401 || status === 403) {
        return {
            type: 'error',
            message: `Kunci API ${name} ditolak. Periksa kunci di /connect.`,
        };
    }
    if (status === 429) {
        return {
            type: 'error',
            message: `Rate limit ${name}. Tunggu sebentar lalu coba lagi.`,
        };
    }
    return {
        type: 'error',
        message: String(
            detail ?? `${name} membalas status ${status}.`,
        ),
    };
}

async function* streamAnthropic(
    req: DirectRequest,
    fetchImpl: typeof fetch,
): AsyncGenerator<AgentEvent> {
    let response: Response;
    try {
        response = await fetchImpl(`${ANTHROPIC_BASE}/v1/messages`, {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
                'x-api-key': req.providerKey,
                'anthropic-version': '2023-06-01',
            },
            body: JSON.stringify({
                model: directModel(req.provider, req.model),
                max_tokens: MAX_TOKENS[req.effort],
                system: SYSTEMS[req.mode],
                messages: toMessages(req),
                stream: true,
            }),
        });
    } catch {
        yield { type: 'error', message: 'Tidak bisa terhubung ke Anthropic.' };
        return;
    }

    if (!response.ok || !response.body) {
        const body = await response.json().catch(() => ({}));
        const detail =
            (body as { error?: { message?: string } }).error?.message ??
            undefined;
        yield providerError(req.provider, response.status, detail);
        return;
    }

    for await (const event of readSseLines(response.body)) {
        const text = (event as { delta?: { text?: string } }).delta?.text;
        if (typeof text === 'string' && text !== '') {
            yield { type: 'token', text };
        }
    }
    yield { type: 'done' };
}

async function* streamOpenAICompat(
    req: DirectRequest,
    fetchImpl: typeof fetch,
): AsyncGenerator<AgentEvent> {
    const base = OPENAI_COMPAT_BASES[req.provider];
    let response: Response;
    try {
        response = await fetchImpl(`${base}/chat/completions`, {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
                Authorization: `Bearer ${req.providerKey}`,
            },
            body: JSON.stringify({
                model: directModel(req.provider, req.model),
                max_tokens: MAX_TOKENS[req.effort],
                messages: [
                    { role: 'system', content: SYSTEMS[req.mode] },
                    ...toMessages(req),
                ],
                stream: true,
            }),
        });
    } catch {
        yield {
            type: 'error',
            message: `Tidak bisa terhubung ke ${req.provider}.`,
        };
        return;
    }

    if (!response.ok || !response.body) {
        const body = await response.json().catch(() => ({}));
        const detail = (body as { error?: { message?: string } }).error
            ?.message;
        if (response.status === 401 || response.status === 403) {
            yield providerError(req.provider, response.status);
            return;
        }
        const message =
            typeof detail === 'string' && detail !== ''
                ? detail
                : undefined;
        yield providerError(req.provider, response.status, message);
        return;
    }

    for await (const event of readSseLines(response.body)) {
        const delta = (
            event as {
                choices?: { delta?: { content?: string } }[];
            }
        ).choices?.[0]?.delta?.content;
        if (typeof delta === 'string' && delta !== '') {
            yield { type: 'token', text: delta };
        }
    }
    yield { type: 'done' };
}

/**
 * Jalankan percakapan langsung ke API provider dengan kunci milik pengguna
 * (BYOK). Tidak lewat server dashboard sehingga tidak memotong kuota.
 */
export async function* streamDirect(
    req: DirectRequest,
    fetchImpl: typeof fetch = fetch,
): AsyncGenerator<AgentEvent> {
    if (directBase(req.provider) === null) {
        yield {
            type: 'error',
            message: 'Provider ini hanya tersedia via server CipinCode.',
        };
        return;
    }
    if (req.provider === 'anthropic') {
        yield* streamAnthropic(req, fetchImpl);
        return;
    }
    yield* streamOpenAICompat(req, fetchImpl);
}

export const modes = ['chat', 'coding', 'debugging'];
export function parseMode(value) {
    return modes.find((m) => m === value) ?? 'chat';
}
export const efforts = ['low', 'medium', 'high', 'xhigh', 'ultra'];
export function parseEffort(value) {
    return efforts.find((e) => e === value) ?? 'medium';
}
const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
/**
 * Jalankan satu giliran percakapan lewat server dashboard, yang meneruskannya ke
 * OpenRouter. Kunci OpenRouter hanya ada di server, tidak pernah di komputer pengguna.
 *
 * Dengan CIPINCODE_AGENT_STUB=1 agent memakai balasan contoh tanpa jaringan
 * (untuk pengembangan dan pengujian).
 */
export async function* run(input, fetchImpl = fetch) {
    if (process.env.CIPINCODE_AGENT_STUB === '1') {
        yield* stub(input);
        return;
    }
    const { credentials } = input;
    if (!credentials) {
        yield {
            type: 'error',
            message: 'Belum login. Ketik /login untuk masuk.',
        };
        return;
    }
    let response;
    try {
        response = await fetchImpl(`${credentials.url.replace(/\/+$/, '')}/api/cli/chat`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${credentials.key}`,
                'Content-Type': 'application/json',
                Accept: 'text/event-stream',
            },
            body: JSON.stringify({
                message: input.message,
                mode: input.mode ?? 'chat',
                effort: input.effort ?? 'medium',
                history: input.history ?? [],
                ...(input.provider
                    ? { provider: input.provider }
                    : null),
                ...(input.model ? { model: input.model } : null),
            }),
        });
    }
    catch {
        yield {
            type: 'error',
            message: `Tidak bisa terhubung ke ${credentials.url}`,
        };
        return;
    }
    if (response.status === 401) {
        yield {
            type: 'error',
            message: 'Login Anda tidak valid lagi. Ketik /login untuk masuk kembali.',
        };
        return;
    }
    // 429: jatah habis. 403 dengan "blocked": masa aktif paket berakhir.
    if (response.status === 429 || response.status === 403) {
        const body = await response.json().catch(() => ({}));
        if (response.status === 429 || body.blocked) {
            yield {
                type: 'blocked',
                ...body,
                message: body.message ?? 'Batas pemakaian tercapai.',
            };
            return;
        }
        yield {
            type: 'error',
            message: String(body.message ?? 'Akses ditolak.'),
        };
        return;
    }
    if (!response.ok || !response.body) {
        const body = await response.json().catch(() => ({}));
        const detail = body.message ?? Object.values(body.errors ?? {}).flat()[0];
        yield {
            type: 'error',
            message: String(detail ?? `Server membalas status ${response.status}`),
        };
        return;
    }
    yield* parseStream(response.body);
}
/** Baca aliran SSE dari server dan teruskan peristiwanya. */
export async function* parseStream(body) {
    const reader = body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    try {
        while (true) {
            const { done, value } = await reader.read();
            if (done)
                break;
            buffer += decoder.decode(value, { stream: true });
            let newline;
            while ((newline = buffer.indexOf('\n')) !== -1) {
                const line = buffer.slice(0, newline).trim();
                buffer = buffer.slice(newline + 1);
                if (!line.startsWith('data:'))
                    continue;
                let event;
                try {
                    event = JSON.parse(line.slice(5).trim());
                }
                catch {
                    continue;
                }
                yield event;
                if (event.type === 'done' || event.type === 'error')
                    return;
            }
        }
    }
    finally {
        reader.releaseLock();
    }
    yield { type: 'done' };
}
async function* stub(input) {
    const reply = `Balasan contoh (mode ${input.mode ?? 'chat'}, effort ${input.effort ?? 'medium'}): "${input.message}"`;
    for (const word of reply.split(/(\s+)/)) {
        if (word === '')
            continue;
        await sleep(5);
        yield { type: 'token', text: word };
    }
    yield { type: 'done' };
}

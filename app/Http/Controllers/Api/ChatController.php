<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Agent;
use App\Models\ChatSession;
use App\Models\User;
use App\Services\UsageTracker;
use App\Support\ChatPrompts;
use App\Support\Effort;
use Illuminate\Http\Client\Response;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Str;
use Illuminate\Validation\Rule;
use Psr\Http\Message\StreamInterface;
use Symfony\Component\HttpFoundation\StreamedResponse;

class ChatController extends Controller
{
    private const HISTORY_LIMIT = 20;

    /**
     * Teruskan percakapan ke OpenRouter dan alirkan jawabannya ke CLI sebagai SSE:
     * `usage` (sekali di awal), lalu `token`, diakhiri `done` atau `error`.
     *
     * Jatah hanya dipotong setelah OpenRouter menerima permintaan, jadi galat di sisi
     * layanan AI tidak menghabiskan jatah pengguna.
     */
    public function store(Request $request, UsageTracker $usage): JsonResponse|StreamedResponse
    {
        $data = $request->validate([
            'message' => ['required', 'string', 'max:8000'],
            'mode' => ['nullable', 'in:'.implode(',', ChatPrompts::MODES)],
            'effort' => ['nullable', 'in:'.implode(',', Effort::LEVELS)],
            'session' => ['nullable', 'uuid'],
            'agent' => ['nullable', 'integer', Rule::exists('agents', 'id')->where('user_id', $request->user()->id)],
            'history' => ['nullable', 'array', 'max:100'],
            'history.*.role' => ['required', 'in:user,assistant'],
            'history.*.content' => ['required', 'string', 'max:8000'],
            'provider' => ['nullable', 'string', 'max:32'],
            'model' => ['nullable', 'string', 'max:200'],
        ]);

        $key = config('services.openrouter.key');
        if (blank($key)) {
            return response()->json([
                'message' => 'Layanan AI belum dikonfigurasi di server (OPENROUTER_API_KEY belum diisi).',
            ], 503);
        }

        $user = $request->user();

        // Tolak lebih awal bila jatah sudah habis, sebelum menghubungi OpenRouter.
        if ($blocked = $usage->blockedBy($usage->snapshot($user))) {
            return $this->blocked($usage, $user, $blocked);
        }

        $mode = $data['mode'] ?? 'chat';
        $effort = Effort::normalize($data['effort'] ?? null);
        $model = $this->model($user, $mode, $data['provider'] ?? null, $data['model'] ?? null);
        $agent = isset($data['agent']) ? $user->agents()->find($data['agent']) : null;
        $upstream = $this->openRouter($key, $model, $mode, $effort, $agent, $data['message'], $data['history'] ?? []);

        if (! $upstream->successful()) {
            return response()->json([
                'message' => $this->upstreamMessage($upstream->status(), $upstream->json('error.message')),
            ], $upstream->status() === 429 ? 503 : 502);
        }

        $stream = $upstream->toPsrResponse()->getBody();

        // Potong jatah sekarang; kalau ternyata habis (permintaan bersamaan), batalkan.
        [$allowed, $snapshot, $blocked] = $usage->consume($user);
        if (! $allowed) {
            $stream->close();

            return $this->blocked($usage, $user, $blocked);
        }

        // Riwayat hanya dicatat untuk pesan yang benar-benar diproses (jatah sudah terpotong).
        $session = isset($data['session'])
            ? $this->recordUserMessage($user, $data['session'], $data['message'], $mode, $model)
            : null;

        return response()->stream(function () use ($stream, $snapshot, $session): void {
            $this->emit(['type' => 'usage', ...$snapshot]);

            $reply = '';
            foreach ($this->tokens($stream) as $event) {
                $this->emit($event);

                if ($event['type'] !== 'token') {
                    break;
                }
                $reply .= $event['text'];
            }

            $stream->close();

            if ($session && $reply !== '') {
                $session->messages()->create(['role' => 'assistant', 'content' => $reply])
                    ->forceFill(['created_at' => now()])->save();
                $session->forceFill([
                    'messages_count' => $session->messages()->count(),
                    'last_activity_at' => now(),
                ])->save();
            }
        }, 200, [
            'Content-Type' => 'text/event-stream',
            'Cache-Control' => 'no-cache',
            'X-Accel-Buffering' => 'no',
        ]);
    }

    /**
     * @param  array<int, array{role: string, content: string}>  $history
     */
    private function openRouter(string $key, string $model, string $mode, string $effort, ?Agent $agent, string $message, array $history): Response
    {
        $messages = [['role' => 'system', 'content' => ChatPrompts::system($mode, $effort, $agent?->name, $agent?->instructions)]];
        foreach (array_slice($history, -self::HISTORY_LIMIT) as $turn) {
            $messages[] = ['role' => $turn['role'], 'content' => $turn['content']];
        }
        $messages[] = ['role' => 'user', 'content' => $message];

        return Http::withToken($key)
            ->withHeaders([
                'HTTP-Referer' => (string) config('app.url'),
                'X-Title' => (string) config('app.name'),
            ])
            ->connectTimeout(10)
            ->timeout(120)
            ->withOptions(['stream' => true])
            ->post(rtrim((string) config('services.openrouter.base_url'), '/').'/chat/completions', [
                'model' => $model,
                'stream' => true,
                'reasoning' => ['effort' => Effort::openRouter($effort)],
                'messages' => $messages,
            ]);
    }

    private function recordUserMessage(User $user, string $uuid, string $text, string $mode, string $model): ChatSession
    {
        $session = $user->chatSessions()->firstOrCreate(
            ['uuid' => $uuid],
            ['title' => Str::limit(trim(preg_replace('/\s+/', ' ', $text)), 80, '…'), 'mode' => $mode, 'model' => $model],
        );

        $message = $session->messages()->create(['role' => 'user', 'content' => $text]);
        $message->forceFill(['created_at' => now()])->save();

        $session->forceFill([
            'mode' => $mode,
            'model' => $model,
            'messages_count' => $session->messages()->count(),
            'last_activity_at' => now(),
        ])->save();

        return $session;
    }

    /**
     * Paket gratis memakai model gratis. Paket berbayar memakai model khusus paket itu
     * bila ada, lalu model per mode, lalu model berbayar bawaan.
     *
     * Pilihan /connect dan /model dari CLI dihormati hanya untuk paket berbayar
     * dan hanya bila modelnya ada di katalog config/ai.php.
     */
    private function model(User $user, string $mode, ?string $provider = null, ?string $requested = null): string
    {
        if (! $user->isFreePlan() && $provider && $requested) {
            $allowed = config("ai.providers.$provider", []);
            if (in_array($requested, $allowed, true)) {
                return $requested;
            }
        }

        if ($user->isFreePlan()) {
            return (string) (config('services.openrouter.free_model') ?: 'openrouter/free');
        }

        return config("usage.plans.{$user->plan}.model")
            ?: config("services.openrouter.models.$mode")
            ?: (string) config('services.openrouter.model');
    }

    /**
     * Baca aliran SSE OpenRouter dan hasilkan peristiwa `token`, lalu `done` atau `error`.
     *
     * @return \Generator<int, array<string, mixed>>
     */
    private function tokens(StreamInterface $stream): \Generator
    {
        $buffer = '';

        while (! $stream->eof()) {
            if (connection_aborted()) {
                return;
            }

            $buffer .= $stream->read(1024);

            while (($newline = strpos($buffer, "\n")) !== false) {
                $line = trim(substr($buffer, 0, $newline));
                $buffer = substr($buffer, $newline + 1);

                // Baris kosong dan komentar keep-alive (": OPENROUTER PROCESSING") diabaikan.
                if (! str_starts_with($line, 'data:')) {
                    continue;
                }

                $payload = trim(substr($line, 5));
                if ($payload === '[DONE]') {
                    yield ['type' => 'done'];

                    return;
                }

                $json = json_decode($payload, true);
                if (! is_array($json)) {
                    continue;
                }

                if (isset($json['error'])) {
                    yield ['type' => 'error', 'message' => $this->upstreamMessage(0, $json['error']['message'] ?? null)];

                    return;
                }

                $text = $json['choices'][0]['delta']['content'] ?? '';
                if (is_string($text) && $text !== '') {
                    yield ['type' => 'token', 'text' => $text];
                }
            }
        }

        yield ['type' => 'done'];
    }

    /**
     * @param  array<string, mixed>  $event
     */
    private function emit(array $event): void
    {
        echo 'data: '.json_encode($event, JSON_UNESCAPED_UNICODE)."\n\n";

        if (ob_get_level() > 0) {
            ob_flush();
        }
        flush();
    }

    private function blocked(UsageTracker $usage, User $user, string $blockedBy): JsonResponse
    {
        $snapshot = $usage->snapshot($user);

        return response()->json([
            'message' => $usage->blockedMessage($blockedBy, $snapshot),
            'blocked' => $blockedBy,
            ...$snapshot,
        ], $blockedBy === 'plan' ? 403 : 429);
    }

    private function upstreamMessage(int $status, ?string $detail): string
    {
        $base = match ($status) {
            401, 403 => 'Kunci OpenRouter di server ditolak.',
            402 => 'Saldo OpenRouter di server tidak cukup.',
            429 => 'Layanan AI sedang sibuk, coba lagi sebentar.',
            default => 'Layanan AI mengalami gangguan.',
        };

        return $detail ? "$base ($detail)" : $base;
    }
}

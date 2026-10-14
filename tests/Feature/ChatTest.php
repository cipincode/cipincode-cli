<?php

namespace Tests\Feature;

use App\Models\ApiKey;
use App\Models\UsageEvent;
use App\Models\User;
use App\Services\UsageTracker;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Http\Client\Request;
use Illuminate\Support\Facades\Http;
use Illuminate\Testing\TestResponse;
use Tests\TestCase;

class ChatTest extends TestCase
{
    use RefreshDatabase;

    private User $user;

    private string $token;

    protected function setUp(): void
    {
        parent::setUp();

        config([
            'services.openrouter.key' => 'sk-or-test',
            'services.openrouter.base_url' => 'https://openrouter.ai/api/v1',
            'services.openrouter.model' => 'openrouter/auto',
            'services.openrouter.models' => ['chat' => null, 'coding' => null, 'debugging' => null],
            'services.openrouter.free_model' => 'vendor/gratis:free',
            'usage.plans.go.limits.session' => 3,
            'usage.plans.go.limits.weekly' => 10,
            'usage.plans.free.limits.session' => 3,
        ]);

        $this->user = User::factory()->create(['plan' => 'go', 'plan_expires_at' => null]);
        [, $this->token] = ApiKey::generate($this->user, 'CLI');
    }

    private function sse(string ...$tokens): string
    {
        $body = ": OPENROUTER PROCESSING\n\n";
        foreach ($tokens as $token) {
            $body .= 'data: '.json_encode(['choices' => [['delta' => ['content' => $token]]]])."\n\n";
        }

        return $body."data: [DONE]\n\n";
    }

    private function fakeOpenRouter(?string $body = null, int $status = 200): void
    {
        $body ??= $this->sse('Halo', ' dunia');

        // Closure: setiap permintaan mendapat aliran baru, seperti OpenRouter sungguhan.
        Http::fake(['openrouter.ai/*' => fn () => Http::response($body, $status, ['Content-Type' => 'text/event-stream'])]);
    }

    /** @return array<int, array<string, mixed>> */
    private function events(string $content): array
    {
        $events = [];
        foreach (explode("\n", $content) as $line) {
            if (str_starts_with($line, 'data: ')) {
                $events[] = json_decode(substr($line, 6), true);
            }
        }

        return $events;
    }

    private function chat(array $payload = []): TestResponse
    {
        return $this->withToken($this->token)->postJson('/api/cli/chat', ['message' => 'halo', ...$payload]);
    }

    public function test_requires_an_api_key(): void
    {
        $this->postJson('/api/cli/chat', ['message' => 'halo'])->assertUnauthorized();
    }

    public function test_streams_tokens_then_done_with_usage_first(): void
    {
        $this->fakeOpenRouter();

        $events = $this->events($this->chat()->assertOk()->streamedContent());

        $this->assertSame('usage', $events[0]['type']);
        $this->assertSame(1, $events[0]['session']['used']);
        $this->assertSame(['token', 'token', 'done'], array_column(array_slice($events, 1), 'type'));
        $this->assertSame('Halo dunia', implode('', array_column(array_filter($events, fn ($e) => $e['type'] === 'token'), 'text')));
    }

    public function test_response_is_an_event_stream(): void
    {
        $this->fakeOpenRouter();

        $this->chat()->assertOk()->assertHeader('Content-Type', 'text/event-stream; charset=utf-8');
    }

    public function test_sends_system_prompt_history_and_message_to_openrouter(): void
    {
        $this->fakeOpenRouter();

        $this->chat([
            'message' => 'lanjutkan',
            'mode' => 'coding',
            'history' => [
                ['role' => 'user', 'content' => 'buat fungsi'],
                ['role' => 'assistant', 'content' => 'ini fungsinya'],
            ],
        ])->streamedContent();

        Http::assertSent(function (Request $request) {
            $messages = $request['messages'];

            return $request->url() === 'https://openrouter.ai/api/v1/chat/completions'
                && $request->hasHeader('Authorization', 'Bearer sk-or-test')
                && $request['stream'] === true
                && $request['model'] === 'openrouter/auto'
                && $messages[0]['role'] === 'system'
                && str_contains($messages[0]['content'], 'Mode coding')
                && array_column(array_slice($messages, 1), 'content') === ['buat fungsi', 'ini fungsinya', 'lanjutkan'];
        });
    }

    public function test_each_mode_gets_its_own_system_prompt(): void
    {
        $this->fakeOpenRouter();

        foreach (['chat' => 'Mode chat', 'debugging' => 'Mode debugging'] as $mode => $needle) {
            $this->chat(['mode' => $mode])->streamedContent();

            Http::assertSent(fn (Request $r) => str_contains($r['messages'][0]['content'], $needle));
        }
    }

    public function test_history_is_trimmed_to_the_latest_turns(): void
    {
        $this->fakeOpenRouter();
        $history = array_map(fn ($i) => ['role' => 'user', 'content' => "pesan $i"], range(1, 40));

        $this->chat(['history' => $history])->streamedContent();

        Http::assertSent(function (Request $r) {
            $messages = $r['messages'];

            return count($messages) === 1 + 20 + 1 && $messages[1]['content'] === 'pesan 21';
        });
    }

    public function test_per_mode_model_overrides_the_default(): void
    {
        config(['services.openrouter.models.coding' => 'vendor/coder-1']);
        $this->fakeOpenRouter();

        $this->chat(['mode' => 'coding'])->streamedContent();
        Http::assertSent(fn (Request $r) => $r['model'] === 'vendor/coder-1');

        $this->chat(['mode' => 'chat'])->streamedContent();
        Http::assertSent(fn (Request $r) => $r['model'] === 'openrouter/auto');
    }

    public function test_usage_is_charged_once_per_message(): void
    {
        $this->fakeOpenRouter();

        $this->chat()->streamedContent();

        $this->assertSame(1, UsageEvent::count());
        $this->assertSame(1, app(UsageTracker::class)->snapshot($this->user->refresh())['session']['used']);
    }

    public function test_missing_server_key_returns_503_without_charging(): void
    {
        config(['services.openrouter.key' => null]);
        Http::fake();

        $this->chat()->assertStatus(503)->assertJsonPath('message', fn ($m) => str_contains($m, 'OPENROUTER_API_KEY'));

        Http::assertNothingSent();
        $this->assertSame(0, UsageEvent::count());
    }

    public function test_exhausted_quota_returns_429_before_calling_openrouter(): void
    {
        app(UsageTracker::class)->consume($this->user, 3);
        Http::fake();

        $this->chat()
            ->assertStatus(429)
            ->assertJsonPath('blocked', 'session')
            ->assertJsonPath('session.exhausted', true);

        Http::assertNothingSent();
        $this->assertSame(1, UsageEvent::count());
    }

    public function test_upstream_failure_is_not_charged(): void
    {
        Http::fake(['openrouter.ai/*' => Http::response(['error' => ['message' => 'Insufficient credits', 'code' => 402]], 402)]);

        $this->chat()
            ->assertStatus(502)
            ->assertJsonPath('message', fn ($m) => str_contains($m, 'Saldo OpenRouter') && str_contains($m, 'Insufficient credits'));

        $this->assertSame(0, UsageEvent::count());
    }

    public function test_upstream_rate_limit_maps_to_503(): void
    {
        Http::fake(['openrouter.ai/*' => Http::response(['error' => ['message' => 'slow down']], 429)]);

        $this->chat()->assertStatus(503);
        $this->assertSame(0, UsageEvent::count());
    }

    public function test_upstream_error_in_the_stream_becomes_an_error_event(): void
    {
        $this->fakeOpenRouter(
            'data: '.json_encode(['choices' => [['delta' => ['content' => 'Se']]]])."\n\n"
            .'data: '.json_encode(['error' => ['message' => 'Provider disconnected']])."\n\n",
        );

        $events = $this->events($this->chat()->streamedContent());

        $this->assertSame(['usage', 'token', 'error'], array_column($events, 'type'));
        $this->assertStringContainsString('Provider disconnected', $events[2]['message']);
    }

    public function test_stream_without_done_marker_still_ends_with_done(): void
    {
        $this->fakeOpenRouter('data: '.json_encode(['choices' => [['delta' => ['content' => 'Halo']]]])."\n\n");

        $events = $this->events($this->chat()->streamedContent());

        $this->assertSame('done', end($events)['type']);
    }

    public function test_validates_the_request(): void
    {
        $this->chat(['message' => ''])->assertStatus(422)->assertJsonValidationErrors('message');
        $this->chat(['mode' => 'hack'])->assertStatus(422)->assertJsonValidationErrors('mode');
        $this->chat(['history' => [['role' => 'system', 'content' => 'x']]])->assertStatus(422);
        $this->chat(['message' => str_repeat('a', 8001)])->assertStatus(422);
    }

    public function test_free_plan_uses_the_free_model_for_every_mode(): void
    {
        config(['services.openrouter.models.coding' => 'vendor/coder-1']);
        $free = User::factory()->create();
        [, $token] = ApiKey::generate($free, 'CLI');
        $this->fakeOpenRouter();

        foreach (['chat', 'coding', 'debugging'] as $mode) {
            $this->withToken($token)->postJson('/api/cli/chat', ['message' => 'halo', 'mode' => $mode])->streamedContent();
        }

        Http::assertSentCount(3);
        Http::assertSent(fn (Request $r) => $r['model'] === 'vendor/gratis:free');
        Http::assertNotSent(fn (Request $r) => $r['model'] !== 'vendor/gratis:free');
    }

    public function test_free_model_falls_back_to_openrouter_free(): void
    {
        config(['services.openrouter.free_model' => null]);
        [, $token] = ApiKey::generate(User::factory()->create(), 'CLI');
        $this->fakeOpenRouter();

        $this->withToken($token)->postJson('/api/cli/chat', ['message' => 'halo'])->streamedContent();

        Http::assertSent(fn (Request $r) => $r['model'] === 'openrouter/free');
    }

    public function test_expired_free_trial_is_blocked_before_calling_openrouter(): void
    {
        $free = User::factory()->create(['plan_expires_at' => now()->subMinute()]);
        [, $token] = ApiKey::generate($free, 'CLI');
        Http::fake();

        $this->withToken($token)->postJson('/api/cli/chat', ['message' => 'halo'])
            ->assertForbidden()
            ->assertJsonPath('blocked', 'plan')
            ->assertJsonPath('plan.expired', true)
            ->assertJsonPath('message', fn ($m) => str_contains($m, 'Masa percobaan gratis'));

        Http::assertNothingSent();
        $this->assertSame(0, UsageEvent::count());
    }

    public function test_default_effort_is_medium_with_no_extra_directive(): void
    {
        $this->fakeOpenRouter();

        $this->chat()->streamedContent();

        Http::assertSent(function (Request $r) {
            return $r['reasoning'] === ['effort' => 'medium']
                && ! str_contains($r['messages'][0]['content'], 'Tingkat effort');
        });
    }

    public function test_each_effort_maps_to_an_openrouter_level_and_a_prompt_directive(): void
    {
        config(['usage.plans.go.limits.session' => 20]); // lima panggilan beruntun
        $this->fakeOpenRouter();
        $expected = [
            'low' => ['low', 'effort rendah'],
            'medium' => ['medium', null],
            'high' => ['high', 'effort tinggi'],
            'xhigh' => ['xhigh', 'effort sangat tinggi'],
            'ultra' => ['max', 'effort maksimum'],
        ];

        foreach ($expected as $effort => [$level, $phrase]) {
            $this->chat(['effort' => $effort])->streamedContent();

            Http::assertSent(function (Request $r) use ($level, $phrase) {
                $system = $r['messages'][0]['content'];

                return $r['reasoning']['effort'] === $level
                    && ($phrase === null ? ! str_contains($system, 'Tingkat effort') : str_contains($system, $phrase));
            });
        }
    }

    public function test_effort_works_on_the_free_plan_too(): void
    {
        [, $token] = ApiKey::generate(User::factory()->create(), 'CLI');
        $this->fakeOpenRouter();

        $this->withToken($token)->postJson('/api/cli/chat', ['message' => 'halo', 'effort' => 'ultra'])->streamedContent();

        Http::assertSent(fn (Request $r) => $r['model'] === 'vendor/gratis:free'
            && $r['reasoning']['effort'] === 'max'
            && str_contains($r['messages'][0]['content'], 'effort maksimum'));
    }

    public function test_effort_combines_with_the_mode_prompt(): void
    {
        $this->fakeOpenRouter();

        $this->chat(['mode' => 'debugging', 'effort' => 'high'])->streamedContent();

        Http::assertSent(fn (Request $r) => str_contains($r['messages'][0]['content'], 'Mode debugging')
            && str_contains($r['messages'][0]['content'], 'effort tinggi'));
    }

    public function test_unknown_effort_is_rejected(): void
    {
        $this->chat(['effort' => 'mega'])->assertStatus(422)->assertJsonValidationErrors('effort');
        $this->assertSame(0, UsageEvent::count());
    }
}

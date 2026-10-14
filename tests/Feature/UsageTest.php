<?php

namespace Tests\Feature;

use App\Models\ApiKey;
use App\Models\UsageEvent;
use App\Models\User;
use App\Services\UsageTracker;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class UsageTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();

        config([
            'usage.plans.free.limits.session' => 4,
            'usage.plans.go.limits.session' => 4,
            'usage.plans.go.limits.weekly' => 10,
            'usage.session_hours' => 5,
        ]);
        $this->travelTo(now()->startOfWeek()->addDays(2)->setTime(10, 0)); // Rabu 10.00
    }

    private function paidUser(): User
    {
        return User::factory()->create(['plan' => 'go', 'plan_expires_at' => null]);
    }

    private function tracker(): UsageTracker
    {
        return app(UsageTracker::class);
    }

    public function test_fresh_user_has_zero_usage_and_no_session_window(): void
    {
        $snapshot = $this->tracker()->snapshot($this->paidUser());

        $this->assertSame(0, $snapshot['session']['used']);
        $this->assertSame(4, $snapshot['session']['limit']);
        $this->assertNull($snapshot['session']['resets_at']);
        $this->assertSame(0, $snapshot['session']['percent']);
        $this->assertFalse($snapshot['session']['exhausted']);
        $this->assertSame(0, $snapshot['weekly']['used']);
        $this->assertNotNull($snapshot['weekly']['resets_at']);
    }

    public function test_first_use_starts_a_five_hour_window(): void
    {
        $user = $this->paidUser();

        [$allowed, $snapshot] = $this->tracker()->consume($user);

        $this->assertTrue($allowed);
        $this->assertSame(1, $snapshot['session']['used']);
        $this->assertSame(25, $snapshot['session']['percent']);
        $this->assertSame(now()->addHours(5)->toIso8601String(), $snapshot['session']['resets_at']);
        $this->assertSame(1, $snapshot['weekly']['used']);
    }

    public function test_window_end_is_fixed_by_the_first_message(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user);
        $expectedReset = now()->addHours(5)->toIso8601String();

        $this->travel(3)->hours();
        [, $snapshot] = $this->tracker()->consume($user);

        $this->assertSame(2, $snapshot['session']['used']);
        $this->assertSame($expectedReset, $snapshot['session']['resets_at']);
    }

    public function test_session_resets_automatically_after_five_hours_but_weekly_keeps_counting(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user, 3);

        $this->travel(5)->hours();
        $snapshot = $this->tracker()->snapshot($user->refresh());

        $this->assertSame(0, $snapshot['session']['used']);
        $this->assertNull($snapshot['session']['resets_at']);
        $this->assertSame(3, $snapshot['weekly']['used']);

        [$allowed, $after] = $this->tracker()->consume($user->refresh());
        $this->assertTrue($allowed);
        $this->assertSame(1, $after['session']['used']);
        $this->assertSame(4, $after['weekly']['used']);
    }

    public function test_session_limit_blocks_and_does_not_record_more(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user, 4);

        [$allowed, $snapshot, $blockedBy] = $this->tracker()->consume($user->refresh());

        $this->assertFalse($allowed);
        $this->assertSame('session', $blockedBy);
        $this->assertTrue($snapshot['session']['exhausted']);
        $this->assertSame(100, $snapshot['session']['percent']);
        $this->assertSame(1, UsageEvent::count());
    }

    public function test_weekly_limit_blocks_even_when_session_is_fresh(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user, 4);
        $this->travel(5)->hours();
        $this->tracker()->consume($user->refresh(), 4);
        $this->travel(5)->hours();
        $this->tracker()->consume($user->refresh(), 2);
        $this->travel(5)->hours();

        [$allowed, $snapshot, $blockedBy] = $this->tracker()->consume($user->refresh());

        $this->assertFalse($allowed);
        $this->assertSame('weekly', $blockedBy);
        $this->assertSame(0, $snapshot['session']['used']);
        $this->assertTrue($snapshot['weekly']['exhausted']);
    }

    public function test_weekly_usage_resets_on_monday(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user, 3);
        $weeklyReset = $this->tracker()->snapshot($user->refresh())['weekly']['resets_at'];

        $this->assertSame(now()->startOfWeek()->addWeek()->toIso8601String(), $weeklyReset);

        $this->travelTo(now()->startOfWeek()->addWeek()->addMinute());
        $this->assertSame(0, $this->tracker()->snapshot($user->refresh())['weekly']['used']);
    }

    public function test_usage_is_per_user(): void
    {
        $a = $this->paidUser();
        $b = $this->paidUser();
        $this->tracker()->consume($a, 3);

        $this->assertSame(0, $this->tracker()->snapshot($b)['session']['used']);
    }

    public function test_api_requires_a_key(): void
    {
        $this->getJson('/api/cli/usage')->assertUnauthorized();
        $this->postJson('/api/cli/usage')->assertUnauthorized();
    }

    public function test_api_reports_and_consumes_usage(): void
    {
        [, $plain] = ApiKey::generate($this->paidUser(), 'CLI');

        $this->withToken($plain)->getJson('/api/cli/usage')
            ->assertOk()
            ->assertJsonPath('session.used', 0)
            ->assertJsonPath('weekly.limit', 10);

        $this->withToken($plain)->postJson('/api/cli/usage')
            ->assertOk()
            ->assertJsonPath('session.used', 1);

        $this->withToken($plain)->postJson('/api/cli/usage', ['units' => 2])
            ->assertOk()
            ->assertJsonPath('session.used', 3);
    }

    public function test_api_returns_429_with_snapshot_when_exhausted(): void
    {
        [, $plain] = ApiKey::generate($this->paidUser(), 'CLI');
        $this->withToken($plain)->postJson('/api/cli/usage', ['units' => 4])->assertOk();

        $this->withToken($plain)->postJson('/api/cli/usage')
            ->assertStatus(429)
            ->assertJsonPath('blocked', 'session')
            ->assertJsonPath('session.exhausted', true)
            ->assertJsonPath('message', 'Batas pemakaian sesi 5 jam tercapai.');
    }

    public function test_api_validates_units(): void
    {
        [, $plain] = ApiKey::generate($this->paidUser(), 'CLI');

        $this->withToken($plain)->postJson('/api/cli/usage', ['units' => 0])->assertStatus(422);
        $this->withToken($plain)->postJson('/api/cli/usage', ['units' => 101])->assertStatus(422);
    }

    public function test_dashboard_page_shows_the_snapshot(): void
    {
        $user = $this->paidUser();
        $this->tracker()->consume($user, 2);

        $this->actingAs($user->refresh())->get(route('usage'))
            ->assertInertia(fn ($page) => $page
                ->component('usage')
                ->where('usage.session.used', 2)
                ->where('usage.session.percent', 50)
                ->where('usage.weekly.used', 2)
                ->where('sessionHours', 5));
    }

    public function test_guests_cannot_view_the_usage_page(): void
    {
        $this->get(route('usage'))->assertRedirect(route('login'));
    }
}

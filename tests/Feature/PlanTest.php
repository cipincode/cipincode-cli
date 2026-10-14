<?php

namespace Tests\Feature;

use App\Models\ApiKey;
use App\Models\User;
use App\Services\UsageTracker;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class PlanTest extends TestCase
{
    use RefreshDatabase;

    protected function setUp(): void
    {
        parent::setUp();

        config([
            'usage.plans.free.limits.session' => 3,
            'usage.plans.go.limits.session' => 5,
            'usage.plans.go.limits.weekly' => 20,
        ]);
    }

    private function tracker(): UsageTracker
    {
        return app(UsageTracker::class);
    }

    public function test_new_users_start_on_the_free_plan_with_a_two_week_trial(): void
    {
        $user = User::factory()->create();

        $this->assertSame('free', $user->plan);
        $this->assertTrue($user->isFreePlan());
        $this->assertSame(now()->addDays(14)->toDateString(), $user->plan_expires_at->toDateString());
    }

    public function test_paid_users_created_directly_have_no_expiry(): void
    {
        $user = User::factory()->create(['plan' => 'go']);

        $this->assertFalse($user->isFreePlan());
        $this->assertNull($user->plan_expires_at);
    }

    public function test_free_plan_only_has_the_five_hour_limit(): void
    {
        $snapshot = $this->tracker()->snapshot(User::factory()->create());

        $this->assertSame(3, $snapshot['session']['limit']);
        $this->assertNull($snapshot['weekly']);
        $this->assertSame('free', $snapshot['plan']['name']);
        $this->assertSame('Gratis', $snapshot['plan']['label']);
        $this->assertFalse($snapshot['plan']['expired']);
        $this->assertNotNull($snapshot['plan']['expires_at']);
    }

    public function test_paid_plan_has_session_and_weekly_limits(): void
    {
        $snapshot = $this->tracker()->snapshot(User::factory()->create(['plan' => 'go']));

        $this->assertSame(5, $snapshot['session']['limit']);
        $this->assertSame(20, $snapshot['weekly']['limit']);
        $this->assertSame('Go', $snapshot['plan']['label']);
        $this->assertNull($snapshot['plan']['expires_at']);
    }

    public function test_free_plan_is_never_blocked_by_a_weekly_limit(): void
    {
        $user = User::factory()->create();

        foreach (range(1, 3) as $i) {
            [$allowed] = $this->tracker()->consume($user->refresh());
            $this->assertTrue($allowed);
        }
        $this->travel(5)->hours();
        // Sesi baru; satu-satunya batas yang berlaku tetap batas 5 jam.
        [$allowed, $snapshot] = $this->tracker()->consume($user->refresh());

        $this->assertTrue($allowed);
        $this->assertNull($snapshot['weekly']);
    }

    public function test_free_plan_session_limit_still_applies_and_resets(): void
    {
        $user = User::factory()->create();
        $this->tracker()->consume($user, 3);

        [$allowed, , $blockedBy] = $this->tracker()->consume($user->refresh());
        $this->assertFalse($allowed);
        $this->assertSame('session', $blockedBy);

        $this->travel(5)->hours();
        [$allowed] = $this->tracker()->consume($user->refresh());
        $this->assertTrue($allowed);
    }

    public function test_free_trial_expires_after_two_weeks_and_blocks_everything(): void
    {
        $user = User::factory()->create();
        [$allowed] = $this->tracker()->consume($user);
        $this->assertTrue($allowed);

        $this->travel(14)->days();
        $snapshot = $this->tracker()->snapshot($user->refresh());
        $this->assertTrue($snapshot['plan']['expired']);

        [$allowed, , $blockedBy] = $this->tracker()->consume($user->refresh());
        $this->assertFalse($allowed);
        $this->assertSame('plan', $blockedBy);
    }

    public function test_trial_is_still_valid_one_minute_before_it_ends(): void
    {
        $user = User::factory()->create();

        $this->travel(14)->days();
        $this->travel(-1)->minutes();

        $this->assertFalse($this->tracker()->snapshot($user->refresh())['plan']['expired']);
    }

    public function test_plan_expiry_takes_priority_over_other_limits(): void
    {
        $user = User::factory()->create(['plan_expires_at' => now()->subSecond()]);

        $this->assertSame('plan', $this->tracker()->blockedBy($this->tracker()->snapshot($user)));
    }

    public function test_blocked_messages_depend_on_the_plan(): void
    {
        $free = $this->tracker()->snapshot(User::factory()->create());
        $paid = $this->tracker()->snapshot(User::factory()->create(['plan' => 'go']));

        $this->assertStringContainsString('percobaan gratis 14 hari', $this->tracker()->blockedMessage('plan', $free));
        $this->assertStringContainsString('Perpanjang', $this->tracker()->blockedMessage('plan', $paid));
        $this->assertStringContainsString('sesi 5 jam', $this->tracker()->blockedMessage('session'));
        $this->assertStringContainsString('mingguan', $this->tracker()->blockedMessage('weekly'));
    }

    public function test_expired_paid_plan_is_blocked_too(): void
    {
        $user = User::factory()->create(['plan' => 'go', 'plan_expires_at' => now()->subDay()]);

        [$allowed, , $blockedBy] = $this->tracker()->consume($user);

        $this->assertFalse($allowed);
        $this->assertSame('plan', $blockedBy);
    }

    public function test_usage_api_exposes_the_plan_and_null_weekly_for_free(): void
    {
        [, $token] = ApiKey::generate(User::factory()->create(), 'CLI');

        $this->withToken($token)->getJson('/api/cli/usage')
            ->assertOk()
            ->assertJsonPath('plan.name', 'free')
            ->assertJsonPath('weekly', null)
            ->assertJsonPath('session.limit', 3);
    }

    public function test_manual_usage_endpoint_reports_an_expired_plan_as_the_blocker(): void
    {
        [, $token] = ApiKey::generate(User::factory()->create(['plan_expires_at' => now()->subDay()]), 'CLI');

        $this->withToken($token)->postJson('/api/cli/usage')
            ->assertStatus(429) // endpoint manual tetap 429; blocked menunjukkan penyebabnya
            ->assertJsonPath('blocked', 'plan');
    }

    public function test_dashboard_usage_page_receives_the_plan(): void
    {
        $user = User::factory()->create();

        $this->actingAs($user)->get(route('usage'))
            ->assertInertia(fn ($page) => $page
                ->component('usage')
                ->where('usage.plan.name', 'free')
                ->where('usage.weekly', null));
    }

    public function test_plan_set_command_switches_plans(): void
    {
        $user = User::factory()->create();

        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'go'])->assertSuccessful();
        $user->refresh();
        $this->assertSame('go', $user->plan);
        $this->assertNull($user->plan_expires_at);

        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'free'])->assertSuccessful();
        $user->refresh();
        $this->assertSame('free', $user->plan);
        $this->assertSame(now()->addDays(14)->toDateString(), $user->plan_expires_at->toDateString());
    }

    public function test_plan_set_command_accepts_custom_days_and_rejects_bad_input(): void
    {
        $user = User::factory()->create();

        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'go', '--days' => 30])->assertSuccessful();
        $this->assertSame(now()->addDays(30)->toDateString(), $user->refresh()->plan_expires_at->toDateString());

        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'gold'])->assertFailed();
        $this->artisan('plan:set', ['email' => 'tidak@ada.test', 'plan' => 'go'])->assertFailed();
        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'go', '--days' => 'abc'])->assertFailed();
        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'go', '--days' => '0'])->assertFailed();
    }
}

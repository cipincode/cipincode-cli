<?php

namespace Tests\Feature;

use App\Models\User;
use App\Services\UsageTracker;
use App\Support\Plans;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class PlanCatalogTest extends TestCase
{
    use RefreshDatabase;

    public function test_catalog_lists_free_go_pro_max_in_ascending_order(): void
    {
        $catalog = Plans::catalog();

        $this->assertSame(['free', 'go', 'pro', 'max'], array_column($catalog, 'key'));
        $this->assertSame(['Gratis', 'Go', 'Pro', 'Max'], array_column($catalog, 'label'));
    }

    public function test_prices_are_in_rupiah_and_formatted_with_dots(): void
    {
        $byKey = array_column(Plans::catalog(), null, 'key');

        $this->assertSame('Gratis', $byKey['free']['price_label']);
        $this->assertNull($byKey['free']['period']);
        $this->assertSame('Rp 49.000', $byKey['go']['price_label']);
        $this->assertSame('Rp 149.000', $byKey['pro']['price_label']);
        $this->assertSame('Rp 499.000', $byKey['max']['price_label']);
        $this->assertSame('per bulan', $byKey['max']['period']);
        $this->assertSame('Rp 1.250.000', Plans::rupiah(1250000));
    }

    public function test_paid_tiers_get_progressively_larger_limits_and_prices(): void
    {
        $paid = array_values(array_filter(Plans::catalog(), fn ($p) => $p['paid']));

        $this->assertCount(3, $paid);
        for ($i = 1; $i < 3; $i++) {
            $this->assertGreaterThan($paid[$i - 1]['price'], $paid[$i]['price']);
            $this->assertGreaterThan($paid[$i - 1]['session_limit'], $paid[$i]['session_limit']);
            $this->assertGreaterThan($paid[$i - 1]['weekly_limit'], $paid[$i]['weekly_limit']);
        }
    }

    public function test_free_plan_has_no_weekly_limit_in_the_catalog(): void
    {
        $free = Plans::catalog()[0];

        $this->assertNull($free['weekly_limit']);
        $this->assertContains('Tanpa batas mingguan', $free['features']);
        $this->assertContains('Masa percobaan 14 hari', $free['features']);
    }

    public function test_features_are_generated_from_the_real_limits(): void
    {
        config(['usage.plans.pro.limits.session' => 123, 'usage.plans.pro.limits.weekly' => 4567]);

        $pro = array_column(Plans::catalog(), null, 'key')['pro'];

        $this->assertContains('123 pesan per 5 jam', $pro['features']);
        $this->assertContains('4.567 pesan per minggu', $pro['features']);
    }

    public function test_each_tier_enforces_its_own_limits(): void
    {
        config([
            'usage.plans.go.limits' => ['session' => 2, 'weekly' => 10],
            'usage.plans.pro.limits' => ['session' => 4, 'weekly' => 20],
            'usage.plans.max.limits' => ['session' => 6, 'weekly' => 30],
        ]);
        $tracker = app(UsageTracker::class);

        foreach (['go' => 2, 'pro' => 4, 'max' => 6] as $plan => $limit) {
            $user = User::factory()->create(['plan' => $plan]);
            $snapshot = $tracker->snapshot($user);

            $this->assertSame($limit, $snapshot['session']['limit'], $plan);
            $this->assertSame(ucfirst($plan), $snapshot['plan']['label']);

            $tracker->consume($user, $limit);
            [$allowed, , $blockedBy] = $tracker->consume($user->refresh());
            $this->assertFalse($allowed, $plan);
            $this->assertSame('session', $blockedBy);
        }
    }

    public function test_plan_set_accepts_every_plan_and_rejects_unknown_ones(): void
    {
        $user = User::factory()->create();

        foreach (['go', 'pro', 'max', 'free'] as $plan) {
            $this->artisan('plan:set', ['email' => $user->email, 'plan' => $plan])->assertSuccessful();
            $this->assertSame($plan, $user->refresh()->plan);
        }

        $this->artisan('plan:set', ['email' => $user->email, 'plan' => 'paid'])->assertFailed();
        $this->assertSame('free', $user->refresh()->plan);
    }

    public function test_usage_page_receives_the_catalog_for_the_upgrade_section(): void
    {
        $this->actingAs(User::factory()->create())->get(route('usage'))
            ->assertInertia(fn ($page) => $page
                ->component('usage')
                ->has('plans', 4)
                ->where('plans.1.key', 'go')
                ->where('plans.1.price_label', 'Rp 49.000'));
    }

    public function test_billing_page_shows_the_catalog_current_plan_and_selection(): void
    {
        $user = User::factory()->create(['plan' => 'go']);

        $this->actingAs($user)->get(route('billing', ['plan' => 'pro']))
            ->assertInertia(fn ($page) => $page
                ->component('billing')
                ->has('plans', 4)
                ->where('current.name', 'go')
                ->where('selected', 'pro'));

        $this->actingAs($user)->get(route('billing', ['plan' => 'hack']))
            ->assertInertia(fn ($page) => $page->where('selected', null));
    }

    public function test_guests_cannot_open_billing(): void
    {
        $this->get(route('billing'))->assertRedirect(route('login'));
    }

    public function test_plan_specific_model_overrides_the_default_paid_model(): void
    {
        config([
            'usage.plans.max.model' => 'vendor/max-model',
            'services.openrouter.key' => 'sk-test',
            'services.openrouter.model' => 'openrouter/auto',
        ]);
        \Illuminate\Support\Facades\Http::fake(['openrouter.ai/*' => fn () => \Illuminate\Support\Facades\Http::response("data: [DONE]\n\n")]);

        foreach (['go' => 'openrouter/auto', 'max' => 'vendor/max-model'] as $plan => $model) {
            [, $token] = \App\Models\ApiKey::generate(User::factory()->create(['plan' => $plan]), 'CLI');
            $this->withToken($token)->postJson('/api/cli/chat', ['message' => 'halo'])->streamedContent();

            \Illuminate\Support\Facades\Http::assertSent(fn ($r) => $r['model'] === $model);
        }
    }
}

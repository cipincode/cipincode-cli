<?php

namespace Tests\Feature;

use App\Models\ApiKey;
use App\Models\CliAuthCode;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Str;
use Tests\TestCase;

class CliAuthTest extends TestCase
{
    use RefreshDatabase;

    private function verifier(): string
    {
        return Str::random(64);
    }

    /** Setujui di browser, kembalikan kode yang tampil. */
    private function approve(User $user, string $verifier): string
    {
        $this->actingAs($user)->post(route('cli.authorize.approve'), [
            'challenge' => CliAuthCode::challengeFor($verifier),
        ]);

        return session('cliCode');
    }

    public function test_guests_are_sent_to_login_and_return_afterwards(): void
    {
        $url = route('cli.authorize', ['challenge' => CliAuthCode::challengeFor($this->verifier())]);

        $this->get($url)->assertRedirect(route('login'));
        $this->assertSame($url, session('url.intended'));
    }

    public function test_challenge_must_be_well_formed(): void
    {
        $this->actingAs(User::factory()->create())
            ->get(route('cli.authorize', ['challenge' => 'pendek']))
            ->assertSessionHasErrors('challenge');
    }

    public function test_authorize_page_shows_the_account(): void
    {
        $user = User::factory()->create(['name' => 'Budi']);

        $this->actingAs($user)
            ->get(route('cli.authorize', ['challenge' => CliAuthCode::challengeFor($this->verifier()), 'device' => 'laptop']))
            ->assertInertia(fn ($page) => $page
                ->component('auth/cli-authorize')
                ->where('user.name', 'Budi')
                ->where('device', 'laptop')
                ->where('code', null));
    }

    public function test_approving_issues_a_readable_code(): void
    {
        $code = $this->approve(User::factory()->create(), $this->verifier());

        $this->assertMatchesRegularExpression('/^[A-Z2-9]{4}-[A-Z2-9]{4}$/', $code);
        $this->assertDatabaseCount('cli_auth_codes', 1);
        $this->assertDatabaseMissing('cli_auth_codes', ['code_hash' => $code]);
    }

    public function test_guests_cannot_approve(): void
    {
        $this->post(route('cli.authorize.approve'), [
            'challenge' => CliAuthCode::challengeFor($this->verifier()),
        ])->assertRedirect(route('login'));

        $this->assertDatabaseCount('cli_auth_codes', 0);
    }

    public function test_full_flow_gives_the_cli_a_working_key(): void
    {
        $user = User::factory()->create(['name' => 'Budi', 'email' => 'budi@example.com']);
        $verifier = $this->verifier();
        $code = $this->approve($user, $verifier);

        $response = $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier, 'device' => 'laptop-saya'])
            ->assertOk()
            ->assertJsonPath('user.email', 'budi@example.com');

        $key = $response->json('key');
        $this->assertStringStartsWith('cpc_', $key);
        $this->assertSame('CLI · laptop-saya', $user->apiKeys()->firstOrFail()->name);
        $this->withToken($key)->getJson('/api/cli/me')->assertOk();
    }

    public function test_code_is_accepted_without_dash_and_in_lowercase(): void
    {
        $verifier = $this->verifier();
        $code = $this->approve(User::factory()->create(), $verifier);

        $this->postJson('/api/cli/exchange', [
            'code' => strtolower(str_replace('-', '', $code)),
            'verifier' => $verifier,
        ])->assertOk();
    }

    public function test_code_can_only_be_used_once(): void
    {
        $verifier = $this->verifier();
        $code = $this->approve(User::factory()->create(), $verifier);

        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier])->assertOk();
        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier])->assertStatus(422);
    }

    public function test_wrong_verifier_is_rejected_and_code_stays_usable(): void
    {
        $verifier = $this->verifier();
        $code = $this->approve(User::factory()->create(), $verifier);

        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $this->verifier()])->assertStatus(422);
        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier])->assertOk();
    }

    public function test_expired_and_unknown_codes_are_rejected(): void
    {
        $verifier = $this->verifier();
        $code = $this->approve(User::factory()->create(), $verifier);

        $this->postJson('/api/cli/exchange', ['code' => 'ZZZZ-ZZZZ', 'verifier' => $verifier])->assertStatus(422);

        $this->travel(CliAuthCode::TTL_MINUTES + 1)->minutes();
        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier])->assertStatus(422);
    }

    public function test_key_limit_blocks_the_exchange(): void
    {
        $user = User::factory()->create();
        foreach (range(1, ApiKey::MAX_PER_USER) as $i) {
            ApiKey::generate($user, "Key $i");
        }
        $verifier = $this->verifier();
        $code = $this->approve($user, $verifier);

        $this->postJson('/api/cli/exchange', ['code' => $code, 'verifier' => $verifier])
            ->assertStatus(422)
            ->assertJsonPath('message', fn ($m) => str_contains($m, 'Batas'));

        $this->assertSame(ApiKey::MAX_PER_USER, $user->apiKeys()->count());
    }

    public function test_exchange_requires_both_fields(): void
    {
        $this->postJson('/api/cli/exchange', [])->assertStatus(422)->assertJsonValidationErrors(['code', 'verifier']);
    }
}

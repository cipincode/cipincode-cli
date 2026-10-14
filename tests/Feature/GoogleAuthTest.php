<?php

namespace Tests\Feature;

use App\Models\CliAuthCode;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Str;
use Laravel\Socialite\Contracts\Provider;
use Laravel\Socialite\Facades\Socialite;
use Laravel\Socialite\Two\User as SocialiteUser;
use Mockery;
use Tests\TestCase;

class GoogleAuthTest extends TestCase
{
    use RefreshDatabase;

    private function enable(): void
    {
        config([
            'services.google.client_id' => 'id',
            'services.google.client_secret' => 'secret',
        ]);
    }

    private function googleReturns(string $email, bool $verified = true, string $name = 'Budi Google'): void
    {
        $google = (new SocialiteUser)->map(['id' => '1', 'name' => $name, 'email' => $email]);
        $google->user = ['email_verified' => $verified];

        $provider = Mockery::mock(Provider::class);
        $provider->shouldReceive('user')->andReturn($google);
        Socialite::shouldReceive('driver')->with('google')->andReturn($provider);
    }

    public function test_routes_are_hidden_until_google_is_configured(): void
    {
        $this->get(route('google.redirect'))->assertNotFound();
        $this->get(route('google.callback'))->assertNotFound();
    }

    public function test_login_page_reports_whether_google_is_enabled(): void
    {
        $this->get(route('login'))->assertInertia(fn ($page) => $page->where('googleEnabled', false));

        $this->enable();
        $this->get(route('login'))->assertInertia(fn ($page) => $page->where('googleEnabled', true));
    }

    public function test_redirect_goes_to_google(): void
    {
        $this->enable();

        $this->get(route('google.redirect'))->assertRedirectContains('accounts.google.com');
    }

    public function test_new_google_user_is_created_and_logged_in(): void
    {
        $this->enable();
        $this->googleReturns('baru@example.com');

        $this->get(route('google.callback'))->assertRedirect(route('dashboard'));

        $user = User::where('email', 'baru@example.com')->firstOrFail();
        $this->assertSame('Budi Google', $user->name);
        $this->assertNotNull($user->email_verified_at);
        $this->assertAuthenticatedAs($user);
    }

    public function test_existing_user_is_logged_in_without_duplicates(): void
    {
        $this->enable();
        $user = User::factory()->create(['email' => 'ada@example.com']);
        $this->googleReturns('ada@example.com');

        $this->get(route('google.callback'))->assertRedirect(route('dashboard'));

        $this->assertAuthenticatedAs($user);
        $this->assertSame(1, User::where('email', 'ada@example.com')->count());
    }

    public function test_unverified_google_email_is_rejected(): void
    {
        $this->enable();
        $this->googleReturns('palsu@example.com', verified: false);

        $this->get(route('google.callback'))->assertRedirect(route('login'))->assertSessionHasErrors('email');

        $this->assertGuest();
        $this->assertDatabaseMissing('users', ['email' => 'palsu@example.com']);
    }

    public function test_accounts_with_two_factor_must_use_email_login(): void
    {
        $this->enable();
        User::factory()->create(['email' => 'aman@example.com', 'two_factor_confirmed_at' => now()]);
        $this->googleReturns('aman@example.com');

        $this->get(route('google.callback'))->assertRedirect(route('login'))->assertSessionHasErrors('email');

        $this->assertGuest();
    }

    public function test_cli_authorize_sends_guests_to_google_when_requested(): void
    {
        $this->enable();
        $url = route('cli.authorize', ['challenge' => CliAuthCode::challengeFor(Str::random(64)), 'method' => 'google']);

        $this->get($url)->assertRedirect(route('google.redirect'));
        $this->assertSame($url, session('url.intended'));
    }

    public function test_cli_authorize_falls_back_to_login_when_google_is_off(): void
    {
        $url = route('cli.authorize', ['challenge' => CliAuthCode::challengeFor(Str::random(64)), 'method' => 'google']);

        $this->get($url)->assertRedirect(route('login'));
    }
}

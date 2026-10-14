<?php

namespace Tests\Feature;

use App\Models\ApiKey;
use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class ApiKeyTest extends TestCase
{
    use RefreshDatabase;

    public function test_user_can_create_a_key_and_sees_it_once(): void
    {
        $user = User::factory()->create();

        $response = $this->actingAs($user)->post(route('api-keys.store'), ['name' => 'Laptop']);

        $response->assertRedirect(route('api-keys'));
        $plain = session('newKey');
        $this->assertStringStartsWith('cpc_', $plain);

        $key = $user->apiKeys()->firstOrFail();
        $this->assertSame('Laptop', $key->name);
        $this->assertSame(ApiKey::hash($plain), $key->key_hash);
        $this->assertNotSame($plain, $key->key_hash);
        $this->assertSame(substr($plain, 0, 8), $key->prefix);
    }

    public function test_name_is_required(): void
    {
        $this->actingAs(User::factory()->create())
            ->post(route('api-keys.store'), ['name' => ''])
            ->assertSessionHasErrors('name');
    }

    public function test_key_limit_is_enforced(): void
    {
        $user = User::factory()->create();
        foreach (range(1, ApiKey::MAX_PER_USER) as $i) {
            ApiKey::generate($user, "Key $i");
        }

        $this->actingAs($user)
            ->post(route('api-keys.store'), ['name' => 'Terlalu banyak'])
            ->assertSessionHasErrors('name');

        $this->assertSame(ApiKey::MAX_PER_USER, $user->apiKeys()->count());
    }

    public function test_user_can_revoke_own_key(): void
    {
        $user = User::factory()->create();
        [$key] = ApiKey::generate($user, 'Lama');

        $this->actingAs($user)->delete(route('api-keys.destroy', $key))->assertRedirect(route('api-keys'));

        $this->assertDatabaseMissing('api_keys', ['id' => $key->id]);
    }

    public function test_user_cannot_revoke_someone_elses_key(): void
    {
        [$key] = ApiKey::generate(User::factory()->create(), 'Milik orang lain');

        $this->actingAs(User::factory()->create())
            ->delete(route('api-keys.destroy', $key))
            ->assertNotFound();

        $this->assertDatabaseHas('api_keys', ['id' => $key->id]);
    }

    public function test_index_lists_only_own_keys_without_secrets(): void
    {
        $user = User::factory()->create();
        ApiKey::generate($user, 'Punya saya');
        ApiKey::generate(User::factory()->create(), 'Punya orang lain');

        $this->actingAs($user)->get(route('api-keys'))
            ->assertInertia(fn ($page) => $page
                ->component('api-keys')
                ->has('keys', 1)
                ->where('keys.0.name', 'Punya saya')
                ->where('keys.0.is_active', true)
                ->missing('keys.0.key_hash'));
    }

    public function test_cli_me_returns_the_key_owner(): void
    {
        $user = User::factory()->create(['name' => 'Budi', 'email' => 'budi@example.com']);
        [$key, $plain] = ApiKey::generate($user, 'CLI');

        $this->withToken($plain)->getJson('/api/cli/me')
            ->assertOk()
            ->assertExactJson(['user' => ['name' => 'Budi', 'email' => 'budi@example.com']]);

        $this->assertNotNull($key->refresh()->last_used_at);
    }

    public function test_cli_me_rejects_missing_or_invalid_or_deleted_keys(): void
    {
        $this->getJson('/api/cli/me')->assertUnauthorized();
        $this->withToken('cpc_salah')->getJson('/api/cli/me')->assertUnauthorized();

        [$key, $plain] = ApiKey::generate(User::factory()->create(), 'Dicabut');
        $key->delete();
        $this->withToken($plain)->getJson('/api/cli/me')->assertUnauthorized();
    }

    public function test_guests_cannot_manage_keys(): void
    {
        $this->post(route('api-keys.store'), ['name' => 'x'])->assertRedirect(route('login'));
    }

    public function test_user_can_rename_a_key(): void
    {
        $user = User::factory()->create();
        [$key] = ApiKey::generate($user, 'Lama');

        $this->actingAs($user)
            ->patch(route('api-keys.update', $key), ['name' => 'Baru'])
            ->assertRedirect(route('api-keys'));

        $this->assertSame('Baru', $key->refresh()->name);
        $this->assertTrue($key->is_active);
    }

    public function test_rename_requires_a_name(): void
    {
        $user = User::factory()->create();
        [$key] = ApiKey::generate($user, 'Lama');

        $this->actingAs($user)
            ->patch(route('api-keys.update', $key), ['name' => ''])
            ->assertSessionHasErrors('name');
    }

    public function test_user_can_deactivate_and_reactivate_a_key(): void
    {
        $user = User::factory()->create();
        [$key, $plain] = ApiKey::generate($user, 'CLI');

        $this->actingAs($user)->patch(route('api-keys.update', $key), ['is_active' => false]);
        $this->assertFalse($key->refresh()->is_active);
        $this->assertSame('CLI', $key->name);
        $this->withToken($plain)->getJson('/api/cli/me')->assertUnauthorized();

        $this->actingAs($user)->patch(route('api-keys.update', $key), ['is_active' => true]);
        $this->assertTrue($key->refresh()->is_active);
        $this->withToken($plain)->getJson('/api/cli/me')->assertOk();
    }

    public function test_user_cannot_update_someone_elses_key(): void
    {
        [$key] = ApiKey::generate(User::factory()->create(), 'Milik orang lain');

        $this->actingAs(User::factory()->create())
            ->patch(route('api-keys.update', $key), ['name' => 'Dibajak', 'is_active' => false])
            ->assertNotFound();

        $this->assertSame('Milik orang lain', $key->refresh()->name);
        $this->assertTrue($key->is_active);
    }

    public function test_index_reports_summary_stats(): void
    {
        $user = User::factory()->create();
        [$used] = ApiKey::generate($user, 'Dipakai');
        [$stale] = ApiKey::generate($user, 'Lama');
        [$off] = ApiKey::generate($user, 'Mati');
        ApiKey::generate($user, 'Baru');

        $used->forceFill(['last_used_at' => now()->subDay()])->save();
        $stale->forceFill(['last_used_at' => now()->subDays(30)])->save();
        $off->update(['is_active' => false]);

        $this->actingAs($user)->get(route('api-keys'))
            ->assertInertia(fn ($page) => $page
                ->where('stats.total', 4)
                ->where('stats.active', 3)
                ->where('stats.inactive', 1)
                ->where('stats.used_7d', 1));
    }
}

<?php

namespace Tests\Feature;

use App\Models\User;
use Illuminate\Foundation\Testing\RefreshDatabase;
use PHPUnit\Framework\Attributes\DataProvider;
use Tests\TestCase;

class MemberPagesTest extends TestCase
{
    use RefreshDatabase;

    /**
     * @return array<string, array{string}>
     */
    public static function pages(): array
    {
        return [
            'setup' => ['setup'],
            'agents' => ['agents'],
            'sessions' => ['sessions'],
            'api keys' => ['api-keys'],
            'usage' => ['usage'],
            'billing' => ['billing'],
        ];
    }

    #[DataProvider('pages')]
    public function test_guests_are_redirected_to_login(string $route): void
    {
        $this->get(route($route))->assertRedirect(route('login'));
    }

    #[DataProvider('pages')]
    public function test_authenticated_users_can_view_page(string $route): void
    {
        $this->actingAs(User::factory()->create())
            ->get(route($route))
            ->assertOk();
    }
}

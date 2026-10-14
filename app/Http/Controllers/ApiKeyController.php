<?php

namespace App\Http\Controllers;

use App\Http\Requests\StoreApiKeyRequest;
use App\Http\Requests\UpdateApiKeyRequest;
use App\Models\ApiKey;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class ApiKeyController extends Controller
{
    public function index(Request $request): Response
    {
        $keys = $request->user()->apiKeys()->latest()->get();

        return Inertia::render('api-keys', [
            'keys' => $keys->map(fn (ApiKey $key) => [
                'id' => $key->id,
                'name' => $key->name,
                'prefix' => $key->prefix,
                'is_active' => $key->is_active,
                'created_at' => $key->created_at->toIso8601String(),
                'last_used_at' => $key->last_used_at?->toIso8601String(),
            ]),
            'stats' => [
                'total' => $keys->count(),
                'active' => $keys->where('is_active', true)->count(),
                'inactive' => $keys->where('is_active', false)->count(),
                'used_7d' => $keys->filter(
                    fn (ApiKey $key) => $key->last_used_at?->gte(now()->subDays(7)) ?? false,
                )->count(),
            ],
            'limit' => ApiKey::MAX_PER_USER,
            'newKey' => $request->session()->get('newKey'),
        ]);
    }

    public function store(StoreApiKeyRequest $request): RedirectResponse
    {
        [, $plain] = ApiKey::generate($request->user(), $request->validated('name'));

        return to_route('api-keys')->with('newKey', $plain);
    }

    public function update(UpdateApiKeyRequest $request, ApiKey $apiKey): RedirectResponse
    {
        abort_unless($apiKey->user_id === $request->user()->id, 404);

        $apiKey->update($request->validated());

        Inertia::flash('toast', ['type' => 'success', 'message' => __('API key diperbarui.')]);

        return to_route('api-keys');
    }

    public function destroy(Request $request, ApiKey $apiKey): RedirectResponse
    {
        abort_unless($apiKey->user_id === $request->user()->id, 404);

        $apiKey->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('API key dihapus.')]);

        return to_route('api-keys');
    }
}

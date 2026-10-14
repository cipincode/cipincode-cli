<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\ApiKey;
use App\Models\CliAuthCode;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class CliController extends Controller
{
    /**
     * Identify the user that owns the API key used by the CLI.
     */
    public function me(Request $request): JsonResponse
    {
        $user = $request->user();

        return response()->json([
            'user' => [
                'name' => $user->name,
                'email' => $user->email,
            ],
        ]);
    }

    /**
     * Tukar kode dari browser dengan API key baru.
     */
    public function exchange(Request $request): JsonResponse
    {
        $data = $request->validate([
            'code' => ['required', 'string', 'max:32'],
            'verifier' => ['required', 'string', 'min:43', 'max:128'],
            'device' => ['nullable', 'string', 'max:60'],
        ]);

        $user = CliAuthCode::consume($data['code'], $data['verifier']);

        if (! $user) {
            return response()->json(['message' => 'Kode tidak valid atau sudah kedaluwarsa.'], 422);
        }

        if ($user->apiKeys()->count() >= ApiKey::MAX_PER_USER) {
            return response()->json([
                'message' => 'Batas '.ApiKey::MAX_PER_USER.' API key tercapai. Hapus key lama di dashboard lalu coba lagi.',
            ], 422);
        }

        $device = trim($data['device'] ?? '');
        [, $plain] = ApiKey::generate($user, $device !== '' ? 'CLI · '.$device : 'CLI');

        return response()->json([
            'key' => $plain,
            'user' => ['name' => $user->name, 'email' => $user->email],
        ]);
    }
}

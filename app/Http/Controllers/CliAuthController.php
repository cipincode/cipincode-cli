<?php

namespace App\Http\Controllers;

use App\Models\CliAuthCode;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class CliAuthController extends Controller
{
    /**
     * Halaman otorisasi yang dibuka dari CLI. Tamu diarahkan login dulu
     * (email atau Google), lalu kembali ke sini.
     */
    public function show(Request $request): Response|RedirectResponse
    {
        $data = $request->validate([
            'challenge' => ['required', 'regex:/^[A-Za-z0-9_-]{43}$/'],
            'device' => ['nullable', 'string', 'max:60'],
            'method' => ['nullable', 'in:email,google'],
        ]);

        $user = $request->user();

        if (! $user) {
            $login = ($data['method'] ?? null) === 'google' && config('services.google.client_id')
                ? route('google.redirect')
                : route('login');

            return redirect()->guest($login);
        }

        return Inertia::render('auth/cli-authorize', [
            'challenge' => $data['challenge'],
            'device' => $data['device'] ?? null,
            'user' => ['name' => $user->name, 'email' => $user->email],
            'code' => $request->session()->get('cliCode'),
            'ttl' => CliAuthCode::TTL_MINUTES,
        ]);
    }

    public function approve(Request $request): RedirectResponse
    {
        $data = $request->validate([
            'challenge' => ['required', 'regex:/^[A-Za-z0-9_-]{43}$/'],
            'device' => ['nullable', 'string', 'max:60'],
        ]);

        $code = CliAuthCode::issue($request->user(), $data['challenge']);

        return to_route('cli.authorize', array_filter([
            'challenge' => $data['challenge'],
            'device' => $data['device'] ?? null,
        ]))->with('cliCode', $code);
    }
}

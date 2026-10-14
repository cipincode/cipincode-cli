<?php

namespace App\Http\Controllers\Auth;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Str;
use Laravel\Socialite\Facades\Socialite;
use Laravel\Socialite\Two\InvalidStateException;
use Symfony\Component\HttpFoundation\RedirectResponse as SymfonyRedirect;

class GoogleAuthController extends Controller
{
    public function redirect(): SymfonyRedirect
    {
        abort_unless(self::enabled(), 404);

        return Socialite::driver('google')->redirect();
    }

    public function callback(Request $request): RedirectResponse
    {
        abort_unless(self::enabled(), 404);

        try {
            $google = Socialite::driver('google')->user();
        } catch (InvalidStateException) {
            return $this->failed(__('Sesi login Google kedaluwarsa. Silakan coba lagi.'));
        }

        // Google menandai email yang sudah diverifikasi. Tanpa itu, jangan percaya.
        if (! $google->getEmail() || ! ($google->user['email_verified'] ?? false)) {
            return $this->failed(__('Email akun Google Anda belum terverifikasi.'));
        }

        $user = User::query()->where('email', $google->getEmail())->first();

        // Login Google tidak melewati tantangan 2FA, jadi akun ber-2FA harus masuk lewat email.
        if ($user?->two_factor_confirmed_at) {
            return $this->failed(__('Akun ini memakai verifikasi dua langkah. Masuklah dengan email dan password.'));
        }

        $user ??= User::forceCreate([
            'name' => $google->getName() ?: Str::before($google->getEmail(), '@'),
            'email' => $google->getEmail(),
            'password' => Str::random(40),
            'email_verified_at' => now(),
        ]);

        Auth::login($user, remember: true);
        $request->session()->regenerate();

        return redirect()->intended(route('dashboard'));
    }

    public static function enabled(): bool
    {
        return filled(config('services.google.client_id')) && filled(config('services.google.client_secret'));
    }

    private function failed(string $message): RedirectResponse
    {
        return to_route('login')->withErrors(['email' => $message]);
    }
}

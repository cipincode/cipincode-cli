<?php

namespace App\Models;

use Carbon\CarbonInterface;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Support\Str;

/**
 * Kode sekali pakai yang ditempel pengguna ke CLI setelah login di browser.
 *
 * @property int $id
 * @property int $user_id
 * @property string $code_hash
 * @property string $challenge
 * @property CarbonInterface $expires_at
 * @property CarbonInterface|null $used_at
 */
class CliAuthCode extends Model
{
    public const TTL_MINUTES = 5;

    /** Tanpa karakter yang mirip (0/O, 1/I/L) agar mudah dibaca dan diketik. */
    private const ALPHABET = 'ABCDEFGHJKMNPQRSTUVWXYZ23456789';

    protected $fillable = ['code_hash', 'challenge', 'expires_at'];

    protected function casts(): array
    {
        return [
            'expires_at' => 'datetime',
            'used_at' => 'datetime',
        ];
    }

    /**
     * Terbitkan kode untuk user. Kode asli hanya dikembalikan di sini.
     */
    public static function issue(User $user, string $challenge): string
    {
        $plain = '';
        for ($i = 0; $i < 8; $i++) {
            $plain .= self::ALPHABET[random_int(0, strlen(self::ALPHABET) - 1)];
        }

        $user->cliAuthCodes()->create([
            'code_hash' => self::hash($plain),
            'challenge' => $challenge,
            'expires_at' => now()->addMinutes(self::TTL_MINUTES),
        ]);

        return substr($plain, 0, 4).'-'.substr($plain, 4);
    }

    /**
     * Tukar kode dengan pemiliknya. Kode hangus setelah dipakai satu kali.
     */
    public static function consume(string $code, string $verifier): ?User
    {
        $record = self::query()
            ->where('code_hash', self::hash($code))
            ->whereNull('used_at')
            ->where('expires_at', '>', now())
            ->first();

        if (! $record || ! hash_equals($record->challenge, self::challengeFor($verifier))) {
            return null;
        }

        // Update bersyarat agar dua permintaan bersamaan tidak sama-sama lolos.
        $claimed = self::query()
            ->whereKey($record->id)
            ->whereNull('used_at')
            ->update(['used_at' => now()]);

        return $claimed === 1 ? $record->user : null;
    }

    /**
     * Tantangan PKCE: base64url(sha256(verifier)).
     */
    public static function challengeFor(string $verifier): string
    {
        return rtrim(strtr(base64_encode(hash('sha256', $verifier, true)), '+/', '-_'), '=');
    }

    public static function hash(string $code): string
    {
        return hash('sha256', Str::upper(preg_replace('/[^A-Za-z0-9]/', '', $code) ?? ''));
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }
}

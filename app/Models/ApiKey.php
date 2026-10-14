<?php

namespace App\Models;

use Carbon\CarbonInterface;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Support\Str;

/**
 * @property int $id
 * @property int $user_id
 * @property string $name
 * @property string $prefix
 * @property string $key_hash
 * @property bool $is_active
 * @property CarbonInterface|null $last_used_at
 */
class ApiKey extends Model
{
    public const PREFIX = 'cpc_';

    public const MAX_PER_USER = 10;

    protected $fillable = ['name', 'prefix', 'key_hash', 'is_active'];

    protected function casts(): array
    {
        return [
            'is_active' => 'boolean',
            'last_used_at' => 'datetime',
        ];
    }

    /**
     * Create a key for the user. The plain text key is only returned here.
     *
     * @return array{0: self, 1: string}
     */
    public static function generate(User $user, string $name): array
    {
        $plain = self::PREFIX.Str::random(40);

        $key = $user->apiKeys()->create([
            'name' => $name,
            'prefix' => substr($plain, 0, 8),
            'key_hash' => self::hash($plain),
        ]);

        return [$key, $plain];
    }

    public static function hash(string $plain): string
    {
        return hash('sha256', $plain);
    }

    public static function findByPlain(string $plain): ?self
    {
        return self::query()->where('key_hash', self::hash($plain))->first();
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }
}

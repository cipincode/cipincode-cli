<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

/**
 * Agent buatan pengguna: instruksi tambahan plus mode dan effort bawaan.
 *
 * @property int $id
 * @property int $user_id
 * @property string $name
 * @property string|null $description
 * @property string $instructions
 * @property string $mode
 * @property string $effort
 */
class Agent extends Model
{
    public const MAX_PER_USER = 20;

    protected $fillable = ['name', 'description', 'instructions', 'mode', 'effort'];

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }
}

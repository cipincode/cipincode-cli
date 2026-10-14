<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasMany;

/**
 * Satu percakapan di CLI (dimulai lagi dengan /new).
 *
 * @property int $id
 * @property string $uuid
 * @property int $user_id
 * @property string $title
 * @property string $mode
 * @property string|null $model
 * @property int $messages_count
 * @property \Carbon\CarbonInterface $last_activity_at
 */
class ChatSession extends Model
{
    protected $fillable = ['uuid', 'title', 'mode', 'model'];

    protected function casts(): array
    {
        return [
            'last_activity_at' => 'datetime',
        ];
    }

    /**
     * @return BelongsTo<User, $this>
     */
    public function user(): BelongsTo
    {
        return $this->belongsTo(User::class);
    }

    /**
     * @return HasMany<ChatMessage, $this>
     */
    public function messages(): HasMany
    {
        return $this->hasMany(ChatMessage::class);
    }

    public function getRouteKeyName(): string
    {
        return 'uuid';
    }
}

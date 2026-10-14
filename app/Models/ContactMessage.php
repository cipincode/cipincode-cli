<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

/**
 * @property int $id
 * @property string $name
 * @property string $email
 * @property string $message
 * @property string $status
 * @property \Illuminate\Support\Carbon $created_at
 * @property \Illuminate\Support\Carbon $updated_at
 */
class ContactMessage extends Model
{
    public const STATUSES = ['baru', 'dibaca', 'dibalas'];

    protected $fillable = ['name', 'email', 'message', 'status'];
}

<?php

namespace App\Console\Commands;

use App\Models\User;
use Illuminate\Console\Command;

class SprAdmin extends Command
{
    protected $signature = 'spr:admin
        {email : Email pengguna}
        {--revoke : Cabut akses superadmin}';

    protected $description = 'Beri atau cabut akses superadmin (/spr)';

    public function handle(): int
    {
        $user = User::query()->where('email', $this->argument('email'))->first();
        if (! $user) {
            $this->error('Pengguna dengan email itu tidak ditemukan.');

            return self::FAILURE;
        }

        $user->forceFill(['is_admin' => ! $this->option('revoke')])->save();

        $this->info(sprintf(
            '%s %s superadmin.',
            $user->email,
            $this->option('revoke') ? 'bukan lagi' : 'sekarang',
        ));

        return self::SUCCESS;
    }
}

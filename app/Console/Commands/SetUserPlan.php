<?php

namespace App\Console\Commands;

use App\Models\User;
use App\Support\Plans;
use Illuminate\Console\Command;

class SetUserPlan extends Command
{
    protected $signature = 'plan:set
        {email : Email pengguna}
        {plan : free, go, pro, atau max}
        {--days= : Masa berlaku dalam hari. Kosong: gratis memakai masa percobaan bawaan, paket berbayar tanpa batas waktu}';

    protected $description = 'Ubah paket pengguna (gratis atau berbayar)';

    public function handle(): int
    {
        $plan = (string) $this->argument('plan');
        if (! in_array($plan, Plans::keys(), true)) {
            $this->error('Paket harus salah satu dari: '.implode(', ', Plans::keys()).'.');

            return self::FAILURE;
        }

        $user = User::query()->where('email', $this->argument('email'))->first();
        if (! $user) {
            $this->error('Pengguna dengan email itu tidak ditemukan.');

            return self::FAILURE;
        }

        $days = $this->option('days');
        if ($days !== null && (! ctype_digit((string) $days) || (int) $days < 1)) {
            $this->error('--days harus bilangan bulat positif.');

            return self::FAILURE;
        }

        $expires = match (true) {
            $days !== null => now()->addDays((int) $days),
            $plan === 'free' => now()->addDays((int) config('usage.plans.free.trial_days')),
            default => null,
        };

        $user->forceFill(['plan' => $plan, 'plan_expires_at' => $expires])->save();

        $this->info(sprintf(
            '%s sekarang paket %s%s.',
            $user->email,
            config("usage.plans.$plan.label"),
            $expires ? ', berlaku sampai '.$expires->format('d M Y H:i') : ' (tanpa batas waktu)',
        ));

        return self::SUCCESS;
    }
}

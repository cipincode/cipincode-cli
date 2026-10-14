<?php

namespace App\Services;

use App\Models\UsageEvent;
use App\Models\User;
use Carbon\CarbonInterface;
use Illuminate\Support\Facades\DB;

/**
 * Menghitung pemakaian dan menegakkan batas sesuai paket pengguna.
 *
 *  - sesi: mulai saat pesan pertama, reset otomatis setelah N jam (semua paket);
 *  - mingguan: reset setiap awal minggu (Senin 00.00), hanya paket yang punya batas mingguan;
 *  - paket: paket gratis punya masa percobaan, setelah itu semua permintaan ditolak.
 */
class UsageTracker
{
    /**
     * @return array{session: array<string, mixed>, weekly: array<string, mixed>|null, plan: array<string, mixed>}
     */
    public function snapshot(User $user): array
    {
        $now = now();
        $plan = $this->planConfig($user);

        $started = $user->usage_window_started_at;
        $sessionEnds = $started?->addHours($this->sessionHours());
        $sessionActive = $sessionEnds && $sessionEnds->gt($now);

        $sessionUsed = $sessionActive
            ? (int) $user->usageEvents()->where('created_at', '>=', $started)->sum('units')
            : 0;

        $weeklyLimit = $plan['limits']['weekly'];
        $weekStart = $now->startOfWeek();
        $weekly = $weeklyLimit === null ? null : $this->meter(
            (int) $user->usageEvents()->where('created_at', '>=', $weekStart)->sum('units'),
            (int) $weeklyLimit,
            $weekStart->addWeek(),
        );

        return [
            'session' => $this->meter($sessionUsed, (int) $plan['limits']['session'], $sessionActive ? $sessionEnds : null),
            'weekly' => $weekly,
            'plan' => [
                'name' => $user->plan,
                'label' => $plan['label'],
                'expires_at' => $user->plan_expires_at?->toIso8601String(),
                'expired' => $user->plan_expires_at !== null && $user->plan_expires_at->lte($now),
            ],
        ];
    }

    /**
     * Pakai $units dari jatah. Mengembalikan [diizinkan, snapshot, penghalang].
     *
     * @return array{0: bool, 1: array<string, mixed>, 2: string|null}
     */
    public function consume(User $user, int $units = 1): array
    {
        return DB::transaction(function () use ($user, $units): array {
            $locked = User::query()->whereKey($user->getKey())->lockForUpdate()->firstOrFail();
            $snapshot = $this->snapshot($locked);

            if ($blockedBy = $this->blockedBy($snapshot)) {
                return [false, $snapshot, $blockedBy];
            }

            $now = now();
            $started = $locked->usage_window_started_at;
            if (! $started || $started->addHours($this->sessionHours())->lte($now)) {
                $locked->forceFill(['usage_window_started_at' => $now])->save();
            }

            $event = new UsageEvent;
            $event->user_id = $locked->getKey();
            $event->units = $units;
            $event->created_at = $now;
            $event->save();

            return [true, $this->snapshot($locked->refresh()), null];
        });
    }

    /**
     * Penghalang pertama: "plan" (masa aktif habis), "session", atau "weekly".
     * Null jika masih boleh memakai.
     *
     * @param  array<string, mixed>  $snapshot
     */
    public function blockedBy(array $snapshot): ?string
    {
        if ($snapshot['plan']['expired']) {
            return 'plan';
        }

        foreach (['session', 'weekly'] as $name) {
            if ($snapshot[$name] !== null && $snapshot[$name]['exhausted']) {
                return $name;
            }
        }

        return null;
    }

    /**
     * @param  array<string, mixed>|null  $snapshot
     */
    public function blockedMessage(string $blockedBy, ?array $snapshot = null): string
    {
        if ($blockedBy === 'plan') {
            return ($snapshot['plan']['name'] ?? 'free') === 'free'
                ? 'Masa percobaan gratis '.config('usage.plans.free.trial_days').' hari telah berakhir. Beralih ke paket berbayar untuk melanjutkan.'
                : 'Masa aktif paket Anda telah berakhir. Perpanjang untuk melanjutkan.';
        }

        return $blockedBy === 'session'
            ? 'Batas pemakaian sesi '.$this->sessionHours().' jam tercapai.'
            : 'Batas pemakaian mingguan tercapai.';
    }

    /**
     * @return array{label: string, limits: array{session: int, weekly: int|null}}
     */
    private function planConfig(User $user): array
    {
        return config("usage.plans.{$user->plan}") ?? config('usage.plans.free');
    }

    /**
     * @return array{used: int, limit: int, percent: int, exhausted: bool, resets_at: string|null}
     */
    private function meter(int $used, int $limit, ?CarbonInterface $resetsAt): array
    {
        return [
            'used' => $used,
            'limit' => $limit,
            'percent' => $limit > 0 ? min(100, (int) round($used / $limit * 100)) : 0,
            'exhausted' => $limit > 0 && $used >= $limit,
            'resets_at' => $resetsAt?->toIso8601String(),
        ];
    }

    private function sessionHours(): int
    {
        return (int) config('usage.session_hours', 5);
    }
}

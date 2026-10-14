<?php

namespace App\Support;

/**
 * Katalog paket untuk ditampilkan (halaman Usage dan Billing).
 * Semua angka dibaca dari config/usage.php agar tidak pernah berbeda dari yang ditegakkan.
 */
class Plans
{
    /**
     * Urutan paket dari terendah ke tertinggi.
     *
     * @return array<int, array<string, mixed>>
     */
    public static function catalog(): array
    {
        $sessionHours = (int) config('usage.session_hours', 5);
        $catalog = [];

        foreach (config('usage.plans') as $key => $plan) {
            $session = $plan['limits']['session'];
            $weekly = $plan['limits']['weekly'];
            $free = $key === 'free';

            $features = [
                self::number($session)." pesan per {$sessionHours} jam",
                $weekly === null ? 'Tanpa batas mingguan' : self::number($weekly).' pesan per minggu',
                'Semua mode dan effort',
                $free ? 'Model gratis' : 'Model berbayar',
            ];
            if ($free) {
                $features[] = 'Masa percobaan '.$plan['trial_days'].' hari';
            }

            $catalog[] = [
                'key' => $key,
                'label' => $plan['label'],
                'price' => $plan['price'],
                'price_label' => $plan['price'] > 0 ? self::rupiah($plan['price']) : 'Gratis',
                'period' => $plan['price'] > 0 ? 'per bulan' : null,
                'session_limit' => $session,
                'weekly_limit' => $weekly,
                'features' => $features,
                'paid' => ! $free,
            ];
        }

        return $catalog;
    }

    /** @return array<int, string> */
    public static function keys(): array
    {
        return array_keys(config('usage.plans'));
    }

    public static function rupiah(int $amount): string
    {
        return 'Rp '.self::number($amount);
    }

    private static function number(int $value): string
    {
        return number_format($value, 0, ',', '.');
    }
}

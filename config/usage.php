<?php

return [

    /*
    |--------------------------------------------------------------------------
    | Paket dan batas pemakaian
    |--------------------------------------------------------------------------
    |
    | Satuan batas adalah jumlah pesan. Sesi dimulai saat pesan pertama dan otomatis
    | reset setelah `session_hours` jam. Batas mingguan reset setiap Senin 00.00.
    | Harga dalam rupiah per bulan. weekly bernilai null berarti tanpa batas mingguan.
    |
    | - free: hanya batas sesi, masa percobaan `trial_days` hari sejak daftar, model gratis.
    | - go, pro, max: paket berbayar dengan batas sesi dan mingguan yang makin besar.
    |   `model` (opsional) menimpa model berbayar bawaan untuk paket itu.
    |
    */

    'session_hours' => 5,

    'plans' => [
        'free' => [
            'label' => 'Gratis',
            'price' => 0,
            'trial_days' => 14,
            'limits' => [
                'session' => (int) env('USAGE_FREE_LIMIT_SESSION', 20),
                'weekly' => null,
            ],
        ],
        'go' => [
            'label' => 'Go',
            'price' => (int) env('PLAN_GO_PRICE', 49000),
            'model' => env('OPENROUTER_MODEL_GO'),
            'limits' => [
                'session' => (int) env('USAGE_GO_LIMIT_SESSION', 40),
                'weekly' => (int) env('USAGE_GO_LIMIT_WEEKLY', 400),
            ],
        ],
        'pro' => [
            'label' => 'Pro',
            'price' => (int) env('PLAN_PRO_PRICE', 149000),
            'model' => env('OPENROUTER_MODEL_PRO'),
            'limits' => [
                'session' => (int) env('USAGE_PRO_LIMIT_SESSION', 100),
                'weekly' => (int) env('USAGE_PRO_LIMIT_WEEKLY', 1500),
            ],
        ],
        'max' => [
            'label' => 'Max',
            'price' => (int) env('PLAN_MAX_PRICE', 499000),
            'model' => env('OPENROUTER_MODEL_MAX'),
            'limits' => [
                'session' => (int) env('USAGE_MAX_LIMIT_SESSION', 300),
                'weekly' => (int) env('USAGE_MAX_LIMIT_WEEKLY', 5000),
            ],
        ],
    ],

];

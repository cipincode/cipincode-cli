<?php

/*
 * Versi CLI terbaru dan daftar perubahannya untuk modal pengingat update
 * di terminal (GET /api/cli/version).
 *
 * Naikkan `version` setiap merilis binary baru (Makefile VERSION dan
 * install.sh mengikuti angka ini) lalu isi `added` (fitur baru) dan
 * `fixed` (perbaikan) dalam Bahasa Indonesia.
 */
return [
    'version' => env('CLI_VERSION', '0.2.9'),

    'changelog' => [
        'added' => [
            'Status shell di bawah kolom input: bulatan kuning + jumlah saat berjalan, hijau saat selesai',
        ],
        'fixed' => [
            'Counter shell di-reset tiap pesan baru agar status selalu akurat',
        ],
    ],
];

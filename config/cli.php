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
    'version' => env('CLI_VERSION', '0.2.4'),

    'changelog' => [
        'added' => [
            '/model kini memuat daftar langsung dari akun provider Anda bila tersambung dengan kunci sendiri (disaring ke model obrolan, nama ramah bila dikenal)',
        ],
        'fixed' => [
            'Bila akun provider tidak bisa dihubungi, /model otomatis memakai katalog bawaan',
        ],
    ],
];

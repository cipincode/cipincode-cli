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
    'version' => env('CLI_VERSION', '0.2.5'),

    'changelog' => [
        'added' => [
            '/usage kini mengikuti akun provider saat mode langsung (akun Claude / kunci sendiri), bukan jatah server 5 jam',
            'Panggilan tool AI (mis. execute_command) tampil sebagai kotak rapi, bukan XML mentah',
        ],
        'fixed' => [
            'Header chat menandai mode langsung agar tidak tertukar dengan kuota server',
        ],
    ],
];

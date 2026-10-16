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
    'version' => env('CLI_VERSION', '0.2.2'),

    'changelog' => [
        'added' => [],
        'fixed' => [
            'Perbaikan proses rilis npm (Trusted Publishing OIDC, Node 22) agar `npm update -g cipincode` lebih andal',
            'Sinkron versi 0.2.2 di binary, installer, dan server agar pengingat update muncul tepat',
        ],
    ],
];

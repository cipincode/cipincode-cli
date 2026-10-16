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
    'version' => env('CLI_VERSION', '0.2.3'),

    'changelog' => [
        'added' => [
            'Pilihan Auth Login dan API Key di /connect → Anthropic: login Pro/Max lewat browser (link resmi claude.ai) atau tempel kunci sendiri',
        ],
        'fixed' => [
            'Opsi API Key Anthropic dikembalikan agar bisa langsung dengan kunci sendiri tanpa lewat server',
        ],
    ],
];

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
    'version' => env('CLI_VERSION', '0.2.6'),

    'changelog' => [
        'added' => [
            '/update otomatis via install.sh tanpa npm, lalu CLI me-restart sendiri — auth dan pengaturan tetap dipakai',
        ],
        'fixed' => [
            'Modal update mendeteksi cara install (npm vs curl) agar perintahnya selalu benar',
        ],
    ],
];

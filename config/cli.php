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
    'version' => env('CLI_VERSION', '0.2.8'),

    'changelog' => [
        'added' => [
            'Mode Coding/Debugging kini benar-benar menjalankan perintah shell AI (execute_command) otomatis dan memakai hasilnya',
        ],
        'fixed' => [
            'System prompt server memberi tahu model soal kemampuan shell; mode Chat tetap pratinjau saja',
        ],
    ],
];

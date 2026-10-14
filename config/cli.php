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
    'version' => env('CLI_VERSION', '0.2.0'),

    'changelog' => [
        'added' => [
            'Perintah /connect — pilih penyedia API (Anthropic, OpenAI, Google, xAI, DeepSeek, Meta) dari modal',
            'Perintah /model — pilih model dari provider yang sedang aktif',
            'Modal pengingat update otomatis dengan checkbox jangan ingatkan lagi',
        ],
        'fixed' => [
            'Pilihan provider dan model kini tersimpan antar sesi',
            'Paket gratis selalu memakai model gratis walau ada pilihan lain',
        ],
    ],
];

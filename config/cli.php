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
    'version' => env('CLI_VERSION', '0.2.1'),

    'changelog' => [
        'added' => [
            'Login akun Claude (Pro/Max) lewat browser di /connect → Anthropic, tanpa kunci API',
            'Perintah /connect — pilih penyedia API (Anthropic, OpenAI, Google, xAI, DeepSeek, Meta) dari modal',
            'Perintah /model — pilih model dari provider yang sedang aktif',
            'Modal pengingat update otomatis dengan checkbox jangan ingatkan lagi',
        ],
        'fixed' => [
            'Opsi kunci API Anthropic dihapus dari /connect; gunakan akun Claude atau via server',
            'Pilihan provider dan model kini tersimpan antar sesi',
            'Paket gratis selalu memakai model gratis walau ada pilihan lain',
        ],
    ],
];

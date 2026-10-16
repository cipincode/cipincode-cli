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
    'version' => env('CLI_VERSION', '0.2.7'),

    'changelog' => [
        'added' => [
            'Chat via server kembali lewat agent lokal sehingga tidak lagi 404',
            'Login ulang mempertahankan kunci API dan akun Claude yang sudah dipilih',
        ],
        'fixed' => [
            'Hapus default agent remote yang tidak ada servernya; --server tetap didukung eksplisit',
        ],
    ],
];

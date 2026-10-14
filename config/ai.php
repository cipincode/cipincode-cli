<?php

/*
 * Katalog provider API dan model yang boleh dipilih lewat /connect dan /model
 * di CLI. Daftar ini harus sama dengan `providers` di cli/tui/provider.go.
 *
 * Paket gratis selalu memakai model gratis (aman dari biaya). Paket berbayar
 * boleh meminta salah satu model di bawah; permintaan lain diabaikan dan
 * kembali ke model bawaan paket.
 */
return [
    'providers' => [
        'anthropic' => [
            'anthropic/claude-sonnet-4-5',
            'anthropic/claude-opus-4-1',
            'anthropic/claude-3-5-haiku',
        ],
        'openai' => [
            'openai/gpt-4o',
            'openai/gpt-4o-mini',
            'openai/o4-mini',
        ],
        'google' => [
            'google/gemini-2.5-pro',
            'google/gemini-2.5-flash',
        ],
        'xai' => [
            'x-ai/grok-3',
            'x-ai/grok-3-mini',
        ],
        'deepseek' => [
            'deepseek/deepseek-chat',
            'deepseek/deepseek-r1',
        ],
        'meta' => [
            'meta-llama/llama-3.3-70b-instruct',
        ],
    ],
];

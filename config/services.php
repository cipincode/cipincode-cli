<?php

return [

    /*
    |--------------------------------------------------------------------------
    | Third Party Services
    |--------------------------------------------------------------------------
    |
    | This file is for storing the credentials for third party services such
    | as Resend, Postmark, AWS, and more. This file provides the de facto
    | location for this type of information, allowing packages to have
    | a conventional file to locate the various service credentials.
    |
    */

    'postmark' => [
        'key' => env('POSTMARK_API_KEY'),
    ],

    'resend' => [
        'key' => env('RESEND_API_KEY'),
    ],

    'ses' => [
        'key' => env('AWS_ACCESS_KEY_ID'),
        'secret' => env('AWS_SECRET_ACCESS_KEY'),
        'region' => env('AWS_DEFAULT_REGION', 'us-east-1'),
    ],

    'slack' => [
        'notifications' => [
            'bot_user_oauth_token' => env('SLACK_BOT_USER_OAUTH_TOKEN'),
            'channel' => env('SLACK_BOT_USER_DEFAULT_CHANNEL'),
        ],
    ],

    'google' => [
        'client_id' => env('GOOGLE_CLIENT_ID'),
        'client_secret' => env('GOOGLE_CLIENT_SECRET'),
        'redirect' => env('GOOGLE_REDIRECT_URI', '/auth/google/callback'),
    ],

    'openrouter' => [
        'key' => env('OPENROUTER_API_KEY'),
        'base_url' => env('OPENROUTER_BASE_URL', 'https://openrouter.ai/api/v1'),
        // Model untuk paket gratis. "openrouter/free" memilih otomatis di antara model gratis.
        'free_model' => env('OPENROUTER_FREE_MODEL', 'openrouter/free'),
        // Model paket berbayar. "openrouter/auto" memilih model otomatis. Ganti dengan id model tertentu bila perlu.
        'model' => env('OPENROUTER_MODEL', 'openrouter/auto'),
        // Opsional: model berbeda per mode. Kosong berarti memakai OPENROUTER_MODEL.
        'models' => [
            'chat' => env('OPENROUTER_MODEL_CHAT'),
            'coding' => env('OPENROUTER_MODEL_CODING'),
            'debugging' => env('OPENROUTER_MODEL_DEBUGGING'),
        ],
    ],

];

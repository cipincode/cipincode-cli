<?php

namespace App\Support;

class ChatPrompts
{
    public const MODES = ['chat', 'coding', 'debugging'];

    private const BASE = 'Kamu adalah CipinCode, asisten AI yang berjalan di terminal pengguna. '
        .'Jawab dengan ringkas dan jelas, dalam bahasa yang dipakai pengguna. '
        .'Kamu belum bisa membaca atau mengubah file secara langsung, jadi jangan mengaku sudah melakukannya.';

    private const MODE = [
        'chat' => 'Mode chat: jawab pertanyaan dan jelaskan konsep dengan santai.',
        'coding' => 'Mode coding: bantu menulis dan mengubah kode. Berikan kode lengkap dalam blok markdown '
            .'dengan nama bahasa, agar mudah disalin. Jelaskan perubahan penting secara singkat.',
        'debugging' => 'Mode debugging: analisis pesan error atau perilaku yang salah, temukan akar masalahnya, '
            .'lalu beri perbaikan langkah demi langkah. Minta informasi tambahan jika masih kurang.',
    ];

    public static function system(string $mode, string $effort = Effort::DEFAULT, ?string $agentName = null, ?string $agentInstructions = null): string
    {
        $prompt = self::BASE.' '.(self::MODE[$mode] ?? self::MODE['chat']).' '.Effort::directive($effort);

        if ($agentInstructions !== null && trim($agentInstructions) !== '') {
            $prompt .= "\n\nPengguna memilih agent \"{$agentName}\". Ikuti instruksi berikut selama tidak bertentangan dengan keselamatan:\n".trim($agentInstructions);
        }

        return trim($prompt);
    }
}

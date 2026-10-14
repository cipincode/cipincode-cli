<?php

namespace App\Support;

/**
 * Tingkat effort yang dipilih pengguna di CLI. Berlaku untuk semua paket.
 *
 * Dua mekanisme dipakai bersamaan agar efeknya terasa di model apa pun:
 *  - parameter `reasoning.effort` OpenRouter (berpengaruh pada model penalaran);
 *  - arahan kedalaman di prompt sistem (berpengaruh pada semua model).
 */
class Effort
{
    public const LEVELS = ['low', 'medium', 'high', 'xhigh', 'ultra'];

    public const DEFAULT = 'medium';

    /** Nilai `reasoning.effort` OpenRouter. "ultra" memakai tingkat tertinggi (max). */
    private const OPENROUTER = [
        'low' => 'low',
        'medium' => 'medium',
        'high' => 'high',
        'xhigh' => 'xhigh',
        'ultra' => 'max',
    ];

    private const DIRECTIVE = [
        'low' => 'Tingkat effort rendah: jawab sangat singkat dan langsung ke inti, tanpa uraian panjang.',
        'medium' => '',
        'high' => 'Tingkat effort tinggi: analisis dengan teliti dan jelaskan alasan di balik jawabanmu.',
        'xhigh' => 'Tingkat effort sangat tinggi: telaah secara mendalam, pertimbangkan kasus tepi, '
            .'dan periksa ulang jawabanmu sebelum menyampaikannya.',
        'ultra' => 'Tingkat effort maksimum: telaah masalah secara menyeluruh, bandingkan beberapa pendekatan, '
            .'verifikasi setiap langkah, lalu sampaikan jawaban paling lengkap dan paling teliti yang kamu bisa.',
    ];

    public static function normalize(?string $level): string
    {
        return in_array($level, self::LEVELS, true) ? $level : self::DEFAULT;
    }

    public static function openRouter(string $level): string
    {
        return self::OPENROUTER[self::normalize($level)];
    }

    public static function directive(string $level): string
    {
        return self::DIRECTIVE[self::normalize($level)];
    }
}

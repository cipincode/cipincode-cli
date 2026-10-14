<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;

class VersionController extends Controller
{
    /**
     * Versi CLI terbaru + changelog untuk modal update di terminal.
     * Publik agar CLI yang belum login pun bisa memeriksa pembaruan.
     */
    public function show(): JsonResponse
    {
        return response()->json([
            'version' => (string) config('cli.version', '0.1.0'),
            'added' => array_values((array) config('cli.changelog.added', [])),
            'fixed' => array_values((array) config('cli.changelog.fixed', [])),
        ]);
    }
}

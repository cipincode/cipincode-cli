<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Services\UsageTracker;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class UsageController extends Controller
{
    public function show(Request $request, UsageTracker $usage): JsonResponse
    {
        return response()->json($usage->snapshot($request->user()));
    }

    /**
     * Catat pemakaian sebelum CLI meminta jawaban. Ditolak (429) jika jatah habis.
     */
    public function consume(Request $request, UsageTracker $usage): JsonResponse
    {
        $data = $request->validate([
            'units' => ['nullable', 'integer', 'min:1', 'max:100'],
        ]);

        [$allowed, $snapshot, $blockedBy] = $usage->consume($request->user(), $data['units'] ?? 1);

        if (! $allowed) {
            return response()->json([
                'message' => $usage->blockedMessage($blockedBy, $snapshot),
                'blocked' => $blockedBy,
                ...$snapshot,
            ], 429);
        }

        return response()->json($snapshot);
    }
}

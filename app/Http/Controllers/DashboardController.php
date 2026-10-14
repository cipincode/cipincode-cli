<?php

namespace App\Http\Controllers;

use App\Services\UsageTracker;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class DashboardController extends Controller
{
    public function index(Request $request, UsageTracker $usage): Response
    {
        $user = $request->user();

        $days = collect(range(6, 0))->map(function (int $ago) use ($user): array {
            $day = now()->subDays($ago);

            return [
                'label' => $day->translatedFormat('D'),
                'date' => $day->toDateString(),
                'count' => (int) $user->usageEvents()
                    ->whereBetween('created_at', [$day->copy()->startOfDay(), $day->copy()->endOfDay()])
                    ->sum('units'),
            ];
        })->values();

        return Inertia::render('dashboard', [
            'usage' => $usage->snapshot($user),
            'stats' => [
                'sessions' => $user->chatSessions()->count(),
                'agents' => $user->agents()->count(),
                'keys' => $user->apiKeys()->where('is_active', true)->count(),
                'messages' => (int) $user->usageEvents()->sum('units'),
            ],
            'activity' => $days,
            'recent' => $user->chatSessions()
                ->latest('last_activity_at')
                ->limit(5)
                ->get(['uuid', 'title', 'mode', 'messages_count', 'last_activity_at']),
        ]);
    }
}

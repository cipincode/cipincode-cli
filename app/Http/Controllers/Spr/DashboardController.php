<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\Agent;
use App\Models\ApiKey;
use App\Models\ChatSession;
use App\Models\UsageEvent;
use App\Models\User;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class DashboardController extends Controller
{
    public function index(Request $request): Response
    {
        $days = collect(range(6, 0))->map(function (int $ago): array {
            $day = now()->subDays($ago);

            return [
                'label' => $day->translatedFormat('D'),
                'date' => $day->toDateString(),
                'count' => (int) UsageEvent::query()
                    ->whereBetween('created_at', [$day->copy()->startOfDay(), $day->copy()->endOfDay()])
                    ->sum('units'),
            ];
        })->values();

        return Inertia::render('spr/dashboard', [
            'stats' => [
                'users' => User::query()->count(),
                'users_week' => User::query()->where('created_at', '>=', now()->subDays(7))->count(),
                'messages' => (int) UsageEvent::query()->sum('units'),
                'sessions' => ChatSession::query()->count(),
                'agents' => Agent::query()->count(),
                'keys' => ApiKey::query()->where('is_active', true)->count(),
            ],
            'plans' => User::query()
                ->selectRaw('plan, COUNT(*) as total')
                ->groupBy('plan')
                ->pluck('total', 'plan'),
            'activity' => $days,
            'recent_users' => User::query()
                ->latest('id')
                ->limit(5)
                ->get(['id', 'name', 'email', 'plan', 'is_admin', 'created_at']),
        ]);
    }
}

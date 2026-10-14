<?php

namespace App\Http\Controllers;

use App\Models\ChatMessage;
use App\Models\ChatSession;
use App\Support\ChatPrompts;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class SessionController extends Controller
{
    public function index(Request $request): Response
    {
        $user = $request->user();
        $search = trim((string) $request->query('q', ''));
        $mode = in_array($request->query('mode'), ChatPrompts::MODES, true) ? $request->query('mode') : null;

        $sessions = $user->chatSessions()
            ->when($search !== '', fn ($q) => $q->where('title', 'like', '%'.addcslashes($search, '%_\\').'%'))
            ->when($mode, fn ($q) => $q->where('mode', $mode))
            ->orderByDesc('last_activity_at')
            ->paginate(15)
            ->withQueryString()
            ->through(fn (ChatSession $s) => $this->summary($s));

        $all = $user->chatSessions();

        return Inertia::render('sessions', [
            'sessions' => $sessions,
            'filters' => ['q' => $search, 'mode' => $mode],
            'modes' => ChatPrompts::MODES,
            'stats' => [
                'total' => (clone $all)->count(),
                'messages' => (int) (clone $all)->sum('messages_count'),
                'today' => (clone $all)->where('last_activity_at', '>=', now()->startOfDay())->count(),
                'last_activity_at' => (clone $all)->max('last_activity_at')
                    ? now()->parse((clone $all)->max('last_activity_at'))->toIso8601String()
                    : null,
            ],
        ]);
    }

    public function show(Request $request, ChatSession $session): Response
    {
        abort_unless($session->user_id === $request->user()->id, 404);

        return Inertia::render('sessions/show', [
            'session' => $this->summary($session),
            'messages' => $session->messages()->orderBy('id')->get()->map(fn (ChatMessage $m) => [
                'id' => $m->id,
                'role' => $m->role,
                'content' => $m->content,
                'created_at' => $m->created_at->toIso8601String(),
            ]),
        ]);
    }

    public function destroy(Request $request, ChatSession $session): RedirectResponse
    {
        abort_unless($session->user_id === $request->user()->id, 404);

        $session->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Sesi dihapus.')]);

        return to_route('sessions');
    }

    /**
     * @return array<string, mixed>
     */
    private function summary(ChatSession $session): array
    {
        return [
            'uuid' => $session->uuid,
            'title' => $session->title,
            'mode' => $session->mode,
            'model' => $session->model,
            'messages_count' => $session->messages_count,
            'created_at' => $session->created_at->toIso8601String(),
            'last_activity_at' => $session->last_activity_at->toIso8601String(),
        ];
    }
}

<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\ChatSession;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class SessionController extends Controller
{
    public function index(Request $request): Response
    {
        $search = trim((string) $request->query('q', ''));

        $sessions = ChatSession::query()
            ->with('user:id,name,email')
            ->when($search !== '', fn ($q) => $q->where(
                fn ($q) => $q->where('title', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhereHas('user', fn ($q) => $q
                        ->where('name', 'like', '%'.addcslashes($search, '%_\\').'%')
                        ->orWhere('email', 'like', '%'.addcslashes($search, '%_\\').'%'))
            ))
            ->orderByDesc('last_activity_at')
            ->paginate(15)
            ->withQueryString()
            ->through(fn (ChatSession $s) => [
                'id' => $s->id,
                'uuid' => $s->uuid,
                'title' => $s->title,
                'mode' => $s->mode,
                'messages_count' => $s->messages_count,
                'last_activity_at' => $s->last_activity_at->toIso8601String(),
                'user' => $s->user ? [
                    'id' => $s->user->id,
                    'name' => $s->user->name,
                    'email' => $s->user->email,
                ] : null,
            ]);

        return Inertia::render('spr/sessions', [
            'sessions' => $sessions,
            'filters' => ['q' => $search],
        ]);
    }

    public function destroy(ChatSession $session): RedirectResponse
    {
        $session->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Sesi dihapus.')]);

        return to_route('spr.sessions');
    }
}

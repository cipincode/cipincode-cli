<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class MemberController extends Controller
{
    public function index(Request $request): Response
    {
        $search = trim((string) $request->query('q', ''));
        $status = in_array($request->query('status'), ['aktif', 'trial', 'kedaluwarsa'], true)
            ? $request->query('status')
            : null;

        $members = User::query()
            ->when($search !== '', fn ($q) => $q->where(
                fn ($q) => $q->where('name', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhere('email', 'like', '%'.addcslashes($search, '%_\\').'%')
            ))
            ->when($status === 'kedaluwarsa', fn ($q) => $q
                ->whereNotNull('plan_expires_at')
                ->where('plan_expires_at', '<=', now()))
            ->when($status === 'trial', fn ($q) => $q
                ->where('plan', 'free')
                ->where(fn ($q) => $q->whereNull('plan_expires_at')->orWhere('plan_expires_at', '>', now())))
            ->when($status === 'aktif', fn ($q) => $q
                ->where('plan', '!=', 'free')
                ->where(fn ($q) => $q->whereNull('plan_expires_at')->orWhere('plan_expires_at', '>', now())))
            ->orderByDesc('id')
            ->paginate(15)
            ->withQueryString()
            ->through(fn (User $u) => [
                'id' => $u->id,
                'name' => $u->name,
                'email' => $u->email,
                'plan' => $u->plan,
                'plan_label' => config("usage.plans.{$u->plan}.label", $u->plan),
                'status' => $this->status($u),
                'expires_at' => $u->plan_expires_at?->toIso8601String(),
                'messages_total' => (int) $u->usageEvents()->sum('units'),
            ]);

        return Inertia::render('spr/members', [
            'members' => $members,
            'filters' => ['q' => $search, 'status' => $status],
        ]);
    }

    public function extend(Request $request, User $user): RedirectResponse
    {
        $base = $user->plan_expires_at && $user->plan_expires_at->isFuture()
            ? $user->plan_expires_at
            : now();

        $user->forceFill(['plan_expires_at' => $base->copy()->addDays(14)])->save();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Masa aktif diperpanjang 14 hari.')]);

        return back();
    }

    /**
     * @return array{key: string, label: string}
     */
    private function status(User $user): array
    {
        if ($user->plan_expires_at && $user->plan_expires_at->isPast()) {
            return ['key' => 'kedaluwarsa', 'label' => 'Kedaluwarsa'];
        }

        if ($user->plan === 'free') {
            return ['key' => 'trial', 'label' => 'Trial'];
        }

        return ['key' => 'aktif', 'label' => 'Aktif'];
    }
}

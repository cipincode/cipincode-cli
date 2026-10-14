<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\User;
use App\Support\Plans;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Inertia\Inertia;
use Inertia\Response;

class UserController extends Controller
{
    public function index(Request $request): Response
    {
        $search = trim((string) $request->query('q', ''));
        $plan = in_array($request->query('plan'), Plans::keys(), true)
            ? $request->query('plan')
            : null;

        $users = User::query()
            ->when($search !== '', fn ($q) => $q->where(
                fn ($q) => $q->where('name', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhere('email', 'like', '%'.addcslashes($search, '%_\\').'%')
            ))
            ->when($plan, fn ($q) => $q->where('plan', $plan))
            ->orderByDesc('id')
            ->paginate(15)
            ->withQueryString()
            ->through(fn (User $u) => [
                'id' => $u->id,
                'name' => $u->name,
                'email' => $u->email,
                'plan' => $u->plan,
                'plan_label' => config("usage.plans.{$u->plan}.label", $u->plan),
                'is_admin' => (bool) $u->is_admin,
                'sessions_count' => $u->chatSessions()->count(),
                'messages_total' => (int) $u->usageEvents()->sum('units'),
                'created_at' => $u->created_at->toIso8601String(),
            ]);

        return Inertia::render('spr/users', [
            'users' => $users,
            'filters' => ['q' => $search, 'plan' => $plan],
            'plans' => Plans::keys(),
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $data = $request->validate([
            'name' => ['required', 'string', 'max:255'],
            'email' => ['required', 'string', 'email', 'max:255', 'unique:users,email'],
            'password' => ['required', 'string', 'min:8'],
            'plan' => ['required', 'string', 'in:'.implode(',', Plans::keys())],
            'is_admin' => ['sometimes', 'boolean'],
        ]);

        $user = new User;
        $user->forceFill([
            'name' => $data['name'],
            'email' => $data['email'],
            'password' => Hash::make($data['password']),
            'plan' => $data['plan'],
            'plan_expires_at' => $data['plan'] === 'free'
                ? now()->addDays((int) config('usage.plans.free.trial_days', 14))
                : null,
            'is_admin' => $data['is_admin'] ?? false,
        ])->save();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Pengguna ditambahkan.')]);

        return to_route('spr.users');
    }

    public function update(Request $request, User $user): RedirectResponse
    {
        $data = $request->validate([
            'plan' => ['required', 'string', 'in:'.implode(',', Plans::keys())],
            'is_admin' => ['sometimes', 'boolean'],
        ]);

        $expires = $data['plan'] === 'free'
            ? now()->addDays((int) config('usage.plans.free.trial_days', 14))
            : null;

        $user->forceFill([
            'plan' => $data['plan'],
            'plan_expires_at' => $expires,
        ]);

        if (array_key_exists('is_admin', $data) && $user->id !== $request->user()->id) {
            $user->forceFill(['is_admin' => $data['is_admin']]);
        }

        $user->save();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Pengguna diperbarui.')]);

        return back();
    }

    public function destroy(Request $request, User $user): RedirectResponse
    {
        if ($user->id === $request->user()->id) {
            abort(403, 'Tidak bisa menghapus akun sendiri.');
        }

        $user->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Pengguna dihapus beserta seluruh datanya.')]);

        return to_route('spr.users');
    }
}

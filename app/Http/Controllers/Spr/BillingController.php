<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\User;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class BillingController extends Controller
{
    public function index(Request $request): Response
    {
        $search = trim((string) $request->query('q', ''));

        $paid = User::query()->where('plan', '!=', 'free');
        $revenue = (clone $paid)->get()->sum(
            fn (User $u) => (int) config("usage.plans.{$u->plan}.price", 0)
        );

        $trial = User::query()
            ->where('plan', 'free')
            ->where(function ($q): void {
                $q->whereNull('plan_expires_at')
                    ->orWhere('plan_expires_at', '>', now());
            })
            ->count();

        $expired = User::query()
            ->whereNotNull('plan_expires_at')
            ->where('plan_expires_at', '<=', now())
            ->count();

        $customers = User::query()
            ->where('plan', '!=', 'free')
            ->when($search !== '', fn ($q) => $q->where(
                fn ($q) => $q->where('name', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhere('email', 'like', '%'.addcslashes($search, '%_\\').'%')
            ))
            ->orderByDesc('id')
            ->paginate(15)
            ->withQueryString()
            ->through(fn (User $u) => [
                'id' => $u->id,
                'name' => $u->name,
                'email' => $u->email,
                'plan' => $u->plan,
                'plan_label' => config("usage.plans.{$u->plan}.label", $u->plan),
                'price_label' => 'Rp '.number_format((int) config("usage.plans.{$u->plan}.price", 0), 0, ',', '.'),
                'messages_total' => (int) $u->usageEvents()->sum('units'),
                'created_at' => $u->created_at->toIso8601String(),
            ]);

        return Inertia::render('spr/billing', [
            'stats' => [
                'revenue' => 'Rp '.number_format($revenue, 0, ',', '.'),
                'customers' => (clone $paid)->count(),
                'trial' => $trial,
                'expired' => $expired,
            ],
            'customers_list' => $customers,
            'filters' => ['q' => $search],
            'note' => 'Pembayaran online belum aktif — angka pemasukan adalah estimasi dari paket aktif.',
        ]);
    }
}

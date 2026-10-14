<?php

namespace App\Http\Controllers;

use App\Services\UsageTracker;
use App\Support\Plans;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class BillingController extends Controller
{
    public function index(Request $request, UsageTracker $usage): Response
    {
        $plan = $usage->snapshot($request->user())['plan'];

        return Inertia::render('billing', [
            'plans' => Plans::catalog(),
            'current' => $plan,
            'selected' => in_array($request->query('plan'), Plans::keys(), true) ? $request->query('plan') : null,
        ]);
    }
}

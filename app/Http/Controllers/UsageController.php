<?php

namespace App\Http\Controllers;

use App\Services\UsageTracker;
use App\Support\Plans;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class UsageController extends Controller
{
    public function index(Request $request, UsageTracker $usage): Response
    {
        return Inertia::render('usage', [
            'usage' => $usage->snapshot($request->user()),
            'sessionHours' => (int) config('usage.session_hours'),
            'plans' => Plans::catalog(),
        ]);
    }
}

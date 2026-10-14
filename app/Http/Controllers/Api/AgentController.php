<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\Agent;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class AgentController extends Controller
{
    /**
     * Daftar agent milik pemilik API key, untuk menu /agent di CLI.
     */
    public function index(Request $request): JsonResponse
    {
        return response()->json([
            'agents' => $request->user()->agents()->orderBy('name')->get()->map(fn (Agent $a) => [
                'id' => $a->id,
                'name' => $a->name,
                'description' => $a->description,
                'mode' => $a->mode,
                'effort' => $a->effort,
            ]),
        ]);
    }
}

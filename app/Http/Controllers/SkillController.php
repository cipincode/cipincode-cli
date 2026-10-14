<?php

namespace App\Http\Controllers;

use App\Models\Agent;
use App\Support\Skills;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class SkillController extends Controller
{
    public function index(Request $request): Response
    {
        $installed = $request->user()->agents()->pluck('name')->all();

        return Inertia::render('skills', [
            'skills' => array_map(
                fn (array $s) => [...$s, 'installed' => in_array($s['name'], $installed, true)],
                Skills::all(),
            ),
            'categories' => Skills::CATEGORIES,
            'agentCount' => count($installed),
            'agentLimit' => Agent::MAX_PER_USER,
        ]);
    }

    public function install(Request $request, string $key): RedirectResponse
    {
        $skill = Skills::find($key);
        abort_unless($skill, 404);

        $user = $request->user();

        if ($user->agents()->where('name', $skill['name'])->exists()) {
            Inertia::flash('toast', ['type' => 'info', 'message' => __('Skill sudah ada di Agents.')]);

            return to_route('skills');
        }

        if ($user->agents()->count() >= Agent::MAX_PER_USER) {
            Inertia::flash('toast', ['type' => 'error', 'message' => __('Batas :max agent tercapai. Hapus agent lama terlebih dahulu.', ['max' => Agent::MAX_PER_USER])]);

            return to_route('skills');
        }

        $user->agents()->create([
            'name' => $skill['name'],
            'description' => $skill['description'],
            'instructions' => $skill['instructions'],
            'mode' => $skill['mode'],
            'effort' => $skill['effort'],
        ]);

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Skill ditambahkan ke Agents.')]);

        return to_route('skills');
    }
}

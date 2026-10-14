<?php

namespace App\Http\Controllers;

use App\Models\Agent;
use App\Support\ChatPrompts;
use App\Support\Effort;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;
use Inertia\Inertia;
use Inertia\Response;

class AgentController extends Controller
{
    public function index(Request $request): Response
    {
        return Inertia::render('agents', [
            'agents' => $request->user()->agents()->orderBy('name')->get()->map(fn (Agent $a) => $this->payload($a)),
            'limit' => Agent::MAX_PER_USER,
            'modes' => ChatPrompts::MODES,
            'efforts' => Effort::LEVELS,
        ]);
    }

    public function store(Request $request): RedirectResponse
    {
        $user = $request->user();
        $data = $this->validated($request, null);

        if ($user->agents()->count() >= Agent::MAX_PER_USER) {
            return back()->withErrors(['name' => __('Batas :max agent tercapai. Hapus agent lama terlebih dahulu.', ['max' => Agent::MAX_PER_USER])]);
        }

        $user->agents()->create($data);

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Agent dibuat.')]);

        return to_route('agents');
    }

    public function update(Request $request, Agent $agent): RedirectResponse
    {
        abort_unless($agent->user_id === $request->user()->id, 404);

        $agent->update($this->validated($request, $agent));

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Agent diperbarui.')]);

        return to_route('agents');
    }

    public function destroy(Request $request, Agent $agent): RedirectResponse
    {
        abort_unless($agent->user_id === $request->user()->id, 404);

        $agent->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Agent dihapus.')]);

        return to_route('agents');
    }

    /**
     * @return array<string, string|null>
     */
    private function validated(Request $request, ?Agent $agent): array
    {
        $data = $request->validate([
            'name' => [
                'required', 'string', 'max:50',
                Rule::unique('agents', 'name')->where('user_id', $request->user()->id)->ignore($agent?->id),
            ],
            'description' => ['nullable', 'string', 'max:160'],
            'instructions' => ['required', 'string', 'max:4000'],
            'mode' => ['required', Rule::in(ChatPrompts::MODES)],
            'effort' => ['required', Rule::in(Effort::LEVELS)],
        ]);

        return [...$data, 'name' => trim($data['name']), 'description' => ($data['description'] ?? null) ?: null];
    }

    /**
     * @return array<string, mixed>
     */
    private function payload(Agent $agent): array
    {
        return [
            'id' => $agent->id,
            'name' => $agent->name,
            'description' => $agent->description,
            'instructions' => $agent->instructions,
            'mode' => $agent->mode,
            'effort' => $agent->effort,
            'updated_at' => $agent->updated_at->toIso8601String(),
        ];
    }
}

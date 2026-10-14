<?php

namespace App\Http\Controllers\Spr;

use App\Http\Controllers\Controller;
use App\Models\ContactMessage;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class ContactController extends Controller
{
    public function index(Request $request): Response
    {
        $search = trim((string) $request->query('q', ''));
        $status = in_array($request->query('status'), ContactMessage::STATUSES, true)
            ? $request->query('status')
            : null;

        $messages = ContactMessage::query()
            ->when($search !== '', fn ($q) => $q->where(
                fn ($q) => $q->where('name', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhere('email', 'like', '%'.addcslashes($search, '%_\\').'%')
                    ->orWhere('message', 'like', '%'.addcslashes($search, '%_\\').'%')
            ))
            ->when($status, fn ($q) => $q->where('status', $status))
            ->orderByDesc('id')
            ->paginate(15)
            ->withQueryString();

        return Inertia::render('spr/contacts', [
            'messages' => $messages,
            'filters' => ['q' => $search, 'status' => $status],
            'statuses' => ContactMessage::STATUSES,
            'stats' => [
                'baru' => ContactMessage::query()->where('status', 'baru')->count(),
                'dibaca' => ContactMessage::query()->where('status', 'dibaca')->count(),
                'dibalas' => ContactMessage::query()->where('status', 'dibalas')->count(),
            ],
        ]);
    }

    public function update(Request $request, ContactMessage $message): RedirectResponse
    {
        $data = $request->validate([
            'status' => ['required', 'string', 'in:'.implode(',', ContactMessage::STATUSES)],
        ]);

        $message->update(['status' => $data['status']]);

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Status pesan diperbarui.')]);

        return back();
    }

    public function destroy(ContactMessage $message): RedirectResponse
    {
        $message->delete();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Pesan dihapus.')]);

        return to_route('spr.contacts');
    }
}

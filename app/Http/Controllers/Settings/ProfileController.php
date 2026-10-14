<?php

namespace App\Http\Controllers\Settings;

use App\Http\Controllers\Controller;
use App\Http\Requests\Settings\ProfileDeleteRequest;
use App\Http\Requests\Settings\ProfileUpdateRequest;
use App\Services\AvatarStore;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Inertia\Inertia;
use Inertia\Response;

class ProfileController extends Controller
{
    /**
     * Show the user's profile settings page.
     */
    public function edit(Request $request): Response
    {
        return Inertia::render('settings/profile', [
            'status' => $request->session()->get('status'),
        ]);
    }

    /**
     * Update the user's profile information.
     */
    public function update(ProfileUpdateRequest $request): RedirectResponse
    {
        $request->user()->fill($request->validated());

        if ($request->user()->isDirty('email')) {
            $request->user()->email_verified_at = null;
        }

        $request->user()->save();

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Profile updated.')]);

        return to_route('profile.edit');
    }

    public function avatar(Request $request, AvatarStore $avatars): RedirectResponse
    {
        $request->validate([
            'avatar' => ['required', 'image', 'mimes:jpg,jpeg,png,webp', 'max:5120'],
        ], [
            'avatar.max' => __('Ukuran foto maksimal 5 MB.'),
            'avatar.image' => __('File harus berupa gambar.'),
            'avatar.mimes' => __('Format yang didukung: JPG, PNG, atau WebP.'),
        ]);

        try {
            $avatars->store($request->user(), $request->file('avatar'));
        } catch (\RuntimeException $e) {
            return back()->withErrors(['avatar' => $e->getMessage()]);
        }

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Foto profil diperbarui.')]);

        return to_route('profile.edit');
    }

    public function destroyAvatar(Request $request, AvatarStore $avatars): RedirectResponse
    {
        $avatars->remove($request->user());

        Inertia::flash('toast', ['type' => 'success', 'message' => __('Foto profil dihapus.')]);

        return to_route('profile.edit');
    }

    /**
     * Delete the user's profile.
     */
    public function destroy(ProfileDeleteRequest $request, AvatarStore $avatars): RedirectResponse
    {
        $user = $request->user();
        $avatars->delete($user);

        Auth::logout();

        $user->delete();

        $request->session()->invalidate();
        $request->session()->regenerateToken();

        return redirect('/');
    }
}

<?php

use App\Http\Controllers\AgentController;
use App\Http\Controllers\ApiKeyController;
use App\Http\Controllers\Auth\GoogleAuthController;
use App\Http\Controllers\BillingController;
use App\Http\Controllers\CliAuthController;
use App\Http\Controllers\ContactController;
use App\Http\Controllers\DashboardController;
use App\Http\Controllers\SkillController;
use App\Http\Controllers\SessionController;
use App\Http\Controllers\Spr\BillingController as SprBillingController;
use App\Http\Controllers\Spr\ContactController as SprContactController;
use App\Http\Controllers\Spr\DashboardController as SprDashboardController;
use App\Http\Controllers\Spr\MemberController as SprMemberController;
use App\Http\Controllers\Spr\SessionController as SprSessionController;
use App\Http\Controllers\Spr\UserController as SprUserController;
use App\Http\Controllers\UsageController;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return \Inertia\Inertia::render('cli');
})->name('home');

Route::inertia('docs', 'docs')->name('docs');

Route::inertia('kontak', 'kontak')->name('kontak');
Route::post('kontak', [ContactController::class, 'store'])
    ->middleware('throttle:10,1')
    ->name('kontak.store');

Route::middleware('guest')->group(function () {
    Route::get('auth/google/redirect', [GoogleAuthController::class, 'redirect'])->name('google.redirect');
    Route::get('auth/google/callback', [GoogleAuthController::class, 'callback'])->name('google.callback');
});

Route::get('cli/authorize', [CliAuthController::class, 'show'])->name('cli.authorize');

Route::middleware(['auth'])->group(function () {
    Route::post('cli/authorize', [CliAuthController::class, 'approve'])->name('cli.authorize.approve');

    Route::get('dashboard', [DashboardController::class, 'index'])->name('dashboard');

    Route::inertia('setup', 'setup')->name('setup');
    Route::get('skills', [SkillController::class, 'index'])->name('skills');
    Route::post('skills/{key}/install', [SkillController::class, 'install'])->name('skills.install');
    Route::get('agents', [AgentController::class, 'index'])->name('agents');
    Route::post('agents', [AgentController::class, 'store'])->name('agents.store');
    Route::patch('agents/{agent}', [AgentController::class, 'update'])->name('agents.update');
    Route::delete('agents/{agent}', [AgentController::class, 'destroy'])->name('agents.destroy');
    Route::get('sessions', [SessionController::class, 'index'])->name('sessions');
    Route::get('sessions/{session}', [SessionController::class, 'show'])->name('sessions.show');
    Route::delete('sessions/{session}', [SessionController::class, 'destroy'])->name('sessions.destroy');
    Route::get('api-keys', [ApiKeyController::class, 'index'])->name('api-keys');
    Route::post('api-keys', [ApiKeyController::class, 'store'])->name('api-keys.store');
    Route::patch('api-keys/{apiKey}', [ApiKeyController::class, 'update'])->name('api-keys.update');
    Route::delete('api-keys/{apiKey}', [ApiKeyController::class, 'destroy'])->name('api-keys.destroy');
    Route::get('usage', [UsageController::class, 'index'])->name('usage');
    Route::get('billing', [BillingController::class, 'index'])->name('billing');

    Route::prefix('spr')->name('spr.')->middleware('spr')->group(function () {
        Route::get('dashboard', [SprDashboardController::class, 'index'])->name('dashboard');
        Route::get('users', [SprUserController::class, 'index'])->name('users');
        Route::post('users', [SprUserController::class, 'store'])->name('users.store');
        Route::patch('users/{user}', [SprUserController::class, 'update'])->name('users.update');
        Route::delete('users/{user}', [SprUserController::class, 'destroy'])->name('users.destroy');
        Route::get('sessions', [SprSessionController::class, 'index'])->name('sessions');
        Route::delete('sessions/{session}', [SprSessionController::class, 'destroy'])->name('sessions.destroy');
        Route::get('billing', [SprBillingController::class, 'index'])->name('billing');
        Route::get('members', [SprMemberController::class, 'index'])->name('members');
        Route::patch('members/{user}', [SprMemberController::class, 'extend'])->name('members.extend');
        Route::get('contacts', [SprContactController::class, 'index'])->name('contacts');
        Route::patch('contacts/{message}', [SprContactController::class, 'update'])->name('contacts.update');
        Route::delete('contacts/{message}', [SprContactController::class, 'destroy'])->name('contacts.destroy');
    });
});

require __DIR__.'/settings.php';

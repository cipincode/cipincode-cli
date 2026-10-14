<?php

use App\Http\Controllers\Api\AgentController;
use App\Http\Controllers\Api\ChatController;
use App\Http\Controllers\Api\CliController;
use App\Http\Controllers\Api\UsageController;
use App\Http\Controllers\Api\VersionController;
use Illuminate\Support\Facades\Route;

Route::post('cli/exchange', [CliController::class, 'exchange'])
    ->middleware('throttle:10,1')
    ->name('cli.exchange');

Route::get('cli/version', [VersionController::class, 'show'])
    ->middleware('throttle:60,1')
    ->name('cli.version');

Route::middleware(['throttle:60,1', 'api.key'])->prefix('cli')->group(function () {
    Route::get('me', [CliController::class, 'me'])->name('cli.me');
    Route::post('chat', [ChatController::class, 'store'])->name('cli.chat');
    Route::get('agents', [AgentController::class, 'index'])->name('cli.agents');
    Route::get('usage', [UsageController::class, 'show'])->name('cli.usage');
    Route::post('usage', [UsageController::class, 'consume'])->name('cli.usage.consume');
});

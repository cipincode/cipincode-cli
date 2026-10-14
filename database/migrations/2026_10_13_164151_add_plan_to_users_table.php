<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->string('plan', 16)->default('free')->index();
            $table->timestamp('plan_expires_at')->nullable();
        });

        // Akun yang sudah ada mendapat masa percobaan penuh, dihitung dari sekarang.
        DB::table('users')->update([
            'plan_expires_at' => now()->addDays((int) config('usage.plans.free.trial_days', 14)),
        ]);
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropIndex(['plan']);
            $table->dropColumn(['plan', 'plan_expires_at']);
        });
    }
};

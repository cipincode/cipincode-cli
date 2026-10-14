<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

return new class extends Migration
{
    /**
     * Paket "paid" tunggal dipecah menjadi go, pro, dan max. Pengguna lama jadi Go.
     */
    public function up(): void
    {
        DB::table('users')->where('plan', 'paid')->update(['plan' => 'go']);
    }

    public function down(): void
    {
        DB::table('users')->whereIn('plan', ['go', 'pro', 'max'])->update(['plan' => 'paid']);
    }
};

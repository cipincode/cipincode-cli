<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}" class="scroll-smooth">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>CipinCode CLI — Autonomous AI Agents in Your Terminal</title>
        <link rel="preconnect" href="https://fonts.bunny.net">
        <link href="https://fonts.bunny.net/css?family=poppins:300,400,500,600,700,800|fira-code:400,500&display=swap" rel="stylesheet">
        @vite(['resources/css/app.css'])
        <style>
            body { font-family: 'Poppins', sans-serif; }
            .font-mono { font-family: 'Fira Code', monospace; }
        </style>
    </head>
    <body class="bg-[#fafafa] text-[#111111] antialiased selection:bg-orange-500 selection:text-white overflow-x-hidden">
        
        <!-- Navbar -->
        <header class="sticky top-0 z-50 bg-[#fafafa]/90 backdrop-blur-md border-b border-[#e5e5e5]">
            <div class="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
                <div class="flex items-center gap-3">
                    <img src="/images/logohitam.png" alt="CipinCode Logo" class="h-9 object-contain" />
                </div>
                <nav class="hidden md:flex items-center gap-8 text-sm font-medium text-[#555555]">
                    <a href="#features" class="hover:text-[#111111] transition-colors">Fitur</a>
                    <a href="#workflow" class="hover:text-[#111111] transition-colors">Cara Kerja</a>
                    <a href="#demo" class="hover:text-[#111111] transition-colors">Live Demo</a>
                    <a href="#pricing" class="hover:text-[#111111] transition-colors">Harga & Plan</a>
                    <a href="#faq" class="hover:text-[#111111] transition-colors">FAQ</a>
                </nav>
                <div class="flex items-center gap-4">
                    @auth
                        <a href="{{ route('dashboard') }}" class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#111111] hover:bg-black text-white font-semibold text-sm transition-all shadow-md">
                            Dashboard &rarr;
                        </a>
                    @else
                        <a href="{{ route('login') }}" class="text-sm font-medium text-[#555555] hover:text-[#111111] transition-colors px-3 py-2">
                            Masuk
                        </a>
                        <a href="{{ route('login') }}" class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-orange-600 hover:bg-orange-700 text-white font-semibold text-sm transition-all shadow-lg shadow-orange-600/20">
                            Mulai Gratis &rarr;
                        </a>
                    @endauth
                </div>
            </div>
        </header>

        <!-- Hero Section -->
        <section class="relative max-w-6xl mx-auto px-6 pt-24 pb-20 text-center">
            <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-orange-500/10 border border-orange-500/20 text-orange-600 text-xs font-semibold uppercase tracking-wider mb-8">
                ✨ Next-Gen AI Terminal Agent Platform
            </div>
            <h1 class="text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-[#111111] mb-6 leading-[1.12]">
                Otomasi Coding & Debugging <br />
                <span class="bg-gradient-to-r from-orange-600 to-amber-600 bg-clip-text text-transparent">
                    Langsung di Terminal Anda.
                </span>
            </h1>
            <p class="text-lg sm:text-xl text-[#555555] max-w-2xl mx-auto mb-12 font-normal leading-relaxed">
                Tingkatkan produktivitas dengan agen AI otonom berbasis Bubbletea TUI. Kelola skrip, refaktor kode, dan pantau penggunaan secara real-time.
            </p>

            <!-- Terminal Install Box -->
            <div class="max-w-xl mx-auto mb-12 bg-[#111111] rounded-2xl p-4 shadow-2xl text-left border border-zinc-800">
                <div class="flex items-center justify-between mb-3 px-2">
                    <div class="flex items-center gap-2">
                        <div class="w-3 h-3 rounded-full bg-red-500"></div>
                        <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                        <div class="w-3 h-3 rounded-full bg-green-500"></div>
                    </div>
                    <span class="text-xs text-zinc-400 font-mono">zsh — terminal</span>
                </div>
                <div class="bg-[#1a1a1c] rounded-xl p-4 font-mono text-sm text-zinc-200 flex items-center justify-between border border-zinc-800">
                    <span class="text-orange-400">$</span>
                    <span class="flex-1 px-3 text-left font-medium select-all">curl -fsSL https://cli.cipincode.com/install.sh | bash</span>
                    <button onclick="navigator.clipboard.writeText('curl -fsSL https://cli.cipincode.com/install.sh | bash'); alert('Perintah berhasil disalin!');" class="px-3 py-1.5 text-xs bg-zinc-800 hover:bg-zinc-700 text-zinc-200 rounded-lg transition-colors font-sans font-semibold">
                        Salin
                    </button>
                </div>
            </div>

            <div class="flex flex-wrap items-center justify-center gap-8 text-sm font-semibold text-[#555555]">
                <div class="flex items-center gap-2">✓ Auto-install Go & Bun</div>
                <div class="flex items-center gap-2">✓ Linux & macOS Native</div>
                <div class="flex items-center gap-2">✓ Secure Token OAuth</div>
            </div>
        </section>

        <!-- Stats / Trust Bar -->
        <section class="border-y border-[#e5e5e5] bg-white py-12">
            <div class="max-w-7xl mx-auto px-6 grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
                <div>
                    <div class="text-3xl font-extrabold text-[#111111] mb-1 font-mono">10x</div>
                    <div class="text-sm text-[#666666]">Produktivitas Coding</div>
                </div>
                <div>
                    <div class="text-3xl font-extrabold text-[#111111] mb-1 font-mono">&lt; 1s</div>
                    <div class="text-sm text-[#666666]">Respon Latensi Agent</div>
                </div>
                <div>
                    <div class="text-3xl font-extrabold text-[#111111] mb-1 font-mono">100%</div>
                    <div class="text-sm text-[#666666]">Kontrol Lokal & Sesi</div>
                </div>
                <div>
                    <div class="text-3xl font-extrabold text-[#111111] mb-1 font-mono">24/7</div>
                    <div class="text-sm text-[#666666]">Otonom di Terminal</div>
                </div>
            </div>
        </section>

        <!-- Features Section -->
        <section id="features" class="max-w-7xl mx-auto px-6 py-24">
            <div class="text-center max-w-2xl mx-auto mb-16">
                <span class="text-orange-600 text-xs font-bold uppercase tracking-wider bg-orange-50 px-3 py-1 rounded-full border border-orange-200">Fitur Unggulan</span>
                <h2 class="text-3xl sm:text-4xl font-extrabold text-[#111111] mt-3 mb-4">Dirancang untuk Developer Profesional</h2>
                <p class="text-[#666666]">Semua yang Anda butuhkan untuk mendelegasikan tugas pemrograman langsung dari baris perintah.</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm hover:shadow-md transition-all">
                    <div class="w-12 h-12 rounded-2xl bg-orange-100 flex items-center justify-center text-orange-600 mb-6 font-bold text-xl">⚡</div>
                    <h3 class="text-xl font-bold text-[#111111] mb-3">Interactive Bubbletea TUI</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Antarmuka terminal interaktif dengan modal popup, navigasi keyboard penuh, dan indikator status agent secara real-time.</p>
                </div>

                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm hover:shadow-md transition-all">
                    <div class="w-12 h-12 rounded-2xl bg-orange-100 flex items-center justify-center text-orange-600 mb-6 font-bold text-xl">🤖</div>
                    <h3 class="text-xl font-bold text-[#111111] mb-3">Multi-Mode AI Execution</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Pilih mode operasional: Chat biasa untuk tanya jawab, Coding dengan auto-approve, atau Debugging untuk analisis error mendalam.</p>
                </div>

                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm hover:shadow-md transition-all">
                    <div class="w-12 h-12 rounded-2xl bg-orange-100 flex items-center justify-center text-orange-600 mb-6 font-bold text-xl">🔒</div>
                    <h3 class="text-xl font-bold text-[#111111] mb-3">Secure Token & Dashboard</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Kelola API keys, pantau penggunaan pesan kuota 5 jam & mingguan, serta atur agent kustom melalui dashboard web.</p>
                </div>
            </div>
        </section>

        <!-- Interactive Demo Preview Section -->
        <section id="demo" class="bg-[#111111] text-white py-24 relative overflow-hidden">
            <div class="max-w-7xl mx-auto px-6 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
                <div>
                    <span class="text-orange-400 text-xs font-bold uppercase tracking-wider bg-orange-500/10 px-3 py-1 rounded-full border border-orange-500/20">Live Experience</span>
                    <h2 class="text-3xl sm:text-4xl font-extrabold text-white mt-3 mb-6">Rasakan Kekuatan TUI di Komputer Anda</h2>
                    <p class="text-zinc-400 text-base leading-relaxed mb-8">
                        CipinCode CLI memadukan kecepatan eksekusi Go dengan kecerdasan model AI mutakhir. Jalankan perintah kompleks tanpa harus berpindah jendela dari editor atau terminal Anda.
                    </p>
                    <ul class="space-y-4 text-zinc-300 text-sm font-medium">
                        <li class="flex items-center gap-3">✅ Sinkronisasi riwayat chat otomatis ke cloud</li>
                        <li class="flex items-center gap-3">✅ Dukungan kustom agent & pengaturan effort level</li>
                        <li class="flex items-center gap-3">✅ Integrasi mulus dengan Git dan build tools lokal</li>
                    </ul>
                </div>
                <div class="bg-[#1a1a1c] border border-zinc-800 rounded-3xl p-6 shadow-2xl font-mono text-sm text-zinc-200">
                    <div class="flex items-center justify-between pb-4 mb-4 border-b border-zinc-800">
                        <span class="text-xs text-orange-400 font-semibold">cipincode agent session v1.2</span>
                        <span class="text-xs text-zinc-500">Connected 🟢</span>
                    </div>
                    <div class="space-y-3 text-xs leading-relaxed">
                        <p class="text-zinc-400">> cipincode --mode=coding --effort=high</p>
                        <p class="text-orange-400">⚡ Initializing autonomous agent environment...</p>
                        <p class="text-zinc-300">🔍 Scanning project structure... 14 files indexed.</p>
                        <p class="text-green-400">✔ Found bug in authentication middleware (token expiry check).</p>
                        <p class="text-zinc-300">✍ Applying patch to app/Http/Middleware/VerifyToken.php...</p>
                        <p class="text-green-400 font-bold">🎉 Successfully fixed! All tests passed.</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Workflow Section -->
        <section id="workflow" class="max-w-7xl mx-auto px-6 py-24">
            <div class="text-center max-w-2xl mx-auto mb-16">
                <span class="text-orange-600 text-xs font-bold uppercase tracking-wider bg-orange-50 px-3 py-1 rounded-full border border-orange-200">Cara Kerja</span>
                <h2 class="text-3xl sm:text-4xl font-extrabold text-[#111111] mt-3 mb-4">Mulai Hanya dalam 3 Langkah</h2>
                <p class="text-[#666666]">Proses setup yang sangat cepat dan ramah untuk semua developer.</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm">
                    <div class="text-4xl font-black text-orange-600/30 mb-4 font-mono">01</div>
                    <h3 class="text-lg font-bold text-[#111111] mb-2">Instalasi Satu Baris</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Jalankan skrip instalasi terminal untuk mengunduh binary CLI dan dependensi secara otomatis.</p>
                </div>
                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm">
                    <div class="text-4xl font-black text-orange-600/30 mb-4 font-mono">02</div>
                    <h3 class="text-lg font-bold text-[#111111] mb-2">Otorisasi Akun</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Masuk dan sinkronkan token akses CLI Anda dengan akun web CipinCode yang aman.</p>
                </div>
                <div class="bg-white border border-[#e5e5e5] rounded-3xl p-8 shadow-sm">
                    <div class="text-4xl font-black text-orange-600/30 mb-4 font-mono">03</div>
                    <h3 class="text-lg font-bold text-[#111111] mb-2">Mulai Beraksi</h3>
                    <p class="text-[#666666] text-sm leading-relaxed">Ketik perintah agent dan biarkan AI menyelesaikan tugas coding dan debugging Anda.</p>
                </div>
            </div>
        </section>

        <!-- Pricing / Plans Section -->
        <section id="pricing" class="bg-white border-t border-[#e5e5e5] py-24">
            <div class="max-w-7xl mx-auto px-6">
                <div class="text-center max-w-2xl mx-auto mb-16">
                    <span class="text-orange-600 text-xs font-bold uppercase tracking-wider bg-orange-50 px-3 py-1 rounded-full border border-orange-200">Harga & Plan</span>
                    <h2 class="text-3xl sm:text-4xl font-extrabold text-[#111111] mt-3 mb-4">Pilih Paket Sesuai Kebutuhan</h2>
                    <p class="text-[#666666]">Mulai secara gratis dan tingkatkan seiring berkembangnya skala proyek Anda.</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
                    <div class="bg-[#fafafa] border border-[#e5e5e5] rounded-3xl p-8 shadow-sm">
                        <h3 class="text-xl font-bold text-[#111111] mb-2">Free Plan</h3>
                        <p class="text-[#666666] text-sm mb-6">Cocok untuk mencoba dan eksplorasi fitur dasar CLI.</p>
                        <div class="text-4xl font-extrabold text-[#111111] mb-6">Rp 0 <span class="text-sm font-normal text-[#666666]">/selamanya</span></div>
                        <ul class="space-y-3 text-sm text-[#555555] mb-8">
                            <li>✅ Akses penuh TUI CLI dasar</li>
                            <li>✅ Kuota pesan standar 5 jam</li>
                            <li>✅ Manajemen 2 API key aktif</li>
                        </ul>
                        <a href="{{ route('login') }}" class="block w-full py-3 text-center rounded-xl bg-[#111111] text-white font-semibold text-sm hover:bg-black transition-colors">Daftar Gratis</a>
                    </div>

                    <div class="bg-[#111111] text-white border border-zinc-800 rounded-3xl p-8 shadow-xl relative overflow-hidden">
                        <div class="absolute top-4 right-4 bg-orange-600 text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wider">Popular</div>
                        <h3 class="text-xl font-bold text-white mb-2">Go Plan (Pro)</h3>
                        <p class="text-zinc-400 text-sm mb-6">Untuk developer profesional dengan produktivitas tinggi.</p>
                        <div class="text-4xl font-extrabold text-white mb-6">Rp 99rb <span class="text-sm font-normal text-zinc-400">/bulan</span></div>
                        <ul class="space-y-3 text-sm text-zinc-300 mb-8">
                            <li>✅ Kuota pesan tinggi & prioritas</li>
                            <li>✅ Multi-mode AI & Effort levels</li>
                            <li>✅ Agent kustom tanpa batas</li>
                            <li>✅ API keys & session analytics lengkap</li>
                        </ul>
                        <a href="{{ route('login') }}" class="block w-full py-3 text-center rounded-xl bg-orange-600 text-white font-semibold text-sm hover:bg-orange-700 transition-colors shadow-lg shadow-orange-600/30">Mulai Sekarang</a>
                    </div>
                </div>
            </div>
        </section>

        <!-- FAQ Section -->
        <section id="faq" class="max-w-4xl mx-auto px-6 py-24">
            <div class="text-center max-w-2xl mx-auto mb-16">
                <span class="text-orange-600 text-xs font-bold uppercase tracking-wider bg-orange-50 px-3 py-1 rounded-full border border-orange-200">FAQ</span>
                <h2 class="text-3xl font-extrabold text-[#111111] mt-3 mb-4">Pertanyaan yang Sering Diajukan</h2>
            </div>

            <div class="space-y-6">
                <div class="bg-white border border-[#e5e5e5] rounded-2xl p-6 shadow-sm">
                    <h3 class="font-bold text-[#111111] mb-2">Bagaimana cara menginstal CipinCode CLI?</h3>
                    <p class="text-[#666666] text-sm">Cukup salin perintah instalasi di bagian atas halaman ini dan tempelkan ke terminal Linux atau macOS Anda.</p>
                </div>
                <div class="bg-white border border-[#e5e5e5] rounded-2xl p-6 shadow-sm">
                    <h3 class="font-bold text-[#111111] mb-2">Apakah data kode saya aman?</h3>
                    <p class="text-[#666666] text-sm">Ya! Sesi percakapan dan token terenkripsi aman, dan Anda memiliki kendali penuh atas file yang dimodifikasi agent.</p>
                </div>
                <div class="bg-white border border-[#e5e5e5] rounded-2xl p-6 shadow-sm">
                    <h3 class="font-bold text-[#111111] mb-2">Bagaimana cara memperbarui kuota pesan?</h3>
                    <p class="text-[#666666] text-sm">Anda dapat meningkatkan ke Go Plan melalui dashboard web setelah masuk dengan akun Anda.</p>
                </div>
            </div>
        </section>

        <!-- Footer -->
        <footer class="border-t border-[#e5e5e5] bg-white py-12 text-center text-sm text-[#666666]">
            <div class="max-w-7xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
                <div class="flex items-center gap-2">
                    <img src="/images/logohitam.png" alt="CipinCode" class="h-6 object-contain" />
                </div>
                <p>© {{ date('Y') }} CipinCode. All rights reserved.</p>
            </div>
        </footer>

    </body>
</html>

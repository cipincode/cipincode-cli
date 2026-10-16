import { Head, Link, usePage } from '@inertiajs/react';
import { dashboard, login } from '@/routes';
import { Terminal, Cpu, Zap, Shield, ArrowRight, CheckCircle2, Code2, Sparkles } from 'lucide-react';

export default function Welcome() {
    const { auth } = usePage().props;

    return (
        <>
            <Head title="CipinCode CLI - AI Agent Terminal" />
            <div className="min-h-screen bg-[#09090b] text-zinc-100 selection:bg-orange-500 selection:text-white font-sans antialiased overflow-x-hidden">
                {/* Background Glow */}
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-7xl h-[500px] bg-gradient-to-b from-orange-500/10 via-orange-600/5 to-transparent blur-3xl pointer-events-none" />

                {/* Navbar */}
                <header className="relative z-50 border-b border-zinc-800/80 bg-[#09090b]/80 backdrop-blur-xl">
                    <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <img src="/images/logocode.png" alt="CipinCode Logo" className="h-10 object-contain rounded-xl shadow-lg ring-1 ring-orange-500/30" />
                        </div>
                        <nav className="flex items-center gap-4">
                            <Link
                                href="/setup"
                                className="text-sm font-medium text-zinc-400 hover:text-zinc-100 transition-colors hidden sm:block"
                            >
                                Setup Guide
                            </Link>
                            {auth.user ? (
                                <Link
                                    href={dashboard()}
                                    className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-medium text-sm transition-all shadow-lg shadow-orange-500/25"
                                >
                                    Dashboard <ArrowRight className="w-4 h-4" />
                                </Link>
                            ) : (
                                <>
                                    <Link
                                        href={login()}
                                        className="px-4 py-2 text-sm font-medium text-zinc-300 hover:text-white transition-colors"
                                    >
                                        Log in
                                    </Link>
                                    <Link
                                        href={login()}
                                        className="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-medium text-sm transition-all shadow-lg shadow-orange-500/25"
                                    >
                                        Get Started <ArrowRight className="w-4 h-4" />
                                    </Link>
                                </>
                            )}
                        </nav>
                    </div>
                </header>

                {/* Hero Section */}
                <section className="relative z-10 max-w-5xl mx-auto px-6 pt-24 pb-20 text-center">
                    <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-orange-500/10 border border-orange-500/20 text-orange-400 text-xs font-semibold uppercase tracking-wider mb-8 animate-pulse">
                        <Sparkles className="w-3.5 h-3.5" /> Next-Gen AI Terminal Agent
                    </div>
                    <h1 className="text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-white mb-6 leading-[1.1]">
                        AI Agent di Terminal Anda, <br />
                        <span className="bg-gradient-to-r from-orange-400 via-orange-500 to-amber-500 bg-clip-text text-transparent">
                            Satu Perintah untuk Semua.
                        </span>
                    </h1>
                    <p className="text-lg sm:text-xl text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed">
                        Pasang CipinCode CLI sekarang dan nikmati kekuatan AI otonom untuk coding, debug, dan mengelola proyek langsung dari terminal komputer Anda.
                    </p>

                    {/* Terminal Install Box */}
                    <div className="max-w-xl mx-auto mb-14 bg-zinc-900/90 border border-zinc-800 rounded-2xl p-4 shadow-2xl backdrop-blur-md">
                        <div className="flex items-center justify-between mb-3 px-2">
                            <div className="flex items-center gap-2">
                                <div className="w-3 h-3 rounded-full bg-red-500/80" />
                                <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
                                <div className="w-3 h-3 rounded-full bg-green-500/80" />
                            </div>
                            <span className="text-xs text-zinc-500 font-mono">terminal</span>
                        </div>
                        <div className="bg-[#050507] rounded-xl p-4 font-mono text-sm text-zinc-200 flex items-center justify-between overflow-x-auto border border-zinc-800/50">
                            <span className="text-orange-400">$</span>
                            <span className="flex-1 px-3 text-left selection:bg-orange-500/30">
                                npm install -g cipincode
                            </span>
                            <button
                                onClick={() => navigator.clipboard.writeText("npm install -g cipincode")}
                                className="px-3 py-1.5 text-xs bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg transition-colors"
                            >
                                Copy
                            </button>
                        </div>
                    </div>

                    <div className="flex flex-wrap items-center justify-center gap-6 text-sm text-zinc-400">
                        <div className="flex items-center gap-2">
                            <CheckCircle2 className="w-4 h-4 text-orange-500" /> Auto-install Go & Bun
                        </div>
                        <div className="flex items-center gap-2">
                            <CheckCircle2 className="w-4 h-4 text-orange-500" /> Linux & macOS Support
                        </div>
                        <div className="flex items-center gap-2">
                            <CheckCircle2 className="w-4 h-4 text-orange-500" /> Secure OAuth Login
                        </div>
                    </div>
                </section>

                {/* Features Grid */}
                <section className="relative z-10 max-w-7xl mx-auto px-6 py-20 border-t border-zinc-800/80">
                    <div className="text-center max-w-2xl mx-auto mb-16">
                        <h2 className="text-3xl font-bold text-white mb-4">Dirancang untuk Produktivitas Tinggi</h2>
                        <p className="text-zinc-400">Fitur canggih yang terintegrasi langsung dengan workflow terminal Anda.</p>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                        <div className="bg-zinc-900/50 border border-zinc-800/80 rounded-2xl p-8 hover:border-orange-500/50 transition-all group">
                            <div className="w-12 h-12 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-400 mb-6 group-hover:scale-110 transition-transform">
                                <Terminal className="w-6 h-6" />
                            </div>
                            <h3 className="text-xl font-semibold text-white mb-3">Interactive TUI</h3>
                            <p className="text-zinc-400 text-sm leading-relaxed">
                                Antarmuka terminal interaktif berbasis Bubbletea dengan modal pop-up, indikator checkbox agent, dan riwayat chat otonom.
                            </p>
                        </div>

                        <div className="bg-zinc-900/50 border border-zinc-800/80 rounded-2xl p-8 hover:border-orange-500/50 transition-all group">
                            <div className="w-12 h-12 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-400 mb-6 group-hover:scale-110 transition-transform">
                                <Cpu className="w-6 h-6" />
                            </div>
                            <h3 className="text-xl font-semibold text-white mb-3">Multi-Mode AI</h3>
                            <p className="text-zinc-400 text-sm leading-relaxed">
                                Ganti mode dengan cepat antara Chat biasa, Coding (auto-approve), dan Debugging dengan tingkat effort yang bisa diatur.
                            </p>
                        </div>

                        <div className="bg-zinc-900/50 border border-zinc-800/80 rounded-2xl p-8 hover:border-orange-500/50 transition-all group">
                            <div className="w-12 h-12 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-400 mb-6 group-hover:scale-110 transition-transform">
                                <Shield className="w-6 h-6" />
                            </div>
                            <h3 className="text-xl font-semibold text-white mb-3">Secure & Tracked</h3>
                            <p className="text-zinc-400 text-sm leading-relaxed">
                                Kelola API keys, pantau penggunaan kuota 5 jam & mingguan, serta atur agent kustom langsung dari dashboard web.
                            </p>
                        </div>
                    </div>
                </section>

                {/* Footer */}
                <footer className="border-t border-zinc-800/80 py-12 text-center text-sm text-zinc-500">
                    <div className="max-w-7xl mx-auto px-6 flex flex-col sm:flex-row items-center justify-between gap-4">
                        <div className="flex items-center gap-2">
                            <img src="/images/logocode.png" alt="CipinCode" className="h-6 object-contain" />
                        </div>
                        <p>© {new Date().getFullYear()} CipinCode. All rights reserved.</p>
                    </div>
                </footer>
            </div>
        </>
    );
}

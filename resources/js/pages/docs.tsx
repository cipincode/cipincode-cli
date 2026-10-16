import { useEffect, useState } from 'react';
import { Head } from '@inertiajs/react';
import { CheckIcon, CopyIcon } from 'lucide-react';
import Header from '../../../components/saas-landing/header';
import Footer from '../../../components/saas-landing/footer';
import { cn } from '@/lib/utils';

function CodeBlock({ command }: { command: string }) {
    const [copied, setCopied] = useState(false);

    const copy = async () => {
        try {
            await navigator.clipboard.writeText(command);
        } catch {
            const ta = document.createElement('textarea');
            ta.value = command;
            document.body.appendChild(ta);
            ta.select();
            document.execCommand('copy');
            document.body.removeChild(ta);
        }
        setCopied(true);
        setTimeout(() => setCopied(false), 1600);
    };

    return (
        <div className="mt-3 overflow-hidden rounded-xl border border-border bg-zinc-950">
            <div className="flex items-center justify-between gap-3 px-4 py-3">
                <code className="min-w-0 flex-1 truncate font-mono text-sm text-zinc-100">
                    <span className="mr-2 text-orange-400">$</span>
                    {command}
                </code>
                <button
                    type="button"
                    onClick={copy}
                    className="inline-flex shrink-0 items-center gap-1.5 rounded-lg bg-zinc-800 px-3 py-1.5 text-xs font-medium text-zinc-200 transition-colors hover:bg-zinc-700"
                >
                    {copied ? (
                        <>
                            <CheckIcon className="size-3.5 text-green-400" />
                            Disalin
                        </>
                    ) : (
                        <>
                            <CopyIcon className="size-3.5" />
                            Salin
                        </>
                    )}
                </button>
            </div>
        </div>
    );
}

const sections = [
    { id: 'mulai-cepat', label: 'Mulai Cepat' },
    { id: 'perintah', label: 'Perintah Utama' },
    { id: 'mode-effort', label: 'Mode & Effort' },
    { id: 'kuota-paket', label: 'Kuota & Paket' },
    { id: 'dashboard', label: 'Dashboard Web' },
];

const quotas: Array<{
    name: string;
    price: string;
    session: string;
    weekly: string;
}> = [
    {
        name: 'Gratis',
        price: 'Gratis',
        session: '20 pesan per 5 jam',
        weekly: 'Tanpa batas mingguan',
    },
    {
        name: 'Go',
        price: 'Rp 49.000/bln',
        session: '40 pesan per 5 jam',
        weekly: '400 pesan per minggu',
    },
    {
        name: 'Pro',
        price: 'Rp 149.000/bln',
        session: '100 pesan per 5 jam',
        weekly: '1.500 pesan per minggu',
    },
    {
        name: 'Max',
        price: 'Rp 499.000/bln',
        session: '300 pesan per 5 jam',
        weekly: '5.000 pesan per minggu',
    },
];

export default function Docs() {
    const [active, setActive] = useState(sections[0].id);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                for (const entry of entries) {
                    if (entry.isIntersecting) setActive(entry.target.id);
                }
            },
            { rootMargin: '-25% 0px -65% 0px' },
        );
        sections.forEach(({ id }) => {
            const el = document.getElementById(id);
            if (el) observer.observe(el);
        });
        return () => observer.disconnect();
    }, []);

    return (
        <>
            <Head title="Dokumentasi" />
            <div className="relative isolate flex min-h-screen w-full flex-col bg-background text-foreground">
                <Header />
                <main className="mx-auto w-full max-w-6xl flex-1 px-4 sm:px-6">
                    <div className="border-b border-border py-10 md:py-14">
                        <p className="font-mono text-xs font-semibold tracking-widest text-orange-600 uppercase">
                            Dokumentasi
                        </p>
                        <h1 className="mt-2 max-w-2xl text-3xl font-semibold tracking-tight text-balance sm:text-4xl lg:text-5xl">
                            Mulai pakai CipinCode CLI dalam hitungan menit
                        </h1>
                        <p className="mt-4 max-w-2xl leading-relaxed text-muted-foreground">
                            Panduan instalasi, perintah utama, mode agent,
                            kuota, dan dashboard — semua yang Anda butuhkan
                            untuk coding dengan agen AI di terminal.
                        </p>
                    </div>

                    <nav
                        aria-label="Bagian dokumentasi"
                        className="sticky top-16 z-20 -mx-4 border-b border-border bg-background/90 px-4 py-3 backdrop-blur lg:hidden"
                    >
                        <div className="flex gap-2 overflow-x-auto">
                            {sections.map((s) => (
                                <a
                                    key={s.id}
                                    href={`#${s.id}`}
                                    className={cn(
                                        'shrink-0 rounded-full border px-3.5 py-1.5 text-sm font-medium transition-colors',
                                        active === s.id
                                            ? 'border-orange-600 bg-orange-600 text-white'
                                            : 'border-border text-muted-foreground hover:text-foreground',
                                    )}
                                >
                                    {s.label}
                                </a>
                            ))}
                        </div>
                    </nav>

                    <div className="flex gap-10 py-10">
                        <aside className="hidden w-60 shrink-0 lg:block">
                            <div className="sticky top-24 rounded-2xl border border-border bg-card p-3 shadow-[0_20px_40px_-24px_rgb(0_0_0/0.25)]">
                                <p className="px-3 pt-2 pb-1 font-mono text-[11px] font-semibold tracking-widest text-muted-foreground uppercase">
                                    Daftar Isi
                                </p>
                                <nav
                                    aria-label="Daftar isi"
                                    className="flex flex-col"
                                >
                                    {sections.map((s, i) => (
                                        <a
                                            key={s.id}
                                            href={`#${s.id}`}
                                            className={cn(
                                                'flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors',
                                                active === s.id
                                                    ? 'bg-orange-600/10 text-orange-600'
                                                    : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                                            )}
                                        >
                                            <span
                                                className={cn(
                                                    'font-mono text-xs tabular-nums',
                                                    active === s.id
                                                        ? 'text-orange-600'
                                                        : 'text-muted-foreground/60',
                                                )}
                                            >
                                                {String(i + 1).padStart(
                                                    2,
                                                    '0',
                                                )}
                                            </span>
                                            {s.label}
                                        </a>
                                    ))}
                                </nav>
                            </div>
                        </aside>

                        <article className="min-w-0 flex-1">
                            <section id="mulai-cepat" className="scroll-mt-32">
                                <h2 className="text-2xl font-semibold tracking-tight">
                                    1. Mulai Cepat
                                </h2>
                                <p className="mt-2 leading-relaxed text-muted-foreground">
                                    Tiga langkah untuk memasang CLI dan
                                    menjalankan agent pertama Anda.
                                </p>
                                <ol className="mt-6 flex flex-col gap-4">
                                    <li className="rounded-2xl border border-border bg-card p-5">
                                        <p className="font-medium">
                                            Pasang CipinCode CLI
                                        </p>
                                        <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                                            Satu perintah ini otomatis memasang
                                            Go, Bun, dan binary CipinCode CLI
                                            di Linux &amp; macOS.
                                        </p>
                                        <CodeBlock command="npm install -g cipincode" />
                                    </li>
                                    <li className="rounded-2xl border border-border bg-card p-5">
                                        <p className="font-medium">
                                            Verifikasi instalasi
                                        </p>
                                        <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                                            Jika nomor versi tampil, berarti
                                            CLI siap dipakai.
                                        </p>
                                        <CodeBlock command="cipincode --version" />
                                    </li>
                                    <li className="rounded-2xl border border-border bg-card p-5">
                                        <p className="font-medium">
                                            Masuk dan mulai
                                        </p>
                                        <p className="mt-1 text-sm leading-relaxed text-muted-foreground">
                                            Jalankan perintah di bawah, pilih
                                            masuk dengan Email atau Google,
                                            buka URL yang tampil, klik Izinkan,
                                            lalu tempel kode yang muncul ke
                                            terminal.
                                        </p>
                                        <CodeBlock command="cipincode login" />
                                    </li>
                                </ol>
                            </section>

                            <section
                                id="perintah"
                                className="mt-14 scroll-mt-32"
                            >
                                <h2 className="text-2xl font-semibold tracking-tight">
                                    2. Perintah Utama
                                </h2>
                                <p className="mt-2 leading-relaxed text-muted-foreground">
                                    Perintah yang paling sering dipakai
                                    sehari-hari.
                                </p>
                                <div className="mt-6 overflow-hidden rounded-2xl border border-border">
                                    <table className="w-full text-left text-sm">
                                        <thead>
                                            <tr className="border-b border-border bg-muted/40">
                                                <th className="px-4 py-3 font-mono text-xs font-semibold tracking-wide uppercase">
                                                    Perintah
                                                </th>
                                                <th className="px-4 py-3 font-mono text-xs font-semibold tracking-wide uppercase">
                                                    Fungsi
                                                </th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-border">
                                            <tr>
                                                <td className="px-4 py-3 font-mono text-[13px]">
                                                    cipincode --version
                                                </td>
                                                <td className="px-4 py-3 text-muted-foreground">
                                                    Cek versi CLI yang
                                                    terpasang
                                                </td>
                                            </tr>
                                            <tr>
                                                <td className="px-4 py-3 font-mono text-[13px]">
                                                    cipincode login
                                                </td>
                                                <td className="px-4 py-3 text-muted-foreground">
                                                    Masuk dan hubungkan akun
                                                    Anda
                                                </td>
                                            </tr>
                                            <tr>
                                                <td className="px-4 py-3 font-mono text-[13px]">
                                                    /
                                                </td>
                                                <td className="px-4 py-3 text-muted-foreground">
                                                    Ketik garis miring di dalam
                                                    CLI untuk melihat semua
                                                    perintah
                                                </td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                            </section>

                            <section
                                id="mode-effort"
                                className="mt-14 scroll-mt-32"
                            >
                                <h2 className="text-2xl font-semibold tracking-tight">
                                    3. Mode &amp; Effort
                                </h2>
                                <p className="mt-2 leading-relaxed text-muted-foreground">
                                    Atur cara agent bekerja sesuai tugas Anda.
                                </p>
                                <div className="mt-6 grid gap-4 sm:grid-cols-3">
                                    {[
                                        {
                                            name: 'Chat',
                                            desc: 'Tanya jawab biasa seputar kode tanpa mengubah file.',
                                        },
                                        {
                                            name: 'Coding',
                                            desc: 'Menulis dan mengubah kode dengan persetujuan otomatis.',
                                        },
                                        {
                                            name: 'Debugging',
                                            desc: 'Analisis error mendalam sampai ketemu akar masalahnya.',
                                        },
                                    ].map((m) => (
                                        <div
                                            key={m.name}
                                            className="rounded-2xl border border-border bg-card p-5"
                                        >
                                            <p className="font-mono text-sm font-semibold text-orange-600">
                                                {m.name}
                                            </p>
                                            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                                                {m.desc}
                                            </p>
                                        </div>
                                    ))}
                                </div>
                                <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
                                    Setiap mode punya tingkat effort yang bisa
                                    diatur — makin tinggi effort, makin dalam
                                    agent menganalisis sebelum bertindak.
                                </p>
                            </section>

                            <section
                                id="kuota-paket"
                                className="mt-14 scroll-mt-32"
                            >
                                <h2 className="text-2xl font-semibold tracking-tight">
                                    4. Kuota &amp; Paket
                                </h2>
                                <p className="mt-2 leading-relaxed text-muted-foreground">
                                    Kuota dihitung per 5 jam dan per minggu.
                                    Paket Gratis punya masa percobaan 14 hari —
                                    paket berbayar segera hadir.
                                </p>
                                <div className="mt-6 overflow-hidden rounded-2xl border border-border">
                                    <table className="w-full text-left text-sm">
                                        <thead>
                                            <tr className="border-b border-border bg-muted/40">
                                                <th className="px-4 py-3 font-mono text-xs font-semibold tracking-wide uppercase">
                                                    Paket
                                                </th>
                                                <th className="px-4 py-3 font-mono text-xs font-semibold tracking-wide uppercase">
                                                    Per 5 jam
                                                </th>
                                                <th className="px-4 py-3 font-mono text-xs font-semibold tracking-wide uppercase">
                                                    Per minggu
                                                </th>
                                            </tr>
                                        </thead>
                                        <tbody className="divide-y divide-border">
                                            {quotas.map((q) => (
                                                <tr key={q.name}>
                                                    <td className="px-4 py-3 font-medium">
                                                        {q.name}
                                                        <span className="ml-2 text-xs font-normal text-muted-foreground">
                                                            {q.price}
                                                        </span>
                                                    </td>
                                                    <td className="px-4 py-3 text-muted-foreground">
                                                        {q.session}
                                                    </td>
                                                    <td className="px-4 py-3 text-muted-foreground">
                                                        {q.weekly}
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                                <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
                                    Pembayaran online belum aktif.{' '}
                                    <a
                                        href="/#harga"
                                        className="font-medium text-orange-600 hover:underline"
                                    >
                                        Lihat paket di halaman utama
                                    </a>
                                    .
                                </p>
                            </section>

                            <section
                                id="dashboard"
                                className="mt-14 scroll-mt-32"
                            >
                                <h2 className="text-2xl font-semibold tracking-tight">
                                    5. Dashboard Web
                                </h2>
                                <p className="mt-2 leading-relaxed text-muted-foreground">
                                    Setelah masuk, kelola semuanya dari browser:
                                    buat dan cabut API key, pantau riwayat
                                    sesi, lihat pemakaian kuota, dan atur
                                    agent kustom. Semua tersinkron otomatis
                                    dengan CLI di terminal Anda.
                                </p>
                                <a
                                    href="/login"
                                    className="mt-5 inline-flex h-10 items-center rounded-xl bg-orange-600 px-5 text-sm font-semibold text-white transition-colors hover:bg-orange-700"
                                >
                                    Buka Dashboard
                                </a>
                            </section>
                        </article>
                    </div>
                </main>
                <div className="border-t border-border/50">
                    <Footer />
                </div>
            </div>
        </>
    );
}

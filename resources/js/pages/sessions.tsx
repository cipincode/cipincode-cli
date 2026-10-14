import { Head, Link, router } from '@inertiajs/react';
import {
    Clock,
    MessagesSquare,
    Search,
    Sun,
    TerminalSquare,
} from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { formatDateTime, formatRelative } from '@/components/api-keys/format';
import StatCard from '@/components/api-keys/stat-card';
import { ModeBadge } from '@/components/mode-badge';
import DeleteSessionDialog from '@/components/sessions/delete-session-dialog';
import { sessions as sessionsRoute, setup } from '@/routes';
import { show } from '@/routes/sessions';

type Session = {
    uuid: string;
    title: string;
    mode: string;
    model: string | null;
    messages_count: number;
    created_at: string;
    last_activity_at: string;
};

type Paginated = {
    data: Session[];
    from: number | null;
    to: number | null;
    total: number;
    prev_page_url: string | null;
    next_page_url: string | null;
};

type Props = {
    sessions: Paginated;
    filters: { q: string; mode: string | null };
    modes: string[];
    stats: {
        total: number;
        messages: number;
        today: number;
        last_activity_at: string | null;
    };
};

const th =
    'px-5 py-3 text-left text-xs font-medium tracking-wide text-muted-foreground uppercase whitespace-nowrap';
const MODE_LABEL: Record<string, string> = {
    chat: 'Chat',
    coding: 'Coding',
    debugging: 'Debugging',
};

export default function Sessions({ sessions, filters, modes, stats }: Props) {
    const [search, setSearch] = useState(filters.q);
    const [mode, setMode] = useState<string | null>(filters.mode);
    const [deleting, setDeleting] = useState<Session | null>(null);
    const first = useRef(true);

    // Pencarian dan filter dijalankan di server, ditunda sebentar saat mengetik.
    useEffect(() => {
        if (first.current) {
            first.current = false;

            return;
        }
        const timer = setTimeout(() => {
            router.get(
                sessionsRoute().url,
                { q: search || undefined, mode: mode ?? undefined },
                { preserveState: true, preserveScroll: true, replace: true },
            );
        }, 300);

        return () => clearTimeout(timer);
    }, [search, mode]);

    const filtering = filters.q !== '' || filters.mode !== null;

    return (
        <>
            <Head title="Sessions" />

            <div className="flex flex-1 flex-col">
                <div className="border-b border-sidebar-border px-6 py-5">
                    <h1 className="text-xl font-semibold tracking-tight">
                        Sessions
                    </h1>
                    <p className="mt-1 text-sm text-muted-foreground">
                        Riwayat percakapan Anda di CLI. Setiap /new memulai sesi
                        baru.
                    </p>
                </div>

                <div className="flex w-full flex-1 flex-col gap-6 p-6">
                    <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
                        <StatCard
                            label="Total sesi"
                            value={stats.total}
                            caption="semua percakapan"
                            icon={TerminalSquare}
                        />
                        <StatCard
                            label="Total pesan"
                            value={stats.messages}
                            caption="pertanyaan dan jawaban"
                            icon={MessagesSquare}
                            tone="brand"
                        />
                        <StatCard
                            label="Aktif hari ini"
                            value={stats.today}
                            caption="sesi dengan aktivitas"
                            icon={Sun}
                            tone="success"
                        />
                        <StatCard
                            label="Aktivitas terakhir"
                            value={
                                stats.last_activity_at
                                    ? formatRelative(stats.last_activity_at)
                                    : '-'
                            }
                            caption={
                                stats.last_activity_at
                                    ? formatDateTime(stats.last_activity_at)
                                    : 'belum ada aktivitas'
                            }
                            icon={Clock}
                            tone="muted"
                        />
                    </div>

                    <section className="flex flex-1 flex-col border border-sidebar-border bg-card">
                        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-sidebar-border p-4">
                            <div className="relative w-full sm:max-w-sm">
                                <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                                <input
                                    type="search"
                                    value={search}
                                    onChange={(e) => setSearch(e.target.value)}
                                    placeholder="Cari judul sesi…"
                                    aria-label="Cari judul sesi"
                                    className="h-9 w-full border border-input bg-transparent pr-3 pl-9 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
                                />
                            </div>
                            <div
                                role="group"
                                aria-label="Filter mode"
                                className="inline-flex border border-sidebar-border"
                            >
                                {[null, ...modes].map((m) => (
                                    <button
                                        key={m ?? 'all'}
                                        type="button"
                                        aria-pressed={mode === m}
                                        onClick={() => setMode(m)}
                                        className={`h-9 px-4 text-sm font-medium transition-colors not-first:border-l not-first:border-sidebar-border ${
                                            mode === m
                                                ? 'bg-orange-600 text-white'
                                                : 'text-muted-foreground hover:bg-sidebar-accent hover:text-foreground'
                                        }`}
                                    >
                                        {m ? MODE_LABEL[m] : 'Semua'}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {sessions.total === 0 ? (
                            <div className="flex flex-1 flex-col items-center justify-center gap-3 px-6 py-16 text-center">
                                <div className="flex size-12 items-center justify-center border border-sidebar-border text-orange-600">
                                    <TerminalSquare className="size-6" />
                                </div>
                                <h2 className="font-medium">
                                    {filtering
                                        ? 'Tidak ada sesi yang cocok'
                                        : 'Belum ada sesi'}
                                </h2>
                                <p className="max-w-md text-sm text-muted-foreground">
                                    {filtering
                                        ? 'Ubah kata kunci pencarian atau filter mode.'
                                        : 'Sesi tercatat otomatis setiap kali Anda memakai cipincode setelah login. Mulai dari panduan instalasi.'}
                                </p>
                                {!filtering && (
                                    <Link
                                        href={setup()}
                                        className="mt-1 inline-flex h-9 items-center bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700"
                                    >
                                        Buka panduan instalasi
                                    </Link>
                                )}
                            </div>
                        ) : (
                            <div className="overflow-x-auto">
                                <table className="w-full min-w-[760px] text-sm">
                                    <thead className="border-b border-sidebar-border bg-muted/40">
                                        <tr>
                                            <th className={th}>Judul</th>
                                            <th className={th}>Mode</th>
                                            <th className={`${th} text-right`}>
                                                Pesan
                                            </th>
                                            <th className={th}>Model</th>
                                            <th className={th}>
                                                Aktivitas terakhir
                                            </th>
                                            <th
                                                className={`${th} w-40 text-right`}
                                            >
                                                <span className="sr-only">
                                                    Aksi
                                                </span>
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-sidebar-border">
                                        {sessions.data.map((s) => (
                                            <tr
                                                key={s.uuid}
                                                className="transition-colors hover:bg-muted/30"
                                            >
                                                <td className="max-w-md px-5 py-3.5">
                                                    <Link
                                                        href={show(s.uuid)}
                                                        className="block truncate font-medium hover:text-orange-600"
                                                    >
                                                        {s.title}
                                                    </Link>
                                                    <span className="text-xs text-muted-foreground">
                                                        Dimulai{' '}
                                                        {formatDateTime(
                                                            s.created_at,
                                                        )}
                                                    </span>
                                                </td>
                                                <td className="px-5 py-3.5">
                                                    <ModeBadge mode={s.mode} />
                                                </td>
                                                <td className="px-5 py-3.5 text-right tabular-nums">
                                                    {s.messages_count}
                                                </td>
                                                <td
                                                    className="max-w-48 truncate px-5 py-3.5 font-mono text-xs text-muted-foreground"
                                                    title={s.model ?? undefined}
                                                >
                                                    {s.model ?? '-'}
                                                </td>
                                                <td
                                                    className="px-5 py-3.5 whitespace-nowrap text-muted-foreground"
                                                    title={formatDateTime(
                                                        s.last_activity_at,
                                                    )}
                                                >
                                                    {formatRelative(
                                                        s.last_activity_at,
                                                    )}
                                                </td>
                                                <td className="px-5 py-3.5 text-right whitespace-nowrap">
                                                    <Link
                                                        href={show(s.uuid)}
                                                        className="mr-2 inline-flex h-8 items-center border border-sidebar-border px-3 text-xs font-medium transition-colors hover:bg-sidebar-accent"
                                                    >
                                                        Lihat
                                                    </Link>
                                                    <button
                                                        type="button"
                                                        onClick={() =>
                                                            setDeleting(s)
                                                        }
                                                        className="inline-flex h-8 items-center border border-sidebar-border px-3 text-xs font-medium transition-colors hover:border-red-600 hover:text-red-600"
                                                    >
                                                        Hapus
                                                    </button>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}

                        {sessions.total > 0 && (
                            <div className="mt-auto flex flex-wrap items-center justify-between gap-3 border-t border-sidebar-border px-5 py-3 text-xs text-muted-foreground">
                                <span>
                                    Menampilkan {sessions.from}–{sessions.to}{' '}
                                    dari {sessions.total} sesi.
                                </span>
                                <div className="flex gap-2">
                                    {[
                                        {
                                            label: 'Sebelumnya',
                                            url: sessions.prev_page_url,
                                        },
                                        {
                                            label: 'Berikutnya',
                                            url: sessions.next_page_url,
                                        },
                                    ].map((p) =>
                                        p.url ? (
                                            <Link
                                                key={p.label}
                                                href={p.url}
                                                preserveScroll
                                                className="inline-flex h-8 items-center border border-sidebar-border px-3 font-medium text-foreground transition-colors hover:bg-sidebar-accent"
                                            >
                                                {p.label}
                                            </Link>
                                        ) : (
                                            <span
                                                key={p.label}
                                                className="inline-flex h-8 items-center border border-sidebar-border px-3 opacity-40"
                                            >
                                                {p.label}
                                            </span>
                                        ),
                                    )}
                                </div>
                            </div>
                        )}
                    </section>
                </div>
            </div>

            <DeleteSessionDialog
                session={deleting}
                onClose={() => setDeleting(null)}
            />
        </>
    );
}

Sessions.layout = {
    breadcrumbs: [{ title: 'Sessions', href: sessionsRoute() }],
};

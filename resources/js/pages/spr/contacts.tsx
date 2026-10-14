import { Head, router } from '@inertiajs/react';
import { useState } from 'react';
import { Inbox, Search, Trash2, Eye, Reply } from 'lucide-react';
import { formatRelative } from '@/components/api-keys/format';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import StatCard from '@/components/api-keys/stat-card';

type Message = {
    id: number;
    name: string;
    email: string;
    message: string;
    status: string;
    created_at: string;
};

type Paginated = {
    data: Message[];
    current_page: number;
    last_page: number;
    total: number;
};

type Props = {
    messages: Paginated;
    filters: { q: string; status: string | null };
    statuses: string[];
    stats: { baru: number; dibaca: number; dibalas: number };
};

const statusLabel: Record<string, string> = {
    baru: 'Baru',
    dibaca: 'Dibaca',
    dibalas: 'Dibalas',
};

export default function SprContacts({
    messages,
    filters,
    statuses,
    stats,
}: Props) {
    const [q, setQ] = useState(filters.q);
    const [status, setStatus] = useState(filters.status ?? '');
    const [expanded, setExpanded] = useState<number | null>(null);

    const search = (page = 1, nextStatus = status) => {
        router.get(
            '/spr/contacts',
            {
                q: q || undefined,
                status: nextStatus || undefined,
                page,
            },
            { preserveState: true, replace: true },
        );
    };

    const setStatusOf = (id: number, next: string) => {
        router.patch(
            `/spr/contacts/${id}`,
            { status: next },
            { preserveScroll: true },
        );
    };

    const destroy = (id: number, name: string) => {
        if (!window.confirm(`Hapus pesan dari ${name}?`)) return;
        router.delete(`/spr/contacts/${id}`, { preserveScroll: true });
    };

    return (
        <>
            <Head title="Contact" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div>
                    <h1 className="text-xl font-semibold tracking-tight">
                        Contact
                    </h1>
                    <p className="text-sm text-muted-foreground">
                        Kotak masuk komplain dan pertanyaan dari halaman
                        kontak publik.
                    </p>
                </div>

                <div className="grid gap-4 sm:grid-cols-3">
                    <StatCard
                        label="Baru masuk"
                        value={stats.baru}
                        caption="Belum dibaca"
                        icon={Inbox}
                        tone={stats.baru ? 'brand' : 'muted'}
                    />
                    <StatCard
                        label="Dibaca"
                        value={stats.dibaca}
                        caption="Menunggu balasan"
                        icon={Eye}
                    />
                    <StatCard
                        label="Dibalas"
                        value={stats.dibalas}
                        caption="Selesai ditangani"
                        icon={Reply}
                        tone="success"
                    />
                </div>

                <form
                    className="flex flex-wrap gap-2"
                    onSubmit={(e) => {
                        e.preventDefault();
                        search(1);
                    }}
                >
                    <div className="relative min-w-52 flex-1">
                        <Search className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            value={q}
                            onChange={(e) => setQ(e.target.value)}
                            placeholder="Cari nama, email, atau isi…"
                            className="pl-9"
                        />
                    </div>
                    <select
                        value={status}
                        onChange={(e) => {
                            setStatus(e.target.value);
                            router.get(
                                '/spr/contacts',
                                {
                                    q: q || undefined,
                                    status: e.target.value || undefined,
                                    page: 1,
                                },
                                { preserveState: true, replace: true },
                            );
                        }}
                        className="h-9 border border-sidebar-border bg-background px-2 text-sm"
                        aria-label="Filter status"
                    >
                        <option value="">Semua status</option>
                        {statuses.map((s) => (
                            <option key={s} value={s}>
                                {statusLabel[s] ?? s}
                            </option>
                        ))}
                    </select>
                    <Button type="submit" variant="outline">
                        Cari
                    </Button>
                </form>

                <div className="flex flex-col gap-3">
                    {messages.data.map((m) => {
                        const open = expanded === m.id;
                        return (
                            <article
                                key={m.id}
                                className="border border-sidebar-border bg-card"
                            >
                                <button
                                    type="button"
                                    onClick={() =>
                                        setExpanded(open ? null : m.id)
                                    }
                                    className="flex w-full flex-wrap items-center gap-3 px-5 py-4 text-left"
                                >
                                    <span className="min-w-0 flex-1">
                                        <span className="block truncate text-sm font-medium">
                                            {m.name}
                                            <span className="ml-2 font-normal text-muted-foreground">
                                                {m.email}
                                            </span>
                                        </span>
                                        <span
                                            className={`mt-0.5 block truncate text-sm text-muted-foreground ${open ? '' : ''}`}
                                        >
                                            {open
                                                ? formatRelative(m.created_at)
                                                : m.message}
                                        </span>
                                    </span>
                                    <Badge
                                        variant={
                                            m.status === 'baru'
                                                ? 'default'
                                                : 'outline'
                                        }
                                        className={
                                            m.status === 'baru'
                                                ? 'bg-orange-600 text-white'
                                                : ''
                                        }
                                    >
                                        {statusLabel[m.status] ?? m.status}
                                    </Badge>
                                </button>
                                {open && (
                                    <div className="border-t border-sidebar-border px-5 py-4">
                                        <p className="text-sm leading-relaxed whitespace-pre-wrap">
                                            {m.message}
                                        </p>
                                        <div className="mt-4 flex flex-wrap gap-2">
                                            <select
                                                value={m.status}
                                                onChange={(e) =>
                                                    setStatusOf(
                                                        m.id,
                                                        e.target.value,
                                                    )
                                                }
                                                className="h-9 border border-sidebar-border bg-background px-2 text-sm"
                                                aria-label="Ubah status"
                                            >
                                                {statuses.map((s) => (
                                                    <option key={s} value={s}>
                                                        {statusLabel[s] ?? s}
                                                    </option>
                                                ))}
                                            </select>
                                            <Button
                                                size="sm"
                                                variant="outline"
                                                className="h-9"
                                                onClick={() =>
                                                    destroy(m.id, m.name)
                                                }
                                            >
                                                <Trash2 className="size-3.5 text-red-600" />
                                                Hapus
                                            </Button>
                                        </div>
                                    </div>
                                )}
                            </article>
                        );
                    })}
                    {messages.data.length === 0 && (
                        <p className="border border-sidebar-border bg-card px-4 py-6 text-center text-sm text-muted-foreground">
                            Kotak masuk kosong.
                        </p>
                    )}
                </div>

                <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        Halaman {messages.current_page} dari{' '}
                        {messages.last_page} · Total {messages.total}
                    </span>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={messages.current_page <= 1}
                            onClick={() => search(messages.current_page - 1)}
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={
                                messages.current_page >= messages.last_page
                            }
                            onClick={() => search(messages.current_page + 1)}
                        >
                            Berikutnya
                        </Button>
                    </div>
                </div>
            </div>
        </>
    );
}

SprContacts.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
        {
            title: 'Contact',
            href: '/spr/contacts',
        },
    ],
};

import { Head, router } from '@inertiajs/react';
import { useState } from 'react';
import { Search, Trash2 } from 'lucide-react';
import { formatRelative } from '@/components/api-keys/format';
import { ModeBadge } from '@/components/mode-badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

type SprSession = {
    id: number;
    uuid: string;
    title: string;
    mode: string;
    messages_count: number;
    last_activity_at: string;
    user: { id: number; name: string; email: string } | null;
};

type Paginated = {
    data: SprSession[];
    current_page: number;
    last_page: number;
    total: number;
};

type Props = {
    sessions: Paginated;
    filters: { q: string };
};

export default function SprSessions({ sessions, filters }: Props) {
    const [q, setQ] = useState(filters.q);

    const search = (page = 1) => {
        router.get(
            '/spr/sessions',
            { q: q || undefined, page },
            { preserveState: true, replace: true },
        );
    };

    const destroy = (id: number, title: string) => {
        if (!window.confirm(`Hapus sesi "${title}"?`)) return;
        router.delete(`/spr/sessions/${id}`, { preserveScroll: true });
    };

    return (
        <>
            <Head title="Sesi" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div>
                    <h1 className="text-xl font-semibold tracking-tight">
                        Sesi
                    </h1>
                    <p className="text-sm text-muted-foreground">
                        {sessions.total} percakapan dari semua pengguna.
                    </p>
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
                            placeholder="Cari judul, nama, atau email…"
                            className="pl-9"
                        />
                    </div>
                    <Button type="submit" variant="outline">
                        Cari
                    </Button>
                </form>

                <div className="overflow-x-auto border border-sidebar-border bg-card">
                    <table className="w-full min-w-180 text-left text-sm">
                        <thead>
                            <tr className="border-b border-sidebar-border bg-muted/40 font-mono text-xs tracking-wide uppercase">
                                <th className="px-4 py-3 font-semibold">
                                    Sesi
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Pengguna
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Mode
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Pesan
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Aktivitas
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Aksi
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {sessions.data.map((s) => (
                                <tr
                                    key={s.id}
                                    className="border-t border-sidebar-border"
                                >
                                    <td className="max-w-56 truncate px-4 py-3 font-medium">
                                        {s.title}
                                    </td>
                                    <td className="px-4 py-3">
                                        {s.user ? (
                                            <>
                                                <p>{s.user.name}</p>
                                                <p className="text-xs text-muted-foreground">
                                                    {s.user.email}
                                                </p>
                                            </>
                                        ) : (
                                            <span className="text-muted-foreground">
                                                —
                                            </span>
                                        )}
                                    </td>
                                    <td className="px-4 py-3">
                                        <ModeBadge mode={s.mode} />
                                    </td>
                                    <td className="px-4 py-3 text-right tabular-nums">
                                        {s.messages_count}
                                    </td>
                                    <td className="px-4 py-3 text-right text-xs text-muted-foreground">
                                        {formatRelative(s.last_activity_at)}
                                    </td>
                                    <td className="px-4 py-3 text-right">
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() =>
                                                destroy(s.id, s.title)
                                            }
                                        >
                                            <Trash2 className="size-3.5 text-red-600" />
                                            Hapus
                                        </Button>
                                    </td>
                                </tr>
                            ))}
                            {sessions.data.length === 0 && (
                                <tr>
                                    <td
                                        colSpan={6}
                                        className="px-4 py-6 text-center text-muted-foreground"
                                    >
                                        Tidak ada sesi yang cocok.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        Halaman {sessions.current_page} dari{' '}
                        {sessions.last_page} · Total {sessions.total}
                    </span>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={sessions.current_page <= 1}
                            onClick={() =>
                                search(sessions.current_page - 1)
                            }
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={
                                sessions.current_page >= sessions.last_page
                            }
                            onClick={() =>
                                search(sessions.current_page + 1)
                            }
                        >
                            Berikutnya
                        </Button>
                    </div>
                </div>
            </div>
        </>
    );
}

SprSessions.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
        {
            title: 'Sesi',
            href: '/spr/sessions',
        },
    ],
};

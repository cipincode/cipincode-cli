import { Head, router } from '@inertiajs/react';
import { useState } from 'react';
import { CalendarPlus, Search } from 'lucide-react';
import { formatRelative } from '@/components/api-keys/format';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

type Member = {
    id: number;
    name: string;
    email: string;
    plan: string;
    plan_label: string;
    status: { key: string; label: string };
    expires_at: string | null;
    messages_total: number;
};

type Paginated = {
    data: Member[];
    current_page: number;
    last_page: number;
    total: number;
};

type Props = {
    members: Paginated;
    filters: { q: string; status: string | null };
};

const statuses = [
    { key: '', label: 'Semua status' },
    { key: 'aktif', label: 'Aktif' },
    { key: 'trial', label: 'Trial' },
    { key: 'kedaluwarsa', label: 'Kedaluwarsa' },
];

export default function SprMembers({ members, filters }: Props) {
    const [q, setQ] = useState(filters.q);
    const [status, setStatus] = useState(filters.status ?? '');

    const search = (page = 1, nextStatus = status) => {
        router.get(
            '/spr/members',
            {
                q: q || undefined,
                status: nextStatus || undefined,
                page,
            },
            { preserveState: true, replace: true },
        );
    };

    const extend = (id: number, name: string) => {
        if (!window.confirm(`Perpanjang masa aktif ${name} 14 hari?`)) return;
        router.patch(`/spr/members/${id}`, {}, { preserveScroll: true });
    };

    return (
        <>
            <Head title="Member" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div>
                    <h1 className="text-xl font-semibold tracking-tight">
                        Member
                    </h1>
                    <p className="text-sm text-muted-foreground">
                        {members.total} keanggotaan. Pantau status langganan
                        dan perpanjang masa aktif.
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
                            placeholder="Cari nama atau email…"
                            className="pl-9"
                        />
                    </div>
                    <select
                        value={status}
                        onChange={(e) => {
                            setStatus(e.target.value);
                            router.get(
                                '/spr/members',
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
                        {statuses.map((s) => (
                            <option key={s.key} value={s.key}>
                                {s.label}
                            </option>
                        ))}
                    </select>
                    <Button type="submit" variant="outline">
                        Cari
                    </Button>
                </form>

                <div className="overflow-x-auto border border-sidebar-border bg-card">
                    <table className="w-full min-w-200 text-left text-sm">
                        <thead>
                            <tr className="border-b border-sidebar-border bg-muted/40 font-mono text-xs tracking-wide uppercase">
                                <th className="px-4 py-3 font-semibold">
                                    Member
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Paket
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Status
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Berakhir
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Pesan
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Aksi
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {members.data.map((m) => (
                                <tr
                                    key={m.id}
                                    className="border-t border-sidebar-border"
                                >
                                    <td className="px-4 py-3">
                                        <p className="font-medium">{m.name}</p>
                                        <p className="text-xs text-muted-foreground">
                                            {m.email}
                                        </p>
                                    </td>
                                    <td className="px-4 py-3 capitalize">
                                        {m.plan_label}
                                    </td>
                                    <td className="px-4 py-3">
                                        <Badge
                                            variant={
                                                m.status.key === 'kedaluwarsa'
                                                    ? 'destructive'
                                                    : 'outline'
                                            }
                                        >
                                            {m.status.label}
                                        </Badge>
                                    </td>
                                    <td className="px-4 py-3 text-xs text-muted-foreground">
                                        {m.expires_at
                                            ? formatRelative(m.expires_at)
                                            : 'Tanpa batas'}
                                    </td>
                                    <td className="px-4 py-3 text-right tabular-nums">
                                        {m.messages_total}
                                    </td>
                                    <td className="px-4 py-3 text-right">
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() =>
                                                extend(m.id, m.name)
                                            }
                                        >
                                            <CalendarPlus className="size-3.5 text-orange-600" />
                                            +14 hari
                                        </Button>
                                    </td>
                                </tr>
                            ))}
                            {members.data.length === 0 && (
                                <tr>
                                    <td
                                        colSpan={6}
                                        className="px-4 py-6 text-center text-muted-foreground"
                                    >
                                        Tidak ada member yang cocok.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        Halaman {members.current_page} dari{' '}
                        {members.last_page} · Total {members.total}
                    </span>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={members.current_page <= 1}
                            onClick={() => search(members.current_page - 1)}
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={
                                members.current_page >= members.last_page
                            }
                            onClick={() => search(members.current_page + 1)}
                        >
                            Berikutnya
                        </Button>
                    </div>
                </div>
            </div>
        </>
    );
}

SprMembers.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
        {
            title: 'Member',
            href: '/spr/members',
        },
    ],
};

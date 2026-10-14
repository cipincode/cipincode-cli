import { Head, Link } from '@inertiajs/react';
import {
    Bot,
    KeyRound,
    MessagesSquare,
    TerminalSquare,
    Users,
} from 'lucide-react';
import StatCard from '@/components/api-keys/stat-card';
import { formatRelative } from '@/components/api-keys/format';
import { Badge } from '@/components/ui/badge';

type Props = {
    stats: {
        users: number;
        users_week: number;
        messages: number;
        sessions: number;
        agents: number;
        keys: number;
    };
    plans: Record<string, number>;
    activity: { label: string; date: string; count: number }[];
    recent_users: {
        id: number;
        name: string;
        email: string;
        plan: string;
        is_admin: boolean;
        created_at: string;
    }[];
};

function Activity({ days }: { days: Props['activity'] }) {
    const max = Math.max(1, ...days.map((d) => d.count));

    return (
        <div className="flex h-32 items-end gap-2">
            {days.map((d) => (
                <div
                    key={d.date}
                    className="flex h-full flex-1 flex-col items-center justify-end gap-1"
                    title={`${d.date}: ${d.count} pesan`}
                >
                    <span className="text-xs text-muted-foreground tabular-nums">
                        {d.count}
                    </span>
                    <div
                        className={`w-full ${d.count ? 'bg-orange-600' : 'bg-muted'}`}
                        style={{
                            height: `${Math.max(d.count ? 6 : 2, (d.count / max) * 70)}%`,
                        }}
                    />
                    <span className="text-xs text-muted-foreground">
                        {d.label}
                    </span>
                </div>
            ))}
        </div>
    );
}

export default function SprDashboard({
    stats,
    plans,
    activity,
    recent_users,
}: Props) {
    return (
        <>
            <Head title="Superadmin" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div>
                    <h1 className="text-xl font-semibold tracking-tight">
                        Superadmin
                    </h1>
                    <p className="text-sm text-muted-foreground">
                        Ringkasan seluruh platform CipinCode.
                    </p>
                </div>

                <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <StatCard
                        label="Total pengguna"
                        value={stats.users}
                        caption={`+${stats.users_week} minggu ini`}
                        icon={Users}
                        tone="brand"
                    />
                    <StatCard
                        label="Total pesan"
                        value={stats.messages}
                        caption="Seluruh pemakaian"
                        icon={MessagesSquare}
                    />
                    <StatCard
                        label="Sesi"
                        value={stats.sessions}
                        caption="Percakapan tersimpan"
                        icon={TerminalSquare}
                    />
                    <StatCard
                        label="Agent"
                        value={stats.agents}
                        caption="Agent kustom"
                        icon={Bot}
                    />
                </div>

                <div className="grid gap-6 lg:grid-cols-2">
                    <section className="border border-sidebar-border bg-card">
                        <div className="border-b border-sidebar-border p-5">
                            <h2 className="font-medium">Pesan 7 hari</h2>
                            <p className="mt-0.5 text-sm text-muted-foreground">
                                Semua pengguna
                            </p>
                        </div>
                        <div className="p-5">
                            <Activity days={activity} />
                        </div>
                    </section>

                    <section className="border border-sidebar-border bg-card">
                        <div className="border-b border-sidebar-border p-5">
                            <h2 className="font-medium">Sebaran paket</h2>
                            <p className="mt-0.5 text-sm text-muted-foreground">
                                Jumlah pengguna per paket
                            </p>
                        </div>
                        <ul className="divide-y divide-sidebar-border">
                            {Object.entries(plans).map(([plan, total]) => (
                                <li
                                    key={plan}
                                    className="flex items-center justify-between px-5 py-3 text-sm"
                                >
                                    <span className="font-medium capitalize">
                                        {plan}
                                    </span>
                                    <span className="text-muted-foreground tabular-nums">
                                        {total} pengguna
                                    </span>
                                </li>
                            ))}
                            {Object.keys(plans).length === 0 && (
                                <li className="px-5 py-3 text-sm text-muted-foreground">
                                    Belum ada pengguna.
                                </li>
                            )}
                        </ul>
                    </section>
                </div>

                <section className="border border-sidebar-border bg-card">
                    <div className="flex items-center justify-between border-b border-sidebar-border p-5">
                        <h2 className="font-medium">Pengguna terbaru</h2>
                        <Link
                            href="/spr/users"
                            className="text-sm text-orange-600 hover:underline"
                        >
                            Semua pengguna
                        </Link>
                    </div>
                    <ul className="divide-y divide-sidebar-border">
                        {recent_users.map((u) => (
                            <li
                                key={u.id}
                                className="flex flex-wrap items-center gap-3 px-5 py-3"
                            >
                                <span className="min-w-0 flex-1">
                                    <span className="block truncate text-sm font-medium">
                                        {u.name}
                                    </span>
                                    <span className="block truncate text-xs text-muted-foreground">
                                        {u.email}
                                    </span>
                                </span>
                                {u.is_admin && (
                                    <Badge variant="outline">Admin</Badge>
                                )}
                                <Badge>{u.plan}</Badge>
                                <span className="text-xs text-muted-foreground">
                                    {formatRelative(u.created_at)}
                                </span>
                            </li>
                        ))}
                    </ul>
                </section>

                <div className="flex flex-wrap gap-3 text-sm">
                    <span className="inline-flex items-center border border-sidebar-border px-4 py-2 text-muted-foreground tabular-nums">
                        {stats.keys} API key aktif
                    </span>
                </div>
            </div>
        </>
    );
}

SprDashboard.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
    ],
};

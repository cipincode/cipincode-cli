import { Head, Link, usePage } from '@inertiajs/react';
import { useEffect, useState } from 'react';
import {
    ArrowRight,
    Bot,
    BookOpen,
    KeyRound,
    MessagesSquare,
    TerminalSquare,
} from 'lucide-react';
import { formatRelative } from '@/components/api-keys/format';
import StatCard from '@/components/api-keys/stat-card';
import { ModeBadge } from '@/components/mode-badge';
import {
    default as OnboardingTour,
    isOnboardingDone,
} from '@/components/onboarding/onboarding-tour';
import PlanBanner from '@/components/plans/plan-banner';
import type { CurrentPlan } from '@/components/plans/types';
import { formatReset } from '@/components/usage/format';
import SegmentBar from '@/components/usage/segment-bar';
import {
    agents,
    apiKeys,
    dashboard,
    sessions as sessionsRoute,
    setup,
    usage as usageRoute,
} from '@/routes';
import { show } from '@/routes/sessions';

type Meter = {
    used: number;
    limit: number;
    percent: number;
    exhausted: boolean;
    resets_at: string | null;
};

type Props = {
    usage: { session: Meter; weekly: Meter | null; plan: CurrentPlan };
    stats: { sessions: number; agents: number; keys: number; messages: number };
    activity: { label: string; date: string; count: number }[];
    recent: {
        uuid: string;
        title: string;
        mode: string;
        messages_count: number;
        last_activity_at: string;
    }[];
};

function MeterRow({ title, meter }: { title: string; meter: Meter }) {
    return (
        <div className="space-y-2">
            <div className="flex items-baseline justify-between text-sm">
                <span className="font-medium">{title}</span>
                <span
                    className={`tabular-nums ${meter.exhausted ? 'font-medium text-red-600' : 'text-muted-foreground'}`}
                >
                    {meter.used} / {meter.limit} pesan · {meter.percent}%
                </span>
            </div>
            <SegmentBar used={meter.used} limit={meter.limit} label={title} />
            <p className="text-xs text-muted-foreground">
                {meter.resets_at
                    ? `Reset ${formatRelative(meter.resets_at)} (${formatReset(meter.resets_at)})`
                    : 'Belum dimulai, jatah mulai dihitung saat pesan pertama.'}
            </p>
        </div>
    );
}

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

export default function Dashboard({ usage, stats, activity, recent }: Props) {
    const { auth } = usePage().props;
    const [tourOpen, setTourOpen] = useState(false);

    useEffect(() => {
        if (!isOnboardingDone(auth.user.id)) {
            const t = setTimeout(() => setTourOpen(true), 600);
            return () => clearTimeout(t);
        }
    }, [auth.user.id]);

    return (
        <>
            <Head title="Dashboard" />
            <OnboardingTour
                open={tourOpen}
                onOpenChange={setTourOpen}
                userId={auth.user.id}
            />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                        <h1 className="text-xl font-semibold tracking-tight">
                            Dashboard
                        </h1>
                        <p className="text-sm text-muted-foreground">
                            Ringkasan pemakaian CipinCode Anda.
                        </p>
                    </div>
                    <button
                        type="button"
                        onClick={() => setTourOpen(true)}
                        className="inline-flex h-9 items-center gap-1.5 border border-sidebar-border px-3 text-sm font-medium transition-colors hover:bg-sidebar-accent"
                    >
                        <BookOpen className="size-4 text-orange-600" />
                        Lihat Panduan
                    </button>
                </div>

                <PlanBanner plan={usage.plan} />

                <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <StatCard
                        label="Total pesan"
                        value={stats.messages}
                        caption="Seluruh pemakaian"
                        icon={MessagesSquare}
                        tone="brand"
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
                    <StatCard
                        label="API key aktif"
                        value={stats.keys}
                        caption="Dapat dipakai CLI"
                        icon={KeyRound}
                        tone={stats.keys ? 'success' : 'muted'}
                    />
                </div>

                <div className="grid gap-6 lg:grid-cols-2">
                    <section className="border border-sidebar-border bg-card">
                        <div className="flex items-center justify-between border-b border-sidebar-border p-5">
                            <h2 className="font-medium">Jatah pemakaian</h2>
                            <Link
                                href={usageRoute()}
                                className="inline-flex items-center gap-1 text-sm text-orange-600 hover:underline"
                            >
                                Detail <ArrowRight className="size-3.5" />
                            </Link>
                        </div>
                        <div className="space-y-6 p-5">
                            <MeterRow
                                title="Sesi 5 jam"
                                meter={usage.session}
                            />
                            {usage.weekly && (
                                <MeterRow
                                    title="Mingguan"
                                    meter={usage.weekly}
                                />
                            )}
                        </div>
                    </section>

                    <section className="border border-sidebar-border bg-card">
                        <div className="border-b border-sidebar-border p-5">
                            <h2 className="font-medium">Aktivitas 7 hari</h2>
                            <p className="mt-0.5 text-sm text-muted-foreground">
                                Pesan per hari
                            </p>
                        </div>
                        <div className="p-5">
                            <Activity days={activity} />
                        </div>
                    </section>
                </div>

                <section className="border border-sidebar-border bg-card">
                    <div className="flex items-center justify-between border-b border-sidebar-border p-5">
                        <h2 className="font-medium">Sesi terbaru</h2>
                        <Link
                            href={sessionsRoute()}
                            className="inline-flex items-center gap-1 text-sm text-orange-600 hover:underline"
                        >
                            Semua sesi <ArrowRight className="size-3.5" />
                        </Link>
                    </div>
                    {recent.length === 0 ? (
                        <div className="flex flex-wrap items-center justify-between gap-3 p-5 text-sm text-muted-foreground">
                            <span>
                                Belum ada percakapan. Pasang CLI dan mulai
                                bertanya.
                            </span>
                            <Link
                                href={setup()}
                                className="inline-flex h-9 items-center bg-orange-600 px-4 font-medium text-white hover:bg-orange-700"
                            >
                                Setup Installation
                            </Link>
                        </div>
                    ) : (
                        <ul className="divide-y divide-sidebar-border">
                            {recent.map((s) => (
                                <li key={s.uuid}>
                                    <Link
                                        href={show(s.uuid)}
                                        className="flex flex-wrap items-center gap-3 px-5 py-3 hover:bg-sidebar-accent"
                                    >
                                        <span className="min-w-0 flex-1 truncate text-sm font-medium">
                                            {s.title}
                                        </span>
                                        <ModeBadge mode={s.mode} />
                                        <span className="text-xs text-muted-foreground tabular-nums">
                                            {s.messages_count} pesan ·{' '}
                                            {formatRelative(s.last_activity_at)}
                                        </span>
                                    </Link>
                                </li>
                            ))}
                        </ul>
                    )}
                </section>

                <div className="flex flex-wrap gap-3 text-sm">
                    <Link
                        href={agents()}
                        className="border border-sidebar-border px-4 py-2 hover:bg-sidebar-accent"
                    >
                        Kelola agent
                    </Link>
                    <Link
                        href={apiKeys()}
                        className="border border-sidebar-border px-4 py-2 hover:bg-sidebar-accent"
                    >
                        Kelola API key
                    </Link>
                </div>
            </div>
        </>
    );
}

Dashboard.layout = {
    breadcrumbs: [
        {
            title: 'Dashboard',
            href: dashboard(),
        },
    ],
};

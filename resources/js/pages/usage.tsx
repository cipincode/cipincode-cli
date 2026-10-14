import { Head, Link, router } from '@inertiajs/react';
import { useEffect } from 'react';
import { formatRelative } from '@/components/api-keys/format';
import PlanBanner from '@/components/plans/plan-banner';
import PlanCard from '@/components/plans/plan-card';
import type { CurrentPlan, PlanInfo } from '@/components/plans/types';
import SegmentBar from '@/components/usage/segment-bar';
import { formatReset } from '@/components/usage/format';
import { billing, usage as usageRoute } from '@/routes';

type Meter = {
    used: number;
    limit: number;
    percent: number;
    exhausted: boolean;
    resets_at: string | null;
};

type Props = {
    usage: { session: Meter; weekly: Meter | null; plan: CurrentPlan };
    sessionHours: number;
    plans: PlanInfo[];
};

function Panel({
    title,
    description,
    meter,
    idle,
}: {
    title: string;
    description: string;
    meter: Meter;
    idle: string;
}) {
    return (
        <section className="border border-sidebar-border bg-card">
            <div className="flex items-start justify-between gap-4 border-b border-sidebar-border p-5">
                <div>
                    <h2 className="font-medium">{title}</h2>
                    <p className="mt-0.5 text-sm text-muted-foreground">
                        {description}
                    </p>
                </div>
                <p
                    className={`text-3xl font-semibold tracking-tight tabular-nums ${
                        meter.exhausted ? 'text-red-600' : 'text-foreground'
                    }`}
                >
                    {meter.percent}%
                </p>
            </div>

            <div className="space-y-3 p-5">
                <SegmentBar
                    used={meter.used}
                    limit={meter.limit}
                    label={title}
                />
                <div className="flex flex-wrap items-center justify-between gap-2 text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        {meter.used} / {meter.limit} pesan
                    </span>
                    <span
                        className={
                            meter.exhausted
                                ? 'font-medium text-red-600'
                                : 'text-muted-foreground'
                        }
                    >
                        {meter.resets_at
                            ? `${meter.exhausted ? 'Habis · reset' : 'Reset'} ${formatRelative(meter.resets_at)} (${formatReset(meter.resets_at)})`
                            : idle}
                    </span>
                </div>
            </div>
        </section>
    );
}

function UpgradeSection({
    plans,
    current,
}: {
    plans: PlanInfo[];
    current: CurrentPlan;
}) {
    const rank = (key: string) => plans.findIndex((p) => p.key === key);
    const currentRank = rank(current.name);
    const options = plans.filter((p) => p.paid && rank(p.key) > currentRank);

    if (options.length === 0) {
        return (
            <section className="border border-sidebar-border bg-card p-5 text-sm text-muted-foreground">
                Anda berada di paket tertinggi, <strong>{current.label}</strong>
                . Terima kasih sudah memakai CipinCode.
            </section>
        );
    }

    const free = current.name === 'free';

    return (
        <section className="border border-sidebar-border bg-card">
            <div className="flex flex-wrap items-center justify-between gap-3 border-b border-sidebar-border p-5">
                <div>
                    <h2 className="font-medium">
                        {current.expired
                            ? 'Lanjutkan dengan paket berbayar'
                            : 'Butuh jatah lebih besar?'}
                    </h2>
                    <p className="mt-0.5 text-sm text-muted-foreground">
                        {free
                            ? 'Paket gratis dibatasi sesi dan masa percobaan. Paket berbayar menambah batas mingguan yang lebih longgar.'
                            : 'Naikkan paket untuk batas sesi dan mingguan yang lebih besar.'}
                    </p>
                </div>
                <Link
                    href={billing()}
                    className="inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent"
                >
                    Bandingkan semua paket
                </Link>
            </div>

            <div className="grid gap-4 p-5 md:grid-cols-3">
                {options.map((plan) => (
                    <PlanCard
                        key={plan.key}
                        plan={plan}
                        isCurrent={false}
                        highlight={plan.key === 'pro'}
                        cta={{
                            label: `Upgrade ke ${plan.label}`,
                            href: billing({ query: { plan: plan.key } }),
                        }}
                    />
                ))}
            </div>
        </section>
    );
}

export default function Usage({ usage, sessionHours, plans }: Props) {
    useEffect(() => {
        const timer = setInterval(
            () => router.reload({ only: ['usage'] }),
            30_000,
        );

        return () => clearInterval(timer);
    }, []);

    return (
        <>
            <Head title="Usage" />

            <div className="flex flex-1 flex-col">
                <div className="border-b border-sidebar-border px-6 py-5">
                    <h1 className="text-xl font-semibold tracking-tight">
                        Usage
                    </h1>
                    <p className="mt-1 text-sm text-muted-foreground">
                        Pantau jatah pemakaian CLI Anda. Data diperbarui
                        otomatis.
                    </p>
                </div>

                <div className="mx-auto flex w-full max-w-5xl flex-col gap-6 p-6">
                    <PlanBanner plan={usage.plan} />
                    <Panel
                        title={`Sesi ${sessionHours} jam`}
                        description={`Mulai saat pesan pertama, reset otomatis ${sessionHours} jam kemudian.`}
                        meter={usage.session}
                        idle="Belum dimulai"
                    />
                    {usage.weekly ? (
                        <Panel
                            title="Mingguan"
                            description="Reset setiap Senin pukul 00.00."
                            meter={usage.weekly}
                            idle="Belum dimulai"
                        />
                    ) : (
                        <p className="border border-dashed border-sidebar-border p-4 text-sm text-muted-foreground">
                            Paket {usage.plan.label} hanya memakai batas sesi{' '}
                            {sessionHours} jam, tanpa batas mingguan. Paket
                            berbayar memiliki batas sesi dan mingguan.
                        </p>
                    )}

                    <UpgradeSection plans={plans} current={usage.plan} />

                    <div className="flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-muted-foreground">
                        <span className="flex items-center gap-2">
                            <span className="size-3 border border-sidebar-border bg-muted/40" />
                            Belum terpakai
                        </span>
                        <span className="flex items-center gap-2">
                            <span className="size-3 border border-orange-600 bg-orange-600" />
                            Terpakai
                        </span>
                        <span className="flex items-center gap-2">
                            <span className="size-3 border border-red-600 bg-red-600" />
                            Habis
                        </span>
                    </div>
                </div>
            </div>
        </>
    );
}

Usage.layout = {
    breadcrumbs: [
        {
            title: 'Usage',
            href: usageRoute(),
        },
    ],
};

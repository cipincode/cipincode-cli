import { Head } from '@inertiajs/react';
import { useState } from 'react';
import PlanBanner from '@/components/plans/plan-banner';
import PlanCard from '@/components/plans/plan-card';
import type { CurrentPlan, PlanInfo } from '@/components/plans/types';
import { billing } from '@/routes';

type Props = {
    plans: PlanInfo[];
    current: CurrentPlan;
    selected: string | null;
};

export default function Billing({ plans, current, selected: initial }: Props) {
    const [selected, setSelected] = useState<string | null>(initial);
    const rank = (key: string) => plans.findIndex((p) => p.key === key);
    const currentRank = rank(current.name);
    const choice = plans.find((p) => p.key === selected && p.paid) ?? null;

    return (
        <>
            <Head title="Billing" />

            <div className="flex flex-1 flex-col">
                <div className="border-b border-sidebar-border px-6 py-5">
                    <h1 className="text-xl font-semibold tracking-tight">
                        Billing
                    </h1>
                    <p className="mt-1 text-sm text-muted-foreground">
                        Pilih paket yang sesuai. Semua harga dalam rupiah per
                        bulan.
                    </p>
                </div>

                <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 p-6">
                    <PlanBanner plan={current} />

                    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                        {plans.map((plan) => {
                            const isCurrent = plan.key === current.name;
                            const upgrade =
                                plan.paid && rank(plan.key) > currentRank;

                            return (
                                <PlanCard
                                    key={plan.key}
                                    plan={plan}
                                    isCurrent={isCurrent}
                                    highlight={plan.key === 'pro'}
                                    selected={selected === plan.key}
                                    cta={
                                        isCurrent
                                            ? { label: 'Paket Anda' }
                                            : upgrade
                                              ? {
                                                    label: `Pilih ${plan.label}`,
                                                    onClick: () =>
                                                        setSelected(plan.key),
                                                }
                                              : undefined
                                    }
                                />
                            );
                        })}
                    </div>

                    {choice && (
                        <section className="border border-orange-600 bg-card">
                            <div className="border-b border-orange-600 bg-orange-600/10 px-5 py-3">
                                <h2 className="font-medium">
                                    Ringkasan: Paket {choice.label}
                                </h2>
                            </div>
                            <div className="flex flex-wrap items-center justify-between gap-4 p-5">
                                <div className="text-sm">
                                    <p className="text-2xl font-semibold tabular-nums">
                                        {choice.price_label}{' '}
                                        <span className="text-sm font-normal text-muted-foreground">
                                            {choice.period}
                                        </span>
                                    </p>
                                    <p className="mt-1 text-muted-foreground">
                                        Pembayaran online belum aktif. Paket
                                        Anda tetap{' '}
                                        <strong>{current.label}</strong> sampai
                                        pembayaran tersedia.
                                    </p>
                                </div>
                                <button
                                    type="button"
                                    disabled
                                    className="inline-flex h-9 items-center bg-orange-600 px-5 text-sm font-medium text-white opacity-60"
                                >
                                    Bayar segera hadir
                                </button>
                            </div>
                        </section>
                    )}
                </div>
            </div>
        </>
    );
}

Billing.layout = {
    breadcrumbs: [
        {
            title: 'Billing',
            href: billing(),
        },
    ],
};

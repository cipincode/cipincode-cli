import { Link } from '@inertiajs/react';
import { Check } from 'lucide-react';
import type { PlanInfo } from './types';

type Cta = {
    label: string;
    href?: Parameters<typeof Link>[0]['href'];
    onClick?: () => void;
};

type Props = {
    plan: PlanInfo;
    /** Paket yang sedang dipakai pengguna. */
    isCurrent: boolean;
    /** Ditandai "Populer". */
    highlight?: boolean;
    /** Sedang dipilih (halaman Billing). */
    selected?: boolean;
    cta?: Cta;
};

const base =
    'inline-flex h-9 w-full items-center justify-center text-sm font-medium transition-colors';

export default function PlanCard({
    plan,
    isCurrent,
    highlight = false,
    selected = false,
    cta,
}: Props) {
    const accent = selected || highlight;

    return (
        <article
            className={`relative flex flex-col border bg-card ${
                accent ? 'border-orange-600' : 'border-sidebar-border'
            } ${selected ? 'ring-1 ring-orange-600' : ''}`}
        >
            {highlight && (
                <span className="absolute -top-px right-4 bg-orange-600 px-2 py-0.5 text-[10px] font-semibold tracking-wider text-white uppercase">
                    Populer
                </span>
            )}

            <div className="border-b border-sidebar-border p-5">
                <div className="flex items-center gap-2">
                    <h3 className="text-lg font-semibold">{plan.label}</h3>
                    {isCurrent && (
                        <span className="border border-orange-600 px-1.5 py-0.5 text-[10px] font-medium tracking-wide text-orange-600 uppercase">
                            Paket Anda
                        </span>
                    )}
                </div>
                <p className="mt-3 flex items-baseline gap-1.5">
                    <span className="text-3xl font-semibold tracking-tight tabular-nums">
                        {plan.price_label}
                    </span>
                    {plan.period && (
                        <span className="text-sm text-muted-foreground">
                            {plan.period}
                        </span>
                    )}
                </p>
            </div>

            <ul className="flex-1 space-y-2.5 p-5 text-sm">
                {plan.features.map((feature) => (
                    <li key={feature} className="flex items-start gap-2.5">
                        <Check className="mt-0.5 size-4 shrink-0 text-orange-600" />
                        <span>{feature}</span>
                    </li>
                ))}
            </ul>

            {cta && (
                <div className="border-t border-sidebar-border p-5">
                    {cta.href ? (
                        <Link
                            href={cta.href}
                            className={`${base} ${
                                accent
                                    ? 'bg-orange-600 text-white hover:bg-orange-700'
                                    : 'border border-sidebar-border hover:bg-sidebar-accent'
                            }`}
                        >
                            {cta.label}
                        </Link>
                    ) : (
                        <button
                            type="button"
                            onClick={cta.onClick}
                            disabled={!cta.onClick}
                            className={`${base} disabled:cursor-default disabled:opacity-60 ${
                                accent
                                    ? 'bg-orange-600 text-white hover:bg-orange-700'
                                    : 'border border-sidebar-border hover:bg-sidebar-accent'
                            }`}
                        >
                            {cta.label}
                        </button>
                    )}
                </div>
            )}
        </article>
    );
}

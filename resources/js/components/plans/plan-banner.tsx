import { formatRelative } from '@/components/api-keys/format';
import { formatReset } from '@/components/usage/format';
import type { CurrentPlan } from './types';

/** Satu baris ringkas tentang paket aktif dan masa berlakunya. */
export default function PlanBanner({ plan }: { plan: CurrentPlan }) {
    const free = plan.name === 'free';
    const tone = plan.expired
        ? 'border-red-600 bg-red-600/10'
        : 'border-sidebar-border bg-card';

    let detail: string;
    if (plan.expired) {
        detail = free
            ? 'Masa percobaan gratis sudah berakhir. Pilih paket berbayar untuk melanjutkan.'
            : 'Masa aktif paket sudah berakhir. Perpanjang untuk melanjutkan.';
    } else if (plan.expires_at) {
        detail = `${free ? 'Masa percobaan berakhir' : 'Aktif sampai'} ${formatRelative(plan.expires_at)} (${formatReset(plan.expires_at)}).`;
    } else {
        detail = 'Aktif tanpa batas waktu.';
    }

    return (
        <section
            className={`flex flex-wrap items-center gap-x-4 gap-y-1 border p-5 ${tone}`}
        >
            <span
                className={`border px-2 py-0.5 text-xs font-medium tracking-wide uppercase ${
                    plan.expired
                        ? 'border-red-600 text-red-600'
                        : 'border-orange-600 text-orange-600'
                }`}
            >
                Paket {plan.label}
            </span>
            <span
                className={`text-sm ${plan.expired ? 'font-medium text-red-600' : 'text-muted-foreground'}`}
            >
                {detail}
            </span>
        </section>
    );
}

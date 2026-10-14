import type { LucideIcon } from 'lucide-react';

type Props = {
    label: string;
    value: number | string;
    caption: string;
    icon: LucideIcon;
    tone?: 'default' | 'success' | 'muted' | 'brand';
};

const tones = {
    default: 'text-foreground',
    success: 'text-emerald-600 dark:text-emerald-500',
    muted: 'text-muted-foreground',
    brand: 'text-orange-600',
};

export default function StatCard({
    label,
    value,
    caption,
    icon: Icon,
    tone = 'default',
}: Props) {
    return (
        <div className="flex flex-col gap-3 border border-sidebar-border bg-card p-5">
            <div className="flex items-center justify-between">
                <span className="text-xs font-medium tracking-wide text-muted-foreground uppercase">
                    {label}
                </span>
                <Icon className={`size-4 ${tones[tone]}`} />
            </div>
            <div>
                <p
                    className={`text-3xl font-semibold tracking-tight tabular-nums ${tones[tone]}`}
                >
                    {value}
                </p>
                <p className="mt-1 text-xs text-muted-foreground">{caption}</p>
            </div>
        </div>
    );
}

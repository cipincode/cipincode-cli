const MODE: Record<string, { label: string; className: string }> = {
    chat: {
        label: 'Chat',
        className:
            'border-sky-500/40 bg-sky-500/10 text-sky-600 dark:text-sky-400',
    },
    coding: {
        label: 'Coding',
        className: 'border-orange-600/40 bg-orange-600/10 text-orange-600',
    },
    debugging: {
        label: 'Debugging',
        className:
            'border-rose-500/40 bg-rose-500/10 text-rose-600 dark:text-rose-400',
    },
};

const EFFORT: Record<string, string> = {
    low: 'border-emerald-500/40 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400',
    medium: 'border-sky-500/40 bg-sky-500/10 text-sky-600 dark:text-sky-400',
    high: 'border-amber-500/40 bg-amber-500/10 text-amber-600 dark:text-amber-400',
    xhigh: 'border-orange-500/40 bg-orange-500/10 text-orange-600 dark:text-orange-400',
    ultra: 'border-fuchsia-500/40 bg-fuchsia-500/10 text-fuchsia-600 dark:text-fuchsia-400',
};

const base =
    'inline-flex items-center border px-2 py-0.5 text-xs font-medium whitespace-nowrap';

export function ModeBadge({ mode }: { mode: string }) {
    const m = MODE[mode] ?? MODE.chat;

    return <span className={`${base} ${m.className}`}>{m.label}</span>;
}

export function EffortBadge({ effort }: { effort: string }) {
    return (
        <span className={`${base} ${EFFORT[effort] ?? EFFORT.medium}`}>
            {effort}
        </span>
    );
}

import { Head, router } from '@inertiajs/react';
import { Check, Plus, Search, Sparkles } from 'lucide-react';
import { useMemo, useState } from 'react';
import { EffortBadge, ModeBadge } from '@/components/mode-badge';
import { agents, skills as skillsRoute } from '@/routes';
import { install } from '@/routes/skills';
import { Link } from '@inertiajs/react';

type Skill = {
    key: string;
    category: string;
    name: string;
    description: string;
    mode: string;
    effort: string;
    instructions: string;
    installed: boolean;
};

type Props = {
    skills: Skill[];
    categories: Record<string, string>;
    agentCount: number;
    agentLimit: number;
};

export default function Skills({
    skills,
    categories,
    agentCount,
    agentLimit,
}: Props) {
    const [q, setQ] = useState('');
    const [category, setCategory] = useState<string | null>(null);
    const [busy, setBusy] = useState<string | null>(null);

    const filtered = useMemo(() => {
        const needle = q.trim().toLowerCase();

        return skills.filter(
            (s) =>
                (!category || s.category === category) &&
                (!needle ||
                    `${s.name} ${s.description}`
                        .toLowerCase()
                        .includes(needle)),
        );
    }, [skills, q, category]);

    const counts = useMemo(() => {
        const c: Record<string, number> = {};
        skills.forEach((s) => (c[s.category] = (c[s.category] ?? 0) + 1));

        return c;
    }, [skills]);

    const add = (s: Skill) =>
        router.post(
            install.url(s.key),
            {},
            {
                preserveScroll: true,
                onStart: () => setBusy(s.key),
                onFinish: () => setBusy(null),
            },
        );

    const chip = (active: boolean) =>
        `h-8 border px-3 text-sm whitespace-nowrap transition-colors ${
            active
                ? 'border-orange-600 bg-orange-600 text-white'
                : 'border-sidebar-border hover:bg-sidebar-accent'
        }`;

    return (
        <>
            <Head title="Skills" />

            <div className="flex flex-1 flex-col">
                <div className="flex flex-wrap items-center justify-between gap-4 border-b border-sidebar-border px-6 py-5">
                    <div>
                        <h1 className="text-xl font-semibold tracking-tight">
                            Skills
                        </h1>
                        <p className="mt-1 text-sm text-muted-foreground">
                            Koleksi skill siap pakai untuk AI. Tambahkan ke{' '}
                            <Link
                                href={agents()}
                                className="text-orange-600 hover:underline"
                            >
                                Agents
                            </Link>
                            , lalu pilih di CLI dengan{' '}
                            <code className="bg-muted px-1.5 py-0.5 font-mono text-xs">
                                /agent
                            </code>
                            .
                        </p>
                    </div>
                    <span className="text-sm text-muted-foreground tabular-nums">
                        {agentCount} / {agentLimit} agent terpakai
                    </span>
                </div>

                <div className="space-y-3 border-b border-sidebar-border px-6 py-4">
                    <div className="relative max-w-md">
                        <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                        <input
                            value={q}
                            onChange={(e) => setQ(e.target.value)}
                            placeholder="Cari skill…"
                            className="h-10 w-full border border-sidebar-border bg-background pr-3 pl-10 text-sm outline-none focus:border-orange-600 focus:ring-1 focus:ring-orange-600/40"
                        />
                    </div>
                    <div className="flex gap-2 overflow-x-auto pb-1">
                        <button
                            type="button"
                            className={chip(category === null)}
                            onClick={() => setCategory(null)}
                        >
                            Semua ({skills.length})
                        </button>
                        {Object.entries(categories).map(([key, label]) => (
                            <button
                                key={key}
                                type="button"
                                className={chip(category === key)}
                                onClick={() => setCategory(key)}
                            >
                                {label} ({counts[key] ?? 0})
                            </button>
                        ))}
                    </div>
                </div>

                <div className="flex-1 p-6">
                    {filtered.length === 0 ? (
                        <div className="border border-dashed border-sidebar-border p-10 text-center text-sm text-muted-foreground">
                            Tidak ada skill yang cocok.{' '}
                            <button
                                type="button"
                                className="text-orange-600 hover:underline"
                                onClick={() => {
                                    setQ('');
                                    setCategory(null);
                                }}
                            >
                                Reset filter
                            </button>
                        </div>
                    ) : (
                        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                            {filtered.map((s) => (
                                <article
                                    key={s.key}
                                    className="flex flex-col border border-sidebar-border bg-card"
                                >
                                    <div className="flex-1 space-y-3 p-5">
                                        <div className="flex items-start justify-between gap-3">
                                            <div className="flex items-center gap-2">
                                                <Sparkles className="size-4 shrink-0 text-orange-600" />
                                                <h2 className="font-medium">
                                                    {s.name}
                                                </h2>
                                            </div>
                                            <span className="text-xs whitespace-nowrap text-muted-foreground">
                                                {categories[s.category]}
                                            </span>
                                        </div>
                                        <p className="text-sm text-muted-foreground">
                                            {s.description}
                                        </p>
                                        <div className="flex flex-wrap gap-2">
                                            <ModeBadge mode={s.mode} />
                                            <EffortBadge effort={s.effort} />
                                        </div>
                                    </div>
                                    <div className="border-t border-sidebar-border px-5 py-3">
                                        {s.installed ? (
                                            <span className="inline-flex h-8 items-center gap-1.5 text-sm text-emerald-600">
                                                <Check className="size-4" />
                                                Sudah di Agents
                                            </span>
                                        ) : (
                                            <button
                                                type="button"
                                                disabled={busy === s.key}
                                                onClick={() => add(s)}
                                                className="inline-flex h-8 items-center gap-1.5 bg-orange-600 px-3 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:opacity-60"
                                            >
                                                <Plus className="size-4" />
                                                {busy === s.key
                                                    ? 'Menambahkan…'
                                                    : 'Tambahkan ke Agents'}
                                            </button>
                                        )}
                                    </div>
                                </article>
                            ))}
                        </div>
                    )}
                </div>
            </div>
        </>
    );
}

Skills.layout = {
    breadcrumbs: [{ title: 'Skills', href: skillsRoute() }],
};

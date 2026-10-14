import { Head, router } from '@inertiajs/react';
import {
    Activity,
    CircleCheck,
    CircleOff,
    KeyRound,
    Plus,
    Search,
} from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import ApiKeyController from '@/actions/App/Http/Controllers/ApiKeyController';
import {
    CreateKeyDialog,
    DeleteKeyDialog,
    EditKeyDialog,
    RevealKeyDialog,
} from '@/components/api-keys/key-dialogs';
import KeysTable from '@/components/api-keys/keys-table';
import StatCard from '@/components/api-keys/stat-card';
import type { ApiKey, ApiKeyStats } from '@/components/api-keys/types';
import { apiKeys } from '@/routes';

type Filter = 'all' | 'active' | 'inactive';

const filters: { value: Filter; label: string }[] = [
    { value: 'all', label: 'Semua' },
    { value: 'active', label: 'Aktif' },
    { value: 'inactive', label: 'Nonaktif' },
];

export default function ApiKeys({
    keys,
    stats,
    limit,
    newKey = null,
}: {
    keys: ApiKey[];
    stats: ApiKeyStats;
    limit: number;
    newKey?: string | null;
}) {
    const [search, setSearch] = useState('');
    const [filter, setFilter] = useState<Filter>('all');
    const [createOpen, setCreateOpen] = useState(false);
    const [revealOpen, setRevealOpen] = useState(!!newKey);
    const [editing, setEditing] = useState<ApiKey | null>(null);
    const [deleting, setDeleting] = useState<ApiKey | null>(null);

    useEffect(() => {
        if (newKey) setRevealOpen(true);
    }, [newKey]);

    const visible = useMemo(() => {
        const term = search.trim().toLowerCase();

        return keys.filter((key) => {
            if (filter === 'active' && !key.is_active) return false;
            if (filter === 'inactive' && key.is_active) return false;

            return term === '' || key.name.toLowerCase().includes(term);
        });
    }, [keys, search, filter]);

    const toggle = (key: ApiKey) =>
        router.patch(
            ApiKeyController.update.url(key.id),
            { is_active: !key.is_active },
            { preserveScroll: true },
        );

    const atLimit = stats.total >= limit;

    return (
        <>
            <Head title="API Keys" />

            <div className="flex flex-1 flex-col">
                <div className="flex flex-wrap items-center justify-between gap-4 border-b border-sidebar-border px-6 py-5">
                    <div>
                        <h1 className="text-xl font-semibold tracking-tight">
                            API Keys
                        </h1>
                        <p className="mt-1 text-sm text-muted-foreground">
                            Kelola token untuk menghubungkan CLI ke akun Anda.
                        </p>
                    </div>
                    <button
                        type="button"
                        disabled={atLimit}
                        onClick={() => setCreateOpen(true)}
                        title={
                            atLimit
                                ? `Batas ${limit} key tercapai. Hapus key lama dulu.`
                                : undefined
                        }
                        className="inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                        <Plus className="size-4" />
                        Buat API key
                    </button>
                </div>

                <div className="flex flex-1 flex-col gap-6 p-6">
                    <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
                        <StatCard
                            label="Total key"
                            value={stats.total}
                            caption={`dari ${limit} slot tersedia`}
                            icon={KeyRound}
                        />
                        <StatCard
                            label="Aktif"
                            value={stats.active}
                            caption="dapat dipakai CLI"
                            icon={CircleCheck}
                            tone="success"
                        />
                        <StatCard
                            label="Nonaktif"
                            value={stats.inactive}
                            caption="ditolak oleh server"
                            icon={CircleOff}
                            tone="muted"
                        />
                        <StatCard
                            label="Dipakai 7 hari"
                            value={stats.used_7d}
                            caption="key dengan aktivitas"
                            icon={Activity}
                            tone="brand"
                        />
                    </div>

                    <section className="flex flex-1 flex-col border border-sidebar-border bg-card">
                        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-sidebar-border p-4">
                            <div className="relative w-full sm:max-w-xs">
                                <Search className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                                <input
                                    type="search"
                                    value={search}
                                    onChange={(e) => setSearch(e.target.value)}
                                    placeholder="Cari nama key…"
                                    aria-label="Cari nama key"
                                    className="h-9 w-full border border-input bg-transparent pr-3 pl-9 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
                                />
                            </div>

                            <div
                                role="group"
                                aria-label="Filter status"
                                className="inline-flex border border-sidebar-border"
                            >
                                {filters.map((f) => (
                                    <button
                                        key={f.value}
                                        type="button"
                                        aria-pressed={filter === f.value}
                                        onClick={() => setFilter(f.value)}
                                        className={`h-9 px-4 text-sm font-medium transition-colors not-first:border-l not-first:border-sidebar-border ${
                                            filter === f.value
                                                ? 'bg-orange-600 text-white'
                                                : 'text-muted-foreground hover:bg-sidebar-accent hover:text-foreground'
                                        }`}
                                    >
                                        {f.label}
                                    </button>
                                ))}
                            </div>
                        </div>

                        {keys.length === 0 ? (
                            <div className="flex flex-1 flex-col items-center justify-center gap-3 px-6 py-16 text-center">
                                <div className="flex size-12 items-center justify-center border border-sidebar-border text-orange-600">
                                    <KeyRound className="size-6" />
                                </div>
                                <h2 className="font-medium">
                                    Belum ada API key
                                </h2>
                                <p className="max-w-sm text-sm text-muted-foreground">
                                    Buat API key pertama Anda untuk masuk ke CLI
                                    dari terminal.
                                </p>
                                <button
                                    type="button"
                                    onClick={() => setCreateOpen(true)}
                                    className="mt-1 inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700"
                                >
                                    <Plus className="size-4" />
                                    Buat API key
                                </button>
                            </div>
                        ) : visible.length === 0 ? (
                            <div className="flex flex-1 flex-col items-center justify-center gap-2 px-6 py-16 text-center">
                                <p className="font-medium">
                                    Tidak ada key yang cocok
                                </p>
                                <p className="text-sm text-muted-foreground">
                                    Ubah kata kunci pencarian atau filter
                                    status.
                                </p>
                            </div>
                        ) : (
                            <KeysTable
                                keys={visible}
                                onEdit={setEditing}
                                onToggle={toggle}
                                onDelete={setDeleting}
                            />
                        )}

                        {keys.length > 0 && (
                            <div className="mt-auto border-t border-sidebar-border px-5 py-3 text-xs text-muted-foreground">
                                Menampilkan {visible.length} dari {keys.length}{' '}
                                key.
                            </div>
                        )}
                    </section>
                </div>
            </div>

            <CreateKeyDialog open={createOpen} onOpenChange={setCreateOpen} />
            <EditKeyDialog apiKey={editing} onClose={() => setEditing(null)} />
            <DeleteKeyDialog
                apiKey={deleting}
                onClose={() => setDeleting(null)}
            />
            <RevealKeyDialog
                newKey={newKey}
                open={revealOpen}
                onOpenChange={setRevealOpen}
            />
        </>
    );
}

ApiKeys.layout = {
    breadcrumbs: [
        {
            title: 'API Keys',
            href: apiKeys(),
        },
    ],
};

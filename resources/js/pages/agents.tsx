import { Head, Link, router } from '@inertiajs/react';
import { Bot, Pencil, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import AgentController from '@/actions/App/Http/Controllers/AgentController';
import { formatRelative } from '@/components/api-keys/format';
import AgentDialog from '@/components/agents/agent-dialog';
import type { Agent } from '@/components/agents/agent-dialog';
import { EffortBadge, ModeBadge } from '@/components/mode-badge';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { agents as agentsRoute, setup } from '@/routes';

type Props = {
    agents: Agent[];
    limit: number;
    modes: string[];
    efforts: string[];
};

export default function Agents({ agents, limit, modes, efforts }: Props) {
    const [dialog, setDialog] = useState<{
        open: boolean;
        agent: Agent | null;
    }>({ open: false, agent: null });
    const [deleting, setDeleting] = useState<Agent | null>(null);
    const atLimit = agents.length >= limit;

    return (
        <>
            <Head title="Agents" />

            <div className="flex flex-1 flex-col">
                <div className="flex flex-wrap items-center justify-between gap-4 border-b border-sidebar-border px-6 py-5">
                    <div>
                        <h1 className="text-xl font-semibold tracking-tight">
                            Agents
                        </h1>
                        <p className="mt-1 text-sm text-muted-foreground">
                            Buat agent dengan instruksi khusus, lalu pakai di
                            CLI lewat perintah{' '}
                            <code className="bg-muted px-1.5 py-0.5 font-mono text-xs">
                                /agent
                            </code>
                            .
                        </p>
                    </div>
                    <div className="flex items-center gap-3">
                        <span className="text-sm text-muted-foreground tabular-nums">
                            {agents.length} / {limit} agent
                        </span>
                        <button
                            type="button"
                            disabled={atLimit}
                            onClick={() =>
                                setDialog({ open: true, agent: null })
                            }
                            title={
                                atLimit
                                    ? `Batas ${limit} agent tercapai.`
                                    : undefined
                            }
                            className="inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:cursor-not-allowed disabled:opacity-50"
                        >
                            <Plus className="size-4" />
                            Buat agent
                        </button>
                    </div>
                </div>

                <div className="flex w-full flex-1 flex-col gap-6 p-6">
                    {agents.length === 0 ? (
                        <section className="flex flex-1 flex-col items-center justify-center gap-3 border border-dashed border-sidebar-border bg-card px-6 py-16 text-center">
                            <div className="flex size-12 items-center justify-center border border-sidebar-border text-orange-600">
                                <Bot className="size-6" />
                            </div>
                            <h2 className="font-medium">Belum ada agent</h2>
                            <p className="max-w-md text-sm text-muted-foreground">
                                Agent adalah peran khusus untuk AI, misalnya
                                peninjau kode atau penjelas sederhana.
                                Instruksinya otomatis ikut di setiap pesan saat
                                agent dipilih. Anda bisa mulai dari contoh yang
                                sudah disiapkan.
                            </p>
                            <div className="mt-1 flex flex-wrap justify-center gap-2">
                                <button
                                    type="button"
                                    onClick={() =>
                                        setDialog({ open: true, agent: null })
                                    }
                                    className="inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700"
                                >
                                    <Plus className="size-4" />
                                    Buat agent pertama
                                </button>
                                <Link
                                    href={setup()}
                                    className="inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent"
                                >
                                    Pasang CLI
                                </Link>
                            </div>
                        </section>
                    ) : (
                        <div className="grid gap-4 md:grid-cols-2 2xl:grid-cols-3">
                            {agents.map((agent) => (
                                <article
                                    key={agent.id}
                                    className="flex flex-col border border-sidebar-border bg-card"
                                >
                                    <div className="flex-1 space-y-3 p-5">
                                        <div className="flex items-start justify-between gap-3">
                                            <div className="min-w-0">
                                                <h2 className="truncate font-semibold">
                                                    {agent.name}
                                                </h2>
                                                <p className="mt-0.5 line-clamp-2 text-sm text-muted-foreground">
                                                    {agent.description ||
                                                        'Tanpa deskripsi.'}
                                                </p>
                                            </div>
                                            <div className="flex shrink-0 flex-col items-end gap-1.5">
                                                <ModeBadge mode={agent.mode} />
                                                <EffortBadge
                                                    effort={agent.effort}
                                                />
                                            </div>
                                        </div>
                                        <pre className="max-h-32 overflow-hidden border border-sidebar-border bg-muted/40 p-3 font-mono text-xs leading-relaxed break-words whitespace-pre-wrap text-muted-foreground">
                                            {agent.instructions}
                                        </pre>
                                    </div>
                                    <div className="flex items-center justify-between gap-2 border-t border-sidebar-border px-5 py-3">
                                        <span className="text-xs text-muted-foreground">
                                            Diubah{' '}
                                            {formatRelative(agent.updated_at)}
                                        </span>
                                        <div className="flex gap-2">
                                            <button
                                                type="button"
                                                onClick={() =>
                                                    setDialog({
                                                        open: true,
                                                        agent,
                                                    })
                                                }
                                                className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-xs font-medium transition-colors hover:bg-sidebar-accent"
                                            >
                                                <Pencil className="size-3.5" />
                                                Ubah
                                            </button>
                                            <button
                                                type="button"
                                                onClick={() =>
                                                    setDeleting(agent)
                                                }
                                                className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-xs font-medium transition-colors hover:border-red-600 hover:text-red-600"
                                            >
                                                <Trash2 className="size-3.5" />
                                                Hapus
                                            </button>
                                        </div>
                                    </div>
                                </article>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            <AgentDialog
                open={dialog.open}
                agent={dialog.agent}
                modes={modes}
                efforts={efforts}
                onClose={() => setDialog({ open: false, agent: null })}
            />

            <Dialog
                open={deleting !== null}
                onOpenChange={(o) => !o && setDeleting(null)}
            >
                <DialogContent className="rounded-none sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Hapus agent?</DialogTitle>
                        <DialogDescription>
                            Agent <strong>{deleting?.name}</strong> akan dihapus
                            permanen. Percakapan lama tidak terpengaruh.
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <button
                            type="button"
                            onClick={() => setDeleting(null)}
                            className="inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent"
                        >
                            Batal
                        </button>
                        <button
                            type="button"
                            onClick={() => {
                                if (deleting) {
                                    router.delete(
                                        AgentController.destroy.url(
                                            deleting.id,
                                        ),
                                        {
                                            preserveScroll: true,
                                            onSuccess: () => setDeleting(null),
                                        },
                                    );
                                }
                            }}
                            className="inline-flex h-9 items-center bg-red-600 px-4 text-sm font-medium text-white transition-colors hover:bg-red-700"
                        >
                            Hapus agent
                        </button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}

Agents.layout = {
    breadcrumbs: [{ title: 'Agents', href: agentsRoute() }],
};

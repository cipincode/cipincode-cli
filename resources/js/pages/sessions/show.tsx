import { Head, Link } from '@inertiajs/react';
import { ArrowLeft, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { formatDateTime } from '@/components/api-keys/format';
import { ModeBadge } from '@/components/mode-badge';
import DeleteSessionDialog from '@/components/sessions/delete-session-dialog';
import { sessions } from '@/routes';

type Session = {
    uuid: string;
    title: string;
    mode: string;
    model: string | null;
    messages_count: number;
    created_at: string;
    last_activity_at: string;
};

type Message = {
    id: number;
    role: 'user' | 'assistant';
    content: string;
    created_at: string;
};

/** Pecah teks pada blok ``` agar kode tampil sebagai kotak tersendiri. */
function Content({ text }: { text: string }) {
    const parts = text.split(/```/);

    return (
        <div className="space-y-3 text-sm leading-relaxed">
            {parts.map((part, i) => {
                if (i % 2 === 0) {
                    return part.trim() === '' ? null : (
                        <p key={i} className="break-words whitespace-pre-wrap">
                            {part.trim()}
                        </p>
                    );
                }
                const newline = part.indexOf('\n');
                const lang = newline > -1 ? part.slice(0, newline).trim() : '';
                const code = newline > -1 ? part.slice(newline + 1) : part;

                return (
                    <figure key={i} className="border border-sidebar-border">
                        <figcaption className="border-b border-sidebar-border bg-muted/40 px-3 py-1.5 text-xs font-semibold text-orange-600">
                            {lang || 'kode'}
                        </figcaption>
                        <pre className="overflow-x-auto bg-zinc-950 p-4 font-mono text-[13px] leading-relaxed text-zinc-100">
                            {code.replace(/\n$/, '')}
                        </pre>
                    </figure>
                );
            })}
        </div>
    );
}

export default function SessionShow({
    session,
    messages,
}: {
    session: Session;
    messages: Message[];
}) {
    const [confirming, setConfirming] = useState(false);

    return (
        <>
            <Head title={session.title} />

            <div className="flex flex-1 flex-col">
                <div className="flex flex-wrap items-start justify-between gap-4 border-b border-sidebar-border px-6 py-5">
                    <div className="min-w-0">
                        <Link
                            href={sessions()}
                            className="mb-2 inline-flex items-center gap-1.5 text-xs text-muted-foreground transition-colors hover:text-foreground"
                        >
                            <ArrowLeft className="size-3.5" />
                            Semua sesi
                        </Link>
                        <h1 className="truncate text-xl font-semibold tracking-tight">
                            {session.title}
                        </h1>
                        <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
                            <ModeBadge mode={session.mode} />
                            <span>{session.messages_count} pesan</span>
                            <span>
                                Dimulai {formatDateTime(session.created_at)}
                            </span>
                            {session.model && (
                                <span className="font-mono">
                                    {session.model}
                                </span>
                            )}
                        </div>
                    </div>
                    <button
                        type="button"
                        onClick={() => setConfirming(true)}
                        className="inline-flex h-9 items-center gap-2 border border-sidebar-border px-4 text-sm font-medium transition-colors hover:border-red-600 hover:text-red-600"
                    >
                        <Trash2 className="size-4" />
                        Hapus sesi
                    </button>
                </div>

                <div className="flex w-full flex-1 flex-col gap-4 p-6">
                    {messages.length === 0 ? (
                        <p className="border border-dashed border-sidebar-border p-6 text-center text-sm text-muted-foreground">
                            Sesi ini belum memiliki pesan.
                        </p>
                    ) : (
                        messages.map((m) => (
                            <article
                                key={m.id}
                                className={`border ${m.role === 'user' ? 'border-sidebar-border bg-muted/30' : 'border-sidebar-border bg-card'}`}
                            >
                                <header className="flex items-center justify-between border-b border-sidebar-border px-4 py-2 text-xs">
                                    <span
                                        className={`font-semibold tracking-wide uppercase ${m.role === 'user' ? 'text-muted-foreground' : 'text-orange-600'}`}
                                    >
                                        {m.role === 'user'
                                            ? 'Anda'
                                            : 'CipinCode'}
                                    </span>
                                    <time
                                        className="text-muted-foreground"
                                        dateTime={m.created_at}
                                    >
                                        {formatDateTime(m.created_at)}
                                    </time>
                                </header>
                                <div className="p-4">
                                    <Content text={m.content} />
                                </div>
                            </article>
                        ))
                    )}
                </div>
            </div>

            <DeleteSessionDialog
                session={confirming ? session : null}
                onClose={() => setConfirming(false)}
            />
        </>
    );
}

SessionShow.layout = {
    breadcrumbs: [{ title: 'Sessions', href: sessions() }],
};

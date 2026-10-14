import { Form } from '@inertiajs/react';
import { useState } from 'react';
import AgentController from '@/actions/App/Http/Controllers/AgentController';
import InputError from '@/components/input-error';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Spinner } from '@/components/ui/spinner';

export type Agent = {
    id: number;
    name: string;
    description: string | null;
    instructions: string;
    mode: string;
    effort: string;
    updated_at: string;
};

type Draft = Pick<
    Agent,
    'name' | 'description' | 'instructions' | 'mode' | 'effort'
>;

const TEMPLATES: { label: string; draft: Draft }[] = [
    {
        label: 'Code reviewer',
        draft: {
            name: 'Code Reviewer',
            description: 'Meninjau kode dan menunjukkan masalah yang nyata.',
            instructions:
                'Kamu meninjau kode seperti rekan senior. Cari bug, risiko keamanan, dan masalah performa yang nyata. Urutkan temuan dari yang paling serius, sebutkan baris atau fungsinya, dan beri perbaikan singkat. Jangan mengomentari gaya penulisan kecuali mengganggu keterbacaan.',
            mode: 'coding',
            effort: 'high',
        },
    },
    {
        label: 'Penjelas sederhana',
        draft: {
            name: 'Penjelas Sederhana',
            description: 'Menjelaskan konsep teknis dengan bahasa awam.',
            instructions:
                'Jelaskan konsep teknis dengan bahasa sehari-hari, pakai analogi, dan hindari istilah sulit kecuali sudah dijelaskan. Akhiri dengan satu contoh singkat.',
            mode: 'chat',
            effort: 'medium',
        },
    },
    {
        label: 'Pemburu bug',
        draft: {
            name: 'Pemburu Bug',
            description: 'Menelusuri penyebab error langkah demi langkah.',
            instructions:
                'Saat diberi pesan error atau perilaku salah, tanyakan informasi yang kurang, lalu telusuri dugaan penyebab dari yang paling mungkin. Beri langkah verifikasi untuk setiap dugaan sebelum menyarankan perbaikan.',
            mode: 'debugging',
            effort: 'xhigh',
        },
    },
];

const field =
    'w-full border border-input bg-transparent px-3 py-2 text-sm outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50';
const primary =
    'inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:opacity-60';
const secondary =
    'inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent';

const EMPTY: Draft = {
    name: '',
    description: '',
    instructions: '',
    mode: 'chat',
    effort: 'medium',
};
const MODE_LABEL: Record<string, string> = {
    chat: 'Chat',
    coding: 'Coding',
    debugging: 'Debugging',
};

export default function AgentDialog({
    open,
    agent,
    modes,
    efforts,
    onClose,
}: {
    open: boolean;
    agent: Agent | null;
    modes: string[];
    efforts: string[];
    onClose: () => void;
}) {
    // key mengganti isi formulir setiap kali agent atau templat berganti
    const [template, setTemplate] = useState<{
        n: number;
        draft: Draft;
    } | null>(null);
    const draft: Draft = template?.draft ?? agent ?? EMPTY;
    const editing = agent !== null;

    const close = () => {
        setTemplate(null);
        onClose();
    };

    return (
        <Dialog open={open} onOpenChange={(o) => !o && close()}>
            <DialogContent className="max-h-[92vh] overflow-y-auto rounded-none sm:max-w-xl">
                <DialogHeader>
                    <DialogTitle>
                        {editing ? 'Ubah agent' : 'Buat agent'}
                    </DialogTitle>
                    <DialogDescription>
                        Agent menambahkan instruksi khusus ke setiap percakapan.
                        Pilih di CLI dengan perintah /agent.
                    </DialogDescription>
                </DialogHeader>

                {!editing && (
                    <div className="flex flex-wrap items-center gap-2 text-xs">
                        <span className="text-muted-foreground">
                            Mulai dari contoh:
                        </span>
                        {TEMPLATES.map((t) => (
                            <button
                                key={t.label}
                                type="button"
                                onClick={() =>
                                    setTemplate({
                                        n: (template?.n ?? 0) + 1,
                                        draft: t.draft,
                                    })
                                }
                                className="border border-sidebar-border px-2.5 py-1 font-medium transition-colors hover:border-orange-600 hover:text-orange-600"
                            >
                                {t.label}
                            </button>
                        ))}
                    </div>
                )}

                <Form
                    key={`${agent?.id ?? 'new'}-${template?.n ?? 0}-${open}`}
                    {...(editing
                        ? AgentController.update.form(agent.id)
                        : AgentController.store.form())}
                    onSuccess={close}
                    className="grid gap-4"
                >
                    {({ processing, errors }) => (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="agent-name">Nama</Label>
                                <Input
                                    id="agent-name"
                                    name="name"
                                    required
                                    maxLength={50}
                                    defaultValue={draft.name}
                                    placeholder="Contoh: Code Reviewer"
                                    className="rounded-none"
                                />
                                <InputError message={errors.name} />
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="agent-description">
                                    Deskripsi singkat
                                </Label>
                                <Input
                                    id="agent-description"
                                    name="description"
                                    maxLength={160}
                                    defaultValue={draft.description ?? ''}
                                    placeholder="Untuk apa agent ini dipakai"
                                    className="rounded-none"
                                />
                                <InputError message={errors.description} />
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="agent-instructions">
                                    Instruksi
                                </Label>
                                <textarea
                                    id="agent-instructions"
                                    name="instructions"
                                    required
                                    rows={7}
                                    maxLength={4000}
                                    defaultValue={draft.instructions}
                                    placeholder="Jelaskan peran, gaya jawaban, dan batasan agent ini."
                                    className={`${field} resize-y font-mono text-[13px] leading-relaxed`}
                                />
                                <InputError message={errors.instructions} />
                            </div>

                            <div className="grid gap-4 sm:grid-cols-2">
                                <div className="grid gap-2">
                                    <Label htmlFor="agent-mode">
                                        Mode bawaan
                                    </Label>
                                    <select
                                        id="agent-mode"
                                        name="mode"
                                        defaultValue={draft.mode}
                                        className={`${field} h-9 py-0`}
                                    >
                                        {modes.map((m) => (
                                            <option key={m} value={m}>
                                                {MODE_LABEL[m] ?? m}
                                            </option>
                                        ))}
                                    </select>
                                    <InputError message={errors.mode} />
                                </div>
                                <div className="grid gap-2">
                                    <Label htmlFor="agent-effort">
                                        Effort bawaan
                                    </Label>
                                    <select
                                        id="agent-effort"
                                        name="effort"
                                        defaultValue={draft.effort}
                                        className={`${field} h-9 py-0`}
                                    >
                                        {efforts.map((e) => (
                                            <option key={e} value={e}>
                                                {e}
                                            </option>
                                        ))}
                                    </select>
                                    <InputError message={errors.effort} />
                                </div>
                            </div>

                            <DialogFooter>
                                <button
                                    type="button"
                                    className={secondary}
                                    onClick={close}
                                >
                                    Batal
                                </button>
                                <button
                                    type="submit"
                                    disabled={processing}
                                    className={primary}
                                >
                                    {processing && <Spinner />}
                                    {editing ? 'Simpan' : 'Buat agent'}
                                </button>
                            </DialogFooter>
                        </>
                    )}
                </Form>
            </DialogContent>
        </Dialog>
    );
}

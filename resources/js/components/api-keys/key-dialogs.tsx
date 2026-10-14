import { Form } from '@inertiajs/react';
import ApiKeyController from '@/actions/App/Http/Controllers/ApiKeyController';
import InputError from '@/components/input-error';
import CopyCommand from '@/components/setup/copy-command';
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
import type { ApiKey } from './types';

const content = 'rounded-none sm:max-w-md';
const primary =
    'inline-flex h-9 items-center gap-2 bg-orange-600 px-4 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:opacity-60';
const secondary =
    'inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent';
const danger =
    'inline-flex h-9 items-center gap-2 bg-red-600 px-4 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-60';

export function CreateKeyDialog({
    open,
    onOpenChange,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}) {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className={content}>
                <DialogHeader>
                    <DialogTitle>Buat API key</DialogTitle>
                    <DialogDescription>
                        Beri nama agar mudah dikenali, misalnya nama proyek atau
                        tujuan pemakaiannya.
                    </DialogDescription>
                </DialogHeader>
                <Form
                    {...ApiKeyController.store.form()}
                    resetOnSuccess
                    onSuccess={() => onOpenChange(false)}
                    className="grid gap-4"
                >
                    {({ processing, errors }) => (
                        <>
                            <div className="grid gap-2">
                                <Label htmlFor="create-name">Nama key</Label>
                                <Input
                                    id="create-name"
                                    name="name"
                                    required
                                    autoFocus
                                    maxLength={50}
                                    placeholder="Contoh: Coding proyek web"
                                    className="rounded-none"
                                />
                                <InputError message={errors.name} />
                            </div>
                            <DialogFooter>
                                <button
                                    type="button"
                                    className={secondary}
                                    onClick={() => onOpenChange(false)}
                                >
                                    Batal
                                </button>
                                <button
                                    type="submit"
                                    disabled={processing}
                                    className={primary}
                                >
                                    {processing && <Spinner />}
                                    Buat key
                                </button>
                            </DialogFooter>
                        </>
                    )}
                </Form>
            </DialogContent>
        </Dialog>
    );
}

export function EditKeyDialog({
    apiKey,
    onClose,
}: {
    apiKey: ApiKey | null;
    onClose: () => void;
}) {
    return (
        <Dialog open={apiKey !== null} onOpenChange={(o) => !o && onClose()}>
            <DialogContent className={content}>
                <DialogHeader>
                    <DialogTitle>Ubah nama key</DialogTitle>
                    <DialogDescription>
                        Mengubah nama tidak mengubah key. CLI yang memakainya
                        tetap bekerja.
                    </DialogDescription>
                </DialogHeader>
                {apiKey && (
                    <Form
                        key={apiKey.id}
                        {...ApiKeyController.update.form(apiKey.id)}
                        onSuccess={onClose}
                        className="grid gap-4"
                    >
                        {({ processing, errors }) => (
                            <>
                                <div className="grid gap-2">
                                    <Label htmlFor="edit-name">Nama key</Label>
                                    <Input
                                        id="edit-name"
                                        name="name"
                                        required
                                        autoFocus
                                        maxLength={50}
                                        defaultValue={apiKey.name}
                                        className="rounded-none"
                                    />
                                    <InputError message={errors.name} />
                                </div>
                                <DialogFooter>
                                    <button
                                        type="button"
                                        className={secondary}
                                        onClick={onClose}
                                    >
                                        Batal
                                    </button>
                                    <button
                                        type="submit"
                                        disabled={processing}
                                        className={primary}
                                    >
                                        {processing && <Spinner />}
                                        Simpan
                                    </button>
                                </DialogFooter>
                            </>
                        )}
                    </Form>
                )}
            </DialogContent>
        </Dialog>
    );
}

export function DeleteKeyDialog({
    apiKey,
    onClose,
}: {
    apiKey: ApiKey | null;
    onClose: () => void;
}) {
    return (
        <Dialog open={apiKey !== null} onOpenChange={(o) => !o && onClose()}>
            <DialogContent className={content}>
                <DialogHeader>
                    <DialogTitle>Hapus API key?</DialogTitle>
                    <DialogDescription>
                        Key <strong>{apiKey?.name}</strong> akan dihapus
                        permanen. CLI yang memakainya akan berhenti bekerja dan
                        tindakan ini tidak bisa dibatalkan. Untuk menahan
                        sementara, nonaktifkan saja.
                    </DialogDescription>
                </DialogHeader>
                {apiKey && (
                    <Form
                        key={apiKey.id}
                        {...ApiKeyController.destroy.form(apiKey.id)}
                        onSuccess={onClose}
                    >
                        {({ processing }) => (
                            <DialogFooter>
                                <button
                                    type="button"
                                    className={secondary}
                                    onClick={onClose}
                                >
                                    Batal
                                </button>
                                <button
                                    type="submit"
                                    disabled={processing}
                                    className={danger}
                                >
                                    {processing && <Spinner />}
                                    Hapus key
                                </button>
                            </DialogFooter>
                        )}
                    </Form>
                )}
            </DialogContent>
        </Dialog>
    );
}

export function RevealKeyDialog({
    newKey,
    open,
    onOpenChange,
}: {
    newKey: string | null;
    open: boolean;
    onOpenChange: (open: boolean) => void;
}) {
    return (
        <Dialog open={open && !!newKey} onOpenChange={onOpenChange}>
            <DialogContent className="w-fit max-w-[calc(100%-2rem)] rounded-none sm:max-w-[calc(100%-2rem)]">
                <DialogHeader className="max-w-xl">
                    <DialogTitle>API key berhasil dibuat</DialogTitle>
                    <DialogDescription>
                        Salin key ini sekarang. Demi keamanan, key hanya
                        ditampilkan sekali dan tidak bisa dilihat lagi setelah
                        dialog ditutup.
                    </DialogDescription>
                </DialogHeader>
                {newKey && (
                    <div className="grid gap-4">
                        <div className="grid gap-2">
                            <Label>Key Anda</Label>
                            <CopyCommand command={newKey} prompt={false} wrap />
                        </div>
                        <div className="grid gap-2">
                            <Label>Masuk dari terminal</Label>
                            <CopyCommand
                                command={`cipincode login --key ${newKey}`}
                                wrap
                            />
                        </div>
                    </div>
                )}
                <DialogFooter>
                    <button
                        type="button"
                        className={primary}
                        onClick={() => onOpenChange(false)}
                    >
                        Saya sudah menyalinnya
                    </button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}

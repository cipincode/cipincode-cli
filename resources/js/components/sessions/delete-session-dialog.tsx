import { Form } from '@inertiajs/react';
import SessionController from '@/actions/App/Http/Controllers/SessionController';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Spinner } from '@/components/ui/spinner';

export default function DeleteSessionDialog({
    session,
    onClose,
}: {
    session: { uuid: string; title: string } | null;
    onClose: () => void;
}) {
    return (
        <Dialog
            open={session !== null}
            onOpenChange={(open) => !open && onClose()}
        >
            <DialogContent className="rounded-none sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>Hapus sesi?</DialogTitle>
                    <DialogDescription>
                        Sesi <strong>{session?.title}</strong> beserta seluruh
                        pesannya akan dihapus permanen dan tidak bisa
                        dikembalikan.
                    </DialogDescription>
                </DialogHeader>
                {session && (
                    <Form
                        key={session.uuid}
                        {...SessionController.destroy.form(session.uuid)}
                        onSuccess={onClose}
                    >
                        {({ processing }) => (
                            <DialogFooter>
                                <button
                                    type="button"
                                    onClick={onClose}
                                    className="inline-flex h-9 items-center border border-sidebar-border px-4 text-sm font-medium transition-colors hover:bg-sidebar-accent"
                                >
                                    Batal
                                </button>
                                <button
                                    type="submit"
                                    disabled={processing}
                                    className="inline-flex h-9 items-center gap-2 bg-red-600 px-4 text-sm font-medium text-white transition-colors hover:bg-red-700 disabled:opacity-60"
                                >
                                    {processing && <Spinner />}
                                    Hapus sesi
                                </button>
                            </DialogFooter>
                        )}
                    </Form>
                )}
            </DialogContent>
        </Dialog>
    );
}

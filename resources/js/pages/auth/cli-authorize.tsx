import { Form, Head, Link } from '@inertiajs/react';
import CliAuthController from '@/actions/App/Http/Controllers/CliAuthController';
import CopyCommand from '@/components/setup/copy-command';
import { Spinner } from '@/components/ui/spinner';
import { dashboard } from '@/routes';

type Props = {
    challenge: string;
    device: string | null;
    user: { name: string; email: string };
    code: string | null;
    ttl: number;
};

export default function CliAuthorize({
    challenge,
    device,
    user,
    code,
    ttl,
}: Props) {
    return (
        <>
            <Head title="Hubungkan CLI" />

            <div className="border border-sidebar-border">
                <div className="border-b border-sidebar-border px-4 py-3 text-sm">
                    <p className="text-muted-foreground">Masuk sebagai</p>
                    <p className="font-medium">{user.name}</p>
                    <p className="text-muted-foreground">{user.email}</p>
                </div>
                {device && (
                    <div className="px-4 py-3 text-sm">
                        <p className="text-muted-foreground">Perangkat</p>
                        <p className="font-medium">{device}</p>
                    </div>
                )}
            </div>

            {code ? (
                <div className="flex flex-col gap-4">
                    <CopyCommand command={code} prompt={false} />
                    <p className="text-center text-sm text-muted-foreground">
                        Salin kode ini, lalu tempel di terminal Anda. Kode
                        berlaku {ttl} menit dan hanya bisa dipakai sekali.
                    </p>
                    <Link
                        href={dashboard()}
                        className="inline-flex h-9 items-center justify-center border border-sidebar-border text-sm font-medium transition-colors hover:bg-sidebar-accent"
                    >
                        Selesai
                    </Link>
                </div>
            ) : (
                <Form
                    {...CliAuthController.approve.form()}
                    className="flex flex-col gap-3"
                >
                    {({ processing }) => (
                        <>
                            <input
                                type="hidden"
                                name="challenge"
                                value={challenge}
                            />
                            {device && (
                                <input
                                    type="hidden"
                                    name="device"
                                    value={device}
                                />
                            )}
                            <button
                                type="submit"
                                disabled={processing}
                                className="inline-flex h-9 items-center justify-center gap-2 bg-orange-600 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:opacity-60"
                            >
                                {processing && <Spinner />}
                                Izinkan
                            </button>
                            <Link
                                href={dashboard()}
                                className="inline-flex h-9 items-center justify-center border border-sidebar-border text-sm font-medium transition-colors hover:bg-sidebar-accent"
                            >
                                Batal
                            </Link>
                        </>
                    )}
                </Form>
            )}
        </>
    );
}

CliAuthorize.layout = {
    title: 'Hubungkan CipinCode CLI',
    description:
        'CLI di terminal Anda meminta izin untuk memakai akun ini. Izinkan hanya jika Anda yang menjalankannya.',
};

import { Form, Head, router, usePage } from '@inertiajs/react';
import { Camera, Mail, Trash2, UserRound } from 'lucide-react';
import { useRef, useState } from 'react';
import ProfileController from '@/actions/App/Http/Controllers/Settings/ProfileController';
import { formatDate } from '@/components/api-keys/format';
import DeleteUser from '@/components/delete-user';
import InputError from '@/components/input-error';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { avatar as avatarRoute, edit } from '@/routes/profile';
import type { Auth } from '@/types';

type PageProps = {
    auth: Auth;
};

const field =
    'h-10 rounded-none pl-10 focus-visible:border-orange-600 focus-visible:ring-1 focus-visible:ring-orange-600/40';

export default function Profile() {
    const { auth } = usePage<PageProps>().props;
    const user = auth.user;
    const initial = (user.name || user.email).trim().charAt(0).toUpperCase();
    const fileInput = useRef<HTMLInputElement>(null);
    const [uploading, setUploading] = useState(false);
    const [avatarError, setAvatarError] = useState<string | null>(null);

    const upload = (file?: File) => {
        if (!file) return;
        setAvatarError(null);
        router.post(
            avatarRoute.url(),
            { avatar: file },
            {
                forceFormData: true,
                preserveScroll: true,
                onStart: () => setUploading(true),
                onError: (e) => setAvatarError(e.avatar ?? 'Gagal mengunggah.'),
                onFinish: () => {
                    setUploading(false);
                    if (fileInput.current) fileInput.current.value = '';
                },
            },
        );
    };

    return (
        <>
            <Head title="Profile settings" />

            <section className="flex flex-wrap items-center gap-5 border border-sidebar-border bg-card p-5">
                <div className="relative size-20 shrink-0 overflow-hidden bg-orange-600">
                    {user.avatar ? (
                        <img
                            src={user.avatar}
                            alt={user.name}
                            className="size-full object-cover"
                        />
                    ) : (
                        <span className="flex size-full items-center justify-center text-3xl font-semibold text-white">
                            {initial}
                        </span>
                    )}
                    {uploading && (
                        <span className="absolute inset-0 flex items-center justify-center bg-black/50 text-xs text-white">
                            Mengunggah…
                        </span>
                    )}
                </div>
                <div className="min-w-0 flex-1">
                    <p className="truncate text-lg font-semibold">
                        {user.name}
                    </p>
                    <p className="truncate text-sm text-muted-foreground">
                        {user.email}
                    </p>
                    {user.created_at && (
                        <p className="mt-0.5 text-xs text-muted-foreground">
                            Bergabung sejak {formatDate(user.created_at)}
                        </p>
                    )}
                    <div className="mt-3 flex flex-wrap items-center gap-2">
                        <input
                            ref={fileInput}
                            type="file"
                            accept="image/jpeg,image/png,image/webp"
                            className="hidden"
                            onChange={(e) => upload(e.target.files?.[0])}
                        />
                        <button
                            type="button"
                            disabled={uploading}
                            onClick={() => fileInput.current?.click()}
                            className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-sm hover:bg-sidebar-accent disabled:opacity-60"
                        >
                            <Camera className="size-4" />
                            {user.avatar ? 'Ganti foto' : 'Unggah foto'}
                        </button>
                        {user.avatar && (
                            <button
                                type="button"
                                disabled={uploading}
                                onClick={() =>
                                    router.delete(avatarRoute.url(), {
                                        preserveScroll: true,
                                    })
                                }
                                className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-sm text-red-600 hover:bg-red-600/10 disabled:opacity-60"
                            >
                                <Trash2 className="size-4" />
                                Hapus
                            </button>
                        )}
                    </div>
                    <p className="mt-2 text-xs text-muted-foreground">
                        JPG, PNG, atau WebP, maks 5 MB. Ukuran file otomatis
                        dikompres, resolusi tetap.
                    </p>
                    {avatarError && (
                        <p className="mt-1 text-xs text-red-600">
                            {avatarError}
                        </p>
                    )}
                </div>
            </section>

            <section className="border border-sidebar-border bg-card">
                <div className="border-b border-sidebar-border p-5">
                    <h2 className="font-medium">Informasi profil</h2>
                    <p className="mt-0.5 text-sm text-muted-foreground">
                        Perbarui nama dan alamat email Anda.
                    </p>
                </div>

                <Form
                    {...ProfileController.update.form()}
                    options={{ preserveScroll: true }}
                >
                    {({ processing, errors, recentlySuccessful }) => (
                        <>
                            <div className="space-y-5 p-5">
                                <div className="grid gap-2">
                                    <Label htmlFor="name">Nama</Label>
                                    <div className="relative">
                                        <UserRound className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                                        <Input
                                            id="name"
                                            className={field}
                                            defaultValue={user.name}
                                            name="name"
                                            required
                                            autoComplete="name"
                                            placeholder="Nama lengkap"
                                        />
                                    </div>
                                    <InputError message={errors.name} />
                                </div>

                                <div className="grid gap-2">
                                    <Label htmlFor="email">Alamat email</Label>
                                    <div className="relative">
                                        <Mail className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                                        <Input
                                            id="email"
                                            type="email"
                                            className={field}
                                            defaultValue={user.email}
                                            name="email"
                                            required
                                            autoComplete="username"
                                            placeholder="nama@email.com"
                                        />
                                    </div>
                                    <InputError message={errors.email} />
                                </div>
                            </div>

                            <div className="flex items-center justify-end gap-4 border-t border-sidebar-border px-5 py-3">
                                {recentlySuccessful && (
                                    <span className="text-sm text-emerald-600">
                                        Tersimpan.
                                    </span>
                                )}
                                <button
                                    type="submit"
                                    disabled={processing}
                                    data-test="update-profile-button"
                                    className="inline-flex h-9 items-center bg-orange-600 px-5 text-sm font-medium text-white transition-colors hover:bg-orange-700 disabled:opacity-60"
                                >
                                    {processing ? 'Menyimpan…' : 'Simpan'}
                                </button>
                            </div>
                        </>
                    )}
                </Form>
            </section>

            <DeleteUser />
        </>
    );
}

Profile.layout = {
    breadcrumbs: [
        {
            title: 'Profile settings',
            href: edit(),
        },
    ],
};

import { Head, router, usePage } from '@inertiajs/react';
import { useState } from 'react';
import { Plus, Search, Trash2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

type UserRow = {
    id: number;
    name: string;
    email: string;
    plan: string;
    plan_label: string;
    is_admin: boolean;
    sessions_count: number;
    messages_total: number;
    created_at: string;
};

type Paginated = {
    data: UserRow[];
    current_page: number;
    last_page: number;
    total: number;
};

type Props = {
    users: Paginated;
    filters: { q: string; plan: string | null };
    plans: string[];
};

function UserEditor({
    user,
    plans,
    isSelf,
}: {
    user: UserRow;
    plans: string[];
    isSelf: boolean;
}) {
    const [plan, setPlan] = useState(user.plan);
    const [isAdmin, setIsAdmin] = useState(user.is_admin);
    const dirty = plan !== user.plan || isAdmin !== user.is_admin;

    return (
        <tr className="border-t border-sidebar-border">
            <td className="px-4 py-3">
                <p className="font-medium">{user.name}</p>
                <p className="text-xs text-muted-foreground">{user.email}</p>
            </td>
            <td className="px-4 py-3">
                <select
                    value={plan}
                    onChange={(e) => setPlan(e.target.value)}
                    className="h-9 border border-sidebar-border bg-background px-2 text-sm"
                    aria-label={`Paket ${user.email}`}
                >
                    {plans.map((p) => (
                        <option key={p} value={p}>
                            {p}
                        </option>
                    ))}
                </select>
            </td>
            <td className="px-4 py-3 text-center">
                <Checkbox
                    checked={isAdmin}
                    disabled={isSelf}
                    onCheckedChange={(v) => setIsAdmin(v === true)}
                    aria-label={`Admin ${user.email}`}
                    title={
                        isSelf
                            ? 'Tidak bisa mencabut admin diri sendiri'
                            : undefined
                    }
                />
            </td>
            <td className="px-4 py-3 text-right tabular-nums">
                {user.messages_total}
            </td>
            <td className="px-4 py-3 text-right tabular-nums">
                {user.sessions_count}
            </td>
            <td className="px-4 py-3 text-right">
                <div className="flex justify-end gap-2">
                    <Button
                        size="sm"
                        variant="outline"
                        disabled={!dirty}
                        onClick={() =>
                            router.patch(
                                `/spr/users/${user.id}`,
                                { plan, is_admin: isAdmin },
                                { preserveScroll: true },
                            )
                        }
                    >
                        Simpan
                    </Button>
                    <Button
                        size="sm"
                        variant="outline"
                        disabled={isSelf}
                        title={
                            isSelf
                                ? 'Tidak bisa menghapus akun sendiri'
                                : `Hapus ${user.email}`
                        }
                        onClick={() => {
                            if (
                                !window.confirm(
                                    `Hapus ${user.name} (${user.email}) beserta seluruh datanya?`,
                                )
                            )
                                return;
                            router.delete(`/spr/users/${user.id}`, {
                                preserveScroll: true,
                            });
                        }}
                    >
                        <Trash2 className="size-3.5 text-red-600" />
                        Hapus
                    </Button>
                </div>
            </td>
        </tr>
    );
}

export default function SprUsers({ users, filters, plans }: Props) {
    const { auth, errors } = usePage().props;
    const [q, setQ] = useState(filters.q);
    const [plan, setPlan] = useState(filters.plan ?? '');
    const [createOpen, setCreateOpen] = useState(false);
    const [form, setForm] = useState({
        name: '',
        email: '',
        password: '',
        plan: 'free',
        is_admin: false,
    });

    const set = (key: keyof typeof form, value: string | boolean) =>
        setForm((f) => ({ ...f, [key]: value }));

    const submitCreate = (e: React.FormEvent) => {
        e.preventDefault();
        router.post('/spr/users', form, {
            preserveScroll: true,
            onSuccess: () => {
                setCreateOpen(false);
                setForm({
                    name: '',
                    email: '',
                    password: '',
                    plan: 'free',
                    is_admin: false,
                });
            },
        });
    };

    const search = (page = 1) => {
        router.get(
            '/spr/users',
            { q, plan: plan || undefined, page },
            { preserveState: true, replace: true },
        );
    };

    const gotoPage = (page: number) => {
        router.get(
            '/spr/users',
            { q: filters.q || undefined, plan: filters.plan ?? undefined, page },
            { preserveState: true, replace: true },
        );
    };

    return (
        <>
            <Head title="Pengguna" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                        <h1 className="text-xl font-semibold tracking-tight">
                            Pengguna
                        </h1>
                        <p className="text-sm text-muted-foreground">
                            {users.total} pengguna terdaftar. Tambah, ubah
                            paket, atau hapus pengguna.
                        </p>
                    </div>
                    <Button
                        onClick={() => setCreateOpen(true)}
                        className="bg-orange-600 text-white hover:bg-orange-700"
                    >
                        <Plus className="size-4" />
                        Tambah Pengguna
                    </Button>
                </div>

                <Dialog open={createOpen} onOpenChange={setCreateOpen}>
                    <DialogContent className="sm:max-w-md">
                        <DialogHeader>
                            <DialogTitle>Tambah Pengguna</DialogTitle>
                            <DialogDescription>
                                Buat akun baru. Paket gratis otomatis dapat
                                masa percobaan 14 hari.
                            </DialogDescription>
                        </DialogHeader>
                        <form
                            onSubmit={submitCreate}
                            className="flex flex-col gap-4"
                        >
                            <div className="flex flex-col gap-1.5">
                                <Label htmlFor="new-name">Nama</Label>
                                <Input
                                    id="new-name"
                                    value={form.name}
                                    onChange={(e) =>
                                        set('name', e.target.value)
                                    }
                                    placeholder="Nama lengkap"
                                />
                                {errors.name && (
                                    <p className="text-xs text-red-600">
                                        {errors.name}
                                    </p>
                                )}
                            </div>
                            <div className="flex flex-col gap-1.5">
                                <Label htmlFor="new-email">Email</Label>
                                <Input
                                    id="new-email"
                                    type="email"
                                    value={form.email}
                                    onChange={(e) =>
                                        set('email', e.target.value)
                                    }
                                    placeholder="nama@email.com"
                                />
                                {errors.email && (
                                    <p className="text-xs text-red-600">
                                        {errors.email}
                                    </p>
                                )}
                            </div>
                            <div className="flex flex-col gap-1.5">
                                <Label htmlFor="new-password">
                                    Kata sandi
                                </Label>
                                <Input
                                    id="new-password"
                                    type="password"
                                    value={form.password}
                                    onChange={(e) =>
                                        set('password', e.target.value)
                                    }
                                    placeholder="Minimal 8 karakter"
                                />
                                {errors.password && (
                                    <p className="text-xs text-red-600">
                                        {errors.password}
                                    </p>
                                )}
                            </div>
                            <div className="flex gap-4">
                                <div className="flex flex-1 flex-col gap-1.5">
                                    <Label htmlFor="new-plan">Paket</Label>
                                    <select
                                        id="new-plan"
                                        value={form.plan}
                                        onChange={(e) =>
                                            set('plan', e.target.value)
                                        }
                                        className="h-9 border border-sidebar-border bg-background px-2 text-sm"
                                    >
                                        {plans.map((p) => (
                                            <option key={p} value={p}>
                                                {p}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                                <label
                                    htmlFor="new-admin"
                                    className="flex cursor-pointer items-end gap-2 pb-2 text-sm select-none"
                                >
                                    <Checkbox
                                        id="new-admin"
                                        checked={form.is_admin}
                                        onCheckedChange={(v) =>
                                            set('is_admin', v === true)
                                        }
                                    />
                                    Superadmin
                                </label>
                            </div>
                            <div className="flex justify-end gap-2 border-t border-sidebar-border pt-4">
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={() => setCreateOpen(false)}
                                >
                                    Batal
                                </Button>
                                <Button
                                    type="submit"
                                    className="bg-orange-600 text-white hover:bg-orange-700"
                                >
                                    Simpan
                                </Button>
                            </div>
                        </form>
                    </DialogContent>
                </Dialog>

                <form
                    className="flex flex-wrap gap-2"
                    onSubmit={(e) => {
                        e.preventDefault();
                        search(1);
                    }}
                >
                    <div className="relative min-w-52 flex-1">
                        <Search className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            value={q}
                            onChange={(e) => setQ(e.target.value)}
                            placeholder="Cari nama atau email…"
                            className="pl-9"
                        />
                    </div>
                    <select
                        value={plan}
                        onChange={(e) => {
                            setPlan(e.target.value);
                            router.get(
                                '/spr/users',
                                {
                                    q: q || undefined,
                                    plan: e.target.value || undefined,
                                    page: 1,
                                },
                                { preserveState: true, replace: true },
                            );
                        }}
                        className="h-9 border border-sidebar-border bg-background px-2 text-sm"
                        aria-label="Filter paket"
                    >
                        <option value="">Semua paket</option>
                        {plans.map((p) => (
                            <option key={p} value={p}>
                                {p}
                            </option>
                        ))}
                    </select>
                    <Button type="submit" variant="outline">
                        Cari
                    </Button>
                </form>

                <div className="overflow-x-auto border border-sidebar-border bg-card">
                    <table className="w-full min-w-180 text-left text-sm">
                        <thead>
                            <tr className="border-b border-sidebar-border bg-muted/40 font-mono text-xs tracking-wide uppercase">
                                <th className="px-4 py-3 font-semibold">
                                    Pengguna
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Paket
                                </th>
                                <th className="px-4 py-3 text-center font-semibold">
                                    Admin
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Pesan
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Sesi
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Aksi
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {users.data.map((u) => (
                                <UserEditor
                                    key={u.id}
                                    user={u}
                                    plans={plans}
                                    isSelf={u.id === auth.user.id}
                                />
                            ))}
                            {users.data.length === 0 && (
                                <tr>
                                    <td
                                        colSpan={6}
                                        className="px-4 py-6 text-center text-muted-foreground"
                                    >
                                        Tidak ada pengguna yang cocok.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        Halaman {users.current_page} dari {users.last_page} ·
                        Total {users.total}
                    </span>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={users.current_page <= 1}
                            onClick={() => gotoPage(users.current_page - 1)}
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={
                                users.current_page >= users.last_page
                            }
                            onClick={() => gotoPage(users.current_page + 1)}
                        >
                            Berikutnya
                        </Button>
                    </div>
                </div>

                <p className="text-xs text-muted-foreground">
                    Mengubah paket ke gratis memberi masa percobaan 14 hari;
                    paket berbayar tanpa batas waktu. Menampilkan{' '}
                    <Badge variant="outline">Admin</Badge> hanya penanda —
                    status admin diatur lewat centang di tabel.
                </p>
            </div>
        </>
    );
}

SprUsers.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
        {
            title: 'Pengguna',
            href: '/spr/users',
        },
    ],
};

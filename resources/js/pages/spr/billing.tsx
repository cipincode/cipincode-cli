import { Head, router } from '@inertiajs/react';
import { useState } from 'react';
import { Search } from 'lucide-react';
import StatCard from '@/components/api-keys/stat-card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { CreditCard, Hourglass, TimerOff, Users } from 'lucide-react';

type Customer = {
    id: number;
    name: string;
    email: string;
    plan: string;
    plan_label: string;
    price_label: string;
    messages_total: number;
    created_at: string;
};

type Paginated = {
    data: Customer[];
    current_page: number;
    last_page: number;
    total: number;
};

type Props = {
    stats: {
        revenue: string;
        customers: number;
        trial: number;
        expired: number;
    };
    customers_list: Paginated;
    filters: { q: string };
    note: string;
};

export default function SprBilling({ stats, customers_list, filters, note }: Props) {
    const [q, setQ] = useState(filters.q);

    const search = (page = 1) => {
        router.get(
            '/spr/billing',
            { q: q || undefined, page },
            { preserveState: true, replace: true },
        );
    };

    return (
        <>
            <Head title="Billing" />
            <div className="flex flex-1 flex-col gap-6 p-6">
                <div>
                    <h1 className="text-xl font-semibold tracking-tight">
                        Billing
                    </h1>
                    <p className="text-sm text-muted-foreground">{note}</p>
                </div>

                <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                    <StatCard
                        label="Estimasi / bulan"
                        value={stats.revenue}
                        caption="Dari paket berbayar aktif"
                        icon={CreditCard}
                        tone="brand"
                    />
                    <StatCard
                        label="Pelanggan berbayar"
                        value={stats.customers}
                        caption="Paket Go, Pro, Max"
                        icon={Users}
                    />
                    <StatCard
                        label="Trial aktif"
                        value={stats.trial}
                        caption="Paket gratis dalam masa coba"
                        icon={Hourglass}
                    />
                    <StatCard
                        label="Kedaluwarsa"
                        value={stats.expired}
                        caption="Masa aktif habis"
                        icon={TimerOff}
                        tone={stats.expired ? 'muted' : 'success'}
                    />
                </div>

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
                            placeholder="Cari pelanggan…"
                            className="pl-9"
                        />
                    </div>
                    <Button type="submit" variant="outline">
                        Cari
                    </Button>
                </form>

                <div className="overflow-x-auto border border-sidebar-border bg-card">
                    <table className="w-full min-w-180 text-left text-sm">
                        <thead>
                            <tr className="border-b border-sidebar-border bg-muted/40 font-mono text-xs tracking-wide uppercase">
                                <th className="px-4 py-3 font-semibold">
                                    Pelanggan
                                </th>
                                <th className="px-4 py-3 font-semibold">
                                    Paket
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Tagihan
                                </th>
                                <th className="px-4 py-3 text-right font-semibold">
                                    Pesan
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {customers_list.data.map((c) => (
                                <tr
                                    key={c.id}
                                    className="border-t border-sidebar-border"
                                >
                                    <td className="px-4 py-3">
                                        <p className="font-medium">{c.name}</p>
                                        <p className="text-xs text-muted-foreground">
                                            {c.email}
                                        </p>
                                    </td>
                                    <td className="px-4 py-3 capitalize">
                                        {c.plan_label}
                                    </td>
                                    <td className="px-4 py-3 text-right font-medium tabular-nums">
                                        {c.price_label}
                                        <span className="ml-1 text-xs font-normal text-muted-foreground">
                                            /bln
                                        </span>
                                    </td>
                                    <td className="px-4 py-3 text-right tabular-nums">
                                        {c.messages_total}
                                    </td>
                                </tr>
                            ))}
                            {customers_list.data.length === 0 && (
                                <tr>
                                    <td
                                        colSpan={4}
                                        className="px-4 py-6 text-center text-muted-foreground"
                                    >
                                        Belum ada pelanggan berbayar.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground tabular-nums">
                        Halaman {customers_list.current_page} dari{' '}
                        {customers_list.last_page} · Total{' '}
                        {customers_list.total}
                    </span>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={customers_list.current_page <= 1}
                            onClick={() => search(customers_list.current_page - 1)}
                        >
                            Sebelumnya
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={
                                customers_list.current_page >=
                                customers_list.last_page
                            }
                            onClick={() => search(customers_list.current_page + 1)}
                        >
                            Berikutnya
                        </Button>
                    </div>
                </div>
            </div>
        </>
    );
}

SprBilling.layout = {
    breadcrumbs: [
        {
            title: 'Superadmin',
            href: '/spr/dashboard',
        },
        {
            title: 'Billing',
            href: '/spr/billing',
        },
    ],
};

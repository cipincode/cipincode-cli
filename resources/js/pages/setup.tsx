import { Head, Link } from '@inertiajs/react';
import { ExternalLink } from 'lucide-react';
import CopyCommand from '@/components/setup/copy-command';
import { apiKeys, setup } from '@/routes';

type Step = {
    n: number;
    title: string;
    description: string;
    command: string;
    links?: { label: string; href: string }[];
};

const steps: Step[] = [
    {
        n: 1,
        title: 'Pasang CipinCode CLI',
        description:
            'Satu perintah ini otomatis memasang Go, Bun, dan binary CipinCode CLI sesuai sistem operasi Anda (Linux & macOS).',
        command: 'npm install -g cipincode',
    },
    {
        n: 2,
        title: 'Verifikasi instalasi',
        description:
            'Pastikan CLI terpasang dengan benar. Jika nomor versi tampil, berarti siap dipakai.',
        command: 'cipincode --version',
    },
    {
        n: 3,
        title: 'Masuk dan mulai',
        description:
            'Jalankan cipincode, lalu pilih masuk dengan Email atau Google. Buka URL yang tampil, klik Izinkan, dan tempel kode yang muncul ke terminal. Ketik / di dalam CLI untuk melihat semua perintah.',
        command: 'cipincode login',
    },
];

export default function Setup() {
    return (
        <>
            <Head title="Setup Installation" />

            <div className="flex flex-1 flex-col">
                <div className="border-b border-sidebar-border px-6 py-5">
                    <h1 className="text-xl font-semibold tracking-tight">
                        Setup Installation
                    </h1>
                    <p className="mt-1 text-sm text-muted-foreground">
                        Tiga langkah untuk memasang CipinCode CLI dan mulai
                        memakai agent dari terminal Anda.
                    </p>
                </div>

                <div className="flex-1">
                    <div className="p-6">
                        <ol className="mx-auto flex max-w-3xl flex-col gap-4">
                            {steps.map((step) => {
                                return (
                                    <li
                                        key={step.n}
                                        className="border border-sidebar-border bg-card"
                                    >
                                        <div className="flex items-stretch border-b border-sidebar-border">
                                            <div className="flex size-12 shrink-0 items-center justify-center border-r border-sidebar-border text-sm font-semibold text-orange-600">
                                                {step.n}
                                            </div>
                                            <h2 className="flex items-center px-4 font-medium">
                                                {step.title}
                                            </h2>
                                        </div>

                                        <div className="space-y-4 p-5">
                                            <p className="text-sm text-muted-foreground">
                                                {step.description}
                                            </p>

                                            <CopyCommand
                                                command={step.command}
                                            />
                                        </div>

                                        {(step.links || step.n === 4) && (
                                            <div className="flex flex-wrap gap-2 border-t border-sidebar-border px-5 py-3">
                                                {step.links?.map((link) => (
                                                    <a
                                                        key={link.href}
                                                        href={link.href}
                                                        target="_blank"
                                                        rel="noreferrer"
                                                        className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-xs font-medium transition-colors hover:bg-sidebar-accent"
                                                    >
                                                        <ExternalLink className="size-3.5" />
                                                        {link.label}
                                                    </a>
                                                ))}

                                                {step.n === 4 && (
                                                    <Link
                                                        href={apiKeys()}
                                                        className="inline-flex h-8 items-center gap-1.5 border border-sidebar-border px-3 text-xs font-medium transition-colors hover:bg-sidebar-accent"
                                                    >
                                                        Buka API Keys
                                                    </Link>
                                                )}
                                            </div>
                                        )}
                                    </li>
                                );
                            })}
                        </ol>
                    </div>
                </div>
            </div>
        </>
    );
}

Setup.layout = {
    breadcrumbs: [
        {
            title: 'Setup Installation',
            href: setup(),
        },
    ],
};

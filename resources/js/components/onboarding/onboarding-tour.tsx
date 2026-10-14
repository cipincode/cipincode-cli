import { useEffect, useState } from 'react';
import { Link } from '@inertiajs/react';
import {
    ArrowLeft,
    ArrowRight,
    Bot,
    CreditCard,
    Download,
    KeyRound,
    PartyPopper,
    Sparkles,
    TerminalSquare,
    TrendingUp,
    X,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
    agents,
    apiKeys,
    billing,
    sessions,
    setup,
    skills,
    usage,
} from '@/routes';

type Step = {
    icon: typeof Download;
    title: string;
    description: string;
    /** Selector CSS menu sidebar yang disorot. Kosong = kartu tengah. */
    target?: string;
    href?: () => { url: string };
    linkLabel?: string;
};

const side = (path: string) => `[data-slot="sidebar"] a[href="${path}"]`;

const steps: Step[] = [
    {
        icon: PartyPopper,
        title: 'Selamat datang di CipinCode!',
        description:
            'Ikuti tur singkat ini — kami akan menunjuk langsung tiap menu di sidebar dan menjelaskan urutan memakainya.',
    },
    {
        icon: Download,
        title: '1. Pasang CLI dulu',
        description:
            'Mulai dari menu Setup Installation: salin satu perintah instalasi ke terminal, lalu verifikasi dengan cipincode --version.',
        target: side(setup().url),
        href: () => setup(),
        linkLabel: 'Buka Setup Installation',
    },
    {
        icon: KeyRound,
        title: '2. Hubungkan dengan API Key',
        description:
            'Buat API key di menu ini, lalu jalankan cipincode login di terminal dan tempel kuncinya. Satu key untuk satu perangkat.',
        target: side(apiKeys().url),
        href: () => apiKeys(),
        linkLabel: 'Buka API Keys',
    },
    {
        icon: Bot,
        title: '3. Bikin agent kustom',
        description:
            'Di menu Agents Anda bisa membuat agent sendiri — atur mode (chat, coding, debugging) dan effort-nya.',
        target: side(agents().url),
        href: () => agents(),
        linkLabel: 'Buka Agents',
    },
    {
        icon: Sparkles,
        title: '4. Tambah Skills',
        description:
            'Menu Skills berisi kemampuan tambahan yang bisa dipasang ke agent supaya makin pintar.',
        target: side(skills().url),
        href: () => skills(),
        linkLabel: 'Buka Skills',
    },
    {
        icon: TerminalSquare,
        title: '5. Lihat semua sesi',
        description:
            'Menu Sessions menyimpan seluruh percakapan CLI Anda — bisa dibuka ulang atau dihapus kapan saja.',
        target: side(sessions().url),
        href: () => sessions(),
        linkLabel: 'Buka Sessions',
    },
    {
        icon: TrendingUp,
        title: '6. Pantau kuota',
        description:
            'Menu Usage menunjukkan sisa jatah pesan per 5 jam dan per minggu supaya tidak kehabisan tiba-tiba.',
        target: side(usage().url),
        href: () => usage(),
        linkLabel: 'Buka Usage',
    },
    {
        icon: CreditCard,
        title: '7. Paket & harga',
        description:
            'Paket Gratis cukup untuk mencoba. Butuh kuota lebih? Menu Billing berisi paket Go, Pro, dan Max dalam rupiah.',
        target: side(billing().url),
        href: () => billing(),
        linkLabel: 'Buka Billing',
    },
    {
        icon: PartyPopper,
        title: 'Siap dipakai!',
        description:
            'Itu dia semua menunya. Selamat ngoding dengan agen AI — centang di bawah kalau tur ini sudah tidak perlu muncul lagi.',
    },
];

export function storageKey(userId: number | string) {
    return `cipincode:onboarding:v1:${userId}`;
}

export function isOnboardingDone(userId: number | string) {
    try {
        return localStorage.getItem(storageKey(userId)) === '1';
    } catch {
        return true;
    }
}

function isVisible(el: HTMLElement) {
    const r = el.getBoundingClientRect();
    if (r.width === 0 || r.height === 0) return false;
    return (
        r.bottom > 0 &&
        r.top < window.innerHeight &&
        r.right > 0 &&
        r.left < window.innerWidth
    );
}

export default function OnboardingTour({
    open,
    onOpenChange,
    userId,
}: {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    userId: number | string;
}) {
    const [step, setStep] = useState(0);
    const [dontShow, setDontShow] = useState(false);
    const [rect, setRect] = useState<DOMRect | null>(null);
    const current = steps[step];
    const last = step === steps.length - 1;

    useEffect(() => {
        if (open) {
            setStep(0);
            try {
                setDontShow(
                    localStorage.getItem(storageKey(userId)) === '1',
                );
            } catch {
                setDontShow(false);
            }
        }
    }, [open, userId]);

    useEffect(() => {
        if (!open) return;
        const selector = current.target;
        if (!selector || window.innerWidth < 1024) {
            setRect(null);
            return;
        }
        let raf = 0;
        let timer: ReturnType<typeof setTimeout>;
        const measure = () => {
            const el = document.querySelector(selector) as HTMLElement | null;
            if (!el || !isVisible(el)) {
                setRect(null);
                return;
            }
            setRect(el.getBoundingClientRect());
        };
        const el = document.querySelector(selector) as HTMLElement | null;
        el?.scrollIntoView({ block: 'nearest' });
        timer = setTimeout(() => {
            raf = requestAnimationFrame(measure);
        }, 350);
        window.addEventListener('resize', measure);
        window.addEventListener('scroll', measure, true);
        return () => {
            clearTimeout(timer);
            cancelAnimationFrame(raf);
            window.removeEventListener('resize', measure);
            window.removeEventListener('scroll', measure, true);
        };
    }, [open, step, current.target]);

    useEffect(() => {
        if (!open) return;
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape') close(dontShow);
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    });

    const close = (persist: boolean) => {
        if (persist) {
            try {
                localStorage.setItem(storageKey(userId), '1');
            } catch {
                /* abaikan */
            }
        }
        onOpenChange(false);
    };

    if (!open) return null;

    const Icon = current.icon;
    const tipWidth = 340;
    const tipStyle: React.CSSProperties | undefined = rect
        ? {
              left: Math.max(
                  12,
                  Math.min(
                      rect.right + 14,
                      window.innerWidth - tipWidth - 12,
                  ),
              ),
              top: Math.max(12, rect.top + rect.height / 2 - 140),
          }
        : undefined;

    return (
        <div className="fixed inset-0 z-[100]" role="dialog" aria-modal="true" aria-label="Tur pengenalan dashboard">
            <div
                className="absolute inset-0 bg-black/60"
                onClick={() => close(dontShow)}
            />

            {rect && (
                <div
                    aria-hidden
                    className="pointer-events-none fixed rounded-xl border-2 border-orange-500 transition-all duration-300"
                    style={{
                        top: rect.top - 6,
                        left: rect.left - 6,
                        width: rect.width + 12,
                        height: rect.height + 12,
                        boxShadow: '0 0 0 9999px rgba(0, 0, 0, 0.6)',
                    }}
                />
            )}

            <div
                className={`fixed z-[101] w-[min(340px,calc(100vw-24px))] rounded-2xl border border-sidebar-border bg-card p-5 shadow-2xl ${
                    rect ? '' : 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2'
                }`}
                style={tipStyle}
            >
                <button
                    type="button"
                    aria-label="Lewati tur"
                    onClick={() => close(dontShow)}
                    className="absolute top-3 right-3 rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                    <X className="size-4" />
                </button>

                <div className="flex items-center gap-3 pr-8">
                    <span className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-orange-600/10 text-orange-600">
                        <Icon className="size-5" />
                    </span>
                    <div>
                        <p className="font-semibold">{current.title}</p>
                        <p className="text-xs text-muted-foreground">
                            Langkah {step + 1} dari {steps.length}
                        </p>
                    </div>
                </div>

                <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                    {current.description}
                </p>

                {current.href && current.linkLabel && (
                    <Link
                        href={current.href().url}
                        onClick={() => close(dontShow)}
                        className="mt-3 inline-flex h-9 items-center gap-1.5 self-start border border-sidebar-border px-3 text-xs font-medium transition-colors hover:bg-sidebar-accent"
                    >
                        {current.linkLabel}
                    </Link>
                )}

                <div className="mt-4 flex items-center gap-1.5">
                    {steps.map((s, i) => (
                        <button
                            key={s.title}
                            type="button"
                            aria-label={`Ke langkah ${i + 1}`}
                            onClick={() => setStep(i)}
                            className={`h-1.5 flex-1 rounded-full transition-colors ${
                                i === step
                                    ? 'bg-orange-600'
                                    : i < step
                                      ? 'bg-orange-600/40'
                                      : 'bg-muted'
                            }`}
                        />
                    ))}
                </div>

                <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-sidebar-border pt-4">
                    <label
                        htmlFor="onboarding-dont-show"
                        className="flex cursor-pointer items-center gap-2 text-sm text-muted-foreground select-none"
                    >
                        <Checkbox
                            id="onboarding-dont-show"
                            checked={dontShow}
                            onCheckedChange={(v) => setDontShow(v === true)}
                        />
                        Jangan tampilkan lagi
                    </label>

                    <div className="flex gap-2">
                        {step > 0 && (
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => setStep((s) => s - 1)}
                            >
                                <ArrowLeft className="size-4" />
                                Kembali
                            </Button>
                        )}
                        {!last ? (
                            <Button
                                size="sm"
                                onClick={() => setStep((s) => s + 1)}
                                className="bg-orange-600 text-white hover:bg-orange-700"
                            >
                                Lanjut
                                <ArrowRight className="size-4" />
                            </Button>
                        ) : (
                            <Button
                                size="sm"
                                onClick={() => close(true)}
                                className="bg-orange-600 text-white hover:bg-orange-700"
                            >
                                Mulai Pakai
                            </Button>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}

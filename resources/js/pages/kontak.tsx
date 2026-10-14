import { useState } from 'react';
import { Head, router, usePage } from '@inertiajs/react';
import { CheckCircle2, Send } from 'lucide-react';
import Header from '../../../components/saas-landing/header';
import Footer from '../../../components/saas-landing/footer';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

export default function Kontak() {
    const { errors } = usePage().props;
    const [form, setForm] = useState({ name: '', email: '', message: '' });
    const [sent, setSent] = useState(false);

    const set = (key: keyof typeof form, value: string) =>
        setForm((f) => ({ ...f, [key]: value }));

    const submit = (e: React.FormEvent) => {
        e.preventDefault();
        router.post('/kontak', form, {
            onSuccess: () => {
                setSent(true);
                setForm({ name: '', email: '', message: '' });
            },
        });
    };

    return (
        <>
            <Head title="Kontak" />
            <div className="relative isolate flex min-h-screen w-full flex-col bg-background text-foreground">
                <Header />
                <main className="mx-auto w-full max-w-2xl flex-1 px-4 py-12 sm:px-6 md:py-16">
                    <p className="font-mono text-xs font-semibold tracking-widest text-orange-600 uppercase">
                        Kontak
                    </p>
                    <h1 className="mt-2 text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                        Ada kendala? Ceritakan ke kami
                    </h1>
                    <p className="mt-3 leading-relaxed text-muted-foreground">
                        Tulis komplain, pertanyaan, atau masukan soal
                        CipinCode CLI. Pesan Anda masuk ke kotak masuk tim
                        kami dan dibalas lewat email.
                    </p>

                    {sent ? (
                        <div className="mt-8 rounded-2xl border border-green-600/30 bg-green-600/10 p-6 text-center">
                            <CheckCircle2 className="mx-auto size-10 text-green-600" />
                            <p className="mt-3 text-lg font-semibold">
                                Pesan terkirim!
                            </p>
                            <p className="mt-1 text-sm text-muted-foreground">
                                Terima kasih sudah menghubungi kami. Tunggu
                                balasan di email Anda.
                            </p>
                            <Button
                                variant="outline"
                                className="mt-5"
                                onClick={() => setSent(false)}
                            >
                                Kirim pesan lain
                            </Button>
                        </div>
                    ) : (
                        <form
                            onSubmit={submit}
                            className="mt-8 rounded-2xl border border-border bg-card p-6 sm:p-8"
                        >
                            <div className="flex flex-col gap-5">
                                <div className="flex flex-col gap-1.5">
                                    <Label htmlFor="kontak-name">Nama</Label>
                                    <Input
                                        id="kontak-name"
                                        value={form.name}
                                        onChange={(e) =>
                                            set('name', e.target.value)
                                        }
                                        placeholder="Nama Anda"
                                    />
                                    {errors.name && (
                                        <p className="text-xs text-red-600">
                                            {errors.name}
                                        </p>
                                    )}
                                </div>
                                <div className="flex flex-col gap-1.5">
                                    <Label htmlFor="kontak-email">Email</Label>
                                    <Input
                                        id="kontak-email"
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
                                    <Label htmlFor="kontak-message">
                                        Pesan / Komplain
                                    </Label>
                                    <Textarea
                                        id="kontak-message"
                                        value={form.message}
                                        onChange={(e) =>
                                            set('message', e.target.value)
                                        }
                                        placeholder="Contoh: instalasi gagal di Ubuntu 24.04 dengan pesan…"
                                        rows={6}
                                    />
                                    {errors.message && (
                                        <p className="text-xs text-red-600">
                                            {errors.message}
                                        </p>
                                    )}
                                </div>
                                <Button
                                    type="submit"
                                    size="lg"
                                    className="w-full rounded-full bg-orange-600 text-white hover:bg-orange-700 hover:text-white"
                                >
                                    <Send className="size-4" />
                                    Kirim Pesan
                                </Button>
                            </div>
                        </form>
                    )}
                </main>
                <div className="border-t border-border/50">
                    <Footer />
                </div>
            </div>
        </>
    );
}

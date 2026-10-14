import type { AuthLayoutProps } from '@/types';

export default function AuthSimpleLayout({
    children,
    title,
    description,
}: AuthLayoutProps) {
    return (
        <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
            <div className="w-full max-w-sm">
                <div className="flex flex-col gap-6">
                    <div className="flex h-20 w-[340px] items-center justify-center self-center">
                        <img
                            src="/images/logo-login-light.png"
                            alt="CipinCode"
                            className="h-full w-full object-contain dark:hidden"
                        />
                        <img
                            src="/images/logo-login-dark.png"
                            alt="CipinCode"
                            className="hidden h-full w-full object-contain dark:block"
                        />
                    </div>

                    <div className="flex flex-col gap-8 rounded-xl border bg-card p-8 text-card-foreground shadow-sm">
                        <div className="space-y-2 text-center">
                            <h1 className="text-xl font-medium">{title}</h1>
                            <p className="text-center text-sm text-muted-foreground">
                                {description}
                            </p>
                        </div>
                        {children}
                    </div>
                </div>
            </div>
        </div>
    );
}

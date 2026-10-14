import { Link } from '@inertiajs/react';
import { Palette, ShieldCheck, UserRound } from 'lucide-react';
import type { PropsWithChildren } from 'react';
import { useCurrentUrl } from '@/hooks/use-current-url';
import { cn, toUrl } from '@/lib/utils';
import { edit as editAppearance } from '@/routes/appearance';
import { edit } from '@/routes/profile';
import { edit as editSecurity } from '@/routes/security';
import type { NavItem } from '@/types';

const sidebarNavItems: NavItem[] = [
    { title: 'Profile', href: edit(), icon: UserRound },
    { title: 'Security', href: editSecurity(), icon: ShieldCheck },
    { title: 'Appearance', href: editAppearance(), icon: Palette },
];

export default function SettingsLayout({ children }: PropsWithChildren) {
    const { isCurrentOrParentUrl } = useCurrentUrl();

    return (
        <div className="flex flex-1 flex-col">
            <div className="border-b border-sidebar-border px-6 py-6">
                <h1 className="text-xl font-semibold tracking-tight">
                    Settings
                </h1>
                <p className="text-sm text-muted-foreground">
                    Kelola profil dan pengaturan akun Anda.
                </p>
            </div>

            <div className="flex flex-1 flex-col lg:flex-row">
                <aside className="border-b border-sidebar-border lg:w-60 lg:border-r lg:border-b-0">
                    <nav
                        className="flex gap-1 overflow-x-auto p-3 lg:flex-col"
                        aria-label="Settings"
                    >
                        {sidebarNavItems.map((item, index) => {
                            const active = isCurrentOrParentUrl(item.href);

                            return (
                                <Link
                                    key={`${toUrl(item.href)}-${index}`}
                                    href={item.href}
                                    className={cn(
                                        'flex items-center gap-2 border-l-2 px-3 py-2 text-sm whitespace-nowrap transition-colors',
                                        active
                                            ? 'border-orange-600 bg-orange-600/10 font-medium text-orange-600'
                                            : 'border-transparent text-muted-foreground hover:bg-sidebar-accent hover:text-foreground',
                                    )}
                                >
                                    {item.icon && (
                                        <item.icon className="size-4" />
                                    )}
                                    {item.title}
                                </Link>
                            );
                        })}
                    </nav>
                </aside>

                <div className="min-w-0 flex-1 p-6">
                    <section className="max-w-3xl space-y-6">
                        {children}
                    </section>
                </div>
            </div>
        </div>
    );
}

import { Link } from '@inertiajs/react';
import {
    Bot,
    CreditCard,
    Download,
    KeyRound,
    LayoutGrid,
    Sparkles,
    TerminalSquare,
    TrendingUp,
} from 'lucide-react';
import AppLogo from '@/components/app-logo';
import { NavMain } from '@/components/nav-main';
import { NavUser } from '@/components/nav-user';
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
} from '@/components/ui/sidebar';
import {
    agents,
    apiKeys,
    billing,
    dashboard,
    sessions,
    setup,
    skills,
    usage,
} from '@/routes';
import type { NavItem } from '@/types';

const platformNavItems: NavItem[] = [
    {
        title: 'Dashboard',
        href: dashboard(),
        icon: LayoutGrid,
    },
    {
        title: 'Setup Installation',
        href: setup(),
        icon: Download,
    },
    {
        title: 'Agents',
        href: agents(),
        icon: Bot,
    },
    {
        title: 'Skills',
        href: skills(),
        icon: Sparkles,
    },
    {
        title: 'Sessions',
        href: sessions(),
        icon: TerminalSquare,
    },
];

const accessNavItems: NavItem[] = [
    {
        title: 'API Keys',
        href: apiKeys(),
        icon: KeyRound,
    },
    {
        title: 'Usage',
        href: usage(),
        icon: TrendingUp,
    },
];

const accountNavItems: NavItem[] = [
    {
        title: 'Billing',
        href: billing(),
        icon: CreditCard,
    },
];

export function AppSidebar() {
    return (
        <Sidebar collapsible="icon" variant="sidebar">
            <SidebarHeader className="h-16 justify-center border-b border-sidebar-border group-data-[collapsible=icon]:h-12">
                <Link
                    href={dashboard()}
                    prefetch
                    className="flex w-full items-center justify-center"
                >
                    <AppLogo />
                </Link>
            </SidebarHeader>

            <SidebarContent>
                <NavMain items={platformNavItems} label="Platform" />
                <NavMain items={accessNavItems} label="Akses CLI" />
                <NavMain items={accountNavItems} label="Akun" />
            </SidebarContent>

            <SidebarFooter className="border-t border-sidebar-border">
                <NavUser />
            </SidebarFooter>
        </Sidebar>
    );
}

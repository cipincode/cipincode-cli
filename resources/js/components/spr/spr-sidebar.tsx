import { Link } from '@inertiajs/react';
import {
    ArrowLeft,
    BadgeCheck,
    CreditCard,
    Inbox,
    LayoutGrid,
    TerminalSquare,
    Users,
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
import type { NavItem } from '@/types';

const sprNavItems: NavItem[] = [
    {
        title: 'Dashboard',
        href: '/spr/dashboard',
        icon: LayoutGrid,
    },
    {
        title: 'Pengguna',
        href: '/spr/users',
        icon: Users,
    },
    {
        title: 'Sesi',
        href: '/spr/sessions',
        icon: TerminalSquare,
    },
];

const backNavItems: NavItem[] = [
    {
        title: 'Kembali ke Dashboard',
        href: '/dashboard',
        icon: ArrowLeft,
    },
];

const manageNavItems: NavItem[] = [
    {
        title: 'Billing',
        href: '/spr/billing',
        icon: CreditCard,
    },
    {
        title: 'Member',
        href: '/spr/members',
        icon: BadgeCheck,
    },
    {
        title: 'Contact',
        href: '/spr/contacts',
        icon: Inbox,
    },
];

export function SprSidebar() {
    return (
        <Sidebar collapsible="icon" variant="sidebar">
            <SidebarHeader className="h-16 justify-center border-b border-sidebar-border group-data-[collapsible=icon]:h-12">
                <Link
                    href="/spr/dashboard"
                    prefetch
                    className="flex w-full items-center justify-center"
                >
                    <AppLogo />
                </Link>
            </SidebarHeader>

            <SidebarContent>
                <NavMain items={sprNavItems} label="Superadmin" />
                <NavMain items={manageNavItems} label="Kelola" />
                <NavMain items={backNavItems} label="Navigasi" />
            </SidebarContent>

            <SidebarFooter className="border-t border-sidebar-border">
                <NavUser />
            </SidebarFooter>
        </Sidebar>
    );
}

import { AppContent } from '@/components/app-content';
import { AppShell } from '@/components/app-shell';
import { AppSidebarHeader } from '@/components/app-sidebar-header';
import { SprSidebar } from '@/components/spr/spr-sidebar';
import { useFlashToast } from '@/hooks/use-flash-toast';
import type { AppLayoutProps } from '@/types';

export default function SprLayout({
    children,
    breadcrumbs = [],
}: AppLayoutProps) {
    useFlashToast();

    return (
        <AppShell variant="sidebar">
            <SprSidebar />
            <AppContent variant="sidebar" className="min-w-0 overflow-x-clip">
                <AppSidebarHeader breadcrumbs={breadcrumbs} />
                {children}
            </AppContent>
        </AppShell>
    );
}

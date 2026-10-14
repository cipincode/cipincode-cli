import type { LucideIcon } from 'lucide-react';
import type { ReactNode } from 'react';

type Props = {
    title: string;
    description: string;
    icon: LucideIcon;
    emptyTitle: string;
    emptyDescription: string;
    children?: ReactNode;
};

export default function PageSection({
    title,
    description,
    icon: Icon,
    emptyTitle,
    emptyDescription,
    children,
}: Props) {
    return (
        <div className="flex flex-1 flex-col">
            <div className="border-b border-sidebar-border px-6 py-5">
                <h1 className="text-xl font-semibold tracking-tight">
                    {title}
                </h1>
                <p className="mt-1 text-sm text-muted-foreground">
                    {description}
                </p>
            </div>

            <div className="flex flex-1 items-center justify-center p-6">
                {children ?? (
                    <div className="flex max-w-sm flex-col items-center gap-3 text-center">
                        <div className="flex size-12 items-center justify-center border border-sidebar-border text-orange-600">
                            <Icon className="size-6" />
                        </div>
                        <h2 className="text-base font-medium">{emptyTitle}</h2>
                        <p className="text-sm text-muted-foreground">
                            {emptyDescription}
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}

import {
    KeyRound,
    MoreHorizontal,
    Pencil,
    Power,
    PowerOff,
    Trash2,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { formatDate, formatDateTime, formatRelative } from './format';
import type { ApiKey } from './types';

type Props = {
    keys: ApiKey[];
    onEdit: (key: ApiKey) => void;
    onToggle: (key: ApiKey) => void;
    onDelete: (key: ApiKey) => void;
};

const th =
    'px-5 py-3 text-left text-xs font-medium tracking-wide text-muted-foreground uppercase whitespace-nowrap';

export default function KeysTable({ keys, onEdit, onToggle, onDelete }: Props) {
    return (
        <div className="overflow-x-auto">
            <table className="w-full min-w-[720px] text-sm">
                <thead className="border-b border-sidebar-border bg-muted/40">
                    <tr>
                        <th className={th}>Nama</th>
                        <th className={th}>Key</th>
                        <th className={th}>Status</th>
                        <th className={th}>Dibuat</th>
                        <th className={th}>Terakhir dipakai</th>
                        <th className={`${th} w-12`}>
                            <span className="sr-only">Aksi</span>
                        </th>
                    </tr>
                </thead>
                <tbody className="divide-y divide-sidebar-border">
                    {keys.map((key) => (
                        <tr
                            key={key.id}
                            className="transition-colors hover:bg-muted/30"
                        >
                            <td className="px-5 py-3.5">
                                <div className="flex items-center gap-3">
                                    <div className="flex size-8 shrink-0 items-center justify-center border border-sidebar-border text-muted-foreground">
                                        <KeyRound className="size-4" />
                                    </div>
                                    <span className="font-medium">
                                        {key.name}
                                    </span>
                                </div>
                            </td>
                            <td className="px-5 py-3.5">
                                <code className="font-mono text-xs text-muted-foreground">
                                    {key.prefix}••••••••
                                </code>
                            </td>
                            <td className="px-5 py-3.5">
                                {key.is_active ? (
                                    <Badge
                                        variant="outline"
                                        className="rounded-none border-emerald-600/40 bg-emerald-600/10 text-emerald-700 dark:text-emerald-400"
                                    >
                                        <span className="size-1.5 bg-emerald-600" />
                                        Aktif
                                    </Badge>
                                ) : (
                                    <Badge
                                        variant="outline"
                                        className="rounded-none text-muted-foreground"
                                    >
                                        <span className="size-1.5 bg-muted-foreground/60" />
                                        Nonaktif
                                    </Badge>
                                )}
                            </td>
                            <td className="px-5 py-3.5 whitespace-nowrap text-muted-foreground">
                                {formatDate(key.created_at)}
                            </td>
                            <td className="px-5 py-3.5 whitespace-nowrap text-muted-foreground">
                                {key.last_used_at ? (
                                    <span
                                        title={formatDateTime(key.last_used_at)}
                                    >
                                        {formatRelative(key.last_used_at)}
                                    </span>
                                ) : (
                                    'Belum pernah'
                                )}
                            </td>
                            <td className="px-5 py-3.5 text-right">
                                <DropdownMenu>
                                    <DropdownMenuTrigger asChild>
                                        <button
                                            type="button"
                                            aria-label={`Aksi untuk ${key.name}`}
                                            className="inline-flex size-8 items-center justify-center border border-transparent text-muted-foreground transition-colors hover:border-sidebar-border hover:text-foreground data-[state=open]:border-sidebar-border"
                                        >
                                            <MoreHorizontal className="size-4" />
                                        </button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent
                                        align="end"
                                        className="w-48 rounded-none"
                                    >
                                        <DropdownMenuItem
                                            className="rounded-none"
                                            onSelect={() => onEdit(key)}
                                        >
                                            <Pencil className="size-4" />
                                            Ubah nama
                                        </DropdownMenuItem>
                                        <DropdownMenuItem
                                            className="rounded-none"
                                            onSelect={() => onToggle(key)}
                                        >
                                            {key.is_active ? (
                                                <>
                                                    <PowerOff className="size-4" />
                                                    Nonaktifkan
                                                </>
                                            ) : (
                                                <>
                                                    <Power className="size-4" />
                                                    Aktifkan
                                                </>
                                            )}
                                        </DropdownMenuItem>
                                        <DropdownMenuSeparator />
                                        <DropdownMenuItem
                                            variant="destructive"
                                            className="rounded-none"
                                            onSelect={() => onDelete(key)}
                                        >
                                            <Trash2 className="size-4" />
                                            Hapus
                                        </DropdownMenuItem>
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}

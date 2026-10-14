import { Check, Copy } from 'lucide-react';
import { useState } from 'react';

export default function CopyCommand({
    command,
    prompt = true,
    wrap = false,
}: {
    command: string;
    prompt?: boolean;
    /** Di layar kecil teks dibungkus ke baris baru, di layar lebar tetap satu baris. */
    wrap?: boolean;
}) {
    const [copied, setCopied] = useState(false);

    const copy = async () => {
        try {
            await navigator.clipboard.writeText(command);
            setCopied(true);
            setTimeout(() => setCopied(false), 1500);
        } catch {
            setCopied(false);
        }
    };

    return (
        <div className="flex items-stretch border border-zinc-800 bg-zinc-950 font-mono text-[13px] text-zinc-100">
            <code
                className={`min-w-0 flex-1 px-3 py-2 ${
                    wrap
                        ? 'break-all whitespace-normal sm:break-normal sm:whitespace-nowrap'
                        : 'overflow-x-auto whitespace-nowrap'
                }`}
            >
                {prompt && <span className="mr-2 text-orange-500">$</span>}
                {command}
            </code>
            <button
                type="button"
                onClick={copy}
                className="flex items-center gap-1.5 border-l border-zinc-800 px-3 text-xs text-zinc-400 transition-colors hover:bg-zinc-900 hover:text-zinc-100"
                aria-label={`Salin: ${command}`}
            >
                {copied ? (
                    <Check className="size-3.5 text-emerald-400" />
                ) : (
                    <Copy className="size-3.5" />
                )}
                {copied ? 'Tersalin' : 'Salin'}
            </button>
        </div>
    );
}

type Props = {
    used: number;
    limit: number;
    segments?: number;
    label: string;
};

/**
 * Batang pemakaian berbentuk kotak-kotak. Kosong saat 0, kotak oranye bertambah
 * seiring pemakaian, dan semuanya merah saat jatah habis.
 */
export function filledSegments(used: number, limit: number, segments: number) {
    if (limit <= 0 || used <= 0) return 0;
    if (used >= limit) return segments;

    return Math.min(segments - 1, Math.ceil((used / limit) * segments));
}

export default function SegmentBar({
    used,
    limit,
    segments = 20,
    label,
}: Props) {
    const exhausted = limit > 0 && used >= limit;
    const filled = filledSegments(used, limit, segments);

    return (
        <div
            role="progressbar"
            aria-label={label}
            aria-valuemin={0}
            aria-valuemax={limit}
            aria-valuenow={Math.min(used, limit)}
            className="flex gap-1"
        >
            {Array.from({ length: segments }, (_, i) => (
                <span
                    key={i}
                    className={`h-6 flex-1 border ${
                        i < filled
                            ? exhausted
                                ? 'border-red-600 bg-red-600'
                                : 'border-orange-600 bg-orange-600'
                            : 'border-sidebar-border bg-muted/40'
                    }`}
                />
            ))}
        </div>
    );
}

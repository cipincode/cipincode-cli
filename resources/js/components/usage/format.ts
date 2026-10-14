/** "pukul 19.45" untuk hari ini, selain itu "Sen, 20 Okt 00.00". */
export function formatReset(iso: string): string {
    const date = new Date(iso);
    const time = date
        .toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
        .replace(':', '.');

    if (date.toDateString() === new Date().toDateString()) {
        return `pukul ${time}`;
    }

    const day = date.toLocaleDateString('id-ID', {
        weekday: 'short',
        day: 'numeric',
        month: 'short',
    });

    return `${day} ${time}`;
}

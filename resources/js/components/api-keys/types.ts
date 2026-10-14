export type ApiKey = {
    id: number;
    name: string;
    prefix: string;
    is_active: boolean;
    created_at: string;
    last_used_at: string | null;
};

export type ApiKeyStats = {
    total: number;
    active: number;
    inactive: number;
    used_7d: number;
};

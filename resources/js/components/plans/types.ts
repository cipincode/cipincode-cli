export type PlanInfo = {
    key: string;
    label: string;
    price: number;
    price_label: string;
    period: string | null;
    session_limit: number;
    weekly_limit: number | null;
    features: string[];
    paid: boolean;
};

export type CurrentPlan = {
    name: string;
    label: string;
    expires_at: string | null;
    expired: boolean;
};

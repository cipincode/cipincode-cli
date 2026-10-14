import { usePage } from '@inertiajs/react';

export default function AppLogo() {
    const { name } = usePage().props;

    return (
        <>
            <span className="flex h-10 w-full items-center justify-center group-data-[collapsible=icon]:hidden">
                <img
                    src="/images/logo-login-light.png"
                    alt={name}
                    className="h-full w-full object-contain dark:hidden"
                />
                <img
                    src="/images/logo-login-dark.png"
                    alt={name}
                    className="hidden h-full w-full object-contain dark:block"
                />
            </span>
            <img
                src="/favicon.png"
                alt={name}
                className="hidden size-7 object-contain group-data-[collapsible=icon]:block"
            />
        </>
    );
}

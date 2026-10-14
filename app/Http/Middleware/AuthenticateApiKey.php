<?php

namespace App\Http\Middleware;

use App\Models\ApiKey;
use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class AuthenticateApiKey
{
    /**
     * Authenticate the request with a bearer API key.
     */
    public function handle(Request $request, Closure $next): Response
    {
        $token = $request->bearerToken();
        $key = $token ? ApiKey::findByPlain($token) : null;

        if (! $key || ! $key->is_active) {
            return response()->json(['message' => 'Unauthenticated.'], 401);
        }

        if ($key->last_used_at === null || $key->last_used_at->lt(now()->subMinute())) {
            $key->forceFill(['last_used_at' => now()])->save();
        }

        $request->setUserResolver(fn () => $key->user);
        $request->attributes->set('api_key', $key);

        return $next($request);
    }
}

<x-mail::message>
# {{ __('Halo, :name!', ['name' => $user->name]) }}

{{ __('Terima kasih sudah mendaftar di :app. Klik tombol di bawah untuk memverifikasi alamat email Anda.', ['app' => config('app.name')]) }}

<x-mail::button :url="$url">
{{ __('Verifikasi Email') }}
</x-mail::button>

{{ __('Tautan ini berlaku selama :minutes menit.', ['minutes' => $expire]) }}

{{ __('Jika Anda tidak merasa membuat akun, abaikan saja email ini.') }}

{{ __('Salam,') }}<br>
{{ config('app.name') }}

<x-slot:subcopy>
{{ __('Jika tombol "Verifikasi Email" tidak berfungsi, salin dan tempel URL berikut ke browser Anda:') }}
<span class="break-all">[{{ $url }}]({{ $url }})</span>
</x-slot:subcopy>
</x-mail::message>

<?php

namespace App\Http\Requests;

use App\Models\ApiKey;
use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Validator;

class StoreApiKeyRequest extends FormRequest
{
    /**
     * @return array<string, ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:50'],
        ];
    }

    /**
     * @return array<int, callable>
     */
    public function after(): array
    {
        return [
            function (Validator $validator): void {
                if ($this->user()->apiKeys()->count() >= ApiKey::MAX_PER_USER) {
                    $validator->errors()->add(
                        'name',
                        __('Batas :max API key tercapai. Cabut key lama terlebih dahulu.', ['max' => ApiKey::MAX_PER_USER]),
                    );
                }
            },
        ];
    }
}

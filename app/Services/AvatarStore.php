<?php

namespace App\Services;

use App\Models\User;
use Illuminate\Http\UploadedFile;
use RuntimeException;

/**
 * Menyimpan foto profil ke public/uploads/images/avatar. Gambar dikompres
 * (WebP) tanpa mengubah resolusi, dan metadata EXIF dibuang.
 */
class AvatarStore
{
    public const DIR = 'uploads/images/avatar';

    public const QUALITY = 78;

    public function store(User $user, UploadedFile $file): string
    {
        $image = @imagecreatefromstring((string) file_get_contents($file->getRealPath()));
        if ($image === false) {
            throw new RuntimeException('Gambar tidak dapat dibaca.');
        }

        $image = $this->orient($image, $file);
        imagepalettetotruecolor($image);
        imagealphablending($image, false);
        imagesavealpha($image, true);

        $dir = public_path(self::DIR);
        if (! is_dir($dir)) {
            mkdir($dir, 0755, true);
        }

        $name = $user->getKey().'-'.bin2hex(random_bytes(6)).'.webp';
        $ok = imagewebp($image, $dir.'/'.$name, self::QUALITY);
        imagedestroy($image);

        if (! $ok) {
            throw new RuntimeException('Gambar gagal disimpan.');
        }

        $this->delete($user);

        $path = '/'.self::DIR.'/'.$name;
        $user->forceFill(['avatar' => $path])->save();

        return $path;
    }

    public function delete(User $user): void
    {
        $old = $user->avatar;

        if ($old && str_starts_with($old, '/'.self::DIR.'/')) {
            $file = public_path(ltrim($old, '/'));
            if (is_file($file)) {
                @unlink($file);
            }
        }
    }

    public function remove(User $user): void
    {
        $this->delete($user);
        $user->forceFill(['avatar' => null])->save();
    }

    private function orient(\GdImage $image, UploadedFile $file): \GdImage
    {
        if (! function_exists('exif_read_data') || $file->getMimeType() !== 'image/jpeg') {
            return $image;
        }

        $exif = @exif_read_data($file->getRealPath());
        $rotated = match ($exif['Orientation'] ?? 1) {
            3 => imagerotate($image, 180, 0),
            6 => imagerotate($image, -90, 0),
            8 => imagerotate($image, 90, 0),
            default => $image,
        };

        return $rotated ?: $image;
    }
}

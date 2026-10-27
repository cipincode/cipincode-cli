$ErrorActionPreference = "Stop"

# Cipin CLI installer (Windows). Installs straight from GitHub, no npm publish needed.
#   irm https://raw.githubusercontent.com/cipincode/cipincode-cli/main/cipin-cli/install.ps1 | iex
# Optional: $env:CIPIN_REF = "branch-or-tag" (default: main), $env:CIPIN_DIR = install folder.

foreach ($name in "node", "npm") {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        throw "$name tidak ditemukan. Pasang Node.js 18 atau lebih baru dulu: https://nodejs.org"
    }
}

$ref = if ($env:CIPIN_REF) { $env:CIPIN_REF } else { "main" }
$dest = if ($env:CIPIN_DIR) { $env:CIPIN_DIR } else { Join-Path $HOME ".cipin\cli" }
$tmp = Join-Path ([IO.Path]::GetTempPath()) ("cipin-" + [guid]::NewGuid().ToString("N"))

Write-Host ("Node: " + (node -v))
Write-Host "Mengunduh Cipin CLI ($ref) dari GitHub..."

New-Item -ItemType Directory -Force $tmp | Out-Null
try {
    $zip = Join-Path $tmp "cipin.zip"
    Invoke-WebRequest -UseBasicParsing "https://github.com/cipincode/cipincode-cli/archive/$ref.zip" -OutFile $zip
    Expand-Archive -Path $zip -DestinationPath $tmp -Force

    $source = Get-ChildItem -Path $tmp -Directory | ForEach-Object { Join-Path $_.FullName "cipin-cli" } | Where-Object { Test-Path $_ } | Select-Object -First 1
    if (-not $source) {
        throw "Folder cipin-cli tidak ditemukan di GitHub (ref: $ref)."
    }

    if (Test-Path $dest) { Remove-Item -Recurse -Force $dest }
    New-Item -ItemType Directory -Force (Split-Path $dest) | Out-Null
    Copy-Item -Recurse -Force $source $dest
}
finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}

Push-Location $dest
try {
    Write-Host "Memasang dependensi (OpenCode)..."
    npm install --omit=dev --no-audit --no-fund
    if ($LASTEXITCODE -ne 0) { throw "npm install gagal (kode $LASTEXITCODE)." }

    # A previous install leaves a registration behind (and its folder was just replaced): drop it first.
    $strict = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    npm uninstall -g cipin-cli *> $null
    $ErrorActionPreference = $strict
    npm install -g .
    if ($LASTEXITCODE -ne 0) { throw "Gagal mendaftarkan perintah cipin (kode $LASTEXITCODE)." }
}
finally {
    Pop-Location
}

Write-Host ""
Write-Host "Selesai. Terpasang di $dest"
Write-Host "Buka terminal baru, lalu:"
Write-Host "  cipin login"
Write-Host "  cipin"

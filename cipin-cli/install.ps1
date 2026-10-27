$ErrorActionPreference = "Stop"

foreach ($name in "node", "npm") {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        throw "$name tidak ditemukan. Pasang Node.js 18 atau lebih baru dulu: https://nodejs.org"
    }
}

Write-Host ("Node: " + (node -v))
Write-Host "Memasang cipin-cli..."
npm install -g cipin-cli
if ($LASTEXITCODE -ne 0) {
    throw "Gagal memasang cipin-cli (npm keluar dengan kode $LASTEXITCODE)."
}

Write-Host ""
Write-Host "Selesai. Langkah berikutnya:"
Write-Host "  cipin login"
Write-Host "  cipin"

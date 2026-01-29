Param(
  [string]$Arch = "amd64"
)

$ErrorActionPreference = "Stop"

Write-Host "==> Installing frontend dependencies"
Push-Location "frontend"
npm install
npm run build
Pop-Location

Write-Host "==> Building Wails binary"
wails build -platform "windows/$Arch"

Write-Host "==> Build complete. Output is in ./build/bin/"

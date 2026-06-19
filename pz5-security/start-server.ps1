$root = $PSScriptRoot
if (-not (Test-Path "$root\certs\server.crt")) {
    Write-Host "Сертификаты не найдены, запускаю generate-certs.ps1..."
    & "$root\generate-certs.ps1"
}

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "Set-Location '$root'; Write-Host 'HTTPS :8443, redirect HTTP :8085' -ForegroundColor Green; go run ./cmd/server"
)

Write-Host "Сервер запущен. HTTPS: https://localhost:8443"

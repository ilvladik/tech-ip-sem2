$root = $PSScriptRoot
Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "Set-Location '$root'; Write-Host 'pz4-monitoring :8084' -ForegroundColor Green; go run ./cmd/server"
)
Write-Host "Go-приложение запускается на :8084"

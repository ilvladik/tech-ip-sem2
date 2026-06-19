$root = $PSScriptRoot

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "Set-Location '$root'; Write-Host 'gRPC server :50051' -ForegroundColor Green; go run ./cmd/server"
)

Write-Host "Сервер запущен в отдельном окне (:50051)."

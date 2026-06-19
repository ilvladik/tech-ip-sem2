$root = $PSScriptRoot

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "Set-Location '$root\user-service'; Write-Host 'user-service :8081' -ForegroundColor Green; go run ./cmd/server"
)

Start-Sleep -Seconds 1

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "Set-Location '$root\order-service'; Write-Host 'order-service :8082' -ForegroundColor Green; go run ./cmd/server"
)

Write-Host "Запущены user-service (:8081) и order-service (:8082) в отдельных окнах."

# Сначала запустите Go-приложение: .\start-server.ps1

1..20 | ForEach-Object { curl.exe -s http://localhost:8084/health | Out-Null }
1..15 | ForEach-Object { curl.exe -s http://localhost:8084/students/1 | Out-Null }
1..5  | ForEach-Object { curl.exe -s http://localhost:8084/students/999 | Out-Null }
Write-Host "Трафик сгенерирован. Обновите дашборд Grafana."

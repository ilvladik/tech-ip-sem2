# Запустите: .\start-server.ps1
#
# В памяти есть студенты id 1–3. id 999 не заведён — проверка 404.

function Show-Response($title, $url) {
    Write-Host "=== $title ==="
    curl.exe -s -w "`nHTTP: %{http_code}`n" $url
    Write-Host ""
}

Show-Response "GET /health" "http://localhost:8084/health"
Show-Response "GET /students/1 — есть в системе" "http://localhost:8084/students/1"
Show-Response "GET /students/999 — нет в системе (ожидается 404)" "http://localhost:8084/students/999"

Write-Host "=== GET /metrics (фрагмент) ==="
curl.exe -s http://localhost:8084/metrics | Select-String "app_http"
Write-Host ""

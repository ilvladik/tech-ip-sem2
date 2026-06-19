# Нужны: PostgreSQL (.\start-db.ps1), сертификаты, .\start-server.ps1
#
# В БД после init.sql есть студенты id 1–3. id 999 не заведён — проверка 404.

function Show-Response($title, $url) {
    Write-Host "=== $title ==="
    curl.exe -k -s -w "`nHTTP: %{http_code}`n" $url
    Write-Host ""
}

Show-Response "HTTPS /health" "https://localhost:8443/health"

Write-Host "=== HTTP redirect -> HTTPS ==="
curl.exe -s -o NUL -w "HTTP:%{http_code} -> %{redirect_url}`n" http://localhost:8085/health
Write-Host ""

Show-Response "GET /students?id=1 — есть в БД" "https://localhost:8443/students?id=1"
Show-Response "GET /students?id=999 — нет в БД (ожидается 404)" "https://localhost:8443/students?id=999"
Show-Response "GET /students/by-email?email=ivanov@example.com" "https://localhost:8443/students/by-email?email=ivanov@example.com"
Show-Response "GET /students/by-email — неверный email (ожидается 400)" "https://localhost:8443/students/by-email?email=not-an-email"
Show-Response "Демо небезопасного SQL /students/unsafe?id=1" "https://localhost:8443/students/unsafe?id=1"

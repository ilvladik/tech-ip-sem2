# Запустите сервер: .\start-server.ps1
#
# В памяти есть студенты id 1–3. id 999 не заведён — проверка 404.

function Show-Response($title, $url, $method = "GET", $body = $null) {
    Write-Host "=== $title ==="
    if ($body) {
        curl.exe -s -w "`nHTTP: %{http_code}`n" -X $method $url -H "Content-Type: application/json" -d $body
    } else {
        curl.exe -s -w "`nHTTP: %{http_code}`n" $url
    }
    Write-Host ""
}

Show-Response "GET /health" "http://localhost:8083/health"
Show-Response "GET /students/1 — есть в системе" "http://localhost:8083/students/1"
Show-Response "GET /students/abc — неверный id (ожидается 400)" "http://localhost:8083/students/abc"
Show-Response "GET /students/999 — нет в системе (ожидается 404)" "http://localhost:8083/students/999"
Show-Response "POST /students — создать" "http://localhost:8083/students" "POST" '{"full_name":"Козлов Дмитрий","group":"ИТТ-04-25","email":"kozlov@example.com"}'
Show-Response "POST /students — дубликат email (ожидается 409)" "http://localhost:8083/students" "POST" '{"full_name":"Test","group":"G","email":"kozlov@example.com"}'

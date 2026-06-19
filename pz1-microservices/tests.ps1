# Запустите user-service (:8081) и order-service (:8082), затем:
# cd "c:\Github\Go_Practice_2_sem\pz1-microservices"
# .\tests.ps1
#
# В памяти есть пользователи id 1–3 и заказы 101–103.
# id 999 нигде не заведён — это проверка ответа «не найдено» (404).

function Show-Response($title, $url) {
    Write-Host "=== $title ==="
    curl.exe -s -w "`nHTTP: %{http_code}`n" $url
    Write-Host ""
}

Show-Response "GET /users — список" "http://localhost:8081/users"
Show-Response "GET /users/1 — есть в системе" "http://localhost:8081/users/1"
Show-Response "GET /users/999 — нет в системе (ожидается 404)" "http://localhost:8081/users/999"

Show-Response "GET /orders/101 — есть в системе" "http://localhost:8082/orders/101"
Show-Response "GET /orders/101/full — заказ + пользователь" "http://localhost:8082/orders/101/full"
Show-Response "GET /orders/by-user/1" "http://localhost:8082/orders/by-user/1"
Show-Response "GET /orders/by-user/2" "http://localhost:8082/orders/by-user/2"
Show-Response "GET /orders/999 — нет в системе (ожидается 404)" "http://localhost:8082/orders/999"

Write-Host "=== 502 (вручную: остановите user-service) ==="
Write-Host 'curl.exe -s -w "`nHTTP: %{http_code}`n" http://localhost:8082/orders/101/full'
Write-Host ""

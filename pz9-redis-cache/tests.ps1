Set-Location $PSScriptRoot
$base = "http://localhost:8089"

echo "1) GET /v1/tasks/1 (first request - may be miss)"
curl.exe -s "$base/v1/tasks/1"
echo ""

echo "2) GET /v1/tasks/1 (second request - expect hit in server logs)"
curl.exe -s "$base/v1/tasks/1"
echo ""

echo "3) GET /v1/tasks?page=1&limit=10"
curl.exe -s "$base/v1/tasks?page=1&limit=10"
echo ""

echo "4) PATCH /v1/tasks/1"
curl.exe -s -X PATCH "$base/v1/tasks/1" -H "Content-Type: application/json" -d "{\"id\":1,\"title\":\"Обновленная задача\",\"description\":\"Новый текст\",\"due_date\":\"2026-01-22T00:00:00Z\"}" -w " HTTP:%{http_code}"
echo ""

echo "5) GET /v1/tasks/1 after patch"
curl.exe -s "$base/v1/tasks/1"
echo ""

echo "6) DELETE /v1/tasks/1"
curl.exe -s -X DELETE "$base/v1/tasks/1" -w " HTTP:%{http_code}"
echo ""

echo "7) GET /v1/tasks/1 after delete (expect 404)"
curl.exe -s -w " HTTP:%{http_code}" "$base/v1/tasks/1"
echo ""

echo "8) GET /v1/tasks/2 (fallback works even if Redis is down)"
curl.exe -s -w " HTTP:%{http_code}" "$base/v1/tasks/2"
echo ""

echo "Done. Check server logs for cache hit/miss/set/invalidated."

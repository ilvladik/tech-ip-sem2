Set-Location $PSScriptRoot
$tasks = "http://localhost:8088"
$auth = "http://localhost:8087"

echo "1) auth health"
curl.exe -s "$auth/health"
echo ""

echo "2) tasks health"
curl.exe -s "$tasks/health"
echo ""

echo "3) tasks without token (expect 401)"
curl.exe -s -w " HTTP:%{http_code}" "$tasks/v1/tasks"
echo ""

echo "4) tasks with valid token"
curl.exe -i "$tasks/v1/tasks" -H "Authorization: Bearer demo-token" -H "X-Request-ID: pz7-001"

echo "5) tasks with invalid token (expect 401)"
curl.exe -s -w " HTTP:%{http_code}" "$tasks/v1/tasks" -H "Authorization: Bearer wrong"

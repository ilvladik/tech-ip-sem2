Set-Location $PSScriptRoot
$lb = "http://localhost:8090"

echo "1) health via NGINX"
curl.exe -i "$lb/health"
echo ""

echo "2) whoami via NGINX (10 requests - watch X-Instance-ID)"
1..10 | ForEach-Object {
    curl.exe -s -D - "$lb/whoami" -o NUL | Select-String "X-Instance-ID|HTTP/"
}

echo "3) tasks with Authorization header"
curl.exe -i "$lb/v1/tasks" -H "Authorization: Bearer demo-token"
echo ""

echo "Done. Stop tasks_1: docker compose -f deploy/lb/docker-compose.yml stop tasks_1"

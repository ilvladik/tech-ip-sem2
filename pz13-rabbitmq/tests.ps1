Set-Location $PSScriptRoot
$base = "http://localhost:8096"

echo "1) health"
curl.exe -s "$base/health"
echo ""

echo "2) POST /v1/tasks (publish task.created)"
curl.exe -i -X POST "$base/v1/tasks" `
  -H "Content-Type: application/json" `
  -H "Authorization: Bearer demo-token" `
  -H "X-Request-ID: pz13-001" `
  -d "{\"title\":\"Rabbit\",\"description\":\"publish event\"}"

echo ""
echo "Check worker logs for: received event=task.created"
echo "RabbitMQ UI: http://localhost:15672 queue task_events"

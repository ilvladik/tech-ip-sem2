Set-Location $PSScriptRoot
$base = "http://localhost:8097"

echo "1) success job t_001"
curl.exe -i -X POST "$base/v1/jobs/process-task" -H "Content-Type: application/json" -H "Authorization: Bearer demo-token" -d "{\"task_id\":\"t_001\"}"

echo ""
echo "2) fail job t_fail (3 retries then DLQ)"
curl.exe -i -X POST "$base/v1/jobs/process-task" -H "Content-Type: application/json" -H "Authorization: Bearer demo-token" -d "{\"task_id\":\"t_fail\"}"

echo ""
echo "Check worker logs and RabbitMQ UI queue task_jobs_dlq"

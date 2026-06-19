Set-Location $PSScriptRoot
$env:RABBIT_URL = "amqp://guest:guest@localhost:5672/"
$env:QUEUE_NAME = "task_events"
$env:PREFETCH = "1"
echo "Starting worker (queue task_events)"
go run ./services/worker/cmd/worker

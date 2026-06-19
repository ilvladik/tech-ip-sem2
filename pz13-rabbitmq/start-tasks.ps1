Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8096"
$env:RABBIT_URL = "amqp://guest:guest@localhost:5672/"
$env:QUEUE_NAME = "task_events"
$env:PUBLISH_MODE = "best_effort"
echo "Starting tasks service on http://localhost:8096"
go run ./services/tasks/cmd/tasks

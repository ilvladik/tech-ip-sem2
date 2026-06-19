Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8089"
$env:REDIS_ADDR = "localhost:6379"
echo "Starting pz9-redis-cache on http://localhost:8089"
go run ./cmd/server

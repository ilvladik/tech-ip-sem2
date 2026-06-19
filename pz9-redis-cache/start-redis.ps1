Set-Location "$PSScriptRoot\deploy\redis"
echo "Starting Redis on localhost:6379"
docker compose up -d
docker compose ps

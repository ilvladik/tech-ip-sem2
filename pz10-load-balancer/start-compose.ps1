Set-Location "$PSScriptRoot\deploy\lb"
echo "Starting NGINX :8090 + tasks_1/2/3"
docker compose up -d --build
docker compose ps

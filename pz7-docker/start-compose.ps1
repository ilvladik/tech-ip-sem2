Set-Location "$PSScriptRoot\deploy"
echo "Building and starting auth:8087 + tasks:8088"
docker compose up -d --build

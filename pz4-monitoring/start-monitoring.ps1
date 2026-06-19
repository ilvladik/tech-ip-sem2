Set-Location $PSScriptRoot
docker compose up -d
Write-Host "Prometheus: http://localhost:9090"
Write-Host "Grafana:    http://localhost:3000  (admin / admin)"

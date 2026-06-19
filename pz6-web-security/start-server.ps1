Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8086"
echo "Starting pz6-web-security on http://localhost:8086"
go run ./cmd/server

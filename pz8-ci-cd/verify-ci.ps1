Set-Location $PSScriptRoot\..
$ErrorActionPreference = "Stop"

echo "=== auth: test + build ==="
Set-Location pz7-docker\services\auth
go mod tidy
go test ./...
go build ./...

echo "=== tasks: test + build ==="
Set-Location ..\tasks
go mod tidy
go test ./...
go build ./...

echo "CI checks passed (docker build skipped - run on GitHub Actions)"

Set-Location $PSScriptRoot
New-Item -ItemType Directory -Force -Path bin | Out-Null
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o bin/tasks ./cmd/tasks
echo "Built bin/tasks for Linux amd64"
echo "Upload: scp bin/tasks user@VPS_IP:/tmp/tasks"

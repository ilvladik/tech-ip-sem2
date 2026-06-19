Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8097"
go run ./services/tasks/cmd/tasks

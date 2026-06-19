Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8095"
echo "REST + GraphQL on http://localhost:8095"
go run ./cmd/server

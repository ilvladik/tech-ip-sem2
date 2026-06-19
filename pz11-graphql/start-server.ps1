Set-Location $PSScriptRoot
$env:SERVER_ADDR = ":8094"
echo "GraphQL Playground: http://localhost:8094/"
go run ./cmd/graphql

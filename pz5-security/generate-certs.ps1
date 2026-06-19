Set-Location $PSScriptRoot
go run ./cmd/gencert
Write-Host "TLS-сертификаты готовы в certs/"

# Генерация Go-кода из proto (нужны protoc, protoc-gen-go, protoc-gen-go-grpc в PATH)
# google/protobuf/empty.proto лежит локально в proto/ — ничего дополнительно качать не нужно

$root = $PSScriptRoot
New-Item -ItemType Directory -Path "$root\gen" -Force | Out-Null

protoc `
  --proto_path="$root\proto" `
  --go_out="$root\gen" --go_opt=module=example.com/pz2-grpc/gen `
  --go-grpc_out="$root\gen" --go-grpc_opt=module=example.com/pz2-grpc/gen `
  "$root\proto\student.proto"

Write-Host "Done. Files in gen/studentpb/"

Set-Location $PSScriptRoot
echo "Build auth image"
docker build -t techip-auth:0.1 ./services/auth
echo "Build tasks image"
docker build -t techip-tasks:0.1 ./services/tasks

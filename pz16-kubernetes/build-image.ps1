Set-Location $PSScriptRoot
docker build -t techip-tasks:0.1 .
echo "Image techip-tasks:0.1 built"
echo "For kind: kind load docker-image techip-tasks:0.1"

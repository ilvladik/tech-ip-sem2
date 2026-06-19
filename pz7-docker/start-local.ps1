Set-Location $PSScriptRoot
$authCmd = "Set-Location '$PSScriptRoot\services\auth'; `$env:AUTH_PORT='8087'; go run ./cmd/auth"
$tasksCmd = "Set-Location '$PSScriptRoot\services\tasks'; `$env:TASKS_PORT='8088'; `$env:AUTH_BASE_URL='http://localhost:8087'; go run ./cmd/tasks"
Start-Process powershell -ArgumentList "-NoExit", "-Command", $authCmd
Start-Sleep -Seconds 2
Start-Process powershell -ArgumentList "-NoExit", "-Command", $tasksCmd
echo "auth :8087 and tasks :8088 started in separate windows"

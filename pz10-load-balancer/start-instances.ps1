Set-Location "$PSScriptRoot\services\tasks"
if (-not (Test-Path .\tasks.exe)) { go build -o tasks.exe ./cmd/server }

$root = (Get-Location).Path
Start-Process cmd.exe -ArgumentList '/c',"set APP_PORT=8091&& set INSTANCE_ID=tasks-1&& tasks.exe" -WorkingDirectory $root -WindowStyle Normal
Start-Sleep 1
Start-Process cmd.exe -ArgumentList '/c',"set APP_PORT=8092&& set INSTANCE_ID=tasks-2&& tasks.exe" -WorkingDirectory $root -WindowStyle Normal
Start-Sleep 1
Start-Process cmd.exe -ArgumentList '/c',"set APP_PORT=8093&& set INSTANCE_ID=tasks-3&& tasks.exe" -WorkingDirectory $root -WindowStyle Normal

echo "Started tasks-1 :8091, tasks-2 :8092, tasks-3 :8093"
echo "Use tests-instances.ps1 (direct) or Docker + tests-lb.ps1 for NGINX on :8090"

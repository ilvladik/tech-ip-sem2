Set-Location $PSScriptRoot
echo "Direct requests to instances (round-robin simulation)"

1..3 | ForEach-Object {
    echo "--- tasks-1 :8091 ---"
    curl.exe -s -D - http://localhost:8091/v1/tasks -o NUL | Select-String "X-Instance-ID"
}

1..3 | ForEach-Object {
    echo "--- tasks-2 :8092 ---"
    curl.exe -s -D - http://localhost:8092/v1/tasks -o NUL | Select-String "X-Instance-ID"
}

1..3 | ForEach-Object {
    echo "--- tasks-3 :8093 ---"
    curl.exe -s -D - http://localhost:8093/v1/tasks -o NUL | Select-String "X-Instance-ID"
}

echo "GET /whoami on each instance"
curl.exe -s http://localhost:8091/whoami
echo ""
curl.exe -s http://localhost:8092/whoami
echo ""
curl.exe -s http://localhost:8093/whoami
echo ""

echo "GET /health"
curl.exe -s http://localhost:8091/health
echo ""

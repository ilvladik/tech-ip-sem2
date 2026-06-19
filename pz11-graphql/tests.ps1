Set-Location $PSScriptRoot
$base = "http://localhost:8094/query"

echo "1) Query tasks (no auth)"
$body1 = '{"query":"{ tasks { id title done } }"}'
curl.exe -s -X POST $base -H "Content-Type: application/json" -d $body1
echo ""

echo "2) Query task by id"
$body2 = '{"query":"query($id: ID!) { task(id: $id) { id title description done } }","variables":{"id":"t_001"}}'
curl.exe -s -X POST $base -H "Content-Type: application/json" -d $body2
echo ""

echo "3) Create task without token (expect error)"
$body3 = '{"query":"mutation($input: CreateTaskInput!) { createTask(input: $input) { id title } }","variables":{"input":{"title":"Test","description":"no auth"}}}'
curl.exe -s -X POST $base -H "Content-Type: application/json" -d $body3
echo ""

echo "4) Create task with Bearer demo-token"
$body4 = '{"query":"mutation($input: CreateTaskInput!) { createTask(input: $input) { id title done } }","variables":{"input":{"title":"Изучить GraphQL","description":"Практика 11"}}}'
curl.exe -s -X POST $base -H "Content-Type: application/json" -H "Authorization: Bearer demo-token" -d $body4
echo ""

echo "5) Update task"
$body5 = '{"query":"mutation($id: ID!, $input: UpdateTaskInput!) { updateTask(id: $id, input: $input) { id title done } }","variables":{"id":"t_001","input":{"done":true}}}'
curl.exe -s -X POST $base -H "Content-Type: application/json" -H "Authorization: Bearer demo-token" -d $body5
echo ""

echo "Done."

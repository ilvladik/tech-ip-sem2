Set-Location $PSScriptRoot
$base = "http://localhost:8095"

echo "=== REST: list (full objects - over-fetching demo) ==="
curl.exe -s "$base/v1/tasks"
echo ""

echo "=== REST: get one ==="
curl.exe -s "$base/v1/tasks/t_001"
echo ""

echo "=== REST: create ==="
curl.exe -s -X POST "$base/v1/tasks" -H "Content-Type: application/json" -d "{\"title\":\"Сравнение REST и GraphQL\",\"description\":\"Практика 12\"}"
echo ""

echo "=== REST: patch done ==="
curl.exe -s -X PATCH "$base/v1/tasks/t_001" -H "Content-Type: application/json" -d "{\"done\":true}"
echo ""

echo "=== REST: not found ==="
curl.exe -s -w " HTTP:%{http_code}" "$base/v1/tasks/unknown"
echo ""
echo ""

echo "=== GraphQL: list (only id, title, done) ==="
$q1 = '{"query":"{ tasks { id title done } }"}'
curl.exe -s -X POST "$base/graphql" -H "Content-Type: application/json" -d $q1
echo ""

echo "=== GraphQL: detail ==="
$q2 = '{"query":"query($id: ID!) { task(id: $id) { id title description done } }","variables":{"id":"t_001"}}'
curl.exe -s -X POST "$base/graphql" -H "Content-Type: application/json" -d $q2
echo ""

echo "=== GraphQL: create ==="
$q3 = '{"query":"mutation($input: CreateTaskInput!) { createTask(input: $input) { id title done } }","variables":{"input":{"title":"GraphQL task","description":"pz12"}}}'
curl.exe -s -X POST "$base/graphql" -H "Content-Type: application/json" -d $q3
echo ""

echo "=== GraphQL: unknown id (null) ==="
$q4 = '{"query":"query($id: ID!) { task(id: $id) { id } }","variables":{"id":"unknown"}}'
curl.exe -s -X POST "$base/graphql" -H "Content-Type: application/json" -d $q4
echo ""

echo "Done."

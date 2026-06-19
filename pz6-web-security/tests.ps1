Set-Location $PSScriptRoot
$port = if ($env:SERVER_ADDR) { $env:SERVER_ADDR.TrimStart(':') } else { "8086" }
$base = "http://localhost:$port"
$cookie = "cookies.txt"

if (Test-Path $cookie) { Remove-Item $cookie }

echo "1) Login and save session cookie"
curl.exe -s -c $cookie "$base/login" -o NUL
curl.exe -s -b $cookie "$base/profile" -o profile.html
echo "profile page saved to profile.html"

echo "2) GET /hello"
curl.exe -s -b $cookie "$base/hello" -o hello.html
echo "hello page saved to hello.html"

echo "3) POST /profile without CSRF (expect 403)"
curl.exe -s -b $cookie -X POST "$base/profile" -d "name=Hacker" -w "`nHTTP %{http_code}`n"

echo "4) POST /profile with CSRF from profile.html"
$csrf = (Select-String -Path profile.html -Pattern 'name="csrf_token" value="([^"]+)"').Matches[0].Groups[1].Value
curl.exe -s -b $cookie -c $cookie -X POST "$base/profile" -d "csrf_token=$csrf&name=Alice" -o NUL -w "HTTP %{http_code}`n"

echo "5) GET /comments"
curl.exe -s -b $cookie "$base/comments" -o comments.html
echo "comments page saved to comments.html"

echo "6) POST comment with CSRF"
$csrf2 = (Select-String -Path comments.html -Pattern 'name="csrf_token" value="([^"]+)"').Matches[0].Groups[1].Value
curl.exe -s -b $cookie -c $cookie -X POST "$base/comments" -d "csrf_token=$csrf2&text=<script>alert(1)</script>" -o NUL -w "HTTP %{http_code}`n"
curl.exe -s -b $cookie "$base/comments" -o comments-after.html
echo "comments-after.html saved (script should be escaped in HTML)"

echo "7) GET /demo/unsafe?name=<b>XSS</b>"
curl.exe -s "$base/demo/unsafe?name=%3Cb%3EXSS%3C/b%3E"

echo "8) GET /logout"
curl.exe -s -b $cookie -c $cookie "$base/logout" -o NUL -w "HTTP %{http_code}`n"

echo "Done."

go test ./... -coverprofile="coverage.out"

Get-Content "coverage.out" | Where-Object { $_ -notmatch "cmd|app|setup|errors|mocks|gen" } | Set-Content "coverage1.out"

go tool cover -func="coverage1.out"

echo "Process complete"
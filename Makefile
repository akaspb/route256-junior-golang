
coverage:
	go test ./cmd/cli -coverprofile=coverage
	go tool cover -html coverage -o index.html

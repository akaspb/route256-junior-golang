
coverage:
	go test ./cmd/cli -coverprofile=coverage1
	go tool cover -html coverage1 -o index1.html

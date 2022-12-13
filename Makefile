coverage:
	go test ./cmd/api -coverprofile=coverage.out && go tool cover -html=coverage.out

PHONY: coverage
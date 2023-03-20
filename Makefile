.PHONY: cli
cli:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/cli ./cmd/cli

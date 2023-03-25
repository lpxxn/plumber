.PHONY: cli
cli:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/cli ./cmd/cli

.PHONY: linux_cli
linux_cli:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/cli ./cmd/cli

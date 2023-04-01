.PHONY: cli
cli:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/plumber-cli ./cmd/plumber-cli

.PHONY: linux_cli
linux_cli:
	env CGO_ENABLED=0 GOOS=linux go build -o bin/plumber-cli ./cmd/plumber-cli

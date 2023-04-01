.PHONY: cli
cli:
	go build -o bin/plumber-cli ./cmd/plumber-cli

.PHONY: linux_cli
linux_arm_cli:
	env GOOS=linux GOARCH=arm64 go build -o bin/plumber-cli ./cmd/plumber-cli

linux_amd_cli:
	env GOOS=linux GOARCH=amd64 go build -o bin/plumber-cli ./cmd/plumber-cli

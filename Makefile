.PHONY: cli
cli:
	go build -o bin/plumber-cli ./cmd/plumber-cli

.PHONY: linux_arm_cli
linux_arm_cli:
	env GOOS=linux GOARCH=arm64 go build -o bin/plumber-cli ./cmd/plumber-cli

.PHONY: linux_amd_cli
linux_amd_cli:
	env GOOS=linux GOARCH=amd64 go build -o bin/plumber-cli ./cmd/plumber-cli

.PHONY: srv
srv:
	go build -o bin/plumber ./cmd/plumber

.PHONY: linux_arm_srv
linux_arm_srv:
	env GOOS=linux GOARCH=arm64 go build -o bin/plumber ./cmd/plumber

.PHONY: linux_amd_srv
linux_amd_srv:
	env GOOS=linux GOARCH=amd64 go build -o bin/plumber ./cmd/plumber

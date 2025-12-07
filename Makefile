# Makefile
.PHONY: build install test test-race test-coverage test-coverage-html clean dev deps fmt lint setup-dev

# Build variables
APP_NAME := how
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X github.com/Codilas/how/pkg/version.Version=$(VERSION) \
                    -X github.com/Codilas/how/pkg/version.GitCommit=$(GIT_COMMIT) \
                    -X github.com/Codilas/how/pkg/version.BuildDate=$(BUILD_DATE)"

# Build the application
build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# Install to local bin
install: build
	@echo "Installing $(APP_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	cp bin/$(APP_NAME) ~/.local/bin/
	@echo "Installation complete!"
	@echo "Add ~/.local/bin to your PATH if not already present"
	@echo "Run '$(APP_NAME) setup' to configure"

# Install to system bin (requires sudo)
install-system: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	sudo cp bin/$(APP_NAME) /usr/local/bin/
	@echo "System installation complete!"

# Development build with race detection
dev:
	go build -race $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# Run tests
test:
	go test -v ./...

# Run tests with race detection
test-race:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "Coverage Summary:"
	@go tool cover -func=coverage.out | tail -1

# Generate HTML coverage report
test-coverage-html: test-coverage
	@go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

# Download dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Create release builds for multiple platforms
release:
	@echo "Building releases..."
	@mkdir -p releases
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o releases/$(APP_NAME)-linux-amd64 ./cmd/$(APP_NAME)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o releases/$(APP_NAME)-darwin-amd64 ./cmd/$(APP_NAME)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o releases/$(APP_NAME)-darwin-arm64 ./cmd/$(APP_NAME)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o releases/$(APP_NAME)-windows-amd64.exe ./cmd/$(APP_NAME)

# Development setup
setup-dev: deps
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development setup complete!"

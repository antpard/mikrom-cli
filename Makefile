# Makefile for Mikrom CLI

# Variables
APP_NAME=mikrom
VERSION?=1.0.0
BUILD_DIR=bin
COVERAGE_DIR=coverage
MAIN_PATH=main.go

# Configuration
.DEFAULT_GOAL := help
.PHONY: help run build build-linux build-all install test test-verbose test-short test-coverage coverage-html bench vet fmt fmt-check lint check deps tidy clean info

## help: Show this help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: Run the CLI
run:
	go run $(MAIN_PATH)

## build: Compile the CLI binary
build:
	@echo "Compiling $(APP_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Binary created: $(BUILD_DIR)/$(APP_NAME)"

## build-linux: Compile for Linux
build-linux:
	@echo "Compiling for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_PATH)

## build-all: Compile for multiple platforms
build-all:
	@echo "Compiling for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64  go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64   $(MAIN_PATH)
	GOOS=linux   GOARCH=arm64  go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64   $(MAIN_PATH)
	GOOS=darwin  GOARCH=amd64  go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64  $(MAIN_PATH)
	GOOS=darwin  GOARCH=arm64  go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64  $(MAIN_PATH)
	GOOS=windows GOARCH=amd64  go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Binaries created in $(BUILD_DIR)/"

## install: Install the binary to GOPATH/bin
install:
	@echo "Installing $(APP_NAME)..."
	go install -ldflags="-s -w -X main.Version=$(VERSION)" .
	@echo "✓ $(APP_NAME) installed to $$(go env GOPATH)/bin"

## test: Run all tests with race detection
test:
	@echo "Running tests..."
	go test ./... -race

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	go test ./... -v -race

## test-short: Run fast tests only
test-short:
	@echo "Running short tests..."
	go test ./... -short

## test-coverage: Run tests and generate coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	go test ./... -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic
	go tool cover -func=$(COVERAGE_DIR)/coverage.out
	@echo "\n✓ To view the HTML report run: make coverage-html"

## coverage-html: Generate HTML coverage report
coverage-html:
	@if [ ! -f $(COVERAGE_DIR)/coverage.out ]; then \
		echo "Run first: make test-coverage"; \
		exit 1; \
	fi
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✓ Report generated at $(COVERAGE_DIR)/coverage.html"

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test ./... -bench=. -benchmem

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## fmt-check: Verify code is formatted
fmt-check:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "The following files need formatting:" && gofmt -l . && exit 1)

## lint: Run linter (requires golangci-lint)
lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running linter..."; \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint is not installed."; \
		echo "Install it from: https://golangci-lint.run/usage/install/"; \
	fi

## check: Run all checks (fmt, vet, lint, test)
check: fmt-check vet lint test
	@echo "✓ All checks passed"

## deps: Download and tidy dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	go mod tidy
	@echo "✓ go.mod tidied"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	go clean -cache -testcache
	@echo "✓ Cleanup complete"

## info: Show project information
info:
	@echo "Project:   $(APP_NAME)"
	@echo "Version:   $(VERSION)"
	@echo "Go:        $$(go version)"
	@echo "Build dir: $(BUILD_DIR)"
	@echo "Go files:  $$(find . -name '*.go' -not -path './vendor/*' | wc -l)"
	@echo "Packages:  $$(go list ./... | wc -l)"

.PHONY: build test clean install dev help

# Build variables
BINARY_NAME=any2pdf
BUILD_DIR=build
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildDate=${BUILD_DATE}"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the project
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/any2pdf/main.go
	@echo "Done! Binary at ${BUILD_DIR}/${BINARY_NAME}"

build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 cmd/any2pdf/main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 cmd/any2pdf/main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 cmd/any2pdf/main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 cmd/any2pdf/main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe cmd/any2pdf/main.go
	@echo "Done! Binaries in ${BUILD_DIR}/"

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v -race -coverprofile=coverage.out

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test ./test -v -tags=integration

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@rm -f coverage.out coverage.html
	@go clean
	@echo "Done!"

install: build ## Install binary to $GOPATH/bin
	@echo "Installing to $(shell go env GOPATH)/bin..."
	@cp ${BUILD_DIR}/${BINARY_NAME} $(shell go env GOPATH)/bin/
	@echo "Done! Run 'any2pdf' from anywhere."

dev: ## Run in development mode (build and run with sample file)
	@go run cmd/any2pdf/main.go

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

check-deps: ## Check if external dependencies are installed
	@echo "Checking external dependencies..."
	@command -v pandoc >/dev/null 2>&1 || echo "⚠️  pandoc not found"
	@command -v libreoffice >/dev/null 2>&1 || echo "⚠️  libreoffice not found"
	@command -v img2pdf >/dev/null 2>&1 || echo "⚠️  img2pdf not found"
	@command -v wkhtmltopdf >/dev/null 2>&1 || echo "⚠️  wkhtmltopdf not found"
	@echo "Done!"

.DEFAULT_GOAL := help

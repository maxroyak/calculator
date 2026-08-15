.PHONY: all build test lint fmt vet race clean help

# Variables
BINARY_NAME ?= calculator
BUILD_DIR ?= ./build
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
PACKAGES := $(shell go list ./... | grep -v /vendor/)

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with race detector
race:
	@echo "Running tests with race detector..."
	@go test -race ./...

# Lint the code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet the code for potential issues
vet:
	@echo "Vetting code..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  make build   - Build the application"
	@echo "  make test    - Run tests"
	@echo "  make race    - Run tests with race detector"
	@echo "  make lint    - Run linter"
	@echo "  make fmt     - Format code"
	@echo "  make vet     - Vet code"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make help    - Show this help"

# Development workflow
dev: fmt vet test
	@echo "Development checks passed!"

# CI/CD pipeline step
ci: fmt vet test race
	@echo "CI checks passed!"
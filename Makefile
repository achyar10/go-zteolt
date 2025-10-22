.PHONY: build run clean test deps dev help

# Variables
APP_NAME := go-zteolt
SERVER_BINARY := bin/server
CLI_BINARY := bin/cli
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Install dependencies
	go mod download
	go mod tidy
	@which air > /dev/null || (echo "Installing air for hot reload..." && go install github.com/cosmtrek/air@latest)

dev: ## Run server in development mode with hot reload
	@echo "🔥 Starting development server with hot reload..."
	air

dev-simple: ## Run server in development mode without hot reload
	@echo "Starting development server..."
	go run cmd/server/main.go -dev

dev-build: ## Build server for development
	@echo "Building development server..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(SERVER_BINARY) cmd/server/main.go

# Build targets
build: build-server build-cli ## Build all binaries

build-server: ## Build server binary
	@echo "Building server..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(SERVER_BINARY) cmd/server/main.go

build-cli: ## Build CLI binary
	@echo "Building CLI..."
	@mkdir -p bin
	go build $(LDFLAGS) -o $(CLI_BINARY) main.go

# Run targets
run-server: build-server ## Run server binary
	./$(SERVER_BINARY)

run-cli: build-cli ## Run CLI binary
	./$(CLI_BINARY) --help

# Testing
test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Quality checks
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Docker targets
docker-build: ## Build Docker image
	docker build -t $(APP_NAME):$(VERSION) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

# Installation
install: build-server ## Install server binary
	sudo cp $(SERVER_BINARY) /usr/local/bin/$(APP_NAME)-server

install-cli: build-cli ## Install CLI binary
	sudo cp $(CLI_BINARY) /usr/local/bin/$(APP_NAME)

# Production
prod-build: ## Build for production
	@echo "Building for production..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo $(LDFLAGS) -o $(SERVER_BINARY) cmd/server/main.go

# Quick start
quickstart: deps build-server ## Quick start development setup
	@echo "Quick start complete!"
	@echo "Run: make dev to start development server"
	@echo "API will be available at: http://localhost:8080"
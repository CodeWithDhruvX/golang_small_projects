# AI Recruiter Assistant - Go Backend Makefile

.PHONY: build run test clean deps lint fmt docker

# Variables
APP_NAME := ai-recruiter-assistant
VERSION := 1.0.0
BUILD_DIR := build
MAIN_FILE := cmd/server/main.go

# Default target
all: deps build

# Install dependencies
deps:
	@echo "📦 Installing Go dependencies..."
	go mod tidy
	go mod download

# Build the application
build:
	@echo "🔨 Building $(APP_NAME)..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run the application
run:
	@echo "🚀 Running $(APP_NAME)..."
	go run $(MAIN_FILE)

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Generate API documentation
docs:
	@echo "📚 Generating API documentation..."
	swag init -g cmd/server/main.go

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .

# Run with hot reload (requires air)
dev:
	@echo "🔥 Running with hot reload..."
	air

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database migrations
migrate-up:
	@echo "⬆️ Running database migrations..."
	# This will be handled by the application startup

# Security scan
security:
	@echo "🔒 Running security scan..."
	gosec ./...

# Performance benchmark
benchmark:
	@echo "📈 Running benchmarks..."
	go test -bench=. -benchmem ./...

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docs          - Generate API documentation"
	@echo "  docker        - Build Docker image"
	@echo "  dev           - Run with hot reload"
	@echo "  install-tools - Install development tools"
	@echo "  security      - Run security scan"
	@echo "  benchmark     - Run benchmarks"
	@echo "  help          - Show this help message"

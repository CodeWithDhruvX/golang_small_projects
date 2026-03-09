# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix

# Build targets
.PHONY: build clean test coverage docker-build docker-run help

# Default target
all: build

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/server/main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run

# Run the application
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v ./cmd/server/main.go
	./bin/$(BINARY_NAME)

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_UNIX) -v ./cmd/server/main.go

# Build Docker image
docker-build:
	docker build -t knowledge-base-go-backend:latest .

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose
docker-stop:
	docker-compose down

# View Docker logs
docker-logs:
	docker-compose logs -f

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go

# Run race condition tests
race:
	$(GOTEST) -race -v ./...

# Run benchmark tests
bench:
	$(GOTEST) -bench=. -v ./...

# Update dependencies
update-deps:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Security check
security:
	gosec ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  run            - Build and run the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose"
	@echo "  docker-logs    - View Docker logs"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  race           - Run race condition tests"
	@echo "  bench          - Run benchmark tests"
	@echo "  update-deps    - Update dependencies"
	@echo "  security       - Run security check"
	@echo "  help           - Show this help message"

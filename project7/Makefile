.PHONY: help build run clean docker-up docker-down deps test

help:
	@echo "Available commands:"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make docker-up - Start Docker containers"
	@echo "  make docker-down - Stop and remove Docker containers"
	@echo "  make deps      - Download dependencies"
	@echo "  make test      - Run tests"

build:
	go build -o bin/project7 cmd/main.go

run:
	go run cmd/main.go

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

deps:
	go mod download
	go mod tidy

test:
	go test ./...

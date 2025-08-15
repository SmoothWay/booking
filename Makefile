.PHONY: test test-short test-verbose test-coverage test-benchmark test-integration test-unit clean build

# Test commands
test:
	go test ./...

test-short:
	go test -short ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

test-benchmark:
	go test -bench=. ./...

test-integration:
	go test -tags=integration ./...

test-unit:
	go test -short ./...

# Build commands
build:
	go build -o bin/booking cmd/api/main.go

build-test:
	go test -c ./...

# Clean commands
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Database commands
db-up:
	docker compose up -d postgres redis kafka

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d postgres redis kafka

# Development commands
dev:
	go run cmd/api/main.go

dev-test:
	go run cmd/api/main.go -test

# Linting and formatting
fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

# Dependencies
deps:
	go mod download
	go mod tidy

# Help
help:
	@echo "Available commands:"
	@echo "  test           - Run all tests"
	@echo "  test-short     - Run only unit tests (skip integration)"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-benchmark - Run benchmark tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  db-up          - Start database services"
	@echo "  db-down        - Stop database services"
	@echo "  db-reset       - Reset database services"
	@echo "  dev            - Run in development mode"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Run linter"
	@echo "  deps           - Download and tidy dependencies" 
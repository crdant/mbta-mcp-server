BINARY_NAME=mbta-mcp-server
VERSION=dev
MAIN_PACKAGE=./cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

.PHONY: all build clean test test-coverage lint vet fmt

all: clean fmt lint vet test build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	@echo "Cleaning up..."
	@rm -rf bin/ coverage.out coverage.html

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping linting"; \
	fi

vet:
	@echo "Running go vet..."
	@go vet ./...

fmt:
	@echo "Running go fmt..."
	@go fmt ./...

run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(LDFLAGS) $(MAIN_PACKAGE)

# Docker targets
image:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .

container:
	@echo "Running in Docker container..."
	@docker run --rm -e MBTA_API_KEY -p 8080:8080 $(BINARY_NAME):$(VERSION)

# Release targets
release:
	@echo "Creating release version $(VERSION)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)
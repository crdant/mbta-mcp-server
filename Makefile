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

# OCI image targets
package:
	@echo "Building package with melange..."
	@mkdir -p ./packages
	@melange build --arch amd64,arm64 \
		--signing-key melange.rsa \
		--keyring-append melange.rsa.pub \
		--out-dir ./packages \
		--repository-append ./packages \
		--version $(VERSION) \
		melange.yaml

image: package
	@echo "Building OCI image with apko..."
	@apko build \
		--keyring melange.rsa.pub \
		--arch amd64,arm64 \
		--repository ./packages \
		apko.yaml \
		$(BINARY_NAME):$(VERSION) \
		image.tar \
		sbom.json

container: image
	@echo "Loading image into Docker..."
	@docker load < image.tar
	@echo "Running container..."
	@docker run --rm -e MBTA_API_KEY -p 8080:8080 $(BINARY_NAME):$(VERSION)

keys:
	@echo "Generating signing keys..."
	@if [ ! -f melange.rsa ]; then \
		openssl genrsa -out melange.rsa 4096; \
		openssl rsa -in melange.rsa -pubout -out melange.rsa.pub; \
	fi

# Release targets
release:
	@echo "Creating release version $(VERSION)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)
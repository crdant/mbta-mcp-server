BINARY_NAME=mbta-mcp-server
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
# Get the version from semver-cli, removing any 'v' prefix if present
VERSION=$(shell semver get release 2>/dev/null | sed 's/^v//' || echo "0.1.0")
BUILD_VERSION=$(VERSION)+build.$(GIT_SHORT_SHA)
MAIN_PACKAGE=./cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS=-ldflags "-X main.Version=$(BUILD_VERSION)"

# Semver-cli should be available through devshell
# Fallback only if not in the development environment
semver-check:
	@command -v semver > /dev/null || (echo "semver not found, installing..." && \
	go install github.com/maykonlsf/semver-cli/cmd/semver@latest)

# Initialize semver if .semver.yaml doesn't exist
.semver.yaml:
	@echo "Initializing semver..."
	@semver init

# Version management targets
init-semver: .semver.yaml

alpha: .semver.yaml
	@echo "Incrementing to next alpha version..."
	@semver up alpha

beta: .semver.yaml
	@echo "Incrementing to next beta version..."
	@semver up beta

rc: .semver.yaml
	@echo "Incrementing to next release candidate version..."
	@semver up rc

patch: .semver.yaml
	@echo "Incrementing to next patch version..."
	@semver up release

minor: .semver.yaml
	@echo "Incrementing to next minor version..."
	@semver up minor

major: .semver.yaml
	@echo "Incrementing to next major version..."
	@semver up major

release: .semver.yaml
	@echo "Creating final release from pre-release..."
	@semver up release

tag: .semver.yaml
	@echo "Tagging current version in git..."
	@git tag -a "v$(VERSION)" -m "Version $(VERSION)"
	@echo "Tagged version $(VERSION)"

.PHONY: all build clean test test-coverage lint vet fmt init-semver alpha beta rc patch minor major release tag

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

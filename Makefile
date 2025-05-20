BINARY_NAME=mbta-mcp-server
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
# Get the version from semver-cli, removing any 'v' prefix if present
VERSION=$(shell semver get release 2>/dev/null | sed 's/^v//' || echo "0.1.0")
BUILD_VERSION=$(VERSION)+build.$(GIT_SHORT_SHA)
MAIN_PACKAGE=./cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS=-ldflags "-X main.Version=$(BUILD_VERSION)"
ARCHS=x86_64
KEYFILE=melange.rsa
KEYSIZE=4096

# Container runtime to use (docker, nerdctl, or podman)
# Can be overridden with either:
# - Environment variable: CONTAINER_RUNTIME=nerdctl make container
# - Make argument: make container CONTAINER_RUNTIME=podman
CONTAINER_RUNTIME ?= docker

# Semver-cli should be available through devshell
# Fallback only if not in the development environment
semver-check:
	@command -v semver > /dev/null || (echo "semver not found, installing..." && \
	go install github.com/maykonlsf/semver-cli/cmd/semver@latest)

# Initialize semver if .semver.yaml doesn't exist
.semver.yaml:
	@semver init

# Version management targets
init-semver: .semver.yaml

alpha: .semver.yaml
	@semver up alpha

beta: .semver.yaml
	@semver up beta

rc: .semver.yaml
	@semver up rc

patch: .semver.yaml
	@semver up release

minor: .semver.yaml
	@semver up minor

major: .semver.yaml
	@semver up major

release: .semver.yaml
	@semver up release

tag: .semver.yaml
	@git tag -a "v$(VERSION)" -m "Version $(VERSION)"

.PHONY: all build clean test test-coverage lint vet fmt init-semver alpha beta rc patch minor major release tag

all: clean fmt lint vet test build

build:
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	@rm -rf bin/ coverage.out coverage.html

test:
	@go test -v ./...

test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping linting"; \
	fi

vet:
	@go vet ./...

fmt:
	@go fmt ./...

run:
	@go run $(LDFLAGS) $(MAIN_PACKAGE)

$(KEYFILE):
	melange keygen --key-size $(KEYSIZE)

# OCI image targets
package: $(KEYFILE)
	@mkdir -p ./packages
	@VERSION=$(VERSION) melange build --arch $(ARCHS) \
		--signing-key $(KEYFILE) \
		--keyring-append $(KEYFILE).pub \
		--out-dir ./packages \
		--repository-append ./packages \
		melange.yaml --debug

image: package
	@apko build \
		--arch $(ARCHS) \
		--keyring-append $(KEYFILE).pub \
		--repository-append ./packages \
		apko.yaml \
		$(BINARY_NAME):$(VERSION) \
		image.tar

container: image
	@$(CONTAINER_RUNTIME) load < image.tar
	@$(CONTAINER_RUNTIME) run --rm -e MBTA_API_KEY -p 8080:8080 $(BINARY_NAME):$(VERSION)

keys:
	@if [ ! -f melange.rsa ]; then \
		openssl genrsa -out melange.rsa 4096; \
		openssl rsa -in melange.rsa -pubout -out melange.rsa.pub; \
	fi

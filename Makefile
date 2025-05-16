BINARY_NAME=mbta-mcp-server
VERSION=$(shell cat version.txt 2>/dev/null || echo "0.1.0")
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_VERSION=$(VERSION)+build.$(GIT_SHORT_SHA)
MAIN_PACKAGE=./cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS=-ldflags "-X main.Version=$(BUILD_VERSION)"

# Semver-cli should be available through devshell
# Fallback only if not in the development environment
semver-check:
	@command -v semver > /dev/null || (echo "semver not found, installing..." && \
	go install github.com/maykonlsf/semver-cli/cmd/semver@latest)

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

# Semver targets
.PHONY: patch minor major alpha beta rc release

patch: semver-check
	@echo "Bumping patch version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver bump patch $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

minor: semver-check
	@echo "Bumping minor version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver bump minor $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

major: semver-check
	@echo "Bumping major version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver bump major $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

alpha: semver-check
	@echo "Creating alpha version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver pre alpha $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

beta: semver-check
	@echo "Creating beta version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver pre beta $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

rc: semver-check
	@echo "Creating release candidate..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver pre rc $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

# Release targets
release: semver-check
	@echo "Creating release from pre-release..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver release $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION" && \
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

tag-version:
	@echo "Tagging version $(VERSION)..."
	@git add version.txt
	@git commit -m "chore: bump version to $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Version $(VERSION)"
	@echo "Tag v$(VERSION) created"
	@echo "Run 'git push && git push --tags' to push changes"
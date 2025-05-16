BINARY_NAME=mbta-mcp-server
VERSION=$(shell cat version.txt 2>/dev/null || echo "0.1.0")
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_VERSION=$(VERSION)+build.$(GIT_SHORT_SHA)
MAIN_PACKAGE=./cmd/server
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
LDFLAGS=-ldflags "-X main.Version=$(BUILD_VERSION)"

# Check if semver-cli is installed
SEMVER_CLI := $(shell command -v semver-cli 2> /dev/null)

# Install semver-cli if not installed
$(SEMVER_CLI):
	@echo "Installing semver-cli..."
	@npm install -g semver-cli

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

patch: $(SEMVER_CLI)
	@echo "Bumping patch version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver-cli inc patch $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

minor: $(SEMVER_CLI)
	@echo "Bumping minor version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver-cli inc minor $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

major: $(SEMVER_CLI)
	@echo "Bumping major version..."
	@CURR_VERSION=$$(cat version.txt) && \
	NEW_VERSION=$$(semver-cli inc major $$CURR_VERSION) && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

alpha: $(SEMVER_CLI)
	@echo "Creating alpha version..."
	@CURR_VERSION=$$(cat version.txt) && \
	BASE_VERSION=$$(semver-cli extract release-version $$CURR_VERSION) && \
	PRERELEASE=$$(semver-cli extract prerelease $$CURR_VERSION) && \
	if [ -z "$$PRERELEASE" ]; then \
		NEW_VERSION="$$BASE_VERSION-alpha.1"; \
	elif [[ "$$PRERELEASE" == alpha.* ]]; then \
		ALPHA_NUM=$$(echo $$PRERELEASE | sed 's/alpha\.//') && \
		NEW_VERSION="$$BASE_VERSION-alpha.$$((ALPHA_NUM+1))"; \
	else \
		NEW_VERSION="$$BASE_VERSION-alpha.1"; \
	fi && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

beta: $(SEMVER_CLI)
	@echo "Creating beta version..."
	@CURR_VERSION=$$(cat version.txt) && \
	BASE_VERSION=$$(semver-cli extract release-version $$CURR_VERSION) && \
	PRERELEASE=$$(semver-cli extract prerelease $$CURR_VERSION) && \
	if [ -z "$$PRERELEASE" ]; then \
		NEW_VERSION="$$BASE_VERSION-beta.1"; \
	elif [[ "$$PRERELEASE" == beta.* ]]; then \
		BETA_NUM=$$(echo $$PRERELEASE | sed 's/beta\.//') && \
		NEW_VERSION="$$BASE_VERSION-beta.$$((BETA_NUM+1))"; \
	else \
		NEW_VERSION="$$BASE_VERSION-beta.1"; \
	fi && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

rc: $(SEMVER_CLI)
	@echo "Creating release candidate..."
	@CURR_VERSION=$$(cat version.txt) && \
	BASE_VERSION=$$(semver-cli extract release-version $$CURR_VERSION) && \
	PRERELEASE=$$(semver-cli extract prerelease $$CURR_VERSION) && \
	if [ -z "$$PRERELEASE" ]; then \
		NEW_VERSION="$$BASE_VERSION-rc.1"; \
	elif [[ "$$PRERELEASE" == rc.* ]]; then \
		RC_NUM=$$(echo $$PRERELEASE | sed 's/rc\.//') && \
		NEW_VERSION="$$BASE_VERSION-rc.$$((RC_NUM+1))"; \
	else \
		NEW_VERSION="$$BASE_VERSION-rc.1"; \
	fi && \
	echo $$NEW_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

# Release targets
release: $(SEMVER_CLI)
	@echo "Creating release from pre-release..."
	@CURR_VERSION=$$(cat version.txt) && \
	BASE_VERSION=$$(semver-cli extract release-version $$CURR_VERSION) && \
	echo $$BASE_VERSION > version.txt && \
	echo "Version bumped from $$CURR_VERSION to $$BASE_VERSION" && \
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

tag-version:
	@echo "Tagging version $(VERSION)..."
	@git add version.txt
	@git commit -m "chore: bump version to $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Version $(VERSION)"
	@echo "Tag v$(VERSION) created"
	@echo "Run 'git push && git push --tags' to push changes"
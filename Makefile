BINARY_NAME=mbta-mcp-server
GIT_SHORT_SHA=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION=$(shell semver get release 2>/dev/null || echo "0.1.0")
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
.PHONY: init-semver patch minor major alpha beta rc release

# Initialize semver if not already initialized
init-semver: semver-check
	@if [ ! -f .semver.yaml ]; then \
		echo "Initializing semver with default version..."; \
		semver init; \
	fi

# Helper to update semver config file
define update-semver-yaml
	@CURR_VERSION=$$(semver get release) && \
	echo "Version before: $$CURR_VERSION" && \
	echo -e "alpha: 0\nbeta: 0\nrc: 0\nrelease: $(1)" > .semver.yaml && \
	echo "Version after: $(1)"
endef

# Patch version increment (implements custom logic since semver-cli doesn't support it directly)
patch: semver-check init-semver
	@echo "Bumping patch version..."
	@CURR_VERSION=$$(semver get release) && \
	$(eval PARTS := $$(subst ., ,$$CURR_VERSION)) \
	$(eval MAJOR := $$(word 1,$$(PARTS))) \
	$(eval MINOR := $$(word 2,$$(PARTS))) \
	$(eval PATCH := $$(word 3,$$(PARTS))) \
	$(eval NEW_PATCH := $$(shell expr $$(PATCH) + 1)) \
	$(eval NEW_VERSION := $$(MAJOR).$$(MINOR).$$(NEW_PATCH)) \
	$(call update-semver-yaml,$$(NEW_VERSION))

# Minor version increment (implements custom logic since semver-cli doesn't support it directly)
minor: semver-check init-semver
	@echo "Bumping minor version..."
	@CURR_VERSION=$$(semver get release) && \
	$(eval PARTS := $$(subst ., ,$$CURR_VERSION)) \
	$(eval MAJOR := $$(word 1,$$(PARTS))) \
	$(eval MINOR := $$(word 2,$$(PARTS))) \
	$(eval NEW_MINOR := $$(shell expr $$(MINOR) + 1)) \
	$(eval NEW_VERSION := $$(MAJOR).$$(NEW_MINOR).0) \
	$(call update-semver-yaml,$$(NEW_VERSION))

# Major version increment (implements custom logic since semver-cli doesn't support it directly)
major: semver-check init-semver
	@echo "Bumping major version..."
	@CURR_VERSION=$$(semver get release) && \
	$(eval PARTS := $$(subst ., ,$$CURR_VERSION)) \
	$(eval MAJOR := $$(word 1,$$(PARTS))) \
	$(eval NEW_MAJOR := $$(shell expr $$(MAJOR) + 1)) \
	$(eval NEW_VERSION := $$(NEW_MAJOR).0.0) \
	$(call update-semver-yaml,$$(NEW_VERSION))

# Alpha version (using built-in semver-cli functionality)
alpha: semver-check init-semver
	@echo "Creating alpha version..."
	@CURR_VERSION=$$(semver get release) && \
	semver up alpha && \
	NEW_VERSION=$$(semver get alpha) && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

# Beta version (using built-in semver-cli functionality)
beta: semver-check init-semver
	@echo "Creating beta version..."
	@CURR_VERSION=$$(semver get release) && \
	semver up beta && \
	NEW_VERSION=$$(semver get beta) && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

# Release candidate (using built-in semver-cli functionality)
rc: semver-check init-semver
	@echo "Creating release candidate..."
	@CURR_VERSION=$$(semver get release) && \
	semver up rc && \
	NEW_VERSION=$$(semver get rc) && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION"

# Final release (using built-in semver-cli functionality)
release: semver-check init-semver
	@echo "Creating final release..."
	@CURR_VERSION=$$(semver get alpha 2>/dev/null || semver get beta 2>/dev/null || semver get rc 2>/dev/null || semver get release) && \
	semver up release && \
	NEW_VERSION=$$(semver get release) && \
	echo "Version bumped from $$CURR_VERSION to $$NEW_VERSION" && \
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

tag-version: init-semver
	@echo "Tagging version $(VERSION)..."
	@git add .semver.yaml
	@git commit -m "chore: bump version to $(VERSION)"
	@git tag -a "v$(VERSION)" -m "Version $(VERSION)"
	@echo "Tag v$(VERSION) created"
	@echo "Run 'git push && git push --tags' to push changes"
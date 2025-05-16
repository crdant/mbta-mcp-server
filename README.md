# MBTA MCP Server

An MCP server that communicates with the MBTA API to provide Boston-area transit information.

This Machine Learning Control Protocol (MCP) server integrates with the Massachusetts Bay Transportation Authority (MBTA) API to provide real-time and scheduled transit information for the Boston area. It enables AI assistants to access MBTA data through a standardized interface.

## Features

- Real-time transit predictions
- Service alerts and disruptions
- Route and schedule information
- Accessibility information
- Trip planning assistance
- Location-based station finding

## Installation

### Docker

```bash
docker pull ghcr.io/crdant/mbta-mcp-server:latest
docker run -e MBTA_API_KEY="your-api-key" ghcr.io/crdant/mbta-mcp-server:latest
```

### Go Installation

```bash
go install github.com/username/mbta-mcp-server@latest
```

## Configuration

Set your MBTA API key in the environment:

```bash
export MBTA_API_KEY="your-api-key"
```

## Usage

The server implements the MCP stdio protocol for local usage with AI assistants.

For more detailed information, see the [specification](spec.md).

## Supply Chain Security

### Container Image Signing

All container images are signed using Sigstore's Cosign with keyless signing. This allows users to verify that the container image was built by our GitHub Actions CI/CD pipeline.

#### Signing Security Practice

We follow the best practice for container image signing:

**We sign only the image digest (content hash)** - This is the most secure approach since the digest is a unique, immutable identifier for the specific content. By signing only the digest, we avoid any potential security issues that could arise from mutable tags like `latest`.

#### Verifying Container Images

To verify our container images, always verify by digest:

```bash
# Get the digest first (using any tag to lookup the image)
DIGEST=$(crane digest ghcr.io/crdant/mbta-mcp-server:1.2.3)

# Verify the image by digest
cosign verify \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/crdant/mbta-mcp-server@$DIGEST
```

### Software Bill of Materials (SBOM)

Each build generates a comprehensive Software Bill of Materials (SBOM) that lists all components included in the container image. The SBOM is:

1. Generated during the build process
2. Signed with a GitHub-issued certificate using the actions/attest-sbom tool
3. Available as a GitHub Actions artifact with each build
4. Attached to the container image as an attestation by digest

To verify the SBOM attestation:

```bash
# Get the digest first (most reliable approach)
DIGEST=$(crane digest ghcr.io/crdant/mbta-mcp-server:1.2.3)

# Verify the SBOM attestation by digest
cosign verify-attestation \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --type spdx \
  ghcr.io/crdant/mbta-mcp-server@$DIGEST
```

### Vulnerability Scanning

We use Trivy to scan our container images for vulnerabilities:

1. Container images are automatically scanned after they're built
2. Results are uploaded to GitHub Security in SARIF format
3. Critical and High severity vulnerabilities are reported
4. Scans focus on vulnerabilities with available fixes

These security measures help ensure our software supply chain is secure and transparent from source code to container deployment.

## License

[MIT License](LICENSE)
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
docker pull username/mbta-mcp-server
docker run -p 8080:8080 username/mbta-mcp-server
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

#### Signing Security Practices

We follow best practices for container image signing:

1. **We sign the image digest (content hash)** - This is the most secure approach since the digest is a unique, immutable identifier for the specific content.

2. **We sign immutable tags only** - We only sign tags that are immutable (e.g., version tags like `v1.2.3` and specific SHA commit tags).

3. **We do not sign mutable tags** - Tags that can move (like `latest` or major version tags like `v1`) are not signed because this could create misleading security assertions.

#### Verifying Container Images

For highest security, verify by digest (this reference can't change):

```bash
# Get the digest first
DIGEST=$(crane digest ghcr.io/crdant/mbta-mcp-server:v1.2.3)

# Verify the image by digest
cosign verify \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/crdant/mbta-mcp-server@$DIGEST
```

For version-specific tags, you can verify directly:

```bash
# Verify a specific version tag
cosign verify \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/crdant/mbta-mcp-server:v1.2.3
```

### Software Bill of Materials (SBOM)

Each build generates a comprehensive Software Bill of Materials (SBOM) that lists all components included in the container image. The SBOM is:

1. Generated during the build process
2. Signed with a GitHub-issued certificate using the actions/attest-sbom tool
3. Available as a GitHub Actions artifact with each build
4. Attached to the container image as an attestation

To verify the SBOM attestation:

```bash
# Get the digest first (most reliable approach)
DIGEST=$(crane digest ghcr.io/crdant/mbta-mcp-server:v1.2.3)

# Verify the SBOM attestation by digest
cosign verify-attestation \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --type cyclonedx \
  ghcr.io/crdant/mbta-mcp-server@$DIGEST
```

You can also verify SBOM attestations for specific immutable tags:

```bash
cosign verify-attestation \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --type cyclonedx \
  ghcr.io/crdant/mbta-mcp-server:v1.2.3
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
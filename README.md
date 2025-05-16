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

To verify a container image:

```bash
# Install cosign
cosign verify \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/crdant/mbta-mcp-server:latest
```

### Software Bill of Materials (SBOM)

Each build generates a comprehensive Software Bill of Materials (SBOM) that lists all components included in the container image. The SBOM is:

1. Generated during the build process
2. Signed with a GitHub-issued certificate using the actions/attest-sbom tool
3. Available as a GitHub Actions artifact with each build
4. Attached to the container image as an attestation

To verify the SBOM attestation:

```bash
cosign verify-attestation \
  --certificate-identity "https://github.com/crdant/mbta-mcp-server/.github/workflows/build.yml@refs/heads/main" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --type cyclonedx \
  ghcr.io/crdant/mbta-mcp-server:latest
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
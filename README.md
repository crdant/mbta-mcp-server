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

## License

[MIT License](LICENSE)
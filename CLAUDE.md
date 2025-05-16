# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Build and Run

```bash
# Build the application
make build

# Run the application
make run

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests for a specific package
go test ./internal/config -v

# Run a specific test
go test ./internal/config -v -run TestNew
```

### Code Quality

```bash
# Format code
make fmt

# Run linter (if golangci-lint is installed)
make lint

# Run Go vet
make vet

# Run all code quality checks and tests
make all
```

### Docker

```bash
# Build Docker image
make image

# Run in Docker container
make container
```

## Architecture Overview

The MBTA MCP Server is a Machine Learning Control Protocol (MCP) server that integrates with the Massachusetts Bay Transportation Authority (MBTA) API to provide Boston-area transit information to AI assistants.

### Key Components

1. **MCP Server**: Implemented using the mcp-go library, it handles the MCP protocol and provides an interface for AI assistants to query transit information.

2. **MBTA API Client**: Connects to the MBTA API v3, handles authentication, rate limiting, and error handling.

3. **Configuration System**: Environment-based configuration system that manages settings like API keys, timeouts, and logging levels.

4. **Data Models**: Representations of MBTA transit data like routes, stops, schedules, and alerts.

5. **Request/Response Handlers**: Transform MCP requests into MBTA API calls and format responses back to MCP protocol.

### Project Structure

- `cmd/server/`: Main application entry point
- `internal/`: Private application code
  - `config/`: Configuration loading and management
  - `testutil/`: Test utilities and helpers
  - `server/`: MCP server implementation (planned)
  - `handlers/`: Request handlers (planned)
- `pkg/`: Public packages that may be used by external applications
  - `mbta/`: MBTA API client (planned)
- `test/`: Test fixtures and utilities

### Configuration

The application is configured using environment variables:

- `MBTA_API_KEY`: API key for the MBTA API
- `DEBUG`: Enable debug mode (true/false)
- `LOG_LEVEL`: Logging level (info, debug, error)
- `TIMEOUT_SECONDS`: API request timeout in seconds
- `MBTA_API_URL`: Base URL for the MBTA API
- `PORT`: Server port
- `ENVIRONMENT`: Deployment environment (development, production)

### Implementation Plan

The project follows a phased implementation approach:

1. Project setup and core structure (completed)
2. MBTA API client development
3. Core MCP protocol implementation
4. Transit information features
5. Enhanced features (trip planning, alerts)
6. Deployment and documentation

## Development Guidelines

- Follow test-driven development practices by writing tests before implementation
- Use Go idiomatic patterns and best practices
- Document all exported functions and types
- Keep the MCP server and MBTA API client decoupled for better testability
- Use mocks for external dependencies in tests
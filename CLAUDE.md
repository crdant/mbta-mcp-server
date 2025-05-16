# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Writing code

- We prefer simple, clean, maintainable solutions over clever or complex ones, even if the latter are more concise or performant. Readability and maintainability are primary concerns.
- Make the smallest reasonable changes to get to the desired outcome. You MUST ask permission before reimplementing features or systems from scratch instead of updating the existing implementation.
- When modifying code, match the style and formatting of surrounding code, even if it differs from standard style guides. Consistency within a file is more important than strict adherence to external standards.
- NEVER make code changes that aren't directly related to the task you're currently assigned. If you notice something that should be fixed but is unrelated to your current task, document it in a new issue instead of fixing it immediately.
- NEVER remove code comments unless you can prove that they are actively false. Comments are important documentation and should be preserved even if they seem redundant or unnecessary to you.
- All code files should start with a brief 2 line comment explaining what the file does. Each line of the comment should start with the string "ABOUTME: " to make it easy to grep for.
- When writing comments, avoid referring to temporal context about refactors or recent changes. Comments should be evergreen and describe the code as it is, not how it evolved or was recently changed.
- NEVER implement a mock mode for testing or for any purpose. We always use real data and real APIs, never mock implementations.
- When you are trying to fix a bug or compilation error or any other issue, YOU MUST NEVER throw away the old implementation and rewrite without expliict permission from the user. If you are going to do this, YOU MUST STOP and get explicit permission from the user.
- NEVER name things as 'improved' or 'new' or 'enhanced', etc. Code naming should be evergreen. What is new today will be "old" someday.

## Getting help

- ALWAYS ask for clarification rather than making assumptions.
- If you're having trouble with something, it's ok to stop and ask for help. Especially if it's something your human might be better at.

## Development Guidelines

- Use Go idiomatic patterns and best practices
- Document all exported functions and types
- Keep the MCP server and MBTA API client decoupled for better testability
- Use mocks for external dependencies in tests

## Testing

- Tests MUST cover the functionality being implemented.
- NEVER ignore the output of the system or the tests - Logs and messages often contain CRITICAL information.
- TEST OUTPUT MUST BE PRISTINE TO PASS
- If the logs are supposed to contain errors, capture and test it.
- NO EXCEPTIONS POLICY: Under no circumstances should you mark any test type as "not applicable". Every project, regardless of size or complexity, MUST have unit tests, integration tests, AND end-to-end tests. If you believe a test type doesn't apply, you need the human to say exactly "I AUTHORIZE YOU TO SKIP WRITING TESTS THIS TIME"

### We practice TDD. That means:

- Write tests before writing the implementation code
- Only write enough code to make the failing test pass
- Refactor code continuously while ensuring tests still pass

#### TDD Implementation Process

- Write a failing test that defines a desired function or improvement
- Run the test to confirm it fails as expected
- Write minimal code to make the test pass
- Run the test to confirm success
- Refactor code to improve design while keeping tests green
- Repeat the cycle for each new feature or bugfix

## Pull Request Guidelines

When creating pull requests, follow these guidelines:

### PR Format
- Titles should be concise (40 characters or less) and use present tense with implied subject
- Body must include two main sections:
  - **TL;DR**: 1-2 line summary of the change
  - **Details**: Paragraph(s) explaining intent and impact (not just listing files changed)
- Use appropriate section headers with dashes underlining each header
- Always use present tense with the PR as an implied subject
- Focus on the "why" and impact of changes, not just the "what"
- Never use phrases like "this PR" or "this change" (the PR is the implied subject)

### PR Process
- Keep PRs focused on a single logical change
- Ensure all tests pass before requesting review
- Address feedback promptly and completely
- Update documentation affected by your changes
- Link to relevant issues with "Fixes #123" or "Relates to #123"

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


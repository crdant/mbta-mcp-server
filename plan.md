# MBTA MCP Server Implementation Plan

This document outlines a step-by-step approach to building the MBTA MCP server using test-driven development practices. Each step builds incrementally on previous steps, ensuring manageable complexity and thorough testing coverage.

## Implementation Phases

### Phase 1: Project Setup and Core Structure
- Initialize Go project with proper module structure
- Set up testing framework and CI configuration
- Implement basic MCP server skeleton

### Phase 2: MBTA API Client
- Create client for MBTA API interactions
- Implement data models matching API responses
- Build authentication and error handling

### Phase 3: Core MCP Implementation
- Implement basic MCP protocol handlers
- Develop request/response formatting
- Connect API client to MCP handlers

### Phase 4: Transit Information Features
- Implement routes and schedules functionality
- Add stops and stations information
- Create vehicle tracking capabilities

### Phase 5: Enhanced Features
- Build trip planning functionality
- Implement service alerts and disruptions
- Add accessibility information queries

### Phase 6: Deployment and Documentation
- Containerize application
- Create comprehensive documentation
- Build CI/CD pipeline

## Detailed Implementation Steps

### Phase 1: Project Setup and Core Structure

#### Step 1.1: Initialize Go Project Structure
- Set up Go module
- Create directory structure
- Add license and initial README

#### Step 1.2: Set Up Testing Framework
- Configure Go testing
- Add test helpers and mocks
- Implement first basic tests

#### Step 1.3: MCP Server Skeleton
- Add basic MCP implementation
- Create server entry point
- Implement configuration loading

### Phase 2: MBTA API Client

#### Step 2.1: API Client Foundation
- Create HTTP client with proper timeout handling
- Implement API authentication
- Add basic error handling

#### Step 2.2: Core Data Models
- Define route models
- Create stop/station models
- Implement schedule data structures

#### Step 2.3: Client Integration Tests
- Add mocked API responses
- Implement integration tests for API client
- Test error handling and rate limiting

### Phase 3: Core MCP Implementation

#### Step 3.1: MCP Request Handling
- Implement request parsing
- Add response formatting
- Create main handler logic

#### Step 3.2: MCP-MBTA Integration
- Connect MCP handlers to MBTA client
- Implement request transformation
- Add response mapping

#### Step 3.3: Error Handling
- Implement MCP-compliant error responses
- Add logging for errors
- Create retry mechanisms

### Phase 4: Transit Information Features

#### Step 4.1: Routes and Schedules
- Implement route listing endpoint
- Add schedule retrieval
- Create route filtering options

#### Step 4.2: Stops and Stations
- Add station information endpoint
- Implement station search
- Create accessibility information

#### Step 4.3: Real-time Updates
- Implement vehicle tracking
- Add arrival predictions
- Create status update endpoint

### Phase 5: Enhanced Features

#### Step 5.1: Trip Planning
- Implement trip planning algorithm
- Add transfer point identification
- Create travel time estimation

#### Step 5.2: Service Alerts
- Add service disruption detection
- Implement alternative route suggestions
- Create alert notification endpoint

#### Step 5.3: Geographic Queries
- Implement station proximity search
- Add route coverage mapping
- Create location-based recommendations

### Phase 6: Deployment and Documentation

#### Step 6.1: Containerization
- Create optimized Docker build
- Implement Docker Compose setup
- Add container health checks

#### Step 6.2: Documentation
- Create API documentation
- Add usage examples
- Write deployment guides

#### Step 6.3: CI/CD Pipeline
- Set up automated testing
- Implement build pipeline
- Add release automation

## TDD Implementation Prompts

Below are detailed prompts for implementing each step following test-driven development principles.

### Prompt 1: Project Setup and Go Module Initialization
```
Implement the initial project setup for an MBTA MCP server written in Go. Follow these specific steps:

1. Initialize a Go module named "github.com/crdant/mbta-mcp-server"
2. Create the basic directory structure following Go best practices:
   - cmd/server/ - for the main application entry point
   - internal/ - for private application code
   - pkg/ - for code that may be used by external applications
   - test/ - for test utilities and fixtures
3. Create a basic main.go file in cmd/server that only logs startup messages
4. Implement a simple configuration loader that reads environment variables including MBTA_API_KEY
5. Add a Makefile with targets for build, test, and lint
6. Create initial Go test files for the configuration loading functionality

Follow test-driven development by writing tests first, then implementing the code to make those tests pass. Ensure proper error handling and logging are included.
```

### Prompt 2: Testing Framework and Utilities
```
Set up a comprehensive testing framework for the MBTA MCP server. Implement the following test utilities:

1. Create a test helper package in the internal/testutil directory with:
   - Mock HTTP server for simulating MBTA API responses
   - Test fixtures for common API responses
   - Utility functions for test setup and teardown

2. Implement the following test fixtures:
   - Sample MBTA API responses for routes, stops, schedules
   - Error response examples
   - Authentication test cases

3. Create a test configuration file for development and testing

4. Write an example test that demonstrates how to use the test fixtures to test API client functionality

Follow test-driven development practices by ensuring all test utilities have their own tests. Use Go's testing package and follow Go best practices for writing testable code.
```

### Prompt 3: Basic MCP Server Implementation
```
Implement the core MCP server structure using the mcp-go library. Follow these steps:

1. Add the mcp-go module as a dependency
2. Create an internal/server package with:
   - Server initialization function
   - Request handling setup
   - Basic logging middleware

3. Implement the stdio protocol handling as specified in the MCP documentation
4. Add configuration options for:
   - Log level
   - Debug mode
   - Timeout settings

5. Create unit tests for all server functionality:
   - Test server initialization
   - Test request handling with mock requests
   - Test configuration loading

6. Update the main.go file to use the new server package

Use test-driven development by writing tests for each component before implementing it. Ensure proper error handling and context cancellation. The server should not yet connect to the MBTA API.
```

### Prompt 4: MBTA API Client Foundation
```
Implement a client for the MBTA API v3. Follow these steps:

1. Create a pkg/mbta package with:
   - Client struct with necessary configuration
   - Authentication handling for API key
   - Base request/response handling

2. Implement the following functionality:
   - Configurable HTTP client with timeouts
   - Rate limiting handling
   - Error parsing and standardization

3. Create interfaces for all external dependencies to allow for mocking in tests

4. Write comprehensive tests using the test utilities created earlier:
   - Test successful API connection
   - Test authentication errors
   - Test rate limiting behavior
   - Test timeout handling

5. Implement graceful handling of API unavailability

Use test-driven development by writing tests for each API endpoint before implementing it. Mock the HTTP responses to avoid actual API calls during tests.
```

### Prompt 5: MBTA Data Models
```
Implement the core data models for the MBTA API responses. Follow these steps:

1. Create model structs in pkg/mbta/models for:
   - Routes (with types, names, directions)
   - Stops (with locations, accessibility information)
   - Trips (with schedule information)
   - Vehicles (with locations, status)
   - Alerts (with effect, description)

2. Add JSON tags for proper serialization/deserialization

3. Implement validation methods for each model

4. Create mapping functions that convert API responses to simplified models
   suitable for MCP responses

5. Write comprehensive tests for:
   - Model validation
   - JSON serialization/deserialization
   - Mapping functions

Follow test-driven development by writing tests using sample responses from the MBTA API documentation. Ensure proper error handling for malformed data.
```

### Prompt 6: MBTA API Client Implementation
```
Expand the MBTA API client with specific endpoint implementations. Follow these steps:

1. Add the following endpoint methods to the client:
   - GetRoutes() - retrieve available routes
   - GetStops() - retrieve stops/stations
   - GetSchedule() - retrieve schedules for routes
   - GetVehicles() - retrieve vehicle locations
   - GetAlerts() - retrieve service alerts

2. Implement parameter handling for each endpoint:
   - Filtering options
   - Pagination support
   - Include parameters for related data

3. Add response parsing for each endpoint that converts JSON to model structs

4. Create integration tests that verify the correct handling of:
   - Successful responses
   - Error responses
   - Malformed responses
   - Authentication issues

5. Implement caching headers support

Follow test-driven development by writing tests for each endpoint before implementing it. Use the test fixtures created earlier for mock responses.
```

### Prompt 7: MCP Request Handler Implementation
```
Implement the core MCP request handlers. Follow these steps:

1. Create an internal/handlers package with:
   - Base handler interface
   - Request parsing functions
   - Response formatting functions

2. Implement handlers for:
   - Capabilities listing (what the MCP server can do)
   - Transit information requests
   - Error responses

3. Add request validation and parameter extraction

4. Create response formatting that follows MCP best practices

5. Implement context handling for cancellation and timeouts

6. Write comprehensive tests for:
   - Request parsing
   - Response formatting
   - Error handling
   - Context cancellation

Follow test-driven development by writing tests for each handler before implementing it. Ensure the handlers follow the MCP protocol specification correctly.
```

### Prompt 8: Connect MCP Handlers to MBTA Client
```
Connect the MCP request handlers to the MBTA API client. Follow these steps:

1. Create an internal/service package that:
   - Initializes both the MCP server and MBTA client
   - Wires request handlers to appropriate client methods
   - Handles error translation between systems

2. Implement request transformation from MCP format to MBTA API parameters

3. Create response mapping from MBTA models to MCP response format

4. Add proper error handling that follows MCP best practices:
   - Translate API errors to appropriate MCP errors
   - Include proper error codes and messages
   - Handle timeouts and service unavailability

5. Write integration tests that verify:
   - End-to-end request handling
   - Error propagation
   - Timeout handling

Follow test-driven development by writing tests for each integration point before implementing it. Use mocks for the MBTA client to test the service layer independently.
```

### Prompt 9: Routes and Schedules Implementation
```
Implement the routes and schedules functionality. Follow these steps:

1. Enhance the MBTA client with detailed methods for:
   - Get route details with all attributes
   - Search routes by type, name, or destination
   - Get schedules with stop times

2. Create MCP handlers for:
   - List all available routes
   - Get route details with schedule information
   - Search routes by criteria

3. Implement response formatting that simplifies MBTA data for MCP consumers:
   - Clear route information with relevant attributes
   - Human-readable schedule information
   - Properly formatted timestamps

4. Write comprehensive tests for:
   - Route listing functionality
   - Schedule retrieval
   - Search functionality
   - Edge cases (no routes, service not running)

Follow test-driven development by writing tests for each feature before implementing it. Use real MBTA API response samples for testing.
```

### Prompt 10: Stops and Stations Implementation
```
Implement the stops and stations functionality. Follow these steps:

1. Enhance the MBTA client with detailed methods for:
   - Get all stops with location data
   - Get stop details with accessibility information
   - Search stops by name or proximity

2. Create MCP handlers for:
   - List all stops/stations
   - Get detailed stop information
   - Search stops by name
   - Find stops near a location

3. Add accessibility information:
   - Elevator status
   - Accessible entrance information
   - Accessibility features

4. Write comprehensive tests for:
   - Stop listing functionality
   - Accessibility information retrieval
   - Search functionality
   - Geographic proximity search
   - Edge cases (invalid locations, no results)

Follow test-driven development by writing tests for each feature before implementing it. Ensure proper handling of geographic data and location queries.
```

### Prompt 11: Real-time Updates Implementation
```
Implement real-time update functionality. Follow these steps:

1. Enhance the MBTA client with methods for:
   - Get vehicle locations
   - Get prediction information
   - Get service alerts

2. Create MCP handlers for:
   - Get real-time vehicle positions
   - Get arrival/departure predictions
   - Get service status updates

3. Implement real-time data formatting:
   - Clear presentation of predictions
   - Vehicle information with status
   - Human-readable alerts

4. Add support for streaming updates (if supported by MCP)

5. Write comprehensive tests for:
   - Vehicle tracking functionality
   - Prediction accuracy
   - Alert retrieval and formatting
   - Edge cases (delayed service, outages)

Follow test-driven development by writing tests for each feature before implementing it. Consider caching strategies for real-time data to reduce API load.
```

### Prompt 12: Trip Planning Implementation
```
Implement trip planning functionality. Follow these steps:

1. Create a pkg/tripplanner package with:
   - Trip planning algorithm
   - Transfer point identification
   - Travel time estimation

2. Enhance the MBTA client with methods for:
   - Get trip options between stops
   - Get transfer points
   - Get travel time estimates

3. Create MCP handlers for:
   - Plan trips between locations
   - Get trip alternatives
   - Get estimated travel times

4. Implement response formatting that presents:
   - Clear step-by-step instructions
   - Transfer information
   - Timing for each segment
   - Accessibility considerations

5. Write comprehensive tests for:
   - Trip planning functionality
   - Multiple route options
   - Transfer suggestions
   - Edge cases (no route, service disruptions)

Follow test-driven development by writing tests for each feature before implementing it. Consider performance optimizations for complex trip planning queries.
```

### Prompt 13: Service Alerts Implementation
```
Implement service alerts and disruption handling. Follow these steps:

1. Enhance the MBTA client with detailed alert methods:
   - Get all active alerts
   - Get alerts filtered by route/stop
   - Get alerts with severity information

2. Create MCP handlers for:
   - Get all service alerts
   - Get alerts for specific routes or stations
   - Get service alternatives during disruptions

3. Implement alert analysis logic:
   - Determine impact on trips
   - Suggest alternative routes
   - Estimate delay times

4. Write comprehensive tests for:
   - Alert retrieval
   - Alert filtering
   - Alternative suggestion
   - Edge cases (major outages, partial service)

Follow test-driven development by writing tests for each feature before implementing it. Ensure alerts are properly categorized by severity and impact.
```

### Prompt 14: Geographic Queries Implementation
```
Implement geographic query functionality. Follow these steps:

1. Create a pkg/geo package with:
   - Geospatial utilities
   - Distance calculation
   - Area coverage determination

2. Enhance the MBTA client with geographic methods:
   - Find stops near coordinates
   - Find routes serving a geographic area
   - Get coverage information

3. Create MCP handlers for:
   - Find nearest stations
   - Get routes serving a location
   - Get coverage information for an area

4. Implement response formatting that presents:
   - Distances to stations
   - Coverage maps (as available)
   - Route proximity information

5. Write comprehensive tests for:
   - Proximity search
   - Coverage determination
   - Edge cases (out of service area, invalid coordinates)

Follow test-driven development by writing tests for each feature before implementing it. Consider performance optimizations for geographic calculations.
```

### Prompt 15: Docker Containerization
```
Implement Docker containerization for the application. Follow these steps:

1. Create a minimal Dockerfile that:
   - Uses a multi-stage build approach
   - Results in a small, secure image
   - Includes only the necessary runtime components

2. Create a docker-compose.yml file for easy deployment with:
   - Environment variable configuration
   - Volume mounting for persistence (if needed)
   - Health checks

3. Implement container entrypoint script that:
   - Validates configuration
   - Sets up proper permissions
   - Handles signals correctly

4. Write tests to verify:
   - Container builds successfully
   - Application runs correctly in container
   - Configuration is properly loaded

Follow best practices for container security and minimizing image size. Include documentation on how to build and run the container.
```

### Prompt 16: Comprehensive Documentation
```
Create comprehensive documentation for the project. Follow these steps:

1. Update the README.md with:
   - Detailed project overview
   - Installation instructions
   - Configuration options
   - Basic usage examples

2. Create an API.md document that describes:
   - All available MCP endpoints
   - Request/response formats
   - Error codes and meanings
   - Rate limiting information

3. Add a DEVELOPMENT.md guide that covers:
   - Setting up a development environment
   - Running tests
   - Contributing guidelines
   - Code style recommendations

4. Create a DEPLOYMENT.md document with:
   - Production deployment instructions
   - Security considerations
   - Performance tuning
   - Monitoring recommendations

5. Add inline code documentation for all exported functions and types

Ensure documentation follows a consistent style and provides enough information for both users and contributors to understand the project.
```

### Prompt 17: CI/CD Pipeline Setup
```
Set up a CI/CD pipeline for the project. Follow these steps:

1. Create a GitHub Actions workflow that:
   - Runs tests on pull requests
   - Builds and verifies Docker images
   - Checks code style and linting
   - Runs security scanning

2. Implement versioning with:
   - Semantic version tagging
   - Automated changelog generation
   - Version embedding in binaries

3. Add release automation that:
   - Builds binaries for multiple platforms
   - Creates Docker images with version tags
   - Generates release notes
   - Publishes artifacts

4. Create quality gates that ensure:
   - Test coverage requirements
   - Performance benchmarks
   - Security scanning passes

Follow best practices for CI/CD security, including proper secret management and minimizing build privileges.
```

## Implementation Order and Dependencies

The implementation should proceed in this order:

1. Project Setup (Prompt 1)
2. Testing Framework (Prompt 2)
3. Basic MCP Server (Prompt 3)
4. MBTA API Client Foundation (Prompt 4)
5. MBTA Data Models (Prompt 5)
6. MBTA API Client Implementation (Prompt 6)
7. MCP Request Handler Implementation (Prompt 7)
8. Connect MCP Handlers to MBTA Client (Prompt 8)
9. Routes and Schedules Implementation (Prompt 9)
10. Stops and Stations Implementation (Prompt 10)
11. Real-time Updates Implementation (Prompt 11)
12. Trip Planning Implementation (Prompt 12)
13. Service Alerts Implementation (Prompt 13)
14. Geographic Queries Implementation (Prompt 14)
15. Docker Containerization (Prompt 15)
16. Comprehensive Documentation (Prompt 16)
17. CI/CD Pipeline Setup (Prompt 17)

This order ensures that each step builds upon the previous ones in a logical sequence, starting with core infrastructure and gradually adding features.

## Iterative Development and Testing

For each implementation prompt:

1. Begin by writing tests that define the expected behavior
2. Implement the minimum code necessary to make tests pass
3. Refactor for clarity, performance, and maintainability
4. Verify all tests still pass after refactoring
5. Document the implemented functionality
6. Integrate with previously implemented components

Following this TDD approach ensures high-quality code with good test coverage and minimal technical debt.

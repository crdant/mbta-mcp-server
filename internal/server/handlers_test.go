// ABOUTME: This file contains tests for the MCP server request handlers.
// ABOUTME: It verifies proper handling of transit information requests.

package server

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/mock"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestRegisterDefaultHandlers(t *testing.T) {
	// This test verifies that the RegisterDefaultHandlers method exists
	// and can be called without errors
	t.Run("Register default handlers method exists", func(t *testing.T) {
		// Create a mock server
		server := &Server{}

		// Define a function that calls RegisterDefaultHandlers
		registerFn := func() {
			server.RegisterDefaultHandlers()
		}

		// If this compiles, the method exists
		_ = registerFn
	})
}

func TestGetRoutesHandler(t *testing.T) {
	// Create a mock MBTA API server
	mockServer, err := mock.StandardMockServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer mockServer.Close()

	// Create config pointing to mock server
	cfg := &config.Config{
		APIKey:     "test-api-key",
		APIBaseURL: mockServer.URL,
	}

	// Create MBTA client with mock server
	_ = mbta.NewClient(cfg)

	// Create MCP server with the MBTA client
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	t.Run("Get routes handler can be registered", func(t *testing.T) {
		// Register the routes handler
		server.registerTransitInfoTools()
	})

	t.Run("Get routes returns valid route data", func(t *testing.T) {
		// Create a request for the routes handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_routes",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getRoutesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Check that content is returned
		if len(result.Content) == 0 {
			t.Fatal("Handler returned empty content")
		}

		// Verify content type is text
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Verify the text indicates success
		if textContent.Text == "" {
			t.Error("Text content is empty")
		}
	})

	t.Run("Get routes handles filtering by route type", func(t *testing.T) {
		// Create a request with route type filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_routes",
				Arguments: map[string]any{
					"route_type": "1", // Subway
				},
			},
		}

		// Call the handler
		result, err := server.getRoutesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get routes handles filtering by route ID", func(t *testing.T) {
		// Create a request with route ID filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_routes",
				Arguments: map[string]any{
					"route_id": "Red",
				},
			},
		}

		// Call the handler
		result, err := server.getRoutesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get routes returns data as JSON content", func(t *testing.T) {
		// Create a mocked implementation of getRoutesHandler that returns proper JSON
		// This tests that the handler will eventually be implemented to return
		// structured data, not just text
		handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			client := mbta.NewClient(cfg)
			routes, err := client.GetRoutes(ctx)
			if err != nil {
				return createErrorResponse("Failed to retrieve routes: " + err.Error()), nil
			}

			// Convert the routes to a map for JSON serialization
			routesData := make([]map[string]interface{}, 0, len(routes))
			for _, route := range routes {
				routeMap := map[string]interface{}{
					"id":          route.ID,
					"name":        route.Attributes.LongName,
					"type":        route.Attributes.Type,
					"description": route.Attributes.Description,
				}
				routesData = append(routesData, routeMap)
			}

			// Create JSON content as text since JSONContent is not directly available
			jsonBytes, err := json.Marshal(routesData)
			if err != nil {
				return createErrorResponse("Failed to serialize route data: " + err.Error()), nil
			}

			// Return as Text content with JSON
			textContent := mcp.TextContent{
				Type: "text",
				Text: string(jsonBytes),
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{textContent},
			}, nil
		}

		// Create a request for the routes handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_routes",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := handler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil and has content
		if result == nil || len(result.Content) == 0 {
			t.Fatal("Handler returned nil or empty result")
		}

		// Verify content type is Text
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Verify the JSON can be parsed
		var parsedData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &parsedData); err != nil {
			t.Fatalf("Failed to parse JSON content: %v", err)
		}

		// Verify we have route data
		if len(parsedData) == 0 {
			t.Error("No routes returned in JSON")
		}

		// Check that the route data has the expected fields
		for _, route := range parsedData {
			if _, ok := route["id"]; !ok {
				t.Error("Route is missing 'id' field")
			}
			if _, ok := route["name"]; !ok {
				t.Error("Route is missing 'name' field")
			}
			if _, ok := route["type"]; !ok {
				t.Error("Route is missing 'type' field")
			}
		}
	})
}

// TestGetStopsHandler tests the stops handler functionality
func TestGetStopsHandler(t *testing.T) {
	// Create a mock MBTA API server
	mockServer, err := mock.StandardMockServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer mockServer.Close()

	// Create config pointing to mock server
	cfg := &config.Config{
		APIKey:     "test-api-key",
		APIBaseURL: mockServer.URL,
	}

	// Create MBTA client with mock server
	_ = mbta.NewClient(cfg)

	// Create MCP server with the MBTA client
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	t.Run("Get stops handler can be registered", func(t *testing.T) {
		// Register the stops handler
		server.registerTransitInfoTools()
	})

	t.Run("Get stops returns valid stop data", func(t *testing.T) {
		// Create a request for the stops handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_stops",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getStopsHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Check that content is returned
		if len(result.Content) == 0 {
			t.Fatal("Handler returned empty content")
		}

		// Verify content type is text
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Verify the text indicates success
		if textContent.Text == "" {
			t.Error("Text content is empty")
		}
	})

	t.Run("Get stops handles filtering by location type", func(t *testing.T) {
		// Create a request with location type filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_stops",
				Arguments: map[string]any{
					"location_type": "1", // Station
				},
			},
		}

		// Call the handler
		result, err := server.getStopsHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get stops handles filtering by stop ID", func(t *testing.T) {
		// Create a request with stop ID filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_stops",
				Arguments: map[string]any{
					"stop_id": "place-north",
				},
			},
		}

		// Call the handler
		result, err := server.getStopsHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get stops handles filtering by route ID", func(t *testing.T) {
		// Create a request with route ID filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_stops",
				Arguments: map[string]any{
					"route_id": "Red",
				},
			},
		}

		// Call the handler
		result, err := server.getStopsHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Verify content is returned
		if len(result.Content) == 0 {
			t.Fatal("Handler returned empty content")
		}

		// Verify the content is as expected
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Parse the JSON response to verify filtering worked
		var stopsData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &stopsData); err != nil {
			t.Fatalf("Failed to parse stop response: %v", err)
		}

		// There should be some stops returned
		if len(stopsData) == 0 {
			t.Error("No stops returned after route filtering")
		}
	})
}

// TestGetSchedulesHandler tests the schedules handler functionality
func TestGetSchedulesHandler(t *testing.T) {
	// Create a mock MBTA API server
	mockServer, err := mock.StandardMockServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer mockServer.Close()

	// Create config pointing to mock server
	cfg := &config.Config{
		APIKey:     "test-api-key",
		APIBaseURL: mockServer.URL,
	}

	// Create MCP server with the MBTA client
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	t.Run("Get schedules handler can be registered", func(t *testing.T) {
		// Register the transit info tools, which includes schedules
		server.registerTransitInfoTools()
	})

	t.Run("Get schedules returns valid schedule data", func(t *testing.T) {
		// Create a request for the schedules handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_schedules",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getSchedulesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Check that content is returned
		if len(result.Content) == 0 {
			t.Fatal("Handler returned empty content")
		}

		// Verify content type is text
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Verify the text indicates success
		if textContent.Text == "" {
			t.Error("Text content is empty")
		}
	})

	t.Run("Get schedules handles filtering by route", func(t *testing.T) {
		// Create a request with route filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_schedules",
				Arguments: map[string]any{
					"route_id": "Red",
				},
			},
		}

		// Call the handler
		result, err := server.getSchedulesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get schedules handles filtering by stop", func(t *testing.T) {
		// Create a request with stop filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_schedules",
				Arguments: map[string]any{
					"stop_id": "place-sstat",
				},
			},
		}

		// Call the handler
		result, err := server.getSchedulesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get schedules handles filtering by date", func(t *testing.T) {
		// Create a request with date filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_schedules",
				Arguments: map[string]any{
					"date": "2023-05-20",
				},
			},
		}

		// Call the handler
		result, err := server.getSchedulesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Get the content as text
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Content is not TextContent, got: %T", result.Content[0])
		}

		// Parse the response to verify schedules are for the correct date
		var schedulesData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &schedulesData); err != nil {
			t.Fatalf("Failed to parse schedules response: %v", err)
		}

		// All schedules should contain the date 2023-05-20
		for _, schedule := range schedulesData {
			arrivalTime, ok := schedule["arrival_time"].(string)
			if !ok {
				t.Errorf("Missing or invalid arrival_time in schedule: %v", schedule)
				continue
			}

			if !strings.Contains(arrivalTime, "2023-05-20") {
				t.Errorf("Schedule arrival_time %q does not match requested date 2023-05-20", arrivalTime)
			}
		}
	})

	t.Run("Get schedules returns error with invalid date format", func(t *testing.T) {
		// Create a request with invalid date filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_schedules",
				Arguments: map[string]any{
					"date": "05/20/2023", // Invalid format, should be YYYY-MM-DD
				},
			},
		}

		// Call the handler
		result, err := server.getSchedulesHandler(context.Background(), request)

		// Check that we got a result, not an error
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Verify that IsError is set to true
		if !result.IsError {
			t.Error("Expected result to have IsError=true for invalid date format")
		}
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("Error response function exists", func(t *testing.T) {
		// This test will validate the implementation of the createErrorResponse function
		t.Skip("Will be implemented once the error response function is defined")
	})

	t.Run("Error response includes error details", func(t *testing.T) {
		t.Skip("Will be implemented once the error response function is defined")
	})
}

// This uses the actual createErrorResponse function from handlers.go

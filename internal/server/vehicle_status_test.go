// ABOUTME: This file contains tests for the vehicle status handlers.
// ABOUTME: It verifies the proper handling of vehicle status requests.

package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/mock"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetVehicleStatusHandler(t *testing.T) {
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

	t.Run("Get vehicle status returns valid status data", func(t *testing.T) {
		// Create a request for the status handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_vehicle_status",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getVehicleStatusHandler(context.Background(), request)

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

		// Verify the text contains status data by checking if it can be parsed as JSON
		var statusData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &statusData); err != nil {
			t.Fatalf("Failed to parse response as JSON: %v", err)
		}

		// Verify we have status data
		if len(statusData) == 0 {
			t.Error("No status updates returned in JSON")
		}

		// Verify the status data has expected fields
		for _, status := range statusData {
			requiredFields := []string{"vehicle_id", "label", "status", "latitude", "longitude"}
			for _, field := range requiredFields {
				if _, ok := status[field]; !ok {
					t.Errorf("Status missing required field '%s'", field)
				}
			}
		}
	})

	t.Run("Get vehicle status handles filtering by route", func(t *testing.T) {
		// Create a request with route filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_status",
				Arguments: map[string]any{
					"route_id": "Red",
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleStatusHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get vehicle status handles filtering by status type", func(t *testing.T) {
		// Create a request with status type filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_status",
				Arguments: map[string]any{
					"status_type": "in_transit",
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleStatusHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get vehicle status handles limit parameter", func(t *testing.T) {
		// Create a request with limit parameter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_status",
				Arguments: map[string]any{
					"limit": 5,
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleStatusHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Verify result content exists - we don't need to parse it for this test
		if len(result.Content) == 0 {
			t.Fatal("Handler returned empty content")
		}
	})

	t.Run("Get vehicle status handles invalid status type", func(t *testing.T) {
		// Create a request with invalid status type
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_status",
				Arguments: map[string]any{
					"status_type": "invalid",
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleStatusHandler(context.Background(), request)

		// We should still get a result, but it might have an error flag
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}

		// Check if the result indicates an error
		if !result.IsError {
			t.Error("Expected IsError to be true for invalid status type")
		}
	})
}
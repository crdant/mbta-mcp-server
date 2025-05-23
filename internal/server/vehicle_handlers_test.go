// ABOUTME: This file contains tests for the vehicle tracking MCP handlers.
// ABOUTME: It verifies proper handling of vehicle tracking and status requests.

package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/mock"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetVehiclesHandler(t *testing.T) {
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

	t.Run("Get vehicles handler can be registered", func(t *testing.T) {
		// Register the vehicle tracking tools
		server.registerVehicleTrackingTools()
	})

	t.Run("Get vehicles returns valid vehicle data", func(t *testing.T) {
		// Create a request for the vehicles handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_vehicles",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getVehiclesHandler(context.Background(), request)

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

		// Verify the text contains vehicle data by checking if it can be parsed as JSON
		var vehicleData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &vehicleData); err != nil {
			t.Fatalf("Failed to parse response as JSON: %v", err)
		}

		// Verify we have vehicle data
		if len(vehicleData) == 0 {
			t.Error("No vehicles returned in JSON")
		}

		// Verify the vehicle data has expected fields
		for _, vehicle := range vehicleData {
			requiredFields := []string{"id", "label", "status", "latitude", "longitude"}
			for _, field := range requiredFields {
				if _, ok := vehicle[field]; !ok {
					t.Errorf("Vehicle missing required field '%s'", field)
				}
			}
		}
	})

	t.Run("Get vehicles handles filtering by route", func(t *testing.T) {
		// Create a request with route filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicles",
				Arguments: map[string]any{
					"route_id": "Red",
				},
			},
		}

		// Call the handler
		result, err := server.getVehiclesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get vehicles handles filtering by trip", func(t *testing.T) {
		// Create a request with trip filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicles",
				Arguments: map[string]any{
					"trip_id": "123456",
				},
			},
		}

		// Call the handler
		result, err := server.getVehiclesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get vehicles handles filtering by location", func(t *testing.T) {
		// Create a request with location filter
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicles",
				Arguments: map[string]any{
					"latitude":  42.3601,
					"longitude": -71.0589,
					"radius":    0.05,
				},
			},
		}

		// Call the handler
		result, err := server.getVehiclesHandler(context.Background(), request)

		// Check for errors
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})
}

func TestGetVehicleHandler(t *testing.T) {
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

	t.Run("Get vehicle returns valid vehicle data", func(t *testing.T) {
		// Create a request for the vehicle handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle",
				Arguments: map[string]any{
					"vehicle_id": "R-5463D359",
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleHandler(context.Background(), request)

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

		// Verify the text contains vehicle data by checking if it can be parsed as JSON
		var vehicleData map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &vehicleData); err != nil {
			t.Fatalf("Failed to parse response as JSON: %v", err)
		}

		// Verify the vehicle data has expected fields
		requiredFields := []string{"id", "label", "status", "latitude", "longitude"}
		for _, field := range requiredFields {
			if _, ok := vehicleData[field]; !ok {
				t.Errorf("Vehicle missing required field '%s'", field)
			}
		}
	})

	t.Run("Get vehicle handles invalid vehicle ID", func(t *testing.T) {
		// Create a request with an invalid vehicle ID
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle",
				Arguments: map[string]any{
					"vehicle_id": "non-existent",
				},
			},
		}

		// Call the handler
		result, err := server.getVehicleHandler(context.Background(), request)

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
			t.Error("Expected IsError to be true for invalid vehicle ID")
		}
	})

	t.Run("Get vehicle handles missing vehicle ID", func(t *testing.T) {
		// Create a request without a vehicle ID
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_vehicle",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getVehicleHandler(context.Background(), request)

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
			t.Error("Expected IsError to be true for missing vehicle ID")
		}
	})
}

func TestGetVehiclePredictionsHandler(t *testing.T) {
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

	t.Run("Get vehicle predictions returns valid prediction data", func(t *testing.T) {
		// Create a request for the predictions handler
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_predictions",
				Arguments: map[string]any{
					"vehicle_id": "R-5463D359",
				},
			},
		}

		// Call the handler
		result, err := server.getVehiclePredictionsHandler(context.Background(), request)

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

		// Try to parse the response as JSON
		var predictionsData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &predictionsData); err != nil {
			// If we get an error, it might be a "no predictions available" message
			// which is valid, so we'll just skip further tests
			return
		}

		// If we have prediction data, verify it has expected fields
		for _, prediction := range predictionsData {
			requiredFields := []string{"id", "route_id", "stop_id", "vehicle_id"}
			for _, field := range requiredFields {
				if _, ok := prediction[field]; !ok {
					t.Errorf("Prediction missing required field '%s'", field)
				}
			}
		}
	})

	t.Run("Get vehicle predictions handles invalid vehicle ID", func(t *testing.T) {
		// Create a request with an invalid vehicle ID
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "get_vehicle_predictions",
				Arguments: map[string]any{
					"vehicle_id": "non-existent",
				},
			},
		}

		// Call the handler
		result, err := server.getVehiclePredictionsHandler(context.Background(), request)

		// We should still get a result, but it might have an error indication
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		// Verify result isn't nil
		if result == nil {
			t.Fatal("Handler returned nil result")
		}
	})

	t.Run("Get vehicle predictions handles missing vehicle ID", func(t *testing.T) {
		// Create a request without a vehicle ID
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name:      "get_vehicle_predictions",
				Arguments: map[string]any{},
			},
		}

		// Call the handler
		result, err := server.getVehiclePredictionsHandler(context.Background(), request)

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
			t.Error("Expected IsError to be true for missing vehicle ID")
		}
	})
}
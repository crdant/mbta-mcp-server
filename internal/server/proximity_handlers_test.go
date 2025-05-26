// ABOUTME: This file contains tests for the station proximity handlers for the MCP server.
// ABOUTME: It tests the findNearbyStationsHandler functionality.

package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestFindNearbyStationsHandler(t *testing.T) {
	// Create a minimal test server for the handler tests
	cfg := &config.Config{
		APIKey:      "test-api-key",
		Debug:       false,
		LogLevel:    "info",
		Timeout:     30,
		APIBaseURL:  "https://api-test.mbta.com",
		Environment: "test",
	}

	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test cases
	t.Run("Validates required parameters", func(t *testing.T) {
		// Create a request with missing parameters
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "find_nearby_stations",
				Arguments: map[string]any{
					// Missing latitude and longitude
					"radius": 1.0,
				},
			},
		}

		// Call the handler
		response, err := server.findNearbyStationsHandler(context.Background(), request)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return an error response due to missing required parameters
		if response == nil {
			t.Fatal("Expected error response, got nil")
		}

		// Check for error message in the response
		textContent, ok := response.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", response.Content[0])
		}

		// Parse the error message
		var errorResponse map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &errorResponse); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		// Verify error message
		if errorMessage, ok := errorResponse["error"].(string); !ok || errorMessage == "" {
			t.Errorf("Expected error message, got: %v", errorResponse)
		}
	})

	t.Run("Validates parameter types", func(t *testing.T) {
		// Create a request with incorrect parameter types
		request := mcp.CallToolRequest{
			Params: struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments,omitempty"`
				Meta      *struct {
					ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
				} `json:"_meta,omitempty"`
			}{
				Name: "find_nearby_stations",
				Arguments: map[string]any{
					"latitude":  "not-a-number", // Should be a number
					"longitude": -71.06,
					"radius":    1.0,
				},
			},
		}

		// Call the handler
		response, err := server.findNearbyStationsHandler(context.Background(), request)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return an error response due to invalid parameter type
		if response == nil {
			t.Fatal("Expected error response, got nil")
		}

		// Check for error message in the response
		textContent, ok := response.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("Expected TextContent, got %T", response.Content[0])
		}

		// Parse the error message
		var errorResponse map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &errorResponse); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}

		// Verify error message
		if errorMessage, ok := errorResponse["error"].(string); !ok || errorMessage == "" {
			t.Errorf("Expected error message, got: %v", errorResponse)
		}
	})
}

// Mock code removed to fix linting errors

func TestFindNearbyStationsResponse(t *testing.T) {
	// Test formatting of station proximity results
	testStops := []models.NearbyStation{
		{
			Stop: models.Stop{
				ID:   "place-pktrm",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Park Street",
					Latitude:           42.35639457,
					Longitude:          -71.0624242,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			DistanceKm: 0.25,
		},
		{
			Stop: models.Stop{
				ID:   "place-dwnxg",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Downtown Crossing",
					Latitude:           42.355518,
					Longitude:          -71.060225,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			DistanceKm: 0.5,
		},
	}

	// Create response with the test data
	response, err := formatNearbyStationsResponse(testStops)
	if err != nil {
		t.Fatalf("Failed to format response: %v", err)
	}

	// Verify response structure
	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	// Check for content in the response
	if len(response.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(response.Content))
	}

	// Check that the content is text
	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", response.Content[0])
	}

	// Parse the JSON response
	var responseData []map[string]interface{}
	if err := json.Unmarshal([]byte(textContent.Text), &responseData); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify the response contains the expected number of stations
	if len(responseData) != len(testStops) {
		t.Errorf("Expected %d stations in response, got %d", len(testStops), len(responseData))
	}

	// Verify each station has the required fields
	for i, station := range responseData {
		// Check required fields
		requiredFields := []string{"id", "name", "distance_km", "latitude", "longitude", "municipality", "wheelchair_accessible"}
		for _, field := range requiredFields {
			if _, ok := station[field]; !ok {
				t.Errorf("Station %d missing required field: %s", i, field)
			}
		}

		// Check ID matches original
		if station["id"] != testStops[i].Stop.ID {
			t.Errorf("Expected station %d ID to be %s, got %s", i, testStops[i].Stop.ID, station["id"])
		}

		// Check distance matches original
		if station["distance_km"] != testStops[i].DistanceKm {
			t.Errorf("Expected station %d distance to be %f, got %v", i, testStops[i].DistanceKm, station["distance_km"])
		}
	}
}

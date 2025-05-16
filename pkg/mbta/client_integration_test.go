package mbta

import (
	"context"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/mock"
)

func TestClientIntegration(t *testing.T) {
	// Skip in short mode since these are integration tests
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Create a mock server
	server, err := mock.StandardMockServer()
	if err != nil {
		t.Fatalf("Failed to create mock server: %v", err)
	}
	defer server.Close()

	// Create client configuration using the mock server
	cfg := &config.Config{
		APIKey:     "valid-key",
		Timeout:    5 * time.Second,
		APIBaseURL: server.URL,
	}

	// Create client
	client := NewClient(cfg)

	// Test context
	ctx := context.Background()

	// Test GetRoutes
	t.Run("GetRoutes", func(t *testing.T) {
		routes, err := client.GetRoutes(ctx)
		if err != nil {
			t.Fatalf("GetRoutes failed: %v", err)
		}

		if len(routes) != 2 {
			t.Errorf("Expected 2 routes, got %d", len(routes))
		}

		// Check if the Red Line is in the routes
		var foundRedLine bool
		for _, route := range routes {
			if route.ID == "Red" {
				foundRedLine = true
				if route.Attributes.LongName != "Red Line" {
					t.Errorf("Expected Red Line long name 'Red Line', got '%s'", route.Attributes.LongName)
				}
				if route.Attributes.Type != 1 {
					t.Errorf("Expected Red Line type 1, got %d", route.Attributes.Type)
				}
			}
		}

		if !foundRedLine {
			t.Error("Red Line route not found in results")
		}
	})

	// Test GetRoute
	t.Run("GetRoute", func(t *testing.T) {
		route, err := client.GetRoute(ctx, "Red")
		if err != nil {
			t.Fatalf("GetRoute failed: %v", err)
		}

		if route.ID != "Red" {
			t.Errorf("Expected route ID 'Red', got '%s'", route.ID)
		}
		if route.Attributes.LongName != "Red Line" {
			t.Errorf("Expected long name 'Red Line', got '%s'", route.Attributes.LongName)
		}
	})

	// Test GetStops
	t.Run("GetStops", func(t *testing.T) {
		stops, err := client.GetStops(ctx)
		if err != nil {
			t.Fatalf("GetStops failed: %v", err)
		}

		if len(stops) != 2 {
			t.Errorf("Expected 2 stops, got %d", len(stops))
		}

		// Check if North Station is in the stops
		var foundNorthStation bool
		for _, stop := range stops {
			if stop.ID == "place-north" {
				foundNorthStation = true
				if stop.Attributes.Name != "North Station" {
					t.Errorf("Expected stop name 'North Station', got '%s'", stop.Attributes.Name)
				}
				if stop.Attributes.LocationType != 1 {
					t.Errorf("Expected location type 1, got %d", stop.Attributes.LocationType)
				}
			}
		}

		if !foundNorthStation {
			t.Error("North Station stop not found in results")
		}
	})

	// Test GetSchedules
	t.Run("GetSchedules", func(t *testing.T) {
		schedules, included, err := client.GetSchedules(ctx, map[string]string{
			"filter[route]": "Red",
		})
		if err != nil {
			t.Fatalf("GetSchedules failed: %v", err)
		}

		if len(schedules) != 2 {
			t.Errorf("Expected 2 schedules, got %d", len(schedules))
		}
		if len(included) != 1 {
			t.Errorf("Expected 1 included object, got %d", len(included))
		}

		// Check the schedule data
		if schedules[0].ID != "schedule-1" {
			t.Errorf("Expected schedule ID 'schedule-1', got '%s'", schedules[0].ID)
		}
		if schedules[0].Attributes.StopHeadsign != "Alewife" {
			t.Errorf("Expected stop headsign 'Alewife', got '%s'", schedules[0].Attributes.StopHeadsign)
		}

		// Check included trip data
		if included[0].ID != "Red-123456-20230520" {
			t.Errorf("Expected included trip ID 'Red-123456-20230520', got '%s'", included[0].ID)
		}
		if included[0].Type != "trip" {
			t.Errorf("Expected included type 'trip', got '%s'", included[0].Type)
		}
	})

	// Test error handling - invalid route ID
	t.Run("InvalidRouteID", func(t *testing.T) {
		_, err := client.GetRoute(ctx, "Invalid")
		if err == nil {
			t.Error("Expected error for invalid route ID, got nil")
		}

		// Check if it's an API error
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Errorf("Expected error type *APIError, got %T", err)
		}
		if apiErr != nil && !apiErr.IsNotFoundError() {
			t.Errorf("Expected not found error, got %v", apiErr)
		}
	})

	// Test authentication error - invalid API key
	t.Run("InvalidAPIKey", func(t *testing.T) {
		// Create client with invalid API key
		invalidCfg := &config.Config{
			APIKey:     "invalid-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		invalidClient := NewClient(invalidCfg)

		_, err := invalidClient.GetRoutes(ctx)
		if err == nil {
			t.Error("Expected error for invalid API key, got nil")
		}

		// Check if it's an API error
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Errorf("Expected error type *APIError, got %T", err)
		}
		if apiErr != nil && !apiErr.IsAuthError() {
			t.Errorf("Expected auth error, got %v", apiErr)
		}
	})
}
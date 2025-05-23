package server

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestFormatTripPlanResponse(t *testing.T) {
	// Create a sample trip plan
	now := time.Now()
	later := now.Add(30 * time.Minute)
	
	origin := &models.Stop{
		ID: "place-harvard",
		Attributes: models.StopAttributes{
			Name:         "Harvard",
			Municipality: "Cambridge",
			Latitude:     42.3736,
			Longitude:    -71.1190,
		},
	}
	
	destination := &models.Stop{
		ID: "place-porter",
		Attributes: models.StopAttributes{
			Name:         "Porter",
			Municipality: "Cambridge",
			Latitude:     42.3884,
			Longitude:    -71.1191,
		},
	}
	
	tripPlan := &models.TripPlan{
		Origin:        origin,
		Destination:   destination,
		DepartureTime: now,
		ArrivalTime:   later,
		Duration:      30 * time.Minute,
		Legs: []models.TripLeg{
			{
				Origin:        origin,
				Destination:   destination,
				RouteID:       "Red",
				RouteName:     "Red Line",
				TripID:        "trip-1234",
				DepartureTime: now,
				ArrivalTime:   later,
				Duration:      30 * time.Minute,
				Distance:      2.5,
				Headsign:      "Alewife",
				DirectionID:   0,
				IsAccessible:  true,
				Instructions:  "Board the Red Line toward Alewife",
			},
		},
		TotalDistance:  2.5,
		AccessibleTrip: true,
	}

	// Format the trip plan
	response, err := formatTripPlanResponse(tripPlan)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that response is not nil and is not an error
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.IsError {
		t.Errorf("Response indicates an error: %v", response)
	}

	// Check content
	if len(response.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(response.Content))
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", response.Content[0])
	}

	// Verify JSON content can be parsed
	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify key fields
	originData, ok := responseData["origin"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing or invalid origin in response")
	} else {
		if originID, ok := originData["id"].(string); !ok || originID != "place-harvard" {
			t.Errorf("Expected origin ID 'place-harvard', got '%v'", originData["id"])
		}
	}

	destData, ok := responseData["destination"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing or invalid destination in response")
	} else {
		if destID, ok := destData["id"].(string); !ok || destID != "place-porter" {
			t.Errorf("Expected destination ID 'place-porter', got '%v'", destData["id"])
		}
	}

	// Check legs
	legs, ok := responseData["legs"].([]interface{})
	if !ok {
		t.Errorf("Missing or invalid legs in response")
	} else {
		if len(legs) != 1 {
			t.Errorf("Expected 1 leg, got %d", len(legs))
		} else {
			leg, ok := legs[0].(map[string]interface{})
			if !ok {
				t.Errorf("Invalid leg format")
			} else {
				if routeID, ok := leg["route_id"].(string); !ok || routeID != "Red" {
					t.Errorf("Expected route_id 'Red', got '%v'", leg["route_id"])
				}
				if headsign, ok := leg["headsign"].(string); !ok || headsign != "Alewife" {
					t.Errorf("Expected headsign 'Alewife', got '%v'", leg["headsign"])
				}
			}
		}
	}
}

func TestFormatTransferPointsResponse(t *testing.T) {
	// Create sample transfer points
	transferPoints := []models.TransferPoint{
		{
			Stop: &models.Stop{
				ID: "place-dwnxg",
				Attributes: models.StopAttributes{
					Name:         "Downtown Crossing",
					Municipality: "Boston",
					Latitude:     42.3554,
					Longitude:    -71.0603,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
				},
			},
			FromRoute:       "Red",
			ToRoute:         "Orange",
			TransferType:    models.TransferTypeRecommended,
			MinTransferTime: 3 * time.Minute,
		},
	}

	// Format the transfer points
	response, err := formatTransferPointsResponse(transferPoints)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that response is not nil and is not an error
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.IsError {
		t.Errorf("Response indicates an error: %v", response)
	}

	// Check content
	if len(response.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(response.Content))
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", response.Content[0])
	}

	// Verify JSON content can be parsed
	var transferData []map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &transferData)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify key fields
	if len(transferData) != 1 {
		t.Fatalf("Expected 1 transfer point, got %d", len(transferData))
	}

	transfer := transferData[0]
	if stopID, ok := transfer["stop_id"].(string); !ok || stopID != "place-dwnxg" {
		t.Errorf("Expected stop_id 'place-dwnxg', got '%v'", transfer["stop_id"])
	}

	if fromRoute, ok := transfer["from_route"].(string); !ok || fromRoute != "Red" {
		t.Errorf("Expected from_route 'Red', got '%v'", transfer["from_route"])
	}

	if toRoute, ok := transfer["to_route"].(string); !ok || toRoute != "Orange" {
		t.Errorf("Expected to_route 'Orange', got '%v'", transfer["to_route"])
	}

	if minTime, ok := transfer["min_transfer_time"].(float64); !ok || minTime != 3.0 {
		t.Errorf("Expected min_transfer_time 3.0, got %v", transfer["min_transfer_time"])
	}

	if accessible, ok := transfer["wheelchair_accessible"].(bool); !ok || !accessible {
		t.Errorf("Expected wheelchair_accessible true, got %v", transfer["wheelchair_accessible"])
	}
}

func TestFormatTravelTimeResponse(t *testing.T) {
	// Create sample stops
	origin := &models.Stop{
		ID: "place-north",
		Attributes: models.StopAttributes{
			Name:         "North Station",
			Municipality: "Boston",
			Latitude:     42.3654,
			Longitude:    -71.0613,
			LocationType: models.LocationTypeStation,
		},
	}
	
	destination := &models.Stop{
		ID: "place-sstat",
		Attributes: models.StopAttributes{
			Name:         "South Station",
			Municipality: "Boston",
			Latitude:     42.3523,
			Longitude:    -71.0551,
			LocationType: models.LocationTypeStation,
		},
	}

	// Format the travel time
	distance := 2.1    // kilometers
	timeMinutes := 15.0 // minutes
	source := "Based on recent schedules"
	
	response, err := formatTravelTimeResponse(origin, destination, distance, timeMinutes, source)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that response is not nil and is not an error
	if response == nil {
		t.Fatal("Expected non-nil response")
	}
	
	if response.IsError {
		t.Errorf("Response indicates an error: %v", response)
	}

	// Check content
	if len(response.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(response.Content))
	}

	textContent, ok := response.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", response.Content[0])
	}

	// Verify JSON content can be parsed
	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(textContent.Text), &responseData)
	if err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify key fields
	originData, ok := responseData["origin"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing or invalid origin in response")
	} else {
		if originID, ok := originData["id"].(string); !ok || originID != "place-north" {
			t.Errorf("Expected origin ID 'place-north', got '%v'", originData["id"])
		}
	}

	destData, ok := responseData["destination"].(map[string]interface{})
	if !ok {
		t.Errorf("Missing or invalid destination in response")
	} else {
		if destID, ok := destData["id"].(string); !ok || destID != "place-sstat" {
			t.Errorf("Expected destination ID 'place-sstat', got '%v'", destData["id"])
		}
	}

	if dist, ok := responseData["distance_km"].(float64); !ok || dist != distance {
		t.Errorf("Expected distance_km %f, got %v", distance, responseData["distance_km"])
	}

	if mins, ok := responseData["estimated_minutes"].(float64); !ok || mins != timeMinutes {
		t.Errorf("Expected estimated_minutes %f, got %v", timeMinutes, responseData["estimated_minutes"])
	}

	if src, ok := responseData["estimation_source"].(string); !ok || src != source {
		t.Errorf("Expected estimation_source '%s', got '%v'", source, responseData["estimation_source"])
	}

	if formatted, ok := responseData["formatted_time"].(string); !ok || formatted != "15 minutes" {
		t.Errorf("Expected formatted_time '15 minutes', got '%v'", responseData["formatted_time"])
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		minutes  float64
		expected string
	}{
		{0, "0 minutes"},
		{1, "1 minute"},
		{5, "5 minutes"},
		{60, "1 hour 0 minutes"},
		{61, "1 hour 1 minute"},
		{90, "1 hour 30 minutes"},
		{120, "2 hours 0 minutes"},
		{135, "2 hours 15 minutes"},
	}

	for _, test := range tests {
		result := formatDuration(test.minutes)
		if result != test.expected {
			t.Errorf("Expected formatDuration(%f) to be '%s', got '%s'", test.minutes, test.expected, result)
		}
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		count    int
		expected string
	}{
		{0, "s"},
		{1, ""},
		{2, "s"},
		{5, "s"},
	}

	for _, test := range tests {
		result := pluralize(test.count)
		if result != test.expected {
			t.Errorf("Expected pluralize(%d) to be '%s', got '%s'", test.count, test.expected, result)
		}
	}
}

// Mock for testing planTripHandler
func mockPlanTripHandler(t *testing.T) {
	// This would be a more complex test requiring mocks for the client
	// For brevity, we're not implementing the full functionality here
	cfg := &config.Config{}
	server, _ := New(cfg)
	
	// Test with valid parameters
	validRequest := mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *struct {
				ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
			} `json:"_meta,omitempty"`
		}{
			Name: "plan_trip",
			Arguments: map[string]any{
				"origin_stop_id":      "place-harvard",
				"destination_stop_id": "place-porter",
			},
		},
	}
	
	// This is just a basic structure test - in a real implementation,
	// we would use mocks to avoid actual API calls
	_, err := server.planTripHandler(context.Background(), validRequest)
	if err != nil {
		t.Logf("Got expected error in test environment: %v", err)
	}
}
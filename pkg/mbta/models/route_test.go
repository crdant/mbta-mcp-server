package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestRouteUnmarshal(t *testing.T) {
	// Load the sample routes data
	fixtureData, err := os.ReadFile("../../../test/fixtures/routes.json")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Unmarshal the JSON
	var response RouteResponse
	if err := json.Unmarshal(fixtureData, &response); err != nil {
		t.Fatalf("Failed to unmarshal route response: %v", err)
	}

	// Verify that we have the correct number of routes
	if len(response.Data) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(response.Data))
	}

	// Verify the first route
	route := response.Data[0]
	if route.ID != "Red" {
		t.Errorf("Expected route ID 'Red', got '%s'", route.ID)
	}
	if route.Type != "route" {
		t.Errorf("Expected route type 'route', got '%s'", route.Type)
	}

	attrs := route.Attributes
	if attrs.LongName != "Red Line" {
		t.Errorf("Expected long name 'Red Line', got '%s'", attrs.LongName)
	}
	if attrs.Color != "DA291C" {
		t.Errorf("Expected color 'DA291C', got '%s'", attrs.Color)
	}
	if attrs.Description != "Rapid Transit" {
		t.Errorf("Expected description 'Rapid Transit', got '%s'", attrs.Description)
	}
	if attrs.Type != 1 {
		t.Errorf("Expected type 1, got %d", attrs.Type)
	}

	// Check direction destinations
	if len(attrs.DirectionDestinations) != 2 {
		t.Errorf("Expected 2 direction destinations, got %d", len(attrs.DirectionDestinations))
	}
	if attrs.DirectionDestinations[0] != "Alewife" {
		t.Errorf("Expected first destination 'Alewife', got '%s'", attrs.DirectionDestinations[0])
	}
	if attrs.DirectionDestinations[1] != "Ashmont/Braintree" {
		t.Errorf("Expected second destination 'Ashmont/Braintree', got '%s'", attrs.DirectionDestinations[1])
	}

	// Check relationships
	if line, ok := route.Relationships["line"]; ok {
		if lineData, ok := line.(map[string]interface{})["data"]; ok {
			lineDataMap, ok := lineData.(map[string]interface{})
			if !ok {
				t.Errorf("Expected line data to be a map, got %T", lineData)
			} else if lineID, ok := lineDataMap["id"]; ok {
				if lineID != "line-Red" {
					t.Errorf("Expected line ID 'line-Red', got '%s'", lineID)
				}
			} else {
				t.Error("Missing line ID in relationship")
			}
		} else {
			t.Error("Missing data in line relationship")
		}
	} else {
		t.Error("Missing line relationship")
	}
}

func TestRouteMarshal(t *testing.T) {
	// Create a route object
	route := Route{
		ID:   "Red",
		Type: "route",
		Attributes: RouteAttributes{
			Color:                 "DA291C",
			Description:           "Rapid Transit",
			DirectionDestinations: []string{"Alewife", "Ashmont/Braintree"},
			DirectionNames:        []string{"Outbound", "Inbound"},
			FareClass:             "Rapid Transit",
			LongName:              "Red Line",
			ShortName:             "",
			SortOrder:             10010,
			TextColor:             "FFFFFF",
			Type:                  1,
		},
		Links: map[string]string{
			"self": "/routes/Red",
		},
		Relationships: map[string]interface{}{
			"line": map[string]interface{}{
				"data": map[string]string{
					"id":   "line-Red",
					"type": "line",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to marshal route: %v", err)
	}

	// Unmarshal back to verify
	var roundTrip Route
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("Failed to unmarshal route: %v", err)
	}

	// Verify key fields
	if roundTrip.ID != route.ID {
		t.Errorf("Expected route ID '%s', got '%s'", route.ID, roundTrip.ID)
	}
	if roundTrip.Attributes.LongName != route.Attributes.LongName {
		t.Errorf("Expected long name '%s', got '%s'", route.Attributes.LongName, roundTrip.Attributes.LongName)
	}
	if roundTrip.Attributes.Type != route.Attributes.Type {
		t.Errorf("Expected type %d, got %d", route.Attributes.Type, roundTrip.Attributes.Type)
	}
}

func TestGetRouteTypeDescription(t *testing.T) {
	tests := []struct {
		routeType int
		expected  string
	}{
		{0, "Light Rail"},
		{1, "Subway"},
		{2, "Commuter Rail"},
		{3, "Bus"},
		{4, "Ferry"},
		{100, "Unknown"},
	}

	for _, test := range tests {
		result := GetRouteTypeDescription(test.routeType)
		if result != test.expected {
			t.Errorf("For route type %d, expected '%s', got '%s'", test.routeType, test.expected, result)
		}
	}
}

func TestRoute_GetDirectionName(t *testing.T) {
	route := Route{
		Attributes: RouteAttributes{
			DirectionNames: []string{"Outbound", "Inbound"},
		},
	}

	if name := route.GetDirectionName(0); name != "Outbound" {
		t.Errorf("Expected direction 0 to be 'Outbound', got '%s'", name)
	}

	if name := route.GetDirectionName(1); name != "Inbound" {
		t.Errorf("Expected direction 1 to be 'Inbound', got '%s'", name)
	}

	// Test invalid direction
	if name := route.GetDirectionName(2); name != "" {
		t.Errorf("Expected direction 2 to be empty string, got '%s'", name)
	}
}

func TestRoute_GetDirectionDestination(t *testing.T) {
	route := Route{
		Attributes: RouteAttributes{
			DirectionDestinations: []string{"Alewife", "Ashmont/Braintree"},
		},
	}

	if dest := route.GetDirectionDestination(0); dest != "Alewife" {
		t.Errorf("Expected destination 0 to be 'Alewife', got '%s'", dest)
	}

	if dest := route.GetDirectionDestination(1); dest != "Ashmont/Braintree" {
		t.Errorf("Expected destination 1 to be 'Ashmont/Braintree', got '%s'", dest)
	}

	// Test invalid direction
	if dest := route.GetDirectionDestination(2); dest != "" {
		t.Errorf("Expected destination 2 to be empty string, got '%s'", dest)
	}
}

func TestRoute_GetTypeDescription(t *testing.T) {
	route := Route{
		Attributes: RouteAttributes{
			Type: 1,
		},
	}

	if desc := route.GetTypeDescription(); desc != "Subway" {
		t.Errorf("Expected type description to be 'Subway', got '%s'", desc)
	}
}

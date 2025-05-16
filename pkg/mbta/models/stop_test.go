package models

import (
	"encoding/json"
	"os"
	"testing"
)

func TestStopUnmarshal(t *testing.T) {
	// Load the sample stops data
	fixtureData, err := os.ReadFile("testdata/stops.json")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Unmarshal the JSON
	var response StopResponse
	if err := json.Unmarshal(fixtureData, &response); err != nil {
		t.Fatalf("Failed to unmarshal stop response: %v", err)
	}

	// Verify that we have the correct number of stops
	if len(response.Data) != 2 {
		t.Errorf("Expected 2 stops, got %d", len(response.Data))
	}

	// Verify the first stop (station)
	station := response.Data[0]
	if station.ID != "place-north" {
		t.Errorf("Expected stop ID 'place-north', got '%s'", station.ID)
	}
	if station.Type != "stop" {
		t.Errorf("Expected stop type 'stop', got '%s'", station.Type)
	}

	attrs := station.Attributes
	if attrs.Name != "North Station" {
		t.Errorf("Expected name 'North Station', got '%s'", attrs.Name)
	}
	if attrs.LocationType != LocationTypeStation {
		t.Errorf("Expected location_type %d, got %d", LocationTypeStation, attrs.LocationType)
	}
	if attrs.Latitude != 42.365577 {
		t.Errorf("Expected latitude 42.365577, got %f", attrs.Latitude)
	}
	if attrs.Longitude != -71.06129 {
		t.Errorf("Expected longitude -71.06129, got %f", attrs.Longitude)
	}
	if attrs.WheelchairBoarding != WheelchairBoardingAccessible {
		t.Errorf("Expected wheelchair_boarding %d, got %d", WheelchairBoardingAccessible, attrs.WheelchairBoarding)
	}

	// Verify the second stop (platform)
	platform := response.Data[1]
	if platform.ID != "70061" {
		t.Errorf("Expected stop ID '70061', got '%s'", platform.ID)
	}
	
	platformAttrs := platform.Attributes
	if platformAttrs.LocationType != LocationTypePlatform {
		t.Errorf("Expected location_type %d, got %d", LocationTypePlatform, platformAttrs.LocationType)
	}
	if platformAttrs.PlatformCode != "1" {
		t.Errorf("Expected platform_code '1', got '%s'", platformAttrs.PlatformCode)
	}
	if platformAttrs.PlatformName != "Orange Line - Forest Hills" {
		t.Errorf("Expected platform_name 'Orange Line - Forest Hills', got '%s'", platformAttrs.PlatformName)
	}

	// Check parent station relationship
	if parent, ok := platform.Relationships["parent_station"]; ok {
		if parentData, ok := parent.(map[string]interface{})["data"]; ok {
			parentDataMap, ok := parentData.(map[string]interface{})
			if !ok {
				t.Errorf("Expected parent_station data to be a map, got %T", parentData)
			} else if parentID, ok := parentDataMap["id"]; ok {
				if parentID != "place-north" {
					t.Errorf("Expected parent station ID 'place-north', got '%s'", parentID)
				}
			} else {
				t.Error("Missing parent station ID in relationship")
			}
		} else {
			t.Error("Missing data in parent_station relationship")
		}
	} else {
		t.Error("Missing parent_station relationship")
	}
}

func TestStopMarshal(t *testing.T) {
	// Create a stop object
	stop := Stop{
		ID:   "place-north",
		Type: "stop",
		Attributes: StopAttributes{
			Address:            "North Station, Boston, MA 02114",
			Description:        "North Station - Commuter Rail, Orange Line, and Green Line",
			Latitude:           42.365577,
			LocationType:       LocationTypeStation,
			Longitude:          -71.06129,
			Municipality:       "Boston",
			Name:               "North Station",
			WheelchairBoarding: WheelchairBoardingAccessible,
		},
		Links: map[string]string{
			"self": "/stops/place-north",
		},
		Relationships: map[string]interface{}{
			"parent_station": map[string]interface{}{
				"data": nil,
			},
			"zone": map[string]interface{}{
				"data": map[string]string{
					"id":   "CR-zone-1A",
					"type": "zone",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(stop)
	if err != nil {
		t.Fatalf("Failed to marshal stop: %v", err)
	}

	// Unmarshal back to verify
	var roundTrip Stop
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("Failed to unmarshal stop: %v", err)
	}

	// Verify key fields
	if roundTrip.ID != stop.ID {
		t.Errorf("Expected stop ID '%s', got '%s'", stop.ID, roundTrip.ID)
	}
	if roundTrip.Attributes.Name != stop.Attributes.Name {
		t.Errorf("Expected name '%s', got '%s'", stop.Attributes.Name, roundTrip.Attributes.Name)
	}
	if roundTrip.Attributes.LocationType != stop.Attributes.LocationType {
		t.Errorf("Expected location_type %d, got %d", stop.Attributes.LocationType, roundTrip.Attributes.LocationType)
	}
}

func TestGetLocationTypeDescription(t *testing.T) {
	tests := []struct {
		locationType int
		expected     string
	}{
		{LocationTypePlatform, "Platform"},
		{LocationTypeStation, "Station"},
		{LocationTypeEntrance, "Entrance"},
		{LocationTypeGenericNode, "Generic Node"},
		{LocationTypeBoardingArea, "Boarding Area"},
		{100, "Unknown"},
	}

	for _, test := range tests {
		result := GetLocationTypeDescription(test.locationType)
		if result != test.expected {
			t.Errorf("For location type %d, expected '%s', got '%s'", test.locationType, test.expected, result)
		}
	}
}

func TestGetWheelchairBoardingDescription(t *testing.T) {
	tests := []struct {
		wheelchairBoarding int
		expected           string
	}{
		{WheelchairBoardingUnknown, "Unknown"},
		{WheelchairBoardingAccessible, "Accessible"},
		{WheelchairBoardingInaccessible, "Inaccessible"},
		{100, "Unknown"},
	}

	for _, test := range tests {
		result := GetWheelchairBoardingDescription(test.wheelchairBoarding)
		if result != test.expected {
			t.Errorf("For wheelchair boarding %d, expected '%s', got '%s'", test.wheelchairBoarding, test.expected, result)
		}
	}
}

func TestStop_IsAccessible(t *testing.T) {
	accessibleStop := Stop{
		Attributes: StopAttributes{
			WheelchairBoarding: WheelchairBoardingAccessible,
		},
	}

	inaccessibleStop := Stop{
		Attributes: StopAttributes{
			WheelchairBoarding: WheelchairBoardingInaccessible,
		},
	}

	unknownStop := Stop{
		Attributes: StopAttributes{
			WheelchairBoarding: WheelchairBoardingUnknown,
		},
	}

	if !accessibleStop.IsAccessible() {
		t.Error("Expected accessible stop to return true for IsAccessible()")
	}

	if inaccessibleStop.IsAccessible() {
		t.Error("Expected inaccessible stop to return false for IsAccessible()")
	}

	if unknownStop.IsAccessible() {
		t.Error("Expected unknown accessibility stop to return false for IsAccessible()")
	}
}

func TestStop_IsStation(t *testing.T) {
	station := Stop{
		Attributes: StopAttributes{
			LocationType: LocationTypeStation,
		},
	}

	platform := Stop{
		Attributes: StopAttributes{
			LocationType: LocationTypePlatform,
		},
	}

	if !station.IsStation() {
		t.Error("Expected station to return true for IsStation()")
	}

	if platform.IsStation() {
		t.Error("Expected platform to return false for IsStation()")
	}
}

func TestStop_IsPlatform(t *testing.T) {
	station := Stop{
		Attributes: StopAttributes{
			LocationType: LocationTypeStation,
		},
	}

	platform := Stop{
		Attributes: StopAttributes{
			LocationType: LocationTypePlatform,
		},
	}

	if station.IsPlatform() {
		t.Error("Expected station to return false for IsPlatform()")
	}

	if !platform.IsPlatform() {
		t.Error("Expected platform to return true for IsPlatform()")
	}
}
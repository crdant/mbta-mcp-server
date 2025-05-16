package models

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestScheduleUnmarshal(t *testing.T) {
	// Load the sample schedules data
	fixtureData, err := os.ReadFile("../../../test/fixtures/schedules.json")
	if err != nil {
		t.Fatalf("Failed to read test fixture: %v", err)
	}

	// Unmarshal the JSON
	var response ScheduleResponse
	if err := json.Unmarshal(fixtureData, &response); err != nil {
		t.Fatalf("Failed to unmarshal schedule response: %v", err)
	}

	// Verify that we have the correct number of schedules
	if len(response.Data) != 2 {
		t.Errorf("Expected 2 schedules, got %d", len(response.Data))
	}

	// Verify the first schedule
	schedule := response.Data[0]
	if schedule.ID != "schedule-1" {
		t.Errorf("Expected schedule ID 'schedule-1', got '%s'", schedule.ID)
	}
	if schedule.Type != "schedule" {
		t.Errorf("Expected schedule type 'schedule', got '%s'", schedule.Type)
	}

	attrs := schedule.Attributes
	// Parse the expected time
	expectedArrival, err := time.Parse(time.RFC3339, "2023-05-20T12:00:00-04:00")
	if err != nil {
		t.Fatalf("Failed to parse expected arrival time: %v", err)
	}
	expectedDeparture, err := time.Parse(time.RFC3339, "2023-05-20T12:02:00-04:00")
	if err != nil {
		t.Fatalf("Failed to parse expected departure time: %v", err)
	}

	// Check the parsed times are in the correct format
	arrivalTime, err := time.Parse(time.RFC3339, attrs.ArrivalTime)
	if err != nil {
		t.Errorf("Failed to parse arrival time '%s': %v", attrs.ArrivalTime, err)
	}
	if !arrivalTime.Equal(expectedArrival) {
		t.Errorf("Expected arrival time %v, got %v", expectedArrival, arrivalTime)
	}

	departureTime, err := time.Parse(time.RFC3339, attrs.DepartureTime)
	if err != nil {
		t.Errorf("Failed to parse departure time '%s': %v", attrs.DepartureTime, err)
	}
	if !departureTime.Equal(expectedDeparture) {
		t.Errorf("Expected departure time %v, got %v", expectedDeparture, departureTime)
	}

	if attrs.StopHeadsign != "Alewife" {
		t.Errorf("Expected stop_headsign 'Alewife', got '%s'", attrs.StopHeadsign)
	}
	if attrs.StopSequence != 1 {
		t.Errorf("Expected stop_sequence 1, got %d", attrs.StopSequence)
	}
	if !attrs.Timepoint {
		t.Errorf("Expected timepoint true, got %t", attrs.Timepoint)
	}

	// Check relationships
	// Route relationship
	if route, ok := schedule.Relationships["route"]; ok {
		if routeData, ok := route.(map[string]interface{})["data"]; ok {
			routeDataMap, ok := routeData.(map[string]interface{})
			if !ok {
				t.Errorf("Expected route data to be a map, got %T", routeData)
			} else if routeID, ok := routeDataMap["id"]; ok {
				if routeID != "Red" {
					t.Errorf("Expected route ID 'Red', got '%s'", routeID)
				}
			} else {
				t.Error("Missing route ID in relationship")
			}
		} else {
			t.Error("Missing data in route relationship")
		}
	} else {
		t.Error("Missing route relationship")
	}

	// Stop relationship
	if stop, ok := schedule.Relationships["stop"]; ok {
		if stopData, ok := stop.(map[string]interface{})["data"]; ok {
			stopDataMap, ok := stopData.(map[string]interface{})
			if !ok {
				t.Errorf("Expected stop data to be a map, got %T", stopData)
			} else if stopID, ok := stopDataMap["id"]; ok {
				if stopID != "place-sstat" {
					t.Errorf("Expected stop ID 'place-sstat', got '%s'", stopID)
				}
			} else {
				t.Error("Missing stop ID in relationship")
			}
		} else {
			t.Error("Missing data in stop relationship")
		}
	} else {
		t.Error("Missing stop relationship")
	}

	// Trip relationship
	if trip, ok := schedule.Relationships["trip"]; ok {
		if tripData, ok := trip.(map[string]interface{})["data"]; ok {
			tripDataMap, ok := tripData.(map[string]interface{})
			if !ok {
				t.Errorf("Expected trip data to be a map, got %T", tripData)
			} else if tripID, ok := tripDataMap["id"]; ok {
				if tripID != "Red-123456-20230520" {
					t.Errorf("Expected trip ID 'Red-123456-20230520', got '%s'", tripID)
				}
			} else {
				t.Error("Missing trip ID in relationship")
			}
		} else {
			t.Error("Missing data in trip relationship")
		}
	} else {
		t.Error("Missing trip relationship")
	}

	// Verify included trip
	if len(response.Included) != 1 {
		t.Errorf("Expected 1 included object, got %d", len(response.Included))
	}

	includedTrip := response.Included[0]
	if includedTrip.ID != "Red-123456-20230520" {
		t.Errorf("Expected included trip ID 'Red-123456-20230520', got '%s'", includedTrip.ID)
	}
	if includedTrip.Type != "trip" {
		t.Errorf("Expected included trip type 'trip', got '%s'", includedTrip.Type)
	}

	tripAttrs, ok := includedTrip.Attributes.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected trip attributes to be a map, got %T", includedTrip.Attributes)
	}

	if headsign, ok := tripAttrs["headsign"]; ok {
		if headsign != "Alewife" {
			t.Errorf("Expected trip headsign 'Alewife', got '%s'", headsign)
		}
	} else {
		t.Error("Missing headsign in trip attributes")
	}
}

func TestScheduleMarshal(t *testing.T) {
	// Create a schedule object
	schedule := Schedule{
		ID:   "schedule-1",
		Type: "schedule",
		Attributes: ScheduleAttributes{
			ArrivalTime:   "2023-05-20T12:00:00-04:00",
			DepartureTime: "2023-05-20T12:02:00-04:00",
			DropOffType:   0,
			PickupType:    0,
			StopHeadsign:  "Alewife",
			StopSequence:  1,
			Timepoint:     true,
		},
		Relationships: map[string]interface{}{
			"route": map[string]interface{}{
				"data": map[string]string{
					"id":   "Red",
					"type": "route",
				},
			},
			"stop": map[string]interface{}{
				"data": map[string]string{
					"id":   "place-sstat",
					"type": "stop",
				},
			},
			"trip": map[string]interface{}{
				"data": map[string]string{
					"id":   "Red-123456-20230520",
					"type": "trip",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(schedule)
	if err != nil {
		t.Fatalf("Failed to marshal schedule: %v", err)
	}

	// Unmarshal back to verify
	var roundTrip Schedule
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("Failed to unmarshal schedule: %v", err)
	}

	// Verify key fields
	if roundTrip.ID != schedule.ID {
		t.Errorf("Expected schedule ID '%s', got '%s'", schedule.ID, roundTrip.ID)
	}
	if roundTrip.Attributes.ArrivalTime != schedule.Attributes.ArrivalTime {
		t.Errorf("Expected arrival time '%s', got '%s'", schedule.Attributes.ArrivalTime, roundTrip.Attributes.ArrivalTime)
	}
	if roundTrip.Attributes.StopHeadsign != schedule.Attributes.StopHeadsign {
		t.Errorf("Expected stop_headsign '%s', got '%s'", schedule.Attributes.StopHeadsign, roundTrip.Attributes.StopHeadsign)
	}
}

func TestSchedule_GetDuration(t *testing.T) {
	tests := []struct {
		name           string
		arrivalTime    string
		departureTime  string
		expectedMinutes int
	}{
		{
			name:           "Two minute duration",
			arrivalTime:    "2023-05-20T12:00:00-04:00",
			departureTime:  "2023-05-20T12:02:00-04:00",
			expectedMinutes: 2,
		},
		{
			name:           "Ten minute duration",
			arrivalTime:    "2023-05-20T12:00:00-04:00",
			departureTime:  "2023-05-20T12:10:00-04:00",
			expectedMinutes: 10,
		},
		{
			name:           "Same time (zero duration)",
			arrivalTime:    "2023-05-20T12:00:00-04:00",
			departureTime:  "2023-05-20T12:00:00-04:00",
			expectedMinutes: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schedule := Schedule{
				Attributes: ScheduleAttributes{
					ArrivalTime:   test.arrivalTime,
					DepartureTime: test.departureTime,
				},
			}

			duration, err := schedule.GetDuration()
			if err != nil {
				t.Fatalf("Failed to get duration: %v", err)
			}

			minutes := int(duration.Minutes())
			if minutes != test.expectedMinutes {
				t.Errorf("Expected duration of %d minutes, got %d minutes", test.expectedMinutes, minutes)
			}
		})
	}
}

func TestSchedule_FormattedTimes(t *testing.T) {
	schedule := Schedule{
		Attributes: ScheduleAttributes{
			ArrivalTime:   "2023-05-20T12:00:00-04:00",
			DepartureTime: "2023-05-20T12:02:00-04:00",
		},
	}

	// Test formatted arrival time
	formattedArrival, err := schedule.FormattedArrivalTime("3:04 PM")
	if err != nil {
		t.Fatalf("Failed to format arrival time: %v", err)
	}
	if formattedArrival != "12:00 PM" {
		t.Errorf("Expected formatted arrival time '12:00 PM', got '%s'", formattedArrival)
	}

	// Test formatted departure time
	formattedDeparture, err := schedule.FormattedDepartureTime("3:04 PM")
	if err != nil {
		t.Fatalf("Failed to format departure time: %v", err)
	}
	if formattedDeparture != "12:02 PM" {
		t.Errorf("Expected formatted departure time '12:02 PM', got '%s'", formattedDeparture)
	}
}

func TestSchedule_PickupDropOff(t *testing.T) {
	tests := []struct {
		name         string
		pickupType   int
		dropOffType  int
		expectPickup bool
		expectDropOff bool
	}{
		{
			name:         "Regular service",
			pickupType:   PickupDropOffRegular,
			dropOffType:  PickupDropOffRegular,
			expectPickup: true,
			expectDropOff: true,
		},
		{
			name:         "No pickup",
			pickupType:   PickupDropOffNotAvailable,
			dropOffType:  PickupDropOffRegular,
			expectPickup: false,
			expectDropOff: true,
		},
		{
			name:         "No drop-off",
			pickupType:   PickupDropOffRegular,
			dropOffType:  PickupDropOffNotAvailable,
			expectPickup: true,
			expectDropOff: false,
		},
		{
			name:         "Phone agency for pickup",
			pickupType:   PickupDropOffPhoneAgency,
			dropOffType:  PickupDropOffRegular,
			expectPickup: false,
			expectDropOff: true,
		},
		{
			name:         "Coordinate with driver for drop-off",
			pickupType:   PickupDropOffRegular,
			dropOffType:  PickupDropOffCoordinateWithDriver,
			expectPickup: true,
			expectDropOff: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			schedule := Schedule{
				Attributes: ScheduleAttributes{
					PickupType:   test.pickupType,
					DropOffType:  test.dropOffType,
				},
			}

			if schedule.IsPickupAvailable() != test.expectPickup {
				t.Errorf("Expected IsPickupAvailable to return %v for pickup type %d",
					test.expectPickup, test.pickupType)
			}

			if schedule.IsDropOffAvailable() != test.expectDropOff {
				t.Errorf("Expected IsDropOffAvailable to return %v for drop-off type %d",
					test.expectDropOff, test.dropOffType)
			}
		})
	}
}

func TestSchedule_IsTimepoint(t *testing.T) {
	timepointSchedule := Schedule{
		Attributes: ScheduleAttributes{
			Timepoint: true,
		},
	}

	nonTimepointSchedule := Schedule{
		Attributes: ScheduleAttributes{
			Timepoint: false,
		},
	}

	if !timepointSchedule.IsTimepoint() {
		t.Error("Expected timepoint schedule to return true for IsTimepoint()")
	}

	if nonTimepointSchedule.IsTimepoint() {
		t.Error("Expected non-timepoint schedule to return false for IsTimepoint()")
	}
}
package models

import (
	"encoding/json"
	"testing"
)

func TestVehicleModel(t *testing.T) {
	// Test JSON unmarshaling of a vehicle response
	t.Run("Unmarshal vehicle response", func(t *testing.T) {
		vehicleJSON := `{
			"data": [
				{
					"attributes": {
						"bearing": 315.0,
						"current_status": "IN_TRANSIT_TO",
						"current_stop_sequence": 310,
						"direction_id": 0,
						"label": "3673-3838",
						"latitude": 42.33982849121094,
						"longitude": -71.15853881835938,
						"speed": null,
						"updated_at": "2023-05-20T11:22:55-04:00"
					},
					"id": "G-10067",
					"links": {
						"self": "/vehicles/G-10067"
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Green-B",
								"type": "route"
							}
						},
						"stop": {
							"data": {
								"id": "70107",
								"type": "stop"
							}
						},
						"trip": {
							"data": {
								"id": "36418064",
								"type": "trip"
							}
						}
					},
					"type": "vehicle"
				}
			]
		}`

		var response VehicleResponse
		err := json.Unmarshal([]byte(vehicleJSON), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal vehicle response: %v", err)
		}

		// Check that we got a vehicle
		if len(response.Data) != 1 {
			t.Fatalf("Expected 1 vehicle, got %d", len(response.Data))
		}

		vehicle := response.Data[0]

		// Check vehicle ID
		if vehicle.ID != "G-10067" {
			t.Errorf("Expected ID G-10067, got %s", vehicle.ID)
		}

		// Check vehicle type
		if vehicle.Type != "vehicle" {
			t.Errorf("Expected type 'vehicle', got %s", vehicle.Type)
		}

		// Check vehicle attributes
		attrs := vehicle.Attributes
		if attrs.Bearing != 315.0 {
			t.Errorf("Expected bearing 315.0, got %f", attrs.Bearing)
		}
		if attrs.CurrentStatus != "IN_TRANSIT_TO" {
			t.Errorf("Expected current_status 'IN_TRANSIT_TO', got %s", attrs.CurrentStatus)
		}
		if attrs.CurrentStopSequence != 310 {
			t.Errorf("Expected current_stop_sequence 310, got %d", attrs.CurrentStopSequence)
		}
		if attrs.DirectionID != 0 {
			t.Errorf("Expected direction_id 0, got %d", attrs.DirectionID)
		}
		if attrs.Label != "3673-3838" {
			t.Errorf("Expected label '3673-3838', got %s", attrs.Label)
		}
		if attrs.Latitude != 42.33982849121094 {
			t.Errorf("Expected latitude 42.33982849121094, got %f", attrs.Latitude)
		}
		if attrs.Longitude != -71.15853881835938 {
			t.Errorf("Expected longitude -71.15853881835938, got %f", attrs.Longitude)
		}
		if attrs.Speed != nil {
			t.Errorf("Expected speed nil, got %v", attrs.Speed)
		}
		if attrs.UpdatedAt != "2023-05-20T11:22:55-04:00" {
			t.Errorf("Expected updated_at '2023-05-20T11:22:55-04:00', got %s", attrs.UpdatedAt)
		}

		// Check relationships
		if route, ok := vehicle.Relationships["route"]; ok {
			data, ok := route.(map[string]interface{})["data"].(map[string]interface{})
			if !ok {
				t.Error("Failed to extract route relationship data")
			} else {
				id, ok := data["id"].(string)
				if !ok || id != "Green-B" {
					t.Errorf("Expected route ID 'Green-B', got %v", id)
				}
			}
		} else {
			t.Error("Route relationship not found")
		}

		if stop, ok := vehicle.Relationships["stop"]; ok {
			data, ok := stop.(map[string]interface{})["data"].(map[string]interface{})
			if !ok {
				t.Error("Failed to extract stop relationship data")
			} else {
				id, ok := data["id"].(string)
				if !ok || id != "70107" {
					t.Errorf("Expected stop ID '70107', got %v", id)
				}
			}
		} else {
			t.Error("Stop relationship not found")
		}

		if trip, ok := vehicle.Relationships["trip"]; ok {
			data, ok := trip.(map[string]interface{})["data"].(map[string]interface{})
			if !ok {
				t.Error("Failed to extract trip relationship data")
			} else {
				id, ok := data["id"].(string)
				if !ok || id != "36418064" {
					t.Errorf("Expected trip ID '36418064', got %v", id)
				}
			}
		} else {
			t.Error("Trip relationship not found")
		}
	})

	// Test vehicle with carriage data
	t.Run("Unmarshal vehicle with carriages", func(t *testing.T) {
		vehicleJSON := `{
			"data": [
				{
					"attributes": {
						"bearing": 125.0,
						"carriages": [
							{
								"label": "1861",
								"occupancy_status": "MANY_SEATS_AVAILABLE",
								"occupancy_percentage": 0
							},
							{
								"label": "1860",
								"occupancy_status": "MANY_SEATS_AVAILABLE",
								"occupancy_percentage": 0
							}
						],
						"current_status": "STOPPED_AT",
						"current_stop_sequence": 10,
						"direction_id": 0,
						"label": "1861-1860",
						"latitude": 42.39537,
						"longitude": -71.14236,
						"speed": 0,
						"updated_at": "2023-05-20T11:22:55-04:00"
					},
					"id": "R-5463D359",
					"links": {
						"self": "/vehicles/R-5463D359"
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Red",
								"type": "route"
							}
						},
						"stop": {
							"data": {
								"id": "70061",
								"type": "stop"
							}
						},
						"trip": {
							"data": {
								"id": "51070305",
								"type": "trip"
							}
						}
					},
					"type": "vehicle"
				}
			]
		}`

		var response VehicleResponse
		err := json.Unmarshal([]byte(vehicleJSON), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal vehicle response: %v", err)
		}

		// Check that we got a vehicle
		if len(response.Data) != 1 {
			t.Fatalf("Expected 1 vehicle, got %d", len(response.Data))
		}

		vehicle := response.Data[0]

		// Check carriage data
		if len(vehicle.Attributes.Carriages) != 2 {
			t.Fatalf("Expected 2 carriages, got %d", len(vehicle.Attributes.Carriages))
		}

		carriage1 := vehicle.Attributes.Carriages[0]
		if carriage1.Label != "1861" {
			t.Errorf("Expected carriage label '1861', got '%s'", carriage1.Label)
		}
		if carriage1.OccupancyStatus != "MANY_SEATS_AVAILABLE" {
			t.Errorf("Expected occupancy status 'MANY_SEATS_AVAILABLE', got '%s'", carriage1.OccupancyStatus)
		}
		if carriage1.OccupancyPercentage != 0 {
			t.Errorf("Expected occupancy percentage 0, got %d", carriage1.OccupancyPercentage)
		}

		carriage2 := vehicle.Attributes.Carriages[1]
		if carriage2.Label != "1860" {
			t.Errorf("Expected carriage label '1860', got '%s'", carriage2.Label)
		}
	})

	// Test vehicle status constants
	t.Run("Vehicle status constants", func(t *testing.T) {
		// Check that the constants match the expected values
		expectedStatuses := map[string]string{
			"INCOMING_AT":     "INCOMING_AT",
			"STOPPED_AT":      "STOPPED_AT",
			"IN_TRANSIT_TO":   "IN_TRANSIT_TO",
		}

		if VehicleStatusIncomingAt != expectedStatuses["INCOMING_AT"] {
			t.Errorf("Expected VehicleStatusIncomingAt to be '%s', got '%s'",
				expectedStatuses["INCOMING_AT"], VehicleStatusIncomingAt)
		}
		if VehicleStatusStoppedAt != expectedStatuses["STOPPED_AT"] {
			t.Errorf("Expected VehicleStatusStoppedAt to be '%s', got '%s'",
				expectedStatuses["STOPPED_AT"], VehicleStatusStoppedAt)
		}
		if VehicleStatusInTransitTo != expectedStatuses["IN_TRANSIT_TO"] {
			t.Errorf("Expected VehicleStatusInTransitTo to be '%s', got '%s'",
				expectedStatuses["IN_TRANSIT_TO"], VehicleStatusInTransitTo)
		}
	})

	// Test GetStatusDescription method
	t.Run("GetStatusDescription", func(t *testing.T) {
		v := Vehicle{
			Attributes: VehicleAttributes{
				CurrentStatus: "IN_TRANSIT_TO",
			},
		}

		expected := "In Transit"
		if v.GetStatusDescription() != expected {
			t.Errorf("Expected status description '%s', got '%s'",
				expected, v.GetStatusDescription())
		}

		v.Attributes.CurrentStatus = "STOPPED_AT"
		expected = "Stopped At"
		if v.GetStatusDescription() != expected {
			t.Errorf("Expected status description '%s', got '%s'",
				expected, v.GetStatusDescription())
		}

		v.Attributes.CurrentStatus = "INCOMING_AT"
		expected = "Arriving"
		if v.GetStatusDescription() != expected {
			t.Errorf("Expected status description '%s', got '%s'",
				expected, v.GetStatusDescription())
		}

		v.Attributes.CurrentStatus = "UNKNOWN"
		expected = "Unknown"
		if v.GetStatusDescription() != expected {
			t.Errorf("Expected status description '%s', got '%s'",
				expected, v.GetStatusDescription())
		}
	})
}
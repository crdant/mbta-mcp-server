package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPredictionModel(t *testing.T) {
	t.Run("Can parse prediction JSON", func(t *testing.T) {
		jsonData := `{
			"data": [
				{
					"id": "prediction-123",
					"type": "prediction",
					"attributes": {
						"arrival_time": "2025-06-01T14:30:00-04:00",
						"departure_time": "2025-06-01T14:32:00-04:00",
						"direction_id": 0,
						"schedule_relationship": "SCHEDULED",
						"status": null,
						"stop_sequence": 5,
						"track": "2"
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
								"id": "place-sstat",
								"type": "stop"
							}
						},
						"trip": {
							"data": {
								"id": "CR-Weekday-Fall-17-515",
								"type": "trip"
							}
						},
						"vehicle": {
							"data": {
								"id": "vehicle-123",
								"type": "vehicle"
							}
						}
					}
				}
			]
		}`

		var predictionResponse PredictionResponse
		err := json.Unmarshal([]byte(jsonData), &predictionResponse)

		if err != nil {
			t.Fatalf("Failed to parse prediction JSON: %v", err)
		}

		if len(predictionResponse.Data) != 1 {
			t.Fatalf("Expected 1 prediction, got %d", len(predictionResponse.Data))
		}

		prediction := predictionResponse.Data[0]
		if prediction.ID != "prediction-123" {
			t.Errorf("Expected prediction ID 'prediction-123', got '%s'", prediction.ID)
		}

		if prediction.Type != "prediction" {
			t.Errorf("Expected prediction type 'prediction', got '%s'", prediction.Type)
		}

		expectedArrivalTime := "2025-06-01T14:30:00-04:00"
		if *prediction.Attributes.ArrivalTime != expectedArrivalTime {
			t.Errorf("Expected arrival time '%s', got '%s'", expectedArrivalTime, *prediction.Attributes.ArrivalTime)
		}

		expectedDepartureTime := "2025-06-01T14:32:00-04:00"
		if *prediction.Attributes.DepartureTime != expectedDepartureTime {
			t.Errorf("Expected departure time '%s', got '%s'", expectedDepartureTime, *prediction.Attributes.DepartureTime)
		}

		if prediction.Attributes.Direction != 0 {
			t.Errorf("Expected direction 0, got %d", prediction.Attributes.Direction)
		}

		if prediction.Attributes.StopSequence != 5 {
			t.Errorf("Expected stop sequence 5, got %d", prediction.Attributes.StopSequence)
		}

		expectedTrack := "2"
		if *prediction.Attributes.Track != expectedTrack {
			t.Errorf("Expected track '%s', got '%s'", expectedTrack, *prediction.Attributes.Track)
		}
	})

	t.Run("Can extract relationship IDs", func(t *testing.T) {
		prediction := Prediction{
			ID:   "prediction-123",
			Type: "prediction",
			Relationships: map[string]interface{}{
				"route": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   "Red",
						"type": "route",
					},
				},
				"stop": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   "place-sstat",
						"type": "stop",
					},
				},
				"trip": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   "CR-Weekday-Fall-17-515",
						"type": "trip",
					},
				},
				"vehicle": map[string]interface{}{
					"data": map[string]interface{}{
						"id":   "vehicle-123",
						"type": "vehicle",
					},
				},
			},
		}

		routeID := prediction.GetRouteID()
		if routeID != "Red" {
			t.Errorf("Expected route ID 'Red', got '%s'", routeID)
		}

		stopID := prediction.GetStopID()
		if stopID != "place-sstat" {
			t.Errorf("Expected stop ID 'place-sstat', got '%s'", stopID)
		}

		tripID := prediction.GetTripID()
		if tripID != "CR-Weekday-Fall-17-515" {
			t.Errorf("Expected trip ID 'CR-Weekday-Fall-17-515', got '%s'", tripID)
		}

		vehicleID := prediction.GetVehicleID()
		if vehicleID != "vehicle-123" {
			t.Errorf("Expected vehicle ID 'vehicle-123', got '%s'", vehicleID)
		}
	})

	t.Run("Can parse time values", func(t *testing.T) {
		arrivalTimeStr := "2025-06-01T14:30:00Z"
		departureTimeStr := "2025-06-01T14:32:00Z"

		prediction := Prediction{
			ID:   "prediction-123",
			Type: "prediction",
			Attributes: PredictionAttributes{
				ArrivalTime:   &arrivalTimeStr,
				DepartureTime: &departureTimeStr,
			},
		}

		arrivalTime, err := prediction.GetArrivalTime()
		if err != nil {
			t.Errorf("Failed to parse arrival time: %v", err)
		}
		if arrivalTime == nil {
			t.Fatal("Expected non-nil arrival time")
		}

		expectedArrival, _ := time.Parse(time.RFC3339, arrivalTimeStr)
		if !arrivalTime.Equal(expectedArrival) {
			t.Errorf("Expected arrival time %v, got %v", expectedArrival, arrivalTime)
		}

		departureTime, err := prediction.GetDepartureTime()
		if err != nil {
			t.Errorf("Failed to parse departure time: %v", err)
		}
		if departureTime == nil {
			t.Fatal("Expected non-nil departure time")
		}

		expectedDeparture, _ := time.Parse(time.RFC3339, departureTimeStr)
		if !departureTime.Equal(expectedDeparture) {
			t.Errorf("Expected departure time %v, got %v", expectedDeparture, departureTime)
		}
	})

	t.Run("Handles nil time values", func(t *testing.T) {
		prediction := Prediction{
			ID:   "prediction-123",
			Type: "prediction",
			Attributes: PredictionAttributes{
				ArrivalTime:   nil,
				DepartureTime: nil,
			},
		}

		arrivalTime, err := prediction.GetArrivalTime()
		if err != nil {
			t.Errorf("Error should be nil for nil arrival time, got: %v", err)
		}
		if arrivalTime != nil {
			t.Errorf("Expected nil arrival time, got %v", arrivalTime)
		}

		departureTime, err := prediction.GetDepartureTime()
		if err != nil {
			t.Errorf("Error should be nil for nil departure time, got: %v", err)
		}
		if departureTime != nil {
			t.Errorf("Expected nil departure time, got %v", departureTime)
		}
	})
}

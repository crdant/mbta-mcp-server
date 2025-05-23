package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTripGetRouteID(t *testing.T) {
	// Create a trip with route relationship
	trip := Trip{
		ID:   "12345",
		Type: "trip",
		Relationships: map[string]interface{}{
			"route": map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "Red",
					"type": "route",
				},
			},
		},
	}

	// Test getting route ID
	routeID := trip.GetRouteID()
	if routeID != "Red" {
		t.Errorf("Expected route ID 'Red', got '%s'", routeID)
	}

	// Test with missing relationship
	tripNoRoute := Trip{
		ID:            "12345",
		Type:          "trip",
		Relationships: map[string]interface{}{},
	}

	routeID = tripNoRoute.GetRouteID()
	if routeID != "" {
		t.Errorf("Expected empty route ID for trip with no route relationship, got '%s'", routeID)
	}
}

func TestTripIsWheelchairAccessible(t *testing.T) {
	// Test with wheelchair accessible trip
	accessibleTrip := Trip{
		Attributes: TripAttributes{
			WheelchairEnabled: true,
		},
	}

	if !accessibleTrip.IsWheelchairAccessible() {
		t.Error("Expected trip to be wheelchair accessible, but it wasn't")
	}

	// Test with non-accessible trip
	nonAccessibleTrip := Trip{
		Attributes: TripAttributes{
			WheelchairEnabled: false,
		},
	}

	if nonAccessibleTrip.IsWheelchairAccessible() {
		t.Error("Expected trip to not be wheelchair accessible, but it was")
	}
}

func TestTripIsBikeAllowed(t *testing.T) {
	// Test with bike allowed trip
	bikeAllowedTrip := Trip{
		Attributes: TripAttributes{
			BikeAllowed: true,
		},
	}

	if !bikeAllowedTrip.IsBikeAllowed() {
		t.Error("Expected trip to allow bikes, but it didn't")
	}

	// Test with bike not allowed trip
	bikeNotAllowedTrip := Trip{
		Attributes: TripAttributes{
			BikeAllowed: false,
		},
	}

	if bikeNotAllowedTrip.IsBikeAllowed() {
		t.Error("Expected trip to not allow bikes, but it did")
	}
}

func TestTripPlanJSON(t *testing.T) {
	// Create a sample trip plan
	now := time.Now()
	later := now.Add(1 * time.Hour)

	plan := TripPlan{
		Origin: &Stop{
			ID: "origin-stop",
			Attributes: StopAttributes{
				Name: "Origin Station",
			},
		},
		Destination: &Stop{
			ID: "destination-stop",
			Attributes: StopAttributes{
				Name: "Destination Station",
			},
		},
		DepartureTime: now,
		ArrivalTime:   later,
		Duration:      1 * time.Hour,
		Legs: []TripLeg{
			{
				Origin: &Stop{
					ID: "origin-stop",
					Attributes: StopAttributes{
						Name: "Origin Station",
					},
				},
				Destination: &Stop{
					ID: "destination-stop",
					Attributes: StopAttributes{
						Name: "Destination Station",
					},
				},
				RouteID:       "Red",
				RouteName:     "Red Line",
				TripID:        "trip-1234",
				DepartureTime: now,
				ArrivalTime:   later,
				Duration:      1 * time.Hour,
				Distance:      5.5,
				Headsign:      "Alewife",
				DirectionID:   0,
				IsAccessible:  true,
				Instructions:  "Board the Red Line toward Alewife",
			},
		},
		TotalDistance:  5.5,
		AccessibleTrip: true,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to serialize trip plan to JSON: %v", err)
	}

	// Deserialize back to a trip plan
	var deserializedPlan TripPlan
	err = json.Unmarshal(jsonData, &deserializedPlan)
	if err != nil {
		t.Fatalf("Failed to deserialize trip plan from JSON: %v", err)
	}

	// Verify basic properties
	if deserializedPlan.Origin.ID != "origin-stop" {
		t.Errorf("Expected origin ID 'origin-stop', got '%s'", deserializedPlan.Origin.ID)
	}

	if deserializedPlan.Destination.ID != "destination-stop" {
		t.Errorf("Expected destination ID 'destination-stop', got '%s'", deserializedPlan.Destination.ID)
	}

	if len(deserializedPlan.Legs) != 1 {
		t.Fatalf("Expected 1 trip leg, got %d", len(deserializedPlan.Legs))
	}

	leg := deserializedPlan.Legs[0]
	if leg.RouteID != "Red" {
		t.Errorf("Expected route ID 'Red', got '%s'", leg.RouteID)
	}

	if leg.RouteName != "Red Line" {
		t.Errorf("Expected route name 'Red Line', got '%s'", leg.RouteName)
	}

	if !leg.IsAccessible {
		t.Error("Expected trip leg to be accessible, but it wasn't")
	}

	if leg.Instructions != "Board the Red Line toward Alewife" {
		t.Errorf("Expected instructions 'Board the Red Line toward Alewife', got '%s'", leg.Instructions)
	}
}

func TestTransferPointJSON(t *testing.T) {
	// Create a sample transfer point
	transfer := TransferPoint{
		Stop: &Stop{
			ID: "transfer-stop",
			Attributes: StopAttributes{
				Name: "Transfer Station",
			},
		},
		FromRoute:       "Red",
		ToRoute:         "Green",
		TransferType:    TransferTypeRecommended,
		MinTransferTime: 5 * time.Minute,
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(transfer)
	if err != nil {
		t.Fatalf("Failed to serialize transfer point to JSON: %v", err)
	}

	// Deserialize back to a transfer point
	var deserializedTransfer TransferPoint
	err = json.Unmarshal(jsonData, &deserializedTransfer)
	if err != nil {
		t.Fatalf("Failed to deserialize transfer point from JSON: %v", err)
	}

	// Verify properties
	if deserializedTransfer.Stop.ID != "transfer-stop" {
		t.Errorf("Expected stop ID 'transfer-stop', got '%s'", deserializedTransfer.Stop.ID)
	}

	if deserializedTransfer.FromRoute != "Red" {
		t.Errorf("Expected from route 'Red', got '%s'", deserializedTransfer.FromRoute)
	}

	if deserializedTransfer.ToRoute != "Green" {
		t.Errorf("Expected to route 'Green', got '%s'", deserializedTransfer.ToRoute)
	}

	if deserializedTransfer.TransferType != TransferTypeRecommended {
		t.Errorf("Expected transfer type %d, got %d", TransferTypeRecommended, deserializedTransfer.TransferType)
	}

	if deserializedTransfer.MinTransferTime != 5*time.Minute {
		t.Errorf("Expected min transfer time 5m0s, got %v", deserializedTransfer.MinTransferTime)
	}
}

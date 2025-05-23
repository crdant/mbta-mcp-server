// Package models contains data models for MBTA API responses
package models

import "time"

// TripResponse represents a response containing trip data from the MBTA API
type TripResponse struct {
	Data []Trip `json:"data"`
}

// Trip represents a transit trip in the MBTA system
type Trip struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    TripAttributes         `json:"attributes"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// TripAttributes contains the attributes of a trip
type TripAttributes struct {
	Name              string  `json:"name"`
	Headsign          string  `json:"headsign"`
	Direction         int     `json:"direction_id"`
	BlockID           *string `json:"block_id"`
	ServiceID         string  `json:"service_id"`
	ShapeID           *string `json:"shape_id"`
	WheelchairEnabled bool    `json:"wheelchair_accessible"`
	BikeAllowed       bool    `json:"bikes_allowed"`
}

// TripPlan represents a complete trip plan with multiple legs
type TripPlan struct {
	Origin         *Stop         `json:"origin"`
	Destination    *Stop         `json:"destination"`
	DepartureTime  time.Time     `json:"departure_time"`
	ArrivalTime    time.Time     `json:"arrival_time"`
	Duration       time.Duration `json:"duration"`
	Legs           []TripLeg     `json:"legs"`
	TotalDistance  float64       `json:"total_distance"`
	AccessibleTrip bool          `json:"accessible_trip"`
}

// TripLeg represents a single leg of a trip plan
type TripLeg struct {
	Origin         *Stop         `json:"origin"`
	Destination    *Stop         `json:"destination"`
	RouteID        string        `json:"route_id"`
	RouteName      string        `json:"route_name"`
	TripID         string        `json:"trip_id"`
	DepartureTime  time.Time     `json:"departure_time"`
	ArrivalTime    time.Time     `json:"arrival_time"`
	Duration       time.Duration `json:"duration"`
	Distance       float64       `json:"distance"`
	Headsign       string        `json:"headsign"`
	DirectionID    int           `json:"direction_id"`
	Stops          []Stop        `json:"intermediate_stops,omitempty"`
	IsAccessible   bool          `json:"is_accessible"`
	Instructions   string        `json:"instructions"`
	PredictedTimes *PredictedTimes `json:"predicted_times,omitempty"`
}

// PredictedTimes contains real-time predictions for a trip leg
type PredictedTimes struct {
	PredictedDeparture time.Time `json:"predicted_departure"`
	PredictedArrival   time.Time `json:"predicted_arrival"`
	IsDelayed          bool      `json:"is_delayed"`
	DelayMinutes       int       `json:"delay_minutes"`
}

// TransferPoint represents a transfer between two trip legs
type TransferPoint struct {
	Stop              *Stop         `json:"stop"`
	FromRoute         string        `json:"from_route"`
	ToRoute           string        `json:"to_route"`
	TransferType      int           `json:"transfer_type"`
	MinTransferTime   time.Duration `json:"min_transfer_time"`
	SuggestedWaitTime time.Duration `json:"suggested_wait_time,omitempty"`
}

// Transfer type constants
const (
	TransferTypeRecommended = 0 // Recommended transfer point
	TransferTypeTimed       = 1 // Timed transfer with guaranteed connection
	TransferTypeMinimum     = 2 // Minimum transfer time needed
	TransferTypeNotPossible = 3 // Transfer not possible
)

// GetRouteID extracts the route ID from the trip's relationships
func (t *Trip) GetRouteID() string {
	if route, ok := t.Relationships["route"]; ok {
		if data, ok := route.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetServiceID returns the service ID for this trip
func (t *Trip) GetServiceID() string {
	return t.Attributes.ServiceID
}

// IsWheelchairAccessible returns whether the trip is accessible to wheelchairs
func (t *Trip) IsWheelchairAccessible() bool {
	return t.Attributes.WheelchairEnabled
}

// IsBikeAllowed returns whether bikes are allowed on this trip
func (t *Trip) IsBikeAllowed() bool {
	return t.Attributes.BikeAllowed
}
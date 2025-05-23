// Package models contains data models for MBTA API responses
package models

import "time"

// PredictionResponse represents a response containing prediction data from the MBTA API
type PredictionResponse struct {
	Data  []Prediction           `json:"data"`
	Links map[string]interface{} `json:"links,omitempty"`
}

// Prediction represents arrival and departure predictions for transit vehicles
type Prediction struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    PredictionAttributes   `json:"attributes"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// PredictionAttributes contains the detailed prediction data
type PredictionAttributes struct {
	ArrivalTime   *string `json:"arrival_time"`
	DepartureTime *string `json:"departure_time"`
	Direction     int     `json:"direction_id"`
	Schedule      string  `json:"schedule_relationship"`
	Status        *string `json:"status"`
	StopSequence  int     `json:"stop_sequence"`
	Track         *string `json:"track"`
}

// GetRouteID extracts the route ID from the prediction's relationships
func (p *Prediction) GetRouteID() string {
	if route, ok := p.Relationships["route"]; ok {
		if data, ok := route.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetStopID extracts the stop ID from the prediction's relationships
func (p *Prediction) GetStopID() string {
	if stop, ok := p.Relationships["stop"]; ok {
		if data, ok := stop.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetTripID extracts the trip ID from the prediction's relationships
func (p *Prediction) GetTripID() string {
	if trip, ok := p.Relationships["trip"]; ok {
		if data, ok := trip.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetVehicleID extracts the vehicle ID from the prediction's relationships
func (p *Prediction) GetVehicleID() string {
	if vehicle, ok := p.Relationships["vehicle"]; ok {
		if data, ok := vehicle.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetArrivalTime parses the arrival time string into a time.Time
func (p *Prediction) GetArrivalTime() (*time.Time, error) {
	if p.Attributes.ArrivalTime == nil {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, *p.Attributes.ArrivalTime)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetDepartureTime parses the departure time string into a time.Time
func (p *Prediction) GetDepartureTime() (*time.Time, error) {
	if p.Attributes.DepartureTime == nil {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, *p.Attributes.DepartureTime)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTimeUntilArrival returns the duration until the predicted arrival
func (p *Prediction) GetTimeUntilArrival() (*time.Duration, error) {
	arrivalTime, err := p.GetArrivalTime()
	if err != nil || arrivalTime == nil {
		return nil, err
	}

	duration := time.Until(*arrivalTime)
	return &duration, nil
}

// GetTimeUntilDeparture returns the duration until the predicted departure
func (p *Prediction) GetTimeUntilDeparture() (*time.Duration, error) {
	departureTime, err := p.GetDepartureTime()
	if err != nil || departureTime == nil {
		return nil, err
	}

	duration := time.Until(*departureTime)
	return &duration, nil
}

// Package models contains data models for MBTA API responses
package models

// VehicleResponse represents a response containing vehicle data from the MBTA API
type VehicleResponse struct {
	Data []Vehicle `json:"data"`
}

// Vehicle represents a transit vehicle in the MBTA system
type Vehicle struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    VehicleAttributes      `json:"attributes"`
	Links         map[string]string      `json:"links,omitempty"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// VehicleAttributes contains the attributes of a vehicle
type VehicleAttributes struct {
	Bearing             float64            `json:"bearing"`
	Carriages           []VehicleCarriage  `json:"carriages,omitempty"`
	CurrentStatus       string             `json:"current_status"`
	CurrentStopSequence int                `json:"current_stop_sequence"`
	DirectionID         int                `json:"direction_id"`
	Label               string             `json:"label"`
	Latitude            float64            `json:"latitude"`
	Longitude           float64            `json:"longitude"`
	Speed               *float64           `json:"speed"`
	UpdatedAt           string             `json:"updated_at"`
}

// VehicleCarriage represents individual carriage data for multi-car vehicles
type VehicleCarriage struct {
	Label               string `json:"label"`
	OccupancyStatus     string `json:"occupancy_status"`
	OccupancyPercentage int    `json:"occupancy_percentage"`
}

// Vehicle status constants as defined by the MBTA API
const (
	VehicleStatusIncomingAt   = "INCOMING_AT"
	VehicleStatusStoppedAt    = "STOPPED_AT"
	VehicleStatusInTransitTo  = "IN_TRANSIT_TO"
)

// GetStatusDescription returns a human-readable description of the vehicle's current status
func (v *Vehicle) GetStatusDescription() string {
	switch v.Attributes.CurrentStatus {
	case VehicleStatusIncomingAt:
		return "Arriving"
	case VehicleStatusStoppedAt:
		return "Stopped At"
	case VehicleStatusInTransitTo:
		return "In Transit"
	default:
		return "Unknown"
	}
}

// GetRouteID extracts the route ID from the vehicle's relationships
func (v *Vehicle) GetRouteID() string {
	if route, ok := v.Relationships["route"]; ok {
		if data, ok := route.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetStopID extracts the stop ID from the vehicle's relationships
func (v *Vehicle) GetStopID() string {
	if stop, ok := v.Relationships["stop"]; ok {
		if data, ok := stop.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// GetTripID extracts the trip ID from the vehicle's relationships
func (v *Vehicle) GetTripID() string {
	if trip, ok := v.Relationships["trip"]; ok {
		if data, ok := trip.(map[string]interface{})["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				return id
			}
		}
	}
	return ""
}

// HasOccupancyData returns true if the vehicle has occupancy data for its carriages
func (v *Vehicle) HasOccupancyData() bool {
	return len(v.Attributes.Carriages) > 0
}
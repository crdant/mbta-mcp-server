// Package models contains data models for MBTA API responses
package models

// StopResponse represents a response containing stop data from the MBTA API
type StopResponse struct {
	Data []Stop `json:"data"`
}

// Stop represents a transit stop or station in the MBTA system
type Stop struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    StopAttributes         `json:"attributes"`
	Links         map[string]string      `json:"links,omitempty"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// StopAttributes contains the attributes of a stop
type StopAttributes struct {
	Address            string  `json:"address"`
	AtStreet           *string `json:"at_street"`
	Description        string  `json:"description"`
	Latitude           float64 `json:"latitude"`
	LocationType       int     `json:"location_type"`
	Longitude          float64 `json:"longitude"`
	Municipality       string  `json:"municipality"`
	Name               string  `json:"name"`
	OnStreet           *string `json:"on_street"`
	PlatformCode       string  `json:"platform_code"`
	PlatformName       string  `json:"platform_name"`
	VehicleType        *int    `json:"vehicle_type"`
	WheelchairBoarding int     `json:"wheelchair_boarding"`
}

// Location type constants as defined by the MBTA API
const (
	LocationTypePlatform     = 0
	LocationTypeStation      = 1
	LocationTypeEntrance     = 2
	LocationTypeGenericNode  = 3
	LocationTypeBoardingArea = 4
)

// Wheelchair boarding constants as defined by the MBTA API
const (
	WheelchairBoardingUnknown      = 0
	WheelchairBoardingAccessible   = 1
	WheelchairBoardingInaccessible = 2
)

// GetLocationTypeDescription returns a string description for a location type value
func GetLocationTypeDescription(locationType int) string {
	switch locationType {
	case LocationTypePlatform:
		return "Platform"
	case LocationTypeStation:
		return "Station"
	case LocationTypeEntrance:
		return "Entrance"
	case LocationTypeGenericNode:
		return "Generic Node"
	case LocationTypeBoardingArea:
		return "Boarding Area"
	default:
		return "Unknown"
	}
}

// GetWheelchairBoardingDescription returns a string description for a wheelchair boarding value
func GetWheelchairBoardingDescription(wheelchairBoarding int) string {
	switch wheelchairBoarding {
	case WheelchairBoardingUnknown:
		return "Unknown"
	case WheelchairBoardingAccessible:
		return "Accessible"
	case WheelchairBoardingInaccessible:
		return "Inaccessible"
	default:
		return "Unknown"
	}
}

// IsAccessible returns whether the stop is accessible to wheelchairs
func (s *Stop) IsAccessible() bool {
	return s.Attributes.WheelchairBoarding == WheelchairBoardingAccessible
}

// IsStation returns whether the stop is a station (location_type = 1)
func (s *Stop) IsStation() bool {
	return s.Attributes.LocationType == LocationTypeStation
}

// IsPlatform returns whether the stop is a platform (location_type = 0)
func (s *Stop) IsPlatform() bool {
	return s.Attributes.LocationType == LocationTypePlatform
}
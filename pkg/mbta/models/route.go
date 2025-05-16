// Package models contains data models for MBTA API responses
package models

// RouteResponse represents a response containing route data from the MBTA API
type RouteResponse struct {
	Data []Route `json:"data"`
}

// Route represents a transit route in the MBTA system
type Route struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    RouteAttributes        `json:"attributes"`
	Links         map[string]string      `json:"links,omitempty"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// RouteAttributes contains the attributes of a route
type RouteAttributes struct {
	Color                 string   `json:"color"`
	Description           string   `json:"description"`
	DirectionDestinations []string `json:"direction_destinations"`
	DirectionNames        []string `json:"direction_names"`
	FareClass             string   `json:"fare_class"`
	LongName              string   `json:"long_name"`
	ShortName             string   `json:"short_name"`
	SortOrder             int      `json:"sort_order"`
	TextColor             string   `json:"text_color"`
	Type                  int      `json:"type"`
}

// Route type constants as defined by the MBTA API
const (
	RouteTypeLightRail    = 0
	RouteTypeSubway       = 1
	RouteTypeCommuterRail = 2
	RouteTypeBus          = 3
	RouteTypeFerry        = 4
)

// GetRouteTypeDescription returns a string description for a route type value
func GetRouteTypeDescription(routeType int) string {
	switch routeType {
	case RouteTypeLightRail:
		return "Light Rail"
	case RouteTypeSubway:
		return "Subway"
	case RouteTypeCommuterRail:
		return "Commuter Rail"
	case RouteTypeBus:
		return "Bus"
	case RouteTypeFerry:
		return "Ferry"
	default:
		return "Unknown"
	}
}

// GetTypeDescription returns a string description for this route's type
func (r *Route) GetTypeDescription() string {
	return GetRouteTypeDescription(r.Attributes.Type)
}

// GetDirectionName returns the name for a given direction (0 = outbound, 1 = inbound)
func (r *Route) GetDirectionName(direction int) string {
	if direction >= 0 && direction < len(r.Attributes.DirectionNames) {
		return r.Attributes.DirectionNames[direction]
	}
	return ""
}

// GetDirectionDestination returns the destination for a given direction (0 = outbound, 1 = inbound)
func (r *Route) GetDirectionDestination(direction int) string {
	if direction >= 0 && direction < len(r.Attributes.DirectionDestinations) {
		return r.Attributes.DirectionDestinations[direction]
	}
	return ""
}
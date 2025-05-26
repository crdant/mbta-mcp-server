// ABOUTME: This file implements the station proximity handlers for the MCP server.
// ABOUTME: It defines tools and request processing logic for finding nearby MBTA stations.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/crdant/mbta-mcp-server/pkg/mbta"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// registerGeographicQueryTools registers the geographic query tools and handlers
func (s *Server) registerGeographicQueryTools() {
	// Tool: FindNearbyStations - finds stations near the specified coordinates
	findNearbyStationsTool := mcp.Tool{
		Name:        "find_nearby_stations",
		Description: "Find MBTA stations near the specified coordinates",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"latitude": map[string]any{
					"type":        "number",
					"description": "Latitude coordinate to search around",
				},
				"longitude": map[string]any{
					"type":        "number",
					"description": "Longitude coordinate to search around",
				},
				"radius": map[string]any{
					"type":        "number",
					"description": "Maximum distance in kilometers (default: 1.0)",
				},
				"max_results": map[string]any{
					"type":        "number",
					"description": "Maximum number of results to return (default: 5)",
				},
				"only_stations": map[string]any{
					"type":        "boolean",
					"description": "If true, only return full stations, not platforms or stops (default: true)",
				},
				"wheelchair_accessible": map[string]any{
					"type":        "boolean",
					"description": "If true, only return wheelchair accessible stations",
				},
			},
			Required: []string{"latitude", "longitude"},
		},
	}

	// Register the station proximity tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(findNearbyStationsTool, s.wrapWithMiddleware(s.findNearbyStationsHandler))
}

// findNearbyStationsHandler handles requests for finding stations near coordinates
func (s *Server) findNearbyStationsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for nearby stations: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract required parameters
	args := request.Params.Arguments
	
	// Get latitude
	latitude, ok := args["latitude"].(float64)
	if !ok {
		// Try to convert from other types
		if latStr, ok := args["latitude"].(string); ok {
			var err error
			if latitude, err = parseFloat(latStr); err != nil {
				return createErrorResponse("Invalid latitude parameter, must be a number"), nil
			}
		} else {
			return createErrorResponse("Missing or invalid latitude parameter"), nil
		}
	}

	// Get longitude
	longitude, ok := args["longitude"].(float64)
	if !ok {
		// Try to convert from other types
		if lonStr, ok := args["longitude"].(string); ok {
			var err error
			if longitude, err = parseFloat(lonStr); err != nil {
				return createErrorResponse("Invalid longitude parameter, must be a number"), nil
			}
		} else {
			return createErrorResponse("Missing or invalid longitude parameter"), nil
		}
	}

	// Extract optional parameters with defaults
	radius := 1.0 // Default: 1 km
	if radiusVal, ok := args["radius"].(float64); ok {
		radius = radiusVal
	}

	maxResults := 5 // Default: 5 results
	if maxResultsVal, ok := args["max_results"].(float64); ok {
		maxResults = int(maxResultsVal)
	}

	onlyStations := true // Default: only full stations
	if onlyStationsVal, ok := args["only_stations"].(bool); ok {
		onlyStations = onlyStationsVal
	}

	wheelchairAccessible := false // Default: include all stations
	if wheelchairAccessibleVal, ok := args["wheelchair_accessible"].(bool); ok {
		wheelchairAccessible = wheelchairAccessibleVal
	}

	log.Printf("Searching for stations near (%f, %f) within %f km", latitude, longitude, radius)

	// Find nearby stations
	nearbyStations, err := client.FindNearbyStations(ctx, latitude, longitude, radius, maxResults, onlyStations)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to find nearby stations: %v", err)), nil
	}

	// Filter for wheelchair accessibility if requested
	if wheelchairAccessible && len(nearbyStations) > 0 {
		accessibleStations := make([]models.NearbyStation, 0, len(nearbyStations))
		for _, station := range nearbyStations {
			if station.Stop.IsAccessible() {
				accessibleStations = append(accessibleStations, station)
			}
		}
		nearbyStations = accessibleStations
	}

	// If no stations are found, inform the user
	if len(nearbyStations) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("No stations found within %.1f km of the specified coordinates.", radius),
				},
			},
		}, nil
	}

	log.Printf("Found %d stations near the specified coordinates", len(nearbyStations))

	// Format the stations for response
	return formatNearbyStationsResponse(nearbyStations)
}

// formatNearbyStationsResponse converts nearby station data to a proper MCP response
func formatNearbyStationsResponse(stations []models.NearbyStation) (*mcp.CallToolResult, error) {
	// Convert the stations to a simplified format for the response
	stationsData := make([]map[string]interface{}, 0, len(stations))

	for _, station := range stations {
		stationMap := map[string]interface{}{
			"id":                   station.Stop.ID,
			"name":                 station.Stop.Attributes.Name,
			"distance_km":          station.DistanceKm,
			"distance_miles":       station.DistanceKm * 0.621371, // Convert to miles
			"latitude":             station.Stop.Attributes.Latitude,
			"longitude":            station.Stop.Attributes.Longitude,
			"municipality":         station.Stop.Attributes.Municipality,
			"location_type":        station.Stop.Attributes.LocationType,
			"location_description": models.GetLocationTypeDescription(station.Stop.Attributes.LocationType),
			"wheelchair_accessible": station.Stop.IsAccessible(),
			"wheelchair_boarding":   station.Stop.Attributes.WheelchairBoarding,
			"accessibility_status":  models.GetWheelchairBoardingDescription(station.Stop.Attributes.WheelchairBoarding),
		}

		// Add address if available
		if station.Stop.Attributes.Address != "" {
			stationMap["address"] = station.Stop.Attributes.Address
		}

		// Add formatted distance
		if station.DistanceKm < 1.0 {
			stationMap["distance_text"] = fmt.Sprintf("%.0f meters", station.DistanceKm*1000)
		} else {
			stationMap["distance_text"] = fmt.Sprintf("%.1f km (%.1f miles)", station.DistanceKm, station.DistanceKm*0.621371)
		}

		stationsData = append(stationsData, stationMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(stationsData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize station data: %v", err)), nil
	}

	// Return data as a text content item
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonBytes),
			},
		},
	}, nil
}

// parseFloat converts a string to a float64
func parseFloat(s string) (float64, error) {
	return json.Number(s).Float64()
}
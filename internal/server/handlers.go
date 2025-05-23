// ABOUTME: This file implements the request handlers for the MCP server.
// ABOUTME: It defines tools and request processing logic for MBTA transit information.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/crdant/mbta-mcp-server/pkg/mbta"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// RegisterDefaultHandlers sets up the default tools and handlers for the MCP server.
// This registers all the standard MBTA transit information tools.
func (s *Server) RegisterDefaultHandlers() {
	// Set up transit information tools
	s.registerTransitInfoTools()
}

// registerTransitInfoTools registers the basic transit information tools.
func (s *Server) registerTransitInfoTools() {
	// Tool: GetRoutes - retrieves MBTA routes information
	getRoutesTool := mcp.Tool{
		Name:        "get_routes",
		Description: "Get available MBTA routes",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_type": map[string]any{
					"type":        "string",
					"description": "Filter by route type (0=Light Rail, 1=Subway, 2=Commuter Rail, 3=Bus, etc.)",
				},
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter by specific route ID",
				},
			},
		},
	}

	// Register the routes tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getRoutesTool, s.wrapWithMiddleware(s.getRoutesHandler))

	// Tool: GetStops - retrieves MBTA stops information
	getStopsTool := mcp.Tool{
		Name:        "get_stops",
		Description: "Get available MBTA stops and stations",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"stop_id": map[string]any{
					"type":        "string",
					"description": "Filter by specific stop ID",
				},
				"location_type": map[string]any{
					"type":        "string",
					"description": "Filter by location type (0=Stop/Platform, 1=Station, 2=Entrance, etc.)",
				},
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter stops by route ID",
				},
			},
		},
	}

	// Register the stops tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getStopsTool, s.wrapWithMiddleware(s.getStopsHandler))

	// Tool: GetSchedules - retrieves MBTA schedule information
	getSchedulesTool := mcp.Tool{
		Name:        "get_schedules",
		Description: "Get MBTA schedules for routes and stops",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter schedules by route ID",
				},
				"stop_id": map[string]any{
					"type":        "string",
					"description": "Filter schedules by stop ID",
				},
				"direction_id": map[string]any{
					"type":        "string",
					"description": "Filter by direction (0=outbound, 1=inbound)",
				},
			},
		},
	}

	// Register the schedules tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getSchedulesTool, s.wrapWithMiddleware(s.getSchedulesHandler))
}

// getRoutesHandler handles requests for MBTA route information.
// It connects to the MBTA API client to retrieve and filter route data.
func (s *Server) getRoutesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for routes: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract optional parameters for filtering
	args := request.Params.Arguments
	routeType, hasRouteType := args["route_type"]
	routeID, hasRouteID := args["route_id"]

	// Log the request details
	if hasRouteType {
		log.Printf("Filtering by route type: %v", routeType)
	}
	if hasRouteID {
		log.Printf("Filtering by route ID: %v", routeID)
	}

	// If filtering by specific route ID, use GetRoute instead of GetRoutes
	if hasRouteID {
		routeIDStr, ok := routeID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid route_id parameter: %v", routeID)), nil
		}

		route, err := client.GetRoute(ctx, routeIDStr)
		if err != nil {
			return createErrorResponse(fmt.Sprintf("Failed to retrieve route %s: %v", routeIDStr, err)), nil
		}

		// If also filtering by route type, check if it matches
		if hasRouteType {
			routeTypeStr, ok := routeType.(string)
			if !ok {
				return createErrorResponse(fmt.Sprintf("Invalid route_type parameter: %v", routeType)), nil
			}

			routeTypeInt, err := strconv.Atoi(routeTypeStr)
			if err != nil {
				return createErrorResponse(fmt.Sprintf("Invalid route_type format: %v", routeType)), nil
			}

			if route.Attributes.Type != routeTypeInt {
				// Route type doesn't match the filter
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("No routes found matching ID %s and type %s", routeIDStr, routeTypeStr),
						},
					},
				}, nil
			}
		}

		// Convert the single route to a formatted response
		return formatRouteResponse([]*models.Route{route})
	}

	// Get all routes
	routes, err := client.GetRoutes(ctx)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve routes: %v", err)), nil
	}

	// Filter by route type if specified
	if hasRouteType {
		routeTypeStr, ok := routeType.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid route_type parameter: %v", routeType)), nil
		}

		routeTypeInt, err := strconv.Atoi(routeTypeStr)
		if err != nil {
			return createErrorResponse(fmt.Sprintf("Invalid route_type format: %v", routeType)), nil
		}

		filteredRoutes := make([]models.Route, 0)
		for _, route := range routes {
			if route.Attributes.Type == routeTypeInt {
				filteredRoutes = append(filteredRoutes, route)
			}
		}

		routes = filteredRoutes
	}

	// Convert slice of value types to slice of pointer types for formatting
	routePtrs := make([]*models.Route, len(routes))
	for i := range routes {
		routePtrs[i] = &routes[i]
	}

	return formatRouteResponse(routePtrs)
}

// getStopsHandler handles requests for MBTA stop information.
// It connects to the MBTA API client to retrieve and filter stop data.
func (s *Server) getStopsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for stops: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract optional parameters for filtering
	args := request.Params.Arguments
	stopID, hasStopID := args["stop_id"]
	locationType, hasLocationType := args["location_type"]
	routeID, hasRouteID := args["route_id"]

	// Log the request details
	if hasStopID {
		log.Printf("Filtering by stop ID: %v", stopID)
	}
	if hasLocationType {
		log.Printf("Filtering by location type: %v", locationType)
	}
	if hasRouteID {
		log.Printf("Filtering by route ID: %v", routeID)
	}

	// If filtering by specific stop ID, use GetStop instead of GetStops
	if hasStopID {
		stopIDStr, ok := stopID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid stop_id parameter: %v", stopID)), nil
		}

		stop, err := client.GetStop(ctx, stopIDStr)
		if err != nil {
			return createErrorResponse(fmt.Sprintf("Failed to retrieve stop %s: %v", stopIDStr, err)), nil
		}

		// If also filtering by location type, check if it matches
		if hasLocationType {
			locationTypeStr, ok := locationType.(string)
			if !ok {
				return createErrorResponse(fmt.Sprintf("Invalid location_type parameter: %v", locationType)), nil
			}

			locationTypeInt, err := strconv.Atoi(locationTypeStr)
			if err != nil {
				return createErrorResponse(fmt.Sprintf("Invalid location_type format: %v", locationType)), nil
			}

			if stop.Attributes.LocationType != locationTypeInt {
				// Location type doesn't match the filter
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						mcp.TextContent{
							Type: "text",
							Text: fmt.Sprintf("No stops found matching ID %s and location type %s", stopIDStr, locationTypeStr),
						},
					},
				}, nil
			}
		}

		// For now, return a text response (in future, this would be proper structured data)
		return formatStopResponse([]*models.Stop{stop})
	}

	// Get all stops
	stops, err := client.GetStops(ctx)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve stops: %v", err)), nil
	}

	// Filter by location type if specified
	if hasLocationType {
		locationTypeStr, ok := locationType.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid location_type parameter: %v", locationType)), nil
		}

		locationTypeInt, err := strconv.Atoi(locationTypeStr)
		if err != nil {
			return createErrorResponse(fmt.Sprintf("Invalid location_type format: %v", locationType)), nil
		}

		filteredStops := make([]models.Stop, 0)
		for _, stop := range stops {
			if stop.Attributes.LocationType == locationTypeInt {
				filteredStops = append(filteredStops, stop)
			}
		}

		stops = filteredStops
	}

	// TODO: Add filtering by route ID (would require additional API call)
	if hasRouteID {
		// For now, just note that this feature is not yet implemented
		log.Printf("Route ID filtering not implemented yet")
	}

	// Convert slice of value types to slice of pointer types for formatting
	stopPtrs := make([]*models.Stop, len(stops))
	for i := range stops {
		stopPtrs[i] = &stops[i]
	}

	return formatStopResponse(stopPtrs)
}

// getSchedulesHandler handles requests for MBTA schedule information.
// It connects to the MBTA API client to retrieve and filter schedule data.
func (s *Server) getSchedulesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for schedules: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract optional parameters for filtering
	args := request.Params.Arguments
	routeID, hasRouteID := args["route_id"]
	stopID, hasStopID := args["stop_id"]
	directionID, hasDirectionID := args["direction_id"]

	// Log the request details
	if hasRouteID {
		log.Printf("Filtering by route ID: %v", routeID)
	}
	if hasStopID {
		log.Printf("Filtering by stop ID: %v", stopID)
	}
	if hasDirectionID {
		log.Printf("Filtering by direction ID: %v", directionID)
	}

	// Build query parameters
	params := make(map[string]string)
	
	if hasRouteID {
		routeIDStr, ok := routeID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid route_id parameter: %v", routeID)), nil
		}
		params["filter[route]"] = routeIDStr
	}
	
	if hasStopID {
		stopIDStr, ok := stopID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid stop_id parameter: %v", stopID)), nil
		}
		params["filter[stop]"] = stopIDStr
	}
	
	if hasDirectionID {
		directionIDStr, ok := directionID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid direction_id parameter: %v", directionID)), nil
		}
		params["filter[direction_id]"] = directionIDStr
	}

	// Get schedules
	schedules, included, err := client.GetSchedules(ctx, params)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve schedules: %v", err)), nil
	}

	return formatScheduleResponse(schedules, included)
}

// formatRouteResponse converts route data to a proper MCP response
func formatRouteResponse(routes []*models.Route) (*mcp.CallToolResult, error) {
	// Convert the routes to a structured format
	routesData := make([]map[string]interface{}, 0, len(routes))
	for _, route := range routes {
		routeMap := map[string]interface{}{
			"id":                   route.ID,
			"name":                 route.Attributes.LongName,
			"short_name":           route.Attributes.ShortName,
			"type":                 route.Attributes.Type,
			"type_description":     route.GetTypeDescription(),
			"description":          route.Attributes.Description,
			"color":                route.Attributes.Color,
			"text_color":           route.Attributes.TextColor,
			"directions":           route.Attributes.DirectionNames,
			"direction_destinations": route.Attributes.DirectionDestinations,
		}
		routesData = append(routesData, routeMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(routesData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize route data: %v", err)), nil
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

// formatStopResponse converts stop data to a proper MCP response
func formatStopResponse(stops []*models.Stop) (*mcp.CallToolResult, error) {
	// Convert the stops to a structured format
	stopsData := make([]map[string]interface{}, 0, len(stops))
	for _, stop := range stops {
		stopMap := map[string]interface{}{
			"id":                   stop.ID,
			"name":                 stop.Attributes.Name,
			"description":          stop.Attributes.Description,
			"location_type":        stop.Attributes.LocationType,
			"location_description": models.GetLocationTypeDescription(stop.Attributes.LocationType),
			"municipality":         stop.Attributes.Municipality,
			"latitude":             stop.Attributes.Latitude,
			"longitude":            stop.Attributes.Longitude,
			"wheelchair_boarding":  stop.Attributes.WheelchairBoarding,
			"is_accessible":        stop.IsAccessible(),
		}

		// Add optional fields if they exist
		if stop.Attributes.PlatformCode != "" {
			stopMap["platform_code"] = stop.Attributes.PlatformCode
		}
		if stop.Attributes.PlatformName != "" {
			stopMap["platform_name"] = stop.Attributes.PlatformName
		}

		stopsData = append(stopsData, stopMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(stopsData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize stop data: %v", err)), nil
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

// formatScheduleResponse converts schedule data to a proper MCP response
func formatScheduleResponse(schedules []models.Schedule, included []models.Included) (*mcp.CallToolResult, error) {
	// Convert the schedules to a structured format
	schedulesData := make([]map[string]interface{}, 0, len(schedules))
	for _, schedule := range schedules {
		// Format arrival and departure times for better readability
		arrivalTime, _ := schedule.FormattedArrivalTime("3:04 PM")
		departureTime, _ := schedule.FormattedDepartureTime("3:04 PM")

		scheduleMap := map[string]interface{}{
			"id":              schedule.ID,
			"arrival_time":    schedule.Attributes.ArrivalTime,
			"departure_time":  schedule.Attributes.DepartureTime,
			"formatted_arrival":    arrivalTime,
			"formatted_departure":  departureTime,
			"stop_sequence":   schedule.Attributes.StopSequence,
			"stop_headsign":   schedule.Attributes.StopHeadsign,
			"pickup_available": schedule.IsPickupAvailable(),
			"dropoff_available": schedule.IsDropOffAvailable(),
			"is_timepoint":    schedule.IsTimepoint(),
		}

		// Extract relationship IDs
		if routeData, ok := schedule.Relationships["route"]; ok {
			if routeMap, ok := routeData.(map[string]interface{}); ok {
				if dataMap, ok := routeMap["data"].(map[string]interface{}); ok {
					if routeID, ok := dataMap["id"].(string); ok {
						scheduleMap["route_id"] = routeID
					}
				}
			}
		}

		if stopData, ok := schedule.Relationships["stop"]; ok {
			if stopMap, ok := stopData.(map[string]interface{}); ok {
				if dataMap, ok := stopMap["data"].(map[string]interface{}); ok {
					if stopID, ok := dataMap["id"].(string); ok {
						scheduleMap["stop_id"] = stopID
					}
				}
			}
		}

		schedulesData = append(schedulesData, scheduleMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(schedulesData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize schedule data: %v", err)), nil
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

// createErrorResponse creates a standardized error response for MCP requests.
func createErrorResponse(message string) *mcp.CallToolResult {
	errorContent := mcp.TextContent{
		Type: "text",
		Text: message,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{errorContent},
		IsError: true,
	}
}
// ABOUTME: This file implements the vehicle tracking handlers for the MCP server.
// ABOUTME: It defines tools and request processing logic for MBTA vehicle tracking.

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

// registerVehicleTrackingTools registers the vehicle tracking tools and handlers
func (s *Server) registerVehicleTrackingTools() {
	// Tool: GetVehicles - retrieves MBTA vehicle information
	getVehiclesTool := mcp.Tool{
		Name:        "get_vehicles",
		Description: "Get MBTA vehicle locations and status information",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter vehicles by route ID",
				},
				"trip_id": map[string]any{
					"type":        "string",
					"description": "Filter vehicles by trip ID",
				},
				"latitude": map[string]any{
					"type":        "number",
					"description": "Latitude for location-based filtering",
				},
				"longitude": map[string]any{
					"type":        "number",
					"description": "Longitude for location-based filtering",
				},
				"radius": map[string]any{
					"type":        "number",
					"description": "Radius in degrees for location-based filtering (default: 0.01)",
				},
			},
		},
	}

	// Register the vehicles tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getVehiclesTool, s.wrapWithMiddleware(s.getVehiclesHandler))

	// Tool: GetVehicle - retrieves information for a specific vehicle
	getVehicleTool := mcp.Tool{
		Name:        "get_vehicle",
		Description: "Get information for a specific MBTA vehicle",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"vehicle_id": map[string]any{
					"type":        "string",
					"description": "The ID of the vehicle to retrieve",
				},
			},
			Required: []string{"vehicle_id"},
		},
	}

	// Register the vehicle tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getVehicleTool, s.wrapWithMiddleware(s.getVehicleHandler))

	// Tool: GetVehiclePredictions - retrieves arrival predictions for a vehicle
	getVehiclePredictionsTool := mcp.Tool{
		Name:        "get_vehicle_predictions",
		Description: "Get arrival predictions for a specific MBTA vehicle",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"vehicle_id": map[string]any{
					"type":        "string",
					"description": "The ID of the vehicle to get predictions for",
				},
			},
			Required: []string{"vehicle_id"},
		},
	}

	// Register the predictions tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getVehiclePredictionsTool, s.wrapWithMiddleware(s.getVehiclePredictionsHandler))

	// Tool: GetVehicleStatus - retrieves real-time status updates for transit vehicles
	getVehicleStatusTool := mcp.Tool{
		Name:        "get_vehicle_status",
		Description: "Get real-time status updates for MBTA transit vehicles",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter status updates by route ID",
				},
				"status_type": map[string]any{
					"type":        "string",
					"description": "Type of status update to retrieve (arriving, stopped, in_transit)",
					"enum":        []string{"arriving", "stopped", "in_transit", "all"},
				},
				"limit": map[string]any{
					"type":        "number",
					"description": "Maximum number of status updates to return (default: 10)",
				},
			},
		},
	}

	// Register the status tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getVehicleStatusTool, s.wrapWithMiddleware(s.getVehicleStatusHandler))
}

// getVehiclesHandler handles requests for MBTA vehicle information
func (s *Server) getVehiclesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for vehicles: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract parameters for filtering
	args := request.Params.Arguments
	params := make(map[string]string)

	// Process route_id filter
	if routeID, ok := args["route_id"]; ok {
		routeIDStr, ok := routeID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid route_id parameter: %v", routeID)), nil
		}
		log.Printf("Filtering vehicles by route ID: %s", routeIDStr)
		params["filter[route]"] = routeIDStr
	}

	// Process trip_id filter
	if tripID, ok := args["trip_id"]; ok {
		tripIDStr, ok := tripID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid trip_id parameter: %v", tripID)), nil
		}
		log.Printf("Filtering vehicles by trip ID: %s", tripIDStr)
		params["filter[trip]"] = tripIDStr
	}

	// Process location-based filtering
	lat, hasLat := args["latitude"]
	lon, hasLon := args["longitude"]
	if hasLat && hasLon {
		// Extract latitude
		latitude, ok := lat.(float64)
		if !ok {
			// Try to convert from string
			latStr, ok := lat.(string)
			if ok {
				latFloat, err := strconv.ParseFloat(latStr, 64)
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Invalid latitude parameter: %v", lat)), nil
				}
				latitude = latFloat
			} else {
				return createErrorResponse(fmt.Sprintf("Invalid latitude parameter: %v", lat)), nil
			}
		}

		// Extract longitude
		longitude, ok := lon.(float64)
		if !ok {
			// Try to convert from string
			lonStr, ok := lon.(string)
			if ok {
				lonFloat, err := strconv.ParseFloat(lonStr, 64)
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Invalid longitude parameter: %v", lon)), nil
				}
				longitude = lonFloat
			} else {
				return createErrorResponse(fmt.Sprintf("Invalid longitude parameter: %v", lon)), nil
			}
		}

		// Extract radius if provided
		radius := 0.01 // default radius
		if rad, ok := args["radius"]; ok {
			radius, ok = rad.(float64)
			if !ok {
				// Try to convert from string
				radStr, ok := rad.(string)
				if ok {
					radFloat, err := strconv.ParseFloat(radStr, 64)
					if err != nil {
						return createErrorResponse(fmt.Sprintf("Invalid radius parameter: %v", rad)), nil
					}
					radius = radFloat
				} else {
					return createErrorResponse(fmt.Sprintf("Invalid radius parameter: %v", rad)), nil
				}
			}
		}

		log.Printf("Filtering vehicles by location: lat=%f, lon=%f, radius=%f", latitude, longitude, radius)
		params["filter[latitude]"] = fmt.Sprintf("%f", latitude)
		params["filter[longitude]"] = fmt.Sprintf("%f", longitude)
		params["filter[radius]"] = fmt.Sprintf("%f", radius)
	} else if hasLat || hasLon {
		// If only one coordinate is provided, return an error
		return createErrorResponse("Both latitude and longitude must be provided for location-based filtering"), nil
	}

	// Request vehicles with the specified filters
	vehicles, err := client.GetVehicles(ctx, params)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve vehicles: %v", err)), nil
	}

	log.Printf("Retrieved %d vehicles", len(vehicles))

	// Format the vehicles for response
	return formatVehiclesResponse(vehicles)
}

// getVehicleHandler handles requests for a specific MBTA vehicle
func (s *Server) getVehicleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for vehicle: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract vehicle ID parameter
	args := request.Params.Arguments
	vehicleID, ok := args["vehicle_id"]
	if !ok {
		return createErrorResponse("Missing required parameter: vehicle_id"), nil
	}

	vehicleIDStr, ok := vehicleID.(string)
	if !ok {
		return createErrorResponse(fmt.Sprintf("Invalid vehicle_id parameter: %v", vehicleID)), nil
	}

	log.Printf("Retrieving vehicle with ID: %s", vehicleIDStr)

	// Get the specified vehicle
	vehicle, err := client.GetVehicle(ctx, vehicleIDStr)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve vehicle %s: %v", vehicleIDStr, err)), nil
	}

	// Format the vehicle for response
	return formatVehicleResponse(vehicle)
}

// getVehiclePredictionsHandler handles requests for predictions for a specific vehicle
func (s *Server) getVehiclePredictionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for vehicle predictions: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract vehicle ID parameter
	args := request.Params.Arguments
	vehicleID, ok := args["vehicle_id"]
	if !ok {
		return createErrorResponse("Missing required parameter: vehicle_id"), nil
	}

	vehicleIDStr, ok := vehicleID.(string)
	if !ok {
		return createErrorResponse(fmt.Sprintf("Invalid vehicle_id parameter: %v", vehicleID)), nil
	}

	log.Printf("Retrieving predictions for vehicle ID: %s", vehicleIDStr)

	// Get predictions for the specified vehicle
	predictions, err := client.GetPredictionsByVehicle(ctx, vehicleIDStr)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve predictions for vehicle %s: %v", vehicleIDStr, err)), nil
	}

	// If no predictions are found, inform the user
	if len(predictions) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("No predictions available for vehicle %s.", vehicleIDStr),
				},
			},
		}, nil
	}

	// Format the predictions for response
	return formatPredictionsResponse(predictions)
}

// formatVehiclesResponse converts vehicle data to a proper MCP response
func formatVehiclesResponse(vehicles []models.Vehicle) (*mcp.CallToolResult, error) {
	// Convert the vehicles to a simplified format for the response
	vehiclesData := make([]map[string]interface{}, 0, len(vehicles))
	for _, vehicle := range vehicles {
		vehicleMap := map[string]interface{}{
			"id":           vehicle.ID,
			"label":        vehicle.Attributes.Label,
			"status":       vehicle.GetStatusDescription(),
			"latitude":     vehicle.Attributes.Latitude,
			"longitude":    vehicle.Attributes.Longitude,
			"bearing":      vehicle.Attributes.Bearing,
			"updated_at":   vehicle.Attributes.UpdatedAt,
			"route_id":     vehicle.GetRouteID(),
			"stop_id":      vehicle.GetStopID(),
			"trip_id":      vehicle.GetTripID(),
			"direction_id": vehicle.Attributes.DirectionID,
		}

		// Add speed if available
		if vehicle.Attributes.Speed != nil {
			vehicleMap["speed"] = *vehicle.Attributes.Speed
		}

		// Add occupancy data if available
		if vehicle.HasOccupancyData() {
			carriagesData := make([]map[string]interface{}, 0, len(vehicle.Attributes.Carriages))
			for _, carriage := range vehicle.Attributes.Carriages {
				carriageMap := map[string]interface{}{
					"label":                carriage.Label,
					"occupancy_status":     carriage.OccupancyStatus,
					"occupancy_percentage": carriage.OccupancyPercentage,
				}
				carriagesData = append(carriagesData, carriageMap)
			}
			vehicleMap["carriages"] = carriagesData
		}

		vehiclesData = append(vehiclesData, vehicleMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(vehiclesData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize vehicle data: %v", err)), nil
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

// formatVehicleResponse converts a single vehicle to a proper MCP response
func formatVehicleResponse(vehicle *models.Vehicle) (*mcp.CallToolResult, error) {
	// Convert the vehicle to a simplified format for the response
	vehicleMap := map[string]interface{}{
		"id":           vehicle.ID,
		"label":        vehicle.Attributes.Label,
		"status":       vehicle.GetStatusDescription(),
		"latitude":     vehicle.Attributes.Latitude,
		"longitude":    vehicle.Attributes.Longitude,
		"bearing":      vehicle.Attributes.Bearing,
		"updated_at":   vehicle.Attributes.UpdatedAt,
		"route_id":     vehicle.GetRouteID(),
		"stop_id":      vehicle.GetStopID(),
		"trip_id":      vehicle.GetTripID(),
		"direction_id": vehicle.Attributes.DirectionID,
	}

	// Add speed if available
	if vehicle.Attributes.Speed != nil {
		vehicleMap["speed"] = *vehicle.Attributes.Speed
	}

	// Add occupancy data if available
	if vehicle.HasOccupancyData() {
		carriagesData := make([]map[string]interface{}, 0, len(vehicle.Attributes.Carriages))
		for _, carriage := range vehicle.Attributes.Carriages {
			carriageMap := map[string]interface{}{
				"label":                carriage.Label,
				"occupancy_status":     carriage.OccupancyStatus,
				"occupancy_percentage": carriage.OccupancyPercentage,
			}
			carriagesData = append(carriagesData, carriageMap)
		}
		vehicleMap["carriages"] = carriagesData
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(vehicleMap, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize vehicle data: %v", err)), nil
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

// formatPredictionsResponse converts prediction data to a proper MCP response
func formatPredictionsResponse(predictions []models.Prediction) (*mcp.CallToolResult, error) {
	// Convert the predictions to a simplified format for the response
	predictionsData := make([]map[string]interface{}, 0, len(predictions))

	for _, prediction := range predictions {
		// Calculate time until arrival/departure
		timeUntilArrival, arrivalErr := prediction.GetTimeUntilArrival()
		timeUntilDeparture, departureErr := prediction.GetTimeUntilDeparture()

		// Base prediction data
		predictionMap := map[string]interface{}{
			"id":           prediction.ID,
			"route_id":     prediction.GetRouteID(),
			"stop_id":      prediction.GetStopID(),
			"trip_id":      prediction.GetTripID(),
			"vehicle_id":   prediction.GetVehicleID(),
			"direction_id": prediction.Attributes.Direction,
			"status":       prediction.Attributes.Status,
			"schedule":     prediction.Attributes.Schedule,
		}

		// Add arrival time if available
		if prediction.Attributes.ArrivalTime != nil {
			predictionMap["arrival_time"] = *prediction.Attributes.ArrivalTime

			// Add minutes until arrival if calculation was successful
			if arrivalErr == nil && timeUntilArrival != nil {
				minutesUntilArrival := timeUntilArrival.Minutes()
				predictionMap["minutes_until_arrival"] = minutesUntilArrival
			}
		}

		// Add departure time if available
		if prediction.Attributes.DepartureTime != nil {
			predictionMap["departure_time"] = *prediction.Attributes.DepartureTime

			// Add minutes until departure if calculation was successful
			if departureErr == nil && timeUntilDeparture != nil {
				minutesUntilDeparture := timeUntilDeparture.Minutes()
				predictionMap["minutes_until_departure"] = minutesUntilDeparture
			}
		}

		// Add track information if available
		if prediction.Attributes.Track != nil {
			predictionMap["track"] = *prediction.Attributes.Track
		}

		predictionsData = append(predictionsData, predictionMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(predictionsData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize prediction data: %v", err)), nil
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

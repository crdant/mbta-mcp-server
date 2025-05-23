// ABOUTME: This file implements the vehicle status update functionality for the MCP server.
// ABOUTME: It provides real-time status updates for transit vehicles.

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

// getVehicleStatusHandler handles requests for real-time status updates for transit vehicles
func (s *Server) getVehicleStatusHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for vehicle status updates: %s", request.Params.Name)

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
		log.Printf("Filtering status updates by route ID: %s", routeIDStr)
		params["filter[route]"] = routeIDStr
	}

	// Process status_type filter
	var statusFilter string
	if statusType, ok := args["status_type"]; ok {
		statusTypeStr, ok := statusType.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid status_type parameter: %v", statusType)), nil
		}
		
		// Convert from user-friendly name to API status value
		switch statusTypeStr {
		case "arriving":
			statusFilter = models.VehicleStatusIncomingAt
		case "stopped":
			statusFilter = models.VehicleStatusStoppedAt
		case "in_transit":
			statusFilter = models.VehicleStatusInTransitTo
		case "all":
			statusFilter = ""
		default:
			return createErrorResponse(fmt.Sprintf("Invalid status_type value: %s. Must be 'arriving', 'stopped', 'in_transit', or 'all'", statusTypeStr)), nil
		}
		
		if statusFilter != "" {
			log.Printf("Filtering status updates by status type: %s", statusFilter)
			params["filter[current_status]"] = statusFilter
		}
	}

	// Process limit parameter
	var limit int = 10 // default limit
	if limitParam, ok := args["limit"]; ok {
		limitInt, ok := limitParam.(float64)
		if !ok {
			// Try to convert from string
			limitStr, ok := limitParam.(string)
			if ok {
				parsedLimit, err := strconv.Atoi(limitStr)
				if err != nil {
					return createErrorResponse(fmt.Sprintf("Invalid limit parameter: %v", limitParam)), nil
				}
				limit = parsedLimit
			} else {
				return createErrorResponse(fmt.Sprintf("Invalid limit parameter: %v", limitParam)), nil
			}
		} else {
			limit = int(limitInt)
		}
		
		if limit < 1 {
			limit = 1
		} else if limit > 50 {
			limit = 50
		}
	}

	// Get vehicles with the specified filters
	vehicles, err := client.GetVehicles(ctx, params)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve vehicle status updates: %v", err)), nil
	}

	// If no vehicles are found, inform the user
	if len(vehicles) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No vehicle status updates found matching the specified criteria.",
				},
			},
		}, nil
	}

	// Limit the number of vehicles if needed
	if len(vehicles) > limit {
		vehicles = vehicles[:limit]
	}

	log.Printf("Retrieved %d vehicle status updates", len(vehicles))

	// Format the status updates for response
	return formatVehicleStatusResponse(vehicles)
}

// formatVehicleStatusResponse formats vehicle status data for the MCP response
func formatVehicleStatusResponse(vehicles []models.Vehicle) (*mcp.CallToolResult, error) {
	// Convert the vehicles to a status update format for the response
	statusUpdates := make([]map[string]interface{}, 0, len(vehicles))
	
	for _, vehicle := range vehicles {
		// Base status data
		statusUpdate := map[string]interface{}{
			"vehicle_id":    vehicle.ID,
			"label":         vehicle.Attributes.Label,
			"status":        vehicle.GetStatusDescription(),
			"status_code":   vehicle.Attributes.CurrentStatus,
			"latitude":      vehicle.Attributes.Latitude,
			"longitude":     vehicle.Attributes.Longitude,
			"route_id":      vehicle.GetRouteID(),
			"stop_id":       vehicle.GetStopID(),
			"trip_id":       vehicle.GetTripID(),
			"direction_id":  vehicle.Attributes.DirectionID,
			"updated_at":    vehicle.Attributes.UpdatedAt,
			"is_moving":     vehicle.Attributes.CurrentStatus == models.VehicleStatusInTransitTo,
		}
		
		// Add bearing and speed if moving
		if vehicle.Attributes.CurrentStatus == models.VehicleStatusInTransitTo {
			statusUpdate["bearing"] = vehicle.Attributes.Bearing
			
			if vehicle.Attributes.Speed != nil {
				statusUpdate["speed"] = *vehicle.Attributes.Speed
			}
		}
		
		// Add occupancy information if available
		if vehicle.HasOccupancyData() {
			totalOccupancy := 0
			for _, carriage := range vehicle.Attributes.Carriages {
				totalOccupancy += carriage.OccupancyPercentage
			}
			
			// Calculate average occupancy
			averageOccupancy := totalOccupancy / len(vehicle.Attributes.Carriages)
			statusUpdate["occupancy_percentage"] = averageOccupancy
			
			// Determine overall occupancy status
			var overallStatus string
			if averageOccupancy < 25 {
				overallStatus = "MANY_SEATS_AVAILABLE"
			} else if averageOccupancy < 50 {
				overallStatus = "SEATS_AVAILABLE"
			} else if averageOccupancy < 75 {
				overallStatus = "FEW_SEATS_AVAILABLE"
			} else {
				overallStatus = "FULL"
			}
			statusUpdate["occupancy_status"] = overallStatus
		}
		
		statusUpdates = append(statusUpdates, statusUpdate)
	}
	
	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(statusUpdates, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize status data: %v", err)), nil
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
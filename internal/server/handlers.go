// ABOUTME: This file implements the request handlers for the MCP server.
// ABOUTME: It defines tools and request processing logic for MBTA transit information.

package server

import (
	"context"
	"log"

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
	// Example tool: GetRoutes - retrieves MBTA routes information
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

	// Register the tool with its handler
	s.mcpServer.AddTool(getRoutesTool, s.getRoutesHandler)

	// Future tools will be registered here for additional functionality
}

// getRoutesHandler handles requests for MBTA route information.
// This is a placeholder that will be connected to the MBTA API client.
func (s *Server) getRoutesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for routes: %s", request.Params.Name)

	// TODO: Implement actual MBTA API client call
	// This is a placeholder response
	// These mock routes would be converted to actual data in the final implementation
	_ = []map[string]interface{}{
		{
			"id":   "Red",
			"name": "Red Line",
			"type": "1", // Subway
		},
		{
			"id":   "Green-B",
			"name": "Green Line B",
			"type": "0", // Light Rail
		},
	}

	// Extract optional parameters for potential filtering
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

	// In a real implementation, we would filter based on these parameters
	// For now, return the mock data regardless of filters

	// Create the text content for our response
	routesContent := mcp.TextContent{
		Type: "text",
		Text: "Routes data has been retrieved successfully.",
	}

	// Return data as a text content item
	return &mcp.CallToolResult{
		Content: []mcp.Content{routesContent},
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
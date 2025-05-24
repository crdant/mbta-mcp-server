// ABOUTME: This file implements the service alert handlers for the MCP server.
// ABOUTME: It defines tools and request processing logic for MBTA service alerts.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/crdant/mbta-mcp-server/pkg/mbta"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// registerServiceAlertTools registers the service alert tools and handlers
func (s *Server) registerServiceAlertTools() {
	// Tool: GetAlerts - retrieves MBTA service alerts
	getAlertsTool := mcp.Tool{
		Name:        "get_alerts",
		Description: "Get MBTA service alerts and disruptions",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter alerts by route ID",
				},
				"stop_id": map[string]any{
					"type":        "string",
					"description": "Filter alerts by stop ID",
				},
				"effect": map[string]any{
					"type":        "string",
					"description": "Filter alerts by effect (e.g., DELAYS, DETOUR, SHUTTLE, NO_SERVICE, STATION_CLOSURE)",
				},
				"active_only": map[string]any{
					"type":        "boolean",
					"description": "Only include currently active alerts",
				},
				"activity": map[string]any{
					"type":        "string",
					"description": "Filter alerts by activity (e.g., BOARD, EXIT, RIDE, USING_WHEELCHAIR)",
				},
			},
		},
	}

	// Register the alerts tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getAlertsTool, s.wrapWithMiddleware(s.getAlertsHandler))

	// Tool: GetServiceDisruptions - retrieves significant service disruptions
	getDisruptionsTool := mcp.Tool{
		Name:        "get_service_disruptions",
		Description: "Get significant MBTA service disruptions like service suspensions, reduced service, and significant delays",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"route_id": map[string]any{
					"type":        "string",
					"description": "Filter disruptions by route ID",
				},
				"severity_min": map[string]any{
					"type":        "number",
					"description": "Minimum severity level (1-9, where 9 is most severe)",
				},
			},
		},
	}

	// Register the disruptions tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getDisruptionsTool, s.wrapWithMiddleware(s.getServiceDisruptionsHandler))

	// Tool: GetAccessibilityAlerts - retrieves accessibility-related alerts
	getAccessibilityAlertsTool := mcp.Tool{
		Name:        "get_accessibility_alerts",
		Description: "Get MBTA alerts related to accessibility issues like elevator outages",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"stop_id": map[string]any{
					"type":        "string",
					"description": "Filter accessibility alerts by stop ID",
				},
			},
		},
	}

	// Register the accessibility alerts tool with its handler, wrapped with middleware
	s.mcpServer.AddTool(getAccessibilityAlertsTool, s.wrapWithMiddleware(s.getAccessibilityAlertsHandler))
}

// getAlertsHandler handles requests for MBTA service alert information
func (s *Server) getAlertsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for service alerts: %s", request.Params.Name)

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
		log.Printf("Filtering alerts by route ID: %s", routeIDStr)
		params["filter[route]"] = routeIDStr
	}

	// Process stop_id filter
	if stopID, ok := args["stop_id"]; ok {
		stopIDStr, ok := stopID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid stop_id parameter: %v", stopID)), nil
		}
		log.Printf("Filtering alerts by stop ID: %s", stopIDStr)
		params["filter[stop]"] = stopIDStr
	}

	// Process effect filter
	if effect, ok := args["effect"]; ok {
		effectStr, ok := effect.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid effect parameter: %v", effect)), nil
		}
		log.Printf("Filtering alerts by effect: %s", effectStr)
		params["filter[effect]"] = effectStr
	}

	// Process activity filter
	if activity, ok := args["activity"]; ok {
		activityStr, ok := activity.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid activity parameter: %v", activity)), nil
		}
		log.Printf("Filtering alerts by activity: %s", activityStr)
		params["filter[activity]"] = activityStr
	}

	// Get alerts with the specified filters
	alerts, err := client.GetAlerts(ctx, params)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve alerts: %v", err)), nil
	}

	// Check for active_only filter
	activeOnly := false
	if activeOnlyParam, ok := args["active_only"]; ok {
		activeOnly, _ = activeOnlyParam.(bool)
	}

	// Filter alerts if active_only is true
	if activeOnly {
		now := time.Now()
		activeAlerts := make([]models.Alert, 0, len(alerts))
		for _, alert := range alerts {
			if alert.IsActive(now) {
				activeAlerts = append(activeAlerts, alert)
			}
		}
		alerts = activeAlerts
	}

	// If no alerts are found, inform the user
	if len(alerts) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No alerts found matching the specified criteria.",
				},
			},
		}, nil
	}

	log.Printf("Retrieved %d alerts", len(alerts))

	// Format the alerts for response
	return formatAlertsResponse(alerts)
}

// getServiceDisruptionsHandler handles requests for significant service disruptions
func (s *Server) getServiceDisruptionsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for service disruptions: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract parameters for filtering
	args := request.Params.Arguments

	// Get service disruptions
	disruptions, err := client.GetServiceDisruptions(ctx)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve service disruptions: %v", err)), nil
	}

	// Process route_id filter if specified
	if routeID, ok := args["route_id"]; ok {
		routeIDStr, ok := routeID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid route_id parameter: %v", routeID)), nil
		}

		log.Printf("Filtering disruptions by route ID: %s", routeIDStr)

		// Filter disruptions for the specified route
		filteredDisruptions := make([]models.Alert, 0, len(disruptions))
		for _, disruption := range disruptions {
			affectedRoutes := disruption.GetAffectedRoutes()
			for _, route := range affectedRoutes {
				if route == routeIDStr {
					filteredDisruptions = append(filteredDisruptions, disruption)
					break
				}
			}
		}
		disruptions = filteredDisruptions
	}

	// Process severity_min filter if specified
	if severityMin, ok := args["severity_min"]; ok {
		var minSeverity int

		// Try to convert from different types
		switch v := severityMin.(type) {
		case float64:
			minSeverity = int(v)
		case int:
			minSeverity = v
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				minSeverity = parsed
			} else {
				return createErrorResponse(fmt.Sprintf("Invalid severity_min parameter: %v", severityMin)), nil
			}
		default:
			return createErrorResponse(fmt.Sprintf("Invalid severity_min parameter type: %T", severityMin)), nil
		}

		if minSeverity < 1 || minSeverity > 9 {
			return createErrorResponse("Severity level must be between 1 and 9"), nil
		}

		log.Printf("Filtering disruptions by minimum severity: %d", minSeverity)

		// Filter disruptions by minimum severity
		filteredDisruptions := make([]models.Alert, 0, len(disruptions))
		for _, disruption := range disruptions {
			if disruption.Attributes.Severity >= minSeverity {
				filteredDisruptions = append(filteredDisruptions, disruption)
			}
		}
		disruptions = filteredDisruptions
	}

	// Filter out inactive disruptions
	now := time.Now()
	activeDisruptions := make([]models.Alert, 0, len(disruptions))
	for _, disruption := range disruptions {
		if disruption.IsActive(now) {
			activeDisruptions = append(activeDisruptions, disruption)
		}
	}
	disruptions = activeDisruptions

	// If no disruptions are found, inform the user
	if len(disruptions) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No active service disruptions found matching the specified criteria.",
				},
			},
		}, nil
	}

	log.Printf("Retrieved %d service disruptions", len(disruptions))

	// Format the disruptions for response
	return formatAlertsResponse(disruptions)
}

// getAccessibilityAlertsHandler handles requests for accessibility-related alerts
func (s *Server) getAccessibilityAlertsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("Received request for accessibility alerts: %s", request.Params.Name)

	// Create MBTA client
	client := mbta.NewClient(s.config)

	// Extract parameters for filtering
	args := request.Params.Arguments

	// Get accessibility alerts
	alerts, err := client.GetAccessibilityAlerts(ctx)
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to retrieve accessibility alerts: %v", err)), nil
	}

	// Process stop_id filter if specified
	if stopID, ok := args["stop_id"]; ok {
		stopIDStr, ok := stopID.(string)
		if !ok {
			return createErrorResponse(fmt.Sprintf("Invalid stop_id parameter: %v", stopID)), nil
		}

		log.Printf("Filtering accessibility alerts by stop ID: %s", stopIDStr)

		// Filter alerts for the specified stop
		filteredAlerts := make([]models.Alert, 0, len(alerts))
		for _, alert := range alerts {
			affectedStops := alert.GetAffectedStops()
			for _, stop := range affectedStops {
				if stop == stopIDStr {
					filteredAlerts = append(filteredAlerts, alert)
					break
				}
			}
		}
		alerts = filteredAlerts
	}

	// Filter out inactive alerts
	now := time.Now()
	activeAlerts := make([]models.Alert, 0, len(alerts))
	for _, alert := range alerts {
		if alert.IsActive(now) {
			activeAlerts = append(activeAlerts, alert)
		}
	}
	alerts = activeAlerts

	// If no alerts are found, inform the user
	if len(alerts) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No active accessibility alerts found matching the specified criteria.",
				},
			},
		}, nil
	}

	log.Printf("Retrieved %d accessibility alerts", len(alerts))

	// Format the alerts for response
	return formatAlertsResponse(alerts)
}

// formatAlertsResponse converts alert data to a proper MCP response
func formatAlertsResponse(alerts []models.Alert) (*mcp.CallToolResult, error) {
	// Convert the alerts to a simplified format for the response
	now := time.Now()
	alertsData := make([]map[string]interface{}, 0, len(alerts))

	for _, alert := range alerts {
		// Format alert data
		alertMap := map[string]interface{}{
			"id":             alert.ID,
			"header":         alert.Attributes.Header,
			"description":    alert.Attributes.Description,
			"effect":         string(alert.Attributes.Effect),
			"effect_name":    models.GetAlertEffectDescription(alert.Attributes.Effect),
			"cause":          string(alert.Attributes.Cause),
			"cause_name":     models.GetAlertCauseDescription(alert.Attributes.Cause),
			"severity":       alert.Attributes.Severity,
			"severity_name":  models.GetSeverityDescription(alert.Attributes.Severity),
			"created_at":     alert.Attributes.CreatedAt,
			"updated_at":     alert.Attributes.UpdatedAt,
			"service_effect": alert.Attributes.ServiceEffect,
			"timeframe":      alert.Attributes.Timeframe,
			"lifecycle":      alert.Attributes.Lifecycle,
			"is_active":      alert.IsActive(now),
		}

		// Add URL if available
		if alert.Attributes.URL != "" {
			alertMap["url"] = alert.Attributes.URL
		}

		// Format active periods
		activePeriods := make([]map[string]interface{}, 0, len(alert.Attributes.ActivePeriod))
		for _, period := range alert.Attributes.ActivePeriod {
			periodMap := map[string]interface{}{}

			if !period.Start.IsZero() {
				periodMap["start"] = period.Start.Format(time.RFC3339)
			}

			if !period.End.IsZero() {
				periodMap["end"] = period.End.Format(time.RFC3339)
			}

			// Check if this period is currently active
			isActive := (period.Start.IsZero() || !now.Before(period.Start)) &&
				(period.End.IsZero() || !now.After(period.End))
			periodMap["is_active"] = isActive

			activePeriods = append(activePeriods, periodMap)
		}
		alertMap["active_periods"] = activePeriods

		// Add affected routes and stops
		alertMap["affected_routes"] = alert.GetAffectedRoutes()
		alertMap["affected_stops"] = alert.GetAffectedStops()

		// Create a readable summary of activities affected
		activities := make([]string, 0)
		if alert.HasActivity("BOARD") {
			activities = append(activities, "boarding")
		}
		if alert.HasActivity("EXIT") {
			activities = append(activities, "exiting")
		}
		if alert.HasActivity("RIDE") {
			activities = append(activities, "riding")
		}
		if alert.HasActivity("USING_WHEELCHAIR") {
			activities = append(activities, "wheelchair access")
		}
		if alert.HasActivity("USING_ESCALATOR") {
			activities = append(activities, "escalator use")
		}

		if len(activities) > 0 {
			alertMap["affected_activities"] = activities
		}

		alertsData = append(alertsData, alertMap)
	}

	// Create JSON string response
	jsonBytes, err := json.MarshalIndent(alertsData, "", "  ")
	if err != nil {
		return createErrorResponse(fmt.Sprintf("Failed to serialize alert data: %v", err)), nil
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

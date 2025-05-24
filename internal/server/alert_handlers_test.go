package server

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
	"github.com/mark3labs/mcp-go/mcp"
)

// These functions are left intentionally commented out for potential future use
// but are removed to avoid lint warnings about unused functions

/*
// mockAlertHandler is a helper to create a mock handler that returns predefined alerts
func mockAlertHandler(mockAlerts []models.Alert) mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := formatAlertsResponse(mockAlerts)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

// mockEmptyAlertHandler is a helper to create a mock handler that returns no alerts
func mockEmptyAlertHandler() mcpserver.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No alerts found matching the specified criteria.",
				},
			},
		}, nil
	}
}
*/

func TestFormatAlertsResponse(t *testing.T) {
	t.Run("Format general alerts", func(t *testing.T) {
		// Create test alerts
		testAlerts := []models.Alert{
			{
				ID:   "123456",
				Type: "alert",
				Attributes: models.AlertAttributes{
					Header:      "Red Line Delays",
					Description: "Red Line trains are running with delays due to construction work near Harvard Square.",
					Effect:      models.AlertEffectDelays,
					Cause:       models.AlertCauseConstruction,
					Severity:    5,
					ActivePeriod: []models.AlertPeriod{
						{
							Start: time.Now().Add(-1 * time.Hour),
							End:   time.Now().Add(2 * time.Hour),
						},
					},
					InformedEntity: []models.AlertEntity{
						{
							Activities: []string{"BOARD", "EXIT", "RIDE"},
							Route:      "Red",
							RouteType:  1,
						},
					},
				},
			},
		}

		// Call formatAlertsResponse directly to test it
		result, err := formatAlertsResponse(testAlerts)
		if err != nil {
			t.Fatalf("formatAlertsResponse returned error: %v", err)
		}

		// Check the result
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		if len(result.Content) != 1 {
			t.Fatalf("Expected 1 content item, got %d", len(result.Content))
		}

		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent type")
		}

		// Parse the JSON response to verify its structure
		var alertsData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &alertsData); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		if len(alertsData) != 1 {
			t.Fatalf("Expected 1 alert in response, got %d", len(alertsData))
		}

		alertData := alertsData[0]

		// Verify alert fields
		if alertData["id"] != "123456" {
			t.Errorf("Expected alert ID '123456', got '%v'", alertData["id"])
		}

		if alertData["header"] != "Red Line Delays" {
			t.Errorf("Expected header 'Red Line Delays', got '%v'", alertData["header"])
		}

		if alertData["effect"] != "DELAYS" {
			t.Errorf("Expected effect 'DELAYS', got '%v'", alertData["effect"])
		}

		if alertData["effect_name"] != "Delays" {
			t.Errorf("Expected effect_name 'Delays', got '%v'", alertData["effect_name"])
		}

		if alertData["cause"] != "CONSTRUCTION" {
			t.Errorf("Expected cause 'CONSTRUCTION', got '%v'", alertData["cause"])
		}

		if alertData["cause_name"] != "Construction" {
			t.Errorf("Expected cause_name 'Construction', got '%v'", alertData["cause_name"])
		}

		if alertData["severity"].(float64) != 5 {
			t.Errorf("Expected severity 5, got '%v'", alertData["severity"])
		}

		if alertData["severity_name"] != "Moderate Impact" {
			t.Errorf("Expected severity_name 'Moderate Impact', got '%v'", alertData["severity_name"])
		}

		// Check affected routes
		affectedRoutes, ok := alertData["affected_routes"].([]interface{})
		if !ok {
			t.Fatal("Expected affected_routes to be an array")
		}

		if len(affectedRoutes) != 1 || affectedRoutes[0] != "Red" {
			t.Errorf("Expected affected_routes to contain 'Red', got %v", affectedRoutes)
		}

		// Check is_active flag
		isActive, ok := alertData["is_active"].(bool)
		if !ok {
			t.Fatal("Expected is_active to be a boolean")
		}

		if !isActive {
			t.Error("Expected alert to be active")
		}
	})

	t.Run("Format accessibility alerts", func(t *testing.T) {
		// Create test accessibility alerts
		testAlerts := []models.Alert{
			{
				ID:   "567890",
				Type: "alert",
				Attributes: models.AlertAttributes{
					Header:      "Downtown Crossing: Elevator Out of Service",
					Description: "The elevator connecting the Orange Line platform to the street level is out of service.",
					Effect:      models.AlertEffectElevatorOutage,
					Cause:       models.AlertCauseEquipmentFailure,
					Severity:    5,
					ActivePeriod: []models.AlertPeriod{
						{
							Start: time.Now().Add(-12 * time.Hour),
							End:   time.Now().Add(48 * time.Hour),
						},
					},
					InformedEntity: []models.AlertEntity{
						{
							Activities: []string{"USING_WHEELCHAIR", "USING_ESCALATOR"},
							Stop:       "place-dwnxg",
						},
					},
				},
			},
		}

		// Call formatAlertsResponse
		result, err := formatAlertsResponse(testAlerts)
		if err != nil {
			t.Fatalf("formatAlertsResponse returned error: %v", err)
		}

		// Parse the JSON response
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent type")
		}

		var alertsData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &alertsData); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		alertData := alertsData[0]

		// Verify alert fields
		if alertData["id"] != "567890" {
			t.Errorf("Expected alert ID '567890', got '%v'", alertData["id"])
		}

		if alertData["effect"] != "ELEVATOR_OUTAGE" {
			t.Errorf("Expected effect 'ELEVATOR_OUTAGE', got '%v'", alertData["effect"])
		}

		if alertData["effect_name"] != "Elevator Outage" {
			t.Errorf("Expected effect_name 'Elevator Outage', got '%v'", alertData["effect_name"])
		}

		// Check affected stops
		affectedStops, ok := alertData["affected_stops"].([]interface{})
		if !ok {
			t.Fatal("Expected affected_stops to be an array")
		}

		if len(affectedStops) != 1 || affectedStops[0] != "place-dwnxg" {
			t.Errorf("Expected affected_stops to contain 'place-dwnxg', got %v", affectedStops)
		}

		// Check affected activities
		affectedActivities, ok := alertData["affected_activities"].([]interface{})
		if !ok {
			t.Fatal("Expected affected_activities to be an array")
		}

		// Check that wheelchair access is included in affected activities
		hasWheelchair := false
		for _, activity := range affectedActivities {
			if activity == "wheelchair access" {
				hasWheelchair = true
				break
			}
		}

		if !hasWheelchair {
			t.Errorf("Expected affected_activities to include 'wheelchair access', got %v", affectedActivities)
		}
	})

	t.Run("Format service disruptions", func(t *testing.T) {
		// Create test service disruptions
		testDisruptions := []models.Alert{
			{
				ID:   "987654",
				Type: "alert",
				Attributes: models.AlertAttributes{
					Header:      "Red Line: No Service Between Alewife and Harvard",
					Description: "Red Line service is suspended between Alewife and Harvard. Shuttle buses are being provided.",
					Effect:      models.AlertEffectNoService,
					Cause:       models.AlertCauseConstruction,
					Severity:    7,
					ActivePeriod: []models.AlertPeriod{
						{
							Start: time.Now().Add(-3 * time.Hour),
							End:   time.Now().Add(24 * time.Hour),
						},
					},
					InformedEntity: []models.AlertEntity{
						{
							Activities:  []string{"BOARD", "EXIT", "RIDE"},
							Route:       "Red",
							RouteType:   1,
							DirectionID: new(int),
						},
					},
				},
			},
		}

		// Call formatAlertsResponse
		result, err := formatAlertsResponse(testDisruptions)
		if err != nil {
			t.Fatalf("formatAlertsResponse returned error: %v", err)
		}

		// Parse the JSON response
		textContent, ok := result.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatal("Expected TextContent type")
		}

		var disruptionsData []map[string]interface{}
		if err := json.Unmarshal([]byte(textContent.Text), &disruptionsData); err != nil {
			t.Fatalf("Failed to parse response JSON: %v", err)
		}

		disruptionData := disruptionsData[0]

		// Verify disruption fields
		if disruptionData["id"] != "987654" {
			t.Errorf("Expected disruption ID '987654', got '%v'", disruptionData["id"])
		}

		if disruptionData["header"] != "Red Line: No Service Between Alewife and Harvard" {
			t.Errorf("Expected header 'Red Line: No Service Between Alewife and Harvard', got '%v'", disruptionData["header"])
		}

		if disruptionData["effect"] != "NO_SERVICE" {
			t.Errorf("Expected effect 'NO_SERVICE', got '%v'", disruptionData["effect"])
		}

		if disruptionData["effect_name"] != "No Service" {
			t.Errorf("Expected effect_name 'No Service', got '%v'", disruptionData["effect_name"])
		}

		if disruptionData["severity"].(float64) != 7 {
			t.Errorf("Expected severity 7, got '%v'", disruptionData["severity"])
		}
	})
}

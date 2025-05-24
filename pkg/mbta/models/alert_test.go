package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAlertUnmarshal(t *testing.T) {
	// Sample MBTA alert JSON response
	alertJSON := `{
		"data": [
			{
				"id": "12345",
				"type": "alert",
				"attributes": {
					"active_period": [
						{
							"start": "2025-05-23T10:00:00-04:00",
							"end": "2025-05-23T14:00:00-04:00"
						}
					],
					"banner": false,
					"cause": "MAINTENANCE",
					"created_at": "2025-05-22T15:30:00-04:00",
					"description": "Shuttle buses replacing Red Line service between Harvard and Alewife",
					"effect": "SHUTTLE",
					"header": "Red Line Shuttle Buses",
					"informed_entity": [
						{
							"route": "Red",
							"route_type": 1,
							"stop": "place-harsq",
							"direction_id": 0,
							"activities": ["BOARD", "EXIT", "RIDE"]
						},
						{
							"route": "Red",
							"route_type": 1,
							"stop": "place-cntsq",
							"direction_id": 0,
							"activities": ["BOARD", "EXIT", "RIDE"]
						},
						{
							"route": "Red",
							"route_type": 1,
							"stop": "place-alewf",
							"direction_id": 0,
							"activities": ["BOARD", "EXIT", "RIDE"]
						}
					],
					"lifecycle": "ONGOING",
					"severity": 7,
					"timeframe": "This weekend",
					"updated_at": "2025-05-23T08:30:00-04:00",
					"url": "https://www.mbta.com/alerts/12345",
					"service_effect": "Shuttle buses replacing Red Line service"
				}
			}
		]
	}`

	// Test unmarshaling
	var alertResponse AlertResponse
	err := json.Unmarshal([]byte(alertJSON), &alertResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal alert response: %v", err)
	}

	// Verify values
	if len(alertResponse.Data) != 1 {
		t.Fatalf("Expected 1 alert, got %d", len(alertResponse.Data))
	}

	alert := alertResponse.Data[0]
	if alert.ID != "12345" {
		t.Errorf("Expected ID 12345, got %s", alert.ID)
	}

	if alert.Type != "alert" {
		t.Errorf("Expected type alert, got %s", alert.Type)
	}

	// Check attributes
	attrs := alert.Attributes
	if attrs.Header != "Red Line Shuttle Buses" {
		t.Errorf("Expected header 'Red Line Shuttle Buses', got '%s'", attrs.Header)
	}

	if attrs.Description != "Shuttle buses replacing Red Line service between Harvard and Alewife" {
		t.Errorf("Expected description to match, got '%s'", attrs.Description)
	}

	if attrs.Effect != AlertEffectShuttle {
		t.Errorf("Expected effect SHUTTLE, got %s", attrs.Effect)
	}

	if attrs.Cause != AlertCauseMaintenance {
		t.Errorf("Expected cause MAINTENANCE, got %s", attrs.Cause)
	}

	if attrs.Severity != 7 {
		t.Errorf("Expected severity 7, got %d", attrs.Severity)
	}

	// Check active period
	if len(attrs.ActivePeriod) != 1 {
		t.Fatalf("Expected 1 active period, got %d", len(attrs.ActivePeriod))
	}

	period := attrs.ActivePeriod[0]
	expectedStart, _ := time.Parse(time.RFC3339, "2025-05-23T10:00:00-04:00")
	if !period.Start.Equal(expectedStart) {
		t.Errorf("Expected start time %v, got %v", expectedStart, period.Start)
	}

	expectedEnd, _ := time.Parse(time.RFC3339, "2025-05-23T14:00:00-04:00")
	if !period.End.Equal(expectedEnd) {
		t.Errorf("Expected end time %v, got %v", expectedEnd, period.End)
	}

	// Check informed entities
	if len(attrs.InformedEntity) != 3 {
		t.Fatalf("Expected 3 informed entities, got %d", len(attrs.InformedEntity))
	}

	entity := attrs.InformedEntity[0]
	if entity.Route != "Red" {
		t.Errorf("Expected route Red, got %s", entity.Route)
	}

	if entity.Stop != "place-harsq" {
		t.Errorf("Expected stop place-harsq, got %s", entity.Stop)
	}

	if entity.RouteType != 1 {
		t.Errorf("Expected route type 1, got %d", entity.RouteType)
	}

	// Check activities
	if len(entity.Activities) != 3 {
		t.Fatalf("Expected 3 activities, got %d", len(entity.Activities))
	}

	activities := entity.Activities
	if activities[0] != "BOARD" || activities[1] != "EXIT" || activities[2] != "RIDE" {
		t.Errorf("Expected activities [BOARD, EXIT, RIDE], got %v", activities)
	}
}

func TestAlertIsActive(t *testing.T) {
	// Create test alert with active period
	now := time.Now()
	pastStart := now.Add(-1 * time.Hour)
	futureEnd := now.Add(1 * time.Hour)
	pastEnd := now.Add(-30 * time.Minute)
	futureStart := now.Add(30 * time.Minute)

	tests := []struct {
		name         string
		periods      []AlertPeriod
		expectActive bool
	}{
		{
			name: "Currently active (now between start and end)",
			periods: []AlertPeriod{
				{
					Start: pastStart,
					End:   futureEnd,
				},
			},
			expectActive: true,
		},
		{
			name: "Not active yet (start in future)",
			periods: []AlertPeriod{
				{
					Start: futureStart,
					End:   futureEnd,
				},
			},
			expectActive: false,
		},
		{
			name: "No longer active (end in past)",
			periods: []AlertPeriod{
				{
					Start: pastStart,
					End:   pastEnd,
				},
			},
			expectActive: false,
		},
		{
			name: "Active in one of multiple periods",
			periods: []AlertPeriod{
				{
					Start: pastStart,
					End:   pastEnd,
				},
				{
					Start: pastStart,
					End:   futureEnd,
				},
			},
			expectActive: true,
		},
		{
			name: "Indefinite end time (null end)",
			periods: []AlertPeriod{
				{
					Start: pastStart,
					End:   time.Time{}, // Zero value represents null/indefinite
				},
			},
			expectActive: true,
		},
		{
			name: "Immediate start (null start)",
			periods: []AlertPeriod{
				{
					Start: time.Time{}, // Zero value represents null/immediate
					End:   futureEnd,
				},
			},
			expectActive: true,
		},
		{
			name:         "No active periods",
			periods:      []AlertPeriod{},
			expectActive: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alert := Alert{
				Attributes: AlertAttributes{
					ActivePeriod: test.periods,
				},
			}

			isActive := alert.IsActive(now)
			if isActive != test.expectActive {
				t.Errorf("Expected IsActive to be %v, got %v", test.expectActive, isActive)
			}
		})
	}
}

func TestAlertGetAffectedRoutes(t *testing.T) {
	alert := Alert{
		Attributes: AlertAttributes{
			InformedEntity: []AlertEntity{
				{Route: "Red", RouteType: 1},
				{Route: "Green-B", RouteType: 0},
				{Stop: "place-pktrm", RouteType: 1}, // No route specified
				{Route: "Red", RouteType: 1},        // Duplicate
			},
		},
	}

	routes := alert.GetAffectedRoutes()

	// We should have unique route IDs
	if len(routes) != 2 {
		t.Errorf("Expected 2 unique routes, got %d", len(routes))
	}

	// Check specific routes
	expected := map[string]bool{
		"Red":     true,
		"Green-B": true,
	}

	for _, route := range routes {
		if !expected[route] {
			t.Errorf("Unexpected route: %s", route)
		}
		// Remove from expected to detect missing routes
		delete(expected, route)
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected routes: %v", expected)
	}
}

func TestAlertGetAffectedStops(t *testing.T) {
	alert := Alert{
		Attributes: AlertAttributes{
			InformedEntity: []AlertEntity{
				{Stop: "place-harsq", RouteType: 1},
				{Stop: "place-cntsq", RouteType: 1},
				{Route: "Red", RouteType: 1},        // No stop specified
				{Stop: "place-harsq", RouteType: 1}, // Duplicate
			},
		},
	}

	stops := alert.GetAffectedStops()

	// We should have unique stop IDs
	if len(stops) != 2 {
		t.Errorf("Expected 2 unique stops, got %d", len(stops))
	}

	// Check specific stops
	expected := map[string]bool{
		"place-harsq": true,
		"place-cntsq": true,
	}

	for _, stop := range stops {
		if !expected[stop] {
			t.Errorf("Unexpected stop: %s", stop)
		}
		// Remove from expected to detect missing stops
		delete(expected, stop)
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected stops: %v", expected)
	}
}

func TestAlertHasActivity(t *testing.T) {
	alert := Alert{
		Attributes: AlertAttributes{
			InformedEntity: []AlertEntity{
				{
					Stop:       "place-harsq",
					Activities: []string{"BOARD", "EXIT"},
				},
				{
					Stop:       "place-cntsq",
					Activities: []string{"RIDE", "USING_WHEELCHAIR"},
				},
			},
		},
	}

	tests := []struct {
		activity string
		expected bool
	}{
		{"BOARD", true},
		{"EXIT", true},
		{"RIDE", true},
		{"USING_WHEELCHAIR", true},
		{"PARK_CAR", false},
		{"UNKNOWN", false},
	}

	for _, test := range tests {
		t.Run(test.activity, func(t *testing.T) {
			result := alert.HasActivity(test.activity)
			if result != test.expected {
				t.Errorf("Expected HasActivity(%s) to be %v, got %v", test.activity, test.expected, result)
			}
		})
	}
}

func TestGetAlertEffectDescription(t *testing.T) {
	tests := []struct {
		effect   AlertEffect
		expected string
	}{
		{AlertEffectNoService, "No Service"},
		{AlertEffectReducedService, "Reduced Service"},
		{AlertEffectSignificantDelays, "Significant Delays"},
		{AlertEffectDelays, "Delays"},
		{AlertEffectDetour, "Detour"},
		{AlertEffectStopMoved, "Stop Moved"},
		{AlertEffectStopClosed, "Stop Closed"},
		{AlertEffectShuttle, "Shuttle Bus Service"},
		{AlertEffectElevatorOutage, "Elevator Outage"},
		{AlertEffectAccessibilityIssue, "Accessibility Issue"},
		{AlertEffectScheduleChange, "Schedule Change"},
		{AlertEffectServiceChange, "Service Change"},
		{AlertEffectSnowRoute, "Snow Route in Effect"},
		{AlertEffectStationClosure, "Station Closure"},
		{AlertEffectTrackChange, "Track Change"},
		{AlertEffectAdditionalService, "Additional Service"},
		{AlertEffectModifiedService, "Modified Service"},
		{AlertEffectOther, "Other Effect"},
		{AlertEffect("UNKNOWN"), "Unknown Effect"},
	}

	for _, test := range tests {
		t.Run(string(test.effect), func(t *testing.T) {
			description := GetAlertEffectDescription(test.effect)
			if description != test.expected {
				t.Errorf("Expected effect description '%s', got '%s'", test.expected, description)
			}
		})
	}
}

func TestGetAlertCauseDescription(t *testing.T) {
	tests := []struct {
		cause    AlertCause
		expected string
	}{
		{AlertCauseUnknownCause, "Unknown Cause"},
		{AlertCauseUnspecifiedCause, "Unspecified Cause"},
		{AlertCauseAccident, "Accident"},
		{AlertCauseConstruction, "Construction"},
		{AlertCauseDemonstation, "Demonstration"},
		{AlertCauseEquipmentFailure, "Equipment Failure"},
		{AlertCauseMedicalEmergency, "Medical Emergency"},
		{AlertCausePoliceActivity, "Police Activity"},
		{AlertCauseMaintenance, "Maintenance"},
		{AlertCauseWeather, "Weather"},
		{AlertCauseTrafficCongestion, "Traffic Congestion"},
		{AlertCauseFireActivity, "Fire Activity"},
		{AlertCauseHoliday, "Holiday Schedule"},
		{AlertCauseStrike, "Strike"},
		{AlertCauseSuspiciousActivity, "Suspicious Activity"},
		{AlertCauseSwitchFailure, "Switch Failure"},
		{AlertCauseOther, "Other Cause"},
		{AlertCause("UNKNOWN"), "Unknown Cause"},
	}

	for _, test := range tests {
		t.Run(string(test.cause), func(t *testing.T) {
			description := GetAlertCauseDescription(test.cause)
			if description != test.expected {
				t.Errorf("Expected cause description '%s', got '%s'", test.expected, description)
			}
		})
	}
}

func TestGetSeverityDescription(t *testing.T) {
	tests := []struct {
		severity int
		expected string
	}{
		{0, "Unknown Severity"},
		{1, "Information"},
		{3, "Minor Impact"},
		{5, "Moderate Impact"},
		{7, "Severe Impact"},
		{9, "Critical Impact"},
		{10, "Unknown Severity"},
		{-1, "Unknown Severity"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Severity %d", test.severity), func(t *testing.T) {
			description := GetSeverityDescription(test.severity)
			if description != test.expected {
				t.Errorf("Expected severity description '%s', got '%s'", test.expected, description)
			}
		})
	}
}

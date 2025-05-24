package mbta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

func TestGetAlerts(t *testing.T) {
	t.Run("Get all alerts", func(t *testing.T) {
		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request
			if r.URL.Path != "/alerts" {
				t.Errorf("Expected URL path '/alerts', got '%s'", r.URL.Path)
			}

			// Check that the accept header is set correctly
			if r.Header.Get("Accept") != "application/vnd.api+json" {
				t.Errorf("Expected Accept header 'application/vnd.api+json', got '%s'", r.Header.Get("Accept"))
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": [
					{
						"attributes": {
							"active_period": [
								{
									"start": "2023-05-20T10:00:00-04:00",
									"end": "2023-05-22T16:00:00-04:00"
								}
							],
							"banner": false,
							"cause": "CONSTRUCTION",
							"created_at": "2023-05-19T14:30:00-04:00",
							"description": "Red Line trains are running with delays due to construction work near Harvard Square.",
							"effect": "DELAYS",
							"header": "Red Line Delays",
							"informed_entity": [
								{
									"activities": ["BOARD", "EXIT", "RIDE"],
									"route": "Red",
									"route_type": 1
								}
							],
							"lifecycle": "ONGOING",
							"severity": 5,
							"service_effect": "Red Line delays",
							"updated_at": "2023-05-20T09:15:00-04:00"
						},
						"id": "123456",
						"links": {
							"self": "/alerts/123456"
						},
						"type": "alert"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Get alerts
		alerts, err := client.GetAlerts(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetAlerts returned error: %v", err)
		}

		// Check alerts
		if len(alerts) != 1 {
			t.Fatalf("Expected 1 alert, got %d", len(alerts))
		}

		alert := alerts[0]
		if alert.ID != "123456" {
			t.Errorf("Expected alert ID '123456', got '%s'", alert.ID)
		}

		if alert.Attributes.Header != "Red Line Delays" {
			t.Errorf("Expected alert header 'Red Line Delays', got '%s'", alert.Attributes.Header)
		}

		if alert.Attributes.Effect != models.AlertEffectDelays {
			t.Errorf("Expected alert effect 'DELAYS', got '%s'", alert.Attributes.Effect)
		}

		if alert.Attributes.Cause != models.AlertCauseConstruction {
			t.Errorf("Expected alert cause 'CONSTRUCTION', got '%s'", alert.Attributes.Cause)
		}

		if len(alert.GetAffectedRoutes()) != 1 || alert.GetAffectedRoutes()[0] != "Red" {
			t.Errorf("Expected affected route 'Red', got %v", alert.GetAffectedRoutes())
		}
	})

	t.Run("Get alerts with filter", func(t *testing.T) {
		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request path
			if r.URL.Path != "/alerts" {
				t.Errorf("Expected URL path '/alerts', got '%s'", r.URL.Path)
			}

			// Check query parameters
			query := r.URL.Query()
			routeFilter := query.Get("filter[route]")
			if routeFilter != "Green-B" {
				t.Errorf("Expected filter[route]=Green-B, got '%s'", routeFilter)
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": [
					{
						"attributes": {
							"active_period": [
								{
									"start": "2023-05-20T08:00:00-04:00",
									"end": "2023-05-20T18:00:00-04:00"
								}
							],
							"banner": true,
							"cause": "MAINTENANCE",
							"created_at": "2023-05-19T16:30:00-04:00",
							"description": "Green Line B Branch service is suspended between Boston College and Kenmore. Use shuttle buses instead.",
							"effect": "SHUTTLE",
							"header": "Green Line B Shuttle Buses",
							"informed_entity": [
								{
									"activities": ["BOARD", "EXIT", "RIDE"],
									"route": "Green-B",
									"route_type": 0
								}
							],
							"lifecycle": "ONGOING",
							"severity": 7,
							"service_effect": "Shuttle bus service in effect",
							"updated_at": "2023-05-20T07:45:00-04:00"
						},
						"id": "789012",
						"links": {
							"self": "/alerts/789012"
						},
						"type": "alert"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Create filter parameters
		params := map[string]string{
			"filter[route]": "Green-B",
		}

		// Get alerts with filter
		alerts, err := client.GetAlerts(context.Background(), params)
		if err != nil {
			t.Fatalf("GetAlerts returned error: %v", err)
		}

		// Check alerts
		if len(alerts) != 1 {
			t.Fatalf("Expected 1 alert, got %d", len(alerts))
		}

		alert := alerts[0]
		if alert.ID != "789012" {
			t.Errorf("Expected alert ID '789012', got '%s'", alert.ID)
		}

		if alert.Attributes.Header != "Green Line B Shuttle Buses" {
			t.Errorf("Expected alert header 'Green Line B Shuttle Buses', got '%s'", alert.Attributes.Header)
		}

		if alert.Attributes.Effect != models.AlertEffectShuttle {
			t.Errorf("Expected alert effect 'SHUTTLE', got '%s'", alert.Attributes.Effect)
		}

		if alert.Attributes.Severity != 7 {
			t.Errorf("Expected severity 7, got %d", alert.Attributes.Severity)
		}

		severityDesc := models.GetSeverityDescription(alert.Attributes.Severity)
		if severityDesc != "Severe Impact" {
			t.Errorf("Expected severity description 'Severe Impact', got '%s'", severityDesc)
		}
	})
}

func TestGetAlert(t *testing.T) {
	t.Run("Get alert by ID", func(t *testing.T) {
		alertID := "123456"

		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request path
			expectedPath := "/alerts/" + alertID
			if r.URL.Path != expectedPath {
				t.Errorf("Expected URL path '%s', got '%s'", expectedPath, r.URL.Path)
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"data": {
					"attributes": {
						"active_period": [
							{
								"start": "2023-05-20T10:00:00-04:00",
								"end": "2023-05-22T16:00:00-04:00"
							}
						],
						"banner": false,
						"cause": "CONSTRUCTION",
						"created_at": "2023-05-19T14:30:00-04:00",
						"description": "Red Line trains are running with delays due to construction work near Harvard Square.",
						"effect": "DELAYS",
						"header": "Red Line Delays",
						"informed_entity": [
							{
								"activities": ["BOARD", "EXIT", "RIDE"],
								"route": "Red",
								"route_type": 1
							}
						],
						"lifecycle": "ONGOING",
						"severity": 5,
						"service_effect": "Red Line delays",
						"updated_at": "2023-05-20T09:15:00-04:00"
					},
					"id": "123456",
					"links": {
						"self": "/alerts/123456"
					},
					"type": "alert"
				}
			}`))
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Get alert by ID
		alert, err := client.GetAlert(context.Background(), alertID)
		if err != nil {
			t.Fatalf("GetAlert returned error: %v", err)
		}

		// Check alert
		if alert.ID != alertID {
			t.Errorf("Expected alert ID '%s', got '%s'", alertID, alert.ID)
		}

		if alert.Attributes.Effect != models.AlertEffectDelays {
			t.Errorf("Expected effect 'DELAYS', got '%s'", alert.Attributes.Effect)
		}

		effectDesc := models.GetAlertEffectDescription(alert.Attributes.Effect)
		if effectDesc != "Delays" {
			t.Errorf("Expected effect description 'Delays', got '%s'", effectDesc)
		}

		causeDesc := models.GetAlertCauseDescription(alert.Attributes.Cause)
		if causeDesc != "Construction" {
			t.Errorf("Expected cause description 'Construction', got '%s'", causeDesc)
		}
	})

	t.Run("Get alert with error", func(t *testing.T) {
		// Create a test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{
				"errors": [
					{
						"status": "404",
						"code": "not_found",
						"title": "Not Found",
						"detail": "The requested alert was not found"
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Get alert with a non-existent ID
		_, err := client.GetAlert(context.Background(), "non-existent")
		if err == nil {
			t.Error("Expected error, got nil")
		}

		// Check if it's an API error
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("Expected *APIError, got %T", err)
		}

		if apiErr.StatusCode != 404 {
			t.Errorf("Expected status code 404, got %d", apiErr.StatusCode)
		}

		if apiErr.Detail != "The requested alert was not found" {
			t.Errorf("Expected detail 'The requested alert was not found', got '%s'", apiErr.Detail)
		}
	})
}

func TestGetActiveAlerts(t *testing.T) {
	// Create a test server with a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/alerts" {
			t.Errorf("Expected URL path '/alerts', got '%s'", r.URL.Path)
		}

		// Check query parameters
		query := r.URL.Query()
		activityFilter := query.Get("filter[activity]")
		if !strings.Contains(activityFilter, "BOARD") || !strings.Contains(activityFilter, "RIDE") {
			t.Errorf("Expected filter[activity] to contain BOARD and RIDE, got '%s'", activityFilter)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"attributes": {
						"active_period": [
							{
								"start": "2023-05-20T00:00:00-04:00",
								"end": "2023-05-25T23:59:59-04:00"
							}
						],
						"banner": true,
						"cause": "MAINTENANCE",
						"created_at": "2023-05-19T10:00:00-04:00",
						"description": "Orange Line trains will bypass State due to station improvements.",
						"effect": "STOP_CLOSED",
						"header": "State: Closed",
						"informed_entity": [
							{
								"activities": ["BOARD", "EXIT"],
								"route": "Orange",
								"route_type": 1,
								"stop": "place-state"
							}
						],
						"lifecycle": "ONGOING",
						"severity": 5,
						"service_effect": "Trains bypass State in both directions",
						"updated_at": "2023-05-20T06:00:00-04:00"
					},
					"id": "567890",
					"type": "alert"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	cfg := &config.Config{
		APIKey:     "test-key",
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	// Get active alerts
	alerts, err := client.GetActiveAlerts(context.Background())
	if err != nil {
		t.Fatalf("GetActiveAlerts returned error: %v", err)
	}

	// Check alerts
	if len(alerts) != 1 {
		t.Fatalf("Expected 1 alert, got %d", len(alerts))
	}

	alert := alerts[0]
	if alert.ID != "567890" {
		t.Errorf("Expected alert ID '567890', got '%s'", alert.ID)
	}

	if alert.Attributes.Effect != models.AlertEffectStopClosed {
		t.Errorf("Expected effect 'STOP_CLOSED', got '%s'", alert.Attributes.Effect)
	}

	if !alert.HasActivity("BOARD") {
		t.Error("Expected alert to have BOARD activity")
	}

	stopIDs := alert.GetAffectedStops()
	if len(stopIDs) != 1 || stopIDs[0] != "place-state" {
		t.Errorf("Expected affected stop 'place-state', got %v", stopIDs)
	}
}

func TestGetServiceDisruptions(t *testing.T) {
	// Create a test server with a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request path
		if r.URL.Path != "/alerts" {
			t.Errorf("Expected URL path '/alerts', got '%s'", r.URL.Path)
		}

		// Check query parameters
		query := r.URL.Query()
		effectFilter := query.Get("filter[effect]")

		// Check that the filter contains the major disruption effects
		if !strings.Contains(effectFilter, string(models.AlertEffectNoService)) ||
			!strings.Contains(effectFilter, string(models.AlertEffectSignificantDelays)) {
			t.Errorf("Expected filter[effect] to contain major disruption effects, got '%s'", effectFilter)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"attributes": {
						"active_period": [
							{
								"start": "2023-05-19T22:00:00-04:00",
								"end": "2023-05-22T04:00:00-04:00"
							}
						],
						"banner": true,
						"cause": "CONSTRUCTION",
						"created_at": "2023-05-18T14:00:00-04:00",
						"description": "Red Line service is suspended between Alewife and Harvard. Shuttle buses are being provided.",
						"effect": "NO_SERVICE",
						"header": "Red Line: No Service Between Alewife and Harvard",
						"informed_entity": [
							{
								"activities": ["BOARD", "EXIT", "RIDE"],
								"route": "Red",
								"route_type": 1,
								"direction_id": 0
							},
							{
								"activities": ["BOARD", "EXIT", "RIDE"],
								"route": "Red",
								"route_type": 1,
								"direction_id": 1
							}
						],
						"lifecycle": "ONGOING",
						"severity": 7,
						"service_effect": "Red Line suspended between Alewife and Harvard",
						"updated_at": "2023-05-19T20:30:00-04:00"
					},
					"id": "987654",
					"type": "alert"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	cfg := &config.Config{
		APIKey:     "test-key",
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	// Get service disruptions
	alerts, err := client.GetServiceDisruptions(context.Background())
	if err != nil {
		t.Fatalf("GetServiceDisruptions returned error: %v", err)
	}

	// Check alerts
	if len(alerts) != 1 {
		t.Fatalf("Expected 1 alert, got %d", len(alerts))
	}

	alert := alerts[0]
	if alert.ID != "987654" {
		t.Errorf("Expected alert ID '987654', got '%s'", alert.ID)
	}

	if alert.Attributes.Effect != models.AlertEffectNoService {
		t.Errorf("Expected effect 'NO_SERVICE', got '%s'", alert.Attributes.Effect)
	}

	if alert.Attributes.Severity != 7 {
		t.Errorf("Expected severity 7, got %d", alert.Attributes.Severity)
	}

	effectDesc := models.GetAlertEffectDescription(alert.Attributes.Effect)
	if effectDesc != "No Service" {
		t.Errorf("Expected effect description 'No Service', got '%s'", effectDesc)
	}
}

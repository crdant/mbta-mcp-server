package mbta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
)

func TestGetVehicles(t *testing.T) {
	t.Run("Get all vehicles", func(t *testing.T) {
		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request
			if r.URL.Path != "/vehicles" {
				t.Errorf("Expected URL path '/vehicles', got '%s'", r.URL.Path)
			}

			// Check that the accept header is set correctly
			if r.Header.Get("Accept") != "application/vnd.api+json" {
				t.Errorf("Expected Accept header 'application/vnd.api+json', got '%s'", r.Header.Get("Accept"))
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"data": [
					{
						"attributes": {
							"bearing": 315.0,
							"current_status": "IN_TRANSIT_TO",
							"current_stop_sequence": 310,
							"direction_id": 0,
							"label": "3673-3838",
							"latitude": 42.33982849121094,
							"longitude": -71.15853881835938,
							"speed": null,
							"updated_at": "2023-05-20T11:22:55-04:00"
						},
						"id": "G-10067",
						"links": {
							"self": "/vehicles/G-10067"
						},
						"relationships": {
							"route": {
								"data": {
									"id": "Green-B",
									"type": "route"
								}
							},
							"stop": {
								"data": {
									"id": "70107",
									"type": "stop"
								}
							},
							"trip": {
								"data": {
									"id": "36418064",
									"type": "trip"
								}
							}
						},
						"type": "vehicle"
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

		// Get vehicles
		vehicles, err := client.GetVehicles(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetVehicles returned error: %v", err)
		}

		// Check vehicles
		if len(vehicles) != 1 {
			t.Fatalf("Expected 1 vehicle, got %d", len(vehicles))
		}

		vehicle := vehicles[0]
		if vehicle.ID != "G-10067" {
			t.Errorf("Expected vehicle ID 'G-10067', got '%s'", vehicle.ID)
		}

		if vehicle.Attributes.Label != "3673-3838" {
			t.Errorf("Expected vehicle label '3673-3838', got '%s'", vehicle.Attributes.Label)
		}

		if vehicle.GetRouteID() != "Green-B" {
			t.Errorf("Expected route ID 'Green-B', got '%s'", vehicle.GetRouteID())
		}
	})

	t.Run("Get vehicles with filter", func(t *testing.T) {
		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request path
			if r.URL.Path != "/vehicles" {
				t.Errorf("Expected URL path '/vehicles', got '%s'", r.URL.Path)
			}

			// Check query parameters
			query := r.URL.Query()
			routeFilter := query.Get("filter[route]")
			if routeFilter != "Red" {
				t.Errorf("Expected filter[route]=Red, got '%s'", routeFilter)
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"data": [
					{
						"attributes": {
							"bearing": 125.0,
							"carriages": [
								{
									"label": "1861",
									"occupancy_status": "MANY_SEATS_AVAILABLE",
									"occupancy_percentage": 0
								},
								{
									"label": "1860",
									"occupancy_status": "MANY_SEATS_AVAILABLE",
									"occupancy_percentage": 0
								}
							],
							"current_status": "STOPPED_AT",
							"current_stop_sequence": 10,
							"direction_id": 0,
							"label": "1861-1860",
							"latitude": 42.39537,
							"longitude": -71.14236,
							"speed": 0,
							"updated_at": "2023-05-20T11:22:55-04:00"
						},
						"id": "R-5463D359",
						"relationships": {
							"route": {
								"data": {
									"id": "Red",
									"type": "route"
								}
							},
							"stop": {
								"data": {
									"id": "70061",
									"type": "stop"
								}
							},
							"trip": {
								"data": {
									"id": "51070305",
									"type": "trip"
								}
							}
						},
						"type": "vehicle"
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
			"filter[route]": "Red",
		}

		// Get vehicles with filter
		vehicles, err := client.GetVehicles(context.Background(), params)
		if err != nil {
			t.Fatalf("GetVehicles returned error: %v", err)
		}

		// Check vehicles
		if len(vehicles) != 1 {
			t.Fatalf("Expected 1 vehicle, got %d", len(vehicles))
		}

		vehicle := vehicles[0]
		if vehicle.ID != "R-5463D359" {
			t.Errorf("Expected vehicle ID 'R-5463D359', got '%s'", vehicle.ID)
		}

		if vehicle.Attributes.Label != "1861-1860" {
			t.Errorf("Expected vehicle label '1861-1860', got '%s'", vehicle.Attributes.Label)
		}

		if !vehicle.HasOccupancyData() {
			t.Error("Expected vehicle to have occupancy data")
		}
	})
}

func TestGetVehicle(t *testing.T) {
	t.Run("Get vehicle by ID", func(t *testing.T) {
		vehicleID := "G-10067"

		// Create a test server with a mock response
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request path
			expectedPath := "/vehicles/" + vehicleID
			if r.URL.Path != expectedPath {
				t.Errorf("Expected URL path '%s', got '%s'", expectedPath, r.URL.Path)
			}

			// Return a mock response
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"data": {
					"attributes": {
						"bearing": 315.0,
						"current_status": "IN_TRANSIT_TO",
						"current_stop_sequence": 310,
						"direction_id": 0,
						"label": "3673-3838",
						"latitude": 42.33982849121094,
						"longitude": -71.15853881835938,
						"speed": null,
						"updated_at": "2023-05-20T11:22:55-04:00"
					},
					"id": "G-10067",
					"links": {
						"self": "/vehicles/G-10067"
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Green-B",
								"type": "route"
							}
						},
						"stop": {
							"data": {
								"id": "70107",
								"type": "stop"
							}
						},
						"trip": {
							"data": {
								"id": "36418064",
								"type": "trip"
							}
						}
					},
					"type": "vehicle"
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

		// Get vehicle by ID
		vehicle, err := client.GetVehicle(context.Background(), vehicleID)
		if err != nil {
			t.Fatalf("GetVehicle returned error: %v", err)
		}

		// Check vehicle
		if vehicle.ID != vehicleID {
			t.Errorf("Expected vehicle ID '%s', got '%s'", vehicleID, vehicle.ID)
		}

		if vehicle.Attributes.CurrentStatus != "IN_TRANSIT_TO" {
			t.Errorf("Expected current_status 'IN_TRANSIT_TO', got '%s'", vehicle.Attributes.CurrentStatus)
		}

		if status := vehicle.GetStatusDescription(); status != "In Transit" {
			t.Errorf("Expected status description 'In Transit', got '%s'", status)
		}

		if vehicle.GetStopID() != "70107" {
			t.Errorf("Expected stop ID '70107', got '%s'", vehicle.GetStopID())
		}

		if vehicle.GetTripID() != "36418064" {
			t.Errorf("Expected trip ID '36418064', got '%s'", vehicle.GetTripID())
		}
	})

	t.Run("Get vehicle with error", func(t *testing.T) {
		// Create a test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{
				"errors": [
					{
						"status": "404",
						"code": "not_found",
						"title": "Not Found",
						"detail": "The requested vehicle was not found"
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

		// Get vehicle with a non-existent ID
		_, err := client.GetVehicle(context.Background(), "non-existent")
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

		if apiErr.Detail != "The requested vehicle was not found" {
			t.Errorf("Expected detail 'The requested vehicle was not found', got '%s'", apiErr.Detail)
		}
	})
}
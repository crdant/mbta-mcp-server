package mbta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
)

func TestGetTrips(t *testing.T) {
	// Setup mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		
		if r.URL.Path != "/trips" {
			t.Errorf("Expected path /trips, got %s", r.URL.Path)
		}

		// Check if filter params were sent correctly
		if routeID := r.URL.Query().Get("filter[route]"); routeID != "Red" {
			t.Errorf("Expected filter[route]=Red, got %s", routeID)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"data": [
				{
					"id": "trip1",
					"type": "trip",
					"attributes": {
						"name": "Test Trip",
						"headsign": "Alewife",
						"direction_id": 0,
						"service_id": "service1",
						"wheelchair_accessible": true,
						"bikes_allowed": true
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Red",
								"type": "route"
							}
						}
					}
				}
			]
		}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	// Start mock server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client with mock server URL
	cfg := &config.Config{
		APIBaseURL: server.URL,
		APIKey:     "test-key",
		Timeout:    5 * time.Second,
	}
	client := NewClient(cfg)

	// Test GetTrips
	params := map[string]string{
		"filter[route]": "Red",
	}
	trips, err := client.GetTrips(context.Background(), params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(trips) != 1 {
		t.Fatalf("Expected 1 trip, got %d", len(trips))
	}

	trip := trips[0]
	if trip.ID != "trip1" {
		t.Errorf("Expected trip ID 'trip1', got '%s'", trip.ID)
	}

	if trip.Attributes.Headsign != "Alewife" {
		t.Errorf("Expected headsign 'Alewife', got '%s'", trip.Attributes.Headsign)
	}

	if routeID := trip.GetRouteID(); routeID != "Red" {
		t.Errorf("Expected route ID 'Red', got '%s'", routeID)
	}

	if !trip.IsWheelchairAccessible() {
		t.Error("Expected trip to be wheelchair accessible, but it wasn't")
	}

	if !trip.IsBikeAllowed() {
		t.Error("Expected trip to allow bikes, but it didn't")
	}
}

func TestGetTrip(t *testing.T) {
	// Setup mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		
		if r.URL.Path != "/trips/trip1" {
			t.Errorf("Expected path /trips/trip1, got %s", r.URL.Path)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"data": {
				"id": "trip1",
				"type": "trip",
				"attributes": {
					"name": "Test Trip",
					"headsign": "Alewife",
					"direction_id": 0,
					"service_id": "service1",
					"wheelchair_accessible": true,
					"bikes_allowed": true
				},
				"relationships": {
					"route": {
						"data": {
							"id": "Red",
							"type": "route"
						}
					}
				}
			}
		}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	// Start mock server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client with mock server URL
	cfg := &config.Config{
		APIBaseURL: server.URL,
		APIKey:     "test-key",
		Timeout:    5 * time.Second,
	}
	client := NewClient(cfg)

	// Test GetTrip
	trip, err := client.GetTrip(context.Background(), "trip1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if trip.ID != "trip1" {
		t.Errorf("Expected trip ID 'trip1', got '%s'", trip.ID)
	}

	if trip.Attributes.Headsign != "Alewife" {
		t.Errorf("Expected headsign 'Alewife', got '%s'", trip.Attributes.Headsign)
	}

	if routeID := trip.GetRouteID(); routeID != "Red" {
		t.Errorf("Expected route ID 'Red', got '%s'", routeID)
	}

	if !trip.IsWheelchairAccessible() {
		t.Error("Expected trip to be wheelchair accessible, but it wasn't")
	}

	if !trip.IsBikeAllowed() {
		t.Error("Expected trip to allow bikes, but it didn't")
	}
}

func TestGetTripsByRoute(t *testing.T) {
	// Setup mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		
		if r.URL.Path != "/trips" {
			t.Errorf("Expected path /trips, got %s", r.URL.Path)
		}

		// Check if filter params were sent correctly
		if routeID := r.URL.Query().Get("filter[route]"); routeID != "Red" {
			t.Errorf("Expected filter[route]=Red, got %s", routeID)
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"data": [
				{
					"id": "trip1",
					"type": "trip",
					"attributes": {
						"name": "Trip 1",
						"headsign": "Alewife",
						"direction_id": 0,
						"service_id": "service1",
						"wheelchair_accessible": true,
						"bikes_allowed": true
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Red",
								"type": "route"
							}
						}
					}
				},
				{
					"id": "trip2",
					"type": "trip",
					"attributes": {
						"name": "Trip 2",
						"headsign": "Ashmont",
						"direction_id": 1,
						"service_id": "service1",
						"wheelchair_accessible": false,
						"bikes_allowed": true
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Red",
								"type": "route"
							}
						}
					}
				}
			]
		}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	// Start mock server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client with mock server URL
	cfg := &config.Config{
		APIBaseURL: server.URL,
		APIKey:     "test-key",
		Timeout:    5 * time.Second,
	}
	client := NewClient(cfg)

	// Test GetTripsByRoute
	trips, err := client.GetTripsByRoute(context.Background(), "Red")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(trips) != 2 {
		t.Fatalf("Expected 2 trips, got %d", len(trips))
	}

	if trips[0].ID != "trip1" || trips[1].ID != "trip2" {
		t.Errorf("Expected trip IDs 'trip1' and 'trip2', got '%s' and '%s'", trips[0].ID, trips[1].ID)
	}

	if trips[0].Attributes.Direction != 0 || trips[1].Attributes.Direction != 1 {
		t.Errorf("Expected directions 0 and 1, got %d and %d", trips[0].Attributes.Direction, trips[1].Attributes.Direction)
	}

	if trips[0].Attributes.Headsign != "Alewife" || trips[1].Attributes.Headsign != "Ashmont" {
		t.Errorf("Expected headsigns 'Alewife' and 'Ashmont', got '%s' and '%s'", trips[0].Attributes.Headsign, trips[1].Attributes.Headsign)
	}

	if !trips[0].IsWheelchairAccessible() || trips[1].IsWheelchairAccessible() {
		t.Error("Unexpected wheelchair accessibility values")
	}
}

func TestFindCommonRoutes(t *testing.T) {
	tests := []struct {
		name     string
		routesA  []string
		routesB  []string
		expected []string
	}{
		{
			name:     "No common routes",
			routesA:  []string{"Red", "Blue"},
			routesB:  []string{"Green", "Orange"},
			expected: []string{},
		},
		{
			name:     "One common route",
			routesA:  []string{"Red", "Blue"},
			routesB:  []string{"Blue", "Orange"},
			expected: []string{"Blue"},
		},
		{
			name:     "Multiple common routes",
			routesA:  []string{"Red", "Blue", "Orange"},
			routesB:  []string{"Blue", "Orange", "Green"},
			expected: []string{"Blue", "Orange"},
		},
		{
			name:     "Identical route lists",
			routesA:  []string{"Red", "Blue", "Orange"},
			routesB:  []string{"Red", "Blue", "Orange"},
			expected: []string{"Red", "Blue", "Orange"},
		},
		{
			name:     "Empty route lists",
			routesA:  []string{},
			routesB:  []string{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := findCommonRoutes(test.routesA, test.routesB)
			
			// Check length
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d common routes, got %d", len(test.expected), len(result))
				return
			}
			
			// Create map for easier comparison
			expectedMap := make(map[string]bool)
			for _, route := range test.expected {
				expectedMap[route] = true
			}
			
			// Verify each result is in expected
			for _, route := range result {
				if !expectedMap[route] {
					t.Errorf("Unexpected route in result: %s", route)
				}
			}
		})
	}
}

func TestFindCommonStops(t *testing.T) {
	tests := []struct {
		name     string
		stopsA   []string
		stopsB   []string
		expected []string
	}{
		{
			name:     "No common stops",
			stopsA:   []string{"stop1", "stop2"},
			stopsB:   []string{"stop3", "stop4"},
			expected: []string{},
		},
		{
			name:     "One common stop",
			stopsA:   []string{"stop1", "stop2"},
			stopsB:   []string{"stop2", "stop3"},
			expected: []string{"stop2"},
		},
		{
			name:     "Multiple common stops",
			stopsA:   []string{"stop1", "stop2", "stop3"},
			stopsB:   []string{"stop2", "stop3", "stop4"},
			expected: []string{"stop2", "stop3"},
		},
		{
			name:     "Identical stop lists",
			stopsA:   []string{"stop1", "stop2", "stop3"},
			stopsB:   []string{"stop1", "stop2", "stop3"},
			expected: []string{"stop1", "stop2", "stop3"},
		},
		{
			name:     "Empty stop lists",
			stopsA:   []string{},
			stopsB:   []string{},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := findCommonStops(test.stopsA, test.stopsB)
			
			// Check length
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d common stops, got %d", len(test.expected), len(result))
				return
			}
			
			// Create map for easier comparison
			expectedMap := make(map[string]bool)
			for _, stop := range test.expected {
				expectedMap[stop] = true
			}
			
			// Verify each result is in expected
			for _, stop := range result {
				if !expectedMap[stop] {
					t.Errorf("Unexpected stop in result: %s", stop)
				}
			}
		})
	}
}

func TestCalculateApproximateDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64 // Allowed deviation
	}{
		{
			name:     "Same point",
			lat1:     42.3601,
			lon1:     -71.0589,
			lat2:     42.3601,
			lon2:     -71.0589,
			expected: 0.0,
			delta:    0.1,
		},
		{
			name:     "Short distance (Harvard Square to Central Square)",
			lat1:     42.3736,
			lon1:     -71.1190,
			lat2:     42.3654,
			lon2:     -71.1037,
			expected: 1.2, // ~1.2 km
			delta:    0.5,
		},
		{
			name:     "Medium distance (Downtown Boston to Alewife)",
			lat1:     42.3554,
			lon1:     -71.0603,
			lat2:     42.3954,
			lon2:     -71.1426,
			expected: 7.8, // ~7.8 km
			delta:    1.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			distance := calculateApproximateDistance(test.lat1, test.lon1, test.lat2, test.lon2)
			
			diff := distance - test.expected
			if diff < 0 {
				diff = -diff
			}
			
			if diff > test.delta {
				t.Errorf("Expected distance around %f km, got %f km (diff: %f)", test.expected, distance, diff)
			}
		})
	}
}
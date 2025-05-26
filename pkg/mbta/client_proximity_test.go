package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

func TestFindNearbyStations(t *testing.T) {
	t.Run("Successfully returns nearby stations", func(t *testing.T) {
		// Create test data for stops
		testStops := []models.Stop{
			{
				ID:   "place-pktrm",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Park Street",
					Latitude:           42.35639457,
					Longitude:          -71.0624242,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			{
				ID:   "place-dwnxg",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Downtown Crossing",
					Latitude:           42.355518,
					Longitude:          -71.060225,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			{
				ID:   "place-boyls",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Boylston",
					Latitude:           42.35302,
					Longitude:          -71.06459,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			{
				ID:   "place-armnl",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Arlington",
					Latitude:           42.351902,
					Longitude:          -71.070893,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			{
				ID:   "place-coecl",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Copley",
					Latitude:           42.349974,
					Longitude:          -71.077447,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
		}

		// Create test server that returns all stops
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request
			if r.URL.Path != "/stops" {
				t.Errorf("Expected request to '/stops', got: %s", r.URL.Path)
			}

			// Check for query parameters
			locationTypeParam := r.URL.Query().Get("filter[location_type]")
			if locationTypeParam != "1" {
				t.Errorf("Expected location_type filter to be '1', got: %s", locationTypeParam)
			}

			// Return all stops - client will filter by distance
			resp := models.StopResponse{
				Data: testStops,
			}

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				t.Fatalf("Failed to encode test response: %v", err)
			}
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-api-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Test coordinates near Downtown Crossing
		lat := 42.355
		lon := -71.060
		radius := 1.0 // 1 km radius
		maxResults := 3
		onlyStations := true

		// Call the function under test
		results, err := client.FindNearbyStations(context.Background(), lat, lon, radius, maxResults, onlyStations)

		// Verify results
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return 3 results within 1km, sorted by distance
		if len(results) != 3 {
			t.Errorf("Expected 3 results, got: %d", len(results))
		}

		// First result should be Downtown Crossing (closest to test coordinates)
		if len(results) > 0 && results[0].Stop.ID != "place-dwnxg" {
			t.Errorf("Expected first result to be Downtown Crossing (place-dwnxg), got: %s", results[0].Stop.ID)
		}

		// Verify distances are calculated and in ascending order
		if len(results) >= 2 {
			if results[0].DistanceKm > results[1].DistanceKm {
				t.Errorf("Expected results to be sorted by distance, but %f > %f",
					results[0].DistanceKm, results[1].DistanceKm)
			}
		}
	})

	t.Run("Handles error from API", func(t *testing.T) {
		// Create test server that returns an error
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintln(w, `{"errors":[{"status":"500","title":"Internal Server Error"}]}`)
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-api-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Call the function under test
		_, err := client.FindNearbyStations(context.Background(), 42.355, -71.060, 1.0, 5, true)

		// Verify error is returned
		if err == nil {
			t.Fatal("Expected error when API returns 500, got nil")
		}
	})

	t.Run("Returns empty results when no stations in range", func(t *testing.T) {
		// Create test data for stops
		testStops := []models.Stop{
			{
				ID:   "place-pktrm",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Park Street",
					Latitude:           42.35639457,
					Longitude:          -71.0624242,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
		}

		// Create test server that returns all stops
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := models.StopResponse{
				Data: testStops,
			}

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				t.Fatalf("Failed to encode test response: %v", err)
			}
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-api-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Test coordinates far from any stations (New York City)
		lat := 40.7128
		lon := -74.0060
		radius := 1.0 // 1 km radius
		maxResults := 5
		onlyStations := true

		// Call the function under test
		results, err := client.FindNearbyStations(context.Background(), lat, lon, radius, maxResults, onlyStations)

		// Verify results
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return 0 results as no stations are within 1km of the test coordinates
		if len(results) != 0 {
			t.Errorf("Expected 0 results, got: %d", len(results))
		}
	})

	t.Run("Includes platforms when onlyStations is false", func(t *testing.T) {
		// Create test data with a mix of stations and platforms
		testStops := []models.Stop{
			{
				ID:   "place-pktrm",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Park Street",
					Latitude:           42.35639457,
					Longitude:          -71.0624242,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypeStation,
				},
			},
			{
				ID:   "70075",
				Type: "stop",
				Attributes: models.StopAttributes{
					Name:               "Park Street Platform",
					Latitude:           42.35639457,
					Longitude:          -71.0624242,
					WheelchairBoarding: models.WheelchairBoardingAccessible,
					Municipality:       "Boston",
					LocationType:       models.LocationTypePlatform,
				},
			},
		}

		// Create test server that returns all stops
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify the request
			if r.URL.Path != "/stops" {
				t.Errorf("Expected request to '/stops', got: %s", r.URL.Path)
			}

			// Check if location_type filter is not present
			locationTypeParam := r.URL.Query().Get("filter[location_type]")
			if locationTypeParam != "" {
				t.Errorf("Expected no location_type filter, got: %s", locationTypeParam)
			}

			resp := models.StopResponse{
				Data: testStops,
			}

			w.Header().Set("Content-Type", "application/vnd.api+json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				t.Fatalf("Failed to encode test response: %v", err)
			}
		}))
		defer server.Close()

		// Create client with test server URL
		cfg := &config.Config{
			APIKey:     "test-api-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Test coordinates near Park Street
		lat := 42.356
		lon := -71.062
		radius := 1.0 // 1 km radius
		maxResults := 5
		onlyStations := false // Include platforms

		// Call the function under test
		results, err := client.FindNearbyStations(context.Background(), lat, lon, radius, maxResults, onlyStations)

		// Verify results
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Should return both the station and platform
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got: %d", len(results))
		}

		// Verify both station and platform types are included
		hasStation := false
		hasPlatform := false
		for _, result := range results {
			if result.Stop.Attributes.LocationType == models.LocationTypeStation {
				hasStation = true
			}
			if result.Stop.Attributes.LocationType == models.LocationTypePlatform {
				hasPlatform = true
			}
		}

		if !hasStation {
			t.Error("Expected results to include a station")
		}
		if !hasPlatform {
			t.Error("Expected results to include a platform")
		}
	})
}
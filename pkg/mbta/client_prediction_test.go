package mbta

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

func TestGetPredictions(t *testing.T) {
	t.Run("GetPredictions returns valid predictions", func(t *testing.T) {
		// Create a test server with a mock response
		mockPredictionsResponse := `{
			"data": [
				{
					"id": "prediction-123",
					"type": "prediction",
					"attributes": {
						"arrival_time": "2025-06-01T14:30:00-04:00",
						"departure_time": "2025-06-01T14:32:00-04:00",
						"direction_id": 0,
						"schedule_relationship": "SCHEDULED",
						"status": null,
						"stop_sequence": 5,
						"track": "2"
					},
					"relationships": {
						"route": {
							"data": {
								"id": "Red",
								"type": "route"
							}
						},
						"stop": {
							"data": {
								"id": "place-sstat",
								"type": "stop"
							}
						},
						"trip": {
							"data": {
								"id": "CR-Weekday-Fall-17-515",
								"type": "trip"
							}
						},
						"vehicle": {
							"data": {
								"id": "vehicle-123",
								"type": "vehicle"
							}
						}
					}
				}
			]
		}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockPredictionsResponse))
		}))
		defer server.Close()

		// Create a client pointing to the test server
		cfg := &config.Config{
			APIKey:     "test-api-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Test basic predictions retrieval
		predictions, err := client.GetPredictions(context.Background(), nil)
		if err != nil {
			t.Fatalf("GetPredictions returned error: %v", err)
		}

		if len(predictions) != 1 {
			t.Fatalf("Expected 1 prediction, got %d", len(predictions))
		}

		if predictions[0].ID != "prediction-123" {
			t.Errorf("Expected prediction ID 'prediction-123', got '%s'", predictions[0].ID)
		}

		if predictions[0].GetRouteID() != "Red" {
			t.Errorf("Expected route ID 'Red', got '%s'", predictions[0].GetRouteID())
		}

		if predictions[0].GetVehicleID() != "vehicle-123" {
			t.Errorf("Expected vehicle ID 'vehicle-123', got '%s'", predictions[0].GetVehicleID())
		}
	})

	t.Run("GetPredictionsByVehicle adds correct filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify that the request has the correct filter parameter
			if r.URL.Path != "/predictions" {
				t.Errorf("Expected path '/predictions', got '%s'", r.URL.Path)
			}

			vehicleFilter := r.URL.Query().Get("filter[vehicle]")
			if vehicleFilter != "test-vehicle-id" {
				t.Errorf("Expected filter[vehicle]='test-vehicle-id', got '%s'", vehicleFilter)
			}

			// Return a valid but empty response
			response := models.PredictionResponse{
				Data: []models.Prediction{},
			}
			jsonBytes, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(jsonBytes)
		}))
		defer server.Close()

		// Create a client pointing to the test server
		cfg := &config.Config{
			APIKey:     "test-api-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Call the method being tested
		_, err := client.GetPredictionsByVehicle(context.Background(), "test-vehicle-id")
		if err != nil {
			t.Fatalf("GetPredictionsByVehicle returned error: %v", err)
		}
	})

	t.Run("GetPredictionsByRoute adds correct filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify that the request has the correct filter parameter
			if r.URL.Path != "/predictions" {
				t.Errorf("Expected path '/predictions', got '%s'", r.URL.Path)
			}

			routeFilter := r.URL.Query().Get("filter[route]")
			if routeFilter != "Red" {
				t.Errorf("Expected filter[route]='Red', got '%s'", routeFilter)
			}

			// Return a valid but empty response
			response := models.PredictionResponse{
				Data: []models.Prediction{},
			}
			jsonBytes, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(jsonBytes)
		}))
		defer server.Close()

		// Create a client pointing to the test server
		cfg := &config.Config{
			APIKey:     "test-api-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Call the method being tested
		_, err := client.GetPredictionsByRoute(context.Background(), "Red")
		if err != nil {
			t.Fatalf("GetPredictionsByRoute returned error: %v", err)
		}
	})

	t.Run("GetPredictionsByStop adds correct filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify that the request has the correct filter parameter
			if r.URL.Path != "/predictions" {
				t.Errorf("Expected path '/predictions', got '%s'", r.URL.Path)
			}

			stopFilter := r.URL.Query().Get("filter[stop]")
			if stopFilter != "place-sstat" {
				t.Errorf("Expected filter[stop]='place-sstat', got '%s'", stopFilter)
			}

			// Return a valid but empty response
			response := models.PredictionResponse{
				Data: []models.Prediction{},
			}
			jsonBytes, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(jsonBytes)
		}))
		defer server.Close()

		// Create a client pointing to the test server
		cfg := &config.Config{
			APIKey:     "test-api-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Call the method being tested
		_, err := client.GetPredictionsByStop(context.Background(), "place-sstat")
		if err != nil {
			t.Fatalf("GetPredictionsByStop returned error: %v", err)
		}
	})

	t.Run("GetPredictionsByLocation adds correct filter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify that the request has the correct filter parameters
			if r.URL.Path != "/predictions" {
				t.Errorf("Expected path '/predictions', got '%s'", r.URL.Path)
			}

			latFilter := r.URL.Query().Get("filter[latitude]")
			lonFilter := r.URL.Query().Get("filter[longitude]")
			radiusFilter := r.URL.Query().Get("filter[radius]")

			if latFilter != "42.360100" {
				t.Errorf("Expected filter[latitude]='42.360100', got '%s'", latFilter)
			}
			if lonFilter != "-71.058900" {
				t.Errorf("Expected filter[longitude]='-71.058900', got '%s'", lonFilter)
			}
			if radiusFilter != "0.050000" {
				t.Errorf("Expected filter[radius]='0.050000', got '%s'", radiusFilter)
			}

			// Return a valid but empty response
			response := models.PredictionResponse{
				Data: []models.Prediction{},
			}
			jsonBytes, _ := json.Marshal(response)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(jsonBytes)
		}))
		defer server.Close()

		// Create a client pointing to the test server
		cfg := &config.Config{
			APIKey:     "test-api-key",
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Call the method being tested
		_, err := client.GetPredictionsByLocation(context.Background(), 42.3601, -71.0589, 0.05)
		if err != nil {
			t.Fatalf("GetPredictionsByLocation returned error: %v", err)
		}
	})
}

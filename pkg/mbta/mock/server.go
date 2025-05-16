// Package mock provides mock HTTP servers and responses for testing the MBTA client
package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

// ResponseDefinition defines how to respond to a specific API request
type ResponseDefinition struct {
	Path       string
	Method     string
	StatusCode int
	Response   string
	Headers    map[string]string
}

// NewMockServer creates a new test server that returns predefined responses
func NewMockServer(definitions []ResponseDefinition) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find a matching response definition
		for _, def := range definitions {
			if def.Method == r.Method && def.Path == r.URL.Path {
				// Set headers
				for key, value := range def.Headers {
					w.Header().Set(key, value)
				}
				// Set status code
				w.WriteHeader(def.StatusCode)
				// Write response
				_, _ = w.Write([]byte(def.Response))
				return
			}
		}

		// If no definition matches, return 404
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"status":"404","code":"not_found","title":"Not Found","detail":"The requested resource was not found"}]}`))
	}))
}

// LoadFixture loads a fixture file from the testdata directory
func LoadFixture(filename string) (string, error) {
	// Use testdata directory which is more reliable than relative paths
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MockValidAPIKeyMiddleware mocks the API key validation middleware
func MockValidAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || apiKey == "invalid-key" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"errors":[{"status":"401","code":"unauthorized","title":"Unauthorized request","detail":"API key missing or invalid"}]}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MockRateLimitMiddleware mocks the rate limit middleware
func MockRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If request path contains "rate-limited", simulate a rate limit response
		if strings.Contains(r.URL.Path, "rate-limited") {
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"errors":[{"status":"429","code":"rate_limited","title":"Rate Limit Exceeded","detail":"You have exceeded your rate limit"}]}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MockTimeoutHandler simulates a timeout response
func MockTimeoutHandler(w http.ResponseWriter, r *http.Request) {
	// Close the connection without writing a response
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = conn.Close()
}

// StandardMockServer creates a standard mock server with routes, stops, and schedules
func StandardMockServer() (*httptest.Server, error) {
	// Create embedded test data
	routesResponse := `{
		"data": [
			{
				"id": "Red",
				"type": "route",
				"attributes": {
					"color": "DA291C",
					"description": "Rapid Transit",
					"direction_destinations": ["Alewife", "Ashmont/Braintree"],
					"direction_names": ["Outbound", "Inbound"],
					"fare_class": "Rapid Transit",
					"long_name": "Red Line",
					"short_name": "",
					"sort_order": 10010,
					"text_color": "FFFFFF",
					"type": 1
				},
				"links": {"self": "/routes/Red"},
				"relationships": {
					"line": {"data": {"id": "line-Red", "type": "line"}}
				}
			},
			{
				"id": "Orange",
				"type": "route",
				"attributes": {
					"color": "ED8B00",
					"description": "Rapid Transit",
					"direction_destinations": ["Oak Grove", "Forest Hills"],
					"direction_names": ["Outbound", "Inbound"],
					"fare_class": "Rapid Transit",
					"long_name": "Orange Line",
					"short_name": "",
					"sort_order": 10020,
					"text_color": "FFFFFF",
					"type": 1
				},
				"links": {"self": "/routes/Orange"},
				"relationships": {
					"line": {"data": {"id": "line-Orange", "type": "line"}}
				}
			}
		]
	}`

	stopsResponse := `{
		"data": [
			{
				"id": "place-north",
				"type": "stop",
				"attributes": {
					"address": "North Station, Boston, MA 02114",
					"description": "North Station - Commuter Rail, Orange Line, and Green Line",
					"latitude": 42.365577,
					"location_type": 1,
					"longitude": -71.06129,
					"municipality": "Boston",
					"name": "North Station",
					"wheelchair_boarding": 1
				},
				"relationships": {
					"parent_station": {"data": null},
					"zone": {"data": {"id": "CR-zone-1A", "type": "zone"}}
				},
				"links": {"self": "/stops/place-north"}
			},
			{
				"id": "70061",
				"type": "stop",
				"attributes": {
					"description": "Orange Line platform for Forest Hills-bound trains",
					"latitude": 42.365486,
					"location_type": 0,
					"longitude": -71.06129,
					"municipality": "Boston",
					"name": "North Station",
					"platform_code": "1",
					"platform_name": "Orange Line - Forest Hills",
					"wheelchair_boarding": 1
				},
				"relationships": {
					"parent_station": {"data": {"id": "place-north", "type": "stop"}}
				},
				"links": {"self": "/stops/70061"}
			}
		]
	}`

	schedulesResponse := `{
		"data": [
			{
				"id": "schedule-1",
				"type": "schedule",
				"attributes": {
					"arrival_time": "2023-05-20T12:00:00-04:00",
					"departure_time": "2023-05-20T12:02:00-04:00",
					"drop_off_type": 0,
					"pickup_type": 0,
					"stop_headsign": "Alewife",
					"stop_sequence": 1,
					"timepoint": true
				},
				"relationships": {
					"route": {"data": {"id": "Red", "type": "route"}},
					"stop": {"data": {"id": "place-sstat", "type": "stop"}},
					"trip": {"data": {"id": "Red-123456-20230520", "type": "trip"}}
				}
			},
			{
				"id": "schedule-2",
				"type": "schedule",
				"attributes": {
					"arrival_time": "2023-05-20T12:10:00-04:00",
					"departure_time": "2023-05-20T12:11:00-04:00",
					"drop_off_type": 0,
					"pickup_type": 0,
					"stop_headsign": "Alewife",
					"stop_sequence": 2,
					"timepoint": true
				},
				"relationships": {
					"route": {"data": {"id": "Red", "type": "route"}},
					"stop": {"data": {"id": "place-dwnxg", "type": "stop"}},
					"trip": {"data": {"id": "Red-123456-20230520", "type": "trip"}}
				}
			}
		],
		"included": [
			{
				"id": "Red-123456-20230520",
				"type": "trip",
				"attributes": {
					"block_id": "R-123456-2023",
					"direction_id": 0,
					"headsign": "Alewife",
					"name": "",
					"wheelchair_accessible": 1
				},
				"relationships": {
					"route": {"data": {"id": "Red", "type": "route"}},
					"service": {"data": {"id": "service-weekday", "type": "service"}}
				}
			}
		]
	}`

	// Define response handler
	handler := http.NewServeMux()

	// Routes endpoint
	handler.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, routesResponse)
	})

	// Single route endpoint
	handler.HandleFunc("/routes/Red", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		// Extract just the Red line from the fixture
		_, _ = io.WriteString(w, `{"data": {"id":"Red","type":"route","attributes":{"color":"DA291C","description":"Rapid Transit","direction_destinations":["Alewife","Ashmont/Braintree"],"direction_names":["Outbound","Inbound"],"fare_class":"Rapid Transit","long_name":"Red Line","short_name":"","sort_order":10010,"text_color":"FFFFFF","type":1},"links":{"self":"/routes/Red"},"relationships":{"line":{"data":{"id":"line-Red","type":"line"}}}}}`)
	})

	// Stops endpoint
	handler.HandleFunc("/stops", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, stopsResponse)
	})

	// Schedules endpoint
	handler.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, schedulesResponse)
	})

	// Not found handler
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"errors":[{"status":"404","code":"not_found","title":"Not Found","detail":"The requested resource was not found"}]}`)
	})

	// Wrap with middleware
	var wrappedHandler http.Handler = handler
	wrappedHandler = MockValidAPIKeyMiddleware(wrappedHandler)
	wrappedHandler = MockRateLimitMiddleware(wrappedHandler)

	// Create and return test server
	return httptest.NewServer(wrappedHandler), nil
}

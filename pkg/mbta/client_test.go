package mbta

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		APIKey:     "test-api-key",
		Timeout:    5 * time.Second,
		APIBaseURL: "https://api-test.mbta.com",
	}

	client := NewClient(cfg)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.apiKey != cfg.APIKey {
		t.Errorf("Expected apiKey to be %q, got %q", cfg.APIKey, client.apiKey)
	}

	if client.baseURL != cfg.APIBaseURL {
		t.Errorf("Expected baseURL to be %q, got %q", cfg.APIBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Fatal("Expected httpClient to be created, got nil")
	}

	transport, ok := client.httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatal("Expected httpClient.Transport to be *http.Transport")
	}

	// Check that the transport has reasonable defaults
	if transport.MaxIdleConns <= 0 {
		t.Error("Expected MaxIdleConns to be set")
	}

	// Check timeout
	if client.httpClient.Timeout != cfg.Timeout {
		t.Errorf("Expected timeout to be %v, got %v", cfg.Timeout, client.httpClient.Timeout)
	}
}

func TestClientTimeout(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // Delay response
		_, _ = fmt.Fprintln(w, `{"data": []}`)
	}))
	defer server.Close()

	// Create client with a very short timeout
	cfg := &config.Config{
		APIKey:     "test-api-key",
		Timeout:    100 * time.Millisecond, // Short timeout to force error
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", server.URL+"/routes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Attempt the request - should timeout
	_, err = client.httpClient.Do(req)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Check that we got a timeout error
	if !isTimeoutError(err) {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

// Note: isTimeoutError has been moved to client.go

func TestClientHeaders(t *testing.T) {
	apiKey := "test-api-key"

	// Create a test server that verifies headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		authHeader := r.Header.Get("X-API-Key")
		if authHeader != apiKey {
			t.Errorf("Expected X-API-Key header to be %q, got %q", apiKey, authHeader)
		}

		// Check accept header
		acceptHeader := r.Header.Get("Accept")
		if acceptHeader != "application/vnd.api+json" {
			t.Errorf("Expected Accept header to be 'application/vnd.api+json', got %q", acceptHeader)
		}

		_, _ = fmt.Fprintln(w, `{"data": []}`)
	}))
	defer server.Close()

	// Create client
	cfg := &config.Config{
		APIKey:     apiKey,
		Timeout:    5 * time.Second,
		APIBaseURL: server.URL,
	}
	client := NewClient(cfg)

	// Make a request that will trigger the header checks in the test server
	resp, err := client.makeRequest(context.Background(), "GET", "/routes", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	_ = resp.Body.Close()
}

func TestAuthentication(t *testing.T) {
	t.Run("With API Key", func(t *testing.T) {
		apiKey := "test-api-key"

		// Create a test server that verifies API key
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check API key header
			authHeader := r.Header.Get("X-API-Key")
			if authHeader != apiKey {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintln(w, `{"errors":[{"status":"401","code":"unauthorized","title":"Unauthorized request","detail":"API key missing or invalid"}]}`)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"data": []}`)
		}))
		defer server.Close()

		// Create client with API key
		cfg := &config.Config{
			APIKey:     apiKey,
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Make a request with valid API key
		resp, err := client.makeRequest(context.Background(), "GET", "/routes", nil)
		if err != nil {
			t.Fatalf("Request with valid API key failed: %v", err)
		}
		_ = resp.Body.Close()
	})

	t.Run("Without API Key", func(t *testing.T) {
		// Some endpoints work without API key but are rate-limited
		// Create a test server that accepts requests without API key
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if API key header is missing
			if r.Header.Get("X-API-Key") != "" {
				t.Errorf("Expected no X-API-Key header, but one was provided")
			}

			// Public endpoint still works but might have lower rate limits
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"data": []}`)
		}))
		defer server.Close()

		// Create client without API key
		cfg := &config.Config{
			APIKey:     "", // No API key
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Make a request without API key to a public endpoint
		resp, err := client.makeRequest(context.Background(), "GET", "/routes", nil)
		if err != nil {
			t.Fatalf("Request to public endpoint without API key failed: %v", err)
		}
		_ = resp.Body.Close()
	})

	t.Run("Invalid API Key", func(t *testing.T) {
		// Create a test server that rejects invalid API keys
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if API key is invalid
			if r.Header.Get("X-API-Key") == "invalid-key" {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintln(w, `{"errors":[{"status":"401","code":"unauthorized","title":"Unauthorized request","detail":"API key invalid"}]}`)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"data": []}`)
		}))
		defer server.Close()

		// Create client with invalid API key
		cfg := &config.Config{
			APIKey:     "invalid-key",
			Timeout:    5 * time.Second,
			APIBaseURL: server.URL,
		}
		client := NewClient(cfg)

		// Make a request with invalid API key
		_, err := client.makeRequest(context.Background(), "GET", "/routes", nil)
		if err == nil {
			t.Fatal("Expected error for invalid API key, got nil")
		}

		// Check that we got an API Error
		apiErr, ok := err.(*APIError)
		if !ok {
			t.Fatalf("Expected *APIError, got %T", err)
		}

		// Verify it's an auth error
		if !apiErr.IsAuthError() {
			t.Errorf("Expected auth error, got: %v", err)
		}
	})
}

// Helper function to check if a string contains a substring
// Currently unused but will be needed for future tests
// func contains(s, substr string) bool {
// 	return s != "" && len(substr) > 0 && strings.Contains(s, substr)
// }

package mock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMockServer(t *testing.T) {
	// Define test responses
	definitions := []ResponseDefinition{
		{
			Path:       "/test/path",
			Method:     http.MethodGet,
			StatusCode: http.StatusOK,
			Response:   `{"data": "test"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			Path:       "/another/path",
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Response:   `{"data": "created"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	// Create mock server
	server := NewMockServer(definitions)
	defer server.Close()

	// Test first path
	resp, err := http.Get(server.URL + "/test/path")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != `{"data": "test"}` {
		t.Errorf("Expected response body %q, got %q", `{"data": "test"}`, string(body))
	}

	// Test second path
	req, err := http.NewRequest(http.MethodPost, server.URL+"/another/path", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Test non-existent path (should return 404)
	resp, err = http.Get(server.URL + "/not/found")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestMockValidAPIKeyMiddleware(t *testing.T) {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap it with the middleware
	wrappedHandler := MockValidAPIKeyMiddleware(handler)

	// Create a test server
	server := httptest.NewServer(wrappedHandler)
	defer server.Close()

	// Test with no API key
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// Test with invalid API key
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "invalid-key")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// Test with valid API key
	req, err = http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "success" {
		t.Errorf("Expected response body %q, got %q", "success", string(body))
	}
}

func TestMockRateLimitMiddleware(t *testing.T) {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap it with the middleware
	wrappedHandler := MockRateLimitMiddleware(handler)

	// Create a test server
	server := httptest.NewServer(wrappedHandler)
	defer server.Close()

	// Test normal request
	resp, err := http.Get(server.URL + "/normal/path")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test rate-limited request
	resp, err = http.Get(server.URL + "/rate-limited/path")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, resp.StatusCode)
	}

	if resp.Header.Get("Retry-After") != "60" {
		t.Errorf("Expected Retry-After header to be '60', got '%s'", resp.Header.Get("Retry-After"))
	}
}

func TestStandardMockServer(t *testing.T) {
	server, err := StandardMockServer()
	if err != nil {
		t.Fatalf("Failed to create standard mock server: %v", err)
	}
	defer server.Close()

	// Test routes endpoint
	req, err := http.NewRequest(http.MethodGet, server.URL+"/routes", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test stops endpoint
	req, err = http.NewRequest(http.MethodGet, server.URL+"/stops", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test schedules endpoint
	req, err = http.NewRequest(http.MethodGet, server.URL+"/schedules", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test 404 endpoint
	req, err = http.NewRequest(http.MethodGet, server.URL+"/not/found", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("X-API-Key", "valid-key")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
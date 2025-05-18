package mbta

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: http.StatusNotFound,
		Status:     "404",
		Code:       "not_found",
		Title:      "Resource Not Found",
		Detail:     "The requested resource could not be found",
	}

	// Test Error() method
	expected := "MBTA API Error (404): Resource Not Found - The requested resource could not be found"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}

	// Test IsNotFoundError
	if !err.IsNotFoundError() {
		t.Error("Expected IsNotFoundError to return true for 404 error")
	}

	// Test IsAuthError
	if err.IsAuthError() {
		t.Error("Expected IsAuthError to return false for non-auth error")
	}

	// Test IsRateLimitError
	if err.IsRateLimitError() {
		t.Error("Expected IsRateLimitError to return false for non-rate-limit error")
	}
}

func TestRateLimitError(t *testing.T) {
	err := &RateLimitError{
		APIError: &APIError{
			StatusCode: http.StatusTooManyRequests,
			Status:     "429",
			Code:       "rate_limited",
			Title:      "Rate Limit Exceeded",
			Detail:     "You have exceeded your rate limit",
		},
		RetryAfter: 60,
	}

	// Test Error() method
	expected := "MBTA API Rate Limit Exceeded: You have exceeded your rate limit. Retry after 60 seconds"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}

	// Test IsRateLimitError
	if !err.IsRateLimitError() {
		t.Error("Expected IsRateLimitError to return true for 429 error")
	}
}

func TestTimeoutError(t *testing.T) {
	originalErr := errors.New("connection timed out")
	err := &TimeoutError{
		NetworkError: &NetworkError{Err: originalErr},
		Timeout:      30 * time.Second,
	}

	// Test that it contains timeout information
	errMsg := err.Error()
	if !strings.Contains(errMsg, "30s") || !strings.Contains(errMsg, "timed out") {
		t.Errorf("Expected error message to contain timeout information, got %q", errMsg)
	}
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("network connection refused")
	err := &NetworkError{Err: originalErr}

	// Test Error() method
	expected := "Network error: network connection refused"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestParseAPIError(t *testing.T) {
	// Test valid JSON error
	jsonBody := []byte(`{
		"errors": [
			{
				"status": "404",
				"code": "not_found",
				"title": "Resource Not Found",
				"detail": "The requested resource could not be found",
				"source": {
					"parameter": "id"
				}
			}
		]
	}`)

	err := parseAPIError(http.StatusNotFound, jsonBody)

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusNotFound, apiErr.StatusCode)
	}

	if apiErr.Status != "404" {
		t.Errorf("Expected Status %q, got %q", "404", apiErr.Status)
	}

	if apiErr.Code != "not_found" {
		t.Errorf("Expected Code %q, got %q", "not_found", apiErr.Code)
	}

	if apiErr.Title != "Resource Not Found" {
		t.Errorf("Expected Title %q, got %q", "Resource Not Found", apiErr.Title)
	}

	if apiErr.Detail != "The requested resource could not be found" {
		t.Errorf("Expected Detail %q, got %q", "The requested resource could not be found", apiErr.Detail)
	}

	source, ok := apiErr.Source["parameter"]
	if !ok || source != "id" {
		t.Errorf("Expected Source.parameter %q, got %v", "id", apiErr.Source)
	}

	// Test rate limit error
	jsonBody = []byte(`{
		"errors": [
			{
				"status": "429",
				"code": "rate_limited",
				"title": "Rate Limit Exceeded",
				"detail": "You have exceeded your rate limit"
			}
		]
	}`)

	err = parseAPIError(http.StatusTooManyRequests, jsonBody)

	rateLimitErr, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("Expected *RateLimitError, got %T", err)
	}

	if rateLimitErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusTooManyRequests, rateLimitErr.StatusCode)
	}

	if rateLimitErr.RetryAfter <= 0 {
		t.Errorf("Expected RetryAfter > 0, got %d", rateLimitErr.RetryAfter)
	}

	// Test invalid JSON
	invalidJSON := []byte(`{"errors": [{"status": 404}`)
	err = parseAPIError(http.StatusNotFound, invalidJSON)

	apiErr, ok = err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusNotFound, apiErr.StatusCode)
	}

	// Test empty errors array
	emptyErrors := []byte(`{"errors": []}`)
	err = parseAPIError(http.StatusBadRequest, emptyErrors)

	apiErr, ok = err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
}

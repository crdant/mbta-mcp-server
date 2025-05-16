package mbta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// APIError represents an error returned by the MBTA API
type APIError struct {
	StatusCode int
	Status     string
	Code       string
	Title      string
	Detail     string
	Source     map[string]interface{}
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return fmt.Sprintf("MBTA API Error (%s): %s - %s", e.Status, e.Title, e.Detail)
}

// IsNotFoundError checks if the error is a not found error
func (e *APIError) IsNotFoundError() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsAuthError checks if the error is an authentication error
func (e *APIError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsRateLimitError checks if the error is a rate limit error
func (e *APIError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// RateLimitError is a specific type of error for rate limiting
type RateLimitError struct {
	*APIError
	RetryAfter int // Seconds to wait before retrying
}

// Error implements the error interface for RateLimitError
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("MBTA API Rate Limit Exceeded: %s. Retry after %d seconds", e.Detail, e.RetryAfter)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Err error
}

// Error implements the error interface for NetworkError
func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error: %v", e.Err)
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	*NetworkError
	Timeout time.Duration
}

// Error implements the error interface for TimeoutError
func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Request timed out after %v: %v", e.Timeout, e.Err)
}

// parseAPIError parses an error response from the MBTA API
func parseAPIError(statusCode int, responseBody []byte) error {
	var errorResponse struct {
		Errors []struct {
			Status string `json:"status"`
			Code   string `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
			Source map[string]interface{} `json:"source,omitempty"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
		// If we can't parse the error response, return a generic error
		return &APIError{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d", statusCode),
			Title:      http.StatusText(statusCode),
			Detail:     string(responseBody),
		}
	}

	// If there are no errors in the response, return a generic error
	if len(errorResponse.Errors) == 0 {
		return &APIError{
			StatusCode: statusCode,
			Status:     fmt.Sprintf("%d", statusCode),
			Title:      http.StatusText(statusCode),
			Detail:     string(responseBody),
		}
	}

	// Get the first error
	apiErr := errorResponse.Errors[0]

	// Create the error
	err := &APIError{
		StatusCode: statusCode,
		Status:     apiErr.Status,
		Code:       apiErr.Code,
		Title:      apiErr.Title,
		Detail:     apiErr.Detail,
		Source:     apiErr.Source,
	}

	// Check for rate limit errors
	if statusCode == http.StatusTooManyRequests {
		// For now, default to 60 seconds if no retry-after header
		retryAfter := 60
		return &RateLimitError{
			APIError:   err,
			RetryAfter: retryAfter,
		}
	}

	return err
}
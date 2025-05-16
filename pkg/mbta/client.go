// Package mbta provides a client for the MBTA API v3.
package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/crdant/mbta-mcp-server/internal/config"
	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

// Client represents an MBTA API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new MBTA API client with the provided configuration
func NewClient(cfg *config.Config) *Client {
	// Create transport with sensible defaults
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   10,
	}

	// Create HTTP client with timeout from config
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	return &Client{
		baseURL:    cfg.APIBaseURL,
		apiKey:     cfg.APIKey,
		httpClient: httpClient,
	}
}

// makeRequest performs an HTTP request with proper headers and handles authentication
func (c *Client) makeRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, &NetworkError{Err: fmt.Errorf("error creating request: %w", err)}
	}

	// Set common headers
	req.Header.Set("Accept", "application/vnd.api+json")

	// Set API key if available
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	// Perform the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Check if it's a timeout error
		if isTimeoutError(err) {
			return nil, &TimeoutError{
				NetworkError: &NetworkError{Err: err},
				Timeout:      c.httpClient.Timeout,
			}
		}
		return nil, &NetworkError{Err: fmt.Errorf("error performing request: %w", err)}
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		// Read error response body
		respBody, readErr := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		if readErr != nil {
			return nil, &NetworkError{Err: fmt.Errorf("HTTP error %d and failed to read error body: %w", resp.StatusCode, readErr)}
		}

		// Parse the API error
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	// Successfully processed the request
	return resp, nil
}

// isTimeoutError checks if an error is a timeout error
func isTimeoutError(err error) bool {
	if err, ok := err.(interface{ Timeout() bool }); ok && err.Timeout() {
		return true
	}
	return false
}

// GetRoutes retrieves all available MBTA routes
func (c *Client) GetRoutes(ctx context.Context) ([]models.Route, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/routes", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var routeResponse models.RouteResponse
	if err := json.NewDecoder(resp.Body).Decode(&routeResponse); err != nil {
		return nil, fmt.Errorf("error decoding route response: %w", err)
	}

	return routeResponse.Data, nil
}

// GetRoute retrieves a specific MBTA route by ID
func (c *Client) GetRoute(ctx context.Context, routeID string) (*models.Route, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf("/routes/%s", routeID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var routeData struct {
		Data models.Route `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&routeData); err != nil {
		return nil, fmt.Errorf("error decoding route response: %w", err)
	}

	return &routeData.Data, nil
}

// GetStops retrieves all available MBTA stops
func (c *Client) GetStops(ctx context.Context) ([]models.Stop, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, "/stops", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var stopResponse models.StopResponse
	if err := json.NewDecoder(resp.Body).Decode(&stopResponse); err != nil {
		return nil, fmt.Errorf("error decoding stop response: %w", err)
	}

	return stopResponse.Data, nil
}

// GetStop retrieves a specific MBTA stop by ID
func (c *Client) GetStop(ctx context.Context, stopID string) (*models.Stop, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf("/stops/%s", stopID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var stopData struct {
		Data models.Stop `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&stopData); err != nil {
		return nil, fmt.Errorf("error decoding stop response: %w", err)
	}

	return &stopData.Data, nil
}

// GetSchedules retrieves schedules by route, stop, or trip ID
func (c *Client) GetSchedules(ctx context.Context, params map[string]string) ([]models.Schedule, []models.Included, error) {
	// Build query parameters
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	path := "/schedules"
	if queryString := query.Encode(); queryString != "" {
		path += "?" + queryString
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var scheduleResponse models.ScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&scheduleResponse); err != nil {
		return nil, nil, fmt.Errorf("error decoding schedule response: %w", err)
	}

	return scheduleResponse.Data, scheduleResponse.Included, nil
}
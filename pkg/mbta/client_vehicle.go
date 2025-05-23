package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

// GetVehicles retrieves all available MBTA vehicles with optional filtering
func (c *Client) GetVehicles(ctx context.Context, params map[string]string) ([]models.Vehicle, error) {
	// Build query parameters
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	path := "/vehicles"
	if queryString := query.Encode(); queryString != "" {
		path += "?" + queryString
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	var vehicleResponse models.VehicleResponse
	if err := json.NewDecoder(resp.Body).Decode(&vehicleResponse); err != nil {
		return nil, fmt.Errorf("error decoding vehicle response: %w", err)
	}

	return vehicleResponse.Data, nil
}

// GetVehicle retrieves a specific MBTA vehicle by ID
func (c *Client) GetVehicle(ctx context.Context, vehicleID string) (*models.Vehicle, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf("/vehicles/%s", vehicleID), nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	var vehicleData struct {
		Data models.Vehicle `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&vehicleData); err != nil {
		return nil, fmt.Errorf("error decoding vehicle response: %w", err)
	}

	return &vehicleData.Data, nil
}

// GetVehiclesByRoute retrieves all vehicles for a specific route
func (c *Client) GetVehiclesByRoute(ctx context.Context, routeID string) ([]models.Vehicle, error) {
	params := map[string]string{
		"filter[route]": routeID,
	}
	return c.GetVehicles(ctx, params)
}

// GetVehiclesByTrip retrieves all vehicles for a specific trip
func (c *Client) GetVehiclesByTrip(ctx context.Context, tripID string) ([]models.Vehicle, error) {
	params := map[string]string{
		"filter[trip]": tripID,
	}
	return c.GetVehicles(ctx, params)
}

// GetVehiclesByLocation retrieves all vehicles near a specific location
func (c *Client) GetVehiclesByLocation(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Vehicle, error) {
	// Default radius to 0.01 if not specified or negative
	if radius <= 0 {
		radius = 0.01
	}

	params := map[string]string{
		"filter[latitude]":  fmt.Sprintf("%f", latitude),
		"filter[longitude]": fmt.Sprintf("%f", longitude),
		"filter[radius]":    fmt.Sprintf("%f", radius),
	}
	return c.GetVehicles(ctx, params)
}
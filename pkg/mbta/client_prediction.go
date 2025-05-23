package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

// GetPredictions retrieves all available MBTA predictions with optional filtering
func (c *Client) GetPredictions(ctx context.Context, params map[string]string) ([]models.Prediction, error) {
	// Build query parameters
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	path := "/predictions"
	if queryString := query.Encode(); queryString != "" {
		path += "?" + queryString
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	var predictionResponse models.PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predictionResponse); err != nil {
		return nil, fmt.Errorf("error decoding prediction response: %w", err)
	}

	return predictionResponse.Data, nil
}

// GetPredictionsByVehicle retrieves predictions for a specific vehicle
func (c *Client) GetPredictionsByVehicle(ctx context.Context, vehicleID string) ([]models.Prediction, error) {
	params := map[string]string{
		"filter[vehicle]": vehicleID,
	}
	return c.GetPredictions(ctx, params)
}

// GetPredictionsByRoute retrieves predictions for a specific route
func (c *Client) GetPredictionsByRoute(ctx context.Context, routeID string) ([]models.Prediction, error) {
	params := map[string]string{
		"filter[route]": routeID,
	}
	return c.GetPredictions(ctx, params)
}

// GetPredictionsByStop retrieves predictions for a specific stop
func (c *Client) GetPredictionsByStop(ctx context.Context, stopID string) ([]models.Prediction, error) {
	params := map[string]string{
		"filter[stop]": stopID,
	}
	return c.GetPredictions(ctx, params)
}

// GetPredictionsByTrip retrieves predictions for a specific trip
func (c *Client) GetPredictionsByTrip(ctx context.Context, tripID string) ([]models.Prediction, error) {
	params := map[string]string{
		"filter[trip]": tripID,
	}
	return c.GetPredictions(ctx, params)
}

// GetPredictionsByLocation retrieves predictions for stops near a specific location
func (c *Client) GetPredictionsByLocation(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Prediction, error) {
	// Default radius to 0.01 if not specified or negative
	if radius <= 0 {
		radius = 0.01
	}

	params := map[string]string{
		"filter[latitude]":  fmt.Sprintf("%f", latitude),
		"filter[longitude]": fmt.Sprintf("%f", longitude),
		"filter[radius]":    fmt.Sprintf("%f", radius),
	}
	return c.GetPredictions(ctx, params)
}
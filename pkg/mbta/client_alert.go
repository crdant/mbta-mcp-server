package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

// GetAlerts retrieves all available MBTA alerts with optional filtering
func (c *Client) GetAlerts(ctx context.Context, params map[string]string) ([]models.Alert, error) {
	// Build query parameters
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	path := "/alerts"
	if queryString := query.Encode(); queryString != "" {
		path += "?" + queryString
	}

	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	var alertResponse models.AlertResponse
	if err := json.NewDecoder(resp.Body).Decode(&alertResponse); err != nil {
		return nil, fmt.Errorf("error decoding alert response: %w", err)
	}

	return alertResponse.Data, nil
}

// GetAlert retrieves a specific MBTA alert by ID
func (c *Client) GetAlert(ctx context.Context, alertID string) (*models.Alert, error) {
	resp, err := c.makeRequest(ctx, http.MethodGet, fmt.Sprintf("/alerts/%s", alertID), nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse response
	var alertData struct {
		Data models.Alert `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&alertData); err != nil {
		return nil, fmt.Errorf("error decoding alert response: %w", err)
	}

	return &alertData.Data, nil
}

// GetActiveAlerts retrieves all currently active alerts
func (c *Client) GetActiveAlerts(ctx context.Context) ([]models.Alert, error) {
	params := map[string]string{
		"filter[activity]": "USING_ACCESSORY,BOARD,EXIT,PARK_CAR,RIDE,STORE_BIKE,USING_ESCALATOR,USING_WHEELCHAIR",
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByRoute retrieves all alerts for a specific route
func (c *Client) GetAlertsByRoute(ctx context.Context, routeID string) ([]models.Alert, error) {
	params := map[string]string{
		"filter[route]": routeID,
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByStop retrieves all alerts for a specific stop
func (c *Client) GetAlertsByStop(ctx context.Context, stopID string) ([]models.Alert, error) {
	params := map[string]string{
		"filter[stop]": stopID,
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByTrip retrieves all alerts for a specific trip
func (c *Client) GetAlertsByTrip(ctx context.Context, tripID string) ([]models.Alert, error) {
	params := map[string]string{
		"filter[trip]": tripID,
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByRoutes retrieves all alerts for a list of routes
func (c *Client) GetAlertsByRoutes(ctx context.Context, routeIDs []string) ([]models.Alert, error) {
	params := map[string]string{
		"filter[route]": strings.Join(routeIDs, ","),
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByStops retrieves all alerts for a list of stops
func (c *Client) GetAlertsByStops(ctx context.Context, stopIDs []string) ([]models.Alert, error) {
	params := map[string]string{
		"filter[stop]": strings.Join(stopIDs, ","),
	}
	return c.GetAlerts(ctx, params)
}

// GetAlertsByEffect retrieves all alerts with a specific effect
func (c *Client) GetAlertsByEffect(ctx context.Context, effect models.AlertEffect) ([]models.Alert, error) {
	params := map[string]string{
		"filter[effect]": string(effect),
	}
	return c.GetAlerts(ctx, params)
}

// GetServiceDisruptions retrieves all alerts that represent significant service disruptions
func (c *Client) GetServiceDisruptions(ctx context.Context) ([]models.Alert, error) {
	// These are the effects that represent significant service disruptions
	disruptions := []models.AlertEffect{
		models.AlertEffectNoService,
		models.AlertEffectReducedService,
		models.AlertEffectSignificantDelays,
		models.AlertEffectDetour,
		models.AlertEffectStationClosure,
		models.AlertEffectShuttle,
	}

	// Convert effects to comma-separated string
	effectStrings := make([]string, 0, len(disruptions))
	for _, effect := range disruptions {
		effectStrings = append(effectStrings, string(effect))
	}

	params := map[string]string{
		"filter[effect]": strings.Join(effectStrings, ","),
	}

	return c.GetAlerts(ctx, params)
}

// GetAccessibilityAlerts retrieves all alerts related to accessibility issues
func (c *Client) GetAccessibilityAlerts(ctx context.Context) ([]models.Alert, error) {
	params := map[string]string{
		"filter[effect]": string(models.AlertEffectAccessibilityIssue) + "," + string(models.AlertEffectElevatorOutage),
	}
	return c.GetAlerts(ctx, params)
}

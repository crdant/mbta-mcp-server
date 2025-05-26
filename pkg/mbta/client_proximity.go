package mbta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/crdant/mbta-mcp-server/pkg/mbta/models"
)

// FindNearbyStations finds stations near the specified coordinates
// lat, lon: Coordinates to search around
// radius: Maximum distance in kilometers
// maxResults: Maximum number of results to return
// onlyStations: If true, only return full stations (not platforms or stops)
func (c *Client) FindNearbyStations(ctx context.Context, lat, lon, radius float64, maxResults int, onlyStations bool) ([]models.NearbyStation, error) {
	// Build query parameters
	query := url.Values{}
	
	// If only interested in stations, filter by location_type=1
	if onlyStations {
		query.Add("filter[location_type]", strconv.Itoa(models.LocationTypeStation))
	}

	// Make the request to get all stops
	path := "/stops?" + query.Encode()
	resp, err := c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get stations: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var stopResponse models.StopResponse
	if err := json.NewDecoder(resp.Body).Decode(&stopResponse); err != nil {
		return nil, fmt.Errorf("error decoding station response: %w", err)
	}

	// Calculate distance to each stop and filter by radius
	var nearbyStations []models.NearbyStation
	for _, stop := range stopResponse.Data {
		distance := calculateApproximateDistance(
			lat, lon,
			stop.Attributes.Latitude, stop.Attributes.Longitude,
		)

		// Only include stops within the radius
		if distance <= radius {
			nearbyStations = append(nearbyStations, models.NearbyStation{
				Stop:       stop,
				DistanceKm: distance,
			})
		}
	}

	// Sort by distance (closest first)
	sort.Slice(nearbyStations, func(i, j int) bool {
		return nearbyStations[i].DistanceKm < nearbyStations[j].DistanceKm
	})

	// Limit to maxResults
	if len(nearbyStations) > maxResults && maxResults > 0 {
		nearbyStations = nearbyStations[:maxResults]
	}

	return nearbyStations, nil
}
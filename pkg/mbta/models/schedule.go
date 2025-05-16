// Package models contains data models for MBTA API responses
package models

import (
	"time"
)

// ScheduleResponse represents a response containing schedule data from the MBTA API
type ScheduleResponse struct {
	Data     []Schedule   `json:"data"`
	Included []Included   `json:"included,omitempty"`
}

// Schedule represents a transit vehicle schedule (stop time) in the MBTA system
type Schedule struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    ScheduleAttributes     `json:"attributes"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// ScheduleAttributes contains the attributes of a schedule
type ScheduleAttributes struct {
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	DropOffType   int    `json:"drop_off_type"`
	PickupType    int    `json:"pickup_type"`
	StopHeadsign  string `json:"stop_headsign"`
	StopSequence  int    `json:"stop_sequence"`
	Timepoint     bool   `json:"timepoint"`
}

// Included represents an included object in the MBTA API response
type Included struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Attributes    interface{} `json:"attributes"`
	Relationships interface{} `json:"relationships,omitempty"`
}

// Pickup and drop-off type constants as defined by the MBTA API
const (
	PickupDropOffRegular = 0
	PickupDropOffNotAvailable = 1
	PickupDropOffPhoneAgency = 2
	PickupDropOffCoordinateWithDriver = 3
)

// GetDuration returns the duration between arrival and departure times
func (s *Schedule) GetDuration() (time.Duration, error) {
	arrivalTime, err := time.Parse(time.RFC3339, s.Attributes.ArrivalTime)
	if err != nil {
		return 0, err
	}

	departureTime, err := time.Parse(time.RFC3339, s.Attributes.DepartureTime)
	if err != nil {
		return 0, err
	}

	return departureTime.Sub(arrivalTime), nil
}

// FormattedArrivalTime returns the arrival time formatted according to the given layout
func (s *Schedule) FormattedArrivalTime(layout string) (string, error) {
	t, err := time.Parse(time.RFC3339, s.Attributes.ArrivalTime)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// FormattedDepartureTime returns the departure time formatted according to the given layout
func (s *Schedule) FormattedDepartureTime(layout string) (string, error) {
	t, err := time.Parse(time.RFC3339, s.Attributes.DepartureTime)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// IsPickupAvailable returns whether pickup is available at this stop
func (s *Schedule) IsPickupAvailable() bool {
	return s.Attributes.PickupType == PickupDropOffRegular
}

// IsDropOffAvailable returns whether drop-off is available at this stop
func (s *Schedule) IsDropOffAvailable() bool {
	return s.Attributes.DropOffType == PickupDropOffRegular
}

// IsTimepoint returns whether this stop is a timepoint
// Timepoints are stops where the schedule is strictly adhered to
func (s *Schedule) IsTimepoint() bool {
	return s.Attributes.Timepoint
}
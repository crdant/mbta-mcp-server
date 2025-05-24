// Package models contains data models for MBTA API responses
package models

import (
	"time"
)

// AlertResponse represents a response containing alert data from the MBTA API
type AlertResponse struct {
	Data     []Alert    `json:"data"`
	Included []Included `json:"included,omitempty"`
}

// Alert represents a service alert in the MBTA system
type Alert struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Attributes    AlertAttributes        `json:"attributes"`
	Relationships map[string]interface{} `json:"relationships,omitempty"`
}

// AlertAttributes contains the attributes of an alert
type AlertAttributes struct {
	ActivePeriod   []AlertPeriod `json:"active_period"`
	Banner         bool          `json:"banner"`
	Cause          AlertCause    `json:"cause"`
	CreatedAt      time.Time     `json:"created_at"`
	Description    string        `json:"description"`
	Effect         AlertEffect   `json:"effect"`
	Header         string        `json:"header"`
	InformedEntity []AlertEntity `json:"informed_entity"`
	Lifecycle      string        `json:"lifecycle"`
	Severity       int           `json:"severity"`
	ServiceEffect  string        `json:"service_effect"`
	Timeframe      string        `json:"timeframe,omitempty"`
	UpdatedAt      time.Time     `json:"updated_at"`
	URL            string        `json:"url,omitempty"`
}

// AlertPeriod represents a time period when an alert is active
type AlertPeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// AlertEntity represents an entity affected by an alert
type AlertEntity struct {
	Activities  []string `json:"activities"`
	DirectionID *int     `json:"direction_id,omitempty"`
	Route       string   `json:"route,omitempty"`
	RouteType   int      `json:"route_type,omitempty"`
	Stop        string   `json:"stop,omitempty"`
	Trip        string   `json:"trip,omitempty"`
	FacilityID  string   `json:"facility,omitempty"`
}

// AlertEffect is the type of effect an alert has on service
type AlertEffect string

// Alert effect constants
const (
	AlertEffectNoService          AlertEffect = "NO_SERVICE"
	AlertEffectReducedService     AlertEffect = "REDUCED_SERVICE"
	AlertEffectSignificantDelays  AlertEffect = "SIGNIFICANT_DELAYS"
	AlertEffectDelays             AlertEffect = "DELAYS"
	AlertEffectDetour             AlertEffect = "DETOUR"
	AlertEffectStopMoved          AlertEffect = "STOP_MOVED"
	AlertEffectStopClosed         AlertEffect = "STOP_CLOSED"
	AlertEffectShuttle            AlertEffect = "SHUTTLE"
	AlertEffectElevatorOutage     AlertEffect = "ELEVATOR_OUTAGE"
	AlertEffectAccessibilityIssue AlertEffect = "ACCESSIBILITY_ISSUE"
	AlertEffectScheduleChange     AlertEffect = "SCHEDULE_CHANGE"
	AlertEffectServiceChange      AlertEffect = "SERVICE_CHANGE"
	AlertEffectSnowRoute          AlertEffect = "SNOW_ROUTE"
	AlertEffectStationClosure     AlertEffect = "STATION_CLOSURE"
	AlertEffectTrackChange        AlertEffect = "TRACK_CHANGE"
	AlertEffectAdditionalService  AlertEffect = "ADDITIONAL_SERVICE"
	AlertEffectModifiedService    AlertEffect = "MODIFIED_SERVICE"
	AlertEffectOther              AlertEffect = "OTHER_EFFECT"
)

// AlertCause is the cause of an alert
type AlertCause string

// Alert cause constants
const (
	AlertCauseUnknownCause       AlertCause = "UNKNOWN_CAUSE"
	AlertCauseUnspecifiedCause   AlertCause = "UNSPECIFIED_CAUSE"
	AlertCauseAccident           AlertCause = "ACCIDENT"
	AlertCauseConstruction       AlertCause = "CONSTRUCTION"
	AlertCauseDemonstation       AlertCause = "DEMONSTRATION"
	AlertCauseEquipmentFailure   AlertCause = "EQUIPMENT_FAILURE"
	AlertCauseMedicalEmergency   AlertCause = "MEDICAL_EMERGENCY"
	AlertCausePoliceActivity     AlertCause = "POLICE_ACTIVITY"
	AlertCauseMaintenance        AlertCause = "MAINTENANCE"
	AlertCauseWeather            AlertCause = "WEATHER"
	AlertCauseTrafficCongestion  AlertCause = "TRAFFIC_CONGESTION"
	AlertCauseFireActivity       AlertCause = "FIRE"
	AlertCauseHoliday            AlertCause = "HOLIDAY"
	AlertCauseStrike             AlertCause = "STRIKE"
	AlertCauseSuspiciousActivity AlertCause = "SUSPICIOUS_ACTIVITY"
	AlertCauseSwitchFailure      AlertCause = "SWITCH_FAILURE"
	AlertCauseOther              AlertCause = "OTHER_CAUSE"
)

// GetAlertEffectDescription returns a human-readable description of an alert effect
func GetAlertEffectDescription(effect AlertEffect) string {
	switch effect {
	case AlertEffectNoService:
		return "No Service"
	case AlertEffectReducedService:
		return "Reduced Service"
	case AlertEffectSignificantDelays:
		return "Significant Delays"
	case AlertEffectDelays:
		return "Delays"
	case AlertEffectDetour:
		return "Detour"
	case AlertEffectStopMoved:
		return "Stop Moved"
	case AlertEffectStopClosed:
		return "Stop Closed"
	case AlertEffectShuttle:
		return "Shuttle Bus Service"
	case AlertEffectElevatorOutage:
		return "Elevator Outage"
	case AlertEffectAccessibilityIssue:
		return "Accessibility Issue"
	case AlertEffectScheduleChange:
		return "Schedule Change"
	case AlertEffectServiceChange:
		return "Service Change"
	case AlertEffectSnowRoute:
		return "Snow Route in Effect"
	case AlertEffectStationClosure:
		return "Station Closure"
	case AlertEffectTrackChange:
		return "Track Change"
	case AlertEffectAdditionalService:
		return "Additional Service"
	case AlertEffectModifiedService:
		return "Modified Service"
	case AlertEffectOther:
		return "Other Effect"
	default:
		return "Unknown Effect"
	}
}

// GetAlertCauseDescription returns a human-readable description of an alert cause
func GetAlertCauseDescription(cause AlertCause) string {
	switch cause {
	case AlertCauseUnknownCause:
		return "Unknown Cause"
	case AlertCauseUnspecifiedCause:
		return "Unspecified Cause"
	case AlertCauseAccident:
		return "Accident"
	case AlertCauseConstruction:
		return "Construction"
	case AlertCauseDemonstation:
		return "Demonstration"
	case AlertCauseEquipmentFailure:
		return "Equipment Failure"
	case AlertCauseMedicalEmergency:
		return "Medical Emergency"
	case AlertCausePoliceActivity:
		return "Police Activity"
	case AlertCauseMaintenance:
		return "Maintenance"
	case AlertCauseWeather:
		return "Weather"
	case AlertCauseTrafficCongestion:
		return "Traffic Congestion"
	case AlertCauseFireActivity:
		return "Fire Activity"
	case AlertCauseHoliday:
		return "Holiday Schedule"
	case AlertCauseStrike:
		return "Strike"
	case AlertCauseSuspiciousActivity:
		return "Suspicious Activity"
	case AlertCauseSwitchFailure:
		return "Switch Failure"
	case AlertCauseOther:
		return "Other Cause"
	default:
		return "Unknown Cause"
	}
}

// GetSeverityDescription returns a human-readable description of an alert severity level
func GetSeverityDescription(severity int) string {
	switch severity {
	case 1:
		return "Information"
	case 3:
		return "Minor Impact"
	case 5:
		return "Moderate Impact"
	case 7:
		return "Severe Impact"
	case 9:
		return "Critical Impact"
	default:
		return "Unknown Severity"
	}
}

// IsActive checks if the alert is currently active based on its active periods
func (a *Alert) IsActive(currentTime time.Time) bool {
	if len(a.Attributes.ActivePeriod) == 0 {
		return false
	}

	for _, period := range a.Attributes.ActivePeriod {
		// If start is zero time, treat as immediate start
		hasStarted := period.Start.IsZero() || !currentTime.Before(period.Start)

		// If end is zero time, treat as indefinite end
		hasNotEnded := period.End.IsZero() || !currentTime.After(period.End)

		if hasStarted && hasNotEnded {
			return true
		}
	}

	return false
}

// GetAffectedRoutes returns a list of route IDs affected by this alert
func (a *Alert) GetAffectedRoutes() []string {
	routeMap := make(map[string]bool)

	for _, entity := range a.Attributes.InformedEntity {
		if entity.Route != "" {
			routeMap[entity.Route] = true
		}
	}

	routes := make([]string, 0, len(routeMap))
	for route := range routeMap {
		routes = append(routes, route)
	}

	return routes
}

// GetAffectedStops returns a list of stop IDs affected by this alert
func (a *Alert) GetAffectedStops() []string {
	stopMap := make(map[string]bool)

	for _, entity := range a.Attributes.InformedEntity {
		if entity.Stop != "" {
			stopMap[entity.Stop] = true
		}
	}

	stops := make([]string, 0, len(stopMap))
	for stop := range stopMap {
		stops = append(stops, stop)
	}

	return stops
}

// HasActivity checks if the alert applies to a specific activity (e.g., BOARD, EXIT, RIDE, USING_WHEELCHAIR)
func (a *Alert) HasActivity(activity string) bool {
	for _, entity := range a.Attributes.InformedEntity {
		for _, act := range entity.Activities {
			if act == activity {
				return true
			}
		}
	}
	return false
}

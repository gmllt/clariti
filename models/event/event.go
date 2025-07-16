package event

import (
	"time"

	"github.com/gmllt/clariti/models/component"
)

// TypeEvent represents the type of an event
type TypeEvent string

const (
	TypePlannedMaintenance TypeEvent = "planned"
	TypeFiringIncident     TypeEvent = "firing"
	TypeKnownIssue         TypeEvent = "known_issue"
)

// Status represents the current status of an event
type Status string

const (
	StatusPlanned      Status = "planned"
	StatusOnGoing      Status = "ongoing"
	StatusResolved     Status = "resolved"
	StatusAcknowledged Status = "acknowledged"
	StatusCanceled     Status = "canceled"
	StatusUnknown      Status = "unknown"
)

// Criticality represents the severity level of an event
type Criticality int

const (
	CriticalityOperational      Criticality = 0
	CriticalityDegraded         Criticality = 1
	CriticalityPartialOutage    Criticality = 2
	CriticalityMajorOutage      Criticality = 3
	CriticalityUnderMaintenance Criticality = 4
	CriticalityUnknown          Criticality = -1
)

// String returns the string representation of the criticality
func (c Criticality) String() string {
	switch c {
	case CriticalityOperational:
		return "operational"
	case CriticalityDegraded:
		return "degraded"
	case CriticalityPartialOutage:
		return "partial outage"
	case CriticalityMajorOutage:
		return "major outage"
	case CriticalityUnderMaintenance:
		return "under maintenance"
	default:
		return "unknown"
	}
}

// Event defines the contract for all event types
type Event interface {
	Type() TypeEvent
	Status() Status
	Criticality() Criticality
}

// BaseEvent provides common fields for all event types
type BaseEvent struct {
	GUID           string                 `json:"guid"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	ExtraFields    map[string]string      `json:"extra_fields"`
	Components     []*component.Component `json:"components,omitempty"`
	StartEffective *time.Time             `json:"start_effective"`
	EndEffective   *time.Time             `json:"end_effective"`
}

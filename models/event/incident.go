package event

import (
	"time"
)

// Incident represents a service incident or known issue
type Incident struct {
	BaseEvent
	Perpetual           bool        `json:"perpetual,omitempty"` // Indicates if the incident is perpetual (known issue)
	IncidentCriticality Criticality `json:"criticality,omitempty"`
}

// Type returns the event type based on whether it's perpetual or not
func (i *Incident) Type() TypeEvent {
	if i.Perpetual {
		return TypeKnownIssue
	}
	return TypeFiringIncident
}

// Status returns the current status based on timing
func (i *Incident) Status() Status {
	if i.EndEffective != nil && i.EndEffective.Before(time.Now()) {
		return StatusResolved
	}
	if i.StartEffective != nil && i.StartEffective.Before(time.Now()) {
		return StatusOnGoing
	}
	return StatusUnknown
}

// Criticality returns the incident criticality level
func (i *Incident) Criticality() Criticality {
	return i.IncidentCriticality
}

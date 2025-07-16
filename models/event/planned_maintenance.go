package event

import (
	"time"
)

// PlannedMaintenance represents a scheduled maintenance event
type PlannedMaintenance struct {
	BaseEvent
	StartPlanned time.Time `json:"start_planned"`
	EndPlanned   time.Time `json:"end_planned"`
	Cancelled    bool      `json:"cancelled,omitempty"`
}

// Type returns the event type for planned maintenance
func (pm *PlannedMaintenance) Type() TypeEvent {
	return TypePlannedMaintenance
}

// Status returns the current status based on timing
func (pm *PlannedMaintenance) Status() Status {
	if pm.Cancelled {
		return StatusCanceled
	}
	if pm.EndEffective != nil && pm.EndEffective.Before(time.Now()) {
		return StatusResolved
	}
	if pm.StartEffective != nil && pm.StartEffective.Before(time.Now()) {
		return StatusOnGoing
	}
	if pm.StartPlanned.After(time.Now()) {
		return StatusPlanned
	}
	return StatusUnknown
}

// Criticality returns the criticality level for maintenance events
func (pm *PlannedMaintenance) Criticality() Criticality {
	return CriticalityUnderMaintenance
}

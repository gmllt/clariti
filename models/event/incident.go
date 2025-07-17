package event

import (
	"time"

	"github.com/gmllt/clariti/models/component"
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

// NewIncident creates a new incident with automatically generated GUID
func NewIncident(title, content string, components []*component.Component, criticality Criticality, perpetual bool) *Incident {
	return &Incident{
		BaseEvent:           NewBaseEvent(title, content, components),
		Perpetual:           perpetual,
		IncidentCriticality: criticality,
	}
}

// NewFiringIncident creates a new firing incident (non-perpetual)
func NewFiringIncident(title, content string, components []*component.Component, criticality Criticality) *Incident {
	return NewIncident(title, content, components, criticality, false)
}

// NewKnownIssue creates a new known issue (perpetual incident)
func NewKnownIssue(title, content string, components []*component.Component, criticality Criticality) *Incident {
	return NewIncident(title, content, components, criticality, true)
}

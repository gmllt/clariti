package drivers

import (
	"github.com/gmllt/clariti/models/event"
)

// EventStorage defines the interface for event storage drivers
type EventStorage interface {
	// Incidents
	CreateIncident(incident *event.Incident) error
	GetIncident(id string) (*event.Incident, error)
	GetAllIncidents() ([]*event.Incident, error)
	UpdateIncident(incident *event.Incident) error
	DeleteIncident(id string) error

	// Planned Maintenances
	CreatePlannedMaintenance(pm *event.PlannedMaintenance) error
	GetPlannedMaintenance(id string) (*event.PlannedMaintenance, error)
	GetAllPlannedMaintenances() ([]*event.PlannedMaintenance, error)
	UpdatePlannedMaintenance(pm *event.PlannedMaintenance) error
	DeletePlannedMaintenance(id string) error
}

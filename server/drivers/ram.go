package drivers

import (
	"errors"
	"sync"

	"github.com/gmllt/clariti/models/event"
)

var (
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("already exists")
)

// RAMStorage implements EventStorage interface using in-memory storage
type RAMStorage struct {
	mu                  sync.RWMutex
	incidents           map[string]*event.Incident
	plannedMaintenances map[string]*event.PlannedMaintenance
}

// NewRAMStorage creates a new in-memory storage driver
func NewRAMStorage() *RAMStorage {
	return &RAMStorage{
		incidents:           make(map[string]*event.Incident),
		plannedMaintenances: make(map[string]*event.PlannedMaintenance),
	}
}

// Incidents implementation
func (r *RAMStorage) CreateIncident(incident *event.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.incidents[incident.GUID]; exists {
		return ErrExists
	}
	r.incidents[incident.GUID] = incident
	return nil
}

func (r *RAMStorage) GetIncident(id string) (*event.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	incident, exists := r.incidents[id]
	if !exists {
		return nil, ErrNotFound
	}
	return incident, nil
}

func (r *RAMStorage) GetAllIncidents() ([]*event.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	incidents := make([]*event.Incident, 0, len(r.incidents))
	for _, incident := range r.incidents {
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (r *RAMStorage) UpdateIncident(incident *event.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.incidents[incident.GUID]; !exists {
		return ErrNotFound
	}
	r.incidents[incident.GUID] = incident
	return nil
}

func (r *RAMStorage) DeleteIncident(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.incidents[id]; !exists {
		return ErrNotFound
	}
	delete(r.incidents, id)
	return nil
}

// Planned Maintenances implementation
func (r *RAMStorage) CreatePlannedMaintenance(pm *event.PlannedMaintenance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plannedMaintenances[pm.GUID]; exists {
		return ErrExists
	}
	r.plannedMaintenances[pm.GUID] = pm
	return nil
}

func (r *RAMStorage) GetPlannedMaintenance(id string) (*event.PlannedMaintenance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pm, exists := r.plannedMaintenances[id]
	if !exists {
		return nil, ErrNotFound
	}
	return pm, nil
}

func (r *RAMStorage) GetAllPlannedMaintenances() ([]*event.PlannedMaintenance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pms := make([]*event.PlannedMaintenance, 0, len(r.plannedMaintenances))
	for _, pm := range r.plannedMaintenances {
		pms = append(pms, pm)
	}
	return pms, nil
}

func (r *RAMStorage) UpdatePlannedMaintenance(pm *event.PlannedMaintenance) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plannedMaintenances[pm.GUID]; !exists {
		return ErrNotFound
	}
	r.plannedMaintenances[pm.GUID] = pm
	return nil
}

func (r *RAMStorage) DeletePlannedMaintenance(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plannedMaintenances[id]; !exists {
		return ErrNotFound
	}
	delete(r.plannedMaintenances, id)
	return nil
}

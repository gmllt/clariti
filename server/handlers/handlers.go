package handlers

import (
	"github.com/gmllt/clariti/server/config"
	"github.com/gmllt/clariti/server/drivers"
)

// Handlers holds all the application handlers
type Handlers struct {
	API                *APIHandler
	Incident           *IncidentHandler
	PlannedMaintenance *PlannedMaintenanceHandler
	Weather            *WeatherHandler
}

// New creates a new handlers instance with all sub-handlers
func New(storage drivers.EventStorage, config *config.Config) *Handlers {
	// Cast to RAMStorage for weather service (we know it's RAMStorage)
	ramStorage := storage.(*drivers.RAMStorage)

	return &Handlers{
		API:                NewAPIHandler(storage, config),
		Incident:           NewIncidentHandler(storage),
		PlannedMaintenance: NewPlannedMaintenanceHandler(storage),
		Weather:            NewWeatherHandler(config, ramStorage),
	}
}

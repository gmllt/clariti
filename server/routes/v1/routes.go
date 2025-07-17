package v1

import (
	"net/http"

	"github.com/gmllt/clariti/server/handlers"
)

// SetupV1Routes configures all v1 API routes
func SetupV1Routes(mux *http.ServeMux, h *handlers.Handlers) {
	// Component information routes (read-only)
	setupV1ComponentRoutes(mux, h)

	// Event management routes
	setupV1EventRoutes(mux, h)

	// Weather/status routes
	setupV1WeatherRoutes(mux, h)
} // setupV1ComponentRoutes configures component information endpoints for v1
func setupV1ComponentRoutes(mux *http.ServeMux, h *handlers.Handlers) {
	// Component information endpoints (read-only, no auth needed)
	mux.HandleFunc("/api/v1/components", h.API.HandleComponents)
	mux.HandleFunc("/api/v1/components/hierarchy", h.API.HandleComponentsHierarchy)
	mux.HandleFunc("/api/v1/platforms", h.API.HandlePlatforms)
	mux.HandleFunc("/api/v1/instances", h.API.HandleInstances)
	mux.HandleFunc("/api/v1/components/list", h.API.HandleComponentsList)
}

// setupV1EventRoutes configures event management endpoints for v1
func setupV1EventRoutes(mux *http.ServeMux, h *handlers.Handlers) {
	// Incident endpoints
	mux.HandleFunc("/api/v1/incidents", h.Incident.HandleIncidents)
	mux.HandleFunc("/api/v1/incidents/", h.Incident.HandleIncidentByID)

	// Planned maintenance endpoints
	mux.HandleFunc("/api/v1/planned-maintenances", h.PlannedMaintenance.HandlePlannedMaintenances)
	mux.HandleFunc("/api/v1/planned-maintenances/", h.PlannedMaintenance.HandlePlannedMaintenanceByID)
}

// setupV1WeatherRoutes configures weather/status endpoints for v1
func setupV1WeatherRoutes(mux *http.ServeMux, h *handlers.Handlers) {
	// Weather endpoint (service status overview)
	mux.HandleFunc("GET /api/v1/weather", h.Weather.HandleWeather)
}

package drivers

import (
	"sort"
	"time"

	"github.com/gmllt/clariti/models"
	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/config"
)

// WeatherService calculates service weather based on incidents and planned maintenances
type WeatherService struct {
	config  *config.Config
	storage *RAMStorage
}

// NewWeatherService creates a new weather service
func NewWeatherService(cfg *config.Config, storage *RAMStorage) *WeatherService {
	return &WeatherService{
		config:  cfg,
		storage: storage,
	}
}

// GetWeatherSummary returns the current weather for all components
func (ws *WeatherService) GetWeatherSummary() (*models.WeatherSummary, error) {
	platformWeather := ws.calculatePlatformWeather()
	instanceWeather := ws.calculateInstanceWeather()
	componentWeather := ws.calculateComponentWeather()

	// Calculate overall status (worst status across all)
	overall := ws.calculateOverallWeather(platformWeather, instanceWeather, componentWeather)

	return &models.WeatherSummary{
		Platforms:  platformWeather,
		Instances:  instanceWeather,
		Components: componentWeather,
		Overall:    overall,
	}, nil
}

// calculatePlatformWeather calculates weather for all platforms
func (ws *WeatherService) calculatePlatformWeather() []models.ServiceWeather {
	var weather []models.ServiceWeather

	for _, platform := range ws.config.Components.Platforms {
		platformWeather := ws.calculateWeatherForPath(platform.Code, "", "")
		platformWeather.Platform = platform.Name
		platformWeather.PlatformCode = platform.Code
		weather = append(weather, platformWeather)
	}

	return weather
}

// calculateInstanceWeather calculates weather for all instances
func (ws *WeatherService) calculateInstanceWeather() []models.ServiceWeather {
	var weather []models.ServiceWeather

	for _, platform := range ws.config.Components.Platforms {
		for _, instance := range platform.Instances {
			instanceWeather := ws.calculateWeatherForPath(platform.Code, instance.Code, "")
			instanceWeather.Platform = platform.Name
			instanceWeather.PlatformCode = platform.Code
			instanceWeather.Instance = instance.Name
			instanceWeather.InstanceCode = instance.Code
			weather = append(weather, instanceWeather)
		}
	}

	return weather
}

// calculateComponentWeather calculates weather for all components
func (ws *WeatherService) calculateComponentWeather() []models.ServiceWeather {
	var weather []models.ServiceWeather

	for _, platform := range ws.config.Components.Platforms {
		for _, instance := range platform.Instances {
			for _, comp := range instance.Components {
				componentWeather := ws.calculateWeatherForPath(platform.Code, instance.Code, comp.Code)
				componentWeather.Platform = platform.Name
				componentWeather.PlatformCode = platform.Code
				componentWeather.Instance = instance.Name
				componentWeather.InstanceCode = instance.Code
				componentWeather.Component = comp.Name
				componentWeather.ComponentCode = comp.Code
				weather = append(weather, componentWeather)
			}
		}
	}

	return weather
}

// calculateWeatherForPath calculates weather for a specific component path
func (ws *WeatherService) calculateWeatherForPath(platformCode, instanceCode, componentCode string) models.ServiceWeather {
	activeEvents := ws.getActiveEventsForPath(platformCode, instanceCode, componentCode)

	// Find the highest criticality among active events
	maxCriticality := event.CriticalityOperational
	for _, evt := range activeEvents {
		if evt.Criticality > maxCriticality {
			maxCriticality = evt.Criticality
		}
	}

	return models.ServiceWeather{
		Status:       maxCriticality,
		StatusLabel:  maxCriticality.String(),
		ActiveEvents: activeEvents,
		LastUpdated:  time.Now().Format(time.RFC3339),
	}
}

// getActiveEventsForPath gets all active events affecting a specific component path
func (ws *WeatherService) getActiveEventsForPath(platformCode, instanceCode, componentCode string) []models.ActiveEvent {
	var activeEvents []models.ActiveEvent

	// Check incidents
	incidents, _ := ws.storage.GetAllIncidents()
	for _, incident := range incidents {
		if ws.isEventActive(incident) && ws.eventAffectsPath(incident.Components, platformCode, instanceCode, componentCode) {
			activeEvents = append(activeEvents, models.ActiveEvent{
				GUID:        incident.GUID,
				Type:        incident.Type(),
				Title:       incident.Title,
				Status:      incident.Status(),
				Criticality: incident.Criticality(),
			})
		}
	}

	// Check planned maintenances
	maintenances, _ := ws.storage.GetAllPlannedMaintenances()
	for _, maintenance := range maintenances {
		if ws.isEventActive(maintenance) && ws.eventAffectsPath(maintenance.Components, platformCode, instanceCode, componentCode) {
			activeEvents = append(activeEvents, models.ActiveEvent{
				GUID:        maintenance.GUID,
				Type:        maintenance.Type(),
				Title:       maintenance.Title,
				Status:      maintenance.Status(),
				Criticality: maintenance.Criticality(),
			})
		}
	}

	// Sort by criticality (highest first)
	sort.Slice(activeEvents, func(i, j int) bool {
		return activeEvents[i].Criticality > activeEvents[j].Criticality
	})

	return activeEvents
}

// isEventActive checks if an event is currently active
func (ws *WeatherService) isEventActive(evt event.Event) bool {
	// For the interface, we need to access the timing through a type assertion
	switch e := evt.(type) {
	case *event.Incident:
		status := e.Status()
		// An incident is active if it's ongoing or acknowledged
		// Also consider incidents without start_effective as active (immediate incidents)
		if status == event.StatusOnGoing || status == event.StatusAcknowledged {
			return true
		}
		// If no start_effective is set and no end_effective, consider it active
		if e.StartEffective == nil && e.EndEffective == nil {
			return true
		}
		return false
	case *event.PlannedMaintenance:
		status := e.Status()
		// A maintenance is active if it's ongoing or planned
		if status == event.StatusOnGoing || status == event.StatusPlanned {
			return true
		}
		// If no start_effective is set and no end_effective, consider it active
		if e.StartEffective == nil && e.EndEffective == nil {
			return true
		}
		return false
	}
	return false
}

// eventAffectsPath checks if an event affects a specific component path
func (ws *WeatherService) eventAffectsPath(components []*component.Component, platformCode, instanceCode, componentCode string) bool {
	for _, comp := range components {
		if ws.componentMatchesPath(comp, platformCode, instanceCode, componentCode) {
			return true
		}
	}
	return false
}

// componentMatchesPath checks if a component matches the given path
func (ws *WeatherService) componentMatchesPath(comp *component.Component, platformCode, instanceCode, componentCode string) bool {
	// Get platform code through the hierarchy
	var compPlatformCode, compInstanceCode, compComponentCode string

	if comp.Instance != nil {
		compInstanceCode = comp.Instance.Code
		if comp.Instance.Platform != nil {
			compPlatformCode = comp.Instance.Platform.Code
		}
	}
	compComponentCode = comp.Code

	// If component code is specified, it must match exactly
	if componentCode != "" {
		return compPlatformCode == platformCode && compInstanceCode == instanceCode && compComponentCode == componentCode
	}

	// If only instance is specified, component must be in that instance
	if instanceCode != "" {
		return compPlatformCode == platformCode && compInstanceCode == instanceCode
	}

	// If only platform is specified, component must be in that platform
	return compPlatformCode == platformCode
}

// calculateOverallWeather calculates the overall system weather
func (ws *WeatherService) calculateOverallWeather(platforms, instances, components []models.ServiceWeather) models.ServiceWeather {
	maxCriticality := event.CriticalityOperational
	var allEvents []models.ActiveEvent

	// Collect all unique events and find max criticality
	eventMap := make(map[string]models.ActiveEvent)

	for _, weather := range append(append(platforms, instances...), components...) {
		if weather.Status > maxCriticality {
			maxCriticality = weather.Status
		}
		for _, evt := range weather.ActiveEvents {
			eventMap[evt.GUID] = evt
		}
	}

	// Convert map to slice
	for _, evt := range eventMap {
		allEvents = append(allEvents, evt)
	}

	// Sort by criticality
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].Criticality > allEvents[j].Criticality
	})

	return models.ServiceWeather{
		Platform:     "Overall System",
		PlatformCode: "ALL",
		Status:       maxCriticality,
		StatusLabel:  maxCriticality.String(),
		ActiveEvents: allEvents,
		LastUpdated:  time.Now().Format(time.RFC3339),
	}
}

package models

import (
	"testing"

	"github.com/gmllt/clariti/models/event"
)

func TestServiceWeather_Creation(t *testing.T) {
	weather := ServiceWeather{
		Platform:      "TestPlatform",
		PlatformCode:  "TST",
		Instance:      "TestInstance",
		InstanceCode:  "test",
		Component:     "TestComponent",
		ComponentCode: "comp",
		Status:        event.CriticalityOperational,
		StatusLabel:   "operational",
		LastUpdated:   "2023-01-01T00:00:00Z",
	}

	if weather.Platform != "TestPlatform" {
		t.Errorf("Expected platform 'TestPlatform', got '%s'", weather.Platform)
	}
	if weather.Status != event.CriticalityOperational {
		t.Errorf("Expected operational status, got %d", weather.Status)
	}
}

func TestActiveEvent_Creation(t *testing.T) {
	activeEvent := ActiveEvent{
		GUID:        "test-guid",
		Type:        event.TypeFiringIncident,
		Title:       "Test Incident",
		Status:      event.StatusOnGoing,
		Criticality: event.CriticalityMajorOutage,
	}

	if activeEvent.GUID != "test-guid" {
		t.Errorf("Expected GUID 'test-guid', got '%s'", activeEvent.GUID)
	}
	if activeEvent.Type != event.TypeFiringIncident {
		t.Errorf("Expected firing incident type, got %s", activeEvent.Type)
	}
	if activeEvent.Criticality != event.CriticalityMajorOutage {
		t.Errorf("Expected major outage criticality, got %d", activeEvent.Criticality)
	}
}

func TestWeatherSummary_Creation(t *testing.T) {
	platform := ServiceWeather{
		Platform:     "TestPlatform",
		PlatformCode: "TST",
		Status:       event.CriticalityOperational,
		StatusLabel:  "operational",
	}

	instance := ServiceWeather{
		Platform:     "TestPlatform",
		PlatformCode: "TST",
		Instance:     "TestInstance",
		InstanceCode: "test",
		Status:       event.CriticalityDegraded,
		StatusLabel:  "degraded",
	}

	component := ServiceWeather{
		Platform:      "TestPlatform",
		PlatformCode:  "TST",
		Instance:      "TestInstance",
		InstanceCode:  "test",
		Component:     "TestComponent",
		ComponentCode: "comp",
		Status:        event.CriticalityMajorOutage,
		StatusLabel:   "major outage",
	}

	overall := ServiceWeather{
		Platform:     "Overall System",
		PlatformCode: "ALL",
		Status:       event.CriticalityMajorOutage,
		StatusLabel:  "major outage",
	}

	summary := WeatherSummary{
		Platforms:  []ServiceWeather{platform},
		Instances:  []ServiceWeather{instance},
		Components: []ServiceWeather{component},
		Overall:    overall,
	}

	if len(summary.Platforms) != 1 {
		t.Errorf("Expected 1 platform, got %d", len(summary.Platforms))
	}
	if len(summary.Instances) != 1 {
		t.Errorf("Expected 1 instance, got %d", len(summary.Instances))
	}
	if len(summary.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(summary.Components))
	}

	if summary.Overall.Status != event.CriticalityMajorOutage {
		t.Errorf("Expected major outage status, got %d", summary.Overall.Status)
	}
}

func TestComponentWeather_Creation(t *testing.T) {
	compWeather := ComponentWeather{
		Level:       "component",
		Name:        "TestComponent",
		Code:        "comp",
		Status:      event.CriticalityPartialOutage,
		StatusLabel: "partial outage",
		Path:        "TST/test/comp",
	}

	if compWeather.Level != "component" {
		t.Errorf("Expected level 'component', got '%s'", compWeather.Level)
	}
	if compWeather.Path != "TST/test/comp" {
		t.Errorf("Expected path 'TST/test/comp', got '%s'", compWeather.Path)
	}
	if compWeather.Status != event.CriticalityPartialOutage {
		t.Errorf("Expected partial outage status, got %d", compWeather.Status)
	}
}

func TestServiceWeather_WithActiveEvents(t *testing.T) {
	activeEvents := []ActiveEvent{
		{
			GUID:        "incident-1",
			Type:        event.TypeFiringIncident,
			Title:       "Database Issue",
			Status:      event.StatusOnGoing,
			Criticality: event.CriticalityDegraded,
		},
		{
			GUID:        "maintenance-1",
			Type:        event.TypePlannedMaintenance,
			Title:       "Server Upgrade",
			Status:      event.StatusPlanned,
			Criticality: event.CriticalityUnderMaintenance,
		},
	}

	weather := ServiceWeather{
		Platform:     "TestPlatform",
		PlatformCode: "TST",
		Status:       event.CriticalityUnderMaintenance, // Highest criticality
		StatusLabel:  "under maintenance",
		ActiveEvents: activeEvents,
		LastUpdated:  "2023-01-01T00:00:00Z",
	}

	if len(weather.ActiveEvents) != 2 {
		t.Errorf("Expected 2 active events, got %d", len(weather.ActiveEvents))
	}

	// Status should reflect highest criticality
	if weather.Status != event.CriticalityUnderMaintenance {
		t.Errorf("Expected under maintenance status, got %d", weather.Status)
	}

	// Check individual events
	incident := weather.ActiveEvents[0]
	if incident.Type != event.TypeFiringIncident {
		t.Errorf("Expected firing incident, got %s", incident.Type)
	}

	maintenance := weather.ActiveEvents[1]
	if maintenance.Type != event.TypePlannedMaintenance {
		t.Errorf("Expected planned maintenance, got %s", maintenance.Type)
	}
}

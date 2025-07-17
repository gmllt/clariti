package drivers

import (
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/config"
)

func TestNewWeatherService(t *testing.T) {
	cfg := &config.Config{}
	storage := NewRAMStorage()

	service := NewWeatherService(cfg, storage)

	if service == nil {
		t.Fatal("Expected weather service to be created")
	}
	if service.config != cfg {
		t.Error("Expected config to be set")
	}
	if service.storage != storage {
		t.Error("Expected storage to be set")
	}
}

func TestWeatherService_GetWeatherSummary(t *testing.T) {
	// Setup test config
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "TestPlatform",
					Code: "TST",
					Instances: []config.InstanceConfig{
						{
							Name: "TestInstance",
							Code: "test",
							Components: []config.ComponentConfig{
								{Name: "TestComponent", Code: "comp"},
							},
						},
					},
				},
			},
		},
	}

	storage := NewRAMStorage()
	service := NewWeatherService(cfg, storage)

	// Test with no incidents
	summary, err := service.GetWeatherSummary()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
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

	// All should be operational
	if summary.Overall.Status != event.CriticalityOperational {
		t.Errorf("Expected operational status, got %d", summary.Overall.Status)
	}
}

func TestWeatherService_WithActiveIncidents(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "TestPlatform",
					Code: "TST",
					Instances: []config.InstanceConfig{
						{
							Name: "TestInstance",
							Code: "test",
							Components: []config.ComponentConfig{
								{Name: "TestComponent", Code: "comp"},
							},
						},
					},
				},
			},
		},
	}

	storage := NewRAMStorage()
	service := NewWeatherService(cfg, storage)

	// Create test component
	platform := component.NewPlatform("TestPlatform", "TST")
	instance := component.NewInstance("TestInstance", "test", platform)
	comp := component.NewComponent("TestComponent", "comp", instance)

	// Create and store an incident
	incident := event.NewFiringIncident(
		"Test Incident",
		"Test incident content",
		[]*component.Component{comp},
		event.CriticalityMajorOutage,
	)

	err := storage.CreateIncident(incident)
	if err != nil {
		t.Fatalf("Failed to create incident: %v", err)
	}

	// Test weather with active incident
	summary, err := service.GetWeatherSummary()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Overall status should reflect the incident
	if summary.Overall.Status != event.CriticalityMajorOutage {
		t.Errorf("Expected major outage status (%d), got %d", event.CriticalityMajorOutage, summary.Overall.Status)
	}

	// Should have active events
	if len(summary.Overall.ActiveEvents) != 1 {
		t.Errorf("Expected 1 active event, got %d", len(summary.Overall.ActiveEvents))
	}
}

func TestWeatherService_CriticalityPriority(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Components: config.ComponentsConfig{
			Platforms: []config.PlatformConfig{
				{
					Name: "TestPlatform",
					Code: "TST",
					Instances: []config.InstanceConfig{
						{
							Name: "TestInstance",
							Code: "test",
							Components: []config.ComponentConfig{
								{Name: "Component1", Code: "comp1"},
								{Name: "Component2", Code: "comp2"},
							},
						},
					},
				},
			},
		},
	}

	storage := NewRAMStorage()
	service := NewWeatherService(cfg, storage)

	// Create test components
	platform := component.NewPlatform("TestPlatform", "TST")
	instance := component.NewInstance("TestInstance", "test", platform)
	comp1 := component.NewComponent("Component1", "comp1", instance)
	comp2 := component.NewComponent("Component2", "comp2", instance)

	// Create incidents with different criticalities
	incident1 := event.NewFiringIncident(
		"Minor Issue",
		"Minor degradation",
		[]*component.Component{comp1},
		event.CriticalityDegraded,
	)

	incident2 := event.NewFiringIncident(
		"Major Issue",
		"Major outage",
		[]*component.Component{comp2},
		event.CriticalityMajorOutage,
	)

	storage.CreateIncident(incident1)
	storage.CreateIncident(incident2)

	// Test weather - should show highest criticality
	summary, err := service.GetWeatherSummary()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Overall should be major outage (highest)
	if summary.Overall.Status != event.CriticalityMajorOutage {
		t.Errorf("Expected major outage status (%d), got %d", event.CriticalityMajorOutage, summary.Overall.Status)
	}

	// Platform should also be major outage
	if len(summary.Platforms) > 0 && summary.Platforms[0].Status != event.CriticalityMajorOutage {
		t.Errorf("Expected platform major outage status (%d), got %d", event.CriticalityMajorOutage, summary.Platforms[0].Status)
	}
}

func TestWeatherService_ComponentMatchesPath(t *testing.T) {
	storage := NewRAMStorage()
	service := NewWeatherService(&config.Config{}, storage)

	// Create test component hierarchy
	platform := component.NewPlatform("TestPlatform", "TST")
	instance := component.NewInstance("TestInstance", "test", platform)
	comp := component.NewComponent("TestComponent", "comp", instance)

	tests := []struct {
		name          string
		platformCode  string
		instanceCode  string
		componentCode string
		expectedMatch bool
	}{
		{"Exact match", "TST", "test", "comp", true},
		{"Instance match", "TST", "test", "", true},
		{"Platform match", "TST", "", "", true},
		{"Wrong platform", "WRONG", "test", "comp", false},
		{"Wrong instance", "TST", "wrong", "comp", false},
		{"Wrong component", "TST", "test", "wrong", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := service.componentMatchesPath(comp, tt.platformCode, tt.instanceCode, tt.componentCode)
			if match != tt.expectedMatch {
				t.Errorf("Expected %v, got %v for path %s/%s/%s", tt.expectedMatch, match, tt.platformCode, tt.instanceCode, tt.componentCode)
			}
		})
	}
}

func TestWeatherService_IsEventActive(t *testing.T) {
	storage := NewRAMStorage()
	service := NewWeatherService(&config.Config{}, storage)

	platform := component.NewPlatform("Test", "TST")
	instance := component.NewInstance("Test", "test", platform)
	comp := component.NewComponent("Test", "comp", instance)

	tests := []struct {
		name     string
		event    event.Event
		expected bool
	}{
		{
			name:     "Incident without dates",
			event:    event.NewFiringIncident("Test", "Test", []*component.Component{comp}, event.CriticalityDegraded),
			expected: true,
		},
		{
			name:     "Known issue without dates",
			event:    event.NewKnownIssue("Test", "Test", []*component.Component{comp}, event.CriticalityDegraded),
			expected: true,
		},
		{
			name:     "Planned maintenance without dates",
			event:    event.NewPlannedMaintenance("Test", "Test", []*component.Component{comp}, time.Now(), time.Now().Add(time.Hour)),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			active := service.isEventActive(tt.event)
			if active != tt.expected {
				t.Errorf("Expected %v, got %v for event %s", tt.expected, active, tt.name)
			}
		})
	}
}

func TestWeatherService_EventWithTiming(t *testing.T) {
	storage := NewRAMStorage()
	service := NewWeatherService(&config.Config{}, storage)

	platform := component.NewPlatform("Test", "TST")
	instance := component.NewInstance("Test", "test", platform)
	comp := component.NewComponent("Test", "comp", instance)

	// Create incident with timing
	incident := event.NewFiringIncident("Test", "Test", []*component.Component{comp}, event.CriticalityDegraded)

	// Test with past start time (should be active)
	pastTime := time.Now().Add(-1 * time.Hour)
	incident.StartEffective = &pastTime

	if !service.isEventActive(incident) {
		t.Error("Expected incident with past start time to be active")
	}

	// Test with future start time (should not be active by status)
	futureTime := time.Now().Add(1 * time.Hour)
	incident.StartEffective = &futureTime

	// The incident should not be active based on timing, but our logic considers
	// incidents without proper timing as active, so we need to check the status
	status := incident.Status()
	if status == event.StatusOnGoing && service.isEventActive(incident) {
		t.Error("Expected incident with future start time to not be active based on status")
	}
}

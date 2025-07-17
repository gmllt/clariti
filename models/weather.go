package models

import (
	"github.com/gmllt/clariti/models/event"
)

// ServiceWeather represents the current "weather" status of a service component
type ServiceWeather struct {
	Platform      string            `json:"platform"`
	PlatformCode  string            `json:"platform_code"`
	Instance      string            `json:"instance,omitempty"`
	InstanceCode  string            `json:"instance_code,omitempty"`
	Component     string            `json:"component,omitempty"`
	ComponentCode string            `json:"component_code,omitempty"`
	Status        event.Criticality `json:"status"`
	StatusLabel   string            `json:"status_label"`
	ActiveEvents  []ActiveEvent     `json:"active_events,omitempty"`
	LastUpdated   string            `json:"last_updated"`
}

// ActiveEvent represents an active event affecting the service
type ActiveEvent struct {
	GUID        string            `json:"guid"`
	Type        event.TypeEvent   `json:"type"`
	Title       string            `json:"title"`
	Status      event.Status      `json:"status"`
	Criticality event.Criticality `json:"criticality"`
}

// WeatherSummary provides aggregated weather information
type WeatherSummary struct {
	Platforms  []ServiceWeather `json:"platforms"`
	Instances  []ServiceWeather `json:"instances"`
	Components []ServiceWeather `json:"components"`
	Overall    ServiceWeather   `json:"overall"`
}

// ComponentWeather represents weather for a specific component hierarchy level
type ComponentWeather struct {
	Level        string            `json:"level"` // "platform", "instance", "component"
	Name         string            `json:"name"`
	Code         string            `json:"code"`
	Status       event.Criticality `json:"status"`
	StatusLabel  string            `json:"status_label"`
	Path         string            `json:"path"` // e.g., "PXD/z1/api"
	ActiveEvents []ActiveEvent     `json:"active_events,omitempty"`
}

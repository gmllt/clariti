package drivers

import (
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/server/config"
)

// BenchmarkWeatherService_GetWeatherSummary tests the performance of weather summary calculation
func BenchmarkWeatherService_GetWeatherSummary(b *testing.B) {
	storage := NewRAMStorage()
	cfg := &config.Config{}
	service := NewWeatherService(cfg, storage)

	// Setup test data
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := component.NewComponent("API Gateway", "api-gateway", instance)

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	incident := &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:           "incident-1",
			Title:          "API Issues",
			Content:        "API experiencing issues",
			Components:     []*component.Component{comp},
			StartEffective: &past,
			EndEffective:   &future,
		},
		IncidentCriticality: event.CriticalityMajorOutage,
	}

	storage.CreateIncident(incident)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetWeatherSummary()
	}
}

// BenchmarkWeatherService_CalculateWeatherForPath tests path-specific weather calculation
func BenchmarkWeatherService_CalculateWeatherForPath(b *testing.B) {
	storage := NewRAMStorage()
	cfg := &config.Config{}
	service := NewWeatherService(cfg, storage)

	// Setup test data
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := component.NewComponent("API Gateway", "api-gateway", instance)

	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	incident := &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:           "incident-1",
			Title:          "API Issues",
			Content:        "API experiencing issues",
			Components:     []*component.Component{comp},
			StartEffective: &past,
			EndEffective:   &future,
		},
		IncidentCriticality: event.CriticalityMajorOutage,
	}

	storage.CreateIncident(incident)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.calculateWeatherForPath("aws-prod", "eks-cluster", "api-gateway")
	}
}

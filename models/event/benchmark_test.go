package event

import (
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
)

// Benchmark tests for events with components
func BenchmarkIncident_Create(b *testing.B) {
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp1 := component.NewComponent("API Gateway", "api-gateway", instance)
	comp2 := component.NewComponent("User Service", "user-service", instance)
	components := []*component.Component{comp1, comp2}

	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		incident := &Incident{
			BaseEvent: BaseEvent{
				GUID:           "test-incident",
				Title:          "Service Outage",
				Content:        "Critical service failure",
				Components:     components,
				StartEffective: &now,
				ExtraFields: map[string]string{
					"severity": "critical",
					"region":   "us-east-1",
				},
			},
			Perpetual:           false,
			IncidentCriticality: CriticalityMajorOutage,
		}
		_ = incident
	}
}

func BenchmarkIncident_Status(b *testing.B) {
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := component.NewComponent("API Gateway", "api-gateway", instance)

	now := time.Now()
	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:           "test-incident",
			Title:          "Service Outage",
			Components:     []*component.Component{comp},
			StartEffective: &now,
		},
		IncidentCriticality: CriticalityMajorOutage,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = incident.Status()
	}
}

func BenchmarkIncident_Type(b *testing.B) {
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := component.NewComponent("API Gateway", "api-gateway", instance)

	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:       "test-incident",
			Title:      "Service Outage",
			Components: []*component.Component{comp},
		},
		IncidentCriticality: CriticalityMajorOutage,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = incident.Type()
	}
}

func BenchmarkPlannedMaintenance_Create(b *testing.B) {
	platform := component.NewPlatform("AWS Production", "aws-prod")
	instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)
	comp1 := component.NewComponent("API Gateway", "api-gateway", instance)
	comp2 := component.NewComponent("User Service", "user-service", instance)
	components := []*component.Component{comp1, comp2}

	now := time.Now()
	startPlanned := now.Add(1 * time.Hour)
	endPlanned := now.Add(2 * time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maintenance := &PlannedMaintenance{
			BaseEvent: BaseEvent{
				GUID:       "test-maintenance",
				Title:      "System Upgrade",
				Content:    "Scheduled maintenance window",
				Components: components,
				ExtraFields: map[string]string{
					"version": "v2.0",
					"region":  "us-east-1",
				},
			},
			StartPlanned: startPlanned,
			EndPlanned:   endPlanned,
			Cancelled:    false,
		}
		_ = maintenance
	}
}

func BenchmarkComponent_WithManyAllocations(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate creating components in a loop with many allocations
		platform := component.NewPlatform("AWS Production", "aws-prod")
		instance := component.NewInstance("EKS Cluster", "eks-cluster", platform)

		// Create multiple components and normalize them
		for j := 0; j < 5; j++ {
			comp := component.NewComponent("Service", "service", instance)
			_ = comp.Normalize()
		}
	}
}

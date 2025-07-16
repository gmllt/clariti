package models_test

import (
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
	"github.com/gmllt/clariti/models/event"
	"github.com/gmllt/clariti/utils"
)

// TestFullIntegration tests the complete integration between components and events
func TestFullIntegration(t *testing.T) {
	// Create a realistic infrastructure hierarchy
	awsPlatform := component.NewPlatform("AWS Production US-East-1", "aws-prod-us-east-1")
	eksInstance := component.NewInstance("EKS Cluster v1.25", "eks-cluster", awsPlatform)

	// Create multiple components representing a microservices architecture
	apiGateway := component.NewComponent("API Gateway", "api-gateway", eksInstance)
	userService := component.NewComponent("User Service", "user-svc", eksInstance)
	orderService := component.NewComponent("Order Service", "order-svc", eksInstance)
	paymentService := component.NewComponent("Payment Service", "payment-svc", eksInstance)

	// Test component stringable functionality
	expectedAPI := "AWS Production US-East-1 - EKS Cluster v1.25 - API Gateway"
	if apiGateway.String() != expectedAPI {
		t.Errorf("API Gateway String() = %v, want %v", apiGateway.String(), expectedAPI)
	}

	// Test component normalization
	expectedAPINorm := "aws-prod-us-east-1-eks-cluster-api-gateway"
	if apiGateway.Normalize() != expectedAPINorm {
		t.Errorf("API Gateway Normalize() = %v, want %v", apiGateway.Normalize(), expectedAPINorm)
	}

	now := time.Now()

	// Create an incident affecting multiple components
	incident := &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:    "incident-payment-outage-001",
			Title:   "Payment Service Outage",
			Content: "Critical payment processing failure affecting user transactions",
			Components: []*component.Component{
				apiGateway, userService, orderService, paymentService,
			},
			StartEffective: &now,
			EndEffective:   nil, // Ongoing
			ExtraFields: map[string]string{
				"severity":        "critical",
				"business_impact": "revenue_loss",
				"region":          "us-east-1",
				"alert_source":    "prometheus",
			},
		},
		Perpetual:           false,
		IncidentCriticality: event.CriticalityMajorOutage,
	}

	// Test incident behavior
	if incident.Type() != event.TypeFiringIncident {
		t.Errorf("Expected TypeFiringIncident, got %v", incident.Type())
	}

	if incident.Status() != event.StatusOnGoing {
		t.Errorf("Expected StatusOnGoing for active incident, got %v", incident.Status())
	}

	if incident.Criticality() != event.CriticalityMajorOutage {
		t.Errorf("Expected CriticalityMajorOutage, got %v", incident.Criticality())
	}

	// Test component integration in incident
	if len(incident.Components) != 4 {
		t.Errorf("Expected 4 affected components, got %v", len(incident.Components))
	}

	// Verify each component maintains its hierarchy
	for i, comp := range incident.Components {
		if comp.String() == "" {
			t.Errorf("Component %d has empty string representation", i)
		}
		if comp.Normalize() == "" {
			t.Errorf("Component %d has empty normalized representation", i)
		}

		// All components should be from the same platform
		if !containsString(comp.String(), "AWS Production US-East-1") {
			t.Errorf("Component %d should contain platform name: %v", i, comp.String())
		}
	}

	// Create a planned maintenance event for the same components
	futureStart := now.Add(24 * time.Hour)
	futureEnd := now.Add(26 * time.Hour)

	maintenance := &event.PlannedMaintenance{
		BaseEvent: event.BaseEvent{
			GUID:    "maintenance-eks-upgrade-001",
			Title:   "EKS Cluster Upgrade to v1.26",
			Content: "Scheduled upgrade of EKS cluster with rolling deployment",
			Components: []*component.Component{
				apiGateway, userService, orderService, paymentService,
			},
			ExtraFields: map[string]string{
				"upgrade_version":   "1.26.0",
				"expected_downtime": "minimal",
				"rollback_plan":     "available",
				"business_approval": "approved",
			},
		},
		StartPlanned: futureStart,
		EndPlanned:   futureEnd,
		Cancelled:    false,
	}

	// Test maintenance event behavior
	if maintenance.Type() != event.TypePlannedMaintenance {
		t.Errorf("Expected TypePlannedMaintenance, got %v", maintenance.Type())
	}

	if maintenance.Status() != event.StatusPlanned {
		t.Errorf("Expected StatusPlanned for future maintenance, got %v", maintenance.Status())
	}

	if maintenance.Criticality() != event.CriticalityUnderMaintenance {
		t.Errorf("Expected CriticalityUnderMaintenance, got %v", maintenance.Criticality())
	}

	// Test that both events share the same components
	if len(maintenance.Components) != len(incident.Components) {
		t.Errorf("Maintenance and incident should affect same number of components")
	}

	// Test component normalization consistency between events
	for i := range incident.Components {
		incidentNorm := incident.Components[i].Normalize()
		maintenanceNorm := maintenance.Components[i].Normalize()

		if incidentNorm != maintenanceNorm {
			t.Errorf("Component %d normalization mismatch between events: %v vs %v",
				i, incidentNorm, maintenanceNorm)
		}
	}
}

// TestComponentInterface tests that components properly implement required interfaces
func TestComponentInterface(t *testing.T) {
	platform := component.NewPlatform("Test Platform", "test-platform")
	instance := component.NewInstance("Test Instance", "test-instance", platform)
	comp := component.NewComponent("Test Component", "test-component", instance)

	// Test Stringable interface
	var stringables []utils.Stringable = []utils.Stringable{
		platform, instance, comp,
	}

	for i, s := range stringables {
		if s.String() == "" {
			t.Errorf("Stringable %d returned empty string", i)
		}
	}

	// Test Normalizable interface
	var normalizables []utils.Normalizable = []utils.Normalizable{
		platform, instance, comp,
	}

	for i, n := range normalizables {
		if n.Normalize() == "" {
			t.Errorf("Normalizable %d returned empty normalized string", i)
		}

		// Normalization should always be lowercase and hyphenated
		norm := n.Normalize()
		if norm != utils.NormalizeFromStringable(n) {
			t.Errorf("Normalizable %d inconsistent with NormalizeFromStringable", i)
		}
	}
}

// TestEventInterface tests that events properly implement the Event interface
func TestEventInterface(t *testing.T) {
	now := time.Now()

	// Create test components
	platform := component.NewPlatform("Interface Test Platform", "interface-test-platform")
	instance := component.NewInstance("Interface Test Instance", "interface-test-instance", platform)
	comp := component.NewComponent("Interface Test Component", "interface-test-component", instance)

	// Create different event types
	incident := &event.Incident{
		BaseEvent: event.BaseEvent{
			GUID:       "interface-test-incident",
			Title:      "Interface Test Incident",
			Components: []*component.Component{comp},
		},
		Perpetual:           false,
		IncidentCriticality: event.CriticalityMajorOutage,
	}

	maintenance := &event.PlannedMaintenance{
		BaseEvent: event.BaseEvent{
			GUID:       "interface-test-maintenance",
			Title:      "Interface Test Maintenance",
			Components: []*component.Component{comp},
		},
		StartPlanned: now.Add(1 * time.Hour),
		EndPlanned:   now.Add(2 * time.Hour),
	}

	// Test Event interface compliance
	var events []event.Event = []event.Event{incident, maintenance}

	for i, e := range events {
		// All events should have a type
		if e.Type() == "" {
			t.Errorf("Event %d has empty type", i)
		}

		// All events should have a status
		status := e.Status()
		validStatuses := []event.Status{
			event.StatusPlanned, event.StatusOnGoing, event.StatusResolved,
			event.StatusAcknowledged, event.StatusCanceled, event.StatusUnknown,
		}

		isValidStatus := false
		for _, validStatus := range validStatuses {
			if status == validStatus {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			t.Errorf("Event %d has invalid status: %v", i, status)
		}

		// All events should have a criticality
		criticality := e.Criticality()
		if criticality < event.CriticalityUnknown || criticality > event.CriticalityUnderMaintenance {
			t.Errorf("Event %d has invalid criticality: %v", i, criticality)
		}
	}
}

// TestCriticalityStringRepresentation tests the string representation of criticality levels
func TestCriticalityStringRepresentation(t *testing.T) {
	tests := []struct {
		criticality event.Criticality
		expected    string
	}{
		{event.CriticalityOperational, "operational"},
		{event.CriticalityDegraded, "degraded"},
		{event.CriticalityPartialOutage, "partial outage"},
		{event.CriticalityMajorOutage, "major outage"},
		{event.CriticalityUnderMaintenance, "under maintenance"},
		{event.CriticalityUnknown, "unknown"},
	}

	for _, tt := range tests {
		if got := tt.criticality.String(); got != tt.expected {
			t.Errorf("Criticality(%d).String() = %v, want %v",
				int(tt.criticality), got, tt.expected)
		}
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(findSubstring(s, substr) != -1)))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

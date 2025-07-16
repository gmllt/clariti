package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
)

func TestIncident_Type(t *testing.T) {
	tests := []struct {
		name      string
		perpetual bool
		expected  TypeEvent
	}{
		{
			name:      "Firing Incident",
			perpetual: false,
			expected:  TypeFiringIncident,
		},
		{
			name:      "Known Issue",
			perpetual: true,
			expected:  TypeKnownIssue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incident := &Incident{
				Perpetual: tt.perpetual,
			}
			if got := incident.Type(); got != tt.expected {
				t.Errorf("Incident.Type() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIncident_Status(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name           string
		startEffective *time.Time
		endEffective   *time.Time
		expected       Status
	}{
		{
			name:           "Unknown Status - No Times",
			startEffective: nil,
			endEffective:   nil,
			expected:       StatusUnknown,
		},
		{
			name:           "Resolved - End in Past",
			startEffective: &past,
			endEffective:   &past,
			expected:       StatusResolved,
		},
		{
			name:           "Ongoing - Started, Not Ended",
			startEffective: &past,
			endEffective:   &future,
			expected:       StatusOnGoing,
		},
		{
			name:           "Ongoing - Started, No End",
			startEffective: &past,
			endEffective:   nil,
			expected:       StatusOnGoing,
		},
		{
			name:           "Unknown - Not Started",
			startEffective: &future,
			endEffective:   nil,
			expected:       StatusUnknown,
		},
		{
			name:           "Resolved - End Only in Past",
			startEffective: nil,
			endEffective:   &past,
			expected:       StatusResolved,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incident := &Incident{
				BaseEvent: BaseEvent{
					StartEffective: tt.startEffective,
					EndEffective:   tt.endEffective,
				},
			}
			if got := incident.Status(); got != tt.expected {
				t.Errorf("Incident.Status() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIncident_Criticality(t *testing.T) {
	tests := []struct {
		name        string
		criticality Criticality
		expected    Criticality
	}{
		{
			name:        "Operational",
			criticality: CriticalityOperational,
			expected:    CriticalityOperational,
		},
		{
			name:        "Major Outage",
			criticality: CriticalityMajorOutage,
			expected:    CriticalityMajorOutage,
		},
		{
			name:        "Unknown",
			criticality: CriticalityUnknown,
			expected:    CriticalityUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incident := &Incident{
				IncidentCriticality: tt.criticality,
			}
			if got := incident.Criticality(); got != tt.expected {
				t.Errorf("Incident.Criticality() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIncident_InterfaceCompliance(t *testing.T) {
	incident := &Incident{
		BaseEvent: BaseEvent{
			Title: "Test Incident",
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityMajorOutage,
	}

	// Test Event interface compliance
	var event Event = incident
	if event.Type() != TypeFiringIncident {
		t.Errorf("Expected TypeFiringIncident, got %v", event.Type())
	}
	if event.Criticality() != CriticalityMajorOutage {
		t.Errorf("Expected CriticalityMajorOutage, got %v", event.Criticality())
	}
}

func TestIncident_FullLifecycle(t *testing.T) {
	now := time.Now()
	past := now.Add(-2 * time.Hour)
	end := now.Add(-1 * time.Hour)

	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:           "incident-123",
			Title:          "Production Database Outage",
			Content:        "Database is experiencing connectivity issues",
			StartEffective: &past,
			EndEffective:   &end,
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityMajorOutage,
	}

	// Test all methods
	if incident.Type() != TypeFiringIncident {
		t.Errorf("Expected TypeFiringIncident, got %v", incident.Type())
	}

	if incident.Status() != StatusResolved {
		t.Errorf("Expected StatusResolved, got %v", incident.Status())
	}

	if incident.Criticality() != CriticalityMajorOutage {
		t.Errorf("Expected CriticalityMajorOutage, got %v", incident.Criticality())
	}
}

func TestIncident_ComponentIntegration(t *testing.T) {
	// Create realistic component hierarchy
	awsPlatform := component.NewPlatform("AWS US-East-1", "aws-us-east-1")
	rdsInstance := component.NewInstance("RDS Production Cluster", "rds-prod", awsPlatform)
	primaryDB := component.NewComponent("Primary Database", "primary-db", rdsInstance)
	replicaDB := component.NewComponent("Read Replica", "read-replica", rdsInstance)

	now := time.Now()
	past := now.Add(-30 * time.Minute)

	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:           "incident-rds-outage-001",
			Title:          "RDS Production Database Outage",
			Content:        "Multiple database instances experiencing connection timeouts and high latency",
			Components:     []*component.Component{primaryDB, replicaDB},
			StartEffective: &past,
			EndEffective:   nil, // Ongoing incident
			ExtraFields: map[string]string{
				"region":            "us-east-1",
				"availability_zone": "us-east-1a",
				"service":           "rds",
				"alert_source":      "cloudwatch",
			},
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityMajorOutage,
	}

	// Test component integration
	if len(incident.Components) != 2 {
		t.Errorf("Expected 2 components, got %v", len(incident.Components))
	}

	// Test component details
	expectedPrimary := "AWS US-East-1 - RDS Production Cluster - Primary Database"
	if incident.Components[0].String() != expectedPrimary {
		t.Errorf("Primary DB String() = %v, want %v", incident.Components[0].String(), expectedPrimary)
	}

	expectedReplica := "AWS US-East-1 - RDS Production Cluster - Read Replica"
	if incident.Components[1].String() != expectedReplica {
		t.Errorf("Replica DB String() = %v, want %v", incident.Components[1].String(), expectedReplica)
	}

	// Test normalization for component identification (now uses codes)
	expectedPrimaryNorm := "aws-us-east-1-rds-prod-primary-db"
	if incident.Components[0].Normalize() != expectedPrimaryNorm {
		t.Errorf("Primary DB Normalize() = %v, want %v", incident.Components[0].Normalize(), expectedPrimaryNorm)
	}

	// Test incident status with ongoing components
	if incident.Status() != StatusOnGoing {
		t.Errorf("Expected ongoing status for incident with components, got %v", incident.Status())
	}

	// Test incident type
	if incident.Type() != TypeFiringIncident {
		t.Errorf("Expected firing incident type, got %v", incident.Type())
	}
}

func TestIncident_MultiPlatformComponents(t *testing.T) {
	// Create components from multiple platforms
	awsPlatform := component.NewPlatform("AWS", "aws")
	azurePlatform := component.NewPlatform("Azure", "azure")

	awsInstance := component.NewInstance("EKS Cluster", "eks-cluster", awsPlatform)
	azureInstance := component.NewInstance("AKS Cluster", "aks-cluster", azurePlatform)

	awsComponent := component.NewComponent("API Gateway", "api-gateway", awsInstance)
	azureComponent := component.NewComponent("Load Balancer", "load-balancer", azureInstance)

	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:       "multi-platform-incident",
			Title:      "Cross-Platform Network Issues",
			Content:    "Connectivity issues between AWS and Azure services",
			Components: []*component.Component{awsComponent, azureComponent},
			ExtraFields: map[string]string{
				"incident_type": "network",
				"scope":         "multi-cloud",
			},
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityPartialOutage,
	}

	// Test multi-platform component handling
	if len(incident.Components) != 2 {
		t.Errorf("Expected 2 components from different platforms, got %v", len(incident.Components))
	}

	// Verify platform separation in component strings
	awsString := incident.Components[0].String()
	azureString := incident.Components[1].String()

	if !contains(awsString, "AWS") {
		t.Errorf("AWS component should contain 'AWS', got %v", awsString)
	}
	if !contains(azureString, "Azure") {
		t.Errorf("Azure component should contain 'Azure', got %v", azureString)
	}
}

func TestIncident_KnownIssueWithComponents(t *testing.T) {
	// Test known issue (perpetual) with long-term component issues
	platform := component.NewPlatform("Legacy System", "legacy-system")
	instance := component.NewInstance("Mainframe Instance", "mainframe", platform)
	legacyComponent := component.NewComponent("Legacy Database", "legacy-db", instance)

	knownIssue := &Incident{
		BaseEvent: BaseEvent{
			GUID:       "known-issue-legacy-001",
			Title:      "Legacy Database Performance Degradation",
			Content:    "Known performance issues with legacy database during peak hours",
			Components: []*component.Component{legacyComponent},
			ExtraFields: map[string]string{
				"known_since": "2023-01-01",
				"workaround":  "Use alternative read replicas",
				"planned_fix": "Q2 2024",
			},
		},
		Perpetual:           true,
		IncidentCriticality: CriticalityDegraded,
	}

	// Test known issue behavior
	if knownIssue.Type() != TypeKnownIssue {
		t.Errorf("Expected TypeKnownIssue for perpetual incident, got %v", knownIssue.Type())
	}

	if knownIssue.Criticality() != CriticalityDegraded {
		t.Errorf("Expected CriticalityDegraded, got %v", knownIssue.Criticality())
	}

	// Test component integration for known issues
	expectedLegacy := "Legacy System - Mainframe Instance - Legacy Database"
	if knownIssue.Components[0].String() != expectedLegacy {
		t.Errorf("Legacy component String() = %v, want %v", knownIssue.Components[0].String(), expectedLegacy)
	}
}

func TestIncident_EmptyComponentsHandling(t *testing.T) {
	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:       "no-components-incident",
			Title:      "General System Alert",
			Content:    "System-wide alert without specific component attribution",
			Components: []*component.Component{},
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityOperational,
	}

	// Test that empty components don't break incident functionality
	if len(incident.Components) != 0 {
		t.Errorf("Expected 0 components, got %v", len(incident.Components))
	}

	if incident.Type() != TypeFiringIncident {
		t.Errorf("Expected TypeFiringIncident even with no components, got %v", incident.Type())
	}
}

func TestIncident_JSONSerializationWithComponents(t *testing.T) {
	// Create incident with components
	platform := component.NewPlatform("Test Platform", "test-platform")
	instance := component.NewInstance("Test Instance", "test-instance", platform)
	comp := component.NewComponent("Test Component", "test-component", instance)

	now := time.Now()
	incident := &Incident{
		BaseEvent: BaseEvent{
			GUID:           "json-incident-test",
			Title:          "JSON Test Incident",
			Content:        "Testing JSON serialization",
			Components:     []*component.Component{comp},
			StartEffective: &now,
			ExtraFields:    map[string]string{"test": "value"},
		},
		Perpetual:           false,
		IncidentCriticality: CriticalityMajorOutage,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(incident)
	if err != nil {
		t.Fatalf("Failed to marshal incident to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Incident
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal incident from JSON: %v", err)
	}

	// Verify critical fields
	if unmarshaled.GUID != incident.GUID {
		t.Errorf("JSON GUID mismatch: got %v, want %v", unmarshaled.GUID, incident.GUID)
	}
	if unmarshaled.Perpetual != incident.Perpetual {
		t.Errorf("JSON Perpetual mismatch: got %v, want %v", unmarshaled.Perpetual, incident.Perpetual)
	}
	if unmarshaled.IncidentCriticality != incident.IncidentCriticality {
		t.Errorf("JSON Criticality mismatch: got %v, want %v", unmarshaled.IncidentCriticality, incident.IncidentCriticality)
	}
	if len(unmarshaled.Components) != len(incident.Components) {
		t.Errorf("JSON Components length mismatch: got %v, want %v", len(unmarshaled.Components), len(incident.Components))
	}
}

func TestIncident_ComponentStatusCorrelation(t *testing.T) {
	// Test how incident status correlates with component states
	platform := component.NewPlatform("Production Environment", "prod-env")
	instance := component.NewInstance("Web Tier", "web-tier", platform)
	webComponent := component.NewComponent("Load Balancer", "load-balancer", instance)

	now := time.Now()
	tests := []struct {
		name           string
		startTime      *time.Time
		endTime        *time.Time
		expectedStatus Status
		description    string
	}{
		{
			name:           "Ongoing incident with components",
			startTime:      &now,
			endTime:        nil,
			expectedStatus: StatusOnGoing,
			description:    "Active incident affecting components",
		},
		{
			name:           "Resolved incident with components",
			startTime:      &now,
			endTime:        &now,
			expectedStatus: StatusResolved,
			description:    "Resolved incident with restored components",
		},
		{
			name:           "Future planned incident with components",
			startTime:      func() *time.Time { future := now.Add(1 * time.Hour); return &future }(),
			endTime:        func() *time.Time { future := now.Add(2 * time.Hour); return &future }(),
			expectedStatus: StatusUnknown,
			description:    "Future maintenance affecting components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incident := &Incident{
				BaseEvent: BaseEvent{
					GUID:           "status-test-" + tt.name,
					Title:          tt.description,
					Components:     []*component.Component{webComponent},
					StartEffective: tt.startTime,
					EndEffective:   tt.endTime,
				},
				Perpetual:           false,
				IncidentCriticality: CriticalityPartialOutage,
			}

			if incident.Status() != tt.expectedStatus {
				t.Errorf("Incident status = %v, want %v", incident.Status(), tt.expectedStatus)
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

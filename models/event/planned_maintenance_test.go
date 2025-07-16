package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
)

func TestPlannedMaintenance_Type(t *testing.T) {
	pm := &PlannedMaintenance{}
	if got := pm.Type(); got != TypePlannedMaintenance {
		t.Errorf("PlannedMaintenance.Type() = %v, want %v", got, TypePlannedMaintenance)
	}
}

func TestPlannedMaintenance_Status(t *testing.T) {
	now := time.Now()
	past := now.Add(-2 * time.Hour)
	pastHour := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name           string
		startPlanned   time.Time
		startEffective *time.Time
		endEffective   *time.Time
		expected       Status
	}{
		{
			name:           "Planned - Future Start",
			startPlanned:   future,
			startEffective: nil,
			endEffective:   nil,
			expected:       StatusPlanned,
		},
		{
			name:           "Resolved - End in Past",
			startPlanned:   past,
			startEffective: &past,
			endEffective:   &pastHour,
			expected:       StatusResolved,
		},
		{
			name:           "Ongoing - Started, Not Ended",
			startPlanned:   past,
			startEffective: &past,
			endEffective:   &future,
			expected:       StatusOnGoing,
		},
		{
			name:           "Ongoing - Started, No End",
			startPlanned:   past,
			startEffective: &past,
			endEffective:   nil,
			expected:       StatusOnGoing,
		},
		{
			name:           "Unknown - Past Planned, Not Started",
			startPlanned:   past,
			startEffective: nil,
			endEffective:   nil,
			expected:       StatusUnknown,
		},
		{
			name:           "Resolved - End Only in Past",
			startPlanned:   future,
			startEffective: nil,
			endEffective:   &past,
			expected:       StatusResolved,
		},
		{
			name:           "Unknown - Zero Time (Default)",
			startPlanned:   time.Time{},
			startEffective: nil,
			endEffective:   nil,
			expected:       StatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PlannedMaintenance{
				BaseEvent: BaseEvent{
					StartEffective: tt.startEffective,
					EndEffective:   tt.endEffective,
				},
				StartPlanned: tt.startPlanned,
			}
			if got := pm.Status(); got != tt.expected {
				t.Errorf("PlannedMaintenance.Status() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlannedMaintenance_Criticality(t *testing.T) {
	pm := &PlannedMaintenance{}
	if got := pm.Criticality(); got != CriticalityUnderMaintenance {
		t.Errorf("PlannedMaintenance.Criticality() = %v, want %v", got, CriticalityUnderMaintenance)
	}
}

func TestPlannedMaintenance_InterfaceCompliance(t *testing.T) {
	pm := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			Title: "Server Maintenance",
		},
		StartPlanned: time.Now().Add(1 * time.Hour),
		EndPlanned:   time.Now().Add(3 * time.Hour),
	}

	// Test Event interface compliance
	var event Event = pm
	if event.Type() != TypePlannedMaintenance {
		t.Errorf("Expected TypePlannedMaintenance, got %v", event.Type())
	}
	if event.Criticality() != CriticalityUnderMaintenance {
		t.Errorf("Expected CriticalityUnderMaintenance, got %v", event.Criticality())
	}
}

func TestPlannedMaintenance_FullLifecycle(t *testing.T) {
	now := time.Now()
	startPlanned := now.Add(1 * time.Hour)
	endPlanned := now.Add(3 * time.Hour)
	startEffective := now.Add(1*time.Hour + 10*time.Minute)
	endEffective := now.Add(2*time.Hour + 50*time.Minute)

	pm := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:           "maintenance-456",
			Title:          "Weekly Database Maintenance",
			Content:        "Routine database optimization and backup",
			StartEffective: &startEffective,
			EndEffective:   &endEffective,
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
	}

	// Test all methods
	if pm.Type() != TypePlannedMaintenance {
		t.Errorf("Expected TypePlannedMaintenance, got %v", pm.Type())
	}

	if pm.Status() != StatusPlanned {
		t.Errorf("Expected StatusPlanned, got %v", pm.Status())
	}

	if pm.Criticality() != CriticalityUnderMaintenance {
		t.Errorf("Expected CriticalityUnderMaintenance, got %v", pm.Criticality())
	}
}

func TestPlannedMaintenance_StatusTransitions(t *testing.T) {
	now := time.Now()

	pm := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			Title: "Test Maintenance",
		},
		StartPlanned: now.Add(1 * time.Hour),
		EndPlanned:   now.Add(2 * time.Hour),
	}

	// Should be planned initially
	if pm.Status() != StatusPlanned {
		t.Errorf("Expected StatusPlanned, got %v", pm.Status())
	}

	// Start the maintenance
	startTime := now.Add(-10 * time.Minute)
	pm.StartEffective = &startTime
	if pm.Status() != StatusOnGoing {
		t.Errorf("Expected StatusOnGoing, got %v", pm.Status())
	}

	// End the maintenance
	endTime := now.Add(-5 * time.Minute)
	pm.EndEffective = &endTime
	if pm.Status() != StatusResolved {
		t.Errorf("Expected StatusResolved, got %v", pm.Status())
	}
}

func TestPlannedMaintenance_ComponentIntegration(t *testing.T) {
	// Create comprehensive component hierarchy for maintenance
	awsPlatform := component.NewPlatform("AWS Production", "aws-prod")
	eksInstance := component.NewInstance("EKS Cluster v1.25", "eks-cluster", awsPlatform)
	apiComponent := component.NewComponent("API Gateway", "api-gateway", eksInstance)
	dbComponent := component.NewComponent("RDS Primary", "rds-primary", eksInstance)
	cacheComponent := component.NewComponent("ElastiCache Redis", "elasticache-redis", eksInstance)

	now := time.Now()
	startPlanned := now.Add(1 * time.Hour)
	endPlanned := now.Add(3 * time.Hour)

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:       "maintenance-eks-upgrade-001",
			Title:      "EKS Cluster Upgrade to v1.26",
			Content:    "Upgrading EKS cluster with rolling node replacement. Services will remain available with possible brief interruptions.",
			Components: []*component.Component{apiComponent, dbComponent, cacheComponent},
			ExtraFields: map[string]string{
				"upgrade_version":   "1.26.0",
				"expected_duration": "2 hours",
				"impact_level":      "minimal",
				"rollback_plan":     "available",
				"notification_sent": "48h_prior",
			},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
	}

	// Test component integration
	if len(maintenance.Components) != 3 {
		t.Errorf("Expected 3 components for maintenance, got %v", len(maintenance.Components))
	}

	// Test component details
	expectedAPI := "AWS Production - EKS Cluster v1.25 - API Gateway"
	if maintenance.Components[0].String() != expectedAPI {
		t.Errorf("API component String() = %v, want %v", maintenance.Components[0].String(), expectedAPI)
	}

	expectedDB := "AWS Production - EKS Cluster v1.25 - RDS Primary"
	if maintenance.Components[1].String() != expectedDB {
		t.Errorf("DB component String() = %v, want %v", maintenance.Components[1].String(), expectedDB)
	}

	// Test maintenance type and status
	if maintenance.Type() != TypePlannedMaintenance {
		t.Errorf("Expected TypePlannedMaintenance, got %v", maintenance.Type())
	}

	if maintenance.Status() != StatusPlanned {
		t.Errorf("Expected StatusPlanned for future maintenance, got %v", maintenance.Status())
	}

	if maintenance.Criticality() != CriticalityUnderMaintenance {
		t.Errorf("Expected CriticalityUnderMaintenance, got %v", maintenance.Criticality())
	}
}

func TestPlannedMaintenance_MaintenanceLifecycleWithComponents(t *testing.T) {
	// Create components for database maintenance
	platform := component.NewPlatform("Azure SQL", "azure-sql")
	instance := component.NewInstance("Production Database", "prod-db", platform)
	primaryDB := component.NewComponent("Primary Database", "primary-db", instance)
	readReplica := component.NewComponent("Read Replica", "read-replica", instance)

	now := time.Now()
	startPlanned := now.Add(-2 * time.Hour)
	endPlanned := now.Add(-1 * time.Hour)
	startEffective := now.Add(-90 * time.Minute)
	endEffective := now.Add(-30 * time.Minute)

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:           "maintenance-db-patch-001",
			Title:          "Database Security Patch Deployment",
			Content:        "Applying critical security patches to database instances",
			Components:     []*component.Component{primaryDB, readReplica},
			StartEffective: &startEffective,
			EndEffective:   &endEffective,
			ExtraFields: map[string]string{
				"patch_version": "2023.12.01",
				"security_cve":  "CVE-2023-12345",
				"downtime":      "minimal",
			},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
	}

	// Test completed maintenance
	if maintenance.Status() != StatusResolved {
		t.Errorf("Expected StatusResolved for completed maintenance, got %v", maintenance.Status())
	}

	// Test that components are properly associated
	if len(maintenance.Components) != 2 {
		t.Errorf("Expected 2 database components, got %v", len(maintenance.Components))
	}

	// Verify component normalization for tracking (now uses codes)
	expectedPrimaryNorm := "azure-sql-prod-db-primary-db"
	if maintenance.Components[0].Normalize() != expectedPrimaryNorm {
		t.Errorf("Primary DB Normalize() = %v, want %v", maintenance.Components[0].Normalize(), expectedPrimaryNorm)
	}
}

func TestPlannedMaintenance_MultiServiceMaintenance(t *testing.T) {
	// Create components spanning multiple services
	kubernetesPlatform := component.NewPlatform("Kubernetes Cluster", "k8s-cluster")
	webInstance := component.NewInstance("Web Services", "web-services", kubernetesPlatform)
	apiInstance := component.NewInstance("API Services", "api-services", kubernetesPlatform)
	dataInstance := component.NewInstance("Data Services", "data-services", kubernetesPlatform)

	frontendComponent := component.NewComponent("Frontend App", "frontend-app", webInstance)
	apiGateway := component.NewComponent("API Gateway", "api-gateway", apiInstance)
	userService := component.NewComponent("User Service", "user-service", apiInstance)
	dataProcessor := component.NewComponent("Data Processor", "data-processor", dataInstance)

	now := time.Now()
	startPlanned := now.Add(6 * time.Hour)
	endPlanned := now.Add(8 * time.Hour)

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:    "maintenance-multi-service-001",
			Title:   "Kubernetes Cluster Node Refresh",
			Content: "Rolling update of all worker nodes with new instance types for improved performance",
			Components: []*component.Component{
				frontendComponent, apiGateway, userService, dataProcessor,
			},
			ExtraFields: map[string]string{
				"node_type_old":   "m5.large",
				"node_type_new":   "m6i.large",
				"strategy":        "rolling_update",
				"expected_impact": "zero_downtime",
				"nodes_affected":  "12",
			},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
	}

	// Test multi-service component handling
	if len(maintenance.Components) != 4 {
		t.Errorf("Expected 4 components across services, got %v", len(maintenance.Components))
	}

	// Verify each component has the correct hierarchy
	expectedFrontend := "Kubernetes Cluster - Web Services - Frontend App"
	if maintenance.Components[0].String() != expectedFrontend {
		t.Errorf("Frontend component String() = %v, want %v", maintenance.Components[0].String(), expectedFrontend)
	}

	expectedAPI := "Kubernetes Cluster - API Services - API Gateway"
	if maintenance.Components[1].String() != expectedAPI {
		t.Errorf("API Gateway component String() = %v, want %v", maintenance.Components[1].String(), expectedAPI)
	}

	// Test maintenance is properly planned
	if maintenance.Status() != StatusPlanned {
		t.Errorf("Expected StatusPlanned for future maintenance, got %v", maintenance.Status())
	}
}

func TestPlannedMaintenance_CancelledMaintenanceWithComponents(t *testing.T) {
	// Test cancelled maintenance scenario
	platform := component.NewPlatform("Production Environment", "prod-env")
	instance := component.NewInstance("Critical Service", "critical-service", platform)
	component1 := component.NewComponent("Service A", "service-a", instance)
	component2 := component.NewComponent("Service B", "service-b", instance)

	now := time.Now()
	startPlanned := now.Add(-1 * time.Hour)
	endPlanned := now.Add(1 * time.Hour)

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:       "maintenance-cancelled-001",
			Title:      "Cancelled Network Maintenance",
			Content:    "Maintenance cancelled due to ongoing critical incident",
			Components: []*component.Component{component1, component2},
			ExtraFields: map[string]string{
				"cancellation_reason": "critical_incident",
				"rescheduled_date":    "TBD",
				"cancelled_by":        "ops_team",
			},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
		Cancelled:    true,
	}

	// Test cancelled maintenance behavior
	if maintenance.Status() != StatusCanceled {
		t.Errorf("Expected StatusCanceled for cancelled maintenance, got %v", maintenance.Status())
	}

	// Verify components are still properly associated
	if len(maintenance.Components) != 2 {
		t.Errorf("Expected 2 components even for cancelled maintenance, got %v", len(maintenance.Components))
	}
}

func TestPlannedMaintenance_EmergencyMaintenanceWithComponents(t *testing.T) {
	// Test emergency maintenance scenario (started immediately)
	platform := component.NewPlatform("Production Infrastructure", "prod-infra")
	instance := component.NewInstance("Security Layer", "security-layer", platform)
	firewallComponent := component.NewComponent("Network Firewall", "network-firewall", instance)
	loadBalancerComponent := component.NewComponent("Load Balancer", "load-balancer", instance)

	now := time.Now()
	startPlanned := now.Add(-10 * time.Minute) // Was planned recently
	endPlanned := now.Add(30 * time.Minute)
	startEffective := now.Add(-5 * time.Minute) // Started almost immediately

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:           "maintenance-emergency-001",
			Title:          "Emergency Security Patch",
			Content:        "Critical security vulnerability requires immediate patching",
			Components:     []*component.Component{firewallComponent, loadBalancerComponent},
			StartEffective: &startEffective,
			ExtraFields: map[string]string{
				"urgency":           "critical",
				"cve_reference":     "CVE-2023-99999",
				"impact_window":     "30_minutes",
				"business_approval": "emergency_protocol",
			},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
	}

	// Test ongoing emergency maintenance
	if maintenance.Status() != StatusOnGoing {
		t.Errorf("Expected StatusOnGoing for active emergency maintenance, got %v", maintenance.Status())
	}

	// Verify component integration for emergency scenarios
	expectedFirewall := "Production Infrastructure - Security Layer - Network Firewall"
	if maintenance.Components[0].String() != expectedFirewall {
		t.Errorf("Firewall component String() = %v, want %v", maintenance.Components[0].String(), expectedFirewall)
	}
}

func TestPlannedMaintenance_JSONSerializationWithComponents(t *testing.T) {
	// Create maintenance with components
	platform := component.NewPlatform("Test Platform", "test-platform")
	instance := component.NewInstance("Test Instance", "test-instance", platform)
	comp := component.NewComponent("Test Component", "test-component", instance)

	now := time.Now()
	startPlanned := now.Add(1 * time.Hour)
	endPlanned := now.Add(2 * time.Hour)

	maintenance := &PlannedMaintenance{
		BaseEvent: BaseEvent{
			GUID:        "json-maintenance-test",
			Title:       "JSON Test Maintenance",
			Content:     "Testing JSON serialization",
			Components:  []*component.Component{comp},
			ExtraFields: map[string]string{"test": "value"},
		},
		StartPlanned: startPlanned,
		EndPlanned:   endPlanned,
		Cancelled:    false,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(maintenance)
	if err != nil {
		t.Fatalf("Failed to marshal maintenance to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled PlannedMaintenance
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal maintenance from JSON: %v", err)
	}

	// Verify critical fields
	if unmarshaled.GUID != maintenance.GUID {
		t.Errorf("JSON GUID mismatch: got %v, want %v", unmarshaled.GUID, maintenance.GUID)
	}
	if unmarshaled.Cancelled != maintenance.Cancelled {
		t.Errorf("JSON Cancelled mismatch: got %v, want %v", unmarshaled.Cancelled, maintenance.Cancelled)
	}
	if !unmarshaled.StartPlanned.Equal(maintenance.StartPlanned) {
		t.Errorf("JSON StartPlanned mismatch: got %v, want %v", unmarshaled.StartPlanned, maintenance.StartPlanned)
	}
	if len(unmarshaled.Components) != len(maintenance.Components) {
		t.Errorf("JSON Components length mismatch: got %v, want %v", len(unmarshaled.Components), len(maintenance.Components))
	}
}

func TestPlannedMaintenance_ComponentSchedulingCorrelation(t *testing.T) {
	// Test how maintenance scheduling correlates with component availability
	platform := component.NewPlatform("E-commerce Platform", "ecommerce-platform")
	instance := component.NewInstance("Shopping Service", "shopping-service", platform)
	paymentComponent := component.NewComponent("Payment Gateway", "payment-gateway", instance)
	cartComponent := component.NewComponent("Shopping Cart", "shopping-cart", instance)

	now := time.Now()

	tests := []struct {
		name           string
		startPlanned   time.Time
		endPlanned     time.Time
		startEffective *time.Time
		endEffective   *time.Time
		cancelled      bool
		expectedStatus Status
		description    string
	}{
		{
			name:           "Scheduled maintenance",
			startPlanned:   now.Add(2 * time.Hour),
			endPlanned:     now.Add(3 * time.Hour),
			startEffective: nil,
			endEffective:   nil,
			cancelled:      false,
			expectedStatus: StatusPlanned,
			description:    "Future scheduled maintenance affecting payment systems",
		},
		{
			name:           "Active maintenance",
			startPlanned:   now.Add(-30 * time.Minute),
			endPlanned:     now.Add(30 * time.Minute),
			startEffective: func() *time.Time { t := now.Add(-20 * time.Minute); return &t }(),
			endEffective:   nil,
			cancelled:      false,
			expectedStatus: StatusOnGoing,
			description:    "Currently active maintenance on payment components",
		},
		{
			name:           "Completed maintenance",
			startPlanned:   now.Add(-2 * time.Hour),
			endPlanned:     now.Add(-1 * time.Hour),
			startEffective: func() *time.Time { t := now.Add(-2 * time.Hour); return &t }(),
			endEffective:   func() *time.Time { t := now.Add(-1 * time.Hour); return &t }(),
			cancelled:      false,
			expectedStatus: StatusResolved,
			description:    "Completed maintenance with restored components",
		},
		{
			name:           "Cancelled maintenance",
			startPlanned:   now.Add(1 * time.Hour),
			endPlanned:     now.Add(2 * time.Hour),
			startEffective: nil,
			endEffective:   nil,
			cancelled:      true,
			expectedStatus: StatusCanceled,
			description:    "Cancelled maintenance due to business requirements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maintenance := &PlannedMaintenance{
				BaseEvent: BaseEvent{
					GUID:           "schedule-test-" + tt.name,
					Title:          tt.description,
					Components:     []*component.Component{paymentComponent, cartComponent},
					StartEffective: tt.startEffective,
					EndEffective:   tt.endEffective,
					ExtraFields: map[string]string{
						"business_impact":       "payment_services",
						"customer_notification": "24h_prior",
					},
				},
				StartPlanned: tt.startPlanned,
				EndPlanned:   tt.endPlanned,
				Cancelled:    tt.cancelled,
			}

			if maintenance.Status() != tt.expectedStatus {
				t.Errorf("Maintenance status = %v, want %v", maintenance.Status(), tt.expectedStatus)
			}

			// Verify components are maintained through status changes
			if len(maintenance.Components) != 2 {
				t.Errorf("Expected 2 components regardless of status, got %v", len(maintenance.Components))
			}
		})
	}
}

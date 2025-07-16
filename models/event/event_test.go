package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gmllt/clariti/models/component"
)

// Mock implementations for testing
type MockPlatform struct {
	name string
}

func (m *MockPlatform) String() string {
	return m.name
}

func (m *MockPlatform) Normalize() string {
	return m.name
}

type MockInstance struct {
	name     string
	platform *MockPlatform
}

func (m *MockInstance) String() string {
	if m.platform != nil {
		return m.platform.String() + " - " + m.name
	}
	return m.name
}

func (m *MockInstance) Normalize() string {
	return m.name
}

type MockComponent struct {
	name     string
	instance *MockInstance
}

func (m *MockComponent) String() string {
	if m.instance != nil {
		return m.instance.String() + " - " + m.name
	}
	return m.name
}

func (m *MockComponent) Normalize() string {
	return m.name
}

func TestCriticality_String(t *testing.T) {
	tests := []struct {
		name     string
		c        Criticality
		expected string
	}{
		{"Operational", CriticalityOperational, "operational"},
		{"Degraded", CriticalityDegraded, "degraded"},
		{"Partial Outage", CriticalityPartialOutage, "partial outage"},
		{"Major Outage", CriticalityMajorOutage, "major outage"},
		{"Under Maintenance", CriticalityUnderMaintenance, "under maintenance"},
		{"Unknown", CriticalityUnknown, "unknown"},
		{"Invalid Value", Criticality(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.expected {
				t.Errorf("Criticality.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCriticality_Values(t *testing.T) {
	tests := []struct {
		name     string
		c        Criticality
		expected int
	}{
		{"Operational", CriticalityOperational, 0},
		{"Degraded", CriticalityDegraded, 1},
		{"Partial Outage", CriticalityPartialOutage, 2},
		{"Major Outage", CriticalityMajorOutage, 3},
		{"Under Maintenance", CriticalityUnderMaintenance, 4},
		{"Unknown", CriticalityUnknown, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.c) != tt.expected {
				t.Errorf("Criticality value = %v, want %v", int(tt.c), tt.expected)
			}
		})
	}
}

func TestBaseEvent_ComponentIntegration(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	// Create test components using real component package
	platform := component.NewPlatform("AWS EKS Production", "aws-eks-prod")
	instance := component.NewInstance("Web Application Pod", "web-app-pod", platform)
	webComponent := component.NewComponent("Nginx Frontend", "nginx-frontend", instance)
	apiComponent := component.NewComponent("API Gateway", "api-gateway", instance)

	event := BaseEvent{
		GUID:           "test-guid-123",
		Title:          "Multiple Component Event",
		Content:        "Integration test with multiple components",
		ExtraFields:    map[string]string{"severity": "high", "region": "us-east-1"},
		Components:     []*component.Component{webComponent, apiComponent},
		StartEffective: &past,
		EndEffective:   &future,
	}

	// Test component integration
	if len(event.Components) != 2 {
		t.Errorf("Expected 2 components, got %v", len(event.Components))
	}

	// Test component details
	expectedWeb := "AWS EKS Production - Web Application Pod - Nginx Frontend"
	if event.Components[0].String() != expectedWeb {
		t.Errorf("First component String() = %v, want %v", event.Components[0].String(), expectedWeb)
	}

	expectedAPI := "AWS EKS Production - Web Application Pod - API Gateway"
	if event.Components[1].String() != expectedAPI {
		t.Errorf("Second component String() = %v, want %v", event.Components[1].String(), expectedAPI)
	}

	// Test normalized component names (now uses codes)
	expectedWebNorm := "aws-eks-prod-web-app-pod-nginx-frontend"
	if event.Components[0].Normalize() != expectedWebNorm {
		t.Errorf("First component Normalize() = %v, want %v", event.Components[0].Normalize(), expectedWebNorm)
	}
}

func TestBaseEvent_EmptyComponents(t *testing.T) {
	event := BaseEvent{
		GUID:        "empty-components",
		Title:       "Event without components",
		Content:     "Test event with no components",
		Components:  []*component.Component{},
		ExtraFields: make(map[string]string),
	}

	if len(event.Components) != 0 {
		t.Errorf("Expected 0 components, got %v", len(event.Components))
	}

	// Test that empty components don't cause issues
	if event.GUID != "empty-components" {
		t.Errorf("Expected GUID 'empty-components', got %v", event.GUID)
	}
}

func TestBaseEvent_NilComponents(t *testing.T) {
	event := BaseEvent{
		GUID:       "nil-components",
		Title:      "Event with nil components",
		Content:    "Test event with nil components slice",
		Components: nil,
	}

	if event.Components != nil {
		t.Errorf("Expected nil components, got %v", event.Components)
	}
}

func TestBaseEvent_JSONSerialization(t *testing.T) {
	// Create components
	platform := component.NewPlatform("Test Platform", "test-platform")
	instance := component.NewInstance("Test Instance", "test-instance", platform)
	comp := component.NewComponent("Test Component", "test-component", instance)

	now := time.Now()
	event := BaseEvent{
		GUID:           "json-test-guid",
		Title:          "JSON Serialization Test",
		Content:        "Testing JSON serialization with components",
		ExtraFields:    map[string]string{"test": "value"},
		Components:     []*component.Component{comp},
		StartEffective: &now,
		EndEffective:   &now,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal BaseEvent to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled BaseEvent
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal BaseEvent from JSON: %v", err)
	}

	// Verify basic fields
	if unmarshaled.GUID != event.GUID {
		t.Errorf("JSON GUID mismatch: got %v, want %v", unmarshaled.GUID, event.GUID)
	}
	if unmarshaled.Title != event.Title {
		t.Errorf("JSON Title mismatch: got %v, want %v", unmarshaled.Title, event.Title)
	}
	if len(unmarshaled.Components) != len(event.Components) {
		t.Errorf("JSON Components length mismatch: got %v, want %v", len(unmarshaled.Components), len(event.Components))
	}
}

func TestBaseEvent_TimestampHandling(t *testing.T) {
	now := time.Now()
	past := now.Add(-2 * time.Hour)
	future := now.Add(2 * time.Hour)

	tests := []struct {
		name           string
		startTime      *time.Time
		endTime        *time.Time
		expectDuration bool
	}{
		{
			name:           "Both timestamps set",
			startTime:      &past,
			endTime:        &future,
			expectDuration: true,
		},
		{
			name:           "Only start time set",
			startTime:      &past,
			endTime:        nil,
			expectDuration: false,
		},
		{
			name:           "Only end time set",
			startTime:      nil,
			endTime:        &future,
			expectDuration: false,
		},
		{
			name:           "No timestamps",
			startTime:      nil,
			endTime:        nil,
			expectDuration: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := BaseEvent{
				GUID:           "time-test",
				Title:          "Timestamp Test",
				StartEffective: tt.startTime,
				EndEffective:   tt.endTime,
			}

			// Test timestamp values
			if (event.StartEffective != nil) != (tt.startTime != nil) {
				t.Errorf("StartEffective presence mismatch")
			}
			if (event.EndEffective != nil) != (tt.endTime != nil) {
				t.Errorf("EndEffective presence mismatch")
			}

			// Test duration calculation if both are set
			if tt.expectDuration && event.StartEffective != nil && event.EndEffective != nil {
				duration := event.EndEffective.Sub(*event.StartEffective)
				expectedDuration := 4 * time.Hour
				if duration != expectedDuration {
					t.Errorf("Duration mismatch: got %v, want %v", duration, expectedDuration)
				}
			}
		})
	}
}

func TestBaseEvent_ExtraFieldsHandling(t *testing.T) {
	tests := []struct {
		name        string
		extraFields map[string]string
		key         string
		expectedVal string
		shouldExist bool
	}{
		{
			name:        "Normal extra fields",
			extraFields: map[string]string{"severity": "high", "region": "us-east-1"},
			key:         "severity",
			expectedVal: "high",
			shouldExist: true,
		},
		{
			name:        "Empty extra fields",
			extraFields: map[string]string{},
			key:         "nonexistent",
			expectedVal: "",
			shouldExist: false,
		},
		{
			name:        "Nil extra fields",
			extraFields: nil,
			key:         "nonexistent",
			expectedVal: "",
			shouldExist: false,
		},
		{
			name:        "Special characters in values",
			extraFields: map[string]string{"description": "Alert with special chars: !@#$%^&*()"},
			key:         "description",
			expectedVal: "Alert with special chars: !@#$%^&*()",
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := BaseEvent{
				GUID:        "extra-fields-test",
				Title:       "Extra Fields Test",
				ExtraFields: tt.extraFields,
			}

			val, exists := event.ExtraFields[tt.key]
			if exists != tt.shouldExist {
				t.Errorf("Key existence mismatch for %v: got %v, want %v", tt.key, exists, tt.shouldExist)
			}
			if exists && val != tt.expectedVal {
				t.Errorf("Value mismatch for %v: got %v, want %v", tt.key, val, tt.expectedVal)
			}
		})
	}
}

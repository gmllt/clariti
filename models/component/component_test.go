package component

import (
	"strings"
	"testing"

	"github.com/gmllt/clariti/utils"
)

// TestComponentCreation tests basic component creation
func TestComponentCreation(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	component := NewComponent("API Gateway", "api-gateway", instance)

	if component.Name != "API Gateway" {
		t.Errorf("Expected name 'API Gateway', got %s", component.Name)
	}
	if component.Code != "api-gateway" {
		t.Errorf("Expected code 'api-gateway', got %s", component.Code)
	}
	if component.Instance != instance {
		t.Errorf("Expected instance to be set correctly")
	}
}

// TestComponentString tests component string representation
func TestComponentString(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	component := NewComponent("API Gateway", "api-gateway", instance)

	expected := "AWS Production - EKS Cluster - API Gateway"
	result := component.String()

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestComponentNormalize tests component normalization
func TestComponentNormalize(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	component := NewComponent("API Gateway", "api-gateway", instance)

	expected := "aws-prod-eks-cluster-api-gateway"
	result := component.Normalize()

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestComponentNormalizeFromName tests normalization when code is empty
func TestComponentNormalizeFromName(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	component := NewComponent("API Gateway", "", instance)

	// Should use normalized name when code is empty
	result := component.Normalize()
	if !strings.Contains(result, "api-gateway") {
		t.Errorf("Expected normalized result to contain 'api-gateway', got '%s'", result)
	}
}

// TestComponentWithEmptyCode tests component with empty code
func TestComponentWithEmptyCode(t *testing.T) {
	platform := NewPlatform("Test Platform", "test-platform")
	instance := NewInstance("Test Instance", "test-instance", platform)
	component := NewComponent("Test Component", "", instance)

	if component.Name != "Test Component" {
		t.Errorf("Expected name 'Test Component', got %s", component.Name)
	}
	if component.Code != "" {
		t.Errorf("Expected empty code, got %s", component.Code)
	}
}

// TestComponentInterface tests that Component implements required interfaces
func TestComponentInterface(t *testing.T) {
	platform := NewPlatform("Test Platform", "test-platform")
	instance := NewInstance("Test Instance", "test-instance", platform)
	component := NewComponent("Test Component", "test-component", instance)

	// Test Stringable interface
	var stringable utils.Stringable = component
	if stringable.String() == "" {
		t.Error("Component should implement Stringable with non-empty result")
	}

	// Test Normalizable interface
	var normalizable utils.Normalizable = component
	if normalizable.Normalize() == "" {
		t.Error("Component should implement Normalizable with non-empty result")
	}
}

package component

import (
	"strings"
	"testing"

	"github.com/gmllt/clariti/utils"
)

// TestInstanceCreation tests basic instance creation
func TestInstanceCreation(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)

	if instance.Name != "EKS Cluster" {
		t.Errorf("Expected name 'EKS Cluster', got %s", instance.Name)
	}
	if instance.Code != "eks-cluster" {
		t.Errorf("Expected code 'eks-cluster', got %s", instance.Code)
	}
	if instance.Platform != platform {
		t.Errorf("Expected platform to be set correctly")
	}
}

// TestInstanceString tests instance string representation
func TestInstanceString(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)

	expected := "AWS Production - EKS Cluster"
	result := instance.String()

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestInstanceNormalize tests instance normalization
func TestInstanceNormalize(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)

	expected := "aws-prod-eks-cluster"
	result := instance.Normalize()

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestInstanceNormalizeFromName tests normalization when code is empty
func TestInstanceNormalizeFromName(t *testing.T) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "", platform)

	// Should use normalized name when code is empty
	result := instance.Normalize()
	if !strings.Contains(result, "eks-cluster") {
		t.Errorf("Expected normalized result to contain 'eks-cluster', got '%s'", result)
	}
}

// TestInstanceWithEmptyCode tests instance with empty code
func TestInstanceWithEmptyCode(t *testing.T) {
	platform := NewPlatform("Test Platform", "test-platform")
	instance := NewInstance("Test Instance", "", platform)

	if instance.Name != "Test Instance" {
		t.Errorf("Expected name 'Test Instance', got %s", instance.Name)
	}
	if instance.Code != "" {
		t.Errorf("Expected empty code, got %s", instance.Code)
	}
}

// TestInstanceInterface tests that Instance implements required interfaces
func TestInstanceInterface(t *testing.T) {
	platform := NewPlatform("Test Platform", "test-platform")
	instance := NewInstance("Test Instance", "test-instance", platform)

	// Test Stringable interface
	var stringable utils.Stringable = instance
	if stringable.String() == "" {
		t.Error("Instance should implement Stringable with non-empty result")
	}

	// Test Normalizable interface
	var normalizable utils.Normalizable = instance
	if normalizable.Normalize() == "" {
		t.Error("Instance should implement Normalizable with non-empty result")
	}
}

// TestInstanceHierarchy tests the hierarchical nature of instances
func TestInstanceHierarchy(t *testing.T) {
	platform := NewPlatform("Production Cloud", "prod-cloud")
	instance := NewInstance("Kubernetes Cluster", "k8s-cluster", platform)

	// Instance string should include platform
	instanceStr := instance.String()
	if !strings.Contains(instanceStr, platform.Name) {
		t.Errorf("Instance string should contain platform name: %s", instanceStr)
	}

	// Instance normalization should include platform
	instanceNorm := instance.Normalize()
	if !strings.Contains(instanceNorm, platform.Code) {
		t.Errorf("Instance normalization should contain platform code: %s", instanceNorm)
	}
}

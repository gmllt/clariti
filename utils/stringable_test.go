package utils

import (
	"testing"
)

// Mock types for testing
type mockPlatform struct {
	name, code string
}

func (p *mockPlatform) String() string {
	if p.name != "" {
		return p.name
	}
	return "unknown platform"
}

func (p *mockPlatform) Normalize() string {
	return StrUtils.NormalizeWithFallback(p.code, p.String(), "unknown-platform")
}

func (p *mockPlatform) GetParent() HierarchicalStringable { return nil }
func (p *mockPlatform) GetName() string                   { return p.name }
func (p *mockPlatform) GetCode() string                   { return p.code }
func (p *mockPlatform) GetDefaultNormalized() string      { return "unknown-platform" }

type mockInstance struct {
	name, code string
	platform   *mockPlatform
}

func (i *mockInstance) String() string {
	return StrUtils.BuildHierarchicalStringForComponent(i)
}

func (i *mockInstance) Normalize() string {
	return StrUtils.NormalizeHierarchicalComponent(i)
}

func (i *mockInstance) GetParent() HierarchicalStringable { return i.platform }
func (i *mockInstance) GetName() string                   { return i.name }
func (i *mockInstance) GetCode() string                   { return i.code }
func (i *mockInstance) GetDefaultNormalized() string      { return "unknown-instance" }

type mockComponent struct {
	name, code string
	instance   *mockInstance
}

func (c *mockComponent) String() string {
	return StrUtils.BuildHierarchicalStringForComponent(c)
}

func (c *mockComponent) Normalize() string {
	return StrUtils.NormalizeHierarchicalComponent(c)
}

func (c *mockComponent) GetParent() HierarchicalStringable { return c.instance }
func (c *mockComponent) GetName() string                   { return c.name }
func (c *mockComponent) GetCode() string                   { return c.code }
func (c *mockComponent) GetDefaultNormalized() string      { return "unknown-component" }

// Mock constructors
func NewMockPlatform(name, code string) *mockPlatform {
	return &mockPlatform{name: name, code: code}
}

func NewMockInstance(name, code string, platform *mockPlatform) *mockInstance {
	return &mockInstance{name: name, code: code, platform: platform}
}

func NewMockComponent(name, code string, instance *mockInstance) *mockComponent {
	return &mockComponent{name: name, code: code, instance: instance}
}

func Test_normalizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty String",
			input:    "",
			expected: "",
		},
		{
			name:     "Simple String",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Alphanumeric Only",
			input:    "Hello123World",
			expected: "hello123world",
		},
		{
			name:     "Special Characters",
			input:    "Hello@#$%World!!!",
			expected: "hello-world",
		},
		{
			name:     "Multiple Consecutive Special Chars",
			input:    "Hello---World___Test",
			expected: "hello-world-test",
		},
		{
			name:     "Only Special Characters",
			input:    "@#$%^&*()",
			expected: "",
		},
		{
			name:     "Real World Example",
			input:    "AWS Production Environment v2.0",
			expected: "aws-production-environment-v2-0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeString(tt.input); got != tt.expected {
				t.Errorf("normalizeString() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNormalizeFromStringable(t *testing.T) {
	tests := []struct {
		name     string
		input    Stringable
		expected string
	}{
		{
			name:     "Nil Stringable",
			input:    nil,
			expected: "",
		},
		{
			name:     "Real Platform",
			input:    NewMockPlatform("AWS Production", "aws-prod"),
			expected: "aws-production",
		},
		{
			name:     "Real Instance",
			input:    NewMockInstance("EKS Cluster", "eks", NewMockPlatform("AWS", "aws")),
			expected: "aws-eks-cluster",
		},
		{
			name:     "Real Component",
			input:    NewMockComponent("API Gateway", "api-gw", NewMockInstance("EKS", "eks", NewMockPlatform("AWS", "aws"))),
			expected: "aws-eks-api-gateway",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeFromStringable(tt.input); got != tt.expected {
				t.Errorf("NormalizeFromStringable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringableInterface(t *testing.T) {
	platform := NewMockPlatform("Test Platform", "test")
	instance := NewMockInstance("Test Instance", "test-inst", platform)
	component := NewMockComponent("Test Component", "test-comp", instance)

	// Test that all types implement Stringable
	var stringables []Stringable = []Stringable{platform, instance, component}

	for i, s := range stringables {
		if s.String() == "" {
			t.Errorf("Stringable %d returned empty string", i)
		}
	}
}

func TestNormalizableInterface(t *testing.T) {
	platform := NewMockPlatform("Test Platform", "test")
	instance := NewMockInstance("Test Instance", "test-inst", platform)
	component := NewMockComponent("Test Component", "test-comp", instance)

	// Test that all types implement Normalizable
	var normalizables []Normalizable = []Normalizable{platform, instance, component}

	for i, n := range normalizables {
		if n.Normalize() == "" {
			t.Errorf("Normalizable %d returned empty normalized string", i)
		}
		if n.String() == "" {
			t.Errorf("Normalizable %d String() method returned empty string", i)
		}
	}
}

func TestNormalization_Consistency(t *testing.T) {
	platform := NewMockPlatform("AWS Production Environment", "aws-prod")
	instance := NewMockInstance("EKS Cluster v1.25", "eks-cluster", platform)
	component := NewMockComponent("API Gateway Service", "api-gw", instance)

	// Test that normalization is consistent
	platformNorm1 := platform.Normalize()
	platformNorm2 := platform.Normalize()
	if platformNorm1 != platformNorm2 {
		t.Errorf("Platform normalization inconsistent: %v != %v", platformNorm1, platformNorm2)
	}

	instanceNorm1 := instance.Normalize()
	instanceNorm2 := instance.Normalize()
	if instanceNorm1 != instanceNorm2 {
		t.Errorf("Instance normalization inconsistent: %v != %v", instanceNorm1, instanceNorm2)
	}

	componentNorm1 := component.Normalize()
	componentNorm2 := component.Normalize()
	if componentNorm1 != componentNorm2 {
		t.Errorf("Component normalization inconsistent: %v != %v", componentNorm1, componentNorm2)
	}
}

// Benchmark tests
func BenchmarkNormalizeString_Simple(b *testing.B) {
	input := "Simple Component"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalizeString(input)
	}
}

func BenchmarkNormalizeString_Complex(b *testing.B) {
	input := "Complex Component!!! With Many Special Characters @#$%^&*() 123"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalizeString(input)
	}
}

func BenchmarkNormalizeString_VeryLong(b *testing.B) {
	input := "Very Long Component Name With Many Words And Special Characters That Should Test Memory Allocation Performance @#$%^&*()_+{}|:<>?[];'\"\\,./`~1234567890-="
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalizeString(input)
	}
}

func BenchmarkNormalizeFromStringable(b *testing.B) {
	platform := NewMockPlatform("Benchmark Platform", "bench")
	instance := NewMockInstance("Benchmark Instance", "bench-inst", platform)
	component := NewMockComponent("Benchmark Component", "bench-comp", instance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NormalizeFromStringable(component)
	}
}

func BenchmarkComponent_Normalize(b *testing.B) {
	platform := NewMockPlatform("AWS Production Environment", "aws-prod")
	instance := NewMockInstance("EKS Cluster v1.25", "eks-cluster", platform)
	component := NewMockComponent("API Gateway Service", "api-gateway", instance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		component.Normalize()
	}
}

func BenchmarkInstance_Normalize(b *testing.B) {
	platform := NewMockPlatform("AWS Production Environment", "aws-prod")
	instance := NewMockInstance("EKS Cluster v1.25", "eks-cluster", platform)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		instance.Normalize()
	}
}

func BenchmarkComponent_String(b *testing.B) {
	platform := NewMockPlatform("AWS Production Environment", "aws-prod")
	instance := NewMockInstance("EKS Cluster v1.25", "eks-cluster", platform)
	component := NewMockComponent("API Gateway Service", "api-gateway", instance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = component.String()
	}
}

func BenchmarkConcurrentNormalize(b *testing.B) {
	platform := NewMockPlatform("AWS Production Environment", "aws-prod")
	instance := NewMockInstance("EKS Cluster v1.25", "eks-cluster", platform)
	component := NewMockComponent("API Gateway Service", "api-gateway", instance)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			component.Normalize()
		}
	})
}

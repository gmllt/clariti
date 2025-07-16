package component

import (
	"testing"

	"github.com/gmllt/clariti/utils"
)

func TestNewPlatform(t *testing.T) {
	tests := []struct {
		name         string
		platformName string
		platformCode string
		wantName     string
		wantCode     string
	}{
		{
			name:         "Valid platform",
			platformName: "AWS Production",
			platformCode: "aws-prod",
			wantName:     "AWS Production",
			wantCode:     "aws-prod",
		},
		{
			name:         "Empty platform name",
			platformName: "",
			platformCode: "empty",
			wantName:     "",
			wantCode:     "empty",
		},
		{
			name:         "Platform with special characters",
			platformName: "Azure@#$%",
			platformCode: "azure",
			wantName:     "Azure@#$%",
			wantCode:     "azure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform := NewPlatform(tt.platformName, tt.platformCode)

			if platform.Name != tt.wantName {
				t.Errorf("Platform.Name = %v, want %v", platform.Name, tt.wantName)
			}
			if platform.Code != tt.wantCode {
				t.Errorf("Platform.Code = %v, want %v", platform.Code, tt.wantCode)
			}
		})
	}
}

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		name     string
		platform *Platform
		expected string
	}{
		{
			name:     "Platform with name",
			platform: NewPlatform("AWS", "aws"),
			expected: "AWS",
		},
		{
			name:     "Platform with empty name",
			platform: NewPlatform("", "empty"),
			expected: "unknown platform",
		},
		{
			name:     "Platform with complex name",
			platform: NewPlatform("AWS Production Environment", "aws-prod"),
			expected: "AWS Production Environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.String(); got != tt.expected {
				t.Errorf("Platform.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlatform_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		platform *Platform
		expected string
	}{
		{
			name:     "Platform with code",
			platform: NewPlatform("AWS Production", "aws-prod"),
			expected: "aws-prod",
		},
		{
			name:     "Platform without code",
			platform: NewPlatform("Simple Platform", ""),
			expected: "simple-platform",
		},
		{
			name:     "Platform with special characters in code",
			platform: NewPlatform("Test Platform", "test@#$"),
			expected: "test",
		},
		{
			name:     "Empty platform name and code",
			platform: NewPlatform("", ""),
			expected: "unknown-platform",
		},
		{
			name:     "Platform with only special characters",
			platform: NewPlatform("@#$%^&*()", ""),
			expected: "unknown-platform",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.platform.Normalize(); got != tt.expected {
				t.Errorf("Platform.Normalize() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPlatform_InterfaceCompliance(t *testing.T) {
	platform := NewPlatform("Test Platform", "test")

	// Test Stringable interface
	var s utils.Stringable = platform
	if s.String() == "" {
		t.Error("Platform should implement Stringable interface")
	}

	// Test Normalizable interface
	var n utils.Normalizable = platform
	if n.Normalize() == "" {
		t.Error("Platform should implement Normalizable interface")
	}
	if n.String() == "" {
		t.Error("Normalizable should also have String() method")
	}
}

// Benchmark tests
func BenchmarkPlatform_String(b *testing.B) {
	platform := NewPlatform("AWS", "aws")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = platform.String()
	}
}

func BenchmarkPlatform_Normalize(b *testing.B) {
	platform := NewPlatform("AWS Production Environment", "aws-prod")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		platform.Normalize()
	}
}

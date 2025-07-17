package utils

import (
	"strings"
	"testing"
)

func TestWithBuilder(t *testing.T) {
	result := WithBuilder(func(builder *strings.Builder) string {
		builder.WriteString("Hello")
		builder.WriteString(" ")
		builder.WriteString("World")
		return builder.String()
	})

	expected := "Hello World"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestWithBuilderCapacity(t *testing.T) {
	input := "Test string with capacity"
	result := WithBuilderCapacity(len(input), func(builder *strings.Builder) string {
		builder.WriteString(input)
		return builder.String()
	})

	if result != input {
		t.Errorf("Expected %q, got %q", input, result)
	}
}

func TestWithBuilderEmpty(t *testing.T) {
	result := WithBuilder(func(builder *strings.Builder) string {
		return builder.String()
	})

	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func BenchmarkWithBuilderSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithBuilder(func(builder *strings.Builder) string {
			builder.WriteString("test")
			return builder.String()
		})
	}
}

func BenchmarkWithBuilderCapacitySmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithBuilderCapacity(4, func(builder *strings.Builder) string {
			builder.WriteString("test")
			return builder.String()
		})
	}
}

func BenchmarkTraditionalPattern(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := GetBuilder()
		builder.Grow(4)
		builder.WriteString("test")
		result := builder.String()
		PutBuilder(builder)
		_ = result
	}
}

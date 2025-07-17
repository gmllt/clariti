package utils

import (
	"testing"
)

// Simple benchmarks for the consolidated pools
func BenchmarkConsolidatedPools_StringBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := GetBuilder()
		builder.WriteString("test-")
		builder.WriteString("consolidated")
		result := builder.String()
		PutBuilder(builder)
		_ = result
	}
}

func BenchmarkConsolidatedPools_JSONMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := GetJSONResponse()
		m["service"] = "clariti"
		m["status"] = "operational"
		PutJSONResponse(m)
	}
}

func BenchmarkConsolidatedPools_StringMap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := GetStringMap()
		m["method"] = "GET"
		m["status"] = "200"
		PutStringMap(m)
	}
}

// Compare optimized normalization
func BenchmarkOptimizedNormalization(b *testing.B) {
	input := "Test String With Spaces"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NormalizeStringPooled(input)
	}
}

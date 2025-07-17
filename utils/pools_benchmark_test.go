package utils

import (
	"testing"
)

// Benchmark the object pools
func BenchmarkJSONResponsePool_Get(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := GetJSONResponse()
		m["test"] = "value"
		m["number"] = 42
		PutJSONResponse(m)
	}
}

func BenchmarkJSONResponsePool_GetPut(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := GetJSONResponse()
		m["service"] = "clariti"
		m["version"] = "1.0.0"
		m["status"] = "operational"
		PutJSONResponse(m)
	}
}

func BenchmarkStringMapPool_Get(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := GetStringMap()
		m["status"] = "healthy"
		m["service"] = "clariti-api"
		PutStringMap(m)
	}
}

// Compare with traditional map allocation
func BenchmarkTraditionalMapAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]interface{})
		m["service"] = "clariti"
		m["version"] = "1.0.0"
		m["status"] = "operational"
		_ = m
	}
}

func BenchmarkTraditionalStringMapAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[string]string)
		m["status"] = "healthy"
		m["service"] = "clariti-api"
		_ = m
	}
}

// Concurrent pool access benchmark
func BenchmarkJSONResponsePool_Concurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := GetJSONResponse()
			m["service"] = "clariti"
			m["concurrent"] = true
			PutJSONResponse(m)
		}
	})
}

func BenchmarkStringMapPool_Concurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := GetStringMap()
			m["status"] = "healthy"
			m["concurrent"] = "true"
			PutStringMap(m)
		}
	})
}

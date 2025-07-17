package component

import (
	"sync"
	"testing"
)

// BenchmarkComponent_NormalizeWithCache tests the performance of cached normalization
func BenchmarkComponent_NormalizeWithCache(b *testing.B) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := NewComponent("API Gateway", "api-gateway", instance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = comp.Normalize()
	}
}

// BenchmarkComponent_NormalizeWithoutCache tests without caching (for comparison)
func BenchmarkComponent_NormalizeWithoutCache(b *testing.B) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := NewComponent("API Gateway", "api-gateway", instance)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Clear cache each time to simulate no caching
		comp.ClearCache()
		_ = comp.Normalize()
	}
}

// BenchmarkComponent_NormalizeConcurrent tests concurrent access to cached normalization
func BenchmarkComponent_NormalizeConcurrent(b *testing.B) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)
	comp := NewComponent("API Gateway", "api-gateway", instance)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = comp.Normalize()
		}
	})
}

// BenchmarkComponent_CreateMany tests creating many components efficiently
func BenchmarkComponent_CreateMany(b *testing.B) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create 10 components efficiently
		components := make([]*Component, 10)
		for j := 0; j < 10; j++ {
			components[j] = NewComponent("Service", "service", instance)
		}
		_ = components
	}
}

// BenchmarkComponent_CreateManyWithNormalize tests creating and normalizing many components
func BenchmarkComponent_CreateManyWithNormalize(b *testing.B) {
	platform := NewPlatform("AWS Production", "aws-prod")
	instance := NewInstance("EKS Cluster", "eks-cluster", platform)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				comp := NewComponent("Service", "service", instance)
				_ = comp.Normalize()
			}()
		}
		wg.Wait()
	}
}

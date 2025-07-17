package utils

import (
	"regexp"
	"testing"
)

func TestNewGUID(t *testing.T) {
	guid := NewGUID()

	if guid.IsEmpty() {
		t.Error("NewGUID should not return empty GUID")
	}

	// Check UUID v4 format: 8-4-4-4-12
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(guid.String()) {
		t.Errorf("GUID format invalid: %s", guid.String())
	}
}

func TestNewGUIDString(t *testing.T) {
	guidStr := NewGUIDString()

	if guidStr == "" {
		t.Error("NewGUIDString should not return empty string")
	}

	// Check length (UUID v4 is always 36 characters)
	if len(guidStr) != 36 {
		t.Errorf("GUID string length should be 36, got %d", len(guidStr))
	}
}

func TestGUIDUniqueness(t *testing.T) {
	// Generate multiple GUIDs and ensure they're unique
	guids := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		guid := NewGUIDString()
		if guids[guid] {
			t.Errorf("Generated duplicate GUID: %s", guid)
		}
		guids[guid] = true
	}
}

func TestGUID_IsEmpty(t *testing.T) {
	emptyGUID := GUID("")
	nonEmptyGUID := NewGUID()

	if !emptyGUID.IsEmpty() {
		t.Error("Empty GUID should return true for IsEmpty()")
	}

	if nonEmptyGUID.IsEmpty() {
		t.Error("Non-empty GUID should return false for IsEmpty()")
	}
}

func TestGUID_String(t *testing.T) {
	guid := NewGUID()
	str1 := guid.String()
	str2 := string(guid)

	if str1 != str2 {
		t.Error("GUID.String() should match string conversion")
	}
}

// Benchmark tests for performance optimization
func BenchmarkNewGUID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewGUID()
	}
}

func BenchmarkNewGUIDString(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewGUIDString()
	}
}

func BenchmarkGUIDConcurrent(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = NewGUID()
		}
	})
}

func BenchmarkGUID_String(b *testing.B) {
	guid := NewGUID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guid.String()
	}
}

func BenchmarkGUID_IsEmpty(b *testing.B) {
	guid := NewGUID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guid.IsEmpty()
	}
}

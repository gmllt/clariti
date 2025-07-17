package utils

import (
	"strings"
	"sync"
)

// Pool of strings.Builder to reduce allocations
var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// JSONResponsePool is a pool for reusing JSON response maps to reduce allocations
var JSONResponsePool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{}, 8) // Pre-allocate with capacity
	},
}

// StringMapPool is a pool for reusing string maps
var StringMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string, 4)
	},
}

// GetBuilder gets a builder from the pool
func GetBuilder() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

// PutBuilder returns a builder to the pool after resetting it
func PutBuilder(b *strings.Builder) {
	b.Reset()
	builderPool.Put(b)
}

// GetJSONResponse gets a JSON response map from the pool
func GetJSONResponse() map[string]interface{} {
	m := JSONResponsePool.Get().(map[string]interface{})
	// Clear the map for reuse
	for k := range m {
		delete(m, k)
	}
	return m
}

// PutJSONResponse returns a JSON response map to the pool
func PutJSONResponse(m map[string]interface{}) {
	if len(m) > 16 { // Don't pool very large maps
		return
	}
	JSONResponsePool.Put(m)
}

// GetStringMap gets a string map from the pool
func GetStringMap() map[string]string {
	m := StringMapPool.Get().(map[string]string)
	// Clear the map for reuse
	for k := range m {
		delete(m, k)
	}
	return m
}

// PutStringMap returns a string map to the pool
func PutStringMap(m map[string]string) {
	if len(m) > 8 { // Don't pool very large maps
		return
	}
	StringMapPool.Put(m)
}

// NormalizeStringPooled uses a pool of builders to avoid allocations
func NormalizeStringPooled(input string) string {
	if input == "" {
		return ""
	}

	return WithBuilderCapacity(len(input), func(builder *strings.Builder) string {
		lastWasDash := true // Start as true to avoid leading dash

		for _, r := range input {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				// Lowercase letter or digit - add directly
				builder.WriteRune(r)
				lastWasDash = false
			} else if r >= 'A' && r <= 'Z' {
				// Uppercase letter - convert to lowercase
				builder.WriteRune(r - 'A' + 'a')
				lastWasDash = false
			} else if !lastWasDash {
				// Non-alphanumeric - add dash only if last wasn't dash
				builder.WriteByte('-')
				lastWasDash = true
			}
		}

		result := builder.String()
		// Remove trailing dash if present
		if len(result) > 0 && result[len(result)-1] == '-' {
			return result[:len(result)-1]
		}
		return result
	})
}

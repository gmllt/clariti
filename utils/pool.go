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

// GetBuilder gets a builder from the pool
func GetBuilder() *strings.Builder {
	return builderPool.Get().(*strings.Builder)
}

// PutBuilder returns a builder to the pool after resetting it
func PutBuilder(b *strings.Builder) {
	b.Reset()
	builderPool.Put(b)
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

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
)

// GUID represents a globally unique identifier
type GUID string

// guidPool is a pool of byte slices for GUID generation to reduce allocations
var guidPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 16) // 128 bits for UUID v4
	},
}

// NewGUID generates a new UUID v4 compatible GUID with optimized performance
func NewGUID() GUID {
	// Get a byte slice from the pool
	b := guidPool.Get().([]byte)
	defer guidPool.Put(b)

	// Generate random bytes
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a deterministic pattern if random fails
		for i := range b {
			b[i] = byte(i * 17) // Simple deterministic pattern
		}
	}

	// Set version (4) and variant bits for UUID v4
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10

	// Format as standard UUID string using builder wrapper
	guidStr := WithBuilderCapacity(36, func(builder *strings.Builder) string {
		// Convert to hex string in UUID format: 8-4-4-4-12
		builder.WriteString(hex.EncodeToString(b[:4]))
		builder.WriteByte('-')
		builder.WriteString(hex.EncodeToString(b[4:6]))
		builder.WriteByte('-')
		builder.WriteString(hex.EncodeToString(b[6:8]))
		builder.WriteByte('-')
		builder.WriteString(hex.EncodeToString(b[8:10]))
		builder.WriteByte('-')
		builder.WriteString(hex.EncodeToString(b[10:16]))
		return builder.String()
	})

	return GUID(guidStr)
}

// String returns the string representation of the GUID
func (g GUID) String() string {
	return string(g)
}

// IsEmpty returns true if the GUID is empty
func (g GUID) IsEmpty() bool {
	return string(g) == ""
}

// NewGUIDString is a convenience function that returns a string directly
func NewGUIDString() string {
	return NewGUID().String()
}

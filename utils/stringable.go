package utils

import (
	"strings"
)

// Stringable defines the contract for objects that can be converted to string
type Stringable interface {
	String() string
}

// Normalizable defines the contract for objects that can normalize their string representation
type Normalizable interface {
	Stringable
	Normalize() string
}

// normalizeString converts a string to a normalized form suitable for identifiers
func normalizeString(input string) string {
	if input == "" {
		return ""
	}

	// Pre-allocate builder with estimated capacity to reduce allocations
	var builder strings.Builder
	builder.Grow(len(input))

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
		result = result[:len(result)-1]
	}

	return result
}

// NormalizeFromStringable normalizes any object that implements Stringable
func NormalizeFromStringable(s Stringable) string {
	if s == nil {
		return ""
	}
	return normalizeString(s.String())
}

package utils

import "strings"

// BuilderFunc is a function type that receives a builder and returns a string
type BuilderFunc func(*strings.Builder) string

// WithBuilder executes a function with a pooled builder and returns the result
// This eliminates the need for manual GetBuilder/defer PutBuilder pattern
func WithBuilder(fn BuilderFunc) string {
	builder := GetBuilder()
	defer PutBuilder(builder)
	return fn(builder)
}

// WithBuilderCapacity executes a function with a pooled builder pre-grown to the specified capacity
func WithBuilderCapacity(capacity int, fn BuilderFunc) string {
	builder := GetBuilder()
	defer PutBuilder(builder)
	builder.Grow(capacity)
	return fn(builder)
}

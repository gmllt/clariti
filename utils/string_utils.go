package utils

import "strings"

// HierarchicalStringable represents an object that can build hierarchical strings
type HierarchicalStringable interface {
	Stringable
	GetParent() HierarchicalStringable
	GetName() string
	GetCode() string
	GetDefaultNormalized() string
}

// stringUtils provides utility functions for efficient string building
type stringUtils struct{}

// BuildHierarchicalString creates a hierarchical string like "parent - child"
func (stringUtils) BuildHierarchicalString(parent, child string) string {
	if parent == "" {
		return child
	}
	if child == "" {
		return parent
	}

	return WithBuilderCapacity(len(parent)+3+len(child), func(builder *strings.Builder) string {
		builder.WriteString(parent)
		builder.WriteString(" - ")
		builder.WriteString(child)
		return builder.String()
	})
}

// BuildNormalizedPath creates a normalized path like "parent-child"
func (stringUtils) BuildNormalizedPath(parent, child string) string {
	if parent == "" {
		return normalizeString(child)
	}
	if child == "" {
		return parent
	}

	return WithBuilderCapacity(len(parent)+1+len(child), func(builder *strings.Builder) string {
		builder.WriteString(parent)
		builder.WriteByte('-')
		builder.WriteString(normalizeString(child))
		return builder.String()
	})
}

// NormalizeWithFallback normalizes a code if not empty, otherwise normalizes the full string
func (stringUtils) NormalizeWithFallback(code, fullString, defaultValue string) string {
	if code != "" {
		return normalizeString(code)
	}
	normalized := normalizeString(fullString)
	if normalized == "" {
		return defaultValue
	}
	return normalized
}

// BuildHierarchicalStringForComponent builds a hierarchical string for any hierarchical component
func (stringUtils) BuildHierarchicalStringForComponent(h HierarchicalStringable) string {
	parent := h.GetParent()
	if parent != nil {
		return StrUtils.BuildHierarchicalString(parent.String(), h.GetName())
	}
	return h.GetName()
}

// NormalizeHierarchicalComponent normalizes any hierarchical component
func (stringUtils) NormalizeHierarchicalComponent(h HierarchicalStringable) string {
	code := h.GetCode()
	if code != "" {
		parent := h.GetParent()
		if parent != nil {
			return StrUtils.BuildNormalizedPath(parent.(Normalizable).Normalize(), code)
		}
		return normalizeString(code)
	}
	return StrUtils.NormalizeWithFallback("", h.String(), h.GetDefaultNormalized())
}

// Global instance of string utilities
var StrUtils = stringUtils{}

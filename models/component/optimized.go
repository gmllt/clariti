package component

import (
	"strings"

	"github.com/gmllt/clariti/utils"
)

// OptimizedComponent optimized version with pool
type OptimizedComponent struct {
	BaseComponent
	Instance *Instance `json:"instance,omitempty"`
}

// NormalizePooled uses the builder pool
func (c *OptimizedComponent) NormalizePooled() string {
	if c.Code != "" {
		if c.Instance != nil {
			instanceNorm := c.Instance.Normalize()
			return utils.WithBuilderCapacity(len(instanceNorm)+1+len(c.Code), func(builder *strings.Builder) string {
				builder.WriteString(instanceNorm)
				builder.WriteByte('-')
				builder.WriteString(utils.NormalizeStringPooled(c.Code))
				return builder.String()
			})
		}
		return utils.NormalizeStringPooled(c.Code)
	}
	normalized := utils.NormalizeStringPooled(c.String())
	if normalized == "" {
		return "unknown-component"
	}
	return normalized
}

// String for OptimizedComponent
func (c *OptimizedComponent) String() string {
	if c.Instance != nil {
		instanceStr := c.Instance.String()
		return utils.WithBuilderCapacity(len(instanceStr)+3+len(c.Name), func(builder *strings.Builder) string {
			builder.WriteString(instanceStr)
			builder.WriteString(" - ")
			builder.WriteString(c.Name)
			return builder.String()
		})
	}
	return c.Name
}

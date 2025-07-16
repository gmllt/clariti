package component

import (
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
			builder := utils.GetBuilder()
			defer utils.PutBuilder(builder)

			instanceNorm := c.Instance.Normalize()
			builder.Grow(len(instanceNorm) + 1 + len(c.Code))
			builder.WriteString(instanceNorm)
			builder.WriteByte('-')
			builder.WriteString(utils.NormalizeStringPooled(c.Code))
			return builder.String()
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
		builder := utils.GetBuilder()
		defer utils.PutBuilder(builder)

		instanceStr := c.Instance.String()
		builder.Grow(len(instanceStr) + 3 + len(c.Name)) // " - " is 3 chars
		builder.WriteString(instanceStr)
		builder.WriteString(" - ")
		builder.WriteString(c.Name)
		return builder.String()
	}
	return c.Name
}

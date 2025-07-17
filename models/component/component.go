package component

import (
	"sync"

	"github.com/gmllt/clariti/utils"
)

// BaseComponent provides common fields for all component types
type BaseComponent struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Component represents a service component that belongs to an instance
type Component struct {
	BaseComponent
	Instance *Instance `json:"instance,omitempty"`

	// Cache for normalization to avoid repeated string building
	normalizedCache string
	cacheMutex      sync.RWMutex
}

// NewComponent creates a new component with the given name, code and instance
func NewComponent(name, code string, instance *Instance) *Component {
	return &Component{
		BaseComponent: BaseComponent{
			Name: name,
			Code: code,
		},
		Instance: instance,
	}
}

// String returns the string representation of the component
func (c *Component) String() string {
	return utils.StrUtils.BuildHierarchicalStringForComponent(c)
}

// Normalize returns a normalized identifier for the component with caching
func (c *Component) Normalize() string {
	// Check cache first (read lock)
	c.cacheMutex.RLock()
	if c.normalizedCache != "" {
		cached := c.normalizedCache
		c.cacheMutex.RUnlock()
		return cached
	}
	c.cacheMutex.RUnlock()

	// Calculate and cache (write lock)
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	// Double-check pattern in case another goroutine calculated it
	if c.normalizedCache != "" {
		return c.normalizedCache
	}

	c.normalizedCache = utils.StrUtils.NormalizeHierarchicalComponent(c)
	return c.normalizedCache
}

// ClearCache clears the normalization cache (useful when component changes)
func (c *Component) ClearCache() {
	c.cacheMutex.Lock()
	c.normalizedCache = ""
	c.cacheMutex.Unlock()
}

// HierarchicalStringable interface implementation
func (c *Component) GetParent() utils.HierarchicalStringable { return c.Instance }
func (c *Component) GetName() string                         { return c.Name }
func (c *Component) GetCode() string                         { return c.Code }
func (c *Component) GetDefaultNormalized() string            { return "unknown-component" }

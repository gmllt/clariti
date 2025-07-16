package component

import "github.com/gmllt/clariti/utils"

// BaseComponent provides common fields for all component types
type BaseComponent struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Component represents a service component that belongs to an instance
type Component struct {
	BaseComponent
	Instance *Instance `json:"instance,omitempty"`
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

// Normalize returns a normalized identifier for the component
func (c *Component) Normalize() string {
	return utils.StrUtils.NormalizeHierarchicalComponent(c)
}

// HierarchicalStringable interface implementation
func (c *Component) GetParent() utils.HierarchicalStringable { return c.Instance }
func (c *Component) GetName() string                         { return c.Name }
func (c *Component) GetCode() string                         { return c.Code }
func (c *Component) GetDefaultNormalized() string            { return "unknown-component" }

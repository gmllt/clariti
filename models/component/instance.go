package component

import "github.com/gmllt/clariti/utils"

// Instance represents a service instance running on a platform
type Instance struct {
	BaseComponent
	Platform *Platform `json:"platform,omitempty"`
}

// NewInstance creates a new instance with the given name, code and platform
func NewInstance(name, code string, platform *Platform) *Instance {
	return &Instance{
		BaseComponent: BaseComponent{
			Name: name,
			Code: code,
		},
		Platform: platform,
	}
}

// String returns the string representation of the instance
func (i *Instance) String() string {
	return utils.StrUtils.BuildHierarchicalStringForComponent(i)
}

// Normalize returns a normalized identifier for the instance
func (i *Instance) Normalize() string {
	return utils.StrUtils.NormalizeHierarchicalComponent(i)
}

// HierarchicalStringable interface implementation
func (i *Instance) GetParent() utils.HierarchicalStringable { return i.Platform }
func (i *Instance) GetName() string                         { return i.Name }
func (i *Instance) GetCode() string                         { return i.Code }
func (i *Instance) GetDefaultNormalized() string            { return "unknown-instance" }

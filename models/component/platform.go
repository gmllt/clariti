package component

import "github.com/gmllt/clariti/utils"

// Platform represents the underlying infrastructure platform
type Platform struct {
	BaseComponent
}

// NewPlatform creates a new platform with the given name and code
func NewPlatform(name, code string) *Platform {
	return &Platform{
		BaseComponent: BaseComponent{
			Name: name,
			Code: code,
		},
	}
}

// String returns the name of the platform
func (p *Platform) String() string {
	if p.Name != "" {
		return p.Name
	}
	return "unknown platform"
}

// Normalize returns a normalized identifier for the platform
func (p *Platform) Normalize() string {
	return utils.StrUtils.NormalizeWithFallback(p.Code, p.String(), "unknown-platform")
}

// HierarchicalStringable interface implementation
func (p *Platform) GetParent() utils.HierarchicalStringable { return nil }
func (p *Platform) GetName() string                         { return p.Name }
func (p *Platform) GetCode() string                         { return p.Code }
func (p *Platform) GetDefaultNormalized() string            { return "unknown-platform" }

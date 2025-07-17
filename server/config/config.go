package config

import (
	"fmt"
	"os"

	"github.com/gmllt/clariti/models/component"
	"gopkg.in/yaml.v3"
)

// Config holds the server configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Auth       AuthConfig       `yaml:"auth"`
	Components ComponentsConfig `yaml:"components"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	CertFile string `yaml:"cert_file,omitempty"` // Optional TLS certificate file
	KeyFile  string `yaml:"key_file,omitempty"`  // Optional TLS private key file
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
}

// ComponentsConfig holds the hierarchical components configuration
type ComponentsConfig struct {
	Platforms []PlatformConfig `yaml:"platforms"`
}

// PlatformConfig represents a platform with its instances
type PlatformConfig struct {
	Name      string           `yaml:"name"`
	Code      string           `yaml:"code"`
	BaseURL   string           `yaml:"base_url"`
	Instances []InstanceConfig `yaml:"instances"`
}

// InstanceConfig represents an instance with its components
type InstanceConfig struct {
	Name       string            `yaml:"name"`
	Code       string            `yaml:"code"`
	Components []ComponentConfig `yaml:"components"`
}

// ComponentConfig represents a component configuration
type ComponentConfig struct {
	Name string `yaml:"name"`
	Code string `yaml:"code"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			// In a production environment, this should be logged properly
			_ = err // Explicitly ignore the error for now
		}
	}()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// Set defaults
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	return &config, nil
}

// GetAddress returns the full server address
func (c *Config) GetAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// GetAllPlatforms returns all platforms as component.Platform slice
func (c *Config) GetAllPlatforms() []component.Platform {
	platforms := make([]component.Platform, len(c.Components.Platforms))
	for i, p := range c.Components.Platforms {
		platforms[i] = component.Platform{
			BaseComponent: component.BaseComponent{
				Name: p.Name,
				Code: p.Code,
			},
		}
	}
	return platforms
}

// GetAllInstances returns all instances as component.Instance slice
func (c *Config) GetAllInstances() []component.Instance {
	var instances []component.Instance
	for _, platform := range c.Components.Platforms {
		platformModel := &component.Platform{
			BaseComponent: component.BaseComponent{
				Name: platform.Name,
				Code: platform.Code,
			},
		}
		for _, instance := range platform.Instances {
			inst := component.Instance{
				BaseComponent: component.BaseComponent{
					Name: instance.Name,
					Code: instance.Code,
				},
				Platform: platformModel,
			}
			instances = append(instances, inst)
		}
	}
	return instances
}

// GetAllComponents returns all components as component.Component slice
func (c *Config) GetAllComponents() []component.Component {
	var components []component.Component
	for _, platform := range c.Components.Platforms {
		platformModel := &component.Platform{
			BaseComponent: component.BaseComponent{
				Name: platform.Name,
				Code: platform.Code,
			},
		}
		for _, instance := range platform.Instances {
			inst := &component.Instance{
				BaseComponent: component.BaseComponent{
					Name: instance.Name,
					Code: instance.Code,
				},
				Platform: platformModel,
			}
			for _, comp := range instance.Components {
				newComp := component.Component{
					BaseComponent: component.BaseComponent{
						Name: comp.Name,
						Code: comp.Code,
					},
					Instance: inst,
				}
				components = append(components, newComp)
			}
		}
	}
	return components
}

// IsHTTPSEnabled returns true if both cert and key files are configured
func (c *Config) IsHTTPSEnabled() bool {
	return c.Server.CertFile != "" && c.Server.KeyFile != ""
}

// GetScheme returns "https" if HTTPS is enabled, "http" otherwise
func (c *Config) GetScheme() string {
	if c.IsHTTPSEnabled() {
		return "https"
	}
	return "http"
}

// GetFullURL returns the full URL for the server (http://host:port or https://host:port)
func (c *Config) GetFullURL() string {
	return fmt.Sprintf("%s://%s", c.GetScheme(), c.GetAddress())
}

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
	Storage    StorageConfig    `yaml:"storage"`
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

// StorageConfig holds the storage driver configuration
type StorageConfig struct {
	Driver string   `yaml:"driver"` // "ram" or "s3"
	S3     S3Config `yaml:"s3,omitempty"`
}

// S3Config holds S3-specific configuration
type S3Config struct {
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"access_key_id,omitempty"`     // Optional, can use IAM roles
	SecretAccessKey string `yaml:"secret_access_key,omitempty"` // Optional, can use IAM roles
	Endpoint        string `yaml:"endpoint,omitempty"`          // Optional, for S3-compatible services
	Prefix          string `yaml:"prefix,omitempty"`            // Optional, prefix for object keys
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
	if config.Storage.Driver == "" {
		config.Storage.Driver = "ram"
	}

	// Validate storage configuration
	if err := config.validateStorageConfig(); err != nil {
		return nil, fmt.Errorf("invalid storage configuration: %w", err)
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

// GetAllComponents returns all components as []*component.Component slice to avoid copying sync.RWMutex
func (c *Config) GetAllComponents() []*component.Component {
	var components []*component.Component
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
				newComp := &component.Component{
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

// validateStorageConfig validates the storage configuration based on the selected driver
func (c *Config) validateStorageConfig() error {
	switch c.Storage.Driver {
	case "ram":
		// RAM driver requires no additional configuration
		return nil
	case "s3":
		if c.Storage.S3.Region == "" {
			return fmt.Errorf("s3 region is required when using s3 driver")
		}
		if c.Storage.S3.Bucket == "" {
			return fmt.Errorf("s3 bucket is required when using s3 driver")
		}
		return nil
	default:
		return fmt.Errorf("unsupported storage driver: %s (supported: ram, s3)", c.Storage.Driver)
	}
}

// GetStorageDriver returns the configured storage driver name
func (c *Config) GetStorageDriver() string {
	return c.Storage.Driver
}

// GetS3Config returns the S3 configuration
func (c *Config) GetS3Config() S3Config {
	return c.Storage.S3
}

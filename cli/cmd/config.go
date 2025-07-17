package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Target represents a saved configuration profile
type Target struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config represents the CLI configuration
type Config struct {
	CurrentTarget string            `json:"current_target"`
	Targets       map[string]Target `json:"targets"`
}

// getConfigDir returns the configuration directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".clariti")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// loadConfig loads the configuration from file
func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Create default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Targets: make(map[string]Target),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Targets == nil {
		config.Targets = make(map[string]Target)
	}

	return &config, nil
}

// saveConfig saves the configuration to file
func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getCurrentTarget returns the current target configuration
func getCurrentTarget() (*Target, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	if config.CurrentTarget == "" {
		return nil, fmt.Errorf("no target set. Use 'clariti-cli target set' to configure a target")
	}

	target, exists := config.Targets[config.CurrentTarget]
	if !exists {
		return nil, fmt.Errorf("current target '%s' not found", config.CurrentTarget)
	}

	return &target, nil
}

// setTarget sets a new target or updates an existing one
func setTarget(name, url, username, password string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.Targets[name] = Target{
		Name:     name,
		URL:      url,
		Username: username,
		Password: password,
	}

	// Set as current target if none is set
	if config.CurrentTarget == "" {
		config.CurrentTarget = name
	}

	return saveConfig(config)
}

// switchTarget switches to an existing target
func switchTarget(name string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Targets[name]; !exists {
		return fmt.Errorf("target '%s' not found", name)
	}

	config.CurrentTarget = name
	return saveConfig(config)
}

// listTargets returns all configured targets
func listTargets() (map[string]Target, string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, "", err
	}

	return config.Targets, config.CurrentTarget, nil
}

// deleteTarget removes a target
func deleteTarget(name string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Targets[name]; !exists {
		return fmt.Errorf("target '%s' not found", name)
	}

	delete(config.Targets, name)

	// Clear current target if it was deleted
	if config.CurrentTarget == name {
		config.CurrentTarget = ""
	}

	return saveConfig(config)
}

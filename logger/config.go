package logger

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// Config représente la configuration du logger
type Config struct {
	Level   string `yaml:"level" json:"level"`       // debug, info, warn, error
	Format  string `yaml:"format" json:"format"`     // json, text
	NoColor bool   `yaml:"no_color" json:"no_color"` // désactive les couleurs pour le format text
}

// DefaultConfig retourne une configuration par défaut
func DefaultConfig() *Config {
	return &Config{
		Level:   "info",
		Format:  "text",
		NoColor: false,
	}
}

// ParseLevel convertit une chaîne de caractères en niveau logrus
func (c *Config) ParseLevel() logrus.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

// IsJSON retourne true si le format est JSON
func (c *Config) IsJSON() bool {
	return strings.ToLower(c.Format) == "json"
}

// Validate valide la configuration
func (c *Config) Validate() error {
	// Validation du niveau
	validLevels := []string{"debug", "info", "warn", "warning", "error", "fatal", "panic"}
	levelValid := false
	for _, level := range validLevels {
		if strings.ToLower(c.Level) == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		c.Level = "info" // Valeur par défaut
	}

	// Validation du format
	validFormats := []string{"json", "text"}
	formatValid := false
	for _, format := range validFormats {
		if strings.ToLower(c.Format) == format {
			formatValid = true
			break
		}
	}
	if !formatValid {
		c.Format = "text" // Valeur par défaut
	}

	return nil
}

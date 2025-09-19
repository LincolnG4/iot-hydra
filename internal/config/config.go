package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/LincolnG4/iot-hydra/internal/utils"
	"gopkg.in/yaml.v3"
)

type ConfigYAML struct {
	APIService     Service            `yaml:"apiService" validate:"required"`
	TelemetryAgent TelemetryAgentYAML `yaml:"telemetryAgent" validate:"required"`
}

// NewConfigFromYAML reads, unmarshals, and validates the YAML configuration file from a given path.
func NewConfigFromYAML(filePath string) (ConfigYAML, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return ConfigYAML{}, fmt.Errorf("failed to read config file '%s': %w", filePath, err)
	}

	var cfg ConfigYAML

	decoder := yaml.NewDecoder(bytes.NewReader(yamlFile))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return ConfigYAML{}, fmt.Errorf("failed to unmarshal config file '%s': %w", filePath, err)
	}

	// Validate the struct using our validator instance
	if err := utils.Validate.Struct(&cfg); err != nil {
		// Use the reusable function to format the error message
		formattedError := utils.FormatValidationErrors(err)
		return ConfigYAML{}, fmt.Errorf("the configuration file is invalid:\n%s", formattedError)
	}

	return cfg, nil
}

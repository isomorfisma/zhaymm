package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads YAML file and returns it as a pointer to Config
func LoadConfig(filename string) (*Config, error) {
	// I/O Operation
	data, err := os.ReadFile(filename)
	if err != nil {
		// %w for error wrapping
		return nil, fmt.Errorf("Failure to read file %s: %w", filename, err)
	}	

	// Empty Config struct
	var cfg Config

	// Translate YAML to cfg
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Invalid YAML file: %w", err)
	}

	// Returns pointer from the filled config
	return &cfg, nil
}
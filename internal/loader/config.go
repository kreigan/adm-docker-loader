// Package loader provides a Docker Compose stack loader for managing multiple stacks.
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the loader configuration.
type Config struct {
	CommonArgs []string `yaml:"common-args"`
	UpArgs     []string `yaml:"up-args"`
	DownArgs   []string `yaml:"down-args"`
	Timeout    int      `yaml:"timeout"`
}

// StackConfig represents per-stack configuration.
type StackConfig struct {
	UpArgs   []string `yaml:"up-args"`
	DownArgs []string `yaml:"down-args"`
}

// LoadConfig loads the global configuration from the base directory.
func LoadConfig(baseDir string) (*Config, error) {
	configPath := filepath.Join(baseDir, "config.yaml")

	// Default configuration
	config := &Config{
		CommonArgs: []string{},
		UpArgs: []string{
			"--detach",
			"--wait-timeout",
			"30",
			"--pull",
			"always",
		},
		DownArgs: []string{},
		Timeout:  10,
	}

	// Add --env-file only if .env exists
	envFile := filepath.Join(baseDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		config.CommonArgs = append(config.CommonArgs, "--env-file", envFile)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath) //nolint:gosec // Config path is constructed from trusted base directory
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Make paths in common_args absolute
	for i := 0; i < len(config.CommonArgs); i++ {
		if config.CommonArgs[i] == "--env-file" && i+1 < len(config.CommonArgs) {
			envFile := config.CommonArgs[i+1]
			if !filepath.IsAbs(envFile) {
				config.CommonArgs[i+1] = filepath.Join(baseDir, envFile)
			}
		}
	}

	return config, nil
}

// LoadStackConfig loads stack-specific configuration if it exists.
func LoadStackConfig(stackDir string) (*StackConfig, error) {
	configPath := filepath.Join(stackDir, "config.yaml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &StackConfig{}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath) //nolint:gosec // Stack config path is from trusted directory
	if err != nil {
		return nil, fmt.Errorf("reading stack config: %w", err)
	}

	// Parse YAML
	var stackConfig StackConfig
	if err := yaml.Unmarshal(data, &stackConfig); err != nil {
		return nil, fmt.Errorf("parsing stack config: %w", err)
	}

	return &stackConfig, nil
}

// MergeStackConfig merges stack-specific config with global config.
func (c *Config) MergeStackConfig(stackConfig *StackConfig) *Config {
	merged := &Config{
		CommonArgs: c.CommonArgs,
		UpArgs:     c.UpArgs,
		DownArgs:   c.DownArgs,
	}

	// Override up_args if provided
	if len(stackConfig.UpArgs) > 0 {
		merged.UpArgs = stackConfig.UpArgs
	}

	// Override down_args if provided
	if len(stackConfig.DownArgs) > 0 {
		merged.DownArgs = stackConfig.DownArgs
	}

	return merged
}

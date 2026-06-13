package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config corresponds to .trellis/config.yaml.
type Config struct {
	Packages  []string          `yaml:"packages"`
	Hooks     map[string]string `yaml:"hooks,omitempty"`
	Codex     CodexConfig       `yaml:"codex,omitempty"`
	Developer string            `yaml:"developer,omitempty"`
}

// CodexConfig holds Codex-specific settings.
type CodexConfig struct {
	DispatchMode string `yaml:"dispatch_mode"` // "inline" | "sub-agent"
}

// Load reads a Config from the given file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found at %s: %w", path, err)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}

// Save writes the Config to the given file path.
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// Validate checks the configuration for correctness.
func (c *Config) Validate() error {
	for i, pkg := range c.Packages {
		if pkg == "" {
			return fmt.Errorf("packages[%d] is empty", i)
		}
	}
	if c.Codex.DispatchMode != "" && c.Codex.DispatchMode != "inline" && c.Codex.DispatchMode != "sub-agent" {
		return fmt.Errorf("codex.dispatch_mode must be 'inline' or 'sub-agent', got %q", c.Codex.DispatchMode)
	}
	return nil
}

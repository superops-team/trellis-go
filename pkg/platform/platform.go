package platform

import "fmt"

// Class represents the workflow mode classification of an AI coding platform.
type Class string

const (
	ClassPushBased Class = "push" // Class-1: Claude, Cursor, etc.
	ClassPullBased Class = "pull" // Class-2: Codex, Copilot, etc.
	ClassAgentless Class = "none" // Class-3: Kilo, Windsurf, etc.
)

// Platform defines the configuration for an AI coding platform.
type Platform struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	ConfigDir           string   `json:"config_dir"`
	TemplateDirs        []string `json:"template_dirs"`
	AgentCapable        bool     `json:"agent_capable"`
	HasHooks            bool     `json:"has_hooks"`
	SupportsAgentSkills bool     `json:"supports_agent_skills"`
	CLIFlag             string   `json:"cli_flag"`
	Class               Class    `json:"class"`
	Aliases             []string `json:"aliases,omitempty"`
}

// Validate checks the platform configuration for completeness and correctness.
func (p *Platform) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("platform ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("platform name is required")
	}
	if p.ConfigDir == "" {
		return fmt.Errorf("platform %q: config_dir is required", p.ID)
	}
	if p.ConfigDir[0] == '/' {
		return fmt.Errorf("platform %q: config_dir must not start with '/': %s", p.ID, p.ConfigDir)
	}
	if p.Class != ClassPushBased && p.Class != ClassPullBased && p.Class != ClassAgentless {
		return fmt.Errorf("platform %q: invalid class %q", p.ID, p.Class)
	}
	return nil
}

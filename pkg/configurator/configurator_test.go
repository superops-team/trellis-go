package configurator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/superops-team/trellis-go/pkg/platform"
)

func TestFor_ReturnsCorrectConfiguratorByClass(t *testing.T) {
	tests := []struct {
		name     string
		platform platform.Platform
		check    func(Configurator) bool
	}{
		{"push-based", platform.Platform{ID: "claude", Class: platform.ClassPushBased}, func(c Configurator) bool { _, ok := c.(*pushConfigurator); return ok }},
		{"pull-based", platform.Platform{ID: "codex", Class: platform.ClassPullBased}, func(c Configurator) bool { _, ok := c.(*pullConfigurator); return ok }},
		{"agentless", platform.Platform{ID: "kilo", Class: platform.ClassAgentless}, func(c Configurator) bool { _, ok := c.(*agentlessConfigurator); return ok }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := For(tt.platform, "trellis")
			if c == nil {
				t.Fatal("expected non-nil configurator")
			}
			if c.Name() != tt.platform.ID {
				t.Errorf("expected name %q, got %q", tt.platform.ID, c.Name())
			}
			if !tt.check(c) {
				t.Errorf("unexpected configurator type: %T", c)
			}
		})
	}
}

func TestFor_UnknownClassReturnsNil(t *testing.T) {
	p := platform.Platform{ID: "unknown", Class: platform.Class("unknown")}
	c := For(p, "trellis")
	if c != nil {
		t.Errorf("expected nil for unknown class, got %T", c)
	}
}

func TestPushConfigurator_Generate(t *testing.T) {
	root := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	c := For(p, "trellis")

	if err := c.Generate(root, Options{}); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check hooks
	for _, hook := range []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"} {
		path := filepath.Join(root, ".claude", hook)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing hook %s: %v", hook, err)
		}
	}

	// Check agents
	for _, agent := range []string{"trellis-implement.md", "trellis-check.md", "trellis-research.md"} {
		path := filepath.Join(root, ".claude/agents", agent)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing agent %s: %v", agent, err)
		}
	}

	// Check skills
	for _, skill := range []string{"trellis-brainstorm", "trellis-before-dev", "trellis-check", "trellis-update-spec", "trellis-break-loop"} {
		path := filepath.Join(root, ".claude/skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing skill %s: %v", skill, err)
		}
	}

	// Check shared skills
	for _, skill := range []string{"trellis-brainstorm", "trellis-before-dev", "trellis-check", "trellis-update-spec", "trellis-break-loop"} {
		path := filepath.Join(root, ".agents/skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing shared skill %s: %v", skill, err)
		}
	}

	// Check commands
	for _, cmd := range []string{"trellis/start.md", "trellis/finish-work.md", "trellis/continue.md"} {
		path := filepath.Join(root, ".claude/commands", cmd)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing command %s: %v", cmd, err)
		}
	}
}

func TestPullConfigurator_Generate(t *testing.T) {
	root := t.TempDir()
	p := platform.Platform{ID: "codex", Name: "Codex", ConfigDir: ".codex", Class: platform.ClassPullBased}
	c := For(p, "trellis")

	if err := c.Generate(root, Options{}); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check agents (TOML format)
	for _, agent := range []string{"trellis-implement.toml", "trellis-check.toml", "trellis-research.toml"} {
		path := filepath.Join(root, ".codex/agents", agent)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing agent %s: %v", agent, err)
		}
	}

	// Check no hooks
	for _, hook := range []string{"session-start.sh", "inject-context.sh"} {
		path := filepath.Join(root, ".codex", hook)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("pull-based should not have hook %s", hook)
		}
	}

	// Check AGENTS.md
	agentsPath := filepath.Join(root, ".codex", "AGENTS.md")
	if _, err := os.Stat(agentsPath); err != nil {
		t.Errorf("missing AGENTS.md: %v", err)
	}
}

func TestAgentlessConfigurator_Generate(t *testing.T) {
	root := t.TempDir()
	p := platform.Platform{ID: "kilo", Name: "Kilo", ConfigDir: ".kilocode", Class: platform.ClassAgentless}
	c := For(p, "trellis")

	if err := c.Generate(root, Options{}); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check skills
	for _, skill := range []string{"trellis-brainstorm", "trellis-before-dev", "trellis-check", "trellis-update-spec", "trellis-break-loop"} {
		path := filepath.Join(root, ".kilocode/skills", skill, "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing skill %s: %v", skill, err)
		}
	}

	// Check no agents
	agentsDir := filepath.Join(root, ".kilocode/agents")
	if _, err := os.Stat(agentsDir); err == nil {
		t.Error("agentless should not have agents directory")
	}

	// Check no hooks
	for _, hook := range []string{"session-start.sh", "inject-context.sh"} {
		path := filepath.Join(root, ".kilocode", hook)
		if _, err := os.Stat(path); err == nil {
			t.Errorf("agentless should not have hook %s", hook)
		}
	}
}

func TestPushConfigurator_DryRun(t *testing.T) {
	root := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	c := For(p, "trellis")

	if err := c.Generate(root, Options{DryRun: true}); err != nil {
		t.Fatalf("Generate with DryRun failed: %v", err)
	}

	// Nothing should be written
	claudeDir := filepath.Join(root, ".claude")
	if _, err := os.Stat(claudeDir); err == nil {
		t.Error("DryRun should not create any files")
	}
}

func TestPushConfigurator_Remove(t *testing.T) {
	root := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	c := For(p, "trellis")

	// Generate first
	if err := c.Generate(root, Options{}); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Then remove
	if err := c.Remove(root); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	// Verify removed
	claudeDir := filepath.Join(root, ".claude")
	if _, err := os.Stat(claudeDir); err == nil {
		t.Error("Remove should delete platform directory")
	}
}

func TestAllPlatformsCanGenerate(t *testing.T) {
	registry := platform.NewRegistry()
	for _, p := range registry.All() {
		t.Run(p.ID, func(t *testing.T) {
			root := t.TempDir()
			c := For(p, "trellis")
			if c == nil {
				t.Fatalf("no configurator for platform %s (class %s)", p.ID, p.Class)
			}
			if err := c.Generate(root, Options{}); err != nil {
				t.Fatalf("Generate failed for %s: %v", p.ID, err)
			}
		})
	}
}

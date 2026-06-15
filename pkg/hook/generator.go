package hook

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/platform"
)

// Generator creates hook scripts and configuration files for a platform.
type Generator struct {
	Platform platform.Platform
	Binary   string
}

// NewGenerator creates a hook generator.
func NewGenerator(p platform.Platform, binary string) *Generator {
	return &Generator{Platform: p, Binary: binary}
}

// GenerateAll generates all hook files for the platform.
func (g *Generator) GenerateAll(dst string) error {
	switch g.Platform.Class {
	case platform.ClassPushBased:
		return g.generatePushBased(dst)
	case platform.ClassPullBased:
		return g.generatePullBased(dst)
	case platform.ClassAgentless:
		return g.generateAgentless(dst)
	}
	return fmt.Errorf("unknown platform class: %s", g.Platform.Class)
}

func (g *Generator) generatePushBased(dst string) error {
	if err := g.GenerateSessionStart(dst); err != nil {
		return err
	}
	if err := g.GenerateInjectContext(dst); err != nil {
		return err
	}
	return g.GenerateInjectWorkflowState(dst)
}

func (g *Generator) generatePullBased(dst string) error {
	// Generate agent definition file (e.g., .toml for Codex)
	return g.GenerateAgentDef(dst)
}

func (g *Generator) generateAgentless(dst string) error {
	// Generate before-dev skill markdown
	return g.GenerateBeforeDevSkill(dst)
}

// GenerateAgentDef creates a sub-agent definition file for pull-based platforms.
func (g *Generator) GenerateAgentDef(dst string) error {
	if g.Platform.Class != platform.ClassPullBased {
		return fmt.Errorf("agent defs only for pull-based platforms")
	}

	content := fmt.Sprintf(`name = "trellis-implement"
description = "Trellis implementation sub-agent"

[features]
multi_agent = false

[features.multi_agent_v2]
enabled = false

developer_instructions = """
You are the trellis-implement sub-agent.
Do NOT spawn another sub-agent.

Active task context will be provided by the parent session.
"""
`)

	path := filepath.Join(dst, "trellis-implement.toml")
	return fsutil.WriteFile(path, []byte(content), 0644)
}

// GenerateBeforeDevSkill creates a before-dev skill for agentless platforms.
func (g *Generator) GenerateBeforeDevSkill(dst string) error {
	content := `# Trellis Before-Dev Skill

Load project specs from .trellis/spec/ before starting development.
`
	path := filepath.Join(dst, "trellis-before-dev.md")
	return fsutil.WriteFile(path, []byte(content), 0644)
}

// GenerateSessionStart creates a session start hook script.
func (g *Generator) GenerateSessionStart(dst string) error {
	return g.generateHookScript(dst, "session-start.sh", "session start", "session-start")
}

// GenerateInjectContext creates a context injection hook script.
func (g *Generator) GenerateInjectContext(dst string) error {
	return g.generateHookScript(dst, "inject-context.sh", "context injection", "inject-context")
}

// GenerateInjectWorkflowState creates a workflow state injection hook script.
func (g *Generator) GenerateInjectWorkflowState(dst string) error {
	return g.generateHookScript(dst, "inject-workflow-state.sh", "workflow state injection", "inject-workflow-state")
}

func (g *Generator) generateHookScript(dst, fileName, description, subcommand string) error {
	script := fmt.Sprintf(`#!/bin/sh
# Trellis %s hook for %s
exec %s hook %s "$@"
`, description, g.Platform.Name, shellQuote(g.Binary), subcommand)
	path := filepath.Join(dst, fileName)
	return fsutil.WriteFile(path, []byte(script), 0755)
}

func shellQuote(word string) string {
	if word == "" {
		return "''"
	}
	if !strings.ContainsAny(word, " \t\n'\"\\$`!*?[]{}();<>|&") {
		return word
	}
	return "'" + strings.ReplaceAll(word, "'", "'\"'\"'") + "'"
}

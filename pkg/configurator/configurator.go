package configurator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/superops-team/trellis-go/pkg/agent"
	"github.com/superops-team/trellis-go/pkg/command"
	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/hook"
	"github.com/superops-team/trellis-go/pkg/platform"
	"github.com/superops-team/trellis-go/pkg/skill"
)

// Options controls configurator behavior.
type Options struct {
	DryRun bool
	Force  bool
	Binary string // trellis binary path for hook scripts
}

// Configurator generates platform-specific Trellis configuration files.
type Configurator interface {
	// Name returns the platform ID.
	Name() string

	// Generate creates all configuration files for the platform in the project root.
	Generate(projectRoot string, opts Options) error

	// Remove cleans up all configuration files for the platform.
	Remove(projectRoot string) error
}

// For returns the appropriate configurator for the given platform.
func For(p platform.Platform, binary string) Configurator {
	switch p.Class {
	case platform.ClassPushBased:
		return &pushConfigurator{platform: p, binary: binary}
	case platform.ClassPullBased:
		return &pullConfigurator{platform: p, binary: binary}
	case platform.ClassAgentless:
		return &agentlessConfigurator{platform: p, binary: binary}
	}
	return nil
}

// writeFile writes content to a file, respecting DryRun and Force options.
func writeFile(path string, content []byte, perm os.FileMode, opts Options) error {
	if opts.DryRun {
		fmt.Printf("[dry-run] write %s (%d bytes)\n", path, len(content))
		return nil
	}
	return writeFileActual(path, content, perm, opts)
}

func writeFileActual(path string, content []byte, perm os.FileMode, opts Options) error {
	dir := filepath.Dir(path)
	if err := fsutil.EnsureDir(dir); err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil && !opts.Force {
		fmt.Printf("[skip] %s already exists (use --force to overwrite)\n", path)
		return nil
	}
	return os.WriteFile(path, content, perm)
}

// writeAgentFiles generates agent definition files for a platform.
func writeAgentFiles(projectRoot, platformDir string, p platform.Platform, opts Options) error {
	agents := agent.AllAgents()
	for _, a := range agents {
		var fn, content string
		switch p.ID {
		case "claude", "codebuddy", "droid":
			fn, content = agent.FormatForClaudeCode(a)
		case "cursor":
			fn, content = agent.FormatForCursor(a)
		case "opencode":
			fn, content = agent.FormatForOpenCode(a)
		case "codex":
			fn, content = agent.FormatForCodex(a)
		case "kiro":
			fn, content = agent.FormatForKiro(a)
		case "gemini":
			fn, content = agent.FormatForGemini(a)
		case "qoder":
			fn, content = agent.FormatForQoder(a)
		case "copilot":
			fn, content = agent.FormatForCopilot(a)
		case "pi":
			fn, content = agent.FormatForPi(a)
		default:
			fn, content = agent.FormatForClaudeCode(a)
		}
		path := filepath.Join(projectRoot, platformDir, fn)
		if err := writeFile(path, []byte(content), 0644, opts); err != nil {
			return fmt.Errorf("write agent %s: %w", a.Name, err)
		}
	}
	return nil
}

// writeSkillFiles generates skill definition files for a platform.
func writeSkillFiles(projectRoot, platformDir string, p platform.Platform, opts Options) error {
	skills := skill.AllSkills()
	for _, s := range skills {
		var fn, content string
		switch p.ID {
		case "cursor":
			fn, content = skill.FormatForCursor(s)
		case "gemini":
			fn, content = skill.FormatForGemini(s)
		case "qoder":
			fn, content = skill.FormatForQoder(s)
		default:
			fn, content = skill.FormatForShared(s)
		}
		path := filepath.Join(projectRoot, platformDir, fn)
		if err := writeFile(path, []byte(content), 0644, opts); err != nil {
			return fmt.Errorf("write skill %s: %w", s.Name, err)
		}
	}

	// Also write to shared .agents/skills/ directory
	sharedDir := filepath.Join(projectRoot, ".agents", "skills")
	for _, s := range skills {
		fn, content := skill.FormatForShared(s)
		path := filepath.Join(sharedDir, fn)
		if err := writeFile(path, []byte(content), 0644, opts); err != nil {
			return fmt.Errorf("write shared skill %s: %w", s.Name, err)
		}
	}

	return nil
}

// writeCommandFiles generates command definition files for a platform.
func writeCommandFiles(projectRoot, platformDir string, p platform.Platform, opts Options) error {
	cmds := command.AllCommands()
	for _, c := range cmds {
		var fn, content string
		switch p.ID {
		case "claude", "opencode", "codebuddy", "droid":
			fn, content = command.FormatForClaudeCode(c)
		case "cursor":
			fn, content = command.FormatForCursor(c)
		case "codex":
			fn, content = command.FormatForCodex(c)
		case "gemini":
			fn, content = command.FormatForGemini(c)
		case "qoder":
			fn, content = command.FormatForQoder(c)
		case "copilot":
			fn, content = command.FormatForCopilot(c)
		case "pi":
			fn, content = command.FormatForCodex(c)
		default:
			fn, content = command.FormatForClaudeCode(c)
		}
		path := filepath.Join(projectRoot, platformDir, fn)
		if err := writeFile(path, []byte(content), 0644, opts); err != nil {
			return fmt.Errorf("write command %s: %w", c.Name, err)
		}
	}
	return nil
}

// writeHookFiles generates hook script files for push-based platforms.
func writeHookFiles(dst string, p platform.Platform, binary string, opts Options) error {
	g := hook.NewGenerator(p, binary)
	if opts.DryRun {
		fmt.Printf("[dry-run] generate hooks for %s in %s\n", p.ID, dst)
		return nil
	}
	return g.GenerateAll(dst)
}

// agentDir returns the agent directory for a platform.
func agentDir(p platform.Platform) string {
	switch p.ID {
	case "claude", "cursor", "opencode", "codebuddy":
		return filepath.Join(p.ConfigDir, "agents")
	case "codex":
		return filepath.Join(p.ConfigDir, "agents")
	case "kiro":
		return filepath.Join(p.ConfigDir, "agents")
	case "gemini":
		return filepath.Join(p.ConfigDir, "agents")
	case "qoder":
		return filepath.Join(p.ConfigDir, "agents")
	case "copilot":
		return filepath.Join(p.ConfigDir, "agents")
	case "droid":
		return filepath.Join(p.ConfigDir, "droids")
	case "pi":
		return filepath.Join(p.ConfigDir, "agents")
	default:
		return filepath.Join(p.ConfigDir, "agents")
	}
}

// skillDir returns the skill directory for a platform.
func skillDir(p platform.Platform) string {
	return filepath.Join(p.ConfigDir, "skills")
}

// commandDir returns the command directory for a platform.
func commandDir(p platform.Platform) string {
	switch p.ID {
	case "claude", "opencode", "codebuddy", "droid":
		return filepath.Join(p.ConfigDir, "commands")
	case "cursor":
		return filepath.Join(p.ConfigDir, "commands")
	case "codex":
		return filepath.Join(p.ConfigDir, "prompts")
	case "gemini":
		return filepath.Join(p.ConfigDir, "commands")
	case "qoder":
		return filepath.Join(p.ConfigDir, "commands")
	case "copilot":
		return filepath.Join(p.ConfigDir, "prompts")
	case "pi":
		return filepath.Join(p.ConfigDir, "prompts")
	default:
		return filepath.Join(p.ConfigDir, "commands")
	}
}

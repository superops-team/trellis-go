package configurator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/platform"
	"github.com/superops-team/trellis-go/pkg/skill"
)

// pullConfigurator handles Class-2 (pull-based) platforms.
// Pull-based platforms don't have hooks; they use agent definitions + skills + commands.
type pullConfigurator struct {
	platform platform.Platform
	binary   string
}

func (c *pullConfigurator) Name() string { return c.platform.ID }

func (c *pullConfigurator) Generate(projectRoot string, opts Options) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := fsutil.EnsureDir(platformDir); err != nil {
		return err
	}

	// 1. Agents (with pull-based prelude)
	ad := agentDir(c.platform)
	if err := writeAgentFiles(projectRoot, ad, c.platform, opts); err != nil {
		return fmt.Errorf("agents: %w", err)
	}

	// 2. Skills
	sd := skillDir(c.platform)
	if err := writeSkillFiles(projectRoot, sd, c.platform, opts); err != nil {
		return fmt.Errorf("skills: %w", err)
	}

	// 3. Commands
	cd := commandDir(c.platform)
	if err := writeCommandFiles(projectRoot, cd, c.platform, opts); err != nil {
		return fmt.Errorf("commands: %w", err)
	}

	// 4. For Codex: write AGENTS.md entry
	if c.platform.ID == "codex" {
		skills := skill.AllSkills()
		entry := skill.CodexAgentEntry(skills)
		entryPath := filepath.Join(platformDir, "AGENTS.md")
		if err := writeFile(entryPath, []byte(entry), 0644, opts); err != nil {
			return fmt.Errorf("codex AGENTS.md: %w", err)
		}
	}

	return nil
}

func (c *pullConfigurator) Remove(projectRoot string) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := os.RemoveAll(platformDir); err != nil {
		return fmt.Errorf("remove %s: %w", c.platform.ID, err)
	}
	return nil
}

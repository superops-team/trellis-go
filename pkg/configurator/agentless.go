package configurator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/platform"
)

// agentlessConfigurator handles Class-3 (agentless) platforms.
// Agentless platforms only get skills and workflow commands — no agents or hooks.
type agentlessConfigurator struct {
	platform platform.Platform
	binary   string
}

func (c *agentlessConfigurator) Name() string { return c.platform.ID }

func (c *agentlessConfigurator) Generate(projectRoot string, opts Options) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := fsutil.EnsureDir(platformDir); err != nil {
		return err
	}

	// 1. Skills (no agents, no hooks)
	sd := skillDir(c.platform)
	if err := writeSkillFiles(projectRoot, sd, c.platform, opts); err != nil {
		return fmt.Errorf("skills: %w", err)
	}

	// 2. Workflow commands (as markdown workflow files)
	cd := commandDir(c.platform)
	if err := writeCommandFiles(projectRoot, cd, c.platform, opts); err != nil {
		return fmt.Errorf("commands: %w", err)
	}

	return nil
}

func (c *agentlessConfigurator) Remove(projectRoot string) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := os.RemoveAll(platformDir); err != nil {
		return fmt.Errorf("remove %s: %w", c.platform.ID, err)
	}
	return nil
}

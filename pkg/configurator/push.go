package configurator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/platform"
)

// pushConfigurator handles Class-1 (push-based) platforms.
type pushConfigurator struct {
	platform platform.Platform
	binary   string
}

func (c *pushConfigurator) Name() string { return c.platform.ID }

func (c *pushConfigurator) Generate(projectRoot string, opts Options) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if !opts.DryRun {
		if err := fsutil.EnsureDir(platformDir); err != nil {
			return err
		}
	}

	// 1. Hook scripts
	hookDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := writeHookFiles(hookDir, c.platform, c.binary, opts); err != nil {
		return fmt.Errorf("hooks: %w", err)
	}

	// 2. Agents
	ad := agentDir(c.platform)
	agentPath := filepath.Join(projectRoot, ad)
	if !opts.DryRun {
		if err := fsutil.EnsureDir(agentPath); err != nil {
			return err
		}
	}
	if err := writeAgentFiles(projectRoot, ad, c.platform, opts); err != nil {
		return fmt.Errorf("agents: %w", err)
	}

	// 3. Skills
	sd := skillDir(c.platform)
	if !opts.DryRun {
		if err := writeSkillFiles(projectRoot, sd, c.platform, opts); err != nil {
			return fmt.Errorf("skills: %w", err)
		}
	} else {
		writeSkillFiles(projectRoot, sd, c.platform, opts)
	}

	// 4. Commands
	cd := commandDir(c.platform)
	if !opts.DryRun {
		if err := writeCommandFiles(projectRoot, cd, c.platform, opts); err != nil {
			return fmt.Errorf("commands: %w", err)
		}
	} else {
		writeCommandFiles(projectRoot, cd, c.platform, opts)
	}

	return nil
}

func (c *pushConfigurator) Remove(projectRoot string) error {
	platformDir := filepath.Join(projectRoot, c.platform.ConfigDir)
	if err := os.RemoveAll(platformDir); err != nil {
		return fmt.Errorf("remove %s: %w", c.platform.ID, err)
	}
	return nil
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mindfold/trellis/pkg/config"
	"github.com/mindfold/trellis/pkg/fsutil"
	"github.com/mindfold/trellis/pkg/platform"
	"github.com/spf13/cobra"
)

var initOpts struct {
	developer string
	platforms []string
	all       bool
}

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Trellis in the current repository",
		RunE:  runInit,
	}
	cmd.Flags().StringVarP(&initOpts.developer, "developer", "u", "", "Developer name")
	cmd.Flags().StringArrayVarP(&initOpts.platforms, "platform", "p", nil, "Platforms to configure (can be specified multiple times)")
	cmd.Flags().BoolVar(&initOpts.all, "all", false, "Configure all supported platforms")
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Verify we're in a Git repository
	gitDir := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("not a git repository; run 'git init' first")
	}

	trellisDir := filepath.Join(cwd, ".trellis")
	if _, err := os.Stat(trellisDir); err == nil {
		return fmt.Errorf("trellis already initialized at %s", trellisDir)
	}

	// Resolve developer name
	developer := initOpts.developer
	if developer == "" {
		developer = os.Getenv("USER")
		if developer == "" {
			developer = "developer"
		}
	}

	// Resolve platforms
	registry := platform.NewRegistry()
	var platforms []platform.Platform
	if initOpts.all {
		platforms = registry.All()
	} else if len(initOpts.platforms) > 0 {
		for _, id := range initOpts.platforms {
			p, ok := registry.Get(id)
			if !ok {
				valid := strings.Join(registry.IDs(), ", ")
				return fmt.Errorf("unknown platform %q; valid platforms: %s", id, valid)
			}
			platforms = append(platforms, p)
		}
	} else {
		// Default: Claude only
		p, _ := registry.Get("claude")
		platforms = []platform.Platform{p}
	}

	// Create .trellis directory structure
	dirs := []string{
		trellisDir,
		filepath.Join(trellisDir, "spec"),
		filepath.Join(trellisDir, "tasks"),
		filepath.Join(trellisDir, "workspace"),
		filepath.Join(trellisDir, ".runtime", "sessions"),
		filepath.Join(trellisDir, "scripts"),
	}
	for _, d := range dirs {
		if err := fsutil.EnsureDir(d); err != nil {
			return err
		}
	}

	// Write config.yaml
	cfg := &config.Config{
		Packages:  []string{},
		Developer: developer,
		Codex:     config.CodexConfig{DispatchMode: "inline"},
	}
	cfgPath := filepath.Join(trellisDir, "config.yaml")
	if err := cfg.Save(cfgPath); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	// Write .version
	versionPath := filepath.Join(trellisDir, ".version")
	if err := os.WriteFile(versionPath, []byte(version+"\n"), 0644); err != nil {
		return fmt.Errorf("write version: %w", err)
	}

	// Write workflow.md
	workflowPath := filepath.Join(trellisDir, "workflow.md")
	workflowContent := `# Trellis Workflow

## Phase 1: Plan
[workflow-state:PLAN]
Brainstorm requirements and write PRD.

## Phase 2: Implement
[workflow-state:IMPLEMENT]
Write code from the PRD.

## Phase 3: Verify
[workflow-state:VERIFY]
Review code against specs and run checks.

## Phase 4: Finish
[workflow-state:FINISH]
Archive the task and update journals.
`
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		return fmt.Errorf("write workflow: %w", err)
	}

	// Create platform configs (placeholder)
	for _, p := range platforms {
		platformDir := filepath.Join(cwd, p.ConfigDir)
		if err := fsutil.EnsureDir(platformDir); err != nil {
			return err
		}
		// TODO: render actual platform templates
		_ = platformDir
	}

	if verbose {
		fmt.Printf("Initialized Trellis for %s\n", developer)
		fmt.Printf("Platforms: ")
		for i, p := range platforms {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(p.Name)
		}
		fmt.Println()
	}
	return nil
}

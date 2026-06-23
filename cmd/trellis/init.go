package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/superops-team/trellis-go/pkg/configurator"
	"github.com/superops-team/trellis-go/pkg/config"
	"github.com/superops-team/trellis-go/pkg/fsutil"
	"github.com/superops-team/trellis-go/pkg/platform"
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
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	// Verify we're in a Git repository
	gitDir := filepath.Join(paths.RepoRoot, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return fmt.Errorf("not a git repository; run 'git init' first")
	}

	trellisDir := paths.TrellisDir
	if _, err := os.Stat(trellisDir); err == nil {
		// .trellis/ already exists — multi-user mode
		return runInitExisting(paths)
	}

	return runInitFresh(paths)
}

// runInitExisting handles re-initialization of an already-initialized project.
func runInitExisting(paths resolvedPaths) error {
	hasPlatforms := len(initOpts.platforms) > 0 || initOpts.all
	hasDeveloper := initOpts.developer != ""

	// Non-interactive mode: --platform and/or --developer specified
	if hasPlatforms || hasDeveloper {
		if hasDeveloper {
			if err := writeDeveloperIdentity(paths.TrellisDir, initOpts.developer); err != nil {
				return err
			}
			fmt.Printf("Developer identity set: %s\n", initOpts.developer)
		}
		if hasPlatforms {
			if err := addPlatforms(paths); err != nil {
				return err
			}
		}
		return nil
	}

	// Interactive mode: show 3 options
	fmt.Println("Trellis is already initialized. What would you like to do?")
	fmt.Println("  1) Add AI platform(s)")
	fmt.Println("  2) Set up developer identity")
	fmt.Println("  3) Full re-initialize")
	fmt.Print("Choose (1-3): ")

	var choice string
	fmt.Scanln(&choice)

	switch strings.TrimSpace(choice) {
	case "1":
		fmt.Print("Platform(s) to add (comma-separated): ")
		var input string
		fmt.Scanln(&input)
		for _, id := range strings.Split(input, ",") {
			id = strings.TrimSpace(id)
			if id != "" {
				initOpts.platforms = append(initOpts.platforms, id)
			}
		}
		if len(initOpts.platforms) == 0 {
			return fmt.Errorf("no platforms specified")
		}
		return addPlatforms(paths)
	case "2":
		fmt.Print("Developer name: ")
		var name string
		fmt.Scanln(&name)
		if name == "" {
			return fmt.Errorf("developer name is required")
		}
		if err := writeDeveloperIdentity(paths.TrellisDir, name); err != nil {
			return err
		}
		fmt.Printf("Developer identity set: %s\n", name)
		return nil
	case "3":
		fmt.Print("This will reset all Trellis configuration. Continue? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			return fmt.Errorf("aborted")
		}
		// Remove existing .trellis/ and re-initialize
		if err := os.RemoveAll(paths.TrellisDir); err != nil {
			return fmt.Errorf("remove existing .trellis: %w", err)
		}
		return runInitFresh(paths)
	default:
		return fmt.Errorf("invalid choice: %s", choice)
	}
}

// runInitFresh performs a fresh initialization.
func runInitFresh(paths resolvedPaths) error {
	trellisDir := paths.TrellisDir

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
				// Try flag/alias lookup (e.g. --windsurf → devin)
				p, ok = registry.ForFlag("--" + id)
			}
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

	// Write .developer
	if err := writeDeveloperIdentity(trellisDir, developer); err != nil {
		return err
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

	// Create platform configs
	for _, p := range platforms {
		cfg := configurator.For(p, os.Args[0])
		if cfg == nil {
			return fmt.Errorf("no configurator for platform %s", p.ID)
		}
		if err := cfg.Generate(paths.RepoRoot, configurator.Options{}); err != nil {
			return fmt.Errorf("configure platform %s: %w", p.ID, err)
		}
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

// writeDeveloperIdentity writes .trellis/.developer file.
func writeDeveloperIdentity(trellisDir, name string) error {
	devPath := filepath.Join(trellisDir, ".developer")
	content := fmt.Sprintf("name: %q\n", name)
	return os.WriteFile(devPath, []byte(content), 0644)
}

// readDeveloperIdentity reads the developer name from .trellis/.developer.
func readDeveloperIdentity(trellisDir string) (string, error) {
	devPath := filepath.Join(trellisDir, ".developer")
	data, err := os.ReadFile(devPath)
	if err != nil {
		return "", err
	}
	// Parse "name: \"value\"" format
	content := strings.TrimSpace(string(data))
	if strings.HasPrefix(content, "name:") {
		name := strings.TrimPrefix(content, "name:")
		name = strings.TrimSpace(name)
		name = strings.Trim(name, "\"")
		return name, nil
	}
	return "", fmt.Errorf("invalid .developer format")
}

// addPlatforms adds new platforms to an existing project without overwriting.
func addPlatforms(paths resolvedPaths) error {
	registry := platform.NewRegistry()
	var platforms []platform.Platform
	if initOpts.all {
		platforms = registry.All()
	} else {
		for _, id := range initOpts.platforms {
			p, ok := registry.Get(id)
			if !ok {
				valid := strings.Join(registry.IDs(), ", ")
				return fmt.Errorf("unknown platform %q; valid platforms: %s", id, valid)
			}
			platforms = append(platforms, p)
		}
	}

	for _, p := range platforms {
		cfg := configurator.For(p, os.Args[0])
		if cfg == nil {
			return fmt.Errorf("no configurator for platform %s", p.ID)
		}
		if err := cfg.Generate(paths.RepoRoot, configurator.Options{}); err != nil {
			return fmt.Errorf("configure platform %s: %w", p.ID, err)
		}
		fmt.Printf("Added platform: %s\n", p.Name)
	}
	return nil
}

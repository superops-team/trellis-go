package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/superops-team/trellis-go/pkg/upgrade"
)

var upgradeOpts struct {
	tag    string
	dryRun bool
}

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Trellis to the latest version",
		RunE:  runUpgrade,
	}
	cmd.Flags().StringVar(&upgradeOpts.tag, "tag", "latest", "Version tag (latest, beta, or specific tag)")
	cmd.Flags().BoolVar(&upgradeOpts.dryRun, "dry-run", false, "Check for updates without downloading")
	return cmd
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	checker := &upgrade.Checker{
		CurrentVersion: version,
		RepoOwner:      "superops-team",
		RepoName:       "trellis-go",
	}

	release, err := checker.CheckLatest(upgradeOpts.tag)
	if err != nil {
		return fmt.Errorf("check latest version: %w", err)
	}

	if !checker.NeedsUpdate(release) {
		fmt.Printf("Already up to date (v%s)\n", version)
		return nil
	}

	fmt.Printf("New version available: %s (current: v%s)\n", release.TagName, version)

	if upgradeOpts.dryRun {
		fmt.Println("(Dry run — no changes made)")
		return nil
	}

	// Download to temp file
	tmpDir, err := os.MkdirTemp("", "trellis-upgrade")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	downloaded := filepath.Join(tmpDir, "trellis-new")
	if err := checker.Download(release, downloaded); err != nil {
		return fmt.Errorf("download update: %w", err)
	}

	// Get current binary path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable path: %w", err)
	}

	// Replace binary
	if err := checker.Replace(downloaded, execPath); err != nil {
		return fmt.Errorf("replace binary: %w", err)
	}

	fmt.Printf("Upgraded to %s\n", release.TagName)
	fmt.Println("Run 'trellis update' to sync templates")
	return nil
}

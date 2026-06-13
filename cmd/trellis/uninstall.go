package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uninstallOpts struct {
	keepTasks bool
}

func newUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove Trellis from the current repository",
		RunE:  runUninstall,
	}
	cmd.Flags().BoolVar(&uninstallOpts.keepTasks, "keep-tasks", false, "Preserve tasks directory")
	return cmd
}

func runUninstall(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	trellisDir := filepath.Join(cwd, ".trellis")
	if _, err := os.Stat(trellisDir); os.IsNotExist(err) {
		return fmt.Errorf("trellis is not initialized in this repository")
	}

	if !uninstallOpts.keepTasks {
		if err := os.RemoveAll(trellisDir); err != nil {
			return fmt.Errorf("remove .trellis: %w", err)
		}
	} else {
		// Remove everything except tasks/
		entries, err := os.ReadDir(trellisDir)
		if err != nil {
			return fmt.Errorf("read .trellis: %w", err)
		}
		for _, entry := range entries {
			if entry.Name() == "tasks" {
				continue
			}
			path := filepath.Join(trellisDir, entry.Name())
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("remove %s: %w", path, err)
			}
		}
	}

	if verbose {
		fmt.Println("Trellis has been uninstalled")
	}
	return nil
}

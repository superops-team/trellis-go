package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/superops-team/trellis-go/internal/embed"
	"github.com/superops-team/trellis-go/pkg/update"
)

var updateOpts struct {
	dryRun  bool
	migrate bool
}

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Trellis templates and configuration",
		RunE:  runUpdate,
	}
	cmd.Flags().BoolVar(&updateOpts.dryRun, "dry-run", false, "Preview changes without writing")
	cmd.Flags().BoolVar(&updateOpts.migrate, "migrate", false, "Force overwrite all template files")
	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	syncer := &update.Syncer{
		EmbedFS:   embed.Templates,
		TargetDir: paths.TrellisDir,
	}

	var result *update.SyncResult
	switch {
	case updateOpts.dryRun:
		result, err = syncer.DryRun()
	case updateOpts.migrate:
		result, err = syncer.Migrate()
	default:
		result, err = syncer.Sync()
	}
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	printUpdateResult(result)
	return nil
}

func printUpdateResult(result *update.SyncResult) {
	hasChanges := len(result.Added) > 0 || len(result.Updated) > 0 || len(result.Sections) > 0

	if !hasChanges && len(result.Skipped) == 0 {
		fmt.Println("Already up to date")
		return
	}

	if len(result.Added) > 0 {
		fmt.Println("Added:")
		for _, f := range result.Added {
			fmt.Printf("  + %s\n", filepath.Join(".trellis", f))
		}
	}
	if len(result.Updated) > 0 {
		fmt.Println("Updated:")
		for _, f := range result.Updated {
			fmt.Printf("  ~ %s\n", filepath.Join(".trellis", f))
		}
	}
	if len(result.Skipped) > 0 {
		fmt.Println("Skipped (user modified):")
		for _, f := range result.Skipped {
			fmt.Printf("  - %s\n", filepath.Join(".trellis", f))
		}
	}
	if len(result.Sections) > 0 {
		fmt.Println("Config sections appended:")
		for _, s := range result.Sections {
			fmt.Printf("  + [%s]\n", s)
		}
	}

	if updateOpts.dryRun {
		fmt.Println("\n(Dry run — no changes written)")
	}
}

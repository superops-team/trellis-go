package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateOpts struct {
	dryRun bool
}

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Trellis templates and configuration",
		RunE:  runUpdate,
	}
	cmd.Flags().BoolVar(&updateOpts.dryRun, "dry-run", false, "Preview changes without writing")
	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// TODO: implement update logic
	if updateOpts.dryRun {
		fmt.Println("Dry run: no changes will be made")
	}
	fmt.Println("Update not yet implemented")
	return nil
}

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage context",
	}
	cmd.AddCommand(
		newContextAddCmd(),
		newContextBuildCmd(),
	)
	return cmd
}

func newContextAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <file>",
		Short: "Add a file to the current task context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Context add not yet implemented")
			return nil
		},
	}
}

func newContextBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "Build and output the current context",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Context build not yet implemented")
			return nil
		},
	}
}

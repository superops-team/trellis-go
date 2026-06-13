package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	root    string
	verbose bool
	noColor bool
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trellis",
		Short: "An engineering framework for AI coding",
		Long:  `Trellis persists specs, tasks, and memory into your repo so any coding agent works to your engineering standards.`,
	}
	cmd.PersistentFlags().StringVar(&root, "root", "", "Trellis root directory (default: auto-detect)")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	cmd.AddCommand(
		newInitCmd(),
		newUpdateCmd(),
		newUninstallCmd(),
		newTaskCmd(),
		newContextCmd(),
		newVersionCmd(),
	)
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("trellis", version)
		},
	}
}

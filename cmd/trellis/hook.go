package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/superops-team/trellis-go/pkg/workflow"
)

var hookOpts struct {
	taskID string
	phase  string
	state  string
}

func newHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Run Trellis agent hook commands",
	}
	cmd.AddCommand(
		newHookSessionStartCmd(),
		newHookInjectContextCmd(),
		newHookInjectWorkflowStateCmd(),
	)
	return cmd
}

func newHookSessionStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "session-start",
		Short: "Print Trellis session start hook context",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := resolveCommandPaths()
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Trellis session started\nRepoRoot: %s\nTrellisDir: %s\n", paths.RepoRoot, paths.TrellisDir)
			return err
		},
	}
}

func newHookInjectContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inject-context",
		Short: "Print context for an agent hook",
		RunE:  runHookInjectContext,
	}
	cmd.Flags().StringVar(&hookOpts.taskID, "task", "", "Task ID")
	cmd.Flags().StringVar(&hookOpts.phase, "phase", "implement", "Context phase: implement, check, or research")
	return cmd
}

func newHookInjectWorkflowStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inject-workflow-state",
		Short: "Print workflow-state prompt for an agent hook",
		RunE:  runHookInjectWorkflowState,
	}
	cmd.Flags().StringVar(&hookOpts.state, "state", "", "Workflow state: plan, implement, verify, or finish")
	return cmd
}

func runHookInjectWorkflowState(cmd *cobra.Command, args []string) error {
	parser := &workflow.Parser{}
	prompt, err := parser.InjectPrompt(workflow.State(hookOpts.state))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), prompt)
	return err
}

func runHookInjectContext(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	phase, err := buildPhase(hookOpts.phase)
	if err != nil {
		return err
	}
	output, err := buildContextOutput(paths, hookOpts.taskID, phase)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), output)
	return err
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/superops-team/trellis-go/pkg/session"
	"github.com/superops-team/trellis-go/pkg/workflow"
)

var hookOpts struct {
	taskID  string
	phase   string
	state   string
	title   string
	commits string
	summary string
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
		newHookRecordSessionCmd(),
		newHookListSessionsCmd(),
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

func newHookRecordSessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-session",
		Short: "Record a session entry in the journal",
		RunE:  runHookRecordSession,
	}
	cmd.Flags().StringVar(&hookOpts.title, "title", "", "Session title")
	cmd.Flags().StringVar(&hookOpts.taskID, "task", "", "Task ID")
	cmd.Flags().StringVar(&hookOpts.commits, "commits", "", "Comma-separated commit hashes")
	cmd.Flags().StringVar(&hookOpts.summary, "summary", "", "Session summary")
	return cmd
}

func runHookRecordSession(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	// Read developer identity
	developer, err := readDeveloperIdentity(paths.TrellisDir)
	if err != nil {
		developer = os.Getenv("USER")
		if developer == "" {
			developer = "developer"
		}
	}

	workspaceDir := filepath.Join(paths.TrellisDir, "workspace", developer)
	recorder := &session.SessionRecorder{
		WorkspaceDir: workspaceDir,
		Config:       session.DefaultSessionConfig(),
	}

	var commitList []string
	if hookOpts.commits != "" {
		for _, c := range strings.Split(hookOpts.commits, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				commitList = append(commitList, c)
			}
		}
	}

	entry := session.SessionEntry{
		Title:   hookOpts.title,
		TaskID:  hookOpts.taskID,
		Commits: commitList,
		Summary: hookOpts.summary,
	}

	if err := recorder.RecordSession(entry); err != nil {
		return fmt.Errorf("record session: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Session recorded: %s\n", entry.Title)
	return nil
}

func newHookListSessionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-sessions",
		Short: "List recorded sessions",
		RunE:  runHookListSessions,
	}
}

func runHookListSessions(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	developer, err := readDeveloperIdentity(paths.TrellisDir)
	if err != nil {
		developer = os.Getenv("USER")
		if developer == "" {
			developer = "developer"
		}
	}

	workspaceDir := filepath.Join(paths.TrellisDir, "workspace", developer)
	recorder := &session.SessionRecorder{
		WorkspaceDir: workspaceDir,
		Config:       session.DefaultSessionConfig(),
	}

	sessions, err := recorder.ListSessions()
	if err != nil {
		return fmt.Errorf("list sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No sessions recorded")
		return nil
	}

	for _, s := range sessions {
		fmt.Fprintf(cmd.OutOrStdout(), "%s  %s  %s\n", s.StartedAt[:10], s.Title, s.JournalFile)
	}
	return nil
}

package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	trelliscontext "github.com/superops-team/trellis-go/pkg/context"
	"github.com/superops-team/trellis-go/pkg/spec"
	trellistask "github.com/superops-team/trellis-go/pkg/task"
)

var contextOpts struct {
	taskID      string
	phase       string
	required    bool
	description string
}

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
	cmd := &cobra.Command{
		Use:   "add <file>",
		Short: "Add a file to the current task context",
		Args:  cobra.ExactArgs(1),
		RunE:  runContextAdd,
	}
	cmd.Flags().StringVar(&contextOpts.taskID, "task", "", "Task ID")
	cmd.Flags().StringVar(&contextOpts.phase, "phase", "implement", "Context phase: implement or check")
	cmd.Flags().BoolVar(&contextOpts.required, "required", false, "Mark context entry as required")
	cmd.Flags().StringVar(&contextOpts.description, "description", "", "Context entry description")
	return cmd
}

func newContextBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build and output the current context",
		RunE:  runContextBuild,
	}
	cmd.Flags().StringVar(&contextOpts.taskID, "task", "", "Task ID")
	cmd.Flags().StringVar(&contextOpts.phase, "phase", "implement", "Context phase: implement, check, or research")
	return cmd
}

func runContextAdd(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	if contextOpts.taskID == "" {
		return fmt.Errorf("--task is required")
	}

	entryPath, err := validateContextPath(args[0])
	if err != nil {
		return err
	}
	phase, err := taskPhase(contextOpts.phase)
	if err != nil {
		return err
	}

	mgr := trellistask.NewManager(paths.TasksDir)
	entry := trellistask.ContextEntry{
		Path:        entryPath,
		Description: contextOpts.description,
		Required:    contextOpts.required,
	}
	if err := mgr.AddContext(contextOpts.taskID, phase, entry); err != nil {
		return err
	}
	fmt.Printf("Added context: %s\n", entryPath)
	return nil
}

func runContextBuild(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	builder := &trelliscontext.Builder{
		SpecLoader: spec.NewLoader(paths.SpecDir),
		Root:       paths.TrellisDir,
	}

	var output string
	switch contextOpts.phase {
	case "implement":
		if contextOpts.taskID == "" {
			return fmt.Errorf("--task is required for implement context")
		}
		taskDir, err := taskDirForID(paths.TasksDir, contextOpts.taskID)
		if err != nil {
			return err
		}
		output, err = builder.BuildImplementContext(taskDir)
		if err != nil {
			return err
		}
	case "check":
		if contextOpts.taskID == "" {
			return fmt.Errorf("--task is required for check context")
		}
		taskDir, err := taskDirForID(paths.TasksDir, contextOpts.taskID)
		if err != nil {
			return err
		}
		output, err = builder.BuildCheckContext(taskDir)
		if err != nil {
			return err
		}
	case "research":
		output, err = builder.BuildResearchContext()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown phase: %s", contextOpts.phase)
	}

	fmt.Println(output)
	return nil
}

func taskPhase(phase string) (trellistask.Phase, error) {
	switch phase {
	case "implement":
		return trellistask.PhaseImplement, nil
	case "check":
		return trellistask.PhaseCheck, nil
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func validateContextPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("context path must be relative: %s", path)
	}
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("context path is required")
	}
	parts := strings.Split(cleaned, string(filepath.Separator))
	for _, part := range parts {
		if part == ".." {
			return "", fmt.Errorf("context path cannot contain ..: %s", path)
		}
	}
	return filepath.ToSlash(cleaned), nil
}

func taskDirForID(tasksDir, taskID string) (string, error) {
	mgr := trellistask.NewManager(tasksDir)
	task, err := mgr.Get(taskID)
	if err != nil {
		return "", err
	}
	return filepath.Join(tasksDir, task.DirName()), nil
}

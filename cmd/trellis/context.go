package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	trelliscontext "github.com/superops-team/trellis-go/pkg/context"
	"github.com/superops-team/trellis-go/pkg/manifest"
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
	entry := manifest.Entry{
		Path:        entryPath,
		Description: contextOpts.description,
		Required:    contextOpts.required,
	}
	if err := mgr.AddContextEntry(contextOpts.taskID, phase, entry); err != nil {
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
	phase, err := buildPhase(contextOpts.phase)
	if err != nil {
		return err
	}
	output, err := buildContextOutput(paths, contextOpts.taskID, phase)
	if err != nil {
		return err
	}

	fmt.Println(output)
	return nil
}

func buildContextOutput(paths resolvedPaths, taskID string, phase trellistask.Phase) (string, error) {

	builder := &trelliscontext.Builder{
		SpecLoader: spec.NewLoader(paths.SpecDir),
		Root:       paths.TrellisDir,
	}

	var output string
	var err error
	switch phase {
	case trellistask.PhaseImplement:
		if taskID == "" {
			return "", fmt.Errorf("--task is required for implement context")
		}
		taskDir, err := taskDirForID(paths.TasksDir, taskID)
		if err != nil {
			return "", err
		}
		output, err = builder.BuildImplementContext(taskDir)
		if err != nil {
			return "", contextBuildError(taskID, err)
		}
	case trellistask.PhaseCheck:
		if taskID == "" {
			return "", fmt.Errorf("--task is required for check context")
		}
		taskDir, err := taskDirForID(paths.TasksDir, taskID)
		if err != nil {
			return "", err
		}
		output, err = builder.BuildCheckContext(taskDir)
		if err != nil {
			return "", contextBuildError(taskID, err)
		}
	case trellistask.PhaseResearch:
		output, err = builder.BuildResearchContext()
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
	return output, nil
}

func contextBuildError(taskID string, err error) error {
	if errors.Is(err, trelliscontext.ErrPRDRequired) {
		return fmt.Errorf("PRD is required for task %s", taskID)
	}
	return err
}

func taskPhase(phase string) (trellistask.Phase, error) {
	switch phase {
	case "implement":
		return trellistask.PhaseImplement, nil
	case "check":
		return trellistask.PhaseCheck, nil
	case "research":
		return trellistask.PhaseResearch, nil
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func buildPhase(phase string) (trellistask.Phase, error) {
	switch phase {
	case string(trellistask.PhaseImplement):
		return trellistask.PhaseImplement, nil
	case string(trellistask.PhaseCheck):
		return trellistask.PhaseCheck, nil
	case string(trellistask.PhaseResearch):
		return trellistask.PhaseResearch, nil
	default:
		return "", fmt.Errorf("unknown phase: %s", phase)
	}
}

func validateContextPath(rawPath string) (string, error) {
	cleaned, err := trelliscontext.NormalizeEntryPath(rawPath)
	if err == nil {
		return cleaned, nil
	}
	errText := err.Error()
	switch {
	case strings.Contains(errText, "must be relative"):
		return "", fmt.Errorf("context path must be relative: %s", rawPath)
	case strings.Contains(errText, "path is required"):
		return "", fmt.Errorf("context path is required")
	case strings.Contains(errText, "cannot contain .."):
		return "", fmt.Errorf("context path cannot contain ..: %s", rawPath)
	default:
		return "", err
	}
}

func taskDirForID(tasksDir, taskID string) (string, error) {
	mgr := trellistask.NewManager(tasksDir)
	return mgr.GetDir(taskID)
}

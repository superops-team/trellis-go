package main

import (
	"fmt"

	"github.com/spf13/cobra"
	trellistask "github.com/superops-team/trellis-go/pkg/task"
)

func newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}
	cmd.AddCommand(
		newTaskCreateCmd(),
		newTaskStartCmd(),
		newTaskArchiveCmd(),
		newTaskListCmd(),
		newTaskCurrentCmd(),
	)
	return cmd
}

func newTaskCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new task",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskCreate,
	}
}

func runTaskCreate(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}

	mgr := trellistask.NewManager(paths.TasksDir)
	_, taskDir, err := mgr.Create(args[0], trellistask.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Created task: %s\n", taskDir)
	return nil
}

func newTaskStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [id]",
		Short: "Start a task (move from planning to in_progress)",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskStart,
	}
}

func runTaskStart(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	if err := mgr.Start(args[0]); err != nil {
		return err
	}
	fmt.Printf("Started task: %s\n", args[0])
	return nil
}

func newTaskArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive [id]",
		Short: "Archive a completed task",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskArchive,
	}
}

func runTaskArchive(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	if err := mgr.Archive(args[0]); err != nil {
		return err
	}
	fmt.Printf("Archived task: %s\n", args[0])
	return nil
}

func newTaskListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE:  runTaskList,
	}
}

func runTaskList(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	tasks, err := mgr.List()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		fmt.Println(task.ID)
	}
	return nil
}

func newTaskCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current active task",
		RunE:  runTaskCurrent,
	}
}

func runTaskCurrent(cmd *cobra.Command, args []string) error {
	// TODO: read from .runtime/sessions/
	fmt.Println("No active task")
	return nil
}

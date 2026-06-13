package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mindfold/trellis/pkg/fsutil"
	"github.com/spf13/cobra"
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
	cwd, _ := os.Getwd()
	tasksDir := filepath.Join(cwd, ".trellis", "tasks")
	if err := fsutil.EnsureDir(tasksDir); err != nil {
		return err
	}

	name := args[0]
	now := time.Now()
	taskDir := filepath.Join(tasksDir, fmt.Sprintf("%02d-%02d-%s", now.Month(), now.Day(), name))
	if err := fsutil.EnsureDir(taskDir); err != nil {
		return err
	}

	taskJSON := filepath.Join(taskDir, "task.json")
	taskData := fmt.Sprintf(`{
  "id": "%s",
  "name": "%s",
  "status": "planning",
  "assignee": "",
  "branch": "",
  "subtasks": [],
  "created_at": "%s",
  "updated_at": "%s"
}
`, name, name, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err := os.WriteFile(taskJSON, []byte(taskData), 0644); err != nil {
		return err
	}

	// Create empty files
	for _, f := range []string{"prd.md", "implement.jsonl", "check.jsonl"} {
		path := filepath.Join(taskDir, f)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			return err
		}
	}
	researchDir := filepath.Join(taskDir, "research")
	if err := fsutil.EnsureDir(researchDir); err != nil {
		return err
	}

	fmt.Printf("Created task: %s\n", taskDir)
	return nil
}

func newTaskStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [id]",
		Short: "Start a task (move from planning to in_progress)",
		RunE:  runTaskStart,
	}
}

func runTaskStart(cmd *cobra.Command, args []string) error {
	// TODO: implement proper task lookup and status transition
	fmt.Println("Task start not yet fully implemented")
	return nil
}

func newTaskArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive [id]",
		Short: "Archive a completed task",
		RunE:  runTaskArchive,
	}
}

func runTaskArchive(cmd *cobra.Command, args []string) error {
	// TODO: implement archive logic
	fmt.Println("Task archive not yet fully implemented")
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
	cwd, _ := os.Getwd()
	tasksDir := filepath.Join(cwd, ".trellis", "tasks")
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		return fmt.Errorf("read tasks: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "archive" {
			fmt.Println(entry.Name())
		}
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

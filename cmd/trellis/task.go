package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	trellistask "github.com/superops-team/trellis-go/pkg/task"
	"github.com/superops-team/trellis-go/pkg/manifest"
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
		newTaskInfoCmd(),
		newTaskEditCmd(),
		newTaskAddSubtaskCmd(),
		newTaskDoneSubtaskCmd(),
		newTaskUndoneSubtaskCmd(),
		newTaskAddContextCmd(),
		newTaskRemoveContextCmd(),
		newTaskListContextCmd(),
		newTaskAddSpecCmd(),
		newTaskListSpecsCmd(),
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
		Use:   "start <id>",
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
		Use:   "archive <id>",
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
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE:  runTaskList,
	}
	cmd.Flags().String("status", "", "Filter by status (planning, in_progress, completed)")
	cmd.Flags().String("format", "table", "Output format (table, json)")
	return cmd
}

func runTaskList(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	status, _ := cmd.Flags().GetString("status")
	format, _ := cmd.Flags().GetString("format")

	var tasks []trellistask.Task
	if status != "" {
		tasks, err = mgr.ListByStatus(trellistask.Status(status))
	} else {
		tasks, err = mgr.List()
	}
	if err != nil {
		return err
	}

	switch format {
	case "json":
		return printTasksJSON(tasks)
	default:
		return printTasksTable(tasks)
	}
}

func printTasksJSON(tasks []trellistask.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func printTasksTable(tasks []trellistask.Task) error {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tPACKAGE\tCREATED")
	for _, t := range tasks {
		pkg := t.Package
		if pkg == "" {
			pkg = "-"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			t.ID, t.Name, t.Status, pkg,
			t.CreatedAt.Format("2006-01-02 15:04"))
	}
	return w.Flush()
}

func newTaskCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the current active task",
		RunE:  runTaskCurrent,
	}
}

func runTaskCurrent(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	tasks, err := mgr.ListByStatus(trellistask.StatusInProgress)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No active task")
		return nil
	}
	t := tasks[0]
	fmt.Printf("Active task: %s (%s)\n", t.ID, t.Status)
	fmt.Printf("  Name:    %s\n", t.Name)
	if t.Package != "" {
		fmt.Printf("  Package: %s\n", t.Package)
	}
	if t.Branch != "" {
		fmt.Printf("  Branch:  %s\n", t.Branch)
	}
	fmt.Printf("  Updated: %s\n", t.UpdatedAt.Format("2006-01-02 15:04"))
	return nil
}

// --- task info ---

func newTaskInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskInfo,
	}
}

func runTaskInfo(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	task, err := mgr.Get(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("ID:        %s\n", task.ID)
	fmt.Printf("Name:      %s\n", task.Name)
	fmt.Printf("Status:    %s\n", task.Status)
	fmt.Printf("Assignee:  %s\n", task.Assignee)
	fmt.Printf("Branch:    %s\n", task.Branch)
	if task.Package != "" {
		fmt.Printf("Package:   %s\n", task.Package)
	}
	fmt.Printf("Created:   %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
	fmt.Printf("Updated:   %s\n", task.UpdatedAt.Format("2006-01-02 15:04"))

	if len(task.Specs) > 0 {
		fmt.Println("Specs:")
		for _, s := range task.Specs {
			fmt.Printf("  - %s\n", s)
		}
	}

	if len(task.Subtasks) > 0 {
		fmt.Println("Subtasks:")
		for _, s := range task.Subtasks {
			done := " "
			if s.Done {
				done = "✓"
			}
			fmt.Printf("  [%s] %s: %s\n", done, s.ID, s.Title)
		}
	}

	dir, err := mgr.GetDir(args[0])
	if err == nil {
		fmt.Printf("Dir:       %s\n", dir)
	}
	return nil
}

// --- task edit ---

func newTaskEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit task fields",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskEdit,
	}
	cmd.Flags().String("name", "", "New task name")
	cmd.Flags().String("assignee", "", "New assignee")
	cmd.Flags().String("branch", "", "New branch")
	cmd.Flags().String("package", "", "New package")
	cmd.Flags().String("status", "", "New status")
	return cmd
}

func runTaskEdit(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	var patch trellistask.TaskPatch
	if cmd.Flags().Changed("name") {
		v, _ := cmd.Flags().GetString("name")
		patch.Name = &v
	}
	if cmd.Flags().Changed("assignee") {
		v, _ := cmd.Flags().GetString("assignee")
		patch.Assignee = &v
	}
	if cmd.Flags().Changed("branch") {
		v, _ := cmd.Flags().GetString("branch")
		patch.Branch = &v
	}
	if cmd.Flags().Changed("package") {
		v, _ := cmd.Flags().GetString("package")
		patch.Package = &v
	}
	if cmd.Flags().Changed("status") {
		v, _ := cmd.Flags().GetString("status")
		s := trellistask.Status(v)
		patch.Status = &s
	}

	if err := mgr.Edit(args[0], patch); err != nil {
		return err
	}
	fmt.Printf("Updated task: %s\n", args[0])
	return nil
}

// --- task subtask ---

func newTaskAddSubtaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-subtask <task-id> <title>",
		Short: "Add a subtask",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskAddSubtask,
	}
}

func runTaskAddSubtask(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	sub, err := mgr.AddSubtask(args[0], args[1])
	if err != nil {
		return err
	}
	fmt.Printf("Added subtask %s to task %s\n", sub.ID, args[0])
	return nil
}

func newTaskDoneSubtaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done-subtask <task-id> <subtask-id>",
		Short: "Mark a subtask as done",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskDoneSubtask,
	}
}

func runTaskDoneSubtask(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	if err := mgr.DoneSubtask(args[0], args[1]); err != nil {
		return err
	}
	fmt.Printf("Marked subtask %s as done in task %s\n", args[1], args[0])
	return nil
}

func newTaskUndoneSubtaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undone-subtask <task-id> <subtask-id>",
		Short: "Mark a subtask as not done",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskUndoneSubtask,
	}
}

func runTaskUndoneSubtask(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	if err := mgr.UndoneSubtask(args[0], args[1]); err != nil {
		return err
	}
	fmt.Printf("Marked subtask %s as not done in task %s\n", args[1], args[0])
	return nil
}

// --- task context ---

func newTaskAddContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-context <task-id> <file>",
		Short: "Add a context file to a task phase",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskAddContext,
	}
	cmd.Flags().String("phase", "implement", "Phase (implement, check, research)")
	cmd.Flags().Bool("required", false, "Mark as required")
	cmd.Flags().String("description", "", "Entry description")
	return cmd
}

func runTaskAddContext(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	phase, _ := cmd.Flags().GetString("phase")
	required, _ := cmd.Flags().GetBool("required")
	desc, _ := cmd.Flags().GetString("description")

	entry := manifest.Entry{
		Path:        args[1],
		Required:    required,
		Description: desc,
	}

	if err := mgr.AddContextEntry(args[0], trellistask.Phase(phase), entry); err != nil {
		return err
	}
	fmt.Printf("Added context %q to task %s (%s)\n", args[1], args[0], phase)
	return nil
}

func newTaskRemoveContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-context <task-id> <file>",
		Short: "Remove a context file from a task phase",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskRemoveContext,
	}
	cmd.Flags().String("phase", "implement", "Phase (implement, check, research)")
	return cmd
}

func runTaskRemoveContext(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	phase, _ := cmd.Flags().GetString("phase")

	if err := mgr.RemoveContextEntry(args[0], trellistask.Phase(phase), args[1]); err != nil {
		return err
	}
	fmt.Printf("Removed context %q from task %s (%s)\n", args[1], args[0], phase)
	return nil
}

func newTaskListContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-context <task-id>",
		Short: "List context files for a task phase",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskListContext,
	}
	cmd.Flags().String("phase", "implement", "Phase (implement, check, research)")
	return cmd
}

func runTaskListContext(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)

	phase, _ := cmd.Flags().GetString("phase")

	entries, err := mgr.ListContextEntries(args[0], trellistask.Phase(phase))
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Printf("No context entries for task %s (%s)\n", args[0], phase)
		return nil
	}
	fmt.Printf("Context for task %s (%s):\n", args[0], phase)
	for _, e := range entries {
		req := " "
		if e.Required {
			req = "*"
		}
		desc := e.Description
		if desc != "" {
			desc = " — " + desc
		}
		fmt.Printf("  [%s] %s%s\n", req, e.Path, desc)
	}
	return nil
}

// --- task spec ---

func newTaskAddSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-spec <task-id> <spec-path>",
		Short: "Associate a spec file with a task",
		Args:  cobra.ExactArgs(2),
		RunE:  runTaskAddSpec,
	}
}

func runTaskAddSpec(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	if err := mgr.AddSpec(args[0], args[1]); err != nil {
		return err
	}
	fmt.Printf("Added spec %q to task %s\n", args[1], args[0])
	return nil
}

func newTaskListSpecsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-specs <task-id>",
		Short: "List associated specs for a task",
		Args:  cobra.ExactArgs(1),
		RunE:  runTaskListSpecs,
	}
}

func runTaskListSpecs(cmd *cobra.Command, args []string) error {
	paths, err := resolveCommandPaths()
	if err != nil {
		return err
	}
	mgr := trellistask.NewManager(paths.TasksDir)
	specs, err := mgr.ListSpecs(args[0])
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		fmt.Printf("No specs associated with task %s\n", args[0])
		return nil
	}
	fmt.Printf("Specs for task %s:\n", args[0])
	for _, s := range specs {
		fmt.Printf("  - %s\n", s)
	}
	return nil
}

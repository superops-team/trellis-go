package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/superops-team/trellis-go/pkg/git"
	"github.com/superops-team/trellis-go/pkg/spec"
	"github.com/superops-team/trellis-go/pkg/task"
	"github.com/superops-team/trellis-go/pkg/workflow"
)

// SessionContext holds the assembled session startup context.
type SessionContext struct {
	Developer    string
	Repository   string
	Branch       string
	IsDirty      bool
	ActiveTask   *task.Task
	Workflow     string
	SpecIndex    string
	RecentTasks  []task.Task
}

// RecordContext holds the assembled session recording context.
type RecordContext struct {
	ActiveTasks        []task.Task
	Branch             string
	RecentCommits      []git.CommitInfo
	UnarchivedComplete []task.Task
}

// FullBuilder extends Builder with task manager, git client, and config.
type FullBuilder struct {
	SpecLoader  *spec.Loader
	TaskManager *task.Manager
	GitClient   *git.Client
	Root        string
}

// BuildSessionContext assembles the session startup context.
func (b *FullBuilder) BuildSessionContext() (*SessionContext, error) {
	ctx := &SessionContext{}

	// Developer identity
	dev, err := readDeveloperFile(b.Root)
	if err == nil {
		ctx.Developer = dev
	}

	// Git status
	if b.GitClient != nil && b.GitClient.IsRepo() {
		if url, err := b.GitClient.RemoteURL(); err == nil {
			ctx.Repository = url
		}
		if branch, err := b.GitClient.CurrentBranch(); err == nil {
			ctx.Branch = branch
		}
		if dirty, err := b.GitClient.HasChanges(); err == nil {
			ctx.IsDirty = dirty
		}
	}

	// Active task
	if b.TaskManager != nil {
		if t, err := b.TaskManager.Current(); err == nil {
			ctx.ActiveTask = t
		}
	}

	// Workflow
	wfPath := filepath.Join(b.Root, "workflow.md")
	if data, err := os.ReadFile(wfPath); err == nil {
		ctx.Workflow = string(data)
	}

	// Spec index
	if b.SpecLoader != nil {
		if idx, err := b.SpecLoader.Index(); err == nil {
			ctx.SpecIndex = idx.ToMarkdown()
		}
	}

	// Recent tasks
	if b.TaskManager != nil {
		recent, err := b.TaskManager.ListRecent(5)
		if err == nil {
			ctx.RecentTasks = recent
		}
	}

	return ctx, nil
}

// FormatSessionContext formats a SessionContext as a string.
func FormatSessionContext(ctx *SessionContext) string {
	var b strings.Builder
	b.WriteString(injectMarker)
	b.WriteString("\n\n")

	if ctx.Developer != "" {
		b.WriteString(fmt.Sprintf("Developer: %s\n", ctx.Developer))
	}
	if ctx.Repository != "" {
		b.WriteString(fmt.Sprintf("Repository: %s\n", ctx.Repository))
	}
	if ctx.Branch != "" {
		b.WriteString(fmt.Sprintf("Branch: %s\n", ctx.Branch))
	}
	b.WriteString("Status: ")
	if ctx.IsDirty {
		b.WriteString("dirty")
	} else {
		b.WriteString("clean")
	}
	b.WriteString("\n")

	if ctx.ActiveTask != nil {
		b.WriteString(fmt.Sprintf("Active task: %s (%s)\n", ctx.ActiveTask.ID, ctx.ActiveTask.Status))
	}

	if ctx.Workflow != "" {
		b.WriteString("\n## Workflow\n")
		b.WriteString(ctx.Workflow)
	}

	if ctx.SpecIndex != "" {
		b.WriteString("\n## Spec Index\n")
		b.WriteString(ctx.SpecIndex)
	}

	if len(ctx.RecentTasks) > 0 {
		b.WriteString("\n## Recent Tasks\n")
		for _, t := range ctx.RecentTasks {
			b.WriteString(fmt.Sprintf("- %s (%s): %s\n", t.ID, t.Status, t.Name))
		}
	}

	return b.String()
}

// BuildRecordContext assembles the session recording context.
func (b *FullBuilder) BuildRecordContext() (*RecordContext, error) {
	ctx := &RecordContext{}

	// Active tasks
	if b.TaskManager != nil {
		active, err := b.TaskManager.ListByStatus(task.StatusInProgress)
		if err == nil {
			ctx.ActiveTasks = active
		}
	}

	// Git status
	if b.GitClient != nil && b.GitClient.IsRepo() {
		if branch, err := b.GitClient.CurrentBranch(); err == nil {
			ctx.Branch = branch
		}
		if commits, err := b.GitClient.RecentCommits(5); err == nil {
			ctx.RecentCommits = commits
		}
	}

	// Unarchived completed tasks
	if b.TaskManager != nil {
		completed, err := b.TaskManager.ListByStatus(task.StatusCompleted)
		if err == nil {
			ctx.UnarchivedComplete = completed
		}
	}

	return ctx, nil
}

// FormatRecordContext formats a RecordContext as a string.
func FormatRecordContext(ctx *RecordContext) string {
	var b strings.Builder
	b.WriteString(injectMarker)
	b.WriteString("\n\n")

	if len(ctx.ActiveTasks) > 0 {
		b.WriteString("Active tasks:\n")
		for _, t := range ctx.ActiveTasks {
			b.WriteString(fmt.Sprintf("  - %s (%s): %s\n", t.ID, t.Status, t.Name))
		}
		b.WriteString("\n")
	}

	b.WriteString("Git status:\n")
	if ctx.Branch != "" {
		b.WriteString(fmt.Sprintf("  Branch: %s\n", ctx.Branch))
	}
	if len(ctx.RecentCommits) > 0 {
		b.WriteString("  Recent commits:\n")
		for _, c := range ctx.RecentCommits {
			b.WriteString(fmt.Sprintf("    %s - %s\n", c.Hash[:7], c.Message))
		}
	}
	b.WriteString("\n")

	if len(ctx.UnarchivedComplete) > 0 {
		b.WriteString("Unarchived completed tasks:\n")
		for _, t := range ctx.UnarchivedComplete {
			b.WriteString(fmt.Sprintf("  - %s: %s\n", t.ID, t.Name))
		}
	}

	return b.String()
}

// BuildPhaseContext extracts a workflow step body.
func (b *FullBuilder) BuildPhaseContext(step string) (string, error) {
	wfPath := filepath.Join(b.Root, "workflow.md")
	data, err := os.ReadFile(wfPath)
	if err != nil {
		return "", fmt.Errorf("read workflow.md: %w", err)
	}

	parser := &workflow.Parser{}
	body, err := parser.ExtractStep(string(data), step)
	if err != nil {
		return "", fmt.Errorf("extract step %s: %w", step, err)
	}

	return fmt.Sprintf("%s\n\n=== file: workflow.md ===\n%s", injectMarker, body), nil
}

// readDeveloperFile reads .trellis/.developer.
func readDeveloperFile(trellisDir string) (string, error) {
	devPath := filepath.Join(trellisDir, ".developer")
	data, err := os.ReadFile(devPath)
	if err != nil {
		return "", err
	}
	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", fmt.Errorf("empty .developer")
	}
	return name, nil
}

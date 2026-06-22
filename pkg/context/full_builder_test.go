package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/git"
	"github.com/superops-team/trellis-go/pkg/spec"
	"github.com/superops-team/trellis-go/pkg/task"
)

func TestFullBuilder_BuildSessionContext_Developer(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, ".developer"), []byte("Alice"), 0644)

	b := &FullBuilder{Root: root}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if ctx.Developer != "Alice" {
		t.Errorf("expected Developer 'Alice', got %q", ctx.Developer)
	}
}

func TestFullBuilder_BuildSessionContext_NoDeveloper(t *testing.T) {
	root := t.TempDir()

	b := &FullBuilder{Root: root}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if ctx.Developer != "" {
		t.Errorf("expected empty Developer, got %q", ctx.Developer)
	}
}

func TestFullBuilder_BuildSessionContext_GitStatus(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "checkout", "-b", "main")
	runGit(t, root, "remote", "add", "origin", "git@github.com:test/repo.git")

	b := &FullBuilder{
		Root:      root,
		GitClient: git.NewClient(root),
	}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if ctx.Branch != "main" {
		t.Errorf("expected Branch 'main', got %q", ctx.Branch)
	}
	if ctx.Repository != "git@github.com:test/repo.git" {
		t.Errorf("expected Repository URL, got %q", ctx.Repository)
	}
	if ctx.IsDirty {
		t.Error("expected clean status")
	}
}

func TestFullBuilder_BuildSessionContext_ActiveTask(t *testing.T) {
	root := t.TempDir()
	tasksDir := filepath.Join(root, "tasks")
	os.MkdirAll(tasksDir, 0755)

	mgr := task.NewManager(tasksDir)
	tk, taskDir, err := mgr.Create("test-task", task.CreateOptions{})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	// Start requires a non-empty PRD
	os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nTest PRD content."), 0644)
	_ = mgr.Start(tk.ID)

	b := &FullBuilder{
		Root:        root,
		TaskManager: mgr,
	}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if ctx.ActiveTask == nil {
		t.Fatal("expected active task")
	}
	if ctx.ActiveTask.ID != tk.ID {
		t.Errorf("expected task ID %s, got %s", tk.ID, ctx.ActiveTask.ID)
	}
}

func TestFullBuilder_BuildSessionContext_Workflow(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "workflow.md"), []byte("# Workflow\n## Plan\n"), 0644)

	b := &FullBuilder{Root: root}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if !strings.Contains(ctx.Workflow, "## Plan") {
		t.Errorf("expected workflow content, got: %s", ctx.Workflow)
	}
}

func TestFullBuilder_BuildSessionContext_SpecIndex(t *testing.T) {
	root := t.TempDir()
	layerDir := filepath.Join(root, "auth", "api")
	os.MkdirAll(layerDir, 0755)
	os.WriteFile(filepath.Join(layerDir, "index.md"), []byte("# Auth API"), 0644)

	b := &FullBuilder{
		Root:       root,
		SpecLoader: spec.NewLoader(root),
	}
	ctx, err := b.BuildSessionContext()
	if err != nil {
		t.Fatalf("BuildSessionContext failed: %v", err)
	}
	if !strings.Contains(ctx.SpecIndex, "auth") {
		t.Errorf("expected spec index to contain 'auth', got: %s", ctx.SpecIndex)
	}
}

func TestFormatSessionContext(t *testing.T) {
	ctx := &SessionContext{
		Developer:  "Alice",
		Repository: "git@github.com:test/repo.git",
		Branch:     "main",
		IsDirty:    false,
		Workflow:   "# Workflow\n## Plan\n",
		SpecIndex:  "# Spec Index\n## auth\n",
		RecentTasks: []task.Task{
			{ID: "abc123", Name: "Login", Status: task.StatusCompleted},
		},
	}

	output := FormatSessionContext(ctx)
	for _, want := range []string{
		injectMarker,
		"Developer: Alice",
		"Repository: git@github.com:test/repo.git",
		"Branch: main",
		"Status: clean",
		"## Workflow",
		"## Spec Index",
		"## Recent Tasks",
		"abc123",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output should contain %q, got:\n%s", want, output)
		}
	}
}

func TestFormatSessionContext_Dirty(t *testing.T) {
	ctx := &SessionContext{IsDirty: true}
	output := FormatSessionContext(ctx)
	if !strings.Contains(output, "Status: dirty") {
		t.Errorf("expected dirty status, got: %s", output)
	}
}

func TestFormatSessionContext_ActiveTask(t *testing.T) {
	ctx := &SessionContext{
		ActiveTask: &task.Task{ID: "abc123", Name: "Login", Status: task.StatusInProgress},
	}
	output := FormatSessionContext(ctx)
	if !strings.Contains(output, "Active task: abc123 (in_progress)") {
		t.Errorf("expected active task, got: %s", output)
	}
}

func TestFullBuilder_BuildRecordContext(t *testing.T) {
	root := t.TempDir()
	tasksDir := filepath.Join(root, "tasks")
	os.MkdirAll(tasksDir, 0755)

	mgr := task.NewManager(tasksDir)
	tk, taskDir1, _ := mgr.Create("active-task", task.CreateOptions{})
	os.WriteFile(filepath.Join(taskDir1, "prd.md"), []byte("# PRD\nTest."), 0644)
	_ = mgr.Start(tk.ID)

	completed, taskDir2, _ := mgr.Create("done-task", task.CreateOptions{})
	os.WriteFile(filepath.Join(taskDir2, "prd.md"), []byte("# PRD\nTest."), 0644)
	_ = mgr.Start(completed.ID)
	_ = mgr.Archive(completed.ID)

	runGit(t, root, "init")
	runGit(t, root, "checkout", "-b", "main")
	runGit(t, root, "remote", "add", "origin", "git@github.com:test/repo.git")
	// Create a commit so RecentCommits returns something
	os.WriteFile(filepath.Join(root, "dummy"), []byte("x"), 0644)
	runGit(t, root, "add", "dummy")
	runGit(t, root, "commit", "-m", "initial commit")

	b := &FullBuilder{
		Root:        root,
		TaskManager: mgr,
		GitClient:   git.NewClient(root),
	}
	ctx, err := b.BuildRecordContext()
	if err != nil {
		t.Fatalf("BuildRecordContext failed: %v", err)
	}
	if len(ctx.ActiveTasks) == 0 {
		t.Error("expected active tasks")
	}
	if ctx.Branch != "main" {
		t.Errorf("expected branch 'main', got %q", ctx.Branch)
	}
	if len(ctx.RecentCommits) == 0 {
		t.Error("expected recent commits")
	}
}

func TestFormatRecordContext(t *testing.T) {
	ctx := &RecordContext{
		ActiveTasks: []task.Task{
			{ID: "abc123", Name: "Login", Status: task.StatusInProgress},
		},
		Branch: "main",
		RecentCommits: []git.CommitInfo{
			{Hash: "abcdef1234567890", Message: "feat: add login"},
		},
		UnarchivedComplete: []task.Task{
			{ID: "def456", Name: "Setup"},
		},
	}

	output := FormatRecordContext(ctx)
	for _, want := range []string{
		injectMarker,
		"Active tasks:",
		"abc123",
		"Branch: main",
		"abcdef1 - feat: add login",
		"Unarchived completed tasks:",
		"def456",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output should contain %q, got:\n%s", want, output)
		}
	}
}

func TestFullBuilder_BuildPhaseContext(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "workflow.md"), []byte(`# Workflow

## Phase 2: Implement

#### 2.1 Design

Design the system architecture.

#### 2.2 Review

Review the design.
`), 0644)

	b := &FullBuilder{Root: root}
	output, err := b.BuildPhaseContext("2.1")
	if err != nil {
		t.Fatalf("BuildPhaseContext failed: %v", err)
	}
	if !strings.Contains(output, injectMarker) {
		t.Errorf("expected injection marker, got: %s", output)
	}
	if !strings.Contains(output, "Design the system architecture") {
		t.Errorf("expected step content, got: %s", output)
	}
}

func TestFullBuilder_BuildPhaseContext_NotFound(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "workflow.md"), []byte("# Workflow\n## Plan\n"), 0644)

	b := &FullBuilder{Root: root}
	_, err := b.BuildPhaseContext("9.9")
	if err == nil {
		t.Fatal("expected error for non-existent step")
	}
}

func TestFullBuilder_BuildPhaseContext_NotFound_NoWorkflow(t *testing.T) {
	root := t.TempDir()

	b := &FullBuilder{Root: root}
	_, err := b.BuildPhaseContext("2.1")
	if err == nil {
		t.Fatal("expected error for missing workflow.md")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestManager_Create(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, err := mgr.Create("add-auth", CreateOptions{Assignee: "alice"})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if task.Status != StatusPlanning {
		t.Errorf("expected status planning, got %s", task.Status)
	}
	if task.Assignee != "alice" {
		t.Errorf("expected assignee alice, got %s", task.Assignee)
	}

	// Verify directory contents
	for _, f := range []string{"task.json", "prd.md", "implement.jsonl", "check.jsonl"} {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing file: %s", f)
		}
	}
}

func TestManager_CreateRejectsUnsafeTaskNames(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
	}{
		{name: "empty", taskName: ""},
		{name: "slash path", taskName: "feature/auth"},
		{name: "backslash path", taskName: `feature\auth`},
		{name: "dot", taskName: "."},
		{name: "dot dot", taskName: ".."},
		{name: "parent traversal", taskName: "../secret"},
		{name: "trimmed whitespace", taskName: " user-auth "},
		{name: "control character", taskName: "user\nauth"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			mgr := NewManager(tmp)

			_, _, err := mgr.Create(tt.taskName, CreateOptions{})
			if err == nil {
				t.Fatalf("expected unsafe task name %q to be rejected", tt.taskName)
			}
			if !strings.Contains(err.Error(), "invalid task name") {
				t.Fatalf("error should mention invalid task name, got: %v", err)
			}

			entries, err := os.ReadDir(tmp)
			if err != nil {
				t.Fatalf("read task root: %v", err)
			}
			if len(entries) != 0 {
				t.Fatalf("invalid create should not leave task directories, got %d entries", len(entries))
			}
		})
	}
}

func TestManager_Start(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, _ := mgr.Create("add-auth", CreateOptions{})
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := mgr.Start(task.ID); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	loaded, _ := mgr.Get(task.ID)
	if loaded.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %s", loaded.Status)
	}
}

func TestManager_StartRequiresNonEmptyPRD(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, err := mgr.Create("add-auth", CreateOptions{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte(" \n\t"), 0644); err != nil {
		t.Fatalf("write blank prd: %v", err)
	}

	err = mgr.Start(task.ID)
	if err == nil {
		t.Fatal("expected blank PRD to block task start")
	}
	if !strings.Contains(err.Error(), "PRD is required") {
		t.Fatalf("error should explain PRD requirement, got: %v", err)
	}
	loaded, err := mgr.Get(task.ID)
	if err != nil {
		t.Fatalf("load task after rejected start: %v", err)
	}
	if loaded.Status != StatusPlanning {
		t.Fatalf("rejected start should keep status %s, got %s", StatusPlanning, loaded.Status)
	}
}

func TestManager_StartRequiresExistingPRD(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, err := mgr.Create("add-auth", CreateOptions{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := os.Remove(filepath.Join(dir, "prd.md")); err != nil {
		t.Fatalf("remove prd: %v", err)
	}

	err = mgr.Start(task.ID)
	if err == nil {
		t.Fatal("expected missing PRD to block task start")
	}
	if !strings.Contains(err.Error(), "PRD is required") {
		t.Fatalf("error should explain PRD requirement, got: %v", err)
	}
	loaded, err := mgr.Get(task.ID)
	if err != nil {
		t.Fatalf("load task after rejected start: %v", err)
	}
	if loaded.Status != StatusPlanning {
		t.Fatalf("rejected start should keep status %s, got %s", StatusPlanning, loaded.Status)
	}
}

func TestManager_Start_InvalidTransition(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, _ := mgr.Create("add-auth", CreateOptions{})
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	mgr.Start(task.ID)
	mgr.Archive(task.ID)

	err := mgr.Start(task.ID)
	if err == nil {
		t.Error("expected error for invalid transition from completed")
	}
}

func TestManager_Archive(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, _ := mgr.Create("add-auth", CreateOptions{})
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	mgr.Start(task.ID)

	if err := mgr.Archive(task.ID); err != nil {
		t.Fatalf("Archive failed: %v", err)
	}

	// Original dir should be gone
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("original task dir should be removed")
	}

	// Should be in archive
	archiveDir := filepath.Join(tmp, "archive")
	entries, _ := os.ReadDir(archiveDir)
	if len(entries) == 0 {
		t.Error("task should be archived")
	}

	archivedTask, err := LoadTask(filepath.Join(archiveDir, entries[0].Name(), task.DirName(), "task.json"))
	if err != nil {
		t.Fatalf("load archived task: %v", err)
	}
	if archivedTask.Status != StatusCompleted {
		t.Errorf("archived task status = %s, want %s", archivedTask.Status, StatusCompleted)
	}
}

func TestManager_ArchiveRenameFailureKeepsTaskInProgress(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, dir, err := mgr.Create("add-auth", CreateOptions{})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := mgr.Start(task.ID); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	archivePath := filepath.Join(tmp, "archive", time.Now().Format("2006-01"), task.DirName())
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		t.Fatalf("create archive path blocker: %v", err)
	}

	err = mgr.Archive(task.ID)
	if err == nil {
		t.Fatal("expected archive to fail when destination already exists")
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("original task dir should remain after failed archive: %v", statErr)
	}
	loaded, err := LoadTask(filepath.Join(dir, "task.json"))
	if err != nil {
		t.Fatalf("load task after failed archive: %v", err)
	}
	if loaded.Status != StatusInProgress {
		t.Fatalf("failed archive should keep status %s, got %s", StatusInProgress, loaded.Status)
	}
}

func TestManager_List(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	mgr.Create("task-a", CreateOptions{})
	mgr.Create("task-b", CreateOptions{})

	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(list))
	}
}

func TestManager_Validate(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("add-auth", CreateOptions{})
	if err := mgr.Validate(task.ID); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestManager_AddContext(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("add-auth", CreateOptions{})
	entry := ContextEntry{Path: "spec/auth.md", Required: true}
	if err := mgr.AddContextEntry(task.ID, PhaseImplement, entry); err != nil {
		t.Fatalf("AddContext failed: %v", err)
	}

	manifestPath := filepath.Join(tmp, task.DirName(), "implement.jsonl")
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if len(manifest.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(manifest.Entries))
	}
}

func TestManager_AddContext_MalformedManifestReturnsError(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("add-auth", CreateOptions{})
	manifestPath := filepath.Join(tmp, task.DirName(), "implement.jsonl")
	if err := os.WriteFile(manifestPath, []byte("{not-json\n"), 0644); err != nil {
		t.Fatalf("write malformed manifest: %v", err)
	}

	err := mgr.AddContextEntry(task.ID, PhaseImplement, ContextEntry{Path: "spec/auth.md", Required: true})
	if err == nil {
		t.Fatal("expected malformed manifest error")
	}
	if !strings.Contains(err.Error(), manifestPath) {
		t.Fatalf("error should mention manifest path %q, got %v", manifestPath, err)
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	if string(data) != "{not-json\n" {
		t.Fatalf("malformed manifest should not be overwritten, got %q", data)
	}
}

func TestManager_ReadsExistingTaskAndManifestFixtures(t *testing.T) {
	tmp := t.TempDir()
	taskDir := filepath.Join(tmp, "03-04-legacy-task")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create legacy task dir: %v", err)
	}
	taskJSON := `{
  "id": "legacy-task",
  "name": "legacy-task",
  "status": "planning",
  "assignee": "alice",
  "branch": "feature/legacy-task",
  "subtasks": [],
  "created_at": "2026-03-04T05:06:07Z",
  "updated_at": "2026-03-04T05:06:07Z"
}`
	if err := os.WriteFile(filepath.Join(taskDir, "task.json"), []byte(taskJSON), 0644); err != nil {
		t.Fatalf("write legacy task.json: %v", err)
	}
	manifestContent := `{"path":"spec/auth.md","description":"Auth spec","required":true}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifestContent), 0644); err != nil {
		t.Fatalf("write legacy implement.jsonl: %v", err)
	}

	mgr := NewManager(tmp)
	task, err := mgr.Get("legacy-task")
	if err != nil {
		t.Fatalf("manager should load existing task.json fixture: %v", err)
	}
	if task.Status != StatusPlanning || task.Assignee != "alice" {
		t.Fatalf("loaded task mismatch: %+v", task)
	}

	manifestPath := filepath.Join(taskDir, "implement.jsonl")
	manifest, err := loadManifest(manifestPath)
	if err != nil {
		t.Fatalf("legacy implement.jsonl should load: %v", err)
	}
	if len(manifest.Entries) != 1 || manifest.Entries[0].Path != "spec/auth.md" || !manifest.Entries[0].Required {
		t.Fatalf("legacy manifest entry mismatch: %+v", manifest.Entries)
	}

	if err := mgr.AddContextEntry("legacy-task", PhaseImplement, ContextEntry{Path: "spec/api.md", Description: "API spec"}); err != nil {
		t.Fatalf("AddContext should append to existing manifest fixture: %v", err)
	}
	reloaded, err := loadManifest(manifestPath)
	if err != nil {
		t.Fatalf("reload appended manifest: %v", err)
	}
	if len(reloaded.Entries) != 2 {
		t.Fatalf("expected appended legacy manifest to contain 2 entries, got %d", len(reloaded.Entries))
	}
}

func TestManager_Current(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	// No tasks in progress
	_, err := mgr.Current()
	if err == nil {
		t.Fatal("expected error when no tasks in progress")
	}
	if err != ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got: %v", err)
	}

	// Create and start a task
	task, dir, _ := mgr.Create("active", CreateOptions{})
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte("# PRD\nActive task."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	mgr.Start(task.ID)

	// Current should return the active task
	current, err := mgr.Current()
	if err != nil {
		t.Fatalf("Current failed: %v", err)
	}
	if current.ID != task.ID {
		t.Errorf("Current ID = %s, want %s", current.ID, task.ID)
	}
	if current.Status != StatusInProgress {
		t.Errorf("Current status = %s, want %s", current.Status, StatusInProgress)
	}
}

func TestManager_GetDir(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("test-dir", CreateOptions{})

	dir, err := mgr.GetDir(task.ID)
	if err != nil {
		t.Fatalf("GetDir failed: %v", err)
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Errorf("GetDir returned non-existent path: %v", statErr)
	}

	// Non-existent task
	_, err = mgr.GetDir("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestManager_Edit(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("edit-test", CreateOptions{Assignee: "alice"})

	name := "renamed-task"
	assignee := "bob"
	branch := "feature/renamed"
	pkg := "pkg/auth"
	status := StatusInProgress

	err := mgr.Edit(task.ID, TaskPatch{
		Name:     &name,
		Assignee: &assignee,
		Branch:   &branch,
		Package:  &pkg,
		Status:   &status,
	})
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	updated, _ := mgr.Get(task.ID)
	if updated.Name != name {
		t.Errorf("Name = %s, want %s", updated.Name, name)
	}
	if updated.Assignee != assignee {
		t.Errorf("Assignee = %s, want %s", updated.Assignee, assignee)
	}
	if updated.Branch != branch {
		t.Errorf("Branch = %s, want %s", updated.Branch, branch)
	}
	if updated.Package != pkg {
		t.Errorf("Package = %s, want %s", updated.Package, pkg)
	}
	if updated.Status != status {
		t.Errorf("Status = %s, want %s", updated.Status, status)
	}
}

func TestManager_EditNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	err := mgr.Edit("nonexistent", TaskPatch{})
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestManager_SubtaskLifecycle(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, _, _ := mgr.Create("subtask-test", CreateOptions{})
	sub1, err := mgr.AddSubtask(task.ID, "Step 1")
	if err != nil {
		t.Fatalf("AddSubtask failed: %v", err)
	}
	if sub1.Title != "Step 1" {
		t.Errorf("Subtask title = %s, want Step 1", sub1.Title)
	}
	mgr.AddSubtask(task.ID, "Step 2")
	mgr.DoneSubtask(task.ID, sub1.ID)
	updated, _ := mgr.Get(task.ID)
	if len(updated.Subtasks) != 2 {
		t.Fatalf("expected 2 subtasks, got %d", len(updated.Subtasks))
	}
	if !updated.Subtasks[0].Done {
		t.Error("subtask 1 should be done")
	}
	mgr.UndoneSubtask(task.ID, sub1.ID)
	updated, _ = mgr.Get(task.ID)
	if updated.Subtasks[0].Done {
		t.Error("subtask 1 should be undone")
	}
}

func TestManager_SubtaskNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("subtask-test", CreateOptions{})
	err := mgr.DoneSubtask(task.ID, "999")
	if err == nil {
		t.Fatal("expected error for non-existent subtask")
	}
}

func TestManager_RemoveContextEntry(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, _, _ := mgr.Create("ctx-test", CreateOptions{})
	mgr.AddContextEntry(task.ID, PhaseImplement, ContextEntry{Path: "spec/auth.md", Required: true})
	mgr.AddContextEntry(task.ID, PhaseImplement, ContextEntry{Path: "spec/api.md", Required: false})
	if err := mgr.RemoveContextEntry(task.ID, PhaseImplement, "spec/auth.md"); err != nil {
		t.Fatalf("RemoveContextEntry failed: %v", err)
	}
	entries, err := mgr.ListContextEntries(task.ID, PhaseImplement)
	if err != nil {
		t.Fatalf("ListContextEntries failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Path != "spec/api.md" {
		t.Errorf("expected [spec/api.md], got %v", entries)
	}
}

func TestManager_RemoveContextEntryNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("ctx-test", CreateOptions{})
	err := mgr.RemoveContextEntry(task.ID, PhaseImplement, "spec/nonexistent.md")
	if err == nil {
		t.Fatal("expected error for non-existent context entry")
	}
}

func TestManager_ListContextEntries_Empty(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("ctx-test", CreateOptions{})
	entries, err := mgr.ListContextEntries(task.ID, PhaseImplement)
	if err != nil {
		t.Fatalf("ListContextEntries failed: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestManager_ListContextEntries_InvalidPhase(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("ctx-test", CreateOptions{})
	_, err := mgr.ListContextEntries(task.ID, "invalid")
	if err == nil {
		t.Fatal("expected error for invalid phase")
	}
}

func TestManager_ListByStatus(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task1, dir1, _ := mgr.Create("task-a", CreateOptions{})
	os.WriteFile(filepath.Join(dir1, "prd.md"), []byte("# PRD\nA"), 0644)
	mgr.Start(task1.ID)
	mgr.Create("task-b", CreateOptions{})
	planningTasks, _ := mgr.ListByStatus(StatusPlanning)
	inProgressTasks, _ := mgr.ListByStatus(StatusInProgress)
	if len(planningTasks) != 1 {
		t.Errorf("expected 1 planning task, got %d", len(planningTasks))
	}
	if len(inProgressTasks) != 1 {
		t.Errorf("expected 1 in_progress task, got %d", len(inProgressTasks))
	}
}

func TestManager_ListRecent(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	mgr.Create("task-a", CreateOptions{})
	mgr.Create("task-b", CreateOptions{})
	mgr.Create("task-c", CreateOptions{})
	recent, err := mgr.ListRecent(2)
	if err != nil {
		t.Fatalf("ListRecent failed: %v", err)
	}
	if len(recent) != 2 {
		t.Errorf("expected 2 recent tasks, got %d", len(recent))
	}
	all, _ := mgr.ListRecent(10)
	if len(all) != 3 {
		t.Errorf("expected 3 tasks when n > total, got %d", len(all))
	}
}

func TestManager_AddSpec(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, _, _ := mgr.Create("spec-test", CreateOptions{})
	if err := mgr.AddSpec(task.ID, "spec/auth.md"); err != nil {
		t.Fatalf("AddSpec failed: %v", err)
	}
	specs, err := mgr.ListSpecs(task.ID)
	if err != nil {
		t.Fatalf("ListSpecs failed: %v", err)
	}
	if len(specs) != 1 || specs[0] != "spec/auth.md" {
		t.Errorf("specs = %v, want [spec/auth.md]", specs)
	}
	mgr.AddSpec(task.ID, "spec/auth.md")
	specs, _ = mgr.ListSpecs(task.ID)
	if len(specs) != 1 {
		t.Errorf("duplicate AddSpec should be idempotent, got %d specs", len(specs))
	}
}

func TestManager_ListEmpty(t *testing.T) {
	mgr := NewManager(t.TempDir())
	list, err := mgr.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{"valid task", Task{ID: "t1", Name: "test", Status: StatusPlanning}, false},
		{"missing ID", Task{Name: "test", Status: StatusPlanning}, true},
		{"missing name", Task{ID: "t1", Status: StatusPlanning}, true},
		{"invalid status", Task{ID: "t1", Name: "test", Status: "unknown"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestTask_Save(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "task.json")
	task := &Task{ID: "t1", Name: "test", Status: StatusPlanning}
	if err := task.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := LoadTask(path)
	if err != nil {
		t.Fatalf("LoadTask failed: %v", err)
	}
	if loaded.ID != task.ID || loaded.Name != task.Name {
		t.Errorf("loaded task mismatch: %+v", loaded)
	}
}

func TestTask_LoadTask_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "task.json")
	os.WriteFile(path, []byte("{invalid"), 0644)
	_, err := LoadTask(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTask_LoadTask_NonExistent(t *testing.T) {
	_, err := LoadTask(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestTask_DirName(t *testing.T) {
	task := &Task{ID: "t1", Name: "my-task"}
	dir := task.DirName()
	if !strings.Contains(dir, "my-task") {
		t.Errorf("DirName should contain task name, got: %s", dir)
	}
	if len(dir) < 10 {
		t.Errorf("DirName too short: %s", dir)
	}
}

func TestManager_AddContextEntry_CheckPhase(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, _, _ := mgr.Create("check-test", CreateOptions{})
	if err := mgr.AddContextEntry(task.ID, PhaseCheck, ContextEntry{Path: "spec/check.md"}); err != nil {
		t.Fatalf("AddContextEntry check phase failed: %v", err)
	}
	entries, err := mgr.ListContextEntries(task.ID, PhaseCheck)
	if err != nil {
		t.Fatalf("ListContextEntries check phase failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestManager_AddContextEntry_ResearchPhase(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, _, _ := mgr.Create("research-test", CreateOptions{})
	if err := mgr.AddContextEntry(task.ID, PhaseResearch, ContextEntry{Path: "spec/research.md"}); err != nil {
		t.Fatalf("AddContextEntry research phase failed: %v", err)
	}
	entries, err := mgr.ListContextEntries(task.ID, PhaseResearch)
	if err != nil {
		t.Fatalf("ListContextEntries research phase failed: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
}

func TestManager_AddContextEntry_InvalidPhase(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("invalid-phase", CreateOptions{})
	err := mgr.AddContextEntry(task.ID, "invalid", ContextEntry{Path: "spec/x.md"})
	if err == nil {
		t.Fatal("expected error for invalid phase")
	}
}

func TestManager_RemoveContextEntry_InvalidPhase(t *testing.T) {
	mgr := NewManager(t.TempDir())
	task, _, _ := mgr.Create("invalid-phase", CreateOptions{})
	err := mgr.RemoveContextEntry(task.ID, "invalid", "spec/x.md")
	if err == nil {
		t.Fatal("expected error for invalid phase")
	}
}

func TestManager_GetTaskNotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	_, err := mgr.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestManager_Validate_MissingFiles(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)
	task, dir, _ := mgr.Create("validate-test", CreateOptions{})
	os.Remove(filepath.Join(dir, "prd.md"))
	err := mgr.Validate(task.ID)
	if err == nil {
		t.Fatal("expected error for missing prd.md")
	}
	if !strings.Contains(err.Error(), "prd.md") {
		t.Errorf("error should mention missing file, got: %v", err)
	}
}

func TestManager_ListSpecs_NotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	_, err := mgr.ListSpecs("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestManager_AddSpec_NotFound(t *testing.T) {
	mgr := NewManager(t.TempDir())
	err := mgr.AddSpec("nonexistent", "spec/x.md")
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}
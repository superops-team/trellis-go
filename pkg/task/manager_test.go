package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestManager_Start(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("add-auth", CreateOptions{})
	if err := mgr.Start(task.ID); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	loaded, _ := mgr.Get(task.ID)
	if loaded.Status != StatusInProgress {
		t.Errorf("expected status in_progress, got %s", loaded.Status)
	}
}

func TestManager_Start_InvalidTransition(t *testing.T) {
	tmp := t.TempDir()
	mgr := NewManager(tmp)

	task, _, _ := mgr.Create("add-auth", CreateOptions{})
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
	if err := mgr.AddContext(task.ID, PhaseImplement, entry); err != nil {
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

	err := mgr.AddContext(task.ID, PhaseImplement, ContextEntry{Path: "spec/auth.md", Required: true})
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

	if err := mgr.AddContext("legacy-task", PhaseImplement, ContextEntry{Path: "spec/api.md", Description: "API spec"}); err != nil {
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

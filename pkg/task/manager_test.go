package task

import (
	"os"
	"path/filepath"
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

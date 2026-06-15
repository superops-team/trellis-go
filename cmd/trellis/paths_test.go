package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePaths(t *testing.T) {
	tests := []struct {
		name           string
		rootFlag       string
		cwd            string
		wantRepoRoot   string
		wantTrellisDir string
		wantTasksDir   string
		wantSpecDir    string
	}{
		{
			name:           "empty root uses cwd as repo root",
			cwd:            filepath.Join("tmp", "repo"),
			wantRepoRoot:   filepath.Join("tmp", "repo"),
			wantTrellisDir: filepath.Join("tmp", "repo", ".trellis"),
			wantTasksDir:   filepath.Join("tmp", "repo", ".trellis", "tasks"),
			wantSpecDir:    filepath.Join("tmp", "repo", ".trellis", "spec"),
		},
		{
			name:           "root flag can point at repo root",
			rootFlag:       filepath.Join("tmp", "repo"),
			cwd:            filepath.Join("tmp", "elsewhere"),
			wantRepoRoot:   filepath.Join("tmp", "repo"),
			wantTrellisDir: filepath.Join("tmp", "repo", ".trellis"),
			wantTasksDir:   filepath.Join("tmp", "repo", ".trellis", "tasks"),
			wantSpecDir:    filepath.Join("tmp", "repo", ".trellis", "spec"),
		},
		{
			name:           "root flag can point at trellis dir",
			rootFlag:       filepath.Join("tmp", "repo", ".trellis"),
			cwd:            filepath.Join("tmp", "elsewhere"),
			wantRepoRoot:   filepath.Join("tmp", "repo"),
			wantTrellisDir: filepath.Join("tmp", "repo", ".trellis"),
			wantTasksDir:   filepath.Join("tmp", "repo", ".trellis", "tasks"),
			wantSpecDir:    filepath.Join("tmp", "repo", ".trellis", "spec"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePaths(tt.rootFlag, func() (string, error) { return tt.cwd, nil })
			if err != nil {
				t.Fatalf("resolvePaths failed: %v", err)
			}
			if got.RepoRoot != tt.wantRepoRoot {
				t.Errorf("RepoRoot = %q, want %q", got.RepoRoot, tt.wantRepoRoot)
			}
			if got.TrellisDir != tt.wantTrellisDir {
				t.Errorf("TrellisDir = %q, want %q", got.TrellisDir, tt.wantTrellisDir)
			}
			if got.TasksDir != tt.wantTasksDir {
				t.Errorf("TasksDir = %q, want %q", got.TasksDir, tt.wantTasksDir)
			}
			if got.SpecDir != tt.wantSpecDir {
				t.Errorf("SpecDir = %q, want %q", got.SpecDir, tt.wantSpecDir)
			}
		})
	}
}

func TestResolvePaths_GetwdError(t *testing.T) {
	_, err := resolvePaths("", func() (string, error) { return "", errors.New("boom") })
	if err == nil {
		t.Fatal("expected getwd error")
	}
	if got := err.Error(); got != "get working directory: boom" {
		t.Fatalf("error = %q", got)
	}
}

func TestValidateContextPathRejectsUnsafePortablePaths(t *testing.T) {
	tests := []string{
		"../secret.txt",
		`..\secret.txt`,
		"/tmp/secret.txt",
		"C:/Users/alice/secret.txt",
		`C:\Users\alice\secret.txt`,
		"//server/share/secret.txt",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := validateContextPath(input)
			if err == nil {
				t.Fatalf("expected %q to be rejected", input)
			}
		})
	}
}

func TestValidateContextPathNormalizesBackslashes(t *testing.T) {
	got, err := validateContextPath(`spec\auth.md`)
	if err != nil {
		t.Fatalf("validateContextPath failed: %v", err)
	}
	if got != "spec/auth.md" {
		t.Fatalf("normalized path = %q, want spec/auth.md", got)
	}
}

func TestTaskDirForIDUsesActualDirectory(t *testing.T) {
	tasksDir := t.TempDir()
	taskDir := filepath.Join(tasksDir, "03-04-renamed-dir")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
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
		t.Fatalf("write task.json: %v", err)
	}

	got, err := taskDirForID(tasksDir, "legacy-task")
	if err != nil {
		t.Fatalf("taskDirForID failed: %v", err)
	}
	if got != taskDir {
		t.Fatalf("taskDirForID = %q, want actual dir %q", got, taskDir)
	}
}

func TestTaskCommandUsageMarksRequiredID(t *testing.T) {
	usage := newTaskCmd().Commands()
	seen := map[string]string{}
	for _, cmd := range usage {
		seen[cmd.Name()] = cmd.Use
	}
	for name, use := range map[string]string{"start": "start <id>", "archive": "archive <id>"} {
		if seen[name] != use {
			t.Fatalf("%s Use = %q, want %q; all uses: %s", name, seen[name], use, strings.Join([]string{seen["start"], seen["archive"]}, ", "))
		}
	}
}

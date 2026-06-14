package main

import (
	"errors"
	"path/filepath"
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

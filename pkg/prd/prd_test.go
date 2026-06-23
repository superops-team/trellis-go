package prd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRequired_ValidPRD(t *testing.T) {
	dir := t.TempDir()
	content := "# PRD\nBuild authentication system."
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte(content), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}

	got, err := LoadRequired(dir)
	if err != nil {
		t.Fatalf("LoadRequired failed: %v", err)
	}
	if got != content {
		t.Errorf("LoadRequired = %q, want %q", got, content)
	}
}

func TestLoadRequired_MissingPRD(t *testing.T) {
	dir := t.TempDir()

	_, err := LoadRequired(dir)
	if err == nil {
		t.Fatal("expected error for missing PRD")
	}
	if !strings.Contains(err.Error(), "PRD is required") {
		t.Errorf("error should mention PRD required, got: %v", err)
	}
	if err != ErrRequired {
		t.Errorf("error should be ErrRequired, got: %v", err)
	}
}

func TestLoadRequired_EmptyPRD(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte(""), 0644); err != nil {
		t.Fatalf("write empty prd: %v", err)
	}

	_, err := LoadRequired(dir)
	if err == nil {
		t.Fatal("expected error for empty PRD")
	}
	if err != ErrRequired {
		t.Errorf("error should be ErrRequired, got: %v", err)
	}
}

func TestLoadRequired_WhitespaceOnlyPRD(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte(" \n\t\n  "), 0644); err != nil {
		t.Fatalf("write whitespace prd: %v", err)
	}

	_, err := LoadRequired(dir)
	if err == nil {
		t.Fatal("expected error for whitespace-only PRD")
	}
	if err != ErrRequired {
		t.Errorf("error should be ErrRequired, got: %v", err)
	}
}

func TestLoadRequired_ReadError(t *testing.T) {
	dir := t.TempDir()
	// Create a directory with the same name as prd.md to cause a read error
	if err := os.MkdirAll(filepath.Join(dir, "prd.md"), 0755); err != nil {
		t.Fatalf("create dir as prd.md: %v", err)
	}

	_, err := LoadRequired(dir)
	if err == nil {
		t.Fatal("expected error for directory-instead-of-file")
	}
	if strings.Contains(err.Error(), "PRD is required") {
		t.Errorf("read error should not be ErrRequired, got: %v", err)
	}
	if strings.Contains(err.Error(), "read PRD") {
		// Good - wraps the read error
	} else {
		t.Errorf("error should wrap read PRD context, got: %v", err)
	}
}

func TestLoadRequired_MultilinePRD(t *testing.T) {
	dir := t.TempDir()
	content := "# PRD\n\n## Overview\nBuild auth.\n\n## Requirements\n- Login\n- Logout\n"
	if err := os.WriteFile(filepath.Join(dir, "prd.md"), []byte(content), 0644); err != nil {
		t.Fatalf("write multiline prd: %v", err)
	}

	got, err := LoadRequired(dir)
	if err != nil {
		t.Fatalf("LoadRequired failed: %v", err)
	}
	if got != content {
		t.Errorf("LoadRequired = %q, want %q", got, content)
	}
}

func TestLoadRequired_NonExistentDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")

	_, err := LoadRequired(dir)
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
	if err != ErrRequired {
		t.Errorf("error should be ErrRequired for missing file, got: %v", err)
	}
}

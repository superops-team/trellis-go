package update

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestSyncer_Sync_NewFiles(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md":     {Data: []byte("# Workflow\n[workflow-state:PLAN]\n")},
		"spec/template.md": {Data: []byte("# Spec Template\n")},
	}

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d: %v", len(result.Added), result.Added)
	}

	// Verify files were written
	for _, path := range []string{"workflow.md", "spec/template.md"} {
		targetPath := filepath.Join(targetDir, path)
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", targetPath)
		}
	}
}

func TestSyncer_Sync_SkipUnchanged(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md": {Data: []byte("# Workflow\n")},
	}

	// Pre-create identical file
	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "workflow.md"), []byte("# Workflow\n"), 0644)

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Added) != 0 || len(result.Updated) != 0 {
		t.Errorf("expected no changes, got added=%v updated=%v", result.Added, result.Updated)
	}
}

func TestSyncer_Sync_SkipUserModified(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md": {Data: []byte("# Workflow\n[workflow-state:PLAN]\n")},
	}

	// Pre-create modified file
	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "workflow.md"), []byte("# My Custom Workflow\n"), 0644)

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d: %v", len(result.Skipped), result.Skipped)
	}
}

func TestSyncer_Migrate_Overwrites(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md": {Data: []byte("# Workflow\n[workflow-state:PLAN]\n")},
	}

	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "workflow.md"), []byte("# My Custom Workflow\n"), 0644)

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
	}

	result, err := s.Migrate()
	if err != nil {
		t.Fatalf("Migrate() error: %v", err)
	}

	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d: %v", len(result.Updated), result.Updated)
	}

	// Verify overwritten
	data, _ := os.ReadFile(filepath.Join(targetDir, "workflow.md"))
	if string(data) != "# Workflow\n[workflow-state:PLAN]\n" {
		t.Errorf("expected migrated content, got %q", string(data))
	}
}

func TestSyncer_DryRun_NoWrite(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md": {Data: []byte("# Workflow\n")},
	}

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
	}

	result, err := s.DryRun()
	if err != nil {
		t.Fatalf("DryRun() error: %v", err)
	}

	if len(result.Added) != 1 {
		t.Errorf("expected 1 added in dry-run, got %d", len(result.Added))
	}

	// Verify NOT written
	if _, err := os.Stat(filepath.Join(targetDir, "workflow.md")); !os.IsNotExist(err) {
		t.Error("expected file NOT to be written in dry-run")
	}
}

func TestSyncer_SkipPaths(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{
		"workflow.md":     {Data: []byte("# Workflow\n")},
		"spec/template.md": {Data: []byte("# Template\n")},
	}

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
		Config: SyncerConfig{
			Skip: []string{"spec/template.md"},
		},
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Skipped) != 1 || result.Skipped[0] != "spec/template.md" {
		t.Errorf("expected spec/template.md skipped, got %v", result.Skipped)
	}
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added))
	}
}

func TestSyncer_ConfigSections(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{}

	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "config.yaml"), []byte("packages: []\n"), 0644)

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
		Config: SyncerConfig{
			Sections: []ConfigSection{
				{Sentinel: "update:", Content: "update:\n  skip: []"},
			},
		},
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Sections) != 1 {
		t.Errorf("expected 1 section appended, got %d", len(result.Sections))
	}

	data, _ := os.ReadFile(filepath.Join(targetDir, "config.yaml"))
	if string(data) != "packages: []\n\nupdate:\n  skip: []\n" {
		t.Errorf("unexpected config content: %q", string(data))
	}
}

func TestSyncer_ConfigSections_Idempotent(t *testing.T) {
	dir := t.TempDir()
	targetDir := filepath.Join(dir, ".trellis")

	embedFS := fstest.MapFS{}

	os.MkdirAll(targetDir, 0755)
	os.WriteFile(filepath.Join(targetDir, "config.yaml"), []byte("packages: []\nupdate:\n  skip: []\n"), 0644)

	s := &Syncer{
		EmbedFS:   embedFS,
		TargetDir: targetDir,
		Config: SyncerConfig{
			Sections: []ConfigSection{
				{Sentinel: "update:", Content: "update:\n  skip: []"},
			},
		},
	}

	result, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error: %v", err)
	}

	if len(result.Sections) != 0 {
		t.Errorf("expected 0 sections (already present), got %d", len(result.Sections))
	}
}

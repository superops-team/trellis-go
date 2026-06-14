package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadMissingManifestReturnsEmptyManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.jsonl")

	manifest, err := Load(path)
	if err != nil {
		t.Fatalf("Load missing manifest failed: %v", err)
	}
	if manifest.Version != "1.0" {
		t.Fatalf("Version = %q, want 1.0", manifest.Version)
	}
	if len(manifest.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(manifest.Entries))
	}
}

func TestSaveAndLoadManifestPreservesJSONLFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.jsonl")
	want := &Manifest{
		Version: "1.0",
		Entries: []Entry{
			{Path: "spec/auth.md", Description: "Auth spec", Required: true},
			{Path: "spec/api.md", Required: false},
		},
	}

	if err := Save(path, want); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSONL lines, got %d: %s", len(lines), data)
	}
	if !strings.Contains(lines[0], `"path":"spec/auth.md"`) || !strings.Contains(lines[0], `"description":"Auth spec"`) || !strings.Contains(lines[0], `"required":true`) {
		t.Fatalf("first line does not preserve fields: %s", lines[0])
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(got.Entries) != len(want.Entries) {
		t.Fatalf("Load entries = %d, want %d", len(got.Entries), len(want.Entries))
	}
	if got.Entries[0] != want.Entries[0] || got.Entries[1] != want.Entries[1] {
		t.Fatalf("loaded entries = %#v, want %#v", got.Entries, want.Entries)
	}
}

func TestLoadMalformedManifestReportsPathAndLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.jsonl")
	if err := os.WriteFile(path, []byte("{not-json\n"), 0644); err != nil {
		t.Fatalf("write malformed manifest: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected malformed manifest error")
	}
	if !strings.Contains(err.Error(), path) || !strings.Contains(err.Error(), "line 1") {
		t.Fatalf("error should mention path and line, got: %v", err)
	}
}

func TestManifestSaveMethod(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.jsonl")
	manifest := &Manifest{Version: "1.0", Entries: []Entry{{Path: "a.md", Required: true}}}

	if err := manifest.Save(path); err != nil {
		t.Fatalf("Manifest.Save failed: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Path != "a.md" {
		t.Fatalf("loaded entries = %#v", loaded.Entries)
	}
}

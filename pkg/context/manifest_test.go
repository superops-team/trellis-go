package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "manifest.jsonl")
	content := `{"path":"spec/auth.md","required":true}
{"path":"spec/api.md","description":"API spec","required":false}
`
	os.WriteFile(path, []byte(content), 0644)

	m, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if len(m.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m.Entries))
	}
	if m.Entries[0].Path != "spec/auth.md" {
		t.Errorf("expected path spec/auth.md, got %s", m.Entries[0].Path)
	}
	if !m.Entries[0].Required {
		t.Error("expected first entry to be required")
	}
}

func TestLoadManifest_NotFound(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "missing.jsonl")

	m, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if len(m.Entries) != 0 {
		t.Errorf("expected 0 entries for missing file, got %d", len(m.Entries))
	}
}

func TestManifest_Save(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "manifest.jsonl")

	m := &Manifest{
		Version: "1.0",
		Entries: []Entry{
			{Path: "a.md", Required: true},
			{Path: "b.md", Required: false},
		},
	}
	if err := m.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, _ := os.ReadFile(path)
	lines := string(data)
	if lines == "" {
		t.Error("expected non-empty file")
	}
}

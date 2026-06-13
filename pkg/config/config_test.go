package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	content := `
packages:
  - frontend
  - backend
codex:
  dispatch_mode: inline
developer: alice
`
	os.WriteFile(path, []byte(content), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(cfg.Packages) != 2 {
		t.Errorf("expected 2 packages, got %d", len(cfg.Packages))
	}
	if cfg.Codex.DispatchMode != "inline" {
		t.Errorf("expected inline, got %s", cfg.Codex.DispatchMode)
	}
	if cfg.Developer != "alice" {
		t.Errorf("expected alice, got %s", cfg.Developer)
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	cfg := &Config{
		Packages:  []string{"api", "web"},
		Developer: "bob",
		Codex:     CodexConfig{DispatchMode: "sub-agent"},
	}
	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Packages) != 2 || loaded.Packages[0] != "api" {
		t.Errorf("packages mismatch: %v", loaded.Packages)
	}
	if loaded.Codex.DispatchMode != "sub-agent" {
		t.Errorf("dispatch_mode mismatch: %s", loaded.Codex.DispatchMode)
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid",
			cfg: Config{
				Packages: []string{"pkg1"},
				Codex:    CodexConfig{DispatchMode: "inline"},
			},
			wantErr: false,
		},
		{
			name:    "empty package",
			cfg:     Config{Packages: []string{"pkg1", ""}},
			wantErr: true,
		},
		{
			name:    "invalid dispatch_mode",
			cfg:     Config{Codex: CodexConfig{DispatchMode: "invalid"}},
			wantErr: true,
		},
		{
			name:    "empty dispatch_mode allowed",
			cfg:     Config{},
			wantErr: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/spec"
)

func TestBuilder_BuildImplementContextRequiredEntryMissing(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	manifest := `{"path":"spec/missing.md","description":"Missing spec","required":true}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write implement manifest: %v", err)
	}

	builder := &Builder{Root: root}
	_, err := builder.BuildImplementContext(taskDir)
	if err == nil {
		t.Fatal("expected missing required context entry to fail")
	}
	if !strings.Contains(err.Error(), "required entry spec/missing.md") {
		t.Errorf("error should identify missing required entry, got: %v", err)
	}
}

func TestBuilder_BuildImplementContextOptionalEntryMissing(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	manifest := `{"path":"spec/missing.md","description":"Missing spec","required":false}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write implement manifest: %v", err)
	}

	builder := &Builder{Root: root}
	output, err := builder.BuildImplementContext(taskDir)
	if err != nil {
		t.Fatalf("optional missing context entry should not fail: %v", err)
	}
	if !strings.Contains(output, injectMarker) {
		t.Errorf("output should contain injection marker, got: %s", output)
	}
	if !strings.Contains(output, "=== skipped optional context ===") || !strings.Contains(output, "spec/missing.md") {
		t.Errorf("output should list missing optional entry, got: %s", output)
	}
}

func TestBuilder_BuildImplementContextRequiresNonEmptyPRD(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte(" \n\t"), 0644); err != nil {
		t.Fatalf("write blank prd: %v", err)
	}

	builder := &Builder{Root: root}
	_, err := builder.BuildImplementContext(taskDir)
	if err == nil {
		t.Fatal("expected blank PRD to fail")
	}
	if !strings.Contains(err.Error(), "PRD is required") {
		t.Fatalf("error should explain PRD requirement, got: %v", err)
	}
}

func TestBuilder_BuildCheckContextRequiresNonEmptyPRD(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte(" \n\t"), 0644); err != nil {
		t.Fatalf("write blank prd: %v", err)
	}

	builder := &Builder{Root: root}
	_, err := builder.BuildCheckContext(taskDir)
	if err == nil {
		t.Fatal("expected blank PRD to fail")
	}
	if !strings.Contains(err.Error(), "PRD is required") {
		t.Fatalf("error should explain PRD requirement, got: %v", err)
	}
}

func TestBuilder_BuildImplementContextRejectsUnsafeManifestPaths(t *testing.T) {
	tests := []struct {
		name      string
		entryPath string
	}{
		{name: "absolute", entryPath: "/etc/passwd"},
		{name: "parent traversal", entryPath: "../secret.txt"},
		{name: "backslash parent traversal", entryPath: `..\secret.txt`},
		{name: "windows drive", entryPath: `C:\Users\alice\secret.txt`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
			if err := os.MkdirAll(taskDir, 0755); err != nil {
				t.Fatalf("create task dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
				t.Fatalf("write prd: %v", err)
			}
			manifest := `{"path":"` + strings.ReplaceAll(tt.entryPath, `\`, `\\`) + `","description":"Unsafe","required":true}` + "\n"
			if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
				t.Fatalf("write manifest: %v", err)
			}

			builder := &Builder{Root: root}
			_, err := builder.BuildImplementContext(taskDir)
			if err == nil {
				t.Fatal("expected unsafe manifest path to fail")
			}
			if !strings.Contains(err.Error(), "invalid context path") {
				t.Fatalf("error should identify invalid context path, got: %v", err)
			}
		})
	}
}

func TestBuilder_BuildCheckContextIncludesPRDAndReferencedFiles(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(filepath.Join(root, "spec"), 0755); err != nil {
		t.Fatalf("create spec dir: %v", err)
	}
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "spec", "auth-check.md"), []byte("# Auth Check\nVerify JWT."), 0644); err != nil {
		t.Fatalf("write check spec: %v", err)
	}
	manifest := `{"path":"spec/auth-check.md","description":"Auth check","required":true}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "check.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write check manifest: %v", err)
	}

	builder := &Builder{Root: root}
	output, err := builder.BuildCheckContext(taskDir)
	if err != nil {
		t.Fatalf("BuildCheckContext failed: %v", err)
	}
	for _, want := range []string{injectMarker, "# PRD\nBuild auth.", "# Auth Check\nVerify JWT."} {
		if !strings.Contains(output, want) {
			t.Errorf("check context should contain %q, got: %s", want, output)
		}
	}
}

func TestBuilder_BuildCheckContextRequiredEntryMissing(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	manifest := `{"path":"spec/missing-check.md","description":"Missing check","required":true}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "check.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write check manifest: %v", err)
	}

	builder := &Builder{Root: root}
	_, err := builder.BuildCheckContext(taskDir)
	if err == nil {
		t.Fatal("expected missing required check entry to fail")
	}
	if !strings.Contains(err.Error(), "required entry spec/missing-check.md") {
		t.Errorf("error should identify missing required check entry, got: %v", err)
	}
}

func TestBuilder_BuildCheckContextOptionalEntryMissing(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	manifest := `{"path":"spec/missing-check.md","description":"Missing check","required":false}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "check.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write check manifest: %v", err)
	}

	builder := &Builder{Root: root}
	output, err := builder.BuildCheckContext(taskDir)
	if err != nil {
		t.Fatalf("optional missing check entry should not fail: %v", err)
	}
	if !strings.Contains(output, injectMarker) {
		t.Errorf("output should contain injection marker, got: %s", output)
	}
	if !strings.Contains(output, "=== skipped optional context ===") || !strings.Contains(output, "spec/missing-check.md") {
		t.Errorf("output should list missing optional check entry, got: %s", output)
	}
}

func TestBuilder_BuildImplementContextRejectsRequiredUnsafeEntries(t *testing.T) {
	tests := []struct {
		name      string
		entryPath string
		content   []byte
		wantError string
	}{
		{
			name:      "oversized",
			entryPath: "spec/large.md",
			content:   []byte(strings.Repeat("a", 256*1024+1)),
			wantError: "too large",
		},
		{
			name:      "binary",
			entryPath: "spec/data.bin",
			content:   []byte{0, 1, 2},
			wantError: "binary file",
		},
		{
			name:      "sensitive",
			entryPath: "spec/.env",
			content:   []byte("TOKEN=value"),
			wantError: "sensitive path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
			if err := os.MkdirAll(filepath.Join(root, "spec"), 0755); err != nil {
				t.Fatalf("create spec dir: %v", err)
			}
			if err := os.MkdirAll(taskDir, 0755); err != nil {
				t.Fatalf("create task dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
				t.Fatalf("write prd: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(tt.entryPath)), tt.content, 0644); err != nil {
				t.Fatalf("write context entry: %v", err)
			}
			manifest := `{"path":"` + tt.entryPath + `","description":"unsafe","required":true}` + "\n"
			if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
				t.Fatalf("write manifest: %v", err)
			}

			_, err := (&Builder{Root: root}).BuildImplementContext(taskDir)
			if err == nil {
				t.Fatal("expected unsafe required entry to fail")
			}
			if !strings.Contains(err.Error(), tt.entryPath) || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error should mention path %q and reason %q, got: %v", tt.entryPath, tt.wantError, err)
			}
		})
	}
}

func TestBuilder_BuildImplementContextListsOptionalUnsafeEntries(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(filepath.Join(root, "spec"), 0755); err != nil {
		t.Fatalf("create spec dir: %v", err)
	}
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "spec", "data.bin"), []byte{0, 1, 2}, 0644); err != nil {
		t.Fatalf("write binary entry: %v", err)
	}
	manifest := `{"path":"spec/data.bin","description":"binary","required":false}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	output, err := (&Builder{Root: root}).BuildImplementContext(taskDir)
	if err != nil {
		t.Fatalf("optional unsafe entry should not fail: %v", err)
	}
	for _, want := range []string{"=== skipped optional context ===", "spec/data.bin", "binary file"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output should contain %q, got: %s", want, output)
		}
	}
}

func TestBuilder_BuildImplementContextListsOptionalOversizedAndSensitiveEntries(t *testing.T) {
	tests := []struct {
		name      string
		entryPath string
		content   []byte
		wantError string
	}{
		{
			name:      "oversized optional",
			entryPath: "spec/large.md",
			content:   []byte(strings.Repeat("a", 256*1024+1)),
			wantError: "too large",
		},
		{
			name:      "sensitive optional",
			entryPath: "spec/.env",
			content:   []byte("TOKEN=value"),
			wantError: "sensitive path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
			if err := os.MkdirAll(filepath.Join(root, "spec"), 0755); err != nil {
				t.Fatalf("create spec dir: %v", err)
			}
			if err := os.MkdirAll(taskDir, 0755); err != nil {
				t.Fatalf("create task dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
				t.Fatalf("write prd: %v", err)
			}
			if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(tt.entryPath)), tt.content, 0644); err != nil {
				t.Fatalf("write context entry: %v", err)
			}
			manifest := `{"path":"` + tt.entryPath + `","description":"unsafe","required":false}` + "\n"
			if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
				t.Fatalf("write manifest: %v", err)
			}

			output, err := (&Builder{Root: root}).BuildImplementContext(taskDir)
			if err != nil {
				t.Fatalf("optional unsafe entry should not fail: %v", err)
			}
			for _, want := range []string{"=== skipped optional context ===", tt.entryPath, tt.wantError} {
				if !strings.Contains(output, want) {
					t.Fatalf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestBuilder_BuildImplementContextOmitsSkippedSectionWhenNothingSkipped(t *testing.T) {
	root := t.TempDir()
	taskDir := filepath.Join(root, "tasks", "01-01-user-auth")
	if err := os.MkdirAll(filepath.Join(root, "spec"), 0755); err != nil {
		t.Fatalf("create spec dir: %v", err)
	}
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		t.Fatalf("create task dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskDir, "prd.md"), []byte("# PRD\nBuild auth."), 0644); err != nil {
		t.Fatalf("write prd: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "spec", "auth.md"), []byte("# Auth Spec\nUse JWT."), 0644); err != nil {
		t.Fatalf("write spec: %v", err)
	}
	manifest := `{"path":"spec/auth.md","description":"auth","required":false}` + "\n"
	if err := os.WriteFile(filepath.Join(taskDir, "implement.jsonl"), []byte(manifest), 0644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	output, err := (&Builder{Root: root}).BuildImplementContext(taskDir)
	if err != nil {
		t.Fatalf("BuildImplementContext failed: %v", err)
	}
	if strings.Contains(output, "=== skipped optional context ===") {
		t.Fatalf("output should not contain skipped optional context section, got: %s", output)
	}
}

func TestBuilder_BuildResearchContextIncludesSpecIndex(t *testing.T) {
	root := t.TempDir()
	layerDir := filepath.Join(root, "auth", "api")
	if err := os.MkdirAll(layerDir, 0755); err != nil {
		t.Fatalf("create spec layer dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(layerDir, "index.md"), []byte("# Auth API"), 0644); err != nil {
		t.Fatalf("write layer index: %v", err)
	}

	builder := &Builder{SpecLoader: spec.NewLoader(root)}
	output, err := builder.BuildResearchContext()
	if err != nil {
		t.Fatalf("BuildResearchContext failed: %v", err)
	}
	for _, want := range []string{injectMarker, "# Spec Index", "### auth", "**api**"} {
		if !strings.Contains(output, want) {
			t.Errorf("research context should contain %q, got: %s", want, output)
		}
	}
}

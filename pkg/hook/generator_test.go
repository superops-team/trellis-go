package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/superops-team/trellis-go/pkg/platform"
)

func TestGenerator_GenerateAllDispatchesByPlatformClass(t *testing.T) {
	tests := []struct {
		name      string
		class     platform.Class
		wantFile  string
		wantBytes string
	}{
		{
			name:      "pull based creates agent definition",
			class:     platform.ClassPullBased,
			wantFile:  "trellis-implement.toml",
			wantBytes: "Do NOT spawn another sub-agent",
		},
		{
			name:      "agentless creates before-dev skill",
			class:     platform.ClassAgentless,
			wantFile:  "trellis-before-dev.md",
			wantBytes: "Load project specs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			p := platform.Platform{ID: "test", Name: "Test", ConfigDir: ".test", Class: tt.class}
			g := NewGenerator(p, "trellis")

			if err := g.GenerateAll(tmp); err != nil {
				t.Fatalf("GenerateAll failed: %v", err)
			}

			data, err := os.ReadFile(filepath.Join(tmp, tt.wantFile))
			if err != nil {
				t.Fatalf("expected generated file %s: %v", tt.wantFile, err)
			}
			if !strings.Contains(string(data), tt.wantBytes) {
				t.Errorf("generated file should contain %q, got: %s", tt.wantBytes, data)
			}
		})
	}
}

func TestGenerator_GenerateAllUnknownClassReturnsError(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "mystery", Name: "Mystery", ConfigDir: ".mystery", Class: platform.Class("mystery")}
	g := NewGenerator(p, "trellis")

	err := g.GenerateAll(tmp)
	if err == nil {
		t.Fatal("expected unknown platform class error")
	}
	if !strings.Contains(err.Error(), "unknown platform class") {
		t.Errorf("error should mention unknown platform class, got: %v", err)
	}
}

func TestGenerator_GenerateAgentDefRejectsNonPullBasedPlatform(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	g := NewGenerator(p, "trellis")

	err := g.GenerateAgentDef(tmp)
	if err == nil {
		t.Fatal("expected non-pull-based platform to be rejected")
	}
	if !strings.Contains(err.Error(), "agent defs only for pull-based platforms") {
		t.Errorf("error should explain pull-based restriction, got: %v", err)
	}
}

func TestGenerator_GenerateAgentDef(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "codex", Name: "Codex", ConfigDir: ".codex", Class: platform.ClassPullBased}
	g := NewGenerator(p, "trellis")

	if err := g.GenerateAgentDef(tmp); err != nil {
		t.Fatalf("GenerateAgentDef failed: %v", err)
	}

	path := filepath.Join(tmp, "trellis-implement.toml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read agent def: %v", err)
	}
	content := string(data)
	if !contains(content, "multi_agent = false") {
		t.Error("agent def should disable multi_agent")
	}
	if !contains(content, "Do NOT spawn another sub-agent") {
		t.Error("agent def should contain anti-recursion warning")
	}
}

func TestGenerator_GenerateBeforeDevSkill(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "windsurf", Name: "Windsurf", ConfigDir: ".windsurf", Class: platform.ClassAgentless}
	g := NewGenerator(p, "trellis")

	if err := g.GenerateBeforeDevSkill(tmp); err != nil {
		t.Fatalf("GenerateBeforeDevSkill failed: %v", err)
	}

	path := filepath.Join(tmp, "trellis-before-dev.md")
	if _, err := os.Stat(path); err != nil {
		t.Error("before-dev skill file not created")
	}
}

func TestGenerator_GenerateSessionStart(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	g := NewGenerator(p, "trellis")

	if err := g.GenerateSessionStart(tmp); err != nil {
		t.Fatalf("GenerateSessionStart failed: %v", err)
	}

	path := filepath.Join(tmp, "session-start.sh")
	data, _ := os.ReadFile(path)
	if !contains(string(data), "trellis hook session-start") {
		t.Error("session start script should call trellis hook")
	}
}

func TestGenerator_GenerateInjectHookScripts(t *testing.T) {
	tests := []struct {
		name      string
		generate  func(*Generator, string) error
		wantFile  string
		wantBytes string
	}{
		{
			name:      "inject context",
			generate:  (*Generator).GenerateInjectContext,
			wantFile:  "inject-context.sh",
			wantBytes: "trellis hook inject-context",
		},
		{
			name:      "inject workflow state",
			generate:  (*Generator).GenerateInjectWorkflowState,
			wantFile:  "inject-workflow-state.sh",
			wantBytes: "trellis hook inject-workflow-state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir()
			p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
			g := NewGenerator(p, "trellis")

			if err := tt.generate(g, tmp); err != nil {
				t.Fatalf("generate failed: %v", err)
			}

			path := filepath.Join(tmp, tt.wantFile)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read generated hook script: %v", err)
			}
			if !strings.Contains(string(data), tt.wantBytes) {
				t.Errorf("hook script should call %q, got: %s", tt.wantBytes, data)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("stat generated hook script: %v", err)
			}
			if info.Mode().Perm() != 0755 {
				t.Errorf("hook script mode = %v, want 0755", info.Mode().Perm())
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

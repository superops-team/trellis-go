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
		wantFiles []string
		wantBytes string
	}{
		{
			name:      "push based creates executable hook scripts",
			class:     platform.ClassPushBased,
			wantFiles: []string{"session-start.sh", "inject-context.sh", "inject-workflow-state.sh"},
			wantBytes: "trellis hook",
		},
		{
			name:      "pull based creates agent definition",
			class:     platform.ClassPullBased,
			wantFiles: []string{"trellis-implement.toml"},
			wantBytes: "Do NOT spawn another sub-agent",
		},
		{
			name:      "agentless creates before-dev skill",
			class:     platform.ClassAgentless,
			wantFiles: []string{"trellis-before-dev.md"},
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

			for _, wantFile := range tt.wantFiles {
				path := filepath.Join(tmp, wantFile)
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("expected generated file %s: %v", wantFile, err)
				}
				if !strings.Contains(string(data), tt.wantBytes) {
					t.Errorf("generated file should contain %q, got: %s", tt.wantBytes, data)
				}
				if tt.class == platform.ClassPushBased {
					info, err := os.Stat(path)
					if err != nil {
						t.Fatalf("stat generated file %s: %v", wantFile, err)
					}
					if info.Mode().Perm() != 0755 {
						t.Fatalf("push hook script %s mode = %v, want 0755", wantFile, info.Mode().Perm())
					}
				}
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
	if !strings.Contains(content, "multi_agent = false") {
		t.Error("agent def should disable multi_agent")
	}
	if !strings.Contains(content, "Do NOT spawn another sub-agent") {
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
	if !strings.Contains(string(data), "trellis hook session-start") {
		t.Error("session start script should call trellis hook")
	}
}

func TestGenerator_HookScriptQuotesBinaryAndForwardsArgs(t *testing.T) {
	tmp := t.TempDir()
	p := platform.Platform{ID: "claude", Name: "Claude", ConfigDir: ".claude", Class: platform.ClassPushBased}
	g := NewGenerator(p, "/tmp/Trellis Bin/trellis")

	if err := g.GenerateInjectContext(tmp); err != nil {
		t.Fatalf("GenerateInjectContext failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "inject-context.sh"))
	if err != nil {
		t.Fatalf("read generated hook script: %v", err)
	}
	want := "exec '/tmp/Trellis Bin/trellis' hook inject-context \"$@\""
	if !strings.Contains(string(data), want) {
		t.Fatalf("hook script should quote binary and forward args with %q, got: %s", want, data)
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

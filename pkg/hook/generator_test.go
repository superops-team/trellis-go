package hook

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mindfold/trellis/pkg/platform"
)

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

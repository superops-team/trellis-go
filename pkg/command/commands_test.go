package command

import (
	"strings"
	"testing"
)

func TestStartCommand(t *testing.T) {
	cmd := StartCommand()
	if cmd.Name != "start" {
		t.Errorf("expected name 'start', got %q", cmd.Name)
	}
	if !strings.Contains(cmd.Content, "trellis hook session-start") {
		t.Error("start command should mention session-start")
	}
}

func TestFinishWorkCommand(t *testing.T) {
	cmd := FinishWorkCommand()
	if cmd.Name != "finish-work" {
		t.Errorf("expected name 'finish-work', got %q", cmd.Name)
	}
	if !strings.Contains(cmd.Content, "trellis task archive") {
		t.Error("finish-work should mention task archive")
	}
	if !strings.Contains(cmd.Content, "trellis hook record-session") {
		t.Error("finish-work should mention record-session")
	}
}

func TestContinueCommand(t *testing.T) {
	cmd := ContinueCommand()
	if cmd.Name != "continue" {
		t.Errorf("expected name 'continue', got %q", cmd.Name)
	}
	if !strings.Contains(cmd.Content, "trellis hook inject-workflow-state") {
		t.Error("continue should mention inject-workflow-state")
	}
}

func TestAllCommands(t *testing.T) {
	cmds := AllCommands()
	if len(cmds) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(cmds))
	}
	names := make(map[string]bool)
	for _, c := range cmds {
		names[c.Name] = true
	}
	for _, n := range []string{"start", "finish-work", "continue"} {
		if !names[n] {
			t.Errorf("missing command: %s", n)
		}
	}
}

func TestFormatForClaudeCode(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForClaudeCode(cmd)
	if filename != "trellis/start.md" {
		t.Errorf("expected 'trellis/start.md', got %q", filename)
	}
	if !strings.HasPrefix(content, "---\nname: trellis:start") {
		t.Error("Claude Code format should start with YAML frontmatter")
	}
}

func TestFormatForCursor(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForCursor(cmd)
	if filename != "trellis-start.md" {
		t.Errorf("expected 'trellis-start.md', got %q", filename)
	}
	if !strings.HasPrefix(content, "---\nname: trellis-start") {
		t.Error("Cursor format should start with YAML frontmatter")
	}
}

func TestFormatForCodex(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForCodex(cmd)
	if filename != "trellis-start.md" {
		t.Errorf("expected 'trellis-start.md', got %q", filename)
	}
	if !strings.HasPrefix(content, "# /trellis:start") {
		t.Error("Codex format should start with heading")
	}
}

func TestFormatForGemini(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForGemini(cmd)
	if filename != "trellis/start.toml" {
		t.Errorf("expected 'trellis/start.toml', got %q", filename)
	}
	if !strings.HasPrefix(content, "[command]") {
		t.Error("Gemini format should start with [command]")
	}
}

func TestFormatForQoder(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForQoder(cmd)
	if filename != "trellis-start.md" {
		t.Errorf("expected 'trellis-start.md', got %q", filename)
	}
	if !strings.HasPrefix(content, "---\nname: trellis-start") {
		t.Error("Qoder format should start with YAML frontmatter")
	}
}

func TestFormatForCopilot(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForCopilot(cmd)
	if filename != "trellis-start.prompt.md" {
		t.Errorf("expected 'trellis-start.prompt.md', got %q", filename)
	}
	if !strings.HasPrefix(content, "# /trellis:start") {
		t.Error("Copilot format should start with heading")
	}
}

func TestFormatForWorkflow(t *testing.T) {
	cmd := StartCommand()
	filename, content := FormatForWorkflow(cmd)
	if filename != "trellis-start.md" {
		t.Errorf("expected 'trellis-start.md', got %q", filename)
	}
	if !strings.Contains(content, "Trellis Start Workflow") {
		t.Errorf("expected workflow title, got: %s", content[:50])
	}
}

func TestPlatformFormats(t *testing.T) {
	formats := PlatformFormats()
	if len(formats) != 10 {
		t.Fatalf("expected 10 platform formats, got %d", len(formats))
	}
	ids := make(map[string]bool)
	for _, f := range formats {
		ids[f.PlatformID] = true
	}
	for _, id := range []string{"claude", "cursor", "opencode", "codex", "gemini", "qoder", "codebuddy", "droid", "pi", "copilot"} {
		if !ids[id] {
			t.Errorf("missing platform: %s", id)
		}
	}
}

func TestAllCommands_ContentNotEmpty(t *testing.T) {
	for _, cmd := range AllCommands() {
		if cmd.Content == "" {
			t.Errorf("command %s has empty content", cmd.Name)
		}
	}
}

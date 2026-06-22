package agent

import (
	"strings"
	"testing"
)

func TestFormatForClaudeCode(t *testing.T) {
	a := ImplementAgent()
	fn, content := FormatForClaudeCode(a)

	if fn != "trellis-implement.md" {
		t.Errorf("expected filename 'trellis-implement.md', got %q", fn)
	}
	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
	if !strings.Contains(content, "tools: Read, Write, Edit, Bash, Glob, Grep") {
		t.Error("expected tools line in frontmatter")
	}
	if !strings.Contains(content, "# trellis-implement") {
		t.Error("expected content section")
	}
}

func TestFormatForOpenCode(t *testing.T) {
	a := ImplementAgent()
	fn, content := FormatForOpenCode(a)

	if fn != "trellis-implement.md" {
		t.Errorf("expected filename 'trellis-implement.md', got %q", fn)
	}
	if !strings.Contains(content, "permission:") {
		t.Error("expected permission object in OpenCode format")
	}
}

func TestFormatForCodex(t *testing.T) {
	a := ImplementAgent()
	fn, content := FormatForCodex(a)

	if !strings.HasSuffix(fn, ".toml") {
		t.Errorf("expected .toml extension, got %q", fn)
	}
	if !strings.Contains(content, "[features]") {
		t.Error("expected [features] section in TOML")
	}
	if !strings.Contains(content, "developer_instructions") {
		t.Error("expected developer_instructions in TOML")
	}
}

func TestFormatForKiro(t *testing.T) {
	a := ImplementAgent()
	fn, content := FormatForKiro(a)

	if !strings.HasSuffix(fn, ".json") {
		t.Errorf("expected .json extension, got %q", fn)
	}
	if !strings.Contains(content, `"name"`) {
		t.Error("expected JSON format")
	}
}

func TestFormatForGemini(t *testing.T) {
	a := ImplementAgent()
	_, content := FormatForGemini(a)

	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
	if !strings.Contains(content, "Before acting") {
		t.Error("expected pull-based prelude")
	}
}

func TestFormatForCopilot(t *testing.T) {
	a := ImplementAgent()
	fn, content := FormatForCopilot(a)

	if !strings.HasSuffix(fn, ".agent.md") {
		t.Errorf("expected .agent.md extension, got %q", fn)
	}
	if !strings.Contains(content, "Before acting") {
		t.Error("expected pull-based prelude")
	}
}

func TestFormatForPi(t *testing.T) {
	a := ImplementAgent()
	_, content := FormatForPi(a)

	if !strings.Contains(content, "model: inherit") {
		t.Error("expected model field in Pi format")
	}
	if !strings.Contains(content, "thinking: inherit") {
		t.Error("expected thinking field in Pi format")
	}
}

func TestFormatForCursor(t *testing.T) {
	a := CheckAgent()
	fn, content := FormatForCursor(a)

	if fn != "trellis-check.md" {
		t.Errorf("expected filename 'trellis-check.md', got %q", fn)
	}
	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
}

func TestFormatForQoder(t *testing.T) {
	a := ResearchAgent()
	_, content := FormatForQoder(a)

	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
	if !strings.Contains(content, "Before acting") {
		t.Error("expected pull-based prelude")
	}
}

func TestFormatForCodeBuddy(t *testing.T) {
	a := ImplementAgent()
	_, content := FormatForCodeBuddy(a)

	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
}

func TestFormatForDroid(t *testing.T) {
	a := ImplementAgent()
	_, content := FormatForDroid(a)

	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
}

func TestAllFormattersProduceContent(t *testing.T) {
	agents := AllAgents()
	formatters := []struct {
		name string
		fn   func(Agent) (string, string)
	}{
		{"ClaudeCode", FormatForClaudeCode},
		{"Cursor", FormatForCursor},
		{"OpenCode", FormatForOpenCode},
		{"Codex", FormatForCodex},
		{"Kiro", FormatForKiro},
		{"Gemini", FormatForGemini},
		{"Qoder", FormatForQoder},
		{"CodeBuddy", FormatForCodeBuddy},
		{"Copilot", FormatForCopilot},
		{"Droid", FormatForDroid},
		{"Pi", FormatForPi},
	}

	for _, a := range agents {
		for _, f := range formatters {
			fn, content := f.fn(a)
			if fn == "" {
				t.Errorf("%s/%s: empty filename", f.name, a.Name)
			}
			if content == "" {
				t.Errorf("%s/%s: empty content", f.name, a.Name)
			}
		}
	}
}

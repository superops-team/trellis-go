package skill

import (
	"strings"
	"testing"
)

func TestFormatForShared(t *testing.T) {
	s := BrainstormSkill()
	fn, content := FormatForShared(s)

	if !strings.HasSuffix(fn, "/SKILL.md") {
		t.Errorf("expected SKILL.md filename, got %q", fn)
	}
	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected YAML frontmatter")
	}
	if !strings.Contains(content, "name: trellis-brainstorm") {
		t.Error("expected name in frontmatter")
	}
	if !strings.Contains(content, "description:") {
		t.Error("expected description in frontmatter")
	}
}

func TestFormatForCursor(t *testing.T) {
	s := BrainstormSkill()
	fn, content := FormatForCursor(s)

	if !strings.HasPrefix(fn, "trellis-") {
		t.Errorf("expected trellis- prefix for Cursor filename, got %q", fn)
	}
	if !strings.Contains(content, "name: trellis-brainstorm") {
		t.Error("expected name in Cursor frontmatter")
	}
}

func TestFormatForGemini(t *testing.T) {
	s := BrainstormSkill()
	fn, _ := FormatForGemini(s)

	if !strings.HasPrefix(fn, "trellis-") {
		t.Errorf("expected trellis- prefix for Gemini, got %q", fn)
	}
}

func TestFormatForQoder(t *testing.T) {
	s := BrainstormSkill()
	fn, _ := FormatForQoder(s)

	if !strings.HasPrefix(fn, "trellis-") {
		t.Errorf("expected trellis- prefix for Qoder, got %q", fn)
	}
}

func TestPlatformFormats(t *testing.T) {
	pfs := PlatformFormats()
	if len(pfs) < 14 {
		t.Errorf("expected at least 14 platform formats, got %d", len(pfs))
	}

	ids := make(map[string]bool)
	for _, pf := range pfs {
		if ids[pf.PlatformID] {
			t.Errorf("duplicate platform ID: %s", pf.PlatformID)
		}
		ids[pf.PlatformID] = true
		if pf.Dir == "" {
			t.Errorf("platform %s: empty Dir", pf.PlatformID)
		}
		if pf.FilenameFn == nil {
			t.Errorf("platform %s: nil FilenameFn", pf.PlatformID)
		}
	}

	for _, want := range []string{
		"claude", "cursor", "opencode", "codex", "kiro",
		"gemini", "qoder", "codebuddy", "copilot", "droid",
		"pi", "kilo", "antigravity", "windsurf",
	} {
		if !ids[want] {
			t.Errorf("missing platform: %s", want)
		}
	}
}

func TestPlatformDir(t *testing.T) {
	if d := PlatformDir("claude"); d != ".claude/skills" {
		t.Errorf("expected '.claude/skills', got %q", d)
	}
	if d := PlatformDir("nonexistent"); d != "" {
		t.Errorf("expected empty for unknown platform, got %q", d)
	}
}

func TestCodexAgentEntry(t *testing.T) {
	skills := AllSkills()
	entry := CodexAgentEntry(skills)

	if !strings.Contains(entry, "# Agents") {
		t.Error("expected Agents heading")
	}
	for _, s := range skills {
		if !strings.Contains(entry, s.Name) {
			t.Errorf("missing skill %s in agent entry", s.Name)
		}
	}
}

func TestAllFormatFunctionsProduceContent(t *testing.T) {
	skills := AllSkills()
	formatters := []struct {
		name string
		fn   func(Skill) (string, string)
	}{
		{"Shared", FormatForShared},
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
		{"Kilo", FormatForKilo},
		{"Antigravity", FormatForAntigravity},
		{"Devin", FormatForDevin},
	}

	for _, s := range skills {
		for _, f := range formatters {
			fn, content := f.fn(s)
			if fn == "" {
				t.Errorf("%s/%s: empty filename", f.name, s.Name)
			}
			if content == "" {
				t.Errorf("%s/%s: empty content", f.name, s.Name)
			}
		}
	}
}

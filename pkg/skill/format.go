package skill

import (
	"fmt"
	"strings"
)

// FormatForShared generates the cross-platform shared skill file (.agents/skills/{name}/SKILL.md).
func FormatForShared(s Skill) (filename, content string) {
	filename = fmt.Sprintf("%s/SKILL.md", s.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
---

%s
`, s.Name, s.Description, s.Content)
	return
}

// FormatForClaudeCode generates a Claude Code skill file.
func FormatForClaudeCode(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForCursor generates a Cursor skill file.
func FormatForCursor(s Skill) (filename, content string) {
	filename = fmt.Sprintf("trellis-%s/SKILL.md", s.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
---

%s
`, s.Name, s.Description, s.Content)
	return
}

// FormatForOpenCode generates an OpenCode skill file.
func FormatForOpenCode(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForCodex generates a Codex skill file.
func FormatForCodex(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForKiro generates a Kiro skill file.
func FormatForKiro(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForGemini generates a Gemini CLI skill file.
func FormatForGemini(s Skill) (filename, content string) {
	filename = fmt.Sprintf("trellis-%s/SKILL.md", s.Name)
	return FormatForShared(s)
}

// FormatForQoder generates a Qoder skill file.
func FormatForQoder(s Skill) (filename, content string) {
	filename = fmt.Sprintf("trellis-%s/SKILL.md", s.Name)
	return FormatForShared(s)
}

// FormatForCodeBuddy generates a CodeBuddy skill file.
func FormatForCodeBuddy(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForCopilot generates a Copilot skill file.
func FormatForCopilot(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForDroid generates a Droid skill file.
func FormatForDroid(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForPi generates a Pi Agent skill file.
func FormatForPi(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForKilo generates a Kilo skill file.
func FormatForKilo(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForAntigravity generates an Antigravity skill file.
func FormatForAntigravity(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// FormatForDevin generates a Devin skill file.
func FormatForDevin(s Skill) (filename, content string) {
	return FormatForShared(s)
}

// PlatformFormat describes a skill file location for a specific platform.
type PlatformFormat struct {
	PlatformID  string
	Dir         string // directory relative to project root
	FilenameFn  func(Skill) string
}

// PlatformFormats returns the skill file locations for all supported platforms.
func PlatformFormats() []PlatformFormat {
	return []PlatformFormat{
		{PlatformID: "claude", Dir: ".claude/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "cursor", Dir: ".cursor/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("trellis-%s/SKILL.md", s.Name) }},
		{PlatformID: "opencode", Dir: ".opencode/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "codex", Dir: ".codex/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "kiro", Dir: ".kiro/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "gemini", Dir: ".gemini/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("trellis-%s/SKILL.md", s.Name) }},
		{PlatformID: "qoder", Dir: ".qoder/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("trellis-%s/SKILL.md", s.Name) }},
		{PlatformID: "codebuddy", Dir: ".codebuddy/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "copilot", Dir: ".github/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "droid", Dir: ".factory/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "pi", Dir: ".pi/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "kilo", Dir: ".kilocode/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "antigravity", Dir: ".agent/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
		{PlatformID: "windsurf", Dir: ".devin/skills", FilenameFn: func(s Skill) string { return fmt.Sprintf("%s/SKILL.md", s.Name) }},
	}
}

// PlatformDir returns the skills directory for a given platform ID.
func PlatformDir(platformID string) string {
	for _, pf := range PlatformFormats() {
		if pf.PlatformID == platformID {
			return pf.Dir
		}
	}
	return ""
}

// CodexAgentEntry generates the AGENTS.md entry for Codex skill registration.
func CodexAgentEntry(skills []Skill) string {
	var b strings.Builder
	b.WriteString("# Agents\n\n")
	for _, s := range skills {
		b.WriteString(fmt.Sprintf("## %s\n\n%s\n\n", s.Name, s.Description))
	}
	return b.String()
}

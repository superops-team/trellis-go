package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatForClaudeCode generates a Claude Code agent file (YAML frontmatter + Markdown).
func FormatForClaudeCode(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.md", a.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
tools: %s
---

%s
`, a.Name, a.Description, strings.Join(a.Tools, ", "), a.Content)
	return
}

// FormatForCursor generates a Cursor agent file (YAML frontmatter + Markdown).
func FormatForCursor(a Agent) (filename, content string) {
	return FormatForClaudeCode(a)
}

// FormatForOpenCode generates an OpenCode agent file (YAML frontmatter + permission object).
func FormatForOpenCode(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.md", a.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
tools: %s
permission:
  allow: [Read, Write, Edit, Bash, Glob, Grep]
  deny: []
---

%s
`, a.Name, a.Description, strings.Join(a.Tools, ", "), a.Content)
	return
}

// FormatForCodex generates a Codex agent file (TOML format).
func FormatForCodex(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.toml", a.Name)
	content = fmt.Sprintf(`name = "%s"
description = "%s"

[features]
multi_agent = false

[features.multi_agent_v2]
enabled = false

developer_instructions = """
%s
"""
`, a.Name, a.Description, a.Content)
	return
}

// FormatForKiro generates a Kiro agent file (JSON format).
func FormatForKiro(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.json", a.Name)
	type kiroAgent struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Tools       string `json:"tools"`
		Content     string `json:"content"`
	}
	ka := kiroAgent{
		Name:        a.Name,
		Description: a.Description,
		Tools:       strings.Join(a.Tools, ", "),
		Content:     a.Content,
	}
	b, _ := json.MarshalIndent(ka, "", "  ")
	content = string(b) + "\n"
	return
}

// FormatForGemini generates a Gemini CLI agent file (YAML frontmatter, pull-based prelude).
func FormatForGemini(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.md", a.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
tools: %s
---

# %s

%s

## Before acting

1. Read the active task's implement.jsonl to discover required files
2. Read prd.md for requirements
3. Read design.md if it exists
4. Read implement.md if it exists

%s
`, a.Name, a.Description, strings.Join(a.Tools, ", "), a.Name, a.Description, a.Content)
	return
}

// FormatForQoder generates a Qoder agent file (YAML frontmatter, pull-based prelude).
func FormatForQoder(a Agent) (filename, content string) {
	return FormatForGemini(a)
}

// FormatForCodeBuddy generates a CodeBuddy agent file (YAML frontmatter + Markdown, CC-compatible).
func FormatForCodeBuddy(a Agent) (filename, content string) {
	return FormatForClaudeCode(a)
}

// FormatForCopilot generates a Copilot agent file (.agent.md format).
func FormatForCopilot(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.agent.md", a.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
tools: %s
---

# %s

%s

## Before acting

1. Read the active task's implement.jsonl to discover required files
2. Read prd.md for requirements
3. Read design.md if it exists
4. Read implement.md if it exists
`, a.Name, a.Description, strings.Join(a.Tools, ", "), a.Name, a.Content)
	return
}

// FormatForDroid generates a Droid agent file (YAML frontmatter + Markdown, CC-compatible).
func FormatForDroid(a Agent) (filename, content string) {
	return FormatForClaudeCode(a)
}

// FormatForPi generates a Pi Agent agent file (YAML frontmatter + model/thinking).
func FormatForPi(a Agent) (filename, content string) {
	filename = fmt.Sprintf("%s.md", a.Name)
	content = fmt.Sprintf(`---
name: %s
description: |
  %s
tools: %s
model: inherit
thinking: inherit
---

%s
`, a.Name, a.Description, strings.Join(a.Tools, ", "), a.Content)
	return
}

package command

import (
	"fmt"
	"strings"
)

// Command represents a Trellis slash command definition.
type Command struct {
	Name        string
	Description string
	Content     string
}

// StartCommand returns the /trellis:start command content.
func StartCommand() Command {
	return Command{
		Name:        "start",
		Description: "Start a Trellis session — load workflow, context, and active tasks",
		Content: `# /trellis:start — Trellis Session Start

Run this command at the start of a new session.

## Steps

1. **Load workflow** — Read .trellis/workflow.md to understand the development workflow
2. **Load context** — Run 'trellis hook session-start' for developer identity, git status, and active tasks
3. **Load spec index** — Read .trellis/spec/index.md for project conventions
4. **Report** — Summarize current state and ask the user what to work on

## Task Classification

| User input | Action |
|------------|--------|
| Question / explanation / lookup | Answer directly, no task needed |
| Small single-round edit | Ask if task should be created |
| Multi-file or persistent work | Create task via 'trellis task create <name>' |
`,
	}
}

// FinishWorkCommand returns the /trellis:finish-work command content.
func FinishWorkCommand() Command {
	return Command{
		Name:        "finish-work",
		Description: "Archive the current task and record the session",
		Content: `# /trellis:finish-work — Archive Task and Record Session

Run this command when implementation is complete and code is committed.

## Prerequisites

- All changes must be committed (git status must be clean)
- An active task must exist

## Steps

1. **Check git status** — Run 'git status --porcelain' to verify clean working tree
2. **If dirty** — Refuse and guide user to commit first
3. **Archive task** — Run 'trellis task archive <task-id>'
4. **Record session** — Run 'trellis hook record-session --title "<task>" --commits "<hashes>" --summary "<changes>"'
5. **Report** — Summarize what was accomplished
`,
	}
}

// ContinueCommand returns the /trellis:continue command content.
func ContinueCommand() Command {
	return Command{
		Name:        "continue",
		Description: "Advance the workflow to the next step",
		Content: `# /trellis:continue — Advance Workflow

Run this command to advance the current task to the next workflow step.

## Steps

1. **Read task status** — Check .trellis/tasks/*/task.json for current status
2. **Read workflow** — Check .trellis/workflow.md for current phase/step
3. **Inject workflow state** — Run 'trellis hook inject-workflow-state --state <status>'
4. **Advance** — Move to the next step based on workflow.md

## Typical Flow

| State | Next action |
|-------|-------------|
| planning | Start implementation |
| in_progress | Implement -> Check -> Update spec |
| completed | Ready for finish-work |
`,
	}
}

// AllCommands returns all three Trellis commands.
func AllCommands() []Command {
	return []Command{
		StartCommand(),
		FinishWorkCommand(),
		ContinueCommand(),
	}
}

// PlatformFormat represents the file format for a specific platform.
type PlatformFormat struct {
	PlatformID  string
	FilePattern string // e.g., ".claude/commands/trellis/{name}.md"
}

// FormatForClaudeCode generates a Claude Code command file.
func FormatForClaudeCode(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis/%s.md", cmd.Name)
	content := fmt.Sprintf(`---
name: trellis:%s
description: %s
---

%s
`, cmd.Name, cmd.Description, cmd.Content)
	return filename, content
}

// FormatForCursor generates a Cursor command file.
func FormatForCursor(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis-%s.md", cmd.Name)
	content := fmt.Sprintf(`---
name: trellis-%s
description: %s
---

%s
`, cmd.Name, cmd.Description, cmd.Content)
	return filename, content
}

// FormatForCodex generates a Codex prompt file.
func FormatForCodex(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis-%s.md", cmd.Name)
	content := fmt.Sprintf(`# /trellis:%s

%s
`, cmd.Name, cmd.Content)
	return filename, content
}

// FormatForGemini generates a Gemini CLI TOML command file.
func FormatForGemini(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis/%s.toml", cmd.Name)
	content := fmt.Sprintf(`[command]
name = "trellis:%s"
description = "%s"

[[prompt]]
content = """
%s
"""
`, cmd.Name, cmd.Description, cmd.Content)
	return filename, content
}

// FormatForQoder generates a Qoder command file with YAML frontmatter.
func FormatForQoder(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis-%s.md", cmd.Name)
	content := fmt.Sprintf(`---
name: trellis-%s
description: %s
---

%s
`, cmd.Name, cmd.Description, cmd.Content)
	return filename, content
}

// FormatForCopilot generates a Copilot prompt file.
func FormatForCopilot(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis-%s.prompt.md", cmd.Name)
	content := fmt.Sprintf(`# /trellis:%s

%s
`, cmd.Name, cmd.Content)
	return filename, content
}

// FormatForWorkflow generates a workflow file (for Kilo, Antigravity, Devin).
func FormatForWorkflow(cmd Command) (string, string) {
	filename := fmt.Sprintf("trellis-%s.md", cmd.Name)
	content := fmt.Sprintf(`# Trellis %s Workflow

%s
`, strings.Title(cmd.Name), cmd.Content)
	return filename, content
}

// PlatformFormats returns all supported platform formats.
func PlatformFormats() []PlatformFormat {
	return []PlatformFormat{
		{PlatformID: "claude", FilePattern: ".claude/commands/trellis/{name}.md"},
		{PlatformID: "cursor", FilePattern: ".cursor/commands/trellis-{name}.md"},
		{PlatformID: "opencode", FilePattern: ".opencode/commands/trellis/{name}.md"},
		{PlatformID: "codex", FilePattern: ".codex/prompts/trellis-{name}.md"},
		{PlatformID: "gemini", FilePattern: ".gemini/commands/trellis/{name}.toml"},
		{PlatformID: "qoder", FilePattern: ".qoder/commands/trellis-{name}.md"},
		{PlatformID: "codebuddy", FilePattern: ".codebuddy/commands/trellis/{name}.md"},
		{PlatformID: "droid", FilePattern: ".factory/commands/trellis/{name}.md"},
		{PlatformID: "pi", FilePattern: ".pi/prompts/trellis-{name}.md"},
		{PlatformID: "copilot", FilePattern: ".github/prompts/trellis-{name}.prompt.md"},
	}
}

package agent

// Agent represents a Trellis sub-agent definition.
type Agent struct {
	Name        string
	Description string
	Tools       []string
	Content     string
}

// ImplementAgent returns the trellis-implement sub-agent definition.
func ImplementAgent() Agent {
	return Agent{
		Name:        "trellis-implement",
		Description: "Implements code changes according to the active Trellis task's PRD, design, and implement plan.",
		Tools:       []string{"Read", "Write", "Edit", "Bash", "Glob", "Grep"},
		Content: `# trellis-implement

You are the trellis-implement sub-agent in the Trellis workflow.

## Context Loading

1. Read ` + "`" + `implement.jsonl` + "`" + ` for the file manifest
2. Read ` + "`" + `prd.md` + "`" + ` for requirements
3. Read ` + "`" + `design.md` + "`" + ` if present
4. Read ` + "`" + `implement.md` + "`" + ` if present

## Implementation Rules

- Follow the project's coding conventions (read from ` + "`" + `.trellis/spec/` + "`" + `)
- Write minimal, focused changes
- Do not commit — the main session handles git
- Report what was changed and why
`,
	}
}

// CheckAgent returns the trellis-check sub-agent definition.
func CheckAgent() Agent {
	return Agent{
		Name:        "trellis-check",
		Description: "Code quality check expert. Reviews diffs against specs, runs lint/typecheck/test, self-fixes.",
		Tools:       []string{"Read", "Write", "Edit", "Bash", "Glob", "Grep"},
		Content: `# trellis-check

You are the trellis-check sub-agent in the Trellis workflow.

## Flow

1. ` + "`" + `git diff --name-only HEAD` + "`" + ` → find changed files
2. Discover applicable spec layers from ` + "`" + `.trellis/spec/` + "`" + `
3. Compare diff against each layer's quality checklist
4. Run lint, typecheck, and tests
5. If issues found: fix → re-verify (max 3 rounds)
6. Report: PASSED or FAILED with details
`,
	}
}

// ResearchAgent returns the trellis-research sub-agent definition.
func ResearchAgent() Agent {
	return Agent{
		Name:        "trellis-research",
		Description: "Read-only codebase research agent. Searches, analyzes structure, and reports findings.",
		Tools:       []string{"Read", "Glob", "Grep", "Bash"},
		Content: `# trellis-research

You are the trellis-research sub-agent in the Trellis workflow.
Read-only agent. No file writes.

## Flow

1. Read ` + "`" + `research.jsonl` + "`" + ` for context
2. Search codebase for relevant patterns
3. Analyze and report findings
`,
	}
}

// AllAgents returns all three Trellis sub-agent definitions.
func AllAgents() []Agent {
	return []Agent{
		ImplementAgent(),
		CheckAgent(),
		ResearchAgent(),
	}
}

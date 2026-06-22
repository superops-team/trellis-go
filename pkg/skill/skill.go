package skill

// Skill represents a Trellis auto-trigger skill definition.
type Skill struct {
	Name        string
	Description string
	Content     string
}

// BrainstormSkill returns the trellis-brainstorm skill.
func BrainstormSkill() Skill {
	return Skill{
		Name: "trellis-brainstorm",
		Description: `Use when the user wants a new feature or the requirements are unclear.
Clarifies requirements, inspects evidence, and drafts planning artifacts.`,
		Content: `# trellis-brainstorm

## Trigger check
Verify the user is requesting new work that needs planning.

## Steps
1. Inspect codebase, existing specs, and task history
2. Propose task name and create via ` + "`" + `task.py create` + "`" + `
3. Draft ` + "`" + `prd.md` + "`" + ` with requirements and acceptance criteria
4. Ask one question at a time with recommended answer
5. For complex tasks, add ` + "`" + `design.md` + "`" + ` and ` + "`" + `implement.md` + "`" + `
`,
	}
}

// BeforeDevSkill returns the trellis-before-dev skill.
func BeforeDevSkill() Skill {
	return Skill{
		Name: "trellis-before-dev",
		Description: `Use before touching code in a task. Reads relevant specs so the AI knows
the conventions before writing, not after.`,
		Content: `# trellis-before-dev

## Steps
1. Identify affected packages from the task context
2. Read spec index for each package
3. Read pre-development checklist guidelines
4. Confirm conventions are understood before proceeding
`,
	}
}

// CheckSkill returns the trellis-check skill.
func CheckSkill() Skill {
	return Skill{
		Name: "trellis-check",
		Description: `Use after implementing code changes. Reviews diffs against specs,
runs lint/typecheck/test, and self-fixes issues.`,
		Content: `# trellis-check

## Steps
1. ` + "`" + `git diff --name-only HEAD` + "`" + ` to find changed files
2. Discover applicable spec layers from ` + "`" + `.trellis/spec/` + "`" + `
3. Compare diff against each layer's quality checklist
4. Run lint, typecheck, and tests
5. If issues found: fix -> re-verify (max 3 rounds)
6. Report: PASSED or FAILED with details
`,
	}
}

// UpdateSpecSkill returns the trellis-update-spec skill.
func UpdateSpecSkill() Skill {
	return Skill{
		Name: "trellis-update-spec",
		Description: `Use when there is knowledge worth capturing after completing work.
Updates spec files with new learnings and conventions.`,
		Content: `# trellis-update-spec

## Steps
1. Identify knowledge worth capturing
2. Determine the target spec layer
3. Update spec file with new knowledge
4. Keep specs concise and actionable
`,
	}
}

// BreakLoopSkill returns the trellis-break-loop skill.
func BreakLoopSkill() Skill {
	return Skill{
		Name: "trellis-break-loop",
		Description: `Use when encountering a stubborn bug or repeatedly fixing the same issue.
Analyzes root cause and proposes preventive measures.`,
		Content: `# trellis-break-loop

## Steps
1. Analyze the bug root cause
2. Identify why existing flow didn't catch it
3. Propose preventive measures
4. Update spec if necessary
`,
	}
}

// AllSkills returns all five Trellis auto-trigger skills.
func AllSkills() []Skill {
	return []Skill{
		BrainstormSkill(),
		BeforeDevSkill(),
		CheckSkill(),
		UpdateSpecSkill(),
		BreakLoopSkill(),
	}
}

# How Trellis Works

## Core Concepts

### The `.trellis/` Directory

Trellis stores everything in `.trellis/` inside your Git repository:

```
.trellis/
├── config.yaml          # Trellis configuration
├── .developer           # Developer identity (gitignored)
├── spec/                # Engineering specs
│   └── gap-analysis/    # Gap analysis specs
├── tasks/               # Active tasks
│   └── MM-DD-task-name/
│       ├── task.json    # Task metadata
│       ├── prd.md       # Product requirements
│       ├── implement.jsonl  # Implementation context
│       └── check.jsonl      # Verification context
├── workspace/           # Developer workspace
│   └── <developer>/
│       ├── journal-1.md # Session journal
│       └── index.json   # Session index
└── templates/           # Spec templates
```

### Task Lifecycle

```
planning → in_progress → completed (archived)
```

1. **planning** — Task created, PRD being written
2. **in_progress** — PRD ready, implementation started
3. **completed** — Task archived to `tasks/archive/YYYY-MM/`

### 4-Phase Workflow

Trellis defines a structured workflow for AI coding agents:

1. **Plan** — Clarify requirements, draft PRD
2. **Implement** — Write code with context injection
3. **Verify** — Review, test, and validate
4. **Finish** — Archive task, record session

### Context Injection

Trellis assembles context for each phase:

- **Implement phase**: PRD + implementation specs + context files
- **Check phase**: PRD + verification specs + context files
- **Research phase**: Research specs + context files

Context is managed through JSONL manifest files (`implement.jsonl`, `check.jsonl`, `research.jsonl`).

### Platform Integration

Trellis generates platform-specific hook files that call `trellis hook` commands:

- **Push-based** (Claude, Cursor, Kiro): Hook scripts in platform config dir
- **Pull-based** (Codex, Gemini, Copilot): Agent/skill definitions
- **Agentless** (Devin, Kilo, Antigravity): Workflow + skill files

## Architecture

Trellis is a single Go binary. See [Architecture](../advanced/architecture.md) for details.

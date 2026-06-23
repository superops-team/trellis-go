# Trellis

An engineering framework for AI coding. Persist specs, tasks, and memory into your repo so any coding agent works to your engineering standards.

## Why Trellis?

AI coding agents are powerful but inconsistent. Without shared context, every session starts from scratch. Trellis gives your AI agents:

- **Persistent specs** — Engineering standards that survive across sessions
- **Task lifecycle** — Structured create → start → archive workflow
- **Context injection** — Automatic context assembly for implement/check/research phases
- **Multi-platform** — Works with Claude Code, Cursor, Codex, Copilot, and 12+ more

## Quick Start

```bash
# Install
go install github.com/superops-team/trellis-go/cmd/trellis@latest

# Initialize in a Git repo
git init my-project && cd my-project
trellis init --developer alice --platform claude

# Create your first task
trellis task create my-feature
echo "# My Feature PRD" > .trellis/tasks/*-my-feature/prd.md
trellis task start my-feature
```

## How It Works

1. **Init** — Trellis sets up `.trellis/` and generates platform hook files
2. **Create** — Each task gets a directory with `task.json`, `prd.md`, and phase manifests
3. **Context** — Before coding, Trellis assembles relevant specs, PRD, and context files
4. **Workflow** — Plan → Implement → Verify → Finish state machine guides the process
5. **Journal** — Session recordings track what happened and why

## Supported Platforms

Trellis supports 16+ AI coding platforms: Claude Code, Cursor, Codex, OpenCode, Gemini CLI, Kiro, Copilot, Devin, Kilo, Pi, CodeBuddy, Droid, Qoder, Antigravity, Reasonix, ZCode.

## Next Steps

- [Installation & First Task](start/install-and-first-task.md)
- [How It Works](start/how-it-works.md)
- [Everyday Use](start/everyday-use.md)
- [Architecture](advanced/architecture.md)
- [Configuration](advanced/configuration.md)

<div align="center">

# Trellis (Go)

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](https://github.com/superops-team/trellis-go/actions)

[English](README.md) | [中文](README.zh-CN.md)

</div>

> An engineering framework for AI coding. Persist specs, tasks, and memory into your repo so any coding agent works to your engineering standards.

Trellis is a port of the original [TypeScript/Python Trellis](https://github.com/mindfold-ai/Trellis) framework to Go, designed for single-binary distribution, high performance, and seamless integration with 15+ AI coding platforms.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Supported Platforms](#supported-platforms)
- [Architecture](#architecture)
- [CLI Commands](#cli-commands)
- [Workflow](#workflow)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **15+ AI Platforms** — Built-in support for Claude, Cursor, Codex, Copilot, Windsurf, and more
- **4-Phase Workflow** — Plan → Implement → Verify → Finish state machine
- **Task Lifecycle** — Create, start, archive tasks with automatic organization
- **Context Builder** — JSONL-based manifest system for AI agent context injection
- **Platform Hooks** — Auto-generate platform-specific configuration files
- **Atomic File Operations** — Safe concurrent writes with SHA256 hashing
- **Single Binary** — Statically compiled Go binary with embedded templates
- **Thread-Safe Registry** — Concurrent-safe platform and spec management

## Quick Start

### Installation

```bash
go install github.com/superops-team/trellis-go/cmd/trellis@latest
```

Or download a pre-built binary from [Releases](https://github.com/superops-team/trellis-go/releases).

### Initialize a Project

```bash
# Inside a Git repository
git init my-project && cd my-project

# Initialize Trellis with default platform (Claude)
trellis init --developer alice

# Or initialize with multiple platforms
trellis init --developer alice --platform claude --platform cursor --platform codex
```

### Create Your First Task

```bash
# Create a new task
trellis task create user-auth

# List all tasks
trellis task list

# Check current active task
trellis task current
```

### Project Structure

```
my-project/
├── .git/
├── .trellis/                 # Trellis workspace
│   ├── config.yaml           # Developer & platform config
│   ├── .version              # Trellis version
│   ├── workflow.md           # 4-phase workflow definition
│   ├── spec/                 # Engineering specs
│   ├── tasks/                # Active tasks
│   │   ├── 06-13-user-auth/
│   │   │   ├── task.json     # Task metadata
│   │   │   ├── prd.md        # Product requirements
│   │   │   ├── implement.jsonl  # Context manifest
│   │   │   ├── check.jsonl      # Verification manifest
│   │   │   └── research/     # Research notes
│   │   └── archive/          # Archived tasks (YYYY-MM/)
│   ├── workspace/            # Shared workspace
│   └── .runtime/sessions/    # Active session tracking
├── .claude/                  # Claude-specific config
├── .cursor/                  # Cursor-specific config
└── ...
```

## Supported Platforms

| Platform | Class | Agent | Hooks |
|----------|-------|-------|-------|
| Claude Code | Push-based | Yes | Yes |
| Cursor | Push-based | Yes | Yes |
| Codex | Pull-based | Yes | Yes |
| OpenCode | Push-based | Yes | Yes |
| Gemini CLI | Pull-based | Yes | No |
| Kiro | Push-based | Yes | Yes |
| Copilot | Pull-based | No | No |
| Windsurf | Agentless | No | No |
| Kilo | Agentless | No | No |
| Pi | Push-based | Yes | Yes |
| CodeBuddy | Push-based | Yes | Yes |
| Droid | Push-based | Yes | Yes |
| Qoder | Pull-based | Yes | No |
| Antigravity | Agentless | No | No |
| Reasonix | Push-based | Yes | Yes |

**Platform Classes:**
- **Push-based** — Agent initiates execution (Claude, Cursor, etc.)
- **Pull-based** — IDE pulls context on demand (Codex, Copilot, etc.)
- **Agentless** — No agent capability, manual workflow (Windsurf, Kilo, etc.)

## Architecture

```
cmd/trellis/          CLI commands (cobra)
pkg/
  platform/           Platform definitions & registry
  fsutil/             Atomic file operations & hashing
  config/             YAML configuration management
  template/           embed.FS template engine
  task/               Task lifecycle & manifest
  workflow/           4-phase state machine
  context/            Context builder (JSONL manifests)
  hook/               Platform hook generator
  spec/               Spec loader & index
  git/                Git command wrapper
internal/
  embed/              Embedded template assets
  testutil/           Test helpers
```

## CLI Commands

```bash
# Initialize Trellis in current repo
trellis init [flags]
  --developer, -u    Developer name
  --platform, -p     Platform to configure (repeatable)
  --all              Configure all platforms

# Task management
trellis task create <name>     Create a new task
trellis task list              List all tasks
trellis task current           Show active task
trellis task start [id]        Start a task
trellis task archive [id]      Archive a completed task

# Context management
trellis context add <file>     Add file to current context
trellis context build          Build and output context

# Maintenance
trellis update                 Update platform configs
trellis uninstall              Remove Trellis (use --keep-tasks to preserve)
trellis version                Show version
```

## Workflow

Trellis enforces a 4-phase engineering workflow via `workflow.md`:

```
[workflow-state:PLAN]
→ Brainstorm requirements and write PRD

[workflow-state:IMPLEMENT]
→ Write code from the PRD

[workflow-state:VERIFY]
→ Review code against specs and run checks

[workflow-state:FINISH]
→ Archive task and update journals
```

State transitions are validated by the state machine:
- `plan` → `implement`
- `implement` → `verify`
- `verify` → `implement` | `finish`

## Development

### Prerequisites

- Go 1.23+
- Git

### Build

```bash
go build -o trellis ./cmd/trellis
```

### Run Tests

```bash
# All tests
go test ./...

# Unit tests only
go test ./pkg/...

# E2E tests
go test ./cmd/trellis -run TestE2E_
```

## Testing

The project includes comprehensive test coverage:

- **Unit Tests** — All `pkg/` packages with table-driven tests
- **E2E Tests** — 8 real-world scenarios:
  1. New project initialization + first task creation
  2. Multi-platform configuration (Claude + Cursor + Codex)
  3. Full task lifecycle (create → start → archive)
  4. Context building for AI agent injection
  5. Task listing and current task query
  6. Uninstall while preserving tasks
  7. Invalid platform error handling
  8. Non-git repository error handling

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure all tests pass and follow the existing code style.

## License

[MIT](LICENSE)

---

<div align="center">

Made with Go for AI-native engineering teams.

</div>

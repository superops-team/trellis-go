<div align="center">

# Trellis (Go)

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](https://github.com/superops-team/trellis-go/actions)

[English](README.md) | [中文](README.zh-CN.md)

</div>

> An engineering framework for AI coding. Persist specs, tasks, and memory into your repo so any coding agent works to your engineering standards.

Trellis is a port of the original [TypeScript/Python Trellis](https://github.com/mindfold-ai/Trellis) framework to Go, designed for single-binary distribution, high performance, and seamless integration with 16 AI coding platforms.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Beginner User Guide](#beginner-user-guide)
- [Documentation](#documentation)
- [Supported Platforms](#supported-platforms)
- [Architecture](#architecture)
- [CLI Commands](#cli-commands)
- [Workflow](#workflow)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **16 AI Platforms** — Built-in support for Claude Code, Cursor, Codex, Devin, ZCode, and more
- **4-Phase Workflow** — Plan → Implement → Verify → Finish state machine
- **Task Lifecycle** — Create, start, archive tasks with automatic organization and subtasks
- **Context Builder** — JSONL-based manifest system for AI agent context injection
- **Platform Hooks** — Auto-generate platform-specific configuration files (hooks, agents, skills, workflows)
- **Session Journal** — Record and recall AI agent sessions with commit tracking
- **Slash Commands** — Built-in `/plan`, `/implement`, `/verify`, `/finish` workflow commands
- **Sub-Agent Support** — Define and dispatch parallel sub-agents for complex tasks
- **Auto-Skills** — Platform-specific skill definitions for AI coding agents
- **Spec Templates** — Reusable engineering spec templates (Next.js, Cloudflare Workers, Electron)
- **Skills Marketplace** — Community skills for frontend optimization, memory recall, and more
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

## Beginner User Guide

If you are new to Trellis, start with the step-by-step guide:

- [English Beginner Guide](docs/USAGE.md)
- [中文小白使用指南](docs/USAGE.zh-CN.md)

The guide explains installation, initialization, task lifecycle, context manifests, troubleshooting, and includes Mermaid workflow diagrams.

## Documentation

Full documentation is available in the `docs/` directory:

| Section | Description |
|---------|-------------|
| [Home](docs/index.md) | Overview and quick start |
| [Installation & First Task](docs/start/install-and-first-task.md) | Setup guide |
| [How It Works](docs/start/how-it-works.md) | Core concepts and data flow |
| [Everyday Use](docs/start/everyday-use.md) | Daily workflow patterns |
| [Real-World Scenarios](docs/start/real-world-scenarios.md) | Usage examples |
| [Architecture](docs/advanced/architecture.md) | Package overview and design |
| [Configuration](docs/advanced/configuration.md) | Config file, env vars, CLI flags |
| [Multi-Platform](docs/advanced/multi-platform.md) | Working with multiple platforms |
| [Roadmap](docs/advanced/roadmap.md) | Current and future plans |
| [Document Index](docs/llms.txt) | LLM-friendly doc index |

## Supported Platforms

| Platform | Class | Agent | Hooks | Skills |
|----------|-------|:-----:|:-----:|:------:|
| Claude Code | Push-based | ✅ | ✅ | ✅ |
| Cursor | Push-based | ✅ | ✅ | ✅ |
| Codex | Pull-based | ✅ | ✅ | ✅ |
| OpenCode | Push-based | ✅ | ✅ | ✅ |
| Gemini CLI | Pull-based | ✅ | — | — |
| Kiro | Push-based | ✅ | ✅ | ✅ |
| Copilot | Pull-based | — | — | — |
| Devin | Agentless | — | — | — |
| Kilo | Agentless | — | — | — |
| Pi | Push-based | ✅ | ✅ | ✅ |
| CodeBuddy | Push-based | ✅ | ✅ | ✅ |
| Droid | Push-based | ✅ | ✅ | ✅ |
| Qoder | Pull-based | ✅ | — | — |
| Antigravity | Agentless | — | — | — |
| Reasonix | Push-based | ✅ | ✅ | ✅ |
| ZCode | Push-based | ✅ | ✅ | ✅ |

> **Note:** Devin replaces the previous Windsurf platform. `--windsurf` is retained as an alias for backward compatibility.

**Platform Classes:**
- **Push-based** — Agent initiates execution (Claude, Cursor, etc.)
- **Pull-based** — IDE pulls context on demand (Codex, Copilot, etc.)
- **Agentless** — No agent capability, manual workflow (Devin, Kilo, etc.)

## Architecture

```
cmd/trellis/          CLI commands (cobra)
pkg/
  agent/              AI agent management
  command/            CLI command definitions
  config/             YAML configuration management
  configurator/       Platform config generation
  context/            Context builder (JSONL manifests)
  fsutil/             Atomic file operations & hashing
  git/                Git command wrapper
  hook/               Platform hook execution
  manifest/           Manifest file management
  platform/           Platform definitions & registry (16 platforms)
  prd/                PRD management
  session/            Session journal
  skill/              Skill management
  spec/               Spec loader & index
  task/               Task lifecycle & subtasks
  template/           embed.FS template engine
  update/             Template sync
  upgrade/            Binary upgrade
  workflow/           4-phase state machine
internal/
  embed/              Embedded template assets
  testutil/           Test helpers
```

See [docs/advanced/architecture.md](docs/advanced/architecture.md) for detailed design documentation.

## CLI Commands

### Initialization

```bash
trellis init [flags]
  --developer, -u    Developer name
  --platform, -p     Platform to configure (repeatable)
  --all              Configure all 16 platforms
```

### Task Management

```bash
trellis task create <name>       Create a new task
trellis task list                List all tasks
trellis task current             Show active task
trellis task info <id>           Show task details
trellis task start <id>          Start a task (planning → in_progress)
trellis task archive <id>        Archive a completed task
trellis task edit <id>           Edit task fields
trellis task add-subtask <id>    Add a subtask
trellis task done-subtask <id>   Mark a subtask as done
trellis task add-spec <id>       Associate a spec with a task
trellis task list-specs <id>     List associated specs
```

### Context Management

```bash
trellis context add <file> --task <id> [--phase implement|check]
trellis context build --task <id> --phase implement|check
trellis context build --phase research
trellis task add-context <file> --task <id> [--phase implement|check]
trellis task list-context --task <id>
trellis task remove-context <file> --task <id>
```

### Agent Hooks

```bash
trellis hook inject-context           Print context for agent hook
trellis hook inject-workflow-state    Print workflow-state prompt
trellis hook session-start            Print session start context
trellis hook record-session           Record a session in the journal
trellis hook list-sessions            List recorded sessions
```

### Maintenance

```bash
trellis update              Sync templates and configuration
trellis upgrade             Upgrade to the latest version
trellis uninstall           Remove Trellis (use --keep-tasks to preserve)
trellis version             Show version
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

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:

- Development environment setup
- Code style and conventions
- Pull request process
- Issue templates
- Testing standards

Quick start:

```bash
git clone git@github.com:superops-team/trellis-go.git
cd trellis-go
go test ./...
go build ./cmd/trellis
```

## License

[MIT](LICENSE)

---

<div align="center">

Made with Go for AI-native engineering teams.

</div>

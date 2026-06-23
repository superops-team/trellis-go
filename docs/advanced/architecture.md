# Architecture

Trellis is a single Go binary with a modular package architecture.

## Package Overview

```
cmd/trellis/          # CLI entry point (Cobra commands)
pkg/
├── agent/            # AI agent management
├── command/          # CLI command definitions
├── config/           # Configuration loading
├── configurator/     # Platform config generation
├── context/          # Context builder (JSONL manifests)
├── fsutil/           # File system utilities
├── git/              # Git integration
├── hook/             # Platform hook execution
├── manifest/         # Manifest file management
├── platform/         # Platform registry (16+ platforms)
├── prd/              # PRD management
├── session/          # Session journal
├── skill/            # Skill management
├── spec/             # Spec management
├── task/             # Task lifecycle
├── template/         # Template engine
├── update/           # Template sync
├── upgrade/          # Binary upgrade
└── workflow/         # 4-phase workflow engine
internal/
├── embed/            # Embedded templates
└── testutil/         # Test utilities
```

## Key Design Decisions

### Thread-Safe Registry

The platform registry uses `sync.RWMutex` for concurrent-safe reads and writes. Platforms are registered at startup and rarely modified.

### Atomic File Operations

All file writes use SHA256 hashing for integrity verification. The `fsutil` package provides safe concurrent writes with temp + rename pattern.

### Embedded Templates

Platform hooks, spec templates, and skill files are embedded into the binary using Go's `embed` package. No external file dependencies at runtime.

### JSONL Context Format

Context manifests use JSONL (JSON Lines) format — one JSON object per line. This allows streaming reads and append-only writes.

## Data Flow

```
CLI Command → Cobra Handler → Package API → File System
                  │
                  ├── config.Load() → config.yaml
                  ├── platform.NewRegistry() → builtins
                  ├── task.Manager → .trellis/tasks/
                  ├── context.Builder → .jsonl files
                  └── configurator.Run() → platform hook files
```

## State Machine

```
Task Lifecycle: planning → in_progress → completed
                     │            │              │
                     ▼            ▼              ▼
                 PRD write    Context build   Archive to
                              Implementation  tasks/archive/
```

## Concurrency Model

- **Registry**: RWMutex (read-heavy)
- **Task manager**: File-based locking (per-task)
- **Context builder**: Append-only JSONL (no locking needed)
- **Configurator**: Sequential per-platform generation

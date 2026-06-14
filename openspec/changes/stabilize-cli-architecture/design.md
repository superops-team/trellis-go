# Design: Stabilize CLI Architecture

## Overview

The design principle is simple: `cmd/trellis` is an adapter, not a second implementation. Domain behavior lives in `pkg/*`; Cobra commands parse flags, resolve the Trellis root, call package APIs, and format output.

## Scope Guard

This plan intentionally avoids a large rewrite. The first implementation PR SHALL be limited to these areas:

- `cmd/trellis/*` for root resolution and CLI delegation.
- `cmd/trellis/e2e_test.go` for hermetic binary setup and real lifecycle assertions.
- `pkg/task/manager.go` only where needed to make archive state consistent.
- Tests for the changed command and task-manager behavior.

Context, manifest, template, spec-loader, and module identity work are staged after P0/P1 unless the first PR remains small and green. If implementation starts touching unrelated packages before task CLI delegation is complete, stop and split the work.

## Architecture Decisions

### 1. Centralize root resolution

Create one unexported root-resolution helper in `cmd/trellis`. Do not add a new application package for this slice; a public `pkg/app` would be more abstraction than the project currently needs.

Contract:

- If `--root` is empty, commands use the current working directory as the repository root and `.trellis` as the Trellis data directory.
- If `--root` points to a directory ending in `.trellis`, treat it as the Trellis data directory and use its parent as repository root.
- If `--root` points to any other directory, treat it as the repository root and use `<root>/.trellis` as the Trellis data directory.
- Every command gets both paths from this helper:
  - `RepoRoot`
  - `TrellisDir`
  - `TasksDir`
  - `SpecDir`

This keeps backward compatibility with the existing E2E usage in `cmd/trellis/e2e_test.go:14`, which currently passes `--root <repo>/.trellis`.

Suggested local type:

```go
type resolvedPaths struct {
    RepoRoot   string
    TrellisDir string
    TasksDir   string
    SpecDir    string
}
```

The helper should not check whether `.trellis` exists. `init` must be able to resolve a future `.trellis` path before creating it. Each command remains responsible for validating the state it needs.

### 2. Make CLI commands delegate to packages

`cmd/trellis/task.go` should stop constructing task files directly. It should parse arguments, call `pkg/task.Manager`, and print stable output.

Mapping:

| CLI command | Package API |
|-------------|-------------|
| `task create <name>` | `task.NewManager(paths.TasksDir).Create(name, opts)` |
| `task list` | `task.NewManager(paths.TasksDir).List()` |
| `task start <id>` | `task.NewManager(paths.TasksDir).Start(id)` |
| `task archive <id>` | `task.NewManager(paths.TasksDir).Archive(id)` |
| `task current` | Keep existing no-active-task behavior unless session tracking is implemented in a separate spec. |

CLI output should be stable and testable:

- Create prints `Created task: <path>`.
- List prints task IDs, one per line.
- Start prints `Started task: <id>`.
- Archive prints `Archived task: <id>`.
- Current prints active task info or a clear no-active-task message.

Do not add session tracking in this change. That would introduce a second lifecycle concept and expand the blast radius.

### 3. Consolidate manifest schema ownership

The manifest concepts in `pkg/context/manifest.go:10` and `pkg/task/manifest.go:10` should have one owner.

Preferred path for P2:

- Create `pkg/manifest` with:
  - `Entry`
  - `Manifest`
  - `Load(path string) (*Manifest, error)`
  - `Save(path string, manifest *Manifest) error`
- Update `pkg/context` and `pkg/task` to import `pkg/manifest`.
- Keep compatibility aliases only if needed to avoid a large mechanical churn.

Constraints:

- `pkg/manifest` must contain only data types and JSONL load/save functions.
- It must not import `pkg/task`, `pkg/context`, `pkg/spec`, or Cobra.
- It must preserve the existing JSON field names: `path`, `description`, `required`.

### 4. Wire context commands to real behavior

`context build` should call `pkg/context.Builder` after task lookup is stable.

MVP behavior:

- `trellis context build --phase implement --task <id>` prints implementation context.
- `trellis context build --phase check --task <id>` prints verification context.
- `trellis context build --phase research` prints research context and does not require `--task`.

If no task is supplied for implement/check, use `task.Current()` once current-session tracking exists. Until then, return a clear error that `--task` is required.

`context add <file>` should append to the current task manifest through `pkg/task.Manager.AddContext`.

MVP behavior:

- Add flags `--task <id>`, `--phase implement|check`, `--required`, and `--description`.
- Resolve relative manifest paths against the Trellis directory, so `spec/auth.md` maps to `<repo>/.trellis/spec/auth.md`.
- Reject absolute paths and `..` segments for manifest entries. This preserves portable task manifests and prevents escaping the Trellis workspace.

### 5. Make template rendering explicit

`pkg/template.Engine.Render` should use `ShouldTemplate` to decide whether to parse a text file as a template.

Behavior:

- Binary files are copied unchanged.
- Files with extensions in `TemplateFileExtensions` are rendered.
- Other text files are copied unchanged.
- Missing template keys return an error using `template.Option("missingkey=error")`.
- `.template-hashes.json` write errors are returned. Keep a marshal-error check for correctness, but do not design special machinery around it because `map[string]string` JSON marshaling should not fail in practice.

### 6. Fix E2E test harness

Replace the hardcoded `/tmp/trellis-test` dependency in `cmd/trellis/e2e_test.go:14`.

Preferred path:

- Build the binary in `TestMain` into a package-level `os.MkdirTemp` directory. `testing.T.TempDir()` is not available in `TestMain`.
- Use that path in `runTrellis`.
- Keep tests hermetic: no dependency on pre-existing `/tmp` binaries.

### 7. Canonicalize project identity

Choose one module/install path and make both `go.mod:1` and `README.md:46` agree. This is P3, not a prerequisite for P0/P1.

Default recommendation for this repository:

- Use `github.com/superops-team/trellis-go` if the public GitHub repo is the intended release path.
- Update internal imports accordingly.

The canonical module path for this repository is `github.com/superops-team/trellis-go`; do not retain the old module path in repository-owned imports or install instructions.

Do not make this change silently. The implementer must state the chosen identity in the PR description because this can affect external importers.

## Data Flow

```text
User command
    |
    v
Cobra command in cmd/trellis
    |
    v
resolvePaths(root flag, cwd)
    |
    +--> RepoRoot    -> git checks / platform dirs
    +--> TrellisDir  -> config, workflow, runtime
    +--> TasksDir    -> pkg/task.Manager
    +--> SpecDir     -> pkg/spec.Loader / context builder
```

Task command flow after P1:

```text
trellis task start <id>
    |
    v
resolvePaths
    |
    v
task.NewManager(paths.TasksDir)
    |
    v
Manager.Start(id)
    |
    v
task.json status: planning -> in_progress
```

## Test Strategy

| Slice | First failing test | Expected production change |
|-------|--------------------|----------------------------|
| P0 | `go test ./...` fails without `/tmp/trellis-test`; add harness test/build path. | Build test binary inside `TestMain`. |
| P0 | Root resolver tests for empty root, repo root, `.trellis` root. | Add unexported resolver helper. |
| P1 | E2E `task start` expects status change instead of placeholder output. | Wire CLI to `Manager.Start`. |
| P1 | E2E `task archive` expects archive path and completed status. | Wire CLI to `Manager.Archive` and fix metadata ordering. |
| P2 | Manifest malformed JSONL test fails with clear error. | Add `pkg/manifest` and update callers. |
| P2 | Context CLI build test expects real injected output. | Wire `context build` to builder. |
| P2 | Template render tests expect missing key errors and non-template copy. | Update `Engine.Render`. |

## Downward Compatibility

- Existing task directory layout stays unchanged: `MM-DD-name` for active tasks and `archive/YYYY-MM/MM-DD-name` for archived tasks.
- Existing manifest JSONL field names stay unchanged.
- Existing `--root <repo>/.trellis` test/user behavior remains supported.
- Existing success output prefix `Created task:` remains supported.
- Placeholder outputs may disappear only when replaced by real behavior.

## Failure Modes to Guard

- Root resolver treats `/tmp/repo/.trellis` as repo root and creates `/tmp/repo/.trellis/.trellis`.
- `init` rejects a valid future `.trellis` path because it checks existence too early.
- `task list` changes from directory-name output to ID output and breaks tests. Choose one output and assert it.
- `archive` moves the directory then fails to save updated metadata.
- `context add` writes absolute local machine paths, making task manifests non-portable.
- Template missing-key behavior breaks existing templates that rely on `<no value>`. Search embedded templates before enabling strict mode; if templates contain optional fields, make those fields explicit in the render context.

## Rollback Plan

- Root resolution and CLI delegation are ordinary code changes; revert the PR if behavior regresses.
- Manifest consolidation should preserve JSONL wire format, so rollback does not require data migration.
- Module path canonicalization is the riskiest part for external users; do it before tagged releases or split it into a separate PR if compatibility is uncertain.

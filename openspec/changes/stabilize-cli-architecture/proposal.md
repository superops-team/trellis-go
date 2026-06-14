# Stabilize CLI Architecture

## Why

Brooks review exposed a conceptual-integrity gap: the repository has reusable `pkg/*` packages, but the `cmd/trellis` CLI often bypasses them or leaves advertised behavior as stubs. This makes the public CLI, README, and package tests disagree about what Trellis actually does.

This change makes Trellis executable behavior match its package architecture and documentation before the project grows more features on top of unstable seams.

## Review Corrections Applied

This proposal was tightened after a full-dimension landing review. The important corrections are:

- Reduce the first implementation slice to root resolution, task CLI delegation, and hermetic tests. Context/template/spec-loader work remains in the same OpenSpec change but is sequenced after the foundation is green.
- Treat module-path canonicalization as a compatibility decision gate, not an automatic code change inside the first PR.
- Avoid introducing broad new architecture. The only allowed new package is a tiny `pkg/manifest` value package, and only if it removes the existing duplicate manifest parser without pulling in context-builder behavior.
- Make every requirement observable through a CLI command, package test, or file-format invariant.
- Make SDD/TDD executable by requiring failing tests before each behavior change.

## Current State

Verified on 2026-06-14.

| Area | Evidence | Current behavior | Gap |
|------|----------|------------------|-----|
| Root resolution | `cmd/trellis/main.go:30`, `cmd/trellis/task.go:38`, `cmd/trellis/uninstall.go:26` | `--root` is declared, but commands mostly use `os.Getwd()` and hardcoded `.trellis` paths. | Public flag and command behavior disagree. |
| Task lifecycle | `cmd/trellis/task.go:37`, `pkg/task/manager.go:46` | CLI manually creates task directories and JSON instead of using `pkg/task.Manager`. | Two task lifecycle implementations can drift. |
| Task start/archive | `cmd/trellis/task.go:91`, `cmd/trellis/task.go:105`, `pkg/task/manager.go:91`, `pkg/task/manager.go:105` | CLI prints placeholder messages while package methods exist. | Advertised lifecycle is not wired to the CLI. |
| Context commands | `cmd/trellis/context.go:21`, `cmd/trellis/context.go:33`, `pkg/context/builder.go:27` | CLI prints placeholder messages while context builder exists. | Context Builder feature is not executable through the CLI. |
| Manifest schema | `pkg/context/manifest.go:10`, `pkg/task/manifest.go:10` | Context manifest models and JSONL parsing are duplicated. | Schema updates require coordinated edits in two packages. |
| Template engine | `pkg/template/engine.go:33`, `pkg/template/engine.go:86`, `pkg/template/engine.go:169`, `pkg/template/engine.go:194` | Directory rendering templates all non-binary files; missing-key validation is a no-op; marshal error is ignored; `ShouldTemplate` is unused. | Template behavior is surprising and undertested. |
| Tests | `cmd/trellis/e2e_test.go:14` | `go test ./...` fails unless `/tmp/trellis-test` is built manually. | Normal contributor workflow is brittle. |
| Module identity | `go.mod:1`, `README.md:46` | Module path and README install path must both use `github.com/superops-team/trellis-go`. | Release/import identity must stay consistent. |

## Proposed Change

Ship the stabilization as staged slices under one OpenSpec change. Each slice must leave `go test ./...` green before the next starts.

| Slice | Scope | Why this order | May ship separately |
|-------|-------|----------------|---------------------|
| P0 | Hermetic E2E harness and root-resolution helper tests | Gives reliable feedback before refactoring command behavior. | Yes |
| P1 | Task CLI delegation: create/list/start/archive through `pkg/task.Manager` | Removes the biggest duplicate implementation and unlocks real lifecycle tests. | Yes |
| P2 | Context/manifest/template/spec-loader cleanup | Builds on stable root/task lookup; avoids mixing schema cleanup with task behavior changes. | Yes |
| P3 | README/module identity decision | Has external compatibility risk; should be performed only after intended public module identity is confirmed. | Yes |

The minimum valuable implementation is P0 + P1. P2 and P3 are follow-up slices unless the implementer can complete them without expanding the diff beyond the files listed in this spec.

## Development Start Scope

Start implementation with P0 + P1 only. This first development pass SHALL deliver:

- hermetic E2E tests with no `/tmp/trellis-test` precondition;
- one internal root-resolution helper used by `init`, `task create`, `task list`, and `uninstall`;
- `task create`, `task list`, `task start`, and `task archive` delegated to `pkg/task.Manager`;
- archive metadata written consistently before the directory is moved;
- `go test ./cmd/trellis ./pkg/task ./...` green.

Do not start P2/P3 until P0/P1 are complete and green.

## Non-Goals

- Do not redesign Trellis workflow phases.
- Do not add new AI platform support.
- Do not build a new template DSL.
- Do not change task file formats except where required to consolidate duplicate manifest code.
- Do not introduce a database or background service.
- Do not implement full session tracking for `task current`; keep a clear no-active-task behavior unless a separate session-tracking spec is created.
- Do not silently change the Go module path in the first implementation slice. Decide and document first, then change in P3 if still appropriate.

## Success Metrics

- `go test ./...` passes from a clean checkout without manually building `/tmp/trellis-test`.
- `trellis --root <repo-or-trellis-root> ...` has documented, consistent behavior across `init`, `task`, `context`, and `uninstall`.
- `task create`, `task list`, `task start`, and `task archive` execute through `pkg/task.Manager` behavior.
- `context build` exercises `pkg/context.Builder` through the CLI.
- Duplicate manifest schema code is removed or reduced to a single owned implementation.
- `pkg/spec`, `pkg/context`, `pkg/template`, and `pkg/git` have targeted tests for success and failure paths.

## Risk and Compatibility Notes

| Risk | Impact | Mitigation |
|------|--------|------------|
| `--root` meaning changes unexpectedly | Existing E2E and users passing `.trellis` paths could break. | Support both repo-root and `.trellis`-root forms; add tests for both. |
| CLI output changes break scripts | Users may parse current output. | Preserve existing success prefixes such as `Created task:`; only replace placeholder outputs for implemented commands. |
| Task archive changes move files differently | Existing task history could be misplaced. | Preserve `tasks/archive/YYYY-MM/<task-dir>` layout. Only change ordering/error handling. |
| Manifest consolidation changes JSONL format | Existing task manifests could become unreadable. | Keep exact JSONL line format and field names: `path`, `description`, `required`. |
| Module path update breaks imports | Downstream users may depend on the previous module identity. | Canonicalize on `github.com/superops-team/trellis-go`; update every import and README together. |

## SDD/TDD Landing Rules

- For every requirement in `specs/*/spec.md`, add or update a failing test first.
- Implement the smallest code change that makes that test pass.
- Refactor only after the test is green.
- Keep structural changes and behavior changes in separate commits where possible.
- Do not mark a task complete until `go test ./...` passes.

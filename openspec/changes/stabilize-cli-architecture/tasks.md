# Tasks: Stabilize CLI Architecture

## Execution Rule

Use SDD + TDD for every slice: pick one scenario from `specs/*/spec.md`, write or update the failing test, implement the smallest change, run the relevant package test, then run `go test ./...` before moving to the next slice.

## P0. Test harness and root resolution foundation

- [x] 0.1 Add a failing E2E harness test or reproduce `go test ./...` failure when `/tmp/trellis-test` is absent.
- [x] 0.2 Replace `/tmp/trellis-test` with a package-level test binary built by `TestMain` using `os.MkdirTemp`.
- [x] 0.3 Add table tests for root resolution: empty root, repo-root root, `.trellis` root. If adding an `os.Getwd` failure seam would require global state, skip that seam and verify wrapped `Getwd` errors through a small helper unit test instead.
- [x] 0.4 Add an unexported root-resolution helper in `cmd/trellis` returning repo root, Trellis dir, tasks dir, and spec dir.
- [x] 0.5 Update `init`, `task create`, `task list`, and `uninstall` to use the helper.
- [x] 0.6 Replace ignored `os.Getwd()` errors in task commands with explicit wrapped errors.
- [x] 0.7 Run `go test ./cmd/trellis` and `go test ./...`.

## P1. Task CLI delegation

- [x] 1.1 Add/adjust E2E coverage proving `task create` creates the same files through CLI behavior.
- [x] 1.2 Refactor `task create` to call `pkg/task.Manager.Create` and preserve `Created task:` output.
- [x] 1.3 Add/adjust E2E coverage proving `task list` excludes archived tasks and includes active tasks exactly once.
- [x] 1.4 Refactor `task list` to call `pkg/task.Manager.List`.
- [x] 1.5 Add failing E2E coverage for `task start <id>` changing `task.json` from `planning` to `in_progress`.
- [x] 1.6 Implement `task start <id>` through `pkg/task.Manager.Start`.
- [x] 1.7 Add failing E2E coverage for `task archive <id>` moving the task to `archive/YYYY-MM/` with `completed` metadata.
- [x] 1.8 Implement `task archive <id>` through `pkg/task.Manager.Archive`.
- [x] 1.9 Make archive metadata update ordering safe: save completed metadata before move or perform a temp-write/rollback-safe equivalent.
- [x] 1.10 Keep `task current` out of scope except for stable no-active-task output.
- [x] 1.11 Run `go test ./cmd/trellis ./pkg/task` and `go test ./...`.

## P2. Manifest and context behavior

- [x] 2.1 Add failing tests that prove malformed manifest JSONL returns an error instead of being replaced with an empty manifest.
- [x] 2.2 Introduce `pkg/manifest` with only `Entry`, `Manifest`, `Load`, and `Save`, preserving existing JSON fields.
- [x] 2.3 Update `pkg/task` and `pkg/context` to share `pkg/manifest` without changing JSONL file format.
- [x] 2.4 Add failing CLI tests for `trellis context add <file> --task <id> --phase implement|check`.
- [x] 2.5 Implement `context add` with `--task`, `--phase`, `--required`, and `--description` flags.
- [x] 2.6 Reject absolute paths and `..` segments in context manifest entries.
- [x] 2.7 Add failing CLI tests for `trellis context build --task <id> --phase implement|check` and `--phase research`.
- [x] 2.8 Implement `context build` through `pkg/context.Builder`.
- [x] 2.9 Convert silent required-data failures into actionable errors; optional entries may still be skipped.
- [x] 2.10 Run `go test ./cmd/trellis ./pkg/manifest ./pkg/context ./pkg/task` and `go test ./...`.

## P2b. Template and spec-loader cleanup

- [x] 2b.1 Add directory-rendering tests that call `Render`, not only `RenderString`.
- [x] 2b.2 Make `Render` use `ShouldTemplate` for non-binary files.
- [x] 2b.3 Configure template execution to fail on missing keys.
- [x] 2b.4 Return `.template-hashes.json` write errors instead of discarding them.
- [x] 2b.5 Add tests for binary copying, non-template text copying, unknown keys, and hash output.
- [x] 2b.6 Add tests for `pkg/spec.Loader` successful load, missing file errors, and unreadable package/layer behavior where practical on the target OS.
- [x] 2b.7 Run `go test ./pkg/template ./pkg/spec` and `go test ./...`.

## P3. Documentation and module identity

- [x] 3.1 Decide the canonical module path and record the decision in the PR description.
- [x] 3.2 If choosing `github.com/superops-team/trellis-go`, update `go.mod`, internal imports, and README install commands together.
- [x] 3.3 Confirm no repository-owned path references the old module path after canonicalization.
- [x] 3.4 Update README command examples only after the corresponding commands are implemented.
- [x] 3.5 Add tests for `pkg/git` wrapper success and command failure behavior.
- [x] 3.6 Run `go test ./...` and `go test -cover ./pkg/...`.

## Final verification

- [x] 4.1 Run `go test ./...` from a clean checkout state where `/tmp/trellis-test` is absent.
- [x] 4.2 Run `go test -cover ./pkg/...` and confirm no CLI-visible core package remains at 0% coverage.
- [x] 4.3 Manually smoke-test `trellis init`, `task create`, `task list`, `task start`, `task archive`, `context add`, and `context build` in a temp Git repository.
- [x] 4.4 Confirm no placeholder output remains for README-advertised commands that remain documented.
- [x] 4.5 Confirm existing task JSON and manifest JSONL files remain readable.

## Recommended schedule

- [x] Day 1: P0, hermetic E2E harness and root-resolution tests.
- [x] Day 2: P1, task create/list/start/archive delegation and task archive consistency.
- [x] Day 3: P2 manifest/context CLI if P1 is green; otherwise split context into a follow-up PR.
- [x] Day 4: P2b template/spec-loader tests and cleanup.
- [x] Day 5: P3 docs/module identity decision, full regression, and manual smoke test.

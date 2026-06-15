# Tasks: Agent Goal Mode Production Readiness

## Execution Rule

Use SDD + TDD for each scenario. Add or update the failing test first, implement the smallest behavior change, run the narrow test, then run `go test ./...` after each phase.

## P0.1 Hook command surface

- [x] 1.1 Add CLI tests proving `trellis hook` is registered on the root command.
- [x] 1.2 Add tests for `trellis hook inject-workflow-state --state implement` printing the implement prompt.
- [x] 1.3 Implement `cmd/trellis/hook.go` with `session-start`, `inject-context`, and `inject-workflow-state` subcommands.
- [x] 1.4 Run `go test ./cmd/trellis ./pkg/workflow`.

## P0.2 Platform init generates real hook files

- [x] 2.1 Add tests proving push-based `hook.Generator.GenerateAll` writes executable hook scripts.
- [x] 2.2 Implement push-based generation by composing `GenerateSessionStart`, `GenerateInjectContext`, and `GenerateInjectWorkflowState`.
- [x] 2.3 Add E2E coverage proving `trellis init --platform claude` creates Trellis-managed hook files under the platform config directory.
- [x] 2.4 Wire `runInit` to call `hook.NewGenerator(...).GenerateAll(...)` for each selected platform.
- [x] 2.5 Run `go test ./cmd/trellis ./pkg/hook`.

## P0.3 PRD goal guard

- [x] 3.1 Add package tests proving `Manager.Start` or the CLI start path rejects a task with missing or blank `prd.md`.
- [x] 3.2 Add E2E coverage proving `trellis task start <id>` fails before PRD content exists and succeeds after PRD content is written.
- [x] 3.3 Implement the PRD non-empty assertion for `task start`.
- [x] 3.4 Add tests proving `context build --phase implement|check --task <id>` fails when `prd.md` is blank.
- [x] 3.5 Implement the PRD assertion in context build paths shared by CLI and hooks.
- [x] 3.6 Run `go test ./cmd/trellis ./pkg/task ./pkg/context`.

## P0.4 Context safety guardrails

- [x] 4.1 Add builder tests for required oversized context entries failing with an actionable error.
- [x] 4.2 Add builder tests for required binary context entries failing with an actionable error.
- [x] 4.3 Add builder tests for required sensitive paths failing with an actionable error.
- [x] 4.4 Add builder tests proving optional missing/unsafe entries are listed under a skipped optional context section.
- [x] 4.5 Implement context entry size, binary, and sensitive-path checks.
- [x] 4.6 Run `go test ./pkg/context`.

## P0.5 End-to-end smoke and regression

- [x] 5.1 Add E2E coverage for `hook inject-context --task <id> --phase implement` returning the same goal-bearing context shape as `context build`.
- [x] 5.2 Manually smoke-test in a temp Git repository: `init --platform claude`, `task create`, write `prd.md`, `task start`, `context add`, `context build`, `hook inject-context`, and `hook inject-workflow-state`.
- [x] 5.3 Run `go test ./...`.
- [x] 5.4 Run `go test -race ./...` if the implementation touches shared file-writing or manifest code.

## P0.6 Spec review closure

- [x] 6.1 Add direct coverage for `hook session-start` deterministic output.
- [x] 6.2 Add direct coverage for missing `prd.md` preserving `planning` status on `task start` failure.
- [x] 6.3 Add hook coverage for blank-PRD implementation context failure and check-context success.
- [x] 6.4 Add direct coverage for optional oversized/sensitive context skips and the no-skip output shape.
- [x] 6.5 Run `openspec validate agent-goal-mode-production-readiness --strict`.

## Out-of-scope follow-ups

- [ ] F1 Specify `.trellis/.runtime/sessions` and make `task current` real.
- [ ] F2 Persist workflow state per task and enforce `plan -> implement -> verify -> finish` transitions.
- [ ] F3 Add manifest write locking or append-only concurrency protection.
- [ ] F4 Add CI and release gates for the full user journey.

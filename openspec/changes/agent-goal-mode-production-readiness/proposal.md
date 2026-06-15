# Agent Goal Mode Production Readiness

## Why

Brooks review found that Trellis is close to a useful local MVP, but its highest-value promise is not yet production-ready: an AI coding agent should receive a clear goal, task context, and workflow instruction automatically. Today the repository has several pieces of that model, but the executable chain is incomplete.

The P0 goal of this change is to make the smallest production-readiness slice real:

- `trellis init --platform ...` creates usable platform hook files instead of only creating placeholder directories.
- Generated hook scripts call real `trellis hook ...` commands.
- Implementation context refuses to run without a non-empty PRD, so agents do not start from an empty goal.
- Context injection has basic safety limits for file size, binary files, sensitive file names, and skipped optional entries.

This change intentionally does not attempt full session tracking, persistent workflow-state transitions, or a full release/CI system. Those are important follow-up slices, but the next implementation should first close the advertised hook/context/goal loop.

## Current State

Verified during Brooks review on 2026-06-15.

| Area | Evidence | Current behavior | Gap |
|------|----------|------------------|-----|
| Platform init | `cmd/trellis/init.go:135`, `cmd/trellis/init.go:141` | `init` creates platform directories and leaves template rendering as TODO. | Users can select platforms, but no usable hook integration is installed. |
| Hook generator | `pkg/hook/generator.go:35`, `pkg/hook/generator.go:37` | Push-based generation is a no-op. | Claude/Cursor-style hook integration has no generated config. |
| Hook commands | `pkg/hook/generator.go:87`, `pkg/hook/generator.go:100`, `pkg/hook/generator.go:110`, `cmd/trellis/main.go:34` | Generated scripts call `trellis hook ...`, but the root command does not register `hook`. | Generated scripts would fail if executed. |
| Goal guard | `pkg/task/manager.go:82`, `pkg/task/manager.go:104`, `pkg/context/builder.go:32` | `task create` creates an empty `prd.md`; `task start` does not require PRD content; context build silently omits empty PRD. | Agents may implement without a clear goal. |
| Context safety | `pkg/context/builder.go:133`, `pkg/fsutil/fsutil.go:116` | Context builder reads whole files into prompts and does not use existing binary detection. | Large, binary, or sensitive files can be injected into agent prompts. |

## Proposed Change

Ship a P0 production-readiness slice under three requirement areas:

| Area | Scope | Out of scope |
|------|-------|--------------|
| Agent hook integration | Add real `hook` CLI commands, wire `init` to hook generation, make generated hook scripts executable and testable. | Deep platform-specific automation for every supported editor. |
| Goal workflow enforcement | Require a non-empty PRD before starting or building implementation/check context. | New `Task.Goal` field, persistent workflow-state machine, session tracking. |
| Context safety | Add size, binary, sensitive-name rejection, and visible optional-skip warnings. | Full secret scanning, token budgeting, summarization, or external policy engines. |

The implementation should prefer one tracer-bullet path over broad platform support: one push-based platform, one pull-based platform, and one agentless platform are enough to prove `GenerateAll` dispatches to real outputs. Additional platform polish can follow after the P0 chain is green.

## Non-Goals

- Do not implement full `.trellis/.runtime/sessions` tracking.
- Do not add a new database, daemon, or background service.
- Do not redesign task statuses beyond the existing `planning`, `in_progress`, and `completed` lifecycle.
- Do not introduce a broad plugin system for hook templates.
- Do not implement full secret detection; use deterministic path/name deny rules for P0.
- Do not require every supported platform to have perfect native configuration in this slice.

## Success Metrics

- `go test ./...` passes.
- `trellis init --platform <push-based-platform>` creates at least one real hook/config file under that platform's config directory.
- `trellis hook session-start`, `trellis hook inject-context`, and `trellis hook inject-workflow-state` are registered commands with deterministic output.
- Every generated shell script references an existing `trellis hook` command.
- `trellis task start <id>` fails when `prd.md` is empty.
- `trellis context build --phase implement --task <id>` fails when `prd.md` is empty.
- Context build rejects binary, oversized, and sensitive-path entries with actionable errors for required entries.
- Optional skipped context entries are visible in output instead of disappearing silently.

## Risk and Compatibility Notes

| Risk | Impact | Mitigation |
|------|--------|------------|
| Existing users created empty PRDs and start tasks immediately. | `task start` will become stricter. | Error message should explain how to fix: write PRD content first. |
| Platform hook files may differ from future native integrations. | P0 generated config may be minimal. | Keep outputs simple and documented as Trellis-managed hook files. |
| Sensitive-file denylist blocks a legitimate context file. | User cannot inject that path directly. | Keep denylist narrow and deterministic; allow future explicit override outside P0. |
| Optional skip warnings change context stdout. | Agents may see extra metadata. | Emit a clearly delimited warning section after injected content. |

## SDD/TDD Landing Rules

- For each scenario in `specs/*/spec.md`, write or update a failing test first.
- Keep each production change scoped to one behavior: hook command registration, init wiring, PRD guard, or context safety.
- Run the narrow package test after each scenario, then `go test ./...` at the end of each phase.
- Do not mark P0 complete until a temp Git repository smoke test covers `init`, `task create`, PRD write, `task start`, `context build`, and the hook commands.

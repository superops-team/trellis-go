# Design: Agent Goal Mode Production Readiness

## Overview

The design keeps Trellis as a single local CLI. The P0 objective is not to build a full orchestration platform; it is to make the current promise executable: generated platform hooks can call real Trellis commands, and those commands inject a safe, goal-bearing context into an AI agent session.

The core rule is: **an implementation/check context is invalid unless it contains a non-empty PRD goal**. This gives the AI Agent goal mode a concrete safety rail without changing the task file format.

## Architecture

```text
trellis init --platform claude
    |
    v
platform registry -> hook.Generator.GenerateAll
    |
    v
platform config dir with Trellis-managed hook files
    |
    v
generated script calls: trellis hook inject-context
    |
    v
hook command resolves task/context/workflow and prints safe agent prompt content
```

## Component Decisions

### 1. Add `trellis hook` as a real CLI adapter

Add `cmd/trellis/hook.go` with these subcommands:

| Command | P0 behavior |
|---------|-------------|
| `hook session-start` | Print a deterministic session-start message and basic Trellis paths. It does not write session state in P0. |
| `hook inject-context --task <id> --phase implement|check|research` | Delegate to the same builder path as `context build`. For implement/check, require `--task`. |
| `hook inject-workflow-state --state plan|implement|verify|finish` | Use `pkg/workflow.Parser.InjectPrompt` and print the phase prompt. |

The hook commands should not duplicate context-building logic. They are adapters over existing package APIs and existing command path resolution.

### 2. Wire `init` to hook generation

`runInit` should call `hook.NewGenerator(platform, binary).GenerateAll(platformDir)` after creating each platform directory.

P0 binary selection should be simple and deterministic:

- Use `os.Args[0]` as the default binary path in generated scripts.
- If this later proves unstable for installed binaries, add a `--binary` flag in a separate change.

Push-based generation must stop being a no-op. For P0 it should generate at least:

- `session-start.sh`
- `inject-context.sh`
- `inject-workflow-state.sh`
- a small platform-local README or config marker only if the target platform needs a config file to discover scripts.

Pull-based and agentless generation can keep their existing minimal outputs, but tests must prove `GenerateAll` writes files for every supported platform class.

### 3. Enforce PRD as the P0 goal source

Do not add a new `Task.Goal` field in this slice. The existing `prd.md` is the goal source.

Rules:

- `task start <id>` SHALL fail if the task's `prd.md` is missing or contains only whitespace.
- `context build --phase implement|check --task <id>` SHALL fail if the task's `prd.md` is missing or blank.
- `hook inject-context --phase implement|check --task <id>` SHALL share the same failure behavior.
- `context build --phase research` SHALL remain task-free and does not require PRD.

This turns the current implicit behavior into an explicit contract: agents may research without a task goal, but they may not implement or verify without one.

### 4. Add context safety guardrails

Extend context entry loading with a narrow policy object or constants in `pkg/context`.

P0 defaults:

- Maximum single context entry size: 256 KiB.
- Binary file rejection: use existing binary detection behavior from `pkg/fsutil` or equivalent package-local helper.
- Sensitive path denylist: reject paths whose slash-normalized basename or path segment matches `.env`, `.env.*`, `id_rsa`, `id_ed25519`, `*.pem`, `*.key`, `credentials.json`, or contains `secret` / `token` as a full segment.

Failure behavior:

- Required unsafe entries fail the build with an actionable error naming the entry path and reason.
- Optional unsafe or missing entries are omitted, but the output includes a `=== skipped optional context ===` section listing path and reason.

This is not full data-loss prevention. It is a deterministic guard that prevents the most obvious prompt-injection and token-explosion failures.

## Data Flow

```text
task start <id>
    -> resolve task dir
    -> assert prd.md non-empty
    -> Manager.Start(id)

context build --phase implement --task <id>
    -> resolve task dir
    -> assert prd.md non-empty
    -> load manifest entries
    -> normalize path within .trellis root
    -> safety check entry
    -> include required content or fail
    -> include optional content or record skip warning
    -> print agent context

hook inject-context --phase implement --task <id>
    -> same builder path as context build
```

## Error Handling

- Empty PRD errors should include `PRD is required` and the task ID.
- Unsafe context errors should include the manifest entry path and one reason: `too large`, `binary file`, or `sensitive path`.
- Hook command argument errors should follow Cobra behavior and fail before reading task files.
- Optional skip warnings should be visible in stdout because hook output is consumed by agents.

## Test Strategy

| Phase | First failing tests | Expected production change |
|-------|---------------------|----------------------------|
| P0.1 hook command registration | Root command contains `hook`; `trellis hook inject-workflow-state --state implement` prints implement prompt. | Add `cmd/trellis/hook.go`. |
| P0.2 init wiring | `trellis init --platform claude` creates executable Trellis hook scripts. | Call `hook.Generator.GenerateAll` from init and implement push-based output. |
| P0.3 goal guard | `task start` fails when `prd.md` is empty; succeeds after PRD content is written. | Add PRD assertion before lifecycle start. |
| P0.4 context guard | Builder rejects required binary/oversized/sensitive entries and reports optional skips. | Add context safety policy to builder. |
| P0.5 hook/context integration | `hook inject-context --task <id> --phase implement` prints the same safe context as `context build`. | Share context build path. |

## Rollback Plan

- Hook commands and generated scripts are additive; revert the hook command file and init wiring if behavior regresses.
- PRD enforcement is stricter behavior; rollback by removing the PRD assertion if it blocks valid users.
- Context safety can be rolled back independently by restoring direct `os.ReadFile` entry loading, but keep tests documenting why the guard existed.

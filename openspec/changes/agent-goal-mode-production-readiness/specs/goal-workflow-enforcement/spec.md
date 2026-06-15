# Delta for Goal Workflow Enforcement

## ADDED Requirements

### Requirement: Implementation starts only with a non-empty PRD goal

Trellis SHALL treat task `prd.md` as the P0 goal source for implementation work.

#### Scenario: Start task with blank PRD
- GIVEN Trellis is initialized
- AND a task exists with an empty or whitespace-only `prd.md`
- WHEN the user runs `trellis task start <task-id>`
- THEN the command SHALL fail
- AND the task status SHALL remain `planning`
- AND the error SHALL contain `PRD is required`

#### Scenario: Start task with non-empty PRD
- GIVEN Trellis is initialized
- AND a task exists with a non-empty `prd.md`
- WHEN the user runs `trellis task start <task-id>`
- THEN the task status SHALL become `in_progress`
- AND the command SHALL print `Started task: <task-id>`

#### Scenario: Missing PRD file
- GIVEN a task directory exists without `prd.md`
- WHEN the user runs `trellis task start <task-id>`
- THEN the command SHALL fail
- AND the error SHALL contain `PRD is required`

### Requirement: Implementation and check context require a PRD goal

Trellis SHALL NOT build implementation or check context for a task unless the task PRD is present and non-empty.

#### Scenario: Build implementation context with blank PRD
- GIVEN Trellis is initialized
- AND a task exists with an empty or whitespace-only `prd.md`
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail
- AND the error SHALL contain `PRD is required`

#### Scenario: Build check context with blank PRD
- GIVEN Trellis is initialized
- AND a task exists with an empty or whitespace-only `prd.md`
- WHEN the user runs `trellis context build --task <task-id> --phase check`
- THEN the command SHALL fail
- AND the error SHALL contain `PRD is required`

#### Scenario: Build research context without PRD
- GIVEN Trellis is initialized
- WHEN the user runs `trellis context build --phase research`
- THEN the command SHALL succeed
- AND the command SHALL NOT require a task PRD

### Requirement: Hook context injection inherits PRD goal enforcement

Hook context injection SHALL enforce the same PRD goal rules as `context build`.

#### Scenario: Hook implementation context with blank PRD
- GIVEN Trellis is initialized
- AND a task exists with an empty or whitespace-only `prd.md`
- WHEN the user runs `trellis hook inject-context --task <task-id> --phase implement`
- THEN the command SHALL fail
- AND the error SHALL contain `PRD is required`

#### Scenario: Hook check context with non-empty PRD
- GIVEN Trellis is initialized
- AND a task exists with a non-empty `prd.md`
- WHEN the user runs `trellis hook inject-context --task <task-id> --phase check`
- THEN the command SHALL use the normal check-context builder behavior

## MODIFIED Requirements

### Requirement: Task lifecycle remains simple in P0

P0 goal enforcement SHALL NOT add new persisted task statuses.

#### Scenario: Start still uses existing status transition
- GIVEN a task has a non-empty PRD and status `planning`
- WHEN the user runs `trellis task start <task-id>`
- THEN the persisted status SHALL be `in_progress`
- AND no new persisted workflow-state field SHALL be required

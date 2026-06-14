# Delta for CLI Behavior

## ADDED Requirements

### Requirement: Centralized Trellis root resolution

All CLI commands SHALL resolve repository and Trellis paths through one shared root-resolution contract.

#### Scenario: Default root uses current repository
- GIVEN a Git repository at the process working directory
- AND no `--root` flag is provided
- WHEN a Trellis command needs paths
- THEN the repository root SHALL be the working directory
- AND the Trellis directory SHALL be `<working-directory>/.trellis`

#### Scenario: Root flag points to repository root
- GIVEN `--root /path/to/repo`
- WHEN a Trellis command needs paths
- THEN the repository root SHALL be `/path/to/repo`
- AND the Trellis directory SHALL be `/path/to/repo/.trellis`

#### Scenario: Root flag points to Trellis directory
- GIVEN `--root /path/to/repo/.trellis`
- WHEN a Trellis command needs paths
- THEN the repository root SHALL be `/path/to/repo`
- AND the Trellis directory SHALL be `/path/to/repo/.trellis`

#### Scenario: Working directory cannot be resolved
- GIVEN the process working directory cannot be resolved
- WHEN no `--root` flag is provided
- THEN the command SHALL fail with an error containing `get working directory`

#### Scenario: Init resolves a future Trellis directory
- GIVEN `/tmp/repo` is a Git repository
- AND `/tmp/repo/.trellis` does not exist yet
- WHEN the user runs `trellis --root /tmp/repo/.trellis init`
- THEN root resolution SHALL succeed before `.trellis` is created
- AND initialization SHALL create `/tmp/repo/.trellis`

### Requirement: CLI task commands delegate to task manager

Task lifecycle commands SHALL use `pkg/task.Manager` as the source of task behavior.

#### Scenario: Create task through CLI
- GIVEN Trellis is initialized
- WHEN the user runs `trellis task create add-auth`
- THEN the CLI SHALL create the task through `pkg/task.Manager.Create`
- AND the created task SHALL have status `planning`
- AND the command SHALL print `Created task: <path>`

#### Scenario: List tasks through CLI
- GIVEN Trellis has three active tasks
- WHEN the user runs `trellis task list`
- THEN the CLI SHALL retrieve tasks through `pkg/task.Manager.List`
- AND every active task ID SHALL appear in the output exactly once
- AND archived tasks SHALL NOT appear

#### Scenario: Start task through CLI
- GIVEN a task with status `planning`
- WHEN the user runs `trellis task start <task-id>`
- THEN the task status SHALL become `in_progress`
- AND the command SHALL print `Started task: <task-id>`

#### Scenario: Archive task through CLI
- GIVEN a task with status `in_progress`
- WHEN the user runs `trellis task archive <task-id>`
- THEN the task status SHALL become `completed`
- AND the task directory SHALL move under `tasks/archive/YYYY-MM/`
- AND the command SHALL print `Archived task: <task-id>`

#### Scenario: Start task with invalid status
- GIVEN a task is not in `planning` status
- WHEN the user runs `trellis task start <task-id>`
- THEN the command SHALL fail
- AND the error SHALL contain `invalid task status transition`

#### Scenario: Archive task with invalid status
- GIVEN a task is not in `in_progress` status
- WHEN the user runs `trellis task archive <task-id>`
- THEN the command SHALL fail
- AND the error SHALL contain `invalid task status transition`

### Requirement: Task archive preserves consistent metadata

Archiving a task SHALL NOT leave a moved task with stale lifecycle metadata.

#### Scenario: Archive save fails
- GIVEN task archiving starts
- AND writing updated task metadata fails
- WHEN the archive operation returns
- THEN the user SHALL receive an error
- AND the task SHALL NOT be silently reported as archived with stale metadata

### Requirement: Task current remains explicitly limited

Until session tracking is specified separately, `trellis task current` SHALL keep a stable no-active-task behavior and SHALL NOT imply that `task start` writes session state.

#### Scenario: Current task without session tracking
- GIVEN Trellis is initialized
- WHEN the user runs `trellis task current`
- THEN the command SHALL exit successfully
- AND the output SHALL state that no active task is available
- AND the command SHALL NOT create or modify task files

## MODIFIED Requirements

### Requirement: Init command respects root resolution

`trellis init` SHALL initialize Trellis at the resolved Trellis directory, not always at `cwd/.trellis`.

#### Scenario: Init with explicit Trellis root
- GIVEN a Git repository at `/tmp/repo`
- AND the user runs `trellis --root /tmp/repo/.trellis init`
- WHEN initialization succeeds
- THEN Trellis files SHALL be created under `/tmp/repo/.trellis`
- AND platform directories SHALL be created under `/tmp/repo`

### Requirement: Uninstall command respects root resolution

`trellis uninstall` SHALL remove files from the resolved Trellis directory, not always from `cwd/.trellis`.

#### Scenario: Uninstall with explicit root
- GIVEN Trellis is initialized under `/tmp/repo/.trellis`
- AND the current working directory is different from `/tmp/repo`
- WHEN the user runs `trellis --root /tmp/repo uninstall`
- THEN the command SHALL remove `/tmp/repo/.trellis`

# Delta for Agent Hook Integration

## ADDED Requirements

### Requirement: Hook commands are executable CLI behavior

Trellis SHALL register hook commands that generated platform scripts can call.

#### Scenario: Root command exposes hook namespace
- GIVEN the Trellis CLI is built
- WHEN the user runs `trellis hook --help`
- THEN the command SHALL exist
- AND help output SHALL list `session-start`, `inject-context`, and `inject-workflow-state`

#### Scenario: Inject workflow state prompt
- GIVEN Trellis is available on PATH
- WHEN the user runs `trellis hook inject-workflow-state --state implement`
- THEN the command SHALL succeed
- AND stdout SHALL contain the implement-phase workflow prompt

#### Scenario: Invalid workflow state
- GIVEN Trellis is available on PATH
- WHEN the user runs `trellis hook inject-workflow-state --state unknown`
- THEN the command SHALL fail
- AND the error SHALL identify the unknown state

### Requirement: Init generates Trellis-managed hook files

`trellis init --platform <id>` SHALL install concrete hook files for supported platform classes instead of only creating empty platform directories.

#### Scenario: Init generates push-based hook scripts
- GIVEN `/tmp/repo` is a Git repository
- AND Trellis is not initialized
- WHEN the user runs `trellis --root /tmp/repo init --platform claude`
- THEN initialization SHALL create the Claude platform config directory
- AND the platform config directory SHALL contain Trellis-managed hook files
- AND those hook files SHALL reference existing `trellis hook` subcommands

#### Scenario: Push-based generator writes executable scripts
- GIVEN a push-based platform
- WHEN `hook.Generator.GenerateAll` is called
- THEN it SHALL write `session-start.sh`, `inject-context.sh`, and `inject-workflow-state.sh`
- AND each script SHALL be executable
- AND each script SHALL call a registered `trellis hook` command

#### Scenario: Pull-based and agentless generation remains supported
- GIVEN a pull-based or agentless platform
- WHEN `hook.Generator.GenerateAll` is called
- THEN it SHALL write at least one Trellis-managed integration file
- AND it SHALL NOT silently succeed without writing any files

### Requirement: Hook context injection delegates to context builder

Hook context injection SHALL reuse the same behavior as `trellis context build`.

#### Scenario: Hook injects implementation context
- GIVEN Trellis is initialized
- AND a task has a non-empty `prd.md`
- AND the task implementation manifest references a required file that exists
- WHEN the user runs `trellis hook inject-context --task <task-id> --phase implement`
- THEN stdout SHALL contain the Trellis injection marker
- AND stdout SHALL contain the PRD content
- AND stdout SHALL contain the referenced file content

#### Scenario: Hook inject context requires task for implementation
- GIVEN Trellis is initialized
- WHEN the user runs `trellis hook inject-context --phase implement`
- THEN the command SHALL fail
- AND the error SHALL explain that `--task` is required

#### Scenario: Hook research context does not require task
- GIVEN Trellis is initialized
- WHEN the user runs `trellis hook inject-context --phase research`
- THEN the command SHALL succeed
- AND the command SHALL NOT require `--task`

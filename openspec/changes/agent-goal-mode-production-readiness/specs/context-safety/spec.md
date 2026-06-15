# Delta for Context Safety

## ADDED Requirements

### Requirement: Context entries are size-limited

Trellis SHALL prevent single context entries from injecting unexpectedly large files into agent prompts.

#### Scenario: Required context entry exceeds maximum size
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references a required file larger than the configured maximum entry size
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail
- AND the error SHALL identify the entry path
- AND the error SHALL contain `too large`

#### Scenario: Optional context entry exceeds maximum size
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references an optional file larger than the configured maximum entry size
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL succeed
- AND stdout SHALL contain a skipped optional context section naming the entry path and `too large`

### Requirement: Binary context entries are rejected

Trellis SHALL NOT inject binary files into agent prompts.

#### Scenario: Required binary context entry
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references a required binary file
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail
- AND the error SHALL identify the entry path
- AND the error SHALL contain `binary file`

#### Scenario: Optional binary context entry
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references an optional binary file
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL succeed
- AND stdout SHALL contain a skipped optional context section naming the entry path and `binary file`

### Requirement: Sensitive context paths are rejected

Trellis SHALL reject obvious sensitive path names before reading file content.

#### Scenario: Required sensitive context path
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references a required path such as `.env`, `id_rsa`, `credentials.json`, `private.key`, or `secrets/token.txt`
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail before injecting file content
- AND the error SHALL identify the entry path
- AND the error SHALL contain `sensitive path`

#### Scenario: Optional sensitive context path
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references an optional sensitive path
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL succeed
- AND stdout SHALL contain a skipped optional context section naming the entry path and `sensitive path`

### Requirement: Optional context skips are visible

Trellis SHALL make skipped optional context entries visible to the consuming agent.

#### Scenario: Optional context entry missing
- GIVEN a task has a non-empty PRD
- AND `implement.jsonl` references an optional file that does not exist
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL succeed
- AND stdout SHALL contain a skipped optional context section naming the missing entry path

#### Scenario: No optional entries skipped
- GIVEN all optional context entries are safe and readable
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN stdout SHALL NOT contain a skipped optional context section

## MODIFIED Requirements

### Requirement: Existing path traversal protection remains enforced

Context safety checks SHALL be applied after manifest paths are normalized and confirmed to stay inside the Trellis root.

#### Scenario: Unsafe traversal path remains rejected
- GIVEN `implement.jsonl` references `../secret.txt`
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail or skip according to the entry's required flag
- AND the file outside the Trellis root SHALL NOT be read

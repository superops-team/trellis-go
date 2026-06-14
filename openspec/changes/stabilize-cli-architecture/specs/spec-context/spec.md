# Delta for Spec and Context Behavior

## ADDED Requirements

### Requirement: Single manifest schema owner

Trellis SHALL define context manifest entries and JSONL load/save behavior in one package used by task and context code.

#### Scenario: Add context entry from task manager
- GIVEN a task manifest file exists
- WHEN task code appends a context entry
- THEN the same manifest entry type SHALL be used by context-building code
- AND no duplicate manifest parser SHALL be required in `pkg/task` and `pkg/context`

#### Scenario: Malformed manifest line
- GIVEN a manifest contains malformed JSONL
- WHEN the manifest is loaded
- THEN loading SHALL fail with an error containing the manifest path or line context
- AND the caller SHALL NOT silently replace the malformed manifest with an empty one

#### Scenario: Existing manifest remains readable
- GIVEN an existing manifest line contains `path`, `description`, and `required`
- WHEN the shared manifest loader reads the file
- THEN the entry SHALL be loaded without changing field names
- AND saving the manifest SHALL preserve one JSON object per line

### Requirement: Context add command writes manifest entries

`trellis context add <file>` SHALL add a file reference to a task context manifest.

#### Scenario: Add required implementation context
- GIVEN Trellis is initialized
- AND a task `user-auth` exists
- AND `.trellis/spec/auth.md` exists
- WHEN the user runs `trellis context add spec/auth.md --task user-auth --phase implement --required --description "Auth spec"`
- THEN `spec/auth.md` SHALL be appended to the task's `implement.jsonl`
- AND the entry SHALL include `required: true`
- AND the entry SHALL include the supplied description

#### Scenario: Add context outside Trellis workspace
- GIVEN Trellis is initialized
- WHEN the user runs `trellis context add ../secret.txt --task user-auth --phase implement`
- THEN the command SHALL fail
- AND no manifest entry SHALL be written

#### Scenario: Add absolute context path
- GIVEN Trellis is initialized
- WHEN the user runs `trellis context add /tmp/secret.txt --task user-auth --phase implement`
- THEN the command SHALL fail
- AND no manifest entry SHALL be written

### Requirement: Context build command prints real builder output

`trellis context build` SHALL call `pkg/context.Builder` and print generated context to stdout.

#### Scenario: Build implementation context
- GIVEN a task has non-empty `prd.md`
- AND `implement.jsonl` references a required file that exists
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN stdout SHALL contain the Trellis injection marker
- AND stdout SHALL contain the PRD content
- AND stdout SHALL contain the referenced file content

#### Scenario: Required context entry is missing
- GIVEN `implement.jsonl` references a required file that does not exist
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL fail
- AND the error SHALL identify the missing required entry path

#### Scenario: Optional context entry is missing
- GIVEN `implement.jsonl` references an optional file that does not exist
- WHEN the user runs `trellis context build --task <task-id> --phase implement`
- THEN the command SHALL succeed
- AND the missing optional file SHALL be omitted from stdout

#### Scenario: Research context does not require task
- GIVEN Trellis is initialized
- WHEN the user runs `trellis context build --phase research`
- THEN the command SHALL succeed
- AND the command SHALL NOT require `--task`

### Requirement: Template rendering has explicit file handling

Template directory rendering SHALL distinguish binary files, known template file extensions, and non-template text files.

#### Scenario: Template file contains unknown key
- GIVEN a rendered template file contains `{{.UnknownKey}}`
- WHEN `Render` is called
- THEN rendering SHALL fail with a missing-key error

#### Scenario: Non-template text file is rendered
- GIVEN a text file has an extension not listed in `TemplateFileExtensions`
- WHEN `Render` is called
- THEN the file SHALL be copied unchanged
- AND template placeholders in that file SHALL NOT be evaluated

#### Scenario: Template hash output cannot be written
- GIVEN the destination cannot accept `.template-hashes.json`
- WHEN `Render` is called
- THEN `Render` SHALL return an error instead of discarding the failure

## MODIFIED Requirements

### Requirement: Spec loader reports data loss risks

Spec loading SHALL report unreadable required locations rather than silently dropping them where doing so hides corrupt data.

#### Scenario: Package layer file cannot be read
- GIVEN a package layer index exists but cannot be read
- WHEN `LoadPackage` is called
- THEN the loader SHALL return an error identifying the unreadable file

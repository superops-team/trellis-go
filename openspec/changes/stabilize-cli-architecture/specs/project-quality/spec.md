# Delta for Project Quality

## ADDED Requirements

### Requirement: E2E tests are hermetic

The full test suite SHALL be runnable from a clean checkout without requiring a manually prebuilt CLI binary.

#### Scenario: Run full test suite from clean checkout
- GIVEN the repository is freshly checked out
- AND `/tmp/trellis-test` does not exist
- WHEN a contributor runs `go test ./...`
- THEN the test suite SHALL build any required test binary itself
- AND the test suite SHALL pass or fail based on product behavior, not missing external binaries

#### Scenario: Test binary path is test-owned
- GIVEN E2E tests need a compiled Trellis binary
- WHEN the tests run
- THEN the binary path SHALL be created by the test process
- AND the path SHALL NOT be hardcoded to `/tmp/trellis-test`

### Requirement: Core packages have targeted failure-path tests

Packages that own CLI-visible behavior SHALL have tests for both success paths and actionable failure paths.

#### Scenario: Spec loader tests
- GIVEN `pkg/spec` loads specs from disk
- WHEN package tests run
- THEN tests SHALL cover successful index/load behavior
- AND missing or unreadable files SHALL be covered

#### Scenario: Context builder tests
- GIVEN `pkg/context` builds command output from manifests
- WHEN package tests run
- THEN required entry success and failure SHALL be covered
- AND optional missing entries SHALL be covered
- AND malformed manifests SHALL be covered

#### Scenario: Template render tests
- GIVEN `pkg/template` renders embedded template directories
- WHEN package tests run
- THEN tests SHALL call `Render`
- AND unknown key, binary copy, non-template copy, and hash output behavior SHALL be covered

#### Scenario: Git wrapper tests
- GIVEN `pkg/git` wraps Git commands
- WHEN package tests run
- THEN command success and command failure behavior SHALL be covered

### Requirement: Public docs only advertise implemented CLI behavior

README command examples SHALL match commands that execute real behavior.

#### Scenario: README lists task lifecycle commands
- GIVEN README documents `trellis task start` and `trellis task archive`
- WHEN a contributor runs those commands in a valid Trellis repository
- THEN the commands SHALL perform the documented lifecycle transition
- AND they SHALL NOT print placeholder text such as `not yet implemented`

#### Scenario: README lists context commands
- GIVEN README documents `trellis context add` and `trellis context build`
- WHEN a contributor runs those commands with valid arguments
- THEN the commands SHALL perform the documented context operation
- AND they SHALL NOT print placeholder text such as `not yet implemented`

### Requirement: Canonical module identity is consistent

The Go module path, internal imports, and README installation command SHALL use one canonical repository identity.

#### Scenario: Module path matches README install path
- GIVEN `go.mod` declares the module path
- AND README provides a `go install` command
- WHEN both are inspected
- THEN they SHALL refer to the same module identity

#### Scenario: Module path change is explicit
- GIVEN the chosen module identity differs from the current `go.mod` value
- WHEN the implementation changes `go.mod`
- THEN all internal imports SHALL be updated in the same PR
- AND the PR description SHALL mention the compatibility impact

## MODIFIED Requirements

### Requirement: Coverage signal reflects protected behavior

Coverage checks SHALL be used as a smoke signal for unprotected core behavior, not as the only success criterion.

#### Scenario: Run package coverage
- GIVEN stabilization work is complete
- WHEN `go test -cover ./pkg/...` is run
- THEN no CLI-visible core package SHALL remain at 0% coverage
- AND any package below 50% coverage SHALL have an explicit reason or follow-up issue

### Requirement: SDD and TDD workflow remains enforceable

Implementation SHALL map each OpenSpec scenario to at least one automated or manual verification step.

#### Scenario: Requirement implemented without test
- GIVEN a scenario in this OpenSpec change is implemented
- WHEN the task is marked complete
- THEN there SHALL be a corresponding automated test or documented manual smoke-test step
- AND `go test ./...` SHALL pass

# Contributing to Trellis

Thank you for considering contributing to Trellis! This guide will help you get started.

## Development Environment

### Prerequisites

- Go 1.23 or later
- Git
- Make (optional, for convenience targets)

### Setup

```bash
# Clone the repository
git clone git@github.com:superops-team/trellis-go.git
cd trellis-go

# Run tests
go test ./...

# Build
go build ./cmd/trellis

# Install locally
go install ./cmd/trellis
```

## Code Style

Trellis follows standard Go conventions:

- **Formatting**: `gofmt` (enforced by CI)
- **Linting**: `golangci-lint` with default config
- **Imports**: Standard library → third-party → internal (grouped with blank lines)
- **Naming**: Go conventions (camelCase, PascalCase for exported)
- **Comments**: Doc comments on all exported symbols
- **Errors**: Descriptive error messages with context (use `fmt.Errorf` with `%w` for wrapping)

### Additional Conventions

- **Package names**: Single word, lowercase, no underscores
- **Test files**: `_test.go` suffix, test functions named `TestXxx`
- **Test helpers**: `internal/testutil` package for shared test utilities
- **Exported functions**: Must have doc comments
- **Error handling**: Check errors, don't use `_` to discard errors
- **Concurrency**: Use `sync` package, prefer RWMutex for read-heavy patterns

## Pull Request Process

1. **Fork the repository** and create a feature branch from `main`
2. **Write tests** for new functionality (aim for ≥ 70% coverage on new code)
3. **Run the full test suite** before submitting:
   ```bash
   go test ./...
   ```
4. **Keep PRs focused** — one feature or fix per PR
5. **Write a clear PR description** explaining what and why
6. **Update documentation** if your change affects the public API or behavior
7. **CI must pass** before merge

### PR Title Convention

```
<type>: <short description>

Types: feat, fix, refactor, test, docs, chore, perf
```

Examples:
- `feat: add ZCode platform support`
- `fix: resolve race condition in registry lookup`
- `docs: add configuration reference`

## Issue Templates

### Bug Report

When filing a bug report, include:
- Trellis version (`trellis --version`)
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error output

### Feature Request

When requesting a feature, include:
- Problem statement (what doesn't work or what's missing)
- Proposed solution
- Alternative approaches considered
- Use case description

## Testing

Trellis uses Go's standard testing package with the following conventions:

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test package-level behavior with real file I/O
- **No external dependencies**: Tests should not require network access
- **Parallel-safe**: Use `t.Parallel()` where safe
- **Test data**: Use `t.TempDir()` for temporary files

### Running Tests

```bash
# All tests
go test ./...

# With race detector
go test -race ./...

# With coverage
go test -cover ./...

# Specific package
go test ./pkg/platform/...
```

## Documentation

- Code changes should include doc comments
- User-facing changes should update the docs in `docs/`
- Spec changes go in `.trellis/spec/`
- Skill changes go in `.agents/skills/`

## Getting Help

- Open an issue on GitHub
- Check existing issues and discussions
- Review the documentation in `docs/`

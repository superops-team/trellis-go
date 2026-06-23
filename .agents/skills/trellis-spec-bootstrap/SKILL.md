---
name: trellis-spec-bootstrap
description: |
  Bootstrap a Trellis spec from an existing codebase.
  Use when you need to generate gap analysis specs, PRDs, or architecture
  documents from existing code without starting from scratch.
---

# Trellis Spec Bootstrap

## Trigger Check

Use this skill when:
- You're starting a new Trellis project from an existing codebase
- You need to document an existing system's architecture
- You want to generate gap analysis specs for a codebase
- You need a PRD for an existing feature

## Steps

### 1. Analyze Codebase

```bash
# Get project overview
find . -type f -name "*.go" -o -name "*.ts" -o -name "*.py" | head -50

# Check package/module structure
go list ./...  # Go
cat package.json  # Node
cat Cargo.toml  # Rust
```

### 2. Identify Components

List all:
- Packages/modules
- Entry points
- Public APIs
- Configuration files
- Test files

### 3. Generate Spec Structure

```bash
mkdir -p .trellis/spec/gap-analysis
```

For each component, create a spec file:
```markdown
# Component Name

## Current State
...

## Requirements
...

## Acceptance Criteria
- [ ] ...
```

### 4. Link to Tasks

```bash
trellis task create <component-name>
trellis context add .trellis/spec/gap-analysis/<component>.md --task <component-name> --phase implement
```

## Output

- `.trellis/spec/gap-analysis/` directory with spec files for each component
- Initial tasks for each gap
- Architecture overview document

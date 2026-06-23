# Installation & First Task

## Prerequisites

- Go 1.23 or later
- Git

## Installation

### Via go install

```bash
go install github.com/superops-team/trellis-go/cmd/trellis@latest
```

### Via pre-built binary

Download from [GitHub Releases](https://github.com/superops-team/trellis-go/releases).

### Verify installation

```bash
trellis --help
```

## First Task

### 1. Initialize a project

```bash
git init my-project && cd my-project
trellis init --developer alice --platform claude
```

This creates:
- `.trellis/` — Trellis data directory (specs, tasks, workspace)
- `.claude/` — Claude Code hook files (for push-based platforms)

### 2. Create a task

```bash
trellis task create user-auth
```

Output: `Created task: /path/to/.trellis/tasks/MM-DD-user-auth`

### 3. Write a PRD

Every task needs a PRD before starting:

```bash
echo "# User Authentication
## Requirements
- Login with email/password
- Session management
- Password reset flow
## Acceptance Criteria
- [ ] Users can register and login
- [ ] Sessions persist across page reloads
- [ ] Password reset emails are sent" > .trellis/tasks/*-user-auth/prd.md
```

### 4. Start the task

```bash
trellis task start user-auth
```

Output: `Started task: user-auth`

### 5. Add context files

```bash
trellis context add docs/architecture.md --task user-auth --phase implement
```

### 6. Build context for your AI agent

```bash
trellis context build --task user-auth --phase implement
```

This outputs the assembled context your AI agent will use.

## Next Steps

- [How Trellis Works](how-it-works.md)
- [Everyday Use](everyday-use.md)

# Everyday Use

## Daily Workflow

### Morning: Check active tasks

```bash
trellis task list
trellis task current
```

### Starting new work

```bash
# Create task
trellis task create fix-login-bug

# Write PRD
vim .trellis/tasks/*-fix-login-bug/prd.md

# Start
trellis task start fix-login-bug
```

### During development

```bash
# Add relevant context files
trellis context add spec/auth.md --task fix-login-bug --phase implement
trellis context add docs/api.md --task fix-login-bug --phase implement

# Build context for your AI agent
trellis context build --task fix-login-bug --phase implement
```

### Finishing work

```bash
# Archive completed task
trellis task archive fix-login-bug

# Record session
trellis hook record-session --title "Fixed login bug" --commit abc1234
```

## Common Patterns

### Multi-platform setup

```bash
trellis init --developer alice --platform claude --platform cursor --platform codex
```

### Adding a platform to existing project

```bash
trellis init --platform zcode
```

### Viewing session history

```bash
trellis hook list-sessions
```

### Updating Trellis

```bash
trellis upgrade
```

### Syncing templates

```bash
trellis update
```

## Tips

- **Write PRDs first** — Trellis requires a non-empty PRD before `task start`
- **Use context add** — Tell your AI agent which files matter for each phase
- **Archive completed tasks** — Keeps `tasks/` clean and builds session history
- **Commit `.trellis/`** — Specs and task metadata belong in version control
- **Don't commit `.trellis/.developer`** — It's gitignored, per-machine identity

---
name: trellis-meta
description: |
  Trellis custom core skill for self-referential operations.
  Use when you need to inspect, update, or debug Trellis's own configuration,
  tasks, specs, or session data.
---

# Trellis Meta

## Trigger Check

Use this skill when:
- You need to inspect `.trellis/` directory structure
- A Trellis command is not working as expected
- You need to manually fix a corrupted task or session
- You're debugging platform hook issues
- You want to understand Trellis's internal state

## Steps

### 1. Inspect State

```bash
# Check current task
trellis task current

# List all tasks
trellis task list

# List all platforms
trellis init --help  # Shows available platforms
```

### 2. Check File Structure

```bash
ls -la .trellis/
ls -la .trellis/tasks/
cat .trellis/config.yaml
```

### 3. Debug Platform Hooks

```bash
# Check which hooks exist for current platform
ls -la .claude/hooks/  # or .cursor/, .codex/, etc.
```

### 4. Manual Fixes

If a task is stuck in `planning` state:
```bash
# Edit task.json directly
vim .trellis/tasks/*-task-name/task.json
# Change status to "in_progress" or "completed"
```

If context files are missing:
```bash
trellis context add spec/current.md --task <task> --phase implement
trellis context build --task <task> --phase implement
```

## Output

- Current Trellis state summary
- Debug information for troubleshooting
- Manual fix instructions when automated commands fail

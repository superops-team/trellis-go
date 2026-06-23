---
name: mem-recall
description: |
  Cross-platform AI conversation recall skill.
  Use when you need to retrieve past AI agent conversations, session journals,
  or task history across platforms using `trellis mem`.
---

# Mem Recall

## Trigger Check

Use this skill when:
- You need to recall a previous AI agent conversation
- You want to find a specific decision or discussion from a past session
- You need context from a task that was archived weeks ago
- You're picking up work after a long break

## Steps

### 1. Search Session Journals

```bash
trellis hook list-sessions
```

Search by:
- Date range: `--since "2026-06-01" --until "2026-06-23"`
- Task name: `--task user-auth`
- Keyword: `--search "token refresh"`

### 2. Read Journal

```bash
trellis hook read-session --id <session-id>
```

### 3. Cross-Reference with Tasks

```bash
trellis task list --archived
trellis task show <task-name>
```

### 4. Rebuild Context

If resuming work on an archived task:

```bash
trellis task unarchive <task-name>
trellis context build --task <task-name> --phase implement
```

## Output

- Session journal entries matching the search criteria
- Reconstructed task context for resuming work
- Timeline of decisions and changes

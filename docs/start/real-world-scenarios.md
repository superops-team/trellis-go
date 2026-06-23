# Real-World Scenarios

## Scenario 1: Solo Developer

Alice works alone on a side project. She uses Trellis with Claude Code.

```bash
# Setup
git init my-app && cd my-app
trellis init --developer alice --platform claude

# Feature work
trellis task create add-dark-mode
# Write PRD, start task, implement with Claude Code
trellis task archive add-dark-mode
```

**Key benefit**: Alice's Claude Code sessions always have the right context — PRD, specs, and relevant files are injected automatically.

## Scenario 2: Team with Multiple Platforms

Bob uses Cursor, Carol uses Codex. Trellis generates config for both.

```bash
trellis init --developer bob --platform cursor
trellis init --developer carol --platform codex
```

**Key benefit**: Each developer's AI agent gets platform-native configuration, but they share the same specs and task tracking.

## Scenario 3: Complex Multi-Phase Feature

A feature requires research, implementation, and verification phases.

```bash
# Research phase
trellis task create new-auth-system
trellis context add research/auth-options.md --task new-auth-system --phase research
trellis context build --task new-auth-system --phase research

# Implementation phase
trellis task start new-auth-system
trellis context add spec/auth.md --task new-auth-system --phase implement
trellis context build --task new-auth-system --phase implement

# Verification phase
trellis context add spec/auth.md --task new-auth-system --phase check
trellis context build --task new-auth-system --phase check

# Done
trellis task archive new-auth-system
```

**Key benefit**: Each phase gets exactly the context it needs — no more, no less.

## Scenario 4: Bug Fix with Root Cause Analysis

```bash
trellis task create fix-token-refresh
# The trellis-break-loop skill helps analyze root causes
trellis context build --task fix-token-refresh --phase implement
# Fix the bug
trellis task archive fix-token-refresh
trellis hook record-session --title "Fixed token refresh race condition" --commit def5678
```

**Key benefit**: Session journal tracks what was fixed and why, building institutional knowledge.

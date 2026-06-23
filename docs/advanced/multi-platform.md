# Multi-Platform Configuration

## Why Multiple Platforms?

Different AI coding agents have different strengths. Trellis supports 16+ platforms so you can:

- Use different agents for different tasks
- Let team members use their preferred platform
- Compare agent performance on the same task

## Initializing with Multiple Platforms

```bash
# Specify platforms explicitly
trellis init --developer alice --platform claude --platform cursor --platform codex

# Or configure all platforms
trellis init --developer alice --all
```

## Adding a Platform Later

```bash
trellis init --platform zcode
```

This generates the platform-specific hook files without reinitializing existing config.

## Platform Hook Files

Each platform gets its own hook directory:

```
.claude/
├── hooks/
│   ├── task-create
│   ├── task-start
│   ├── task-archive
│   └── record-session

.cursor/
├── rules/
│   └── trellis.mdc

.codex/
├── agents/
│   └── trellis.md

.devin/
├── workflow/
│   └── trellis.yaml
└── skills/
    └── trellis.md
```

## Platform-Specific Features

| Feature | Push-Based | Pull-Based | Agentless |
|---------|:----------:|:----------:|:---------:|
| Auto hooks | ✅ | ❌ | ❌ |
| Agent skills | ✅ | ✅ | ❌ |
| Workflow files | ❌ | ❌ | ✅ |
| Context injection | ✅ | ✅ | ❌ |

## Best Practices

1. **Use `--all` for team projects** — Everyone gets their platform config
2. **Use `--platform` for personal projects** — Only what you need
3. **Don't commit platform dirs to `.gitignore`** — They're per-developer
4. **Commit `.trellis/`** — Shared specs and task tracking

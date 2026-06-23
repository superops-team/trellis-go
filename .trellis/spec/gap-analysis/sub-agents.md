# sub-agents: 子代理系统

> 创建 3 个 Trellis 子代理定义，支持多平台格式。

## 当前状态

| 组件 | 状态 |
|------|:---:|
| 子代理定义 | ❌ |
| 子代理模板 | ❌ |
| 多平台格式适配 | ❌ |
| Context injection (PreToolUse) | ❌ |

## 需求

### 1. trellis-implement（代码实现代理）

**职责**:
1. 读取 `implement.jsonl` 获取上下文文件列表
2. 读取 `prd.md`、`design.md`（如果存在）、`implement.md`（如果存在）
3. 按照 spec 和设计文档编写代码
4. 不执行 git commit（由主 session 控制）

**Claude Code 格式**:
```markdown
---
name: trellis-implement
description: |
  Implements code changes according to the active Trellis task's PRD, design, and implement plan.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# trellis-implement

You are the trellis-implement sub-agent in the Trellis workflow.

## Context Loading

1. Read `implement.jsonl` for the file manifest
2. Read `prd.md` for requirements
3. Read `design.md` if present
4. Read `implement.md` if present

## Implementation Rules

- Follow the project's coding conventions (read from `.trellis/spec/`)
- Write minimal, focused changes
- Do not commit — the main session handles git
- Report what was changed and why
```

### 2. trellis-check（验证代理）

**职责**:
1. `git diff --name-only HEAD` 获取变更
2. 发现适用的 spec 层
3. 逐层对比 diff 和 spec 的质量检查清单
4. 运行 lint/typecheck/test
5. 发现问题自动修复（最多 3 轮）
6. 报告通过/失败状态

**Claude Code 格式**:
```markdown
---
name: trellis-check
description: |
  Code quality check expert. Reviews diffs against specs, runs lint/typecheck/test, self-fixes.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# trellis-check

## Flow

1. `git diff --name-only HEAD` → find changed files
2. Discover applicable spec layers from `.trellis/spec/`
3. Compare diff against each layer's quality checklist
4. Run lint, typecheck, and tests
5. If issues found: fix → re-verify (max 3 rounds)
6. Report: PASSED or FAILED with details
```

### 3. trellis-research（调研代理）

**职责**:
1. 只读模式搜索代码库
2. 分析代码结构、依赖关系、模式
3. 输出调研报告

**Claude Code 格式**:
```markdown
---
name: trellis-research
description: |
  Read-only codebase research agent. Searches, analyzes structure, and reports findings.
tools: Read, Glob, Grep, Bash
---

# trellis-research

Read-only agent. No file writes.

## Flow

1. Read `research.jsonl` for context
2. Search codebase for relevant patterns
3. Analyze and report findings
```

### 4. 多平台格式

| 平台 | 格式 | 文件路径 |
|------|------|---------|
| Claude Code | YAML frontmatter + Markdown | `.claude/agents/{name}.md` |
| Cursor | YAML frontmatter + Markdown | `.cursor/agents/{name}.md` |
| OpenCode | YAML frontmatter + `permission:` object | `.opencode/agents/{name}.md` |
| Codex | TOML | `.codex/agents/{name}.toml` |
| Kiro | JSON | `.kiro/agents/{name}.json` |
| Gemini CLI | YAML frontmatter + Markdown (pull-based prelude) | `.gemini/agents/{name}.md` |
| Qoder | YAML frontmatter + Markdown (pull-based prelude) | `.qoder/agents/{name}.md` |
| CodeBuddy | YAML frontmatter + Markdown | `.codebuddy/agents/{name}.md` |
| Copilot | YAML frontmatter + Markdown | `.github/agents/{name}.agent.md` |
| Droid | YAML frontmatter + Markdown | `.factory/droids/{name}.md` |
| Pi Agent | YAML frontmatter + `model`/`thinking` | `.pi/agents/{name}.md` |

### 5. Pull-Based Prelude（Class-2 平台）

对于没有 PreToolUse hook 的平台（Codex, Kiro, Gemini, Qoder, Copilot），子代理需要内置上下文加载逻辑：

```markdown
## Before acting

1. Read `{task_dir}/implement.jsonl` to discover required files
2. Read `{task_dir}/prd.md`
3. Read `{task_dir}/design.md` if it exists
4. Read `{task_dir}/implement.md` if it exists
```

### 6. Context Injection Hook（Class-1 平台）

`inject-subagent-context.py` 在 PreToolUse 事件触发时：
1. 检测子代理类型（implement/check/research）
2. 加载对应的 JSONL manifest
3. 注入 `prd.md`、`design.md`、`implement.md` 内容
4. 注入相关 spec 文件

## 验收标准

- [ ] 3 个子代理定义在所有 11 个 agent-capable 平台上可用
- [ ] Class-1 平台通过 hook 自动注入上下文
- [ ] Class-2 平台通过 pull-based prelude 加载上下文
- [ ] `trellis-implement` 正确读取 manifest 和 spec
- [ ] `trellis-check` 自修复循环最多 3 轮
- [ ] `trellis-research` 只读，不修改文件

## 依赖

- [platform-configurators](./platform-configurators.md) — 模板生成
- [context-system](./context-system.md) — JSONL manifest 和上下文注入
- [auto-skills](./auto-skills.md) — skill 和 sub-agent 的协作关系

## 关联 Spec

- [workflow-engine](./workflow-engine.md) — 工作流状态驱动 sub-agent 调度

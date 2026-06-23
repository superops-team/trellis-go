# slash-commands: 用户命令

> 创建 3 个 Trellis 用户命令，支持多平台格式。

## 当前状态

| 组件 | 状态 |
|------|:---:|
| 命令定义 | ❌ |
| 命令模板 | ❌ |
| 多平台格式适配 | ❌ |

## 需求

### 1. `/trellis:start`（会话启动）

**适用平台**: 无 SessionStart hook/extension 的平台（Kilo, Antigravity, Devin, Kiro）

**流程**:
1. 读取 `.trellis/workflow.md` → AI 了解工作流契约
2. 运行 `get_context.py` → 开发者身份、git 状态、活跃任务
3. 读取 spec index（monorepo 中按 package）
4. 报告上下文并询问要做什么

**任务分类**:

| 类型 | 条件 | 流程 |
|------|------|------|
| 简单对话 | 问题/解释/查找，无代码变更 | 默认不创建任务 |
| 内联小任务 | 单轮可理解和验证的编辑 | 询问是否创建任务 |
| 完整 Trellis 任务 | 多文件或持久规划工作 | 询问是否进入 planning |

### 2. `/trellis:finish-work`（归档 + 记录）

**前置条件**: 代码已提交

**流程**:
1. `get_context.py --mode record` → 活跃任务、git 状态、最近提交
2. `git status --porcelain` → 排除 `.trellis/workspace/` 和 `.trellis/tasks/`
3. 如有未提交变更 → 拒绝执行，引导用户回 Phase 3.4
4. `task.py archive <name>` → 归档任务
5. `add_session.py --title ... --commit ...` → 记录 session

**Git log 顺序**: `<work commits>` → `chore(task): archive ...` → `chore: record journal`

### 3. `/trellis:continue`（推进工作流）

**功能**: 任务内推进，非跨任务

**流程**:
1. 读取 `task.json.status` + workflow-state breadcrumb
2. 查阅 `workflow.md` 定位当前 phase/step
3. 推进到下一步

**典型对话流**:
```
User: "继续" → AI 判断是轻量任务还是需要 design.md/implement.md
User: "继续" → AI 开始 implement/check
User: "继续" → AI 路由到 trellis-update-spec，然后 finish-work
```

### 4. 多平台格式

#### 显式 Slash Commands

| 平台 | 路径 | 调用方式 |
|------|------|---------|
| Claude Code | `.claude/commands/trellis/{name}.md` | `/trellis:{name}` |
| Cursor | `.cursor/commands/trellis-{name}.md` | `/trellis-{name}` |
| OpenCode | `.opencode/commands/trellis/{name}.md` | `/trellis:{name}` |
| Codex | `.codex/prompts/trellis-{name}.md` | `/trellis-{name}` |
| Gemini CLI | `.gemini/commands/trellis/{name}.toml` | `/trellis:{name}` |
| Qoder | `.qoder/commands/trellis-{name}.md` | `/trellis-{name}` |
| CodeBuddy | `.codebuddy/commands/trellis/{name}.md` | `/trellis:{name}` |
| Droid | `.factory/commands/trellis/{name}.md` | `/trellis:{name}` |
| Pi Agent | `.pi/prompts/trellis-{name}.md` | `/trellis-{name}` |
| Copilot | `.github/prompts/trellis-{name}.prompt.md` | prompt-file picker |

#### Workflow Files（无 slash command 原语的平台）

| 平台 | 路径 | 调用方式 |
|------|------|---------|
| Kilo | `.kilocode/workflows/{name}.md` | workflow UI |
| Antigravity | `.agent/workflows/{name}.md` | workflow UI |
| Devin | `.devin/workflows/trellis-{name}.md` | `/trellis-{name}` |

#### Skill-Only（无命令原语的平台）

| 平台 | 路径 |
|------|------|
| Kiro | `.kiro/skills/trellis-{name}/SKILL.md` |

### 5. 命令 vs 技能选择规则

| 用命令 | 用技能 |
|--------|--------|
| 用户决定何时运行 | AI 基于意图自动触发 |
| 标记会话边界（start/finish/continue） | 任务内阶段（before-dev/check/update-spec） |
| 无可靠匹配的触发短语 | 有可预测的用户意图 |
| 无活跃任务时也需要 | 只在活跃任务上下文中有意义 |

## 验收标准

- [ ] 3 个命令在所有支持平台上可用
- [ ] `finish-work` 在 dirty working tree 时拒绝执行
- [ ] `continue` 正确推进工作流状态
- [ ] `start` 正确加载上下文并分类任务
- [ ] Qoder 命令包含正确的 YAML frontmatter
- [ ] Gemini 命令使用正确的 TOML 格式

## 依赖

- [session-journal](./session-journal.md) — finish-work 的 session 记录
- [workflow-engine](./workflow-engine.md) — continue 的工作流推进
- [context-system](./context-system.md) — start 的上下文加载
- [platform-configurators](./platform-configurators.md) — 模板生成

## 关联 Spec

- [auto-skills](./auto-skills.md) — 命令可手动触发技能
- [task-lifecycle](./task-lifecycle.md) — finish-work 触发 task archive

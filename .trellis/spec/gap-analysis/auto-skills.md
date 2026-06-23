# auto-skills: 自动触发技能

> 创建 5 个 Trellis 自动触发技能，支持多平台格式。

## 当前状态

| 组件 | 状态 |
|------|:---:|
| 技能定义 | ❌ |
| 技能模板 | ❌ |
| 多平台格式适配 | ❌ |
| Skills marketplace | ❌ |

## 需求

### 1. trellis-brainstorm（需求澄清 + PRD 起草）

**触发条件**: 用户想要新功能 / 需求不明确

**流程**:
1. 检查代码、测试、配置、文档、已有 spec、任务历史
2. 提出任务名称和 slug
3. 必要时通过 `task.py create` 创建任务
4. 起草和迭代 `prd.md`（需求 + 验收标准）
5. 每次只问一个问题，包含推荐答案
6. 复杂任务：添加 `design.md` 和 `implement.md`

**SKILL.md 结构**:
```markdown
---
name: trellis-brainstorm
description: |
  Use when the user wants a new feature or the requirements are unclear.
  Clarifies requirements, inspects evidence, and drafts planning artifacts.
---

# trellis-brainstorm

## Trigger check
Verify the user is requesting new work that needs planning.

## Steps
1. Inspect codebase, existing specs, and task history
2. Propose task name and create via `task.py create`
3. Draft `prd.md` with requirements and acceptance criteria
4. Ask one question at a time with recommended answer
5. For complex tasks, add `design.md` and `implement.md`
```

### 2. trellis-before-dev（开发前准备）

**触发条件**: 即将开始编写代码

**流程**:
1. 读取受影响 package 的 spec index
2. 读取开发前检查清单中的具体指南文件
3. 确保 AI 在写代码前了解约定

**SKILL.md 结构**:
```markdown
---
name: trellis-before-dev
description: |
  Use before touching code in a task. Reads relevant specs so the AI knows
  the conventions before writing, not after.
---

# trellis-before-dev

## Steps
1. Identify affected packages from the task context
2. Read spec index for each package
3. Read pre-development checklist guidelines
4. Confirm conventions are understood before proceeding
```

### 3. trellis-check（实现后验证）

**触发条件**: 完成代码实现后

**流程**:
1. `git diff --name-only HEAD` 找到变更
2. 发现适用的 spec 层
3. 逐层对比 diff 和质量检查清单
4. 运行 lint/typecheck/test
5. 发现问题 → 修复 → 重新验证（最多 3 轮）
6. 报告结果

### 4. trellis-update-spec（知识沉淀）

**触发条件**: 有值得捕获的学习/经验

**流程**:
1. 识别值得沉淀的知识
2. 确定目标 spec 层
3. 用新知识更新 spec 文件
4. 保持 spec 简洁可操作

### 5. trellis-break-loop（Bug 根因分析）

**触发条件**: 遇到棘手的 bug / 反复修复同一问题

**流程**:
1. 分析 bug 根因
2. 识别为什么现有流程没有捕获它
3. 提出预防措施
4. 必要时更新 spec

### 6. 多平台格式

所有 14 个平台都需要这 5 个技能：

| 平台 | 路径 |
|------|------|
| Claude Code | `.claude/skills/{name}/SKILL.md` |
| Cursor | `.cursor/skills/trellis-{name}/SKILL.md` |
| OpenCode | `.opencode/skills/{name}/SKILL.md` |
| Codex | `.codex/skills/{name}/SKILL.md` |
| Kiro | `.kiro/skills/{name}/SKILL.md` |
| Gemini CLI | `.gemini/skills/trellis-{name}/SKILL.md` |
| Qoder | `.qoder/skills/trellis-{name}/SKILL.md` |
| CodeBuddy | `.codebuddy/skills/trellis-{name}/SKILL.md` |
| Copilot | `.github/skills/{name}/SKILL.md` |
| Droid | `.factory/skills/{name}/SKILL.md` |
| Pi Agent | `.pi/skills/{name}/SKILL.md` |
| Kilo | `.kilocode/skills/{name}/SKILL.md` |
| Antigravity | `.agent/skills/{name}/SKILL.md` |
| Devin | `.devin/skills/{name}/SKILL.md` |

**跨平台共享层**: `.agents/skills/{name}/SKILL.md`（所有平台写入）

### 7. Skills Marketplace（P3）

后续可添加的社区技能：
- `frontend-fullchain-optimization` — Web Vitals 驱动的前端优化
- `mem-recall` — 跨平台 AI 对话回忆
- `trellis-meta` — Trellis 自定义技能
- `trellis-spec-bootstrap` — 从代码库引导生成 spec

## 验收标准

- [ ] 5 个技能在所有 14 个平台上可用
- [ ] `description` 字段正确触发自动匹配
- [ ] 每个技能有明确的触发条件检查
- [ ] 每个技能有固定的输出格式
- [ ] `.agents/skills/` 跨平台共享层正确写入
- [ ] Skill routing table 正确映射用户意图到技能

## 依赖

- [workflow-engine](./workflow-engine.md) — skill routing table
- [platform-configurators](./platform-configurators.md) — 模板生成

## 关联 Spec

- [sub-agents](./sub-agents.md) — skill 和 sub-agent 协作
- [slash-commands](./slash-commands.md) — command 可手动触发 skill

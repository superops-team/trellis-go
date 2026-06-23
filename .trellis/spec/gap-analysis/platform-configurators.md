# platform-configurators: 平台配置器

> 为 16 个 AI 编码平台实现完整的配置器，生成 hooks、agents、skills、commands。

## 当前状态

| 组件 | 状态 | 说明 |
|------|:---:|------|
| `platform.Registry` | ✅ | 15 平台注册 |
| `hook.Generator` | ⚠️ | 只生成基础 hook 脚本，无平台差异 |
| 平台模板目录 | ❌ | `internal/embed` 存在但模板内容为空 |
| 配置器（per-platform） | ❌ | 完全缺失 |

## 平台分类

### Class-1: Push-Based（完整自动化）
Hooks + Sub-agents + Skills + Commands 全部支持。

| 平台 | ConfigDir | 关键差异 |
|------|-----------|---------|
| Claude Code | `.claude/` | Python hooks, agents/*.md, commands/trellis/*.md |
| Cursor | `.cursor/` | hooks.json, 扁平命令命名 trellis-*.md |
| OpenCode | `.opencode/` | JS plugins 替代 Python hooks |
| Kiro | `.kiro/` | JSON agents, `.kiro.hook` IDE hook |
| CodeBuddy | `.codebuddy/` | CC 兼容格式 |
| Droid | `.factory/` | CC 兼容格式 |
| Pi Agent | `.pi/` | TypeScript extension, prompts/*.md |
| Reasonix | `.reasonix/` | CC 兼容格式 |

### Class-2: Pull-Based（部分自动化）
Sub-agents 使用 pull-based prelude，无 sub-agent PreToolUse hook。

| 平台 | ConfigDir | 关键差异 |
|------|-----------|---------|
| Codex | `.codex/` | TOML agents, AGENTS.md 入口, prompts/*.md |
| Gemini CLI | `.gemini/` | TOML commands, settings.json |
| Qoder | `.qoder/` | YAML frontmatter commands, settings.json |
| Copilot | `.github/` | `.agent.md` 格式 agents, prompts/*.prompt.md |

### Class-3: Agentless（无代理）
无 sub-agents，无 hooks，仅 workflows + skills。

| 平台 | ConfigDir | 关键差异 |
|------|-----------|---------|
| Kilo | `.kilocode/` | workflows/*.md |
| Antigravity | `.agent/` | workflows/*.md |
| Devin | `.devin/` | workflows/*.md |

### 缺失平台

| 平台 | Python 版 | Go 版 |
|------|:---:|:---:|
| ZCode | ✅ | ❌ |
| Devin (原 Windsurf) | ✅ | ⚠️ 注册为 "windsurf" |

## 需求

### 1. 模板内容

每个平台需要以下模板文件：

```
templates/{platform}/
├── hooks/
│   ├── session-start.py          # SessionStart hook
│   ├── inject-workflow-state.py  # UserPromptSubmit hook
│   └── inject-subagent-context.py # PreToolUse hook (Class-1 only)
├── agents/
│   ├── trellis-implement.md      # 实现子代理
│   ├── trellis-check.md          # 验证子代理
│   └── trellis-research.md       # 调研子代理
├── skills/
│   ├── trellis-brainstorm/SKILL.md
│   ├── trellis-before-dev/SKILL.md
│   ├── trellis-check/SKILL.md
│   ├── trellis-update-spec/SKILL.md
│   └── trellis-break-loop/SKILL.md
├── commands/
│   ├── finish-work.md
│   └── continue.md
└── settings.json                 # 平台 hook 配置
```

### 2. 平台特定适配

**Claude Code** (参考实现):
- `.claude/hooks/` — Python 脚本
- `.claude/settings.json` — hook 注册
- `.claude/agents/*.md` — YAML frontmatter + Markdown body
- `.claude/commands/trellis/*.md` — 子目录结构
- `.claude/skills/*/SKILL.md` — 标准格式

**Cursor**:
- `.cursor/hooks.json` — 独立 hook 配置文件
- `.cursor/commands/trellis-{name}.md` — 扁平命名
- 其他同 Claude Code

**OpenCode**:
- `.opencode/plugins/*.js` — JS 插件替代 Python hooks
- `.opencode/agents/*.md` — `permission:` 对象格式

**Codex**:
- `.codex/agents/*.toml` — TOML 格式
- `.codex/prompts/trellis-{name}.md` — prompt 文件
- `AGENTS.md` — 仓库根目录入口
- `.codex/hooks.json` — hook 配置

**Kiro**:
- `.kiro/agents/trellis.json` — 主代理 + per-turn hook
- `.kiro/agents/trellis-{implement,check,research}.json` — 子代理
- `.kiro/hooks/trellis-workflow-state.kiro.hook` — IDE hook

**Pi Agent**:
- `.pi/extensions/trellis/index.ts` — TypeScript 扩展
- `.pi/prompts/trellis-{name}.md` — prompt 文件

### 3. Configurator 接口

```go
type Configurator interface {
    // Name returns the platform identifier.
    Name() string
    
    // Generate writes all platform-specific files to the project.
    Generate(projectRoot string, opts GenerateOptions) error
    
    // Update syncs templates, preserving user modifications.
    Update(projectRoot string, opts UpdateOptions) error
    
    // Remove cleans up platform files (for uninstall).
    Remove(projectRoot string) error
}
```

### 4. 跨平台共享层

- `.agents/skills/` — 所有平台写入此目录（agentskills.io 标准）
- 由 Codex 配置器负责写入

## 验收标准

- [ ] `trellis init --claude` 生成完整 Claude Code 配置
- [ ] `trellis init --cursor` 生成完整 Cursor 配置
- [ ] `trellis init --codex` 生成完整 Codex 配置（含 AGENTS.md）
- [ ] `trellis init --all` 生成所有平台配置
- [ ] `trellis update` 正确同步模板变更
- [ ] 每个平台的 hook 脚本可独立执行
- [ ] 每个平台的 agent 定义格式正确
- [ ] `.agents/skills/` 在所有平台上正确写入

## 依赖

- [sub-agents](./sub-agents.md) — agent 定义内容
- [auto-skills](./auto-skills.md) — skill 定义内容
- [slash-commands](./slash-commands.md) — command 定义内容
- [workflow-engine](./workflow-engine.md) — hook 脚本中的 workflow 逻辑

## 工作量

~2 周（16 平台 × 平均 0.8 天/平台）

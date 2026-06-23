# context-system: 上下文系统

> 实现 `get_context.py` 等价功能：多模式上下文构建、JSONL manifest、上下文注入。

## 当前状态

| 组件 | 状态 | 说明 |
|------|:---:|------|
| `context.Builder` | ⚠️ | 只有 BuildImplementContext / BuildCheckContext / BuildResearchContext |
| `context.Manifest` | ✅ | JSONL 读写 |
| Developer identity | ❌ | `.trellis/.developer` 未读取 |
| Git status 注入 | ❌ | 未集成到上下文 |
| Active task 检测 | ❌ | 未实现 |
| Phase/step 上下文 | ❌ | 未实现 |
| Session recording 模式 | ❌ | 未实现 |
| Monorepo package 感知 | ❌ | 未实现 |

## 需求

### 1. `get_context.py` 等价实现

**模式**:

| 模式 | 用途 | 调用方 |
|------|------|--------|
| `--mode session` | SessionStart 上下文 | session-start hook |
| `--mode implement` | 实现阶段上下文 | inject-subagent-context hook |
| `--mode check` | 验证阶段上下文 | inject-subagent-context hook |
| `--mode research` | 调研阶段上下文 | inject-subagent-context hook |
| `--mode record` | Session 记录上下文 | finish-work command |
| `--mode phase --step X.Y` | Phase/step 提取 | 通用 |

### 2. Session 模式输出

```
Developer: <name>
Repository: <git remote>
Branch: <current-branch>
Status: clean | N files modified
Active task: <task-id> (<status>)

## Workflow
<workflow.md Phase Index content>

## Spec Index
<per-package spec index>

## Recent Tasks
<list of recent tasks>
```

### 3. Implement/Check/Research 模式

**Implement 模式**:
1. 读取 `{task_dir}/implement.jsonl`
2. 加载每个 JSONL entry 指向的文件
3. 注入 `prd.md`
4. 注入 `design.md`（如果存在）
5. 注入 `implement.md`（如果存在）
6. 注入相关 spec 文件

**Check 模式**:
1. 读取 `{task_dir}/check.jsonl`
2. 加载相关 spec 质量检查清单
3. 注入 `prd.md`

**Research 模式**:
1. 读取 `{task_dir}/research.jsonl`
2. 加载相关 spec 文件

### 4. Record 模式

**输出**:
```
Active tasks:
  - <task-id> (<status>): <name>

Git status:
  Branch: <branch>
  Recent commits:
    abc1234 - <message>
    def5678 - <message>

Unarchived completed tasks:
  - <task-id>: <name>
```

### 5. Phase 模式

**用法**: `get_context.py --mode phase --step 2.1`

**输出**: `workflow.md` 中 `#### 2.1` 步骤的正文内容

### 6. JSONL Manifest 格式

```jsonl
{"path": "src/auth/login.ts", "description": "Login page component", "required": true}
{"path": "src/auth/api.ts", "description": "Auth API client", "required": true}
{"path": "docs/auth-flow.md", "description": "Auth flow documentation", "required": false}
```

### 7. 上下文注入格式

```
<!-- trellis-hook-injected -->

=== file: prd.md ===
<content>

=== file: design.md ===
<content>

=== file: src/auth/login.ts ===
<content>
```

### 8. Monorepo 感知

- 读取 `config.yaml` 中的 `packages` 配置
- 根据 `default_package` 或 task 指定的 package 确定上下文范围
- Spec index 按 package 分组

## API 设计

```go
// ContextBuilder builds context for different modes.
type ContextBuilder struct {
    SpecLoader   *spec.Loader
    TaskManager  *task.Manager
    GitClient    *git.Client
    Config       *config.Config
    Root         string
}

func (b *ContextBuilder) BuildSessionContext() (*SessionContext, error)
func (b *ContextBuilder) BuildImplementContext(taskID string) (string, error)
func (b *ContextBuilder) BuildCheckContext(taskID string) (string, error)
func (b *ContextBuilder) BuildResearchContext(taskID string) (string, error)
func (b *ContextBuilder) BuildRecordContext() (*RecordContext, error)
func (b *ContextBuilder) BuildPhaseContext(step string) (string, error)
```

## 验收标准

- [ ] 6 种模式全部实现
- [ ] Session 模式输出包含 developer identity + git status + active tasks
- [ ] Implement/Check/Research 模式正确加载 JSONL manifest
- [ ] Record 模式正确列出活跃任务和最近提交
- [ ] Phase 模式正确提取步骤内容
- [ ] Monorepo package 感知正常工作
- [ ] 上下文注入格式与 Python 版兼容

## 依赖

- [session-journal](./session-journal.md) — developer identity
- [workflow-engine](./workflow-engine.md) — phase/step 解析
- [task-lifecycle](./task-lifecycle.md) — 任务状态查询

## 关联 Spec

- [sub-agents](./sub-agents.md) — 子代理上下文注入
- [slash-commands](./slash-commands.md) — start/finish-work 的上下文加载

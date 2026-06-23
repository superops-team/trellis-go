# task-lifecycle: 任务生命周期

> 实现完整的任务生命周期管理：hooks、monorepo 支持、spec 关联。

## 当前状态

| 组件 | 状态 | 说明 |
|------|:---:|------|
| Task CRUD | ✅ | create/start/archive/list |
| Task hooks | ❌ | `config.Hooks` 字段存在但未执行 |
| Monorepo 支持 | ❌ | `config.Packages` 字段存在但未使用 |
| Spec 关联 | ❌ | 无 task-spec 关联机制 |
| 子任务管理 | ⚠️ | 数据结构存在，CLI 缺失 |

## 需求

### 1. Task Lifecycle Hooks

**配置**:
```yaml
hooks:
  after_create:
    - "echo 'Task created'"
  after_start:
    - "echo 'Task started'"
  after_finish:
    - "echo 'Task finished'"
  after_archive:
    - "echo 'Task archived'"
```

**执行规则**:
- 每个 hook 接收 `TASK_JSON_PATH` 环境变量
- Hook 失败打印警告但不阻塞主操作
- Shell 命令通过 `exec.Command("sh", "-c", cmd)` 执行

**实现**:
```go
type HookRunner struct {
    Hooks map[string][]string
}

func (r *HookRunner) Run(event string, taskJSONPath string) error {
    cmds, ok := r.Hooks[event]
    if !ok {
        return nil
    }
    for _, cmd := range cmds {
        c := exec.Command("sh", "-c", cmd)
        c.Env = append(os.Environ(), "TASK_JSON_PATH="+taskJSONPath)
        if err := c.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "Warning: hook %q failed: %v\n", cmd, err)
        }
    }
    return nil
}
```

### 2. Monorepo Package 支持

**配置**:
```yaml
packages:
  frontend:
    path: packages/frontend
  backend:
    path: packages/backend
  docs:
    path: docs-site
    type: submodule
  webapp:
    path: ./webapp
    git: true

default_package: frontend
```

**功能**:
- Task 关联 package（`task.json` 增加 `package` 字段）
- 命令和上下文按 package 过滤
- `default_package` 回退逻辑
- `type: submodule` 和 `git: true` 的路径处理

### 3. Task-Spec 关联

**功能**:
- `task add-spec <id> <spec-path>` — 关联 spec 文件
- `task list-specs <id>` — 列出关联 spec
- 上下文构建时自动加载关联 spec

**task.json 扩展**:
```json
{
  "id": "06-22-add-login",
  "name": "add-login",
  "status": "in_progress",
  "package": "frontend",
  "specs": [
    ".trellis/spec/frontend/security/index.md",
    ".trellis/spec/frontend/api/index.md"
  ],
  "subtasks": [
    {"id": "1", "title": "Login form UI", "done": true},
    {"id": "2", "title": "API integration", "done": false}
  ]
}
```

### 4. 子任务 CLI

```bash
trellis task add-subtask <task-id> "Login form UI"
trellis task done-subtask <task-id> <subtask-id>
trellis task undone-subtask <task-id> <subtask-id>
```

### 5. 任务状态流转

```
planning ──[start]──▶ in_progress ──[archive]──▶ completed
                         │
                         └── 自定义状态（blocked, in_review 等）
```

**自定义状态**:
- 通过直接编辑 `task.json` 设置
- `task list --status <any-string>` 可过滤任意状态
- Breadcrumb 系统自动匹配 `[workflow-state:<status>]`

### 6. 任务目录结构

```
.trellis/tasks/06-22-add-login/
├── task.json
├── prd.md
├── design.md          # 可选
├── implement.md       # 可选
├── implement.jsonl    # 实现上下文 manifest
├── check.jsonl        # 验证上下文 manifest
├── research.jsonl     # 调研上下文 manifest
└── info.md            # 任务元信息
```

## 验收标准

- [ ] 4 个 lifecycle hooks 正确执行
- [ ] Hook 失败不阻塞任务操作
- [ ] `TASK_JSON_PATH` 环境变量正确传递
- [ ] Monorepo package 关联正常工作
- [ ] `default_package` 回退正确
- [ ] Task-spec 关联可添加和列出
- [ ] 子任务 CLI 完整可用
- [ ] 自定义状态可设置和过滤

## 依赖

- [cli-core](./cli-core.md) — task 子命令
- [workflow-engine](./workflow-engine.md) — 状态变更触发 breadcrumb 切换
- [context-system](./context-system.md) — spec 关联影响上下文构建

## 关联 Spec

- [session-journal](./session-journal.md) — task archive 触发 session 记录
- [slash-commands](./slash-commands.md) — finish-work 触发 task archive

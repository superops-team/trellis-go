# session-journal: 会话记录系统

> 实现 `add_session.py` 等价功能：session journal 记录、轮转、自动提交。

## 当前状态

| 组件 | 状态 |
|------|:---:|
| Session journal 记录 | ❌ |
| Journal 轮转 | ❌ |
| Session 自动提交 | ❌ |
| Session index | ❌ |
| Developer identity | ❌ |

## 需求

### 1. Developer Identity

**文件**: `.trellis/.developer`（gitignored，每机器独立）

```yaml
name: "developer-name"
```

**CLI 集成**:
- `trellis init` 在已存在项目中提供 "Set up developer identity" 选项
- 默认值：`git config user.name` 或 `$USER`

### 2. Journal 文件结构

```
.trellis/workspace/<developer>/
├── journal-1.md
├── journal-2.md
└── index.json
```

**journal-N.md 格式**:
```markdown
# Session 2026-06-22 14:30

- Task: add-login-feature
- Commits: abc1234, def5678
- Summary: Implemented login page and API integration

---

# Session 2026-06-22 10:00

- Task: fix-auth-bug
- Commits: xyz9012
- Summary: Fixed token refresh race condition
```

**index.json 格式**:
```json
{
  "sessions": [
    {
      "id": "2026-06-22T14:30:00+08:00",
      "title": "add-login-feature",
      "commits": ["abc1234", "def5678"],
      "journal_file": "journal-1.md",
      "started_at": "2026-06-22T14:30:00+08:00",
      "finished_at": "2026-06-22T16:45:00+08:00"
    }
  ]
}
```

### 3. Journal 轮转

**规则**:
- `max_journal_lines` 配置项（默认 2000）
- 当前 journal 文件行数超过限制时，创建新文件（journal-2.md, journal-3.md...）
- 轮转时在旧文件末尾添加 `(continued in journal-N.md)`

### 4. Session 自动提交

**配置**: `session_auto_commit`（默认 true）

**行为**:
- `true`: 写入 journal 后自动 `git add` + `git commit`
- `false`: 只写入文件，不提交
- 提交信息：`session_commit_message` 配置（默认 "chore: record journal"）

**排除路径**: `.trellis/workspace/` 和 `.trellis/tasks/` 由脚本管理，不触发 dirty check

### 5. API 设计

```go
// SessionRecorder manages session journal recording.
type SessionRecorder struct {
    WorkspaceDir string
    Developer    string
    Config       SessionConfig
}

type SessionConfig struct {
    MaxJournalLines     int
    AutoCommit          bool
    CommitMessage       string
}

// RecordSession appends a session entry to the current journal.
func (r *SessionRecorder) RecordSession(entry SessionEntry) error

// SessionEntry represents one recorded session.
type SessionEntry struct {
    Title    string
    Commits  []string
    Summary  string
    TaskID   string
}
```

### 6. CLI 集成

```bash
# 记录 session（由 /finish-work 命令调用）
trellis hook record-session --title "..." --commit abc1234 --commit def5678

# 查看 session 历史
trellis hook list-sessions
```

## 验收标准

- [ ] `.trellis/.developer` 正确写入和读取
- [ ] Journal 文件按 `max_journal_lines` 自动轮转
- [ ] `session_auto_commit: true` 时自动提交
- [ ] `session_auto_commit: false` 时只写文件
- [ ] `trellis init` 多人协作模式正确设置 developer identity
- [ ] Session index 正确追踪所有 session

## 依赖

- [cli-core](./cli-core.md) — `init` 多人协作模式
- [context-system](./context-system.md) — `get_context.py --mode record`

## 关联 Spec

- [task-lifecycle](./task-lifecycle.md) — `task archive` 触发 session 记录

# test-e2e-subtask: Subtask 端到端测试

> trellis-go 新增功能，原版 trellis 无此功能，需自行设计测试。

## 当前状态

trellis-go 已实现 subtask CLI 命令（`add-subtask <task-id> <title>`、`done-subtask <task-id> <subtask-id>`、`undone-subtask <task-id> <subtask-id>`），`pkg/task/manager.go` 有 `AddSubtask`/`DoneSubtask`/`UndoneSubtask` 方法。

`pkg/task/manager_test.go` 有 13 个测试，但无 Subtask 相关测试。无 Subtask E2E。

## 需求

### 1. Subtask CRUD E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_SubtaskAddAndList(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. trellis task add-subtask <id> "Setup DB schema"
    // 4. trellis task add-subtask <id> "Create API endpoint"
    // 5. trellis task info <id>
    // 6. 验证输出包含两个子任务
}

func TestE2E_SubtaskDone(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. trellis task add-subtask <id> "Setup DB"
    // 4. trellis task done-subtask <id> <subtask-id>
    // 5. trellis task info <id>
    // 6. 验证子任务标记为完成
}

func TestE2E_SubtaskProgress(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. 添加 3 个子任务
    // 4. 完成 1 个
    // 5. trellis task info <id>
    // 6. 验证进度显示 1/3
}

func TestE2E_SubtaskAllDone(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. 添加 2 个子任务
    // 4. 全部完成
    // 5. trellis task info <id>
    // 6. 验证任务状态或提示（全部子任务完成）
}

func TestE2E_SubtaskDuplicateName(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. trellis task add-subtask <id> "Setup DB"
    // 4. trellis task add-subtask <id> "Setup DB"  ← 重复
    // 5. 验证报错或警告
}

func TestE2E_SubtaskDoneNotFound(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "Main feature"
    // 3. trellis task done-subtask <id> "nonexistent-id"
    // 4. 验证报错
}
```

### 2. Subtask 单元测试补充

在 `pkg/task/manager_test.go` 中补充：

```go
func TestManager_AddSubtask(t *testing.T) { ... }
func TestManager_DoneSubtask(t *testing.T) { ... }
func TestManager_SubtaskProgress(t *testing.T) { ... }
func TestManager_SubtaskAllDoneTransitions(t *testing.T) { ... }
```

## 验收标准

- [ ] 6 个 Subtask E2E 场景全部通过
- [ ] 4 个 Subtask 单元测试通过
- [ ] 覆盖正常流程 + 边界（重复、不存在、全部完成）

## 参考

- trellis-go 现有：`pkg/task/manager_test.go`、`cmd/trellis/e2e_test.go`

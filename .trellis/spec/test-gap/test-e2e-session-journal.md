# test-e2e-session-journal: Session Journal 端到端测试

> trellis-go 新增功能，原版 trellis 无此功能，需自行设计测试。

## 当前状态

trellis-go 有 `pkg/session/recorder_test.go` 单元测试，但无 E2E。

## 需求

### 1. Session 记录 E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_SessionRecordAndList(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. trellis hook record-session --title "Fix login bug" --commit abc1234
    // 3. trellis hook record-session --title "Add payment API" --commit def5678
    // 4. trellis hook list-sessions
    // 5. 验证输出包含两条记录
}

func TestE2E_SessionListSearch(t *testing.T) {
    // 1. trellis init
    // 2. 记录多条 session
    // 3. trellis hook list-sessions --search "payment"
    // 4. 验证只返回匹配的记录
}

func TestE2E_SessionListEmpty(t *testing.T) {
    // 1. trellis init
    // 2. trellis hook list-sessions
    // 3. 验证输出 "No sessions recorded" 或空列表
}

func TestE2E_SessionJournalFileFormat(t *testing.T) {
    // 1. trellis init
    // 2. trellis hook record-session --title "Test" --commit abc1234
    // 3. 读取 .trellis/sessions/ 下的 journal 文件
    // 4. 验证 JSON 格式正确，包含 title、commit、timestamp
}
```

### 2. Session 记录器单元测试补充

在 `pkg/session/recorder_test.go` 中补充：

```go
func TestRecorder_DuplicateCommit(t *testing.T) {
    // 同一 commit 重复记录的行为
}

func TestRecorder_EmptyTitle(t *testing.T) {
    // 空标题的处理
}

func TestRecorder_JournalFilePermissions(t *testing.T) {
    // 生成文件的权限正确
}
```

## 验收标准

- [ ] `TestE2E_SessionRecordAndList` 通过
- [ ] `TestE2E_SessionListSearch` 通过
- [ ] `TestE2E_SessionListEmpty` 通过
- [ ] `TestE2E_SessionJournalFileFormat` 通过
- [ ] 单元测试补充 3 个边界用例

## 参考

- trellis-go 现有：`pkg/session/recorder_test.go`

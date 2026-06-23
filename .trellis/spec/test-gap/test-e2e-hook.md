# test-e2e-hook: Hook 端到端测试

> 参考原版 trellis 模板测试（`templates/*.test.ts`），补齐 trellis-go hook 输出验证。

## 当前状态

trellis-go 有 `pkg/hook/generator_test.go` 单元测试，但无 E2E 验证 hook 实际输出内容。

## 需求

### 1. Hook 输出内容 E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_HookInjectWorkflowState(t *testing.T) {
    // 1. trellis init
    // 2. 写入 workflow.md
    // 3. trellis task create "test"
    // 4. trellis task start <id> --phase research
    // 5. trellis hook inject-workflow-state
    // 6. 验证输出包含：
    //    - 当前阶段名称
    //    - 任务 ID
    //    - 阶段描述
}

func TestE2E_HookSessionStart(t *testing.T) {
    // 1. trellis init
    // 2. trellis task create "test"
    // 3. trellis task start <id> --phase research
    // 4. trellis hook session-start
    // 5. 验证输出包含：
    //    - 当前任务信息
    //    - 当前阶段
    //    - 上下文提示
}

func TestE2E_HookInjectContext(t *testing.T) {
    // 1. trellis init
    // 2. trellis context add spec.md
    // 3. trellis hook inject-context
    // 4. 验证输出包含 spec.md 内容
}

func TestE2E_HookRecordSession(t *testing.T) {
    // 1. trellis init
    // 2. trellis hook record-session --title "Test" --commit abc1234
    // 3. 验证输出确认记录成功
}

func TestE2E_HookListSessions(t *testing.T) {
    // 1. trellis init
    // 2. 记录 2 条 session
    // 3. trellis hook list-sessions
    // 4. 验证输出格式正确
}
```

### 2. 多平台 Hook 文件生成验证

```go
func TestE2E_HookFilesForAllPlatforms(t *testing.T) {
    // 1. trellis init --platform claude --platform cursor --platform codex --platform gemini
    // 2. 验证每个平台的 hook 目录存在
    // 3. 验证每个平台有 session-start hook
    // 4. 验证每个平台有 inject-context hook
    // 5. 验证 hook 文件有执行权限
}
```

### 3. Hook 生成器单元测试补充

在 `pkg/hook/generator_test.go` 中补充：

```go
func TestGenerator_AllPlatformsHaveSessionStart(t *testing.T) { ... }
func TestGenerator_AllPlatformsHaveInjectContext(t *testing.T) { ... }
func TestGenerator_HookScriptPermissions(t *testing.T) { ... }
func TestGenerator_HookScriptShebang(t *testing.T) { ... }
```

## 验收标准

- [ ] 5 个 Hook 输出内容 E2E 通过
- [ ] `TestE2E_HookFilesForAllPlatforms` 通过
- [ ] 4 个 Hook 生成器单元测试通过

## 参考

- 原版模板测试：`/tmp/trellis-orig/packages/cli/test/templates/*.test.ts`
- trellis-go 现有：`pkg/hook/generator_test.go`

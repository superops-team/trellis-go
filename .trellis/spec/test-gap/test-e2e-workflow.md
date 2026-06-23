# test-e2e-workflow: Workflow 端到端测试

> 参考原版 trellis `workflow.integration.test.ts`（302 行）。

## 当前状态

trellis-go 有 `pkg/workflow/parser_test.go` 单元测试，但无 E2E。

## 需求

### 1. Workflow 状态转换 E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_WorkflowFullCycle(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 写入 workflow.md（4 阶段：research → plan → implement → verify）
    // 3. trellis task create "test"
    // 4. trellis task start <id> --phase research
    // 5. 验证 inject-workflow-state 输出当前阶段
    // 6. trellis task start <id> --phase plan
    // 7. 验证阶段切换
    // 8. trellis task start <id> --phase implement
    // 9. trellis task start <id> --phase verify
    // 10. trellis task archive <id>
}

func TestE2E_WorkflowInvalidPhaseTransition(t *testing.T) {
    // 1. trellis init
    // 2. 写入 workflow.md
    // 3. trellis task create "test"
    // 4. trellis task start <id> --phase plan  ← 跳过 research
    // 5. 验证报错或警告
}

func TestE2E_WorkflowStateInjection(t *testing.T) {
    // 1. trellis init
    // 2. 写入 workflow.md
    // 3. trellis task create "test"
    // 4. trellis task start <id> --phase research
    // 5. trellis hook inject-workflow-state
    // 6. 验证输出包含当前阶段、任务 ID、状态
}
```

### 2. Workflow 解析单元测试补充

在 `pkg/workflow/parser_test.go` 中补充：

```go
func TestParser_EmptyWorkflow(t *testing.T) {
    // 空 workflow.md 返回默认 4 阶段
}

func TestParser_CustomPhases(t *testing.T) {
    // 自定义阶段名称
}

func TestParser_PhaseOrder(t *testing.T) {
    // 阶段顺序正确
}

func TestParser_MalformedWorkflow(t *testing.T) {
    // 格式错误的 workflow.md 优雅降级
}
```

## 验收标准

- [ ] `TestE2E_WorkflowFullCycle` 覆盖 4 阶段完整流转
- [ ] `TestE2E_WorkflowInvalidPhaseTransition` 验证阶段跳转拒绝
- [ ] `TestE2E_WorkflowStateInjection` 验证状态注入输出
- [ ] 单元测试补充 4 个边界用例

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/commands/workflow.integration.test.ts`
- trellis-go 现有：`pkg/workflow/parser_test.go`

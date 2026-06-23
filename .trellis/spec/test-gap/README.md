# trellis-go Test Gap — Spec Index

> 基于 trellis-go vs 原版 trellis (TypeScript) 测试对比分析，规划测试补齐方案。

## 对比基准

| 维度 | trellis-go | trellis (原版) |
|------|:---:|:---:|
| 测试文件 | 27 `_test.go` | 67 `.test.ts` |
| E2E 场景 | 15 | 大量集成测试 |
| 不变量测试 | ❌ | ✅ 11 项 |
| 回归测试 | ❌ | ✅ 6572 行 |
| Update/Upgrade E2E | ❌ | ✅ 1652 行 |

## Spec 清单

| # | Spec | 分类 | 优先级 | 估算 |
|---|------|------|:---:|:---:|
| 1 | [test-registry-invariants](./test-registry-invariants.md) | 不变量测试 | 🔴 P0 | 0.3 周 |
| 2 | [test-regression](./test-regression.md) | 回归测试套件 | 🔴 P0 | 0.5 周 |
| 3 | [test-e2e-update-upgrade](./test-e2e-update-upgrade.md) | Update/Upgrade E2E | 🔴 P0 | 0.3 周 |
| 4 | [test-e2e-workflow](./test-e2e-workflow.md) | Workflow E2E | 🟡 P1 | 0.3 周 |
| 5 | [test-e2e-session-journal](./test-e2e-session-journal.md) | Session Journal E2E | 🟡 P1 | 0.2 周 |
| 6 | [test-e2e-subtask](./test-e2e-subtask.md) | Subtask E2E | 🟡 P1 | 0.2 周 |
| 7 | [test-e2e-hook](./test-e2e-hook.md) | Hook E2E | 🟡 P1 | 0.2 周 |
| 8 | [test-e2e-init-multi-scenario](./test-e2e-init-multi-scenario.md) | Init 多场景集成 | 🔵 P2 | 0.3 周 |
| 9 | [test-template-integrity](./test-template-integrity.md) | 模板内容完整性 | 🔵 P2 | 0.2 周 |
| 10 | [test-unit-coverage](./test-unit-coverage.md) | 单元测试覆盖率 | 🔵 P2 | 0.5 周 |

## 总估算

| 优先级 | Spec 数 | 工作量 |
|:---|:---:|:---|
| 🔴 P0 | 3 | ~1.1 周 |
| 🟡 P1 | 4 | ~0.9 周 |
| 🔵 P2 | 3 | ~1.0 周 |
| **合计** | **10** | **~3.0 周** |

## 参考

- 原版 trellis 测试：`/tmp/trellis-orig/packages/cli/test/`、`/tmp/trellis-orig/packages/core/test/`
- trellis-go 现有测试：`cmd/trellis/e2e_test.go`、`pkg/*/`
- 对比报告：`dev-loop/reports/03-test-gap-analysis.md`

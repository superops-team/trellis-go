# test-regression: 回归测试套件

> 参考原版 trellis `regression.test.ts`（6572 行），建立 trellis-go 回归测试体系。

## 当前状态

trellis-go 没有任何回归测试。已修复的 bug 可能在未来重构中复现。

## 需求

### 1. 回归测试框架

在 `pkg/platform/regression_test.go` 中建立回归测试模式：

```go
// 每个回归测试标注引入/修复版本
func TestRegression_v0_4_0_PlatformCount(t *testing.T) {
    // 回归：v0.4.0 引入 ZCode 后平台数应为 16
    reg := BuiltinRegistry()
    if len(reg.All()) != 16 {
        t.Errorf("platform count changed: got %d, want 16", len(reg.All()))
    }
}
```

### 2. 平台注册表回归

| 测试 | 版本 | 检查内容 |
|------|:---:|------|
| PlatformCount | v0.4.0 | 平台总数 = 16 |
| DevinAlias | v0.4.0 | Devin 有 `--windsurf` 别名 |
| ZCodeExists | v0.4.0 | ZCode 平台已注册 |
| NoDuplicateFlags | v0.4.0 | 无重复 CLI flag |
| ForFlagResolvesAlias | v0.4.0 | `ForFlag("windsurf")` 返回 Devin |

### 3. 模板内容回归

| 测试 | 版本 | 检查内容 |
|------|:---:|------|
| HookScriptHasShebang | v0.4.0 | 所有 hook 脚本有 `#!/usr/bin/env` |
| HookScriptExecutable | v0.4.0 | 生成 hook 有执行权限 |
| ConfigYAMLHasVersion | v0.4.0 | config.yaml 包含 version 字段 |
| TemplateCommonExists | v0.4.0 | common 模板目录存在 |

### 4. CLI 行为回归

| 测试 | 版本 | 检查内容 |
|------|:---:|------|
| InitCreatesDotTrellis | v0.4.0 | init 创建 .trellis 目录 |
| TaskCreateGeneratesID | v0.4.0 | task create 生成有效 ID |
| ContextAddRejectsUnsafe | v0.4.0 | 拒绝 ../ 和绝对路径 |
| UninstallKeepsTasks | v0.4.0 | --keep-tasks 保留任务文件 |

### 5. 路径/编码回归

trellis-go 不需要原版的 Windows 编码回归（Go 原生 UTF-8），但需要：

| 测试 | 版本 | 检查内容 |
|------|:---:|------|
| PathsAreRelative | v0.4.0 | 所有路径操作用相对路径 |
| NoHardcodedHomeDir | v0.4.0 | 无硬编码 `/home/` 路径 |
| TempDirCleanup | v0.4.0 | E2E 测试后清理临时目录 |

## 验收标准

- [ ] 建立 `regression_test.go` 文件，每个测试标注版本
- [ ] 覆盖 5 大类回归：平台注册表、模板内容、CLI 行为、路径/编码、E2E 清理
- [ ] 回归测试在 CI 中自动运行
- [ ] 新增 bug 修复时同步添加回归测试

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/regression.test.ts`（6572 行）
- trellis-go 现有 E2E：`cmd/trellis/e2e_test.go`

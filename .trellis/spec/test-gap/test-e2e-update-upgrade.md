# test-e2e-update-upgrade: Update/Upgrade 端到端测试

> 参考原版 trellis `update.integration.test.ts`（1265 行）、`update-internals.test.ts`（282 行）、`upgrade.test.ts`（105 行）。

## 当前状态

trellis-go 有 `pkg/update/syncer_test.go` 和 `pkg/upgrade/checker_test.go` 单元测试，但无 E2E。

`trellis update` 和 `trellis upgrade` 端到端行为未验证。

## 需求

### 1. Update E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_UpdateSyncsTemplates(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 删除某个模板文件模拟过期
    // 3. trellis update 恢复缺失文件
    // 4. 验证文件已恢复
}

func TestE2E_UpdatePreservesUserEdits(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 修改某个 hook 文件
    // 3. trellis update
    // 4. 验证用户修改未被覆盖
}

func TestE2E_UpdateDryRun(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 删除某个模板文件
    // 3. trellis update --dry-run
    // 4. 验证文件未被实际恢复（dry-run 不写入）
}

func TestE2E_UpdateSkipPaths(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 在 config.yaml 设置 update.skip
    // 3. 删除被 skip 的文件
    // 4. trellis update
    // 5. 验证被 skip 的文件未被恢复
}

func TestE2E_UpdateConfigSectionAppend(t *testing.T) {
    // 1. trellis init 创建项目
    // 2. 从 config.yaml 删除某个配置段
    // 3. trellis update
    // 4. 验证缺失配置段被追加
    // 5. 再次 update 验证幂等（不重复追加）
}
```

### 2. Upgrade E2E

```go
func TestE2E_UpgradeChecksVersion(t *testing.T) {
    // 1. trellis upgrade --dry-run
    // 2. 验证输出版本检查信息
    // 3. 不实际下载/替换二进制
}

func TestE2E_UpgradeCurrentVersion(t *testing.T) {
    // 1. trellis upgrade
    // 2. 当前已是最新版本时输出提示
}
```

### 3. Update 内部逻辑单元测试

在 `pkg/update/syncer_test.go` 中补充：

```go
func TestSyncer_ConfigSectionAppend(t *testing.T) {
    // 测试配置段追加逻辑
}

func TestSyncer_SkipPaths(t *testing.T) {
    // 测试 skip 路径过滤
}

func TestSyncer_Idempotent(t *testing.T) {
    // 测试多次 update 幂等
}
```

## 验收标准

- [ ] `TestE2E_UpdateSyncsTemplates` 通过
- [ ] `TestE2E_UpdatePreservesUserEdits` 通过
- [ ] `TestE2E_UpdateDryRun` 通过
- [ ] `TestE2E_UpdateSkipPaths` 通过
- [ ] `TestE2E_UpdateConfigSectionAppend` 通过
- [ ] `TestE2E_UpgradeChecksVersion` 通过
- [ ] `TestE2E_UpgradeCurrentVersion` 通过
- [ ] Update 单元测试补充 3 个用例

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/commands/update.integration.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/commands/update-internals.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/commands/upgrade.test.ts`
- trellis-go 现有：`pkg/update/syncer_test.go`、`pkg/upgrade/checker_test.go`

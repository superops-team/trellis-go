# test-e2e-init-multi-scenario: Init 多场景集成测试

> 参考原版 trellis 4 个 init 测试文件（共 2269 行）。

## 当前状态

trellis-go 只有 1 个 `TestE2E_InitNewProject`，覆盖基础 init 流程。

原版 trellis 有 4 个 init 测试文件：`init.integration.test.ts`（1338 行）、`init-internals.test.ts`（207 行）、`init-joiner.integration.test.ts`（361 行）、`init-uninstall-overdelete.integration.test.ts`（363 行）。

## 需求

### 1. Init 多场景 E2E

在 `cmd/trellis/e2e_test.go` 中新增：

```go
func TestE2E_InitSinglePlatform(t *testing.T) {
    // 1. trellis init --platform claude
    // 2. 验证只生成 claude 平台文件
    // 3. 验证不生成其他平台文件
}

func TestE2E_InitAllPlatforms(t *testing.T) {
    // 1. trellis init（默认所有平台）
    // 2. 验证所有已注册平台的文件都生成
}

func TestE2E_InitWithCustomName(t *testing.T) {
    // 1. trellis init --name "my-project"
    // 2. 验证 config.yaml 中 name 字段为 "my-project"
}

func TestE2E_InitOverwriteProtection(t *testing.T) {
    // 1. trellis init
    // 2. 再次 trellis init
    // 3. 验证不覆盖已有文件，或提示确认
}

func TestE2E_InitCreatesCorrectDirectoryStructure(t *testing.T) {
    // 1. trellis init
    // 2. 验证目录结构：
    //    .trellis/
    //    .trellis/config.yaml
    //    .trellis/spec/
    //    .trellis/templates/
    //    .agents/skills/
    //    .claude/
    //    .cursor/
    //    ... 等平台目录
}

func TestE2E_InitGeneratesValidConfigYAML(t *testing.T) {
    // 1. trellis init
    // 2. 解析 .trellis/config.yaml
    // 3. 验证必填字段：version, name, platforms
}

func TestE2E_InitHookFilesExecutable(t *testing.T) {
    // 1. trellis init --platform claude
    // 2. 验证 .claude/hooks/session-start.sh 有执行权限
}
```

### 2. Init 单元测试补充

在 `pkg/configurator/configurator_test.go` 中补充：

```go
func TestConfigurator_InitDirectoryStructure(t *testing.T) { ... }
func TestConfigurator_InitConfigYAML(t *testing.T) { ... }
func TestConfigurator_InitOverwriteBehavior(t *testing.T) { ... }
```

## 验收标准

- [ ] 7 个 Init 多场景 E2E 通过
- [ ] 3 个 Init 单元测试通过
- [ ] 覆盖单平台、全平台、覆盖保护、目录结构、权限

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/commands/init.integration.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/commands/init-internals.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/commands/init-joiner.integration.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/commands/init-uninstall-overdelete.integration.test.ts`

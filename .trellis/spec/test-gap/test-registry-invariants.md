# test-registry-invariants: 平台注册表不变量测试

> 参考原版 trellis `registry-invariants.test.ts`，补齐 trellis-go 平台注册表一致性检查。

## 当前状态

trellis-go 没有任何注册表不变量测试。新增平台时容易遗漏更新，产生隐蔽 bug。

原版 trellis 有 196 行不变量测试，覆盖 11 项检查。

## 需求

### 1. 注册表内部一致性

在 `pkg/platform/registry_test.go` 中新增：

```go
func TestRegistry_InternalConsistency(t *testing.T) {
    reg := BuiltinRegistry()
    platforms := reg.All()

    // 1. 每个平台有非空 Name
    // 2. 每个平台有非空 CLIFlag
    // 3. 每个平台有非空 ConfigDir
    // 4. 所有 CLIFlag 唯一
    // 5. 所有 ConfigDir 唯一
    // 6. 所有 ConfigDir 以 "." 开头
    // 7. 无 ConfigDir 与 ".trellis" 冲突
    // 8. 无 CLIFlag 与保留字冲突（help, version, V, h）
    // 9. supportsAgentSkills 平台不用 ".agents/skills" 作 ConfigDir
    // 10. 每个平台 TemplateDirs 包含 "common"
}
```

### 2. 别名一致性

```go
func TestRegistry_AliasConsistency(t *testing.T) {
    reg := BuiltinRegistry()

    // 1. 所有 Alias 唯一（不与其他平台 CLIFlag 冲突）
    // 2. 所有 Alias 唯一（不与其他平台 Alias 冲突）
    // 3. Alias 不与保留字冲突
    // 4. ForFlag() 能通过别名找到正确平台
}
```

### 3. 平台类型一致性

```go
func TestRegistry_PlatformTypeConsistency(t *testing.T) {
    reg := BuiltinRegistry()

    // 1. Pull-based 平台有非空 ConfigDir
    // 2. Push-based 平台在 Generator 中不被拒绝
    // 3. 每个平台 Class 是有效值（push/pull/agentless）
}
```

## 验收标准

- [ ] `TestRegistry_InternalConsistency` 覆盖 10 项不变量
- [ ] `TestRegistry_AliasConsistency` 覆盖 4 项别名检查
- [ ] `TestRegistry_PlatformTypeConsistency` 覆盖 3 项类型检查
- [ ] 新增平台时自动触发检查（测试失败即提醒更新）

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/registry-invariants.test.ts`
- trellis-go 现有：`pkg/platform/registry_test.go`

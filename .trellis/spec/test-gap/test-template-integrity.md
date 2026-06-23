# test-template-integrity: 模板内容完整性测试

> 参考原版 trellis 12 个平台模板测试（`templates/*.test.ts`）。

## 当前状态

trellis-go 有 `pkg/template/engine_test.go`，只测试模板引擎机制（变量替换、条件渲染），不验证模板内容。

原版 trellis 有 12 个平台模板测试文件，验证每个平台的 hook 内容、settings 配置、模板完整性。

## 需求

### 1. 模板内容验证

在 `pkg/template/content_test.go` 中新增：

```go
func TestTemplateContent_CommonExists(t *testing.T) {
    // 验证 common 模板目录存在且非空
}

func TestTemplateContent_AllPlatformsHaveHookTemplates(t *testing.T) {
    // 验证每个 pull-based 平台有 hook 模板
}

func TestTemplateContent_AllPlatformsHaveConfigTemplates(t *testing.T) {
    // 验证每个平台有 settings/config 模板
}

func TestTemplateContent_HookScriptHasShebang(t *testing.T) {
    // 验证所有 hook 脚本模板以 #!/usr/bin/env 开头
}

func TestTemplateContent_NoHardcodedPaths(t *testing.T) {
    // 验证模板中无硬编码路径（用 {{.ProjectDir}} 等变量）
}

func TestTemplateContent_ConfigYAMLTemplate(t *testing.T) {
    // 验证 config.yaml 模板包含必填字段
    // version, name, platforms, update.skip
}
```

### 2. 模板渲染输出验证

```go
func TestTemplateRender_ClaudeSettingsJSON(t *testing.T) {
    // 渲染 claude settings.json 模板
    // 验证输出是合法 JSON
    // 验证包含 hooks 配置
}

func TestTemplateRender_CursorSettingsJSON(t *testing.T) {
    // 渲染 cursor settings.json 模板
    // 验证输出是合法 JSON
}

func TestTemplateRender_GeminiSettingsYAML(t *testing.T) {
    // 渲染 gemini settings 模板
    // 验证输出是合法 YAML
}

func TestTemplateRender_AllPlatformsRenderWithoutError(t *testing.T) {
    // 遍历所有平台，渲染所有模板
    // 验证无错误
}
```

### 3. Spec 模板验证

```go
func TestSpecTemplates_Exist(t *testing.T) {
    // 验证 .trellis/templates/ 下 3 个 spec 模板存在
}

func TestSpecTemplates_ValidMarkdown(t *testing.T) {
    // 验证 spec 模板是合法 Markdown（有标题、有内容）
}
```

### 4. Skill 模板验证

```go
func TestSkillTemplates_Exist(t *testing.T) {
    // 验证 .agents/skills/ 下 4 个 skill 存在
}

func TestSkillTemplates_ValidSKILLMD(t *testing.T) {
    // 验证每个 skill 的 SKILL.md 有 name 字段
}
```

## 验收标准

- [ ] 6 个模板内容验证通过
- [ ] 4 个模板渲染输出验证通过
- [ ] 2 个 Spec 模板验证通过
- [ ] 2 个 Skill 模板验证通过
- [ ] 新增平台时自动验证模板完整性

## 参考

- 原版：`/tmp/trellis-orig/packages/cli/test/templates/claude.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/templates/codex.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/templates/copilot.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/templates/cursor.test.ts`
- 原版：`/tmp/trellis-orig/packages/cli/test/templates/shared-hooks.test.ts`
- trellis-go 现有：`pkg/template/engine_test.go`

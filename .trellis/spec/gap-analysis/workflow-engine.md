# workflow-engine: 工作流引擎

> 实现完整的 workflow.md 解析和运行时注入系统。

## 当前状态

| 组件 | 状态 | 说明 |
|------|:---:|------|
| `workflow.Parser` | ⚠️ | 返回默认状态机，不解析自定义内容 |
| `workflow.StateMachine` | ⚠️ | 只有 4 个硬编码状态 |
| `[workflow-state:STATUS]` breadcrumb | ⚠️ | 只支持 4 个内置状态 |
| Skill routing table | ❌ | 未解析 |
| Phase/step 解析 | ❌ | 未实现 |
| 自定义状态 | ❌ | 不支持 |
| `get_context.py --mode phase` | ❌ | 未实现 |

## 需求

### 1. workflow.md 完整解析

**解析目标**:

```markdown
## Phase Index
[workflow-state:no_task]...[/workflow-state:no_task]
[workflow-state:planning]...[/workflow-state:planning]
[workflow-state:in_progress]...[/workflow-state:in_progress]
[workflow-state:completed]...[/workflow-state:completed]

### Phase 1: Plan
#### 1.1 Brainstorm
#### 1.2 Write PRD

### Phase 2: Implement
#### 2.1 Before Dev
#### 2.2 Write Code
#### 2.3 Verify

### Phase 3: Finish
#### 3.1 Update Spec
#### 3.2 Commit
#### 3.3 Archive

### Skill Routing
| User intent | Skill |
|-------------|-------|
| ... | ... |

### DO NOT skip skills
...

### Task System
...
```

**数据结构**:

```go
type Workflow struct {
    Phases       []Phase
    Breadcrumbs  map[State]string      // state → breadcrumb text
    SkillRouting []SkillRoute
    DoNotSkip    string
    TaskSystem   string
}

type Phase struct {
    Number      int
    Title       string
    Steps       []Step
}

type Step struct {
    Number      string  // "1.1"
    Title       string
    Body        string
}

type SkillRoute struct {
    Intent string
    Skill  string
}
```

### 2. Breadcrumb 注入

**规则**:
- 每个 `UserPromptSubmit` 事件注入对应状态的 breadcrumb
- 状态来自 `task.json.status`
- 默认状态：`planning`, `in_progress`, `completed`, `no_task`
- 支持自定义状态（如 `blocked`, `in_review`）
- 未知状态：注入通用提示 "Refer to workflow.md for current step."
- 每个 breadcrumb 保持简短（~200 bytes）

**注入格式**:
```
<workflow-state>
You are in the IMPLEMENT phase. Flow: implement → check → update-spec → finish.
Check conversation history + git status to determine current step; do NOT skip check.
</workflow-state>
```

### 3. Phase/Step 上下文提取

**等价于**: `get_context.py --mode phase --step X.Y`

**功能**:
- 根据 `## Phase X` 和 `#### X.Y` 标题定位步骤
- 提取步骤正文内容
- 注入到 agent 上下文

### 4. Skill Routing Table

**解析** `### Skill Routing` 下的 Markdown 表格：
- 映射用户意图 → 技能名称
- 在 SessionStart 时注入给 AI
- AI 根据此表决定何时触发技能

### 5. 自定义工作流支持

**能力**:
- 添加新 Phase（如 Phase 4: Review）
- 拆分 Phase（1A/1B 分支）
- 删除不需要的步骤
- 自定义 breadcrumb 文本
- 自定义状态名称

**约束**:
- 不可修改：`[workflow-state:STATUS]` 标签格式
- 不可修改：`## Phase X` + `#### X.Y` 标题深度
- 不可修改：`task.py` 子命令名称

## 验收标准

- [ ] 自定义 `workflow.md` 被完整解析
- [ ] 自定义 breadcrumb 文本正确注入
- [ ] 自定义状态（如 `blocked`）正确匹配
- [ ] Skill routing table 正确解析
- [ ] Phase/step 上下文提取正确
- [ ] 未知状态回退到通用提示
- [ ] `trellis hook inject-workflow-state --state blocked` 返回自定义 breadcrumb

## 依赖

- [context-system](./context-system.md) — breadcrumb 注入依赖上下文系统
- [platform-configurators](./platform-configurators.md) — hook 脚本使用此引擎

## 关联 Spec

- [auto-skills](./auto-skills.md) — skill routing table 触发技能
- [task-lifecycle](./task-lifecycle.md) — 状态变更触发 breadcrumb 切换

# cli-core: CLI 命令完整性

> 补齐 trellis CLI 所有缺失命令，达到 Python 版功能对等。

## 当前状态

| 命令 | 状态 | 说明 |
|------|:---:|------|
| `init` | ✅ | 基础实现，缺多人协作模式 |
| `update` | ❌ | Stub，只打印 "not yet implemented" |
| `uninstall` | ✅ | 基础实现 |
| `task create` | ✅ | 基础实现 |
| `task start` | ✅ | 基础实现 |
| `task archive` | ✅ | 基础实现 |
| `task list` | ✅ | 只打印 ID |
| `task current` | ❌ | Stub，只打印 "No active task" |
| `context add` | ✅ | 基础实现 |
| `context build` | ✅ | 基础实现 |
| `hook session-start` | ✅ | 基础实现 |
| `hook inject-context` | ✅ | 基础实现 |
| `hook inject-workflow-state` | ✅ | 基础实现 |
| `version` | ✅ | 基础实现 |
| `upgrade` | ❌ | 完全缺失 |

## 需求

### 1. `trellis update` 完整实现

**功能**:
- 从嵌入式文件系统同步模板到项目目录
- 保留用户本地编辑（不覆盖已有文件，除非是模板更新）
- 支持 `--dry-run` 预览模式
- 支持 `--migrate` 强制迁移模式
- `update.skip` 配置支持：跳过指定路径
- 配置段追加机制：通过 sentinel 文本检测新配置段，追加到 `config.yaml`
- 输出变更摘要

**配置段追加逻辑**:
```
1. 读取当前 config.yaml
2. 对每个 configSectionsAdded 中的 sentinel，检查是否已存在于文件中
3. 如果缺失，从模板追加对应配置段
4. 幂等：已存在则跳过
```

**验收标准**:
- [ ] `trellis update` 同步所有平台模板到项目
- [ ] `trellis update --dry-run` 预览不写入
- [ ] 已有用户修改的文件不被覆盖
- [ ] 新配置段正确追加到 config.yaml
- [ ] `update.skip` 路径被正确排除

### 2. `trellis upgrade` 自升级

**功能**:
- 检查 npm registry 最新版本（或 GitHub releases）
- 下载并替换当前二进制
- 支持 `--tag latest|beta|rc` 版本选择
- 支持 `--dry-run` 预览
- 升级后提示运行 `trellis update`

**验收标准**:
- [ ] `trellis upgrade` 升级到最新版本
- [ ] `trellis upgrade --dry-run` 只打印不执行
- [ ] 升级后二进制可正常执行

### 3. `task` 子命令补齐

**缺失子命令**:

| 子命令 | 功能 |
|--------|------|
| `task info <id>` | 显示任务详情（名称、状态、分配人、分支、子任务、时间） |
| `task edit <id> --name/--assignee/--branch` | 编辑任务字段 |
| `task add-subtask <id> <title>` | 添加子任务 |
| `task done-subtask <id> <subtask-id>` | 标记子任务完成 |
| `task add-context <id> <file> --phase/--required/--description` | 添加上下文 |
| `task remove-context <id> <file>` | 移除上下文 |
| `task list-context <id>` | 列出上下文 |
| `task list --status <status>` | 按状态过滤 |
| `task list --format json\|table` | 输出格式 |
| `task current` | 读取 `.runtime/sessions/` 返回活跃任务 |

**验收标准**:
- [ ] 所有子命令可用且输出格式一致
- [ ] `task current` 正确读取活跃任务状态
- [ ] `task list --status` 过滤正确

### 4. `init` 多人协作增强

**功能**:
- 检测已存在的 `.trellis/` 目录
- 提供交互式选择：
  1. Add AI platform(s) — 添加新平台
  2. Set up developer identity — 写入 `.trellis/.developer`
  3. Full re-initialize — 完全重置
- 非交互模式支持 `--platform` 和 `--developer` 参数

**验收标准**:
- [ ] 已存在项目运行 `init` 显示 3 个选项
- [ ] "Set up developer identity" 正确写入 `.trellis/.developer`
- [ ] "Add AI platform(s)" 只添加新平台不覆盖已有配置

## 依赖

- `update` 依赖 `internal/embed` 模板系统
- `upgrade` 依赖 GitHub releases API 或 npm registry
- `task current` 依赖 session 追踪系统

## 关联 Spec

- [platform-configurators](./platform-configurators.md) — update 的模板同步
- [session-journal](./session-journal.md) — task current 的 session 追踪

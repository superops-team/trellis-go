<div align="center">

# Trellis (Go)

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](https://github.com/superops-team/trellis-go/actions)

[English](README.md) | [中文](README.zh-CN.md)

</div>

> 面向 AI 编程的工程框架。将规范、任务和记忆持久化到代码仓库中，让任何编码智能体都遵循你的工程标准。

Trellis 是原版 [TypeScript/Python Trellis](https://github.com/mindfold-ai/Trellis) 框架的 Go 语言移植版本，专为单二进制分发、高性能和与 15+ AI 编程平台无缝集成而设计。

## 目录

- [特性](#特性)
- [快速开始](#快速开始)
- [小白使用指南](#小白使用指南)
- [支持的平台](#支持的平台)
- [架构](#架构)
- [CLI 命令](#cli-命令)
- [工作流](#工作流)
- [开发](#开发)
- [测试](#测试)
- [贡献](#贡献)
- [许可证](#许可证)

## 特性

- **15+ AI 平台** — 内置支持 Claude、Cursor、Codex、Copilot、Windsurf 等
- **四阶段工作流** — Plan → Implement → Verify → Finish 状态机
- **任务生命周期** — 创建、启动、归档任务，自动组织管理
- **上下文构建器** — 基于 JSONL 的清单系统，用于 AI 智能体上下文注入
- **平台钩子** — 自动生成平台特定的配置文件
- **原子文件操作** — 安全的并发写入，支持 SHA256 哈希校验
- **单二进制** — 静态编译的 Go 二进制文件，内置模板资源
- **线程安全注册表** — 并发安全的平台和规范管理

## 快速开始

### 安装

```bash
go install github.com/superops-team/trellis-go/cmd/trellis@latest
```

或从 [Releases](https://github.com/superops-team/trellis-go/releases) 下载预编译二进制文件。

### 初始化项目

```bash
# 在 Git 仓库内
git init my-project && cd my-project

# 使用默认平台（Claude）初始化 Trellis
trellis init --developer alice

# 或使用多个平台初始化
trellis init --developer alice --platform claude --platform cursor --platform codex
```

### 创建第一个任务

```bash
# 创建新任务
trellis task create user-auth

# 列出所有任务
trellis task list

# 查看当前活跃任务
trellis task current
```

### 项目结构

```
my-project/
├── .git/
├── .trellis/                 # Trellis 工作区
│   ├── config.yaml           # 开发者与平台配置
│   ├── .version              # Trellis 版本
│   ├── workflow.md           # 四阶段工作流定义
│   ├── spec/                 # 工程规范
│   ├── tasks/                # 活跃任务
│   │   ├── 06-13-user-auth/
│   │   │   ├── task.json     # 任务元数据
│   │   │   ├── prd.md        # 产品需求文档
│   │   │   ├── implement.jsonl  # 上下文清单
│   │   │   ├── check.jsonl      # 验证清单
│   │   │   └── research/     # 调研笔记
│   │   └── archive/          # 已归档任务 (YYYY-MM/)
│   ├── workspace/            # 共享工作区
│   └── .runtime/sessions/    # 活跃会话追踪
├── .claude/                  # Claude 专属配置
├── .cursor/                  # Cursor 专属配置
└── ...
```

## 小白使用指南

如果你是第一次使用 Trellis，建议先阅读完整的分步指南：

- [中文小白使用指南](docs/USAGE.zh-CN.md)
- [English Beginner Guide](docs/USAGE.md)

指南包含安装、初始化、任务生命周期、上下文清单、常见问题，以及 Mermaid 流程图。

## 支持的平台

| 平台 | 类别 | 智能体 | 钩子 |
|------|------|--------|------|
| Claude Code | Push-based | 支持 | 支持 |
| Cursor | Push-based | 支持 | 支持 |
| Codex | Pull-based | 支持 | 支持 |
| OpenCode | Push-based | 支持 | 支持 |
| Gemini CLI | Pull-based | 支持 | 不支持 |
| Kiro | Push-based | 支持 | 支持 |
| Copilot | Pull-based | 不支持 | 不支持 |
| Windsurf | Agentless | 不支持 | 不支持 |
| Kilo | Agentless | 不支持 | 不支持 |
| Pi | Push-based | 支持 | 支持 |
| CodeBuddy | Push-based | 支持 | 支持 |
| Droid | Push-based | 支持 | 支持 |
| Qoder | Pull-based | 支持 | 不支持 |
| Antigravity | Agentless | 不支持 | 不支持 |
| Reasonix | Push-based | 支持 | 支持 |

**平台类别说明：**
- **Push-based** — 智能体主动发起执行（Claude、Cursor 等）
- **Pull-based** — IDE 按需拉取上下文（Codex、Copilot 等）
- **Agentless** — 无智能体能力，手动工作流（Windsurf、Kilo 等）

## 架构

```
cmd/trellis/          CLI 命令 (cobra)
pkg/
  platform/           平台定义与注册表
  fsutil/             原子文件操作与哈希
  config/             YAML 配置管理
  template/           embed.FS 模板引擎
  task/               任务生命周期与清单
  workflow/           四阶段状态机
  context/            上下文构建器 (JSONL 清单)
  hook/               平台钩子生成器
  spec/               规范加载器与索引
  git/                Git 命令封装
internal/
  embed/              嵌入式模板资源
  testutil/           测试辅助工具
```

## CLI 命令

```bash
# 在当前仓库初始化 Trellis
trellis init [flags]
  --developer, -u    开发者名称
  --platform, -p     要配置的平台（可重复指定）
  --all              配置所有支持的平台

# 任务管理
trellis task create <name>     创建新任务
trellis task list              列出所有任务
trellis task current           显示当前活跃任务
trellis task start <id>        启动任务
trellis task archive <id>      归档已完成任务

# 上下文管理
trellis context add <file> --task <id> [--phase implement|check]
trellis context build --task <id> --phase implement|check
trellis context build --phase research

# 维护
trellis uninstall              移除 Trellis（使用 --keep-tasks 保留任务）
trellis version                显示版本
```

## 工作流

Trellis 通过 `workflow.md` 强制执行四阶段工程工作流：

```
[workflow-state:PLAN]
→ 头脑风暴需求并编写 PRD

[workflow-state:IMPLEMENT]
→ 根据 PRD 编写代码

[workflow-state:VERIFY]
→ 对照规范审查代码并运行检查

[workflow-state:FINISH]
→ 归档任务并更新日志
```

状态转换由状态机验证：
- `plan` → `implement`
- `implement` → `verify`
- `verify` → `implement` | `finish`

## 开发

### 前置条件

- Go 1.23+
- Git

### 构建

```bash
go build -o trellis ./cmd/trellis
```

### 运行测试

```bash
# 全部测试
go test ./...

# 仅单元测试
go test ./pkg/...

# E2E 测试
go test ./cmd/trellis -run TestE2E_
```

## 测试

项目包含全面的测试覆盖：

- **单元测试** — 所有 `pkg/` 包，采用表驱动测试
- **E2E 测试** — 8 个真实场景：
  1. 新项目初始化 + 首次任务创建
  2. 多平台配置（Claude + Cursor + Codex）
  3. 完整任务生命周期（创建 → 启动 → 归档）
  4. AI 智能体上下文构建与注入
  5. 任务列表与当前任务查询
  6. 卸载时保留任务
  7. 无效平台错误处理
  8. 非 Git 仓库错误处理

## 贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 发起 Pull Request

请确保所有测试通过，并遵循现有代码风格。

## 许可证

[MIT](LICENSE)

---

<div align="center">

用 Go 为 AI 原生工程团队打造。

</div>

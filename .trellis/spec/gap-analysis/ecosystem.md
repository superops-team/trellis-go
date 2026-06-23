# ecosystem: 生态与平台扩展

> 补齐缺失平台、文档站点、OpenAPI、skills marketplace、spec templates。

## 当前状态

| 组件 | 状态 | 说明 |
|------|:---:|------|
| 平台数 | ⚠️ | 15 个，缺 ZCode，Devin 命名过时 |
| 文档站点 | ❌ | 只有 README.md |
| OpenAPI Spec | ❌ | 无 |
| Skills marketplace | ❌ | 无 |
| Spec templates | ❌ | 无 |
| 社区贡献指南 | ❌ | 无 |

## 需求

### 1. 新增平台

#### ZCode
```go
{
    ID:                  "zcode",
    Name:                "ZCode",
    ConfigDir:           ".zcode",
    TemplateDirs:        []string{"common", "zcode"},
    AgentCapable:        true,
    HasHooks:            true,
    SupportsAgentSkills: true,
    CLIFlag:             "--zcode",
    Class:               ClassPushBased,
}
```

#### Devin（重命名 Windsurf）
- 将现有 "windsurf" 平台重命名为 "devin"
- ConfigDir: `.devin`
- 保持向后兼容（接受 `--windsurf` 作为别名）

### 2. 文档站点

**结构**（对标 `docs.trytrellis.app`）:

```
docs/
├── index.md                          # 首页
├── start/
│   ├── install-and-first-task.md     # 安装与第一个任务
│   ├── how-it-works.md               # 工作原理
│   ├── everyday-use.md               # 日常使用
│   └── real-world-scenarios.md       # 真实场景
├── advanced/
│   ├── architecture.md               # 架构概述
│   ├── configuration.md              # 配置指南
│   ├── multi-platform.md             # 多平台配置
│   ├── custom-agents.md              # 自定义子代理
│   ├── custom-skills.md              # 自定义技能
│   ├── custom-commands.md            # 自定义命令
│   ├── custom-hooks.md               # 自定义 hooks
│   ├── custom-workflow.md            # 自定义工作流
│   ├── custom-spec-template-marketplace.md
│   ├── appendix-a.md                 # 关键路径参考
│   ├── appendix-b.md                 # 命令/技能速查表
│   ├── appendix-c.md                 # task.json Schema
│   ├── appendix-d.md                 # JSONL 格式参考
│   ├── appendix-f.md                 # FAQ
│   ├── roadmap.md                    # 路线图
│   └── resources.md                  # 资源与致谢
├── blog/
│   ├── index.md
│   ├── ai-collaborative-dev-system.md
│   └── use-k8s-to-know-trellis.md
├── changelog/                        # 版本日志
├── showcase/
│   ├── index.md
│   ├── terminal-demo.md
│   ├── trellis-cursor.md
│   └── open-typeless.md
├── templates/
│   ├── specs-index.md
│   ├── specs-nextjs.md
│   ├── specs-cf-workers.md
│   └── specs-electron.md
├── skills-market/
│   ├── index.md
│   ├── frontend-fullchain-optimization.md
│   ├── mem-recall.md
│   ├── trellis-meta.md
│   └── trellis-spec-bootstrap.md
├── contribute/
│   ├── trellis.md
│   └── docs.md
├── use-cases/
│   └── open-typeless.md
├── api-reference/
│   └── openapi.json
└── llms.txt                          # 文档索引
```

### 3. OpenAPI Spec

**用途**: 提供 Trellis CLI 的 API 参考（如果未来有 server 模式）

**内容**:
```json
{
  "openapi": "3.0.0",
  "info": {
    "title": "Trellis API",
    "version": "0.1.0"
  },
  "paths": {}
}
```

### 4. Spec Templates

#### Next.js + oRPC + PostgreSQL
- 前端：Next.js App Router, React, TypeScript
- API：oRPC
- 数据库：PostgreSQL + Drizzle ORM
- 包含：项目结构约定、命名规范、测试策略、部署清单

#### Cloudflare Workers + Hono + Turso
- 运行时：Cloudflare Workers
- 框架：Hono
- 数据库：Turso (SQLite)
- 包含：边缘计算最佳实践、冷启动优化

#### Electron + React + TypeScript
- 桌面框架：Electron
- 前端：React + TypeScript
- 包含：IPC 通信约定、安全模型、打包配置

### 5. Skills Marketplace

**社区技能**:
- `frontend-fullchain-optimization` — Web Vitals 驱动的前端性能优化
- `mem-recall` — 跨平台 AI 对话回忆（通过 trellis mem）
- `trellis-meta` — Trellis 自定义核心技能
- `trellis-spec-bootstrap` — 从代码库引导生成 spec

### 6. 社区贡献

**CONTRIBUTING.md**:
- 开发环境搭建
- 代码风格指南
- PR 流程
- Issue 模板

**文档贡献**:
- 文档站点仓库（`superops-team/trellis-go-docs`）
- 文档贡献流程

## 验收标准

- [ ] ZCode 平台注册并可用
- [ ] Devin 重命名完成，`--windsurf` 别名可用
- [ ] 文档站点至少包含 start/advanced 核心页面
- [ ] OpenAPI spec 文件存在
- [ ] 3 个 spec template 可用
- [ ] 4 个 marketplace skill 可用
- [ ] CONTRIBUTING.md 存在

## 依赖

- [platform-configurators](./platform-configurators.md) — 新平台配置器
- [auto-skills](./auto-skills.md) — marketplace skills

## 工作量

~1.5 周

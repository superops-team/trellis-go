# trellis-go Gap Analysis — Spec Index

基于 Python Trellis v0.6.3 与 trellis-go (dev) 的完整对比分析。

## Spec 清单

| # | Spec | 分类 | 优先级 | 估算 |
|---|------|------|:---:|:---:|
| 1 | [cli-core](./cli-core.md) | CLI 命令完整性 | 🔴 P0 | 1 周 |
| 2 | [platform-configurators](./platform-configurators.md) | 平台配置器 | 🔴 P0 | 2 周 |
| 3 | [session-journal](./session-journal.md) | 会话记录系统 | 🔴 P0 | 0.5 周 |
| 4 | [workflow-engine](./workflow-engine.md) | 工作流引擎 | 🟡 P1 | 1 周 |
| 5 | [sub-agents](./sub-agents.md) | 子代理系统 | 🟡 P1 | 0.5 周 |
| 6 | [auto-skills](./auto-skills.md) | 自动触发技能 | 🟡 P1 | 0.5 周 |
| 7 | [slash-commands](./slash-commands.md) | 用户命令 | 🟡 P1 | 0.3 周 |
| 8 | [context-system](./context-system.md) | 上下文系统 | 🟡 P1 | 0.5 周 |
| 9 | [task-lifecycle](./task-lifecycle.md) | 任务生命周期 | 🟢 P2 | 0.5 周 |
| 10 | [ecosystem](./ecosystem.md) | 生态与平台扩展 | 🔵 P3 | 1.5 周 |

## 总估算

| 优先级 | Spec 数 | 工作量 |
|:---|:---:|:---|
| 🔴 P0 | 3 | ~3.5 周 |
| 🟡 P1 | 5 | ~2.8 周 |
| 🟢 P2 | 1 | ~0.5 周 |
| 🔵 P3 | 1 | ~1.5 周 |
| **合计** | **10** | **~8.3 周** |

## 当前状态

- trellis-go: dev 版本，~5,947 行 Go 代码
- Python Trellis: v0.6.3，~50,000+ 行（含模板/配置）
- 已实现：基础 CLI 框架、15 平台注册、任务 CRUD、上下文构建器、spec 加载器、模板引擎、hook 生成器、workflow 解析器骨架
- 差距：约 85% 功能待实现

# graphql-sqlc-go-demo

本项目最初目标是构建一个结构清晰、可长期维护的 Blog 系统，并在设计层面预留向 **CMS / 内容平台** 演进的能力。

设计重点不在“快速 CRUD”，而在：

* 业务规则内聚
* 权限与流程清晰
* 架构可演进而不推翻

> 技术栈：**Golang + GraphQL + sqlc**
> 架构风格：**UseCase 驱动 + 读写分离 + SQL First**
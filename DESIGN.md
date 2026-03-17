# Blog / CMS 系统设计说明书（DESIGN.md）

> 技术栈：**Golang + GraphQL + sqlc**
> 架构风格：**UseCase 驱动 + 读写分离 + SQL First**

---

## 1. 项目背景与目标

本项目最初目标是构建一个结构清晰、可长期维护的 Blog 系统，并在设计层面预留向 **CMS / 内容平台** 演进的能力。

设计重点不在“快速 CRUD”，而在：

* 业务规则内聚
* 权限与流程清晰
* 架构可演进而不推翻

**注意事项 / 防踩坑提示**：

* 禁止绕过 UseCase 直接访问 Repository
* 禁止在 UseCase 内调用其他 UseCase 执行业务逻辑

---

## 2. 核心设计原则

1. **UseCase 是唯一业务入口**
2. **GraphQL 只是协议适配层**
3. **SQL First，数据库是事实来源**
4. **读写分离（Command / Query）**
5. **权限与流程在 UseCase 内裁决**
6. **前端不做业务判断**

**防踩坑提示**：

* 严格禁止直接在 GraphQL 或 Repository 中修改状态
* 所有业务逻辑必须集中在 UseCase

---

## 3. 总体架构

### 3.1 分层结构

* API 层：GraphQL Resolver
* UseCase 层：业务规则与流程
* Repository 层：sqlc 生成的数据访问
* Database：PostgreSQL

### 3.2 调用约束

* Resolver → UseCase → Repository → DB
* 禁止反向调用或绕过层级访问

**防踩坑提示**：

* Resolver 仅做 Adapter，不做业务逻辑
* UseCase 不调用 UseCase

---

## 4. UseCase 设计

### 4.1 UseCase 定义

UseCase 表示“用户在系统中能够完成的一次明确业务行为”，例如：

* 创建文章
* 发布文章
* 提交审核

**防踩坑提示**：

* UseCase 不等于 Service，不共享状态
* 禁止承载多个业务动作或互相调用

---

### 4.2 Command / Query 分离

#### Command UseCase

* 改变系统状态
* 需要权限校验
* 可能触发 Workflow 与审计
* **只返回单个实体或执行结果，不返回列表**

#### Query UseCase

* 只读
* 面向前端视图
* 可聚合、裁剪、缓存
* **不修改任何状态字段**

---

## 5. GraphQL 设计原则

### 5.1 UseCase 驱动 Schema

* 每个 Mutation 对应一个 Command UseCase
* 每个 Query 对应一个 Query UseCase
* Schema 中不允许存在“无 UseCase 对应”的能力

**防踩坑提示**：

* 禁止 Resolver 内拼接或组合其他 UseCase 数据
* 所有字段新增必须有明确 UseCase 对应
* Error / canXXX 权限字段统一规范

### 5.2 读模型（ViewModel）

GraphQL 返回的是**读模型**，而不是数据库结构。

---

## 6. 权限模型

### 6.1 权限立场

* 权限是 **UseCase 级别** 的
* GraphQL 不参与权限判断
* **Fail Closed 原则：权限不明确即拒绝**

### 6.2 模型选型

* RBAC：角色（User / Author / Editor / Admin）
* ABAC：资源属性（作者、状态、归属）

**防踩坑提示**：

* 权限变动触发审计日志
* UseCase 必须测试成功与失败路径

---

## 7. 内容 Workflow（状态机）

### 7.1 内容生命周期

* Draft
* PendingReview
* Published
* Unpublished

### 7.2 状态迁移原则

* 所有状态迁移必须通过 Command UseCase
* 禁止通过 UpdateContent 或直接修改状态字段绕过 Workflow
* 建议版本化状态机，支持历史追踪

---

## 8. 数据访问层（sqlc）

### 8.1 sqlc 定位

* 只做 SQL → Go 映射
* 不包含业务逻辑
* 不处理权限

### 8.2 SQL 设计原则

* Query SQL 可冗余，Command SQL 尽量最少
* 事务由 UseCase 控制

---

## 9. 搜索 / 标签 / 分类设计

**防踩坑提示**：

* 搜索默认只针对 Published 内容
* 搜索能力封装，便于替换底层引擎（如 Elasticsearch）
* 标签 / 分类深度查询限制分页或缓存

---

## 10. 审计日志设计

* 每个 Command UseCase 都必须产生审计记录
* Query UseCase 不产生审计
* 审计写入失败不阻塞业务，触发告警
* 高频操作可异步或批量归档

---

## 11. 非功能性设计

* 可测试性：以 UseCase 为单元
* 可扩展性：按 UseCase 拆分
* 可维护性：规则集中、边界清晰

---

## 12. 演进路线

1. 单作者 Blog
2. 多作者 Blog
3. CMS（审核 / 管理后台）
4. 内容平台（多类型内容）

---

## 13. 结语

本设计目标不是“最少代码”，而是：

> **最少推翻重写**

UseCase 驱动、读写分离、权限内聚，使系统可以在长期演进中保持结构稳定。

**防踩坑总结**：

* UseCase 禁止套 UseCase
* Command UseCase 仅修改状态 / 返回结果，不返回列表
* GraphQL Schema 必须映射 UseCase，禁止拼接 Resolver
* Workflow 状态迁移必须走 UseCase，禁止外部指定状态
* 权限 Fail Closed，测试覆盖成功与失败
* 审计日志异步化，核心业务优先
* 搜索 / 标签 / 分类能力封装，便于替换和限制查询深度

---

## 14. 开发任务拆解（Task List）

> 本章节用于将设计文档直接映射为**可执行开发任务**，适合拆成 Issue / TODO / Milestone。

### 阶段 A：项目基础

* [ ] 初始化 Go Module
* [ ] 定义项目目录结构（cmd / internal / usecase / repository 等）
* [ ] 引入基础依赖（GraphQL / sqlc / 配置 / 日志）

### 阶段 B：数据库与 Schema

* [ ] 设计 users 表
* [ ] 设计 posts（content）表
* [ ] 设计 comments 表
* [ ] 设计 tags / categories 表
* [ ] 设计 audit_logs 表
* [ ] 定义内容 workflow 状态枚举

### 阶段 C：SQL 与 sqlc

* [ ] 编写基础 CRUD SQL
* [ ] 编写内容状态流转相关 SQL
* [ ] 编写列表 / 详情查询 SQL
* [ ] 配置并生成 sqlc 代码

### 阶段 D：UseCase（Command）

* [ ] 定义统一业务错误模型
* [ ] CreateContent UseCase
* [ ] UpdateContent UseCase
* [ ] SubmitForReview UseCase
* [ ] ApproveContent / RejectContent UseCase
* [ ] Publish / Unpublish Content UseCase
* [ ] CreateComment UseCase

### 阶段 E：UseCase（Query）

* [ ] ListContent Query UseCase
* [ ] GetContentDetail Query UseCase
* [ ] AdminContentList Query UseCase
* [ ] SearchContent Query UseCase

### 阶段 F：权限与 Workflow

* [ ] 定义角色模型（User / Author / Editor / Admin）
* [ ] 实现 RBAC + ABAC 权限判断（Fail Closed）
* [ ] 实现内容状态机校验

### 阶段 G：审计日志

* [ ] 定义审计事件类型
* [ ] 在所有 Command UseCase 中记录审计日志（异步/告警）
* [ ] 提供审计日志查询能力 / 批量归档

### 阶段 H：GraphQL 接入

* [ ] 定义 GraphQL Schema（Query / Mutation）
* [ ] Resolver → UseCase 适配
* [ ] Context 注入用户与角色信息

### 阶段 I：稳定化

* [ ] 为 UseCase 编写单元测试（覆盖权限/状态/成功失败）
* [ ] 统一错误码与错误语义
* [ ] 补充 README / 使用说明


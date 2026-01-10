## 用户管理系统（Go-zero + GORM + PostgreSQL）

本项目实现了一个具备完整 RBAC 权限控制的用户管理后端，使用 Go-zero 作为 REST 框架，结合 GORM 操作 PostgreSQL。系统覆盖注册/登录、JWT 认证、个人资料维护、密码修改、用户分页查询、启停状态管理以及角色分配等核心能力，并在日志、错误处理与数据校验方面提供统一封装，便于进一步扩展。

### 功能特性
- **注册与登录**：输入校验、唯一约束检测、密码 Bcrypt 加密存储，登录成功后返回短期 Access Token 与可选 Refresh Token。
- **JWT 认证**：`Authorization: Bearer <token>` 头部经过中间件校验，自动把用户 Claims 注入请求上下文供业务使用。
- **个人中心**：支持查询当前用户资料、更新邮箱/姓名以及修改密码（需校验旧密码一致性）。
- **RBAC 权限控制**：基于角色的守卫中间件，仅允许 `admin` 角色访问后台接口；用户-角色、角色-权限均采用多对多表设计。
- **后台运营能力**：
  - 用户分页查询（关键字、状态过滤 + 创建时间倒序）。
  - 用户状态切换（启用/禁用）。
  - 为指定用户重新分配角色，自动在事务内重建关联。
- **安全与合规**：全链路参数校验、统一错误码、详细日志、SQL 占位符防注入、敏感信息加密保存。

### 技术栈
- **语言**：Go 1.24+
- **框架**：Go-zero（REST Server / 中间件 / httpx）
- **ORM**：GORM + `gorm.io/driver/postgres`
- **数据库**：PostgreSQL（示例数据库 `user_mgmt`）
- **验证**：`github.com/go-playground/validator/v10`
- **鉴权**：`github.com/golang-jwt/jwt/v5`

### 目录结构
- `cmd/api/user.go`：服务入口，加载配置、初始化上下文、注册路由并启动 HTTP Server。
- `etc/user-api.yaml`：运行时配置（端口、数据库、JWT、分页、CORS 等）。
- `internal/config`：配置结构体定义。
- `internal/svc`：`ServiceContext`，集中初始化 GORM、Validator、JWT/角色中间件，提供 `AutoMigrate`。
- `internal/model`：用户、角色、权限及关联表模型。
- `internal/handler`：按领域划分的 HTTP Handler（Auth、User Self-Service、Admin）。
- `internal/logic`：业务逻辑层，含公共 DTO 映射、用户与管理员相关逻辑、错误抽象。
- `internal/middleware`：JWT 鉴权与角色守卫中间件。
- `db/migrations`：手写 SQL，用于初始化 PostgreSQL 架构与索引。
- `pkg/*`：通用能力（JWT/密码工具、HTTP 响应包装、上下文 Claims 注入）。

### 快速开始
1. **准备环境**
   - 安装 Go 1.24+ 与 PostgreSQL 14+。
   - 设置 `GOPATH` 并启用 Go modules。
2. **克隆代码并安装依赖**
   ```bash
   git clone <repo-url>
   cd usermgmt
   go mod tidy
   ```
3. **配置数据库**
   - 创建数据库：`createdb user_mgmt`。
   - 修改 `etc/user-api.yaml` 中的 `Database.DSN`、`JWT.AccessSecret`、CORS 白名单等敏感项。
   - 可直接运行 `db/migrations/001_init.sql`，或依赖程序启动时的 `AutoMigrate()` 自动建表（推荐先执行 SQL 以确保 ENUM/索引被创建）。
4. **运行服务**
   ```bash
   go run cmd/api/user.go -f etc/user-api.yaml
   ```
   默认监听 `http://0.0.0.0:8888`。

### API 概览
| 模块 | 方法 & 路径 | 描述 | 认证 | 备注 |
| --- | --- | --- | --- | --- |
| Auth | `POST /api/v1/auth/register` | 用户注册 | 否 | 返回基本 `UserDTO`。
| Auth | `POST /api/v1/auth/login` | 用户登录 | 否 | 返回 Access/Refresh Token + 用户信息。
| Profile | `GET /api/v1/me` | 获取当前用户资料 | 是 | 需携带 JWT。
| Profile | `PUT /api/v1/me` | 更新邮箱/姓名 | 是 | 通过 validator 做格式校验。
| Profile | `POST /api/v1/me/password` | 修改密码 | 是 | 校验旧密码后写入 Bcrypt。
| Admin | `GET /api/v1/admin/users` | 分页查询用户 | 是（Admin） | 支持 `keyword`、`status`、`page`、`pageSize`。
| Admin | `PATCH /api/v1/admin/users/:id/status` | 修改用户启用/禁用状态 | 是（Admin） | 请求体 `{"status":"enabled"|"disabled"}`。
| Admin | `POST /api/v1/admin/users/:id/roles` | 重新分配用户角色 | 是（Admin） | 需传入 `roles` 字符串数组。

> **提示**：所有受保护接口都需要 `Authorization: Bearer <access-token>`，而管理员接口还需当前用户 Claims 中包含 `admin` 角色。

### 数据库与 RBAC
- `users`：记录基础资料、状态、最后登录时间，状态枚举 `enabled/disabled`。
- `roles` / `permissions`：角色与权限元数据表。
- `user_roles`、`role_permissions`：多对多关联表，均配置了外键级联删除。
- 初始化角色 & 超级管理员账户可通过执行 SQL，例如：
  ```sql
  INSERT INTO roles (name, description) VALUES ('admin', 'Platform administrator');
  INSERT INTO users (username, email, password_hash, full_name, status)
  VALUES ('admin', 'admin@example.com', '<bcrypt-hash>', 'Administrator', 'enabled');
  INSERT INTO user_roles (user_id, role_id) VALUES (<admin_user_id>, <admin_role_id>);
  ```

### 安全实践
- **密钥管理**：`JWT.AccessSecret` 必须使用足够复杂的随机字符串，并可通过环境变量注入后写入配置文件。
- **HTTPS / 反向代理**：生产环境建议置于 Nginx、Envoy 等 HTTPS 入口之后。
- **密码策略**：默认最小 8 位，可在 `types` 的 validator 标签中继续增强复杂度要求。
- **审计**：`last_login_at`、`status` 字段可用于风控，亦可扩展加入登录 IP、设备指纹等。

### 开发与测试
- **代码风格**：使用 `gofmt`（已在项目中运行）。
- **自动迁移**：`ServiceContext.AutoMigrate()` 在每次启动时执行，适合开发环境；生产建议使用版本化迁移工具。
- **测试**：当前仓库尚未包含单元测试骨架，可直接运行 `go test ./...` 进行编译级校验，并在后续补充 mock/集成测试。

### 常见问题
- **JWT 失效**：确认 Access Token 与 Refresh Token 的过期时间是否符合需求，必要时刷新并更新客户端缓存。
- **跨域**：默认放开全部 Origin，可在 `Security.AllowOrigins` 中列出受信域名。
- **并发修改角色**：角色分配通过数据库事务保证数据一致，若需审计可在 `user_roles` 表扩展操作人字段。

欢迎在此基础上继续拓展（例如操作日志、权限细粒度控制、OpenAPI 文档等），以满足更复杂的业务场景。
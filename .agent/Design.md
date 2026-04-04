# Design.md - Wenxi Cloud Disk

## 架构设计

### 技术栈 (已确定)

| 组件 | 技术 | 版本/备注 |
|------|------|----------|
| 后端框架 | Go + Gin | 1.21 |
| 数据库 | SQLite | modernc.org/sqlite (纯Go,无CGO) |
| 认证 | JWT | HS256, golang-jwt/jwt/v5 |
| 密码加密 | bcrypt | golang.org/x/crypto |
| 前端框架 | Vue 3 + Vite | ^3.4 |
| UI 库 | Naive UI | ^2.x |
| CSS | Tailwind CSS | 4.x (@tailwindcss/vite) |
| 状态管理 | Pinia | ^2.x |
| HTTP客户端 | Axios | ^1.x |
| 图标 | @vicons/ionicons5 | Ion Icons 5 |

### 项目结构 (当前)

```
wenxi-cloud/
├── backend/                          # Go + Gin 后端
│   ├── cmd/server/main.go           # 入口 (端口8080)
│   ├── config/config.go             # 配置加载
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go           # Gin路由注册
│   │   │   └── handlers/            # HTTP handlers
│   │   │       ├── auth.go         # 认证 (register/login/me)
│   │   │       ├── file.go         # 文件CRUD/上传
│   │   │       ├── folder.go       # 文件夹CRUD
│   │   │       ├── share.go        # 分享
│   │   │       └── lanzou.go       # 蓝奏云token管理
│   │   │   └── middleware/          # Gin中间件
│   │   │       ├── auth.go         # JWT验证
│   │   │       ├── cors.go         # 跨域
│   │   │       └── logger.go       # 请求日志
│   │   ├── model/                   # 数据模型 (GORM)
│   │   │   ├── user.go
│   │   │   ├── file.go
│   │   │   ├── folder.go
│   │   │   ├── share.go
│   │   │   ├── lanzou_token.go
│   │   │   └── upload_session.go
│   │   ├── repository/              # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── file.go
│   │   │   ├── folder.go
│   │   │   ├── share.go
│   │   │   └── lanzou.go
│   │   ├── service/                 # 业务逻辑层
│   │   │   ├── auth.go
│   │   │   ├── file.go
│   │   │   ├── folder.go
│   │   │   ├── share.go
│   │   │   └── lanzou.go
│   │   └── pkg/                     # 工具包
│   │       ├── crypto/password.go   # bcrypt密码
│   │       ├── database/sqlite.go   # 纯Go SQLite
│   │       ├── jwt/jwt.go          # JWT管理
│   │       ├── middleware/          # Gin中间件
│   │       └── response/           # 统一响应
│   ├── go.mod / go.sum
│   └── data/wenxi.db               # SQLite数据库
├── frontend/                         # Vue 3 前端
│   ├── src/
│   │   ├── api/index.js            # axios封装+拦截器
│   │   ├── router/index.js          # Vue Router+守卫
│   │   ├── stores/
│   │   │   ├── auth.js             # Pinia auth store
│   │   │   └── file.js            # Pinia file store
│   │   ├── views/
│   │   │   ├── Login.vue           # 登录页
│   │   │   ├── Register.vue        # 注册页
│   │   │   └── Dashboard.vue       # 主面板
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css               # Tailwind入口
│   ├── vite.config.js              # Vite配置 (代理/api到:8080)
│   └── dist/                        # 构建输出
├── .agent/                          # 开发文档
│   ├── Requirement.md              # 需求规范
│   ├── Design.md                    # 本文件
│   └── Task.md                     # 任务列表
└── AGENTS.md                       # 项目记忆
```

### 分层架构

```
┌─────────────────────────────────┐
│         HTTP Layer              │  Gin Handlers (handlers/)
├─────────────────────────────────┤
│       Middleware Layer          │  Auth, CORS, Logger
├─────────────────────────────────┤
│       Service Layer             │  Business Logic (service/)
├─────────────────────────────────┤
│     Repository Layer            │  Data Access (repository/)
├─────────────────────────────────┤
│        Data Layer               │  GORM + SQLite (model/ + database/)
└─────────────────────────────────┘
```

### 前端架构

```
┌─────────────────────────────────┐
│        Views (页面)              │  Login, Register, Dashboard
├─────────────────────────────────┤
│        Stores (状态)             │  Pinia (auth, file)
├─────────────────────────────────┤
│     Router (路由)                │  Vue Router + 导航守卫
├─────────────────────────────────┤
│       API Layer                 │  Axios (请求拦截, token)
├─────────────────────────────────┤
│      UI Components              │  Naive UI + Tailwind CSS
└─────────────────────────────────┘
```

---

## API 设计 (已实现)

### 认证模块 `/api/auth`

| 方法 | 端点 | 请求体 | 响应 |
|------|------|--------|------|
| POST | /register | `{email, password}` | `{user, token}` |
| POST | /login | `{email, password}` | `{user, token}` |
| GET | /me | - | `{user}` |

### 文件模块 `/api/files`

| 方法 | 端点 | 请求体 | 响应 |
|------|------|--------|------|
| GET | / | `?folder_id=` | `[{file}]` |
| POST | /upload | multipart | `{file}` |
| POST | /upload-url | `{filename, size}` | `{upload_url, file_id}` |
| GET | /:id | - | `{file}` |
| PUT | /:id | `{name}` | `{file}` |
| DELETE | /:id | - | `{message}` |
| GET | /:id/versions | - | `[{version}]` |
| POST | /:id/versions/:vid/restore | - | `{file}` |

### 文件夹模块 `/api/folders`

| 方法 | 端点 | 请求体 | 响应 |
|------|------|--------|------|
| GET | / | `?parent_id=` | `[{folder}]` |
| POST | / | `{name, parent_id?}` | `{folder}` |
| PUT | /:id | `{name}` | `{folder}` |
| DELETE | /:id | - | `{message}` |
| PUT | /:id/move | `{parent_id}` | `{folder}` |

### 分享模块 `/api/shares`

| 方法 | 端点 | 请求体 | 响应 |
|------|------|--------|------|
| GET | / | - | `[{share}]` |
| POST | / | `{file_id, password?, expires_at?}` | `{share}` |
| GET | /:token | - | `{share}` (公开) |
| DELETE | /:id | - | `{message}` |

---

## 数据库设计 (已实现)

### users
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| email | string(100) | UNIQUE, NOT NULL |
| password_hash | string(100) | NOT NULL |
| created_at | datetime | |
| updated_at | datetime | |

### folders
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| user_id | uint | FK(users), NOT NULL |
| parent_id | uint | FK(folders), NULL |
| name | string(255) | NOT NULL |
| created_at | datetime | |
| updated_at | datetime | |

### files
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| user_id | uint | FK(users), NOT NULL |
| folder_id | uint | FK(folders), NULL |
| name | string(255) | NOT NULL |
| size | int64 | NOT NULL |
| lanzou_file_id | string(100) | |
| encryption_key | string(64) | |
| encryption_nonce | string(24) | |
| mime_type | string(100) | |
| created_at | datetime | |
| updated_at | datetime | |

### shares
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| user_id | uint | FK(users), NOT NULL |
| file_id | uint | FK(files), NOT NULL |
| share_token | string(64) | UNIQUE, NOT NULL |
| password_hash | string(100) | NULL |
| expires_at | datetime | NULL |
| created_at | datetime | |

### lanzou_tokens
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| user_id | uint | FK(users), NOT NULL |
| cookie | text | |
| token_value | string(100) | |
| expires_at | datetime | |
| updated_at | datetime | |

### upload_sessions
| 字段 | 类型 | 约束 |
|------|------|------|
| id | uint | PK, AUTO |
| user_id | uint | FK(users), NOT NULL |
| file_name | string(255) | NOT NULL |
| file_size | int64 | NOT NULL |
| file_hash | string(64) | |
| chunks_total | int | NOT NULL |
| chunks_uploaded | int | DEFAULT 0 |
| status | string(20) | NOT NULL |
| created_at | datetime | |
| updated_at | datetime | |

---

## 关键设计决策

### 1. 纯Go SQLite (无CGO)

**问题**: `mattn/go-sqlite3` 需要CGO，Windows无GCC

**方案**: 使用 `modernc.org/sqlite` 纯Go实现

**影响**:
- 编译简单，无外部依赖
- 性能略低于C版本，但足够中小型应用
- 跨平台编译简单

### 2. 客户端直连蓝奏云

**原因**: 节省服务器带宽成本

**架构**:
```
Client ←metadata→ Server ←token→ LanZouCloud
                (带宽节省)  (直传文件)
```

**上传流程**:
1. Client → Server: 初始化上传会话
2. Server → Client: 返回 upload_url
3. Client → LanZouCloud: 直传文件
4. Client → Server: 保存元数据

**下载流程**:
1. Client → Server: 请求下载URL
2. Server → Client: 返回 signed_url
3. Client → LanZouCloud: 直链下载

### 3. 零圆角 UI 设计

**规范**:
- 全局 border-radius: 0px
- 背景: #0f0f0f (深黑)
- 卡片: #1a1a1a (浅黑)
- 强调色: #00D4FF (霓虹青)

**实现**:
- Tailwind CSS 4 无配置文件方式
- Naive UI 全局样式覆盖

### 4. JWT 认证

**配置**:
- 算法: HS256
- 过期: 24小时
- 存储: localStorage

**流程**:
1. 登录 → Server返回token → Client存localStorage
2. 请求 → Client在header带token
3. Server验证token → 允许/拒绝

---

## 测试策略

### 后端单元测试

| 包 | 测试内容 | 框架 |
|----|---------|------|
| pkg/crypto | 密码哈希/验证 | testify |
| pkg/jwt | Token生成/验证 | testify |
| pkg/response | 响应结构 | testify |

### 前端测试 (待实现)

| 类型 | 工具 |
|------|------|
| 组件测试 | Vitest + @vue/test-utils |
| Store测试 | Vitest |
| E2E测试 | Playwright |

---

**文档版本**: 4.0.0
**最后更新**: 2026-04-04
**状态**: 活跃开发中

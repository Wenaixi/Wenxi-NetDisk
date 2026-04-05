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
│   │   ├── api/lanzou.js           # 蓝奏云API封装
│   │   ├── api/recycle.js          # 回收站API封装
│   │   ├── router/index.js          # Vue Router+守卫
│   │   ├── stores/
│   │   │   ├── auth.js            # Pinia auth store
│   │   │   ├── auth.test.js       # auth store测试 (7)
│   │   │   ├── file.js            # Pinia file store
│   │   │   ├── file.test.js       # file store测试 (13)
│   │   │   ├── upload.js           # Pinia upload store
│   │   │   ├── upload.test.js     # upload store测试 (5)
│   │   │   ├── share.js           # Pinia share store
│   │   │   ├── share.test.js     # share store测试 (5)
│   │   │   ├── download.js       # Pinia download store
│   │   │   ├── download.test.js  # download store测试 (4)
│   │   │   ├── shareParse.js     # Pinia shareParse store
│   │   │   ├── shareParse.test.js # shareParse store测试 (16)
│   │   │   ├── recycle.js        # Pinia recycle store
│   │   │   ├── recycle.test.js   # recycle store测试 (15)
│   │   │   ├── sync.js           # Pinia sync store
│   │   │   ├── sync.test.js      # sync store测试 (11)
│   │   │   └── fileDescription.test.js # fileDescription store测试 (4)
│   │   ├── components/
│   │   │   ├── AppHeader.vue      # 通用页面头部组件
│   │   │   └── FileDetailModal.vue # 文件详情弹窗
│   │   ├── utils/
│   │   │   ├── crypto.js         # 客户端加密(AES-GCM)
│   │   │   ├── crypto.test.js   # crypto工具测试 (4)
│   │   │   ├── fileSplit.js     # 文件分割/合并工具
│   │   │   └── fileSplit.test.js # fileSplit工具测试 (22)
│   │   ├── views/
│   │   │   ├── Login.vue         # 登录页
│   │   │   ├── Register.vue      # 注册页
│   │   │   ├── Dashboard.vue     # 主面板(本地文件)
│   │   │   ├── LanzouSettings.vue # 蓝奏云设置页
│   │   │   ├── LanZouBrowser.vue  # 蓝奏云浏览器
│   │   │   ├── RecycleBin.vue     # 回收站
│   │   │   ├── ShareParse.vue     # 链接解析
│   │   │   └── Sync.vue           # 同步资源
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css             # Tailwind入口
│   ├── vite.config.js            # Vite配置 (代理/api到:8080)
│   ├── vitest.config.js          # Vitest测试配置
│   └── dist/                     # 构建输出
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

### 5. Service层接口化

**问题**: 单元测试困难，难以Mock依赖

**方案**: 所有Service层使用接口而非具体类型依赖

**实现的接口**:
```
UserRepository            - 用户仓库 (auth_service)
FileRepository            - 文件仓库 (file_service)
ShareRepository           - 分享仓库 (share_service)
ShareFileRepository       - 分享用文件仓库 (share_service)
LanZouTokenRepository     - 蓝奏云Token仓库 (lanzou_service)
UploadSessionRepository   - 上传会话仓库 (upload_service)
LanZouClientProvider      - 蓝奏云客户端提供者 (upload_service)
FileMetadataCreator       - 文件元数据创建器 (upload_service)
FileVersionRepository     - 文件版本仓库 (file_version_service)
VersionFileRepository     - 版本用文件仓库 (file_version_service)
RecycleBinRepository      - 回收站仓库 (recycle_service)
RecycleFileRepo           - 回收站用文件仓库 (recycle_service)
RecycleFolderRepo         - 回收站用文件夹仓库 (recycle_service)
```

**影响**:
- 100%纯Go测试，无数据库依赖
- 所有Service层可独立测试
- 便于未来替换实现

---

## 测试策略

### 后端单元测试

| 包 | 测试内容 | 测试数 | 框架 |
|----|---------|--------|------|
| pkg/crypto | 密码哈希/验证 | 5 | go test |
| pkg/jwt | Token生成/验证 | 7 | go test |
| pkg/response | 响应结构 | 5 | go test |
| pkg/lanzou | 蓝奏云API解析 | 6 | go test |
| repository | 数据访问层 | 10 | go test + gorm+sqlite |
| service/auth | 注册/登录/验证 | 8 | go test |
| service/file | 文件CRUD | 12 | go test |
| service/folder | 文件夹CRUD/移动 | 22 | go test |
| service/lanzou | 蓝奏云连接 | 9 | go test |
| service/upload | 断点上传 | 6 | go test |
| service/share | 分享管理 | 21 | go test |
| service/download | 下载直链 | 5 | go test |
| service/recycle | 回收站操作 | 11 | go test |
| service/file_version | 版本管理 | 18 | go test |

| pkg/middleware | 中间件 (Auth) | 10 | go test + httptest |
| handlers | Handler验证层 | 69 | go test + httptest |

### 前端测试

| 类型 | 模块 | 测试数 | 状态 |
|------|------|--------|------|
| Store | auth | 7 | ✅ |
| Store | file | 13 | ✅ |
| Store | upload | 7 | ✅ |
| Store | share | 10 | ✅ |
| Store | download | 10 | ✅ |
| Store | recycle | 15 | ✅ |
| Store | sync | 11 | ✅ |
| Store | fileDescription | 4 | ✅ |
| Store | shareParse | 16 | ✅ |
| Utils | crypto | 4 | ✅ |
| Utils | fileSplit | 22 | ✅ |
| Component | FileDetailModal | 17 | ✅ |
| Component | AppHeader | 8 | ✅ |
| Component | Login | 8 | ✅ |
| Component | Register | 10 | ✅ |
| Component | LanzouSettings | 12 | ✅ |
| Component | RecycleBin | 16 | ✅ |
| Component | Sync | 16 | ✅ |
| Component | ShareParse | 15 | ✅ |
| Component | LanZouBrowser | 26 | ✅ |
| View | Dashboard | 27 | ✅ |
| API | auth/index | 3 | ✅ |
| API | file/index | 12 | ✅ |
| API | folder/index | 7 | ✅ |
| API | share/index | 5 | ✅ |
| API | lanzou | 12 | ✅ |
| API | recycle | 5 | ✅ |
| API | interceptors | 15 | ✅ |
| Router | guards | 17 | ✅ |
| E2E | Playwright | 15 | ⏭️ |

---

**文档版本**: 7.0.0
**最后更新**: 2026-04-05
**状态**: Phase 6 完成, 测试 923 全部通过 (~567 Go + 356 Vue), Handler覆盖率 96.4%

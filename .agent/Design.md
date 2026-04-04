# Design.md - Wenxi Cloud Disk

## 架构设计

### 技术栈

| 组件 | 技术 |
|------|------|
| 后端框架 | Go + Gin |
| ORM | GORM |
| 数据库 | SQLite |
| 认证 | JWT (HS256) |
| 前端框架 | Vue 3 + Vite |
| UI 库 | Naive UI |
| CSS | Tailwind CSS |

### 项目结构

```
backend/
├── cmd/server/main.go           # 入口
├── config/config.go            # 配置
├── internal/
│   ├── api/
│   │   ├── router.go         # 路由
│   │   ├── handlers/          # HTTP handlers
│   │   └── middleware/       # 中间件
│   ├── model/                 # GORM 模型
│   ├── repository/           # 数据访问层
│   ├── service/             # 业务逻辑层
│   └── pkg/                 # 工具包
│       ├── crypto/          # 加密工具
│       ├── database/        # 数据库
│       ├── jwt/             # JWT
│       ├── middleware/      # 中间件
│       └── response/        # 响应
└── go.mod

frontend/
├── src/
│   ├── api/                # API 调用
│   ├── components/         # 组件
│   ├── composables/        # 组合函数
│   ├── stores/             # Pinia 状态
│   ├── views/              # 页面
│   └── styles/             # 样式
```

### 分层架构

```
┌─────────────────────────────────┐
│         HTTP Layer              │  Gin Handlers
├─────────────────────────────────┤
│       Middleware Layer          │  Auth, CORS, Logger
├─────────────────────────────────┤
│       Service Layer             │  Business Logic
├─────────────────────────────────┤
│     Repository Layer            │  Data Access
├─────────────────────────────────┤
│        Data Layer               │  GORM + SQLite
└─────────────────────────────────┘
```

### API 设计

#### 认证流程
```
POST /api/auth/register → {id, username, email}
POST /api/auth/login → {access_token, token_type}
GET /api/auth/me → {id, username, email, created_at}
```

#### 文件流程
```
GET /api/files → [{id, name, size, ...}]
POST /api/files → {id, name, ...}
GET /api/files/:id → {id, name, ...}
DELETE /api/files/:id → {message}
POST /api/files/:id/share → {share_token, share_url}
```

### 数据库设计

#### Users
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| username | string(50) | 用户名，唯一 |
| email | string(100) | 邮箱，唯一 |
| password_hash | string(100) | 密码哈希 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

#### Files
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| user_id | uint | 所属用户 |
| name | string(255) | 文件名 |
| size | int64 | 文件大小 |
| lanzou_file_id | string(100) | 蓝奏云文件ID |
| encryption_key | string(64) | 加密密钥 |
| encryption_nonce | string(24) | 加密随机数 |
| created_at | datetime | 创建时间 |

---

**文档版本**: 3.0.0
**最后更新**: 2026-04-04
**状态**: 新架构设计

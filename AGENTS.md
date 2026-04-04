# AGENTS.md - Wenxi Cloud Disk

## 项目概述

**Wenxi Cloud Disk** - 现代化云盘系统，采用客户端直连蓝奏云存储，服务器仅管理元数据的架构设计。

### 核心架构

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│  Client │◄───────►│  Go Server   │◄───────►│  LanZouCloud│
│ (Browser)│ metadata│  (Gin + GORM)│  token  │  (Storage)  │
└─────────┘   only  └──────────────┘  only   └─────────────┘
                       │
                  ┌────┴────┐
                  │ SQLite  │
                  └─────────┘
```

### 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| **后端** | Go + Gin + GORM | 1.21 |
| **前端** | Vue 3 + Vite + Naive UI | ^3.4 |
| **数据库** | SQLite | - |
| **认证** | JWT | HS256 |
| **加密** | ChaCha20-Poly1305 | 客户端 |
| **部署** | VPS | Linux |

---

## 开发历史

### 2026-04-04 - 新架构启动

**决策**: 重构项目架构，采用新方案

**原因**:
- 原 Python FastAPI 方案：文件经过服务器，带宽成本高
- 新 Go + Gin 方案：客户端直连蓝奏云，服务器只存 metadata

**新架构优势**:
- 带宽成本极低（只传 metadata）
- 文件走蓝奏云 CDN，速度快
- 保留 E2E 加密（服务器看不到文件内容）

**技术选型**:
- 后端: Go + Gin + GORM + SQLite
- 前端: Vue 3 + Vite + Naive UI + Tailwind CSS
- UI风格: 零圆角设计系统（保留原有特色）

### 2026-03-20 - 架构调整

**背景**: 项目从 Web 应用演进到 Tauri 桌面应用（解决蓝奏云 CORS 问题），后又考虑 Vue 3 前端

**结论**: 最终选择保持 Web 架构，但改用 Go 后端直连蓝奏云

---

## 当前进度

### Phase 1: Go 后端搭建 ✅

| 任务 | 状态 | 文件 |
|------|------|------|
| 目录结构 | ✅ | backend/cmd/server, backend/internal/... |
| Go 模块 | ✅ | backend/go.mod |
| 配置模块 | ✅ | backend/config/config.go |
| 数据库模块 | ✅ | backend/internal/pkg/database/sqlite.go |
| 响应模块 | ✅ | backend/internal/pkg/response/response.go |
| 数据模型 | ✅ | backend/internal/model/*.go |
| Repository 层 | ✅ | backend/internal/repository/*.go |
| Service 层 | ✅ | backend/internal/service/*.go |
| JWT 工具 | ✅ | backend/internal/pkg/jwt/jwt.go |
| 密码加密 | ✅ | backend/internal/pkg/crypto/password.go |
| 中间件 | ✅ | backend/internal/pkg/middleware/*.go |
| Handlers | ✅ | backend/internal/api/handlers/*.go |
| 路由 | ✅ | backend/internal/api/router.go |
| 主入口 | ✅ | backend/cmd/server/main.go |

### Phase 2: Vue 前端搭建 ⏳ (待开始)

### Phase 3: 蓝奏云对接 ⏳ (待开始)

### Phase 4: 功能完善 ⏳ (待开始)

---

## 项目结构

```
wenxi-cloud/
├── backend/                          # Go + Gin 后端
│   ├── cmd/server/main.go           # 入口
│   ├── config/config.go             # 配置
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go           # 路由
│   │   │   ├── handlers/           # HTTP handlers
│   │   │   └── middleware/        # 中间件
│   │   ├── model/                  # GORM 模型
│   │   ├── repository/            # 数据访问层
│   │   ├── service/               # 业务逻辑层
│   │   └── pkg/                   # 工具包
│   │       ├── crypto/            # 密码加密
│   │       ├── database/          # 数据库
│   │       ├── jwt/               # JWT
│   │       ├── middleware/        # 中间件
│   │       └── response/          # 响应
│   └── go.mod / go.sum
├── frontend/                         # Vue 3 前端 (待搭建)
│   └── src/
│       ├── api/
│       ├── components/
│       ├── composables/
│       ├── stores/
│       ├── views/
│       └── styles/
├── .agent/                          # 开发文档
│   ├── Requirement.md              # 需求规范
│   ├── Design.md                  # 架构设计
│   └── Task.md                    # 任务列表
└── AGENTS.md                       # 项目记忆
```

---

## 数据库表设计

```sql
-- users: 用户表
users (id, username, email, password_hash, created_at, updated_at)

-- files: 文件元数据表
files (id, user_id, name, size, lanzou_file_id, lanzou_folder_id,
       encryption_key, encryption_nonce, mime_type, created_at, updated_at)

-- shares: 分享表
shares (id, file_id, share_token, password_hash, expires_at, created_at)

-- lanzou_tokens: 蓝奏云 cookie/token
lanzou_tokens (id, user_id, cookie, token_value, expires_at, updated_at)

-- upload_sessions: 断点续传会话
upload_sessions (id, user_id, file_name, file_size, file_hash,
                 chunks_total, chunks_uploaded, status, created_at, updated_at)
```

---

## API 接口设计

### 认证模块 `/api/auth`
| 方法 | 端点 | 说明 |
|------|------|------|
| POST | /register | 用户注册 |
| POST | /login | 用户登录 |
| GET | /me | 获取当前用户 |

### 文件模块 `/api/files`
| 方法 | 端点 | 说明 |
|------|------|------|
| GET | / | 文件列表 |
| POST | / | 创建文件元数据 |
| GET | /:id | 获取文件详情 |
| DELETE | /:id | 删除文件 |
| POST | /:id/share | 生成分享链接 |

### 分享模块 `/api/shares`
| 方法 | 端点 | 说明 |
|------|------|------|
| GET | /:token | 获取分享信息 |
| DELETE | /:id | 删除分享 |

### 蓝奏云模块 `/api/lanzou`
| 方法 | 端点 | 说明 |
|------|------|------|
| POST | /connect | 连接蓝奏云 |
| GET | /status | 连接状态 |

---

## 关键设计决策

### 1. 客户端直连蓝奏云

**原因**: 节省服务器带宽成本

**流程**:
1. Client → Server: 初始化上传会话
2. Server → Client: 返回 upload_url
3. Client → LanZouCloud: 直传加密后的文件块
4. Client → Server: 完成上传，保存元数据

### 2. E2E 加密

**方案**: 客户端 ChaCha20-Poly1305 加密

**密钥管理**:
- 主密钥: 用户密码 + Argon2id 派生
- 文件密钥: 每个文件独立，随机 nonce
- 密钥存储: 加密后存服务器

### 3. 零圆角 UI 设计

**规范**:
- border-radius: 0px
- 颜色: #0A0E17 背景, #00D4FF 强调色
- 阴影: 锐利无模糊 (4px 4px 0px)

---

## 待办事项

- [ ] 安装 Go 依赖 (go mod tidy)
- [ ] 测试后端编译和运行
- [ ] 创建 Vue 3 前端项目
- [ ] 集成蓝奏云 API
- [ ] 实现客户端加密
- [ ] 完善分享功能
- [ ] 添加断点续传
- [ ] 编写单元测试

---

**最后更新**: 2026-04-04
**状态**: Phase 1 完成，Phase 2 待开始

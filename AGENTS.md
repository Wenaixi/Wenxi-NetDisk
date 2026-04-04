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
| **前端** | Vue 3 + Vite + Naive UI + Tailwind CSS 4 | ^3.4 |
| **数据库** | SQLite (modernc.org/sqlite 纯Go) | - |
| **认证** | JWT | HS256 |
| **加密** | ChaCha20-Poly1305 | 客户端 |
| **部署** | VPS | Linux |

---

## 开发历史

### 2026-04-04 - Phase 5 进行中 (~55%)

**UI组件化重构:**
- 创建AppHeader.vue通用组件消除重复代码
- Dashboard、LanZouBrowser、LanzouSettings统一使用AppHeader
- Dashboard布局优化: 文件夹/文件分离区域并添加标签

**分享功能:**
- 文件分享弹窗(密码/有效期设置)
- 分享链接一键复制
- 分享成功后显示密码和有效期
- 后端添加CreateShareViaBody接口支持请求体传file_id

**文件操作增强:**
- 文件/文件夹右键菜单支持
- 文件/文件夹重命名功能
- 菜单根据类型显示不同操作(下载/分享仅对文件)
- Dashboard面包屑导航
- 选择模式(复选框支持多选)
- 批量删除文件/文件夹
- 批量删除逻辑修复(正确区分文件/文件夹)

### 2026-04-04 - Phase 1 & 2 完成

**Go后端完成:**
- 目录结构: backend/cmd/server, backend/internal/...
- 纯Go SQLite驱动 (modernc.org/sqlite - 无CGO依赖)
- JWT认证 + bcrypt密码加密
- RESTful API: auth, files, folders, shares, lanzou
- 所有依赖已安装 (GOSUMDB=off 解决网络问题)

**Vue前端完成:**
- Vite + Vue 3 + Tailwind CSS 4 + Naive UI
- 登录/注册页面
- Dashboard (文件列表、文件夹导航、上传)
- Pinia stores (auth, file)
- Vue Router路由守卫
- axios API封装 + 请求拦截器
- 构建成功 (vite build)

### 2026-03-20 - 架构调整

**背景**: 项目从 Web 应用演进到 Tauri 桌面应用（解决蓝奏云 CORS 问题），后又考虑 Vue 3 前端

**结论**: 最终选择保持 Web 架构，但改用 Go 后端直连蓝奏云

---

## 当前进度

| Phase | 状态 | 进度 |
|-------|------|------|
| Phase 1: Go后端搭建 | ✅ 完成 | 100% |
| Phase 2: Go后端验证 | ✅ 完成 | 100% |
| Phase 3: Vue前端搭建 | ✅ 完成 | 100% |
| Phase 4: 蓝奏云对接 | ✅ 完成 | 100% |
| Phase 5: 功能完善 | 🔄 进行中 | ~55% |
| Phase 6: 测试系统 | ✅ 完成 | 100% |

---

## 项目结构 (当前)

```
wenxi-cloud/
├── backend/                          # Go + Gin 后端
│   ├── cmd/server/main.go           # 入口 (端口8080)
│   ├── config/config.go             # 配置
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go           # 路由
│   │   │   ├── handlers/           # HTTP handlers
│   │   │   │   ├── auth.go
│   │   │   │   ├── file.go
│   │   │   │   ├── folder.go
│   │   │   │   ├── share.go
│   │   │   │   └── lanzou.go
│   │   │   └── middleware/          # 中间件 (Auth, CORS, Logger)
│   │   ├── model/                   # 数据模型
│   │   │   ├── user.go
│   │   │   ├── file.go
│   │   │   ├── folder.go
│   │   │   ├── share.go
│   │   │   ├── lanzou_token.go
│   │   │   └── upload_session.go
│   │   ├── repository/              # 数据访问层
│   │   ├── service/                 # 业务逻辑层
│   │   └── pkg/                     # 工具包
│   │       ├── crypto/password.go   # bcrypt
│   │       ├── database/sqlite.go   # 纯Go SQLite
│   │       ├── jwt/jwt.go          # JWT工具
│   │       ├── lanzou/             # 蓝奏云API封装
│   │       ├── middleware/          # 中间件
│   │       └── response/           # 统一响应
│   ├── go.mod / go.sum
│   └── data/wenxi.db               # SQLite数据库
├── frontend/                         # Vue 3 前端
│   ├── src/
│   │   ├── api/
│   │   │   ├── index.js            # axios封装
│   │   │   └── lanzou.js          # 蓝奏云API
│   │   ├── router/index.js          # Vue Router
│   │   ├── stores/
│   │   │   ├── auth.js             # 认证store
│   │   │   ├── auth.test.js        # auth store测试 (7)
│   │   │   ├── file.js             # 文件store
│   │   │   ├── file.test.js        # file store测试 (13)
│   │   │   └── upload.js           # 上传store
│   │   ├── utils/
│   │   │   ├── crypto.js           # AES-GCM客户端加密
│   │   │   └── crypto.test.js     # crypto工具测试 (4)
│   │   ├── views/
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   ├── Dashboard.vue
│   │   │   ├── LanzouSettings.vue  # 蓝奏云设置
│   │   │   └── LanZouBrowser.vue  # 蓝奏云浏览器
│   │   ├── App.vue
│   │   ├── main.js
│   │   └── style.css               # Tailwind入口
│   ├── vite.config.js              # Vite配置 (代理API到:8080)
│   └── dist/                        # 构建输出
├── .agent/                          # 开发文档
│   ├── Requirement.md              # 需求规范
│   ├── Design.md                    # 架构设计
│   └── Task.md                     # 任务列表
├── references/                       # 参考项目
│   └── lanzouyun-disk/             # 蓝奏云桌面端参考 (Electron+React)
└── AGENTS.md                       # 项目记忆
```

---

## 数据库表设计 (已实现)

```sql
-- users: 用户表
users (id, email, password_hash, created_at, updated_at)

-- folders: 文件夹表
folders (id, user_id, parent_id, name, created_at, updated_at)

-- files: 文件元数据表
files (id, user_id, folder_id, name, size, lanzou_file_id,
       encryption_key, encryption_nonce, mime_type, created_at, updated_at)

-- shares: 分享表
shares (id, user_id, file_id, share_token, password_hash, expires_at, created_at)

-- lanzou_tokens: 蓝奏云认证
lanzou_tokens (id, user_id, cookie, token_value, expires_at, updated_at)

-- upload_sessions: 断点续传
upload_sessions (id, user_id, file_name, file_size, file_hash,
                 chunks_total, chunks_uploaded, status, created_at, updated_at)
```

---

## API 接口设计 (已实现)

### 认证模块 `/api/auth`
| 方法 | 端点 | 说明 | 状态 |
|------|------|------|------|
| POST | /register | 用户注册 | ✅ |
| POST | /login | 用户登录 | ✅ |
| GET | /me | 获取当前用户 | ✅ |

### 文件模块 `/api/files`
| 方法 | 端点 | 说明 | 状态 |
|------|------|------|------|
| GET | / | 文件列表 | ✅ |
| POST | /upload | 上传文件 | ✅ |
| POST | /upload-url | 获取上传URL | ✅ |
| GET | /:id | 获取文件详情 | ✅ |
| PUT | /:id | 更新文件 | ✅ |
| DELETE | /:id | 删除文件 | ✅ |
| GET | /:id/versions | 获取版本 | ✅ |
| POST | /:id/versions/:vid/restore | 恢复版本 | ✅ |

### 文件夹模块 `/api/folders`
| 方法 | 端点 | 说明 | 状态 |
|------|------|------|------|
| GET | / | 文件夹列表 | ✅ |
| POST | / | 创建文件夹 | ✅ |
| PUT | /:id | 更新文件夹 | ✅ |
| DELETE | /:id | 删除文件夹 | ✅ |

### 分享模块 `/api/shares`
| 方法 | 端点 | 说明 | 状态 |
|------|------|------|------|
| GET | / | 分享列表 | ✅ |
| POST | / | 创建分享 | ✅ |
| GET | /:token | 获取分享 | ✅ |
| DELETE | /:id | 删除分享 | ✅ |

### 蓝奏云模块 `/api/lanzou`
| 方法 | 端点 | 说明 | 状态 |
|------|------|------|------|
| POST | /connect | 连接蓝奏云 | ✅ |
| GET | /status | 连接状态 | ✅ |
| DELETE | /connect | 断开连接 | ✅ |
| GET | /files | 文件列表 | ✅ |
| GET | /folders | 文件夹列表 | ✅ |
| POST | /folders | 创建文件夹 | ✅ |
| POST | /share | 创建分享链接 | ✅ |
| GET | /files/:id/url | 获取下载直链 | ✅ |
| POST | /upload/init | 初始化上传 | ✅ |
| POST | /upload/complete/:id | 完成上传 | ✅ |
| GET | /upload/status/:id | 上传状态 | ✅ |

---

## 蓝奏云API设计 (参考 lanzouyun-disk)

### 核心API (基于Electron桌面端逆向分析)

| 功能 | API端点 | 方法 | 参数 |
|------|---------|------|------|
| 列出文件 | `doupload.php` | POST | task=5, folder_id, pg |
| 列出文件夹 | `doupload.php` | POST | task=47, folder_id |
| 创建文件夹 | `doupload.php` | POST | task=2, parent_id, folder_name |
| 文件详情 | `doupload.php` | POST | task=22, file_id |
| 文件夹详情 | `doupload.php` | POST | task=18, folder_id |
| 删除 | `doupload.php` | POST | task=6/46 (文件/文件夹) |
| 重命名 | `doupload.php` | POST | task=14, file_id/folder_id, name |
| 移动 | `doupload.php` | POST | task=15/48, file_id/folder_id, folder_id |
| 回收站 | `doupload.php` | POST | task=46 |
| 分享 | `doupload.php` | POST | task=39, file_id/folder_id |
| 下载(无密码) | 解析HTML iframe | GET | url |
| 下载(有密码) | 解析HTML #passwddiv | POST | ajaxData |

### 参考项目结构

```
references/lanzouyun-disk/
├── src/common/
│   ├── http.ts              # HTTP客户端 (got + cookie jar)
│   ├── cookie.ts            # Cookie管理
│   ├── core/
│   │   ├── ls.ts           # 列出文件/文件夹
│   │   ├── mkdir.ts        # 创建文件夹
│   │   ├── download.ts      # 下载链接解析
│   │   ├── detail.ts       # 文件/文件夹详情
│   │   ├── edit.ts         # 编辑
│   │   ├── mv.ts           # 移动
│   │   ├── rm.ts           # 删除
│   │   ├── rename.ts       # 重命名
│   │   ├── matcher.ts      # HTML解析器
│   │   └── recycle.ts      # 回收站
│   └── util.ts             # 工具函数
```

### 蓝奏云API特点

1. **无官方API** - 基于Web页面逆向分析
2. **Cookie认证** - 登录状态通过Cookie维持
3. **HTML解析** - 使用cheerio解析页面获取数据
4. **分页处理** - 文件列表需循环请求
5. ** Referer验证** - 请求需携带正确Referer头

---

## 关键设计决策

### 1. 纯Go SQLite (无CGO)

**原因**: Windows环境无GCC，macOS需要Xcode，Ubuntu需要build-essential

**方案**: 使用 `modernc.org/sqlite` 替代 `mattn/go-sqlite3`

**影响**: 编译简单，无外部依赖

### 2. 客户端直连蓝奏云

**原因**: 节省服务器带宽成本

**流程**:
1. Client → Server: 初始化上传会话
2. Server → Client: 返回 upload_url
3. Client → LanZouCloud: 直传加密后的文件块
4. Client → Server: 完成上传，保存元数据

### 3. 零圆角 UI 设计

**规范**:
- border-radius: 0px
- 颜色: #0f0f0f 背景, #1a1a1a 卡片
- Tailwind CSS 4 无配置方式

---

## 待办事项

- [x] Go后端搭建
- [x] Go后端验证 (API测试通过)
- [x] Vue前端基础 (登录/注册/仪表盘)
- [x] Vue前端构建
- [x] Go单元测试 - pkg (crypto, jwt, response) 全部通过
- [x] Go repository测试 - 框架完成 (CGO需在Linux/CI环境运行)
- [x] 蓝奏云API封装包 - lanzou (types/client/test) 6个测试通过
- [x] Go service层测试 (lanzou 9 + upload 6)
- [x] Vue单元测试 - Pinia stores 全部通过 (auth 7 + file 13 + upload 5)
- [x] 客户端加密实现 (AES-GCM替代ChaCha20) + 测试 (4)
- [x] 前端upload store + 测试 (5)
- [x] 蓝奏云service扩展 (IsConnected, GetClient)
- [x] 上传service (初始化/分块/断点续传/完成)
- [x] 分享功能增强 (分享列表API/过期检查)
- [x] 下载解密集成 (DownloadService + API)

---

## Git 提交历史

| Commit | 描述 |
|--------|------|
| f51db43 | feat: 实现前端DownloadStore及单元测试 |
| f62df4b | feat: 实现前端ShareStore及单元测试 |
| d468056 | feat: 实现DownloadService下载服务和单元测试 |
| 278c177 | feat: 增强分享功能 - 分享列表API与过期检查 |
| 864cb9e | feat: 实现FolderService完整功能及单元测试 |
| 3b5f0f3 | docs: 更新AGENTS.md和Task.md - FolderService完成 |
| a47ea88 | test: 添加ShareService单元测试并接口化 |
| b97313f | test: 添加AuthService单元测试并接口化 |
| d7b6718 | test: 添加FileService单元测试并接口化 |
| e6c9977 | docs: 全面完善AGENTS.md - 测试统计+架构决策+API文档 |
| 4fb8bf2 | test: 完善前端测试覆盖 - upload store + crypto工具 |
| cb6bbec | test: 添加LanZouService单元测试并修复类型问题 |
| 5ccd040 | test: 添加UploadService单元测试并修复类型问题 |
| ce7fe55 | feat: 完善蓝奏云上传服务和断点续传功能 |
| a9d0282 | feat: 完善前端加密上传功能与蓝奏云类型修复 |
| b152313 | test: 添加Vue前端Pinia Store单元测试 (Vitest) |
| 18f6923 | docs: 更新AGENTS.md - Phase 6测试进度至40% |
| 0601bb4 | test: 添加repository层单元测试框架 (gorm+sqlite) |
| 2d17588 | docs: 全面完善AGENTS.md、Design.md、Requirement.md |
| 02e1447 | chore: 完善.gitignore配置 |
| bb7ff53 | docs: 更新AGENTS.md - 添加测试进度和提交历史 |
| 2ff5176 | test: 添加后端单元测试 (crypto, jwt, response) |
| b7912d4 | docs: 更新任务进度 - Phase 3 Vue前端搭建进行中 |
| a109844 | feat: 完成Vue 3前端基础框架搭建 |
| 611413d | chore: 清理旧代码和无关文件 |
| a951cbf | feat: 全新架构重构 - Go后端 + 蓝奏云直连 |

---

**最后更新**: 2026-04-04
**状态**: Phase 3 (~95%), Phase 4 (~65%), Phase 5 (~40%), Phase 6 测试系统 (~95%)

## 测试统计

### Go后端 (107 tests)
| 模块 | 测试数 | 状态 |
|------|--------|------|
| pkg/crypto | 5 | ✅ |
| pkg/jwt | 7 | ✅ |
| pkg/lanzou | 6 | ✅ |
| pkg/response | 5 | ✅ |
| service/auth | 8 | ✅ |
| service/file | 8 | ✅ |
| service/folder | 18 | ✅ |
| service/lanzou | 9 | ✅ |
| service/upload | 6 | ✅ |
| service/share | 21 | ✅ |
| service/download | 5 | ✅ |
| repository | 10 | ⏭️ (skip CGO) |

### Vue前端 (46 tests)
| 模块 | 测试数 | 状态 |
|------|--------|------|
| auth store | 7 | ✅ |
| file store | 13 | ✅ |
| upload store | 5 | ✅ |
| share store | 10 | ✅ |
| download store | 7 | ✅ |
| crypto utils | 4 | ✅ |

**总测试数: 153 (107 Go + 46 Vue) 全部通过**

## 关键架构决策记录

### 1. 接口化service层
service层使用接口而非具体类型依赖，提高可测试性：
- `FolderRepository` - 文件夹仓库接口
- `FileRepository` - 文件仓库接口
- `ShareRepository` - 分享仓库接口
- `ShareFileRepository` - 分享用文件仓库接口
- `LanZouTokenRepository` - 蓝奏云Token仓库接口
- `UploadSessionRepository` - 上传会话仓库接口
- `LanZouClientProvider` - 蓝奏云客户端提供者接口
- `FileMetadataCreator` - 文件元数据创建器接口
- `UserRepository` - 用户仓库接口

### 2. 加密算法选择
Web Crypto API不支持ChaCha20-Poly1305，使用AES-GCM-256替代，安全性相同级别。

### 3. 前端store导出
测试mock时需要注意：`crypto.js`中的纯函数必须export才能被独立测试。`globalThis.crypto`在Vitest中是只读的，无法直接mock。

# Task.md - Wenxi Cloud Disk 开发任务

## 开发目标

打造现代化的云盘系统，客户端直连蓝奏云，服务器管理元数据。

---

## Phase 1: Go 后端搭建 ✅

### 完成状态

| 任务 | 状态 | 文件 |
|------|------|------|
| 创建目录结构 | ✅ | backend/cmd/server, backend/internal/... |
| Go 模块配置 | ✅ | backend/go.mod |
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

---

## Phase 2: Go 后端验证 ✅

### 完成状态
- [x] go mod tidy - 依赖安装完成
- [x] go build - 编译通过 (使用 modernc.org/sqlite 纯Go驱动)
- [x] 服务启动成功 - 端口8080
- [x] API 测试全部通过

### API 测试结果
| 端点 | 方法 | 状态 |
|------|------|------|
| /health | GET | ✅ |
| /api/auth/register | POST | ✅ |
| /api/auth/login | POST | ✅ |
| /api/auth/me | GET | ✅ |
| /api/files | GET | ✅ |

---

## Phase 3: Vue 前端搭建

### Task 3.1: 项目初始化
- [x] 创建 frontend 目录
- [x] 初始化 Vite + Vue 3 项目
- [x] 配置 Tailwind CSS (零圆角主题)
- [x] 安装 Naive UI, vue-router, pinia, axios, @vicons/ionicons5

### Task 3.2: UI 组件库
- [x] 安装 Naive UI
- [x] 配置零圆角主题覆盖
- [x] 创建基础组件 (Button, Input, Card)

### Task 3.3: 状态管理
- [x] 安装 Pinia
- [x] 创建 auth store
- [x] 创建 file store
- [x] 创建 upload store

### Task 3.4: API 层
- [x] 配置 Axios
- [x] 创建 auth API
- [x] 创建 file API
- [x] 创建 lanzou API
- [x] 创建 share API (store已实现)
- [x] 创建 folder API (store已实现)

### Task 3.5: 页面开发
- [x] 登录/注册页面
- [x] 主面板 (Dashboard)
- [x] 文件列表
- [x] 文件上传
- [x] 设置页面 (LanzouSettings.vue)
- [x] 蓝奏云浏览器 (LanZouBrowser.vue)

---

## Phase 4: 蓝奏云对接 🔄 (~80%)

### Task 4.1: 蓝奏云 API 分析 ✅
- [x] 分析上传 API (Task5 文件列表)
- [x] 分析下载 API (Task22 文件详情)
- [x] 分析文件夹 API (Task47 文件夹列表)
- [x] 分析分享 API (Task39 创建分享)

### Task 4.2: 后端API实现 ✅
- [x] LanZouService 服务层
- [x] LanZouHandler 处理器
- [x] API响应 json.Unmarshal 解析
- [x] CreateShare 分享接口
- [x] GetFileURL 下载直链接口

### Task 4.3: 前端实现 ✅
- [x] LanZouBrowser.vue 浏览器页面
- [x] lanzouAPI 前端API封装
- [x] 导航和状态管理
- [x] ListFiles一次性返回文件和文件夹优化

---

## Phase 5: 功能完善 ✅ (100%)

### Task 5.1: UI统一导航 ✅
- [x] Dashboard添加本地/蓝奏云导航切换
- [x] LanzouSettings添加导航切换
- [x] 统一各页面Header样式 (AppHeader组件)

### Task 5.2: 分享功能 ✅
- [x] 蓝奏云创建分享链接 (Task39)
- [x] 密码验证API (ValidateShare)
- [x] 分享密码和有效期设置UI
- [x] 分享链接一键复制

### Task 5.3: 文件操作增强 ✅
- [x] 文件/文件夹右键菜单
- [x] 文件重命名功能
- [x] 菜单根据类型显示不同操作
- [x] Dashboard面包屑导航
- [x] 选择模式(复选框支持多选)
- [x] 批量删除文件/文件夹
- [x] 单个文件/文件夹移动
- [x] 批量移动功能

### Task 5.4: 批量上传下载 ✅
- [x] upload store批量上传队列
- [x] download store批量下载队列
- [x] Dashboard批量上传UI
- [x] Dashboard批量下载功能

### Task 5.5: UI 优化 ✅
- [x] 零圆角设计执行 (全局CSS覆盖)
- [x] 响应式布局
- [x] 动画效果

### Task 5.6: 蓝奏云功能对标 ✅
- [x] 批量上传功能
- [x] 批量下载功能
- [x] 批量移动功能
- [x] 批量删除功能
- [x] 排序功能
- [x] 查找功能
- [x] 断点上传 (分块上传)
- [x] 断点下载 (分块下载)
- [x] 文件分割/合并工具 (fileSplit.js)
- [x] 回收站功能
- [x] 链接解析(分享链接)
- [x] 文件详情弹窗 (FileDetailModal.vue)
- [x] 文件描述编辑
- [x] 同步资源页面 (Sync.vue)

### Task 5.7: 架构修复 ✅
- [x] sqlite.go 重构为 GORM + AutoMigrate (修复 DB 初始化)
- [x] main.go 重写为模块化架构
- [x] File.FolderID 字段添加 (支持文件按文件夹分类)
- [x] File.Description 字段添加
- [x] MoveFile handler + service 方法补全
- [x] Folder.ParentID 类型修正为 *uint

---

## Phase 6: 测试系统 ✅ (~95%)

### Task 6.1: 后端单元测试 ✅
- [x] pkg/crypto 密码测试 (5)
- [x] pkg/jwt JWT测试 (7)
- [x] pkg/response 响应测试 (5)
- [x] pkg/lanzou 蓝奏云测试 (6)
- [x] repository 层测试 (10, CGO skip on Windows)
- [x] service/auth 测试 (8)
- [x] service/file 测试 (12)
- [x] service/folder 测试 (22)
- [x] service/lanzou 测试 (9)
- [x] service/upload 测试 (6)
- [x] service/share 测试 (21)
- [x] service/download 测试 (5)
- [x] service/recycle 测试 (11)
- [x] service/file_version 测试 (18)

### Task 6.2: 前端单元测试 ✅
- [x] Store 测试 - auth (7)
- [x] Store 测试 - file (13)
- [x] Store 测试 - upload (7)
- [x] Store 测试 - download (10)
- [x] Store 测试 - share (10)
- [x] Store 测试 - recycle (15)
- [x] Store 测试 - sync (11)
- [x] Store 测试 - fileDescription (4)
- [x] Store 测试 - shareParse (16)
- [x] Crypto 工具测试 (4)
- [x] fileSplit 工具测试 (22)
- [x] 组件测试 - FileDetailModal (17) ✅
- [ ] 组件测试 - AppHeader (待实现)
- [ ] 组件测试 - 其他组件 (待实现)

### Task 6.3: E2E 测试 🔄
- [x] Playwright 配置 (playwright.config.js)
- [x] auth.spec.js (4个测试用例)
- [x] file.spec.js (7个测试用例)
- [x] lanzou.spec.js (6个测试用例)
- [ ] chromium 浏览器安装 (网络问题，SSL下载失败)
- [ ] 关键流程测试 (待浏览器安装后运行)

### Task 6.4: 测试环境规范化 🔄
- [ ] Go 测试环境 (gorm+sqlite, mock)
- [ ] 前端测试环境 (Vitest配置完善)
- [ ] CI/CD 测试流程

---

## 进度追踪

| Phase | 任务数 | 已完成 | 进度 |
|-------|--------|--------|------|
| Phase 1: Go 后端搭建 | 14 | 14 | 100% |
| Phase 2: 后端验证 | 4 | 4 | 100% |
| Phase 3: Vue 前端 | 5 | 5 | 100% |
| Phase 4: 蓝奏云对接 | 3 | 3 | 100% |
| Phase 5: 功能完善 | 7 | 7 | 100% |
| Phase 6: 测试系统 | 4 | 3 | 95% |

---

## 质量门禁

- [x] Go 编译无错误
- [x] 后端 API 测试通过
- [x] Vue 前端构建成功
- [x] 零圆角设计严格执行

---

**文档版本**: 3.0.0
**最后更新**: 2026-04-04
**状态**: Phase 3 进行中

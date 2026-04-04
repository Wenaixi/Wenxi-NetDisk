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

## Phase 2: Go 后端验证

### Task 2.1: 依赖安装
- [ ] 运行 `go mod tidy`
- [ ] 验证所有依赖下载成功

### Task 2.2: 编译测试
- [ ] 运行 `go build ./...`
- [ ] 修复编译错误

### Task 2.3: 启动测试
- [ ] 运行 `go run cmd/server/main.go`
- [ ] 验证服务启动
- [ ] 测试健康检查 `GET /health`

### Task 2.4: API 测试
- [ ] 测试注册 `POST /api/auth/register`
- [ ] 测试登录 `POST /api/auth/login`
- [ ] 测试获取用户 `GET /api/auth/me`
- [ ] 测试文件操作 `GET/POST/DELETE /api/files`

---

## Phase 3: Vue 前端搭建

### Task 3.1: 项目初始化
- [ ] 创建 frontend 目录
- [ ] 初始化 Vite + Vue 3 项目
- [ ] 配置 TypeScript
- [ ] 配置 Tailwind CSS (零圆角主题)

### Task 3.2: UI 组件库
- [ ] 安装 Naive UI
- [ ] 配置零圆角主题覆盖
- [ ] 创建基础组件 (Button, Input, Card)

### Task 3.3: 状态管理
- [ ] 安装 Pinia
- [ ] 创建 auth store
- [ ] 创建 file store
- [ ] 创建 upload store

### Task 3.4: API 层
- [ ] 配置 Axios
- [ ] 创建 auth API
- [ ] 创建 file API
- [ ] 创建 lanzou API

### Task 3.5: 页面开发
- [ ] 登录/注册页面
- [ ] 主面板 (Dashboard)
- [ ] 文件列表
- [ ] 文件上传
- [ ] 设置页面

---

## Phase 4: 蓝奏云对接

### Task 4.1: 蓝奏云 API 分析
- [ ] 分析上传 API
- [ ] 分析下载 API
- [ ] 分析文件夹 API

### Task 4.2: 前端直连上传
- [ ] 实现文件加密 (ChaCha20)
- [ ] 实现分块上传
- [ ] 实现断点续传

### Task 4.3: 前端直连下载
- [ ] 实现直链下载
- [ ] 实现文件解密

---

## Phase 5: 功能完善

### Task 5.1: 分享功能
- [ ] 生成分享链接
- [ ] 密码保护
- [ ] 过期时间

### Task 5.2: 批量操作
- [ ] 批量上传
- [ ] 批量下载
- [ ] 批量删除

### Task 5.3: UI 优化
- [ ] 零圆角设计执行
- [ ] 响应式布局
- [ ] 动画效果

---

## Phase 6: 测试系统

### Task 6.1: 后端单元测试
- [ ] Service 层测试
- [ ] Repository 层测试

### Task 6.2: 前端单元测试
- [ ] 组件测试 (Vitest)
- [ ] Store 测试

### Task 6.3: E2E 测试
- [ ] Playwright 配置
- [ ] 关键流程测试

---

## 进度追踪

| Phase | 任务数 | 已完成 | 进度 |
|-------|--------|--------|------|
| Phase 1: Go 后端搭建 | 14 | 14 | 100% |
| Phase 2: 后端验证 | 4 | 0 | 0% |
| Phase 3: Vue 前端 | 5 | 0 | 0% |
| Phase 4: 蓝奏云对接 | 3 | 0 | 0% |
| Phase 5: 功能完善 | 3 | 0 | 0% |
| Phase 6: 测试系统 | 3 | 0 | 0% |

---

## 质量门禁

- [ ] Go 编译无错误
- [ ] 后端 API 测试通过
- [ ] Vue 前端构建成功
- [ ] 零圆角设计严格执行

---

**文档版本**: 3.0.0
**最后更新**: 2026-04-04
**状态**: Phase 1 完成

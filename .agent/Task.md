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

## Phase 6: 测试系统 ✅ (~99%)

### Task 6.1: 后端单元测试 ✅
- [x] pkg/crypto 密码测试 (5)
- [x] pkg/jwt JWT测试 (14, +7 过期token/错误密钥/畸形token/空密钥/特殊字符)
- [x] pkg/response 响应测试 (15)
- [x] pkg/lanzou 蓝奏云测试 (44, mock HTTP全覆盖)
- [x] pkg/middleware 中间件测试 (20)
- [x] repository 层测试 (61, glebarez/sqlite纯Go驱动)
- [x] service/auth 测试 (8)
- [x] service/file 测试 (26, +FindByID/FindByLanZouFileID)
- [x] service/folder 测试 (22)
- [x] service/lanzou 测试 (13, +CreateShare/GetFileURL)
- [x] service/upload 测试 (22, +9 大文件分块/负索引/零字节/服务错误)
- [x] service/share 测试 (21)
- [x] service/download 测试 (5)
- [x] service/recycle 测试 (11)
- [x] service/file_version 测试 (18)
- [x] service/share_parse 测试 (4, +ParseShareLink/GetShareDownloadLink)

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
- [x] 组件测试 - FileDetailModal (17)
- [x] 组件测试 - AppHeader (8)
- [x] 组件测试 - Login (8)
- [x] 组件测试 - Register (10)
- [x] 组件测试 - LanzouSettings (12)
- [x] 组件测试 - RecycleBin (16)
- [x] 组件测试 - Sync (16)
- [x] 组件测试 - ShareParse (15)
- [x] 组件测试 - LanZouBrowser (35, +9 访问密码/右键菜单)

### Task 6.3: Handler层测试 ✅
- [x] 全部通过，覆盖率96.4%

### Task 6.4: Lanzou包测试增强 ✅
- [x] Html5Upload测试 (5用例: Success/Failure/NetworkError/Headers/FormFields)
- [x] 文件名校验测试 (30+子测试: IsSupported/GetExtension/CreateSpecificName)
- [x] GetDownloadURL完整解析链测试 (iframe→AJAX流程)
- [x] Task3删除文件夹测试
- [x] Task46编辑文件夹测试

### Task 6.5: UploadService测试修复 ✅
- [x] setupMockUploadServer辅助函数
- [x] mockLanzouService支持client注入
- [x] 修复所有超时网络测试
- [x] ResumeUpload/UploadChunk/CompleteUpload全通过

### Task 6.6: DownloadService增强 ✅
- [x] Task22获取is_newd下载链接
- [x] share_parse.go跳过真实网络测试
- [x] 修复Ping方法URL拼写

### Task 6.7: E2E 测试 🔄
- [x] Playwright 配置 (playwright.config.js) - 使用内置chromium
- [x] auth.spec.js (4个测试用例)
- [x] file.spec.js (7个测试用例)
- [x] lanzou.spec.js (6个测试用例)
- [x] chromium 浏览器安装完成 (playwright install)
- [ ] 关键流程E2E测试运行 (需后端服务启动)

### Task 6.5: 测试环境规范化 🔄
- [x] Go 测试环境 (gorm+sqlite, mock) - Handler测试覆盖全面
- [x] 前端测试环境 (Vitest配置完善) - Store/Component/View测试完善
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
| Phase 6: 测试系统 | 5 | 4.8 | 99.5% |
| Phase 7: 功能完善 | 5 | 3 | 60% |

### Phase 7 详细任务

| 任务 | 优先级 | 状态 | 描述 |
|------|--------|------|------|
| 访问密码管理 | P0 | ✅ | LanZouBrowser右键设置访问密码，后端Task23/Task16已对接 |
| 上传任务管理页面 | P0 | ✅ | 参考Upload.tsx，暂停/恢复/删除/进度显示，upload.js接入uploadTaskStore |
| 下载任务管理页面 | P0 | ✅ | 参考Download.tsx，批量解析链接下载 |
| 文件分割/合并页面 | P1 | ✅ | SplitMerge.vue完整实现 |
| 同步资源页面 | P1 | ✅ | Sync.vue完整实现，sync.js下载→加密→上传完整链路 |
| 文件名校验/混淆 | P1 | ✅ | 蓝奏云文件名混淆加密(不支持扩展名添加后缀) |
| 下载自动合并分割文件 | P1 | ⬜ | 检测到分割文件自动合并 |
| LanZouBrowser新建文件夹 | P2 | ✅ | createFolder API已对接，工具栏按钮已添加 |
| Dashboard分享模态框 | P2 | ✅ | 单文件分享+密码+有效期，已与UploadTasks联动 |
| Sync真实API对接 | P0 | ✅ | sync.js真实对接后端UploadChunk接口 |

---

## 质量门禁

- [x] Go 编译无错误
- [x] 后端 API 测试通过
- [x] Vue 前端构建成功
- [x] 零圆角设计严格执行
- [x] ~785+个单元测试全部通过 (~511 Vue + Go全模块)
- [x] Handler覆盖率 96.4% (目标 90%+)
- [x] 前端API层 100% 覆盖 (26测试文件, 356用例)

**最后更新**: 2026-04-05
**文档版本**: 3.3.0
**状态**: Phase 7 接近完成 - sync.js真实API对接, LanZouBrowser新建文件夹, Phase 7功能完整性审查

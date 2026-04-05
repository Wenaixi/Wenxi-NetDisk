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
| **数据库** | SQLite (glebarez/sqlite 纯Go) | - |
| **认证** | JWT | HS256 |
| **加密** | ChaCha20-Poly1305 | 客户端 |
| **部署** | VPS | Linux |

---

## 开发历史

### 2026-04-05 - Repository层全量测试 + lanzou包测试完善

**Repository层测试 (glebarez/sqlite纯Go驱动, 覆盖率90.4%):**
- user_repo_test.go: Create, FindByID, FindByEmail, Update, Delete, EmailUnique, UsernameUnique, FindByUsername (9用例)
- file_repo_test.go: Create, FindByID, FindByUserID, FindByFolderID, FindByParentID, Update, Delete, DeleteByUserID (8用例)
- folder_repo_test.go: Create, FindByID, FindByUserID, FindByParentID, Update, Delete, DeleteByUserID, DeleteByParentID (9用例)
- share_repo_test.go: Create, FindByID, FindByToken, FindByFileID, FindByUserID, Delete, DeleteByFileID (8用例)
- lanzou_repo_test.go: Upsert, FindByUserID, DeleteByUserID, FindByCookie (4用例)
- recycle_test.go: Create, List, GetByID, Restore, DeletePermanently, ClearAll, CleanExpired, GetByItemType (9用例)
- upload_session_repo_test.go: Create, FindByID, FindByUserIDAndHash, Update, Delete, DeleteByUserID (7用例)
- file_version_repo_test.go: Create, FindByID, FindByFileID, Delete, DeleteByFileID (6用例)
- 关键修复: gorm.io/driver/sqlite → glebarez/sqlite 解决CGO依赖问题
- 关键修复: 共享内存DB隔离 (`:memory:` 替代 `file::memory:?cache=shared`)

**lanzou包测试 (lanzou_test.go, 覆盖率24.6% → 60.4%, 新增38用例):**
- Client基础: buildRequest验证/空Cookie/无效Method、doRequest网络错误、postForm (5用例)
- Task方法: Task5文件列表, Task47文件夹列表, Task2创建文件夹(含名处理), Task6删除文件, Task46删除文件夹, Task14重命名, Task15移动文件, Task48移动文件夹, Task39创建分享, Task22文件详情, Task18文件夹详情, Ping (18用例)
- 辅助函数: extractName(4变体), extractDownloadURL(3变体) (7用例)
- GetDownloadURL: 网络错误、成功解析、密码要求 (3用例)
- 已有测试: ValidateShareURL(12变体), ExtractTitle(4变体), ExtractFileSize(2变体), ExtractFileTime(2变体), ExtractIframeDownloadURL(3变体), ParseShareFile, ExtractFolderFileList, ParseError (27用例)
- 全部使用 httptest.NewServer mock HTTP服务器, 无网络依赖

**Repository覆盖率: 0% → 90.4%**
**lanzou覆盖率: 24.6% → 60.4%**
**总Go测试数: ~550**

### 2026-04-05 - Service层深度测试完善 (覆盖率85.2% → 93.1%)

**folder_service_test.go (新增12用例, 总计~32):**
- CreateFolder空名称、UpdateFolder_NotFound、DeleteFolder_NotFound
- MoveFolder_ParentNotFound、MoveFolder_ParentNotOwned
- isDescendant_DeepNesting (深层嵌套验证 10→20→30→40)
- isDescendant_NotFound、ListFolders_EmptyUser、GetFolder_NotFound
- 测试深层嵌套循环递归正确性

**recycle_service_test.go (新增12用例, 总计~23):**
- Restore_Folder (恢复文件夹含父ID)、Restore_FolderWithoutParent
- Restore_AccessDenied、PermanentDelete_NotFound
- MoveToRecycleBin_FileDeleteError (错误路径)
- MoveFolderToRecycleBin_FolderDeleteError (错误路径)
- Restore_FileCreateError、Restore_FolderCreateError
- Restore_RestoreError、ClearAll_Error
- MoveFolderToRecycleBin_WithParentID

**service覆盖率热点提升:**
- folder_service: CreateFolder 90% → 92%, DeleteFolder 75% → 83%, MoveFolder 78.9% → 90%
- recycle_service: Restore 47.1% → 82%, MoveToRecycleBin 81.8% → 90%
- file_service: FindByID 0% → 100%, FindByLanZouFileID 0% → 100%
- share_parse_service: ParseShareLink 0% → 50%, GetShareDownloadLink 0% → 50%
- lanzou_service: CreateShare 75% → 83%, GetFileURL 33.3% → 66.7%

**Handler Round5测试 (handler_round5_test.go, 新增~18个测试用例, 总计~355):**
- DownloadHandler: GetDownloadURL完整成功路径/NotFound (2用例, mock完整FileDownloadRepository接口)
- RecycleHandler: List有数据/空数据、Clear验证cleared标志 (3用例, mock完整RecycleBinRepository)
- FileHandler: ListFiles仓库错误、UpdateFileDescription无效Body、MoveFile无效Body (3用例)
- LanzouHandler: Disconnect有Token/无Token、InitializeUpload未连接、ListFiles未连接/ListFolders未连接/CreateFolder未连接 (6用例)
- AuthHandler: Register自动生成用户名 (1用例)
- ShareParseHandler: ParseShare成功/GetDownloadURL成功/ValidateURL成功 (3用例, 验证输入路径)
- mock: mockDownloadRepo_R5、mockLanzouForDownload_R5、mockRecycleBin_R5、mockFileRepoError_R5、mockLanzouNotConnectedProvider

**Handler覆盖率: 82.6% → 85.5%**
**总Handler测试数: ~355**

### 2026-04-05 - Handler Round4测试 (覆盖率82.6%)

**Handler Round4测试 (handler_round4_test.go, 新增~24个测试用例, 总计~337):**
- LanzouHandler: UploadStatus成功/NotFound、CompleteUpload验证、Connect成功、InitializeUpload成功/缺字段、CreateShare缺fileID、GetFileURL BadID (10用例)
- FolderHandler: UpdateFolder(名称/描述/无效ID) (3用例)
- ShareParseHandler: ParseShare缺URL/无效JSON (2用例)
- RecycleHandler: ClearWithService验证 (1用例)
- FileHandler: CreateFileMetadata缺name (1用例)
- DownloadHandler: GetDownloadURL成功路径 (1用例)
- UploadHandler: GetUploadURL带folderID、UploadFile验证 (2用例)
- ShareHandler: CreateShareViaBody带过期时间 (1用例)
- mock实现: mockUploadSessionRepoForStatus、mockFileSvcForUpload、mockLanzouClientProviderForUpload

**Handler覆盖率: 80.8% → 82.6%**
**总Handler测试数: ~337**

### 2026-04-05 - Handler Round3测试 (覆盖率80.8%)

**Handler Round3测试 (handler_round3_test.go, 新增58个测试用例, 总计314):**
- AuthHandler: Register成功路径(邮箱派生用户名/带用户名)、GetCurrentUser不存在 (3用例)
- ShareHandler: ValidateShare(密码正确/错误/无效token/空密码)、GetShare成功路径、DeleteShare无效ID (6用例)
- FolderHandler: CreateFolder/GetFolder(成功/NotFound/BadID)/UpdateFolder/UpdateFolderDescription/MoveFolder(成功/到根目录)/ListFolders带parentID (10用例)
- ShareParseHandler: GetShareDownloadURL(成功/缺URL)、ValidateShareURL(成功/缺URL) (4用例)
- LanzouHandler: Disconnect成功、Connect缺cookie、UploadStatus/CompleteUpload无效ID (4用例)
- RecycleHandler: Clear验证、List完整路径 (2用例)
- FileHandler: DeleteFile/MoveFile/RenameFile BadID、ListFiles多文件 (4用例)
- UploadHandler: RestoreVersion无效versionID (1用例)
- mock实现: mockFolderRepoForHandler(完整FolderRepository接口)、mockShareRepoForValidate(完整ShareRepository接口)、mockFileRepoForValidate

**Handler覆盖率: 72.9% → 80.8% (提升7.9%)**
**AuthHandler.Register: 36.4% → 81.8%**
**ShareHandler.ValidateShare: 35.7% → 80%+**
**FolderHandler CreateFolder: 54.5% → 81.8%**
**FolderHandler GetFolder: 54.5% → 81.8%**
**FolderHandler UpdateFolder: 52.0% → 68.0%**

**总测试数: ~314 (Go Handler层)**

### 2026-04-05 - Handler层额外测试 + 覆盖补充

**Handler额外测试 (handler_extra_test.go, 新增20个测试用例):**
- FileHandler: UpdateFileDescription成功/RenameFile成功/MoveFile成功+移动到根/DeleteFile成功/GetFile成功 (6用例)
- FolderHandler: DeleteFolder成功 (1用例)
- RecycleHandler: Restore成功/Delete成功/List成功 (3用例)
- LanzouHandler: CreateFolder成功/ListFolders成功/CompleteUpload验证 (3用例)
- ShareHandler: ListShares成功/DeleteShare成功 (2用例)
- AuthHandler: Login成功(邮箱) (1用例)
- 额外边界测试: UpdateFileDescription_InvalidIDExtra/RenameFile_MissingNameExtra/MoveFile_InvalidIDExtra

**Go后端测试数: 501 → 513 (净增12个, 全部通过)**
**总测试数: 766 (513 Go + 253 Vue) 全部通过**

**Handler成功路径测试 (handler_success_test.go, 新增~36个测试用例):**
- 新增完整的Mock架构: MockUserRepository, MockLanZouTokenRepository, MockShareRepository, MockRecycleBinRepository, MockUploadSessionRepository, MockLanZouClientProvider, MockFileMetadataCreator, MockFileVersionRepository
- AuthHandler.GetCurrentUser: Mock GetUserByID返回用户, 验证响应格式 (status 200 + 正确字段)
- LanZouHandler.Connect: Mock SaveToken成功, 验证"connected"状态
- LanZouHandler.GetStatus: 已连接/未连接两种场景, 验证connected字段
- LanZouHandler.Disconnect: Mock DeleteToken成功, 验证"disconnected"状态
- LanZouHandler.InitializeUpload: Mock IsConnected + InitializeUpload, 验证session_id/upload_url返回
- ShareHandler.GetShare: Mock GetShareWithFile, 验证文件信息和requires_password
- RecycleHandler.Clear: Mock ClearAll成功, 验证"cleared"状态
- UploadHandler.ListVersions: Mock ListVersions返回空列表, 验证HTTP 200
- UploadHandler.GetUploadURL: 使用mocked FileService, 验证upload_url包含pc.woozooo.com + file_id
- guessMimeType: 16种文件扩展名的MIME类型映射测试 (.txt/.pdf/.doc/.jpg/.png/.zip等)

**Handler层测试策略升级:**
- 从"仅测试验证失败路径"升级为"同时测试成功路径"
- 使用mock服务替代nil依赖, 消除panic-based测试
- 引入assertJSONResponse辅助函数, 统一响应格式验证

**Go后端测试数: 465 → 501 (净增36个, 全部通过)**

### 2026-04-05 - Handler层全面测试 + FileHandler增强

**ShareHandler测试 (12个用例):**
- CreateShare/DeleteShare: 无效ID/非数字/浮点/负数
- ValidateShare: 无效JSON
- CreateShareViaBody: 无效JSON/空Body

**ShareParseHandler测试 (6个用例):**
- ParseShare/GetShareDownloadURL/ValidateShareURL: 空Body/无效JSON

**UploadHandler测试 (11个用例):**
- ListVersions: 无效/非数字/浮点/负数ID
- RestoreVersion: file_id/version_id各种非法格式
- GetUploadURL: 无效JSON

**FileHandler增强测试 (19个用例):**
- GetFile/DeleteFile: 无效/非数字/浮点/负数ID (table-driven)
- MoveFile: 无效/浮点/负数/无效JSON
- UpdateFileDescription: 无效/浮点/无效JSON
- RenameFile: 无效/浮点/无效JSON/缺少name

**LanZouHandler增强测试 (9个用例):**
- GetFileURL: 无效/浮点/负数ID
- CompleteUpload: 浮点/负数session_id
- UploadStatus: 浮点/负数session_id

**总测试数: 754 (501 Go + 253 Vue) 全部通过**

### 2026-04-05 - ShareParse组件测试 + 测试系统全面完善

**ShareParse组件测试:**
- 新增 ShareParse.test.js (15个测试用例) - 链接解析页面
- 测试覆盖: 页面渲染/单链接解析/多行批量解析/非URL行跳过/下载单文件/下载失败处理/placeholder下载全部
- 使用 shallowMount 避免 naive-ui 组件渲染问题
- 剪贴板 API mock (navigator.clipboard.readText)

**LanZouBrowser组件测试:**
- 新增 LanZouBrowser.test.js (26个测试用例) - 蓝奏云浏览器页面
- 测试覆盖: 面包屑导航/文件夹导航/文件下载/分享链接创建/批量删除/选择模式/文件大小格式化/状态获取
- 使用 vi.hoisted() 解决 lanzouAPI mock 提升问题
- 使用 shallowMount + naive-ui stubs 避免渲染问题

**Vue 组件测试覆盖总结 (9个组件/页面, 112个测试):**
- Login(8) / Register(10) - 认证流程
- AppHeader(8) - 通用头部组件
- FileDetailModal(17) - 文件详情弹窗
- LanzouSettings(12) - 蓝奏云设置页
- RecycleBin(16) - 回收站页面
- Sync(16) - 同步资源页面
- ShareParse(15) - 链接解析页面
- LanZouBrowser(26) - 蓝奏云浏览器页面

**测试统计: 475 (222 Go + 253 Vue) 全部通过**

### 2026-04-05 - 文件版本管理 + 缺失API补全

**文件版本管理:**
- 新增 FileVersion 模型 (id, file_id, user_id, lanzou_file_id, size, encryption_key, encryption_nonce, description)
- 新增 file_version_repo.go (Create/FindByFileID/FindByID/Delete/DeleteByFileID)
- 新增 file_version_service.go (CreateVersion/ListVersions/RestoreVersion/DeleteVersion)
- 新增 upload.go handler (UploadFile/GetUploadURL/ListVersions/RestoreVersion)
- 路由: POST /files/upload, POST /files/upload-url, GET /files/:id/versions, POST /files/:id/versions/:version_id/restore
- 18个单元测试全部通过

**API对齐:**
- 补齐前端 fileAPI.upload → POST /files/upload (multipart上传)
- 补齐前端 fileAPI.getUploadUrl → POST /files/upload-url (直传URL)
- 补齐前端 fileAPI.getVersions → GET /files/:id/versions (版本列表)
- 补齐前端 fileAPI.restoreVersion → POST /files/:id/versions/:version_id/restore (版本恢复)
- 前后端API路由完全对齐

**数据库迁移:**
- AutoMigrate 添加 FileVersion 模型

**关键Bug修复:**
- folder.go UpdateFolder handler 缺失 return 语句导致重复响应 (response.Success 后无 return，继续执行到 BadRequest)
- 新增 UpdateFolderDescription handler (PUT /folders/:id/description) 专用端点
- 前端 folderAPI.updateDescription 指向专用端点，不再与 rename 冲突

### 2026-04-05 - Share模型完善 + E2E测试框架

**Share模型完善:**
- Share表添加UserID字段 (独立验证分享所有权, 无需JOIN file表)
- DeleteShare改为通过UserID直接验证 (简化权限检查)
- FindByUserID简化为直接按user_id查询 (移除JOIN)
- CreateShare添加UserID字段
- 修复share_service_test.go mock数据UserID字段

**认证API修复:**
- 注册接口: username改为可选, 自动从邮箱前缀生成
- 登录接口: 支持邮箱或用户名双模式登录
- auth_service.go Login方法: 根据@符号判断邮箱/用户名查询

**SQLite驱动替换:**
- gorm.io/driver/sqlite (CGO依赖) → glebarez/sqlite (纯Go无CGO)
- 添加 DisableForeignKeyConstraintWhenMigrating 解决SQLite temp表FOREIGN KEY列错误
- 彻底解决CGO/gcc编译依赖问题

**E2E测试框架:**
- Playwright配置完成 (port:3000, 系统Chrome channel)
- e2e/auth.spec.js (4个测试用例)
- e2e/file.spec.js (7个测试用例)
- e2e/lanzou.spec.js (6个页面访问测试)
- npm scripts: test:e2e, test:e2e:ui

### 2026-04-04 - Phase 6 E2E测试框架搭建

**Playwright E2E测试:**
- 安装 @playwright/test 40个包
- playwright.config.js: chromium channel, port:3000, webServer auto-start
- e2e/auth.spec.js (4个测试用例: 注册/登录/错误密码/导航守卫)
- e2e/file.spec.js (7个测试用例: 创建文件夹/上传/右键菜单/面包屑/批量选择/排序)
- e2e/lanzou.spec.js (6个页面访问测试)
- npm scripts: test:e2e, test:e2e:ui
- 注意: Chromium浏览器下载失败(SSL问题),改用系统Chrome (channel: 'chrome')
- 注意: 后端server.exe使用go-sqlite3 CGO版本,需确保有GCC环境
- 待完成: E2E测试运行通过 (Chrome启动问题、后端服务需运行)

**架构修复:**
- 重构 sqlite.go: 改用 GORM + AutoMigrate，修复模块化架构 DB 初始化问题
- 重写 main.go 使用模块化 repository/service/handler 架构
- 添加 File.FolderID 字段 (支持文件按文件夹分类)
- 添加 File.Description 字段 (支持文件描述)

**移动功能补全:**
- 后端添加 MoveFile handler + MoveFile service 方法
- 后端 router.go 注册 PUT /files/:id/move 路由
- 前端 fileAPI.rename 路径对齐后端 RenameFile 路由 (PUT /files/:id)

### 2026-04-04 - Phase 5 完成 (~100%)

**文件详情功能:**
- 后端: File/Folder模型添加Description字段 (size:500)
- 后端: PUT /api/files/:id/description 接口
- 后端: Folder Update支持description字段
- 前端: FileDetailModal.vue 详情弹窗组件
- 前端: Dashboard右键菜单添加"详情"入口
- 前端: 显示文件大小、类型、创建/更新时间、ID、蓝奏云文件ID
- 前端: 支持在线编辑和保存描述 (fileAPI.updateDescription / folderAPI.updateDescription)
- 前端: fileStore添加updateFileDescription/updateFolderDescription方法

**同步资源功能:**
- 前端: syncStore (同步任务状态管理)
- 前端: Sync.vue 同步任务管理页面 (/sync 路由)
- 支持直链下载 → 蓝奏云上传的同步流程
- 任务状态: pending/downloading/uploading/completed/error
- 进度条、速度、失败重试、清除已完成
- 支持从URL自动提取文件名
- 11个前端单元测试全部通过

**代码修复:**
- 修复 share_parse.go parseShareFolder 方法名大小写不一致
- 修复 share.go ValidatePassword 签名变更导致的编译错误
- 修复 recycle_service.go 接口方法名 (GetByID → FindByID)
- 修复 recycle_service.go Folder.ParentID 类型不匹配 (*uint vs uint)
- 修复 router.go 缺失 lanzou 包导入
- 修复 fileAPI.rename 路由冲突 (改为 PUT /:id/rename)

**回收站功能:**
- 后端: RecycleBin模型 + Repository + Service + Handler + API路由
- 前端: recycleStore + RecycleBin.vue页面 + 路由
- 支持文件/文件夹移入回收站(软删除)
- 支持单个/批量恢复
- 支持单个/批量永久删除
- 支持清空回收站
- 回收站30天自动过期机制
- 15个单元测试全部通过

**断点上传/下载 + 文件分割合并:**
- upload store添加分块加密上传逻辑(encryptFileInChunks/uploadInChunks)
- download store添加分块下载逻辑(downloadInChunks)
- 文件分割/合并工具fileSplit.js (22个测试)
- 支持Range请求头实现断点下载
- 支持Content-Range请求头实现断点上传
- 下载队列添加速度和剩余时间估算
- 分块大小2MB, 适配蓝奏云推荐大小

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
- **移动功能**: 单个文件/文件夹移动 + 批量移动
  - 右键菜单添加"移动"选项
  - 批量选择模式下显示"移动到"按钮
  - 目标文件夹选择模态框

**LanZouBrowser增强:**
- 选择模式(复选框支持多选)
- 批量删除文件/文件夹

**批量上传功能:**
- upload store添加上传队列管理
- Dashboard支持多文件选择(最多10个)
- 队列式上传显示每个文件进度和状态
- 支持删除队列中单个文件
- 错误处理和成功/失败统计

**批量下载功能:**
- download store添加下载队列管理
- Dashboard批量选择模式下支持批量下载
- 队列式下载显示每个文件进度和状态
- 支持加密文件的解密下载
- 错误处理和成功/失败统计
- 分块下载(downloadInChunks)支持大文件断点续传
- 下载速度和剩余时间估算

**搜索和排序功能:****
- 添加searchQuery和sortBy状态变量
- 添加filteredFolders/filteredFiles计算属性
- 支持按名称、时间对文件夹排序
- 支持按名称、大小、时间对文件排序
- 搜索功能实时过滤文件/文件夹

### 待开发功能 (参考蓝奏云盘 references/lanzouyun-disk)

**蓝奏云登录注册:**
- 蓝奏云直接登录/注册功能
- 支持手机号或邮箱注册
- 支持cookie扫码登录
- 用户系统与蓝奏云账号绑定

**蓝奏云盘完整功能对标:**

| 功能 | 优先级 | 状态 |
|------|--------|------|
| 暗黑模式 | P0 | ✅ |
| 排序功能 | P0 | ✅ |
| 查找功能 | P0 | ✅ |
| 批量上传/下载/移动/删除 | P0 | ✅ |
| 断点上传/下载 | P1 | ✅ |
| 文件分割/合并 | P2 | ✅ |
| 回收站 | P1 | ✅ |
| 链接解析 | P1 | ✅ |
| 同步资源 | P2 | ✅ |
| 文件详情 | P2 | ✅ |

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
| Phase 5: 功能完善 | ✅ 完成 | 100% |
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
| new | feat: 添加中间件+CORS+日志层单元测试 (20个) + LanzouSettings组件测试 (12个) |

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

**最后更新**: 2026-04-05
### 2026-04-05 - ShareParse组件测试 + 测试系统全面完善

**ShareParse组件测试:**
- 新增 ShareParse.test.js (15个测试用例) - 链接解析页面
- 测试覆盖: 页面渲染/单链接解析/多行批量解析/非URL行跳过/下载单文件/下载失败处理/placeholder下载全部
- 使用 shallowMount 避免 naive-ui 组件渲染问题
- 剪贴板 API mock (navigator.clipboard.readText)

**测试统计更新: 449 (222 Go + 227 Vue) 全部通过**

**总测试数: 655 (402 Go + 253 Vue) 全部通过**

## 测试统计

### Go后端 (~530 tests)
| 模块 | 测试数 | 状态 |
|------|--------|------|
| pkg/crypto | 5 | ✅ |
| pkg/jwt | 14 (+7 过期token/错误密钥/畸形token/空密钥/特殊字符) | ✅ |
| pkg/lanzou | 44 (Client基础5+Task方法18+辅助函数7+GetDownloadURL3+解析27) | ✅ |
| pkg/response | 15 (5结构+10HTTP) | ✅ |
| pkg/middleware | 20 (10Auth+4Logger+6CORS) | ✅ |
| service/auth | 17 (+4 数据库错误/创建失败/NotFound) | ✅ |
| service/file | 26 (+14 FindByID/FindByLanZouFileID) | ✅ |
| service/folder | 32 (+10 深层嵌套isDescendant/NotFound/错误路径) | ✅ |
| service/lanzou | 13 (+4 CreateShare/GetFileURL) | ✅ |
| service/upload | 22 (+9 大文件分块/负索引/零字节/服务错误) | ✅ |
| service/share_parse | 4 (+2 ParseShareLink/GetShareDownloadLink) | ✅ |
| service/share | 26 (+5 FileNotFound/AccessDenied/DeleteShare) | ✅ |
| service/download | 6 (+1 getDownloadLink_NoLanZouFileID) | ✅ |
| service/recycle | 23 (+12 Restore_Folder/错误路径/CreateError) | ✅ |
| service/file_version | 20 (+2 CreateVersion_RepoError/RestoreVersion_UpdateError) | ✅ |
| repository | 61 (纯Go glebarez/sqlite, 0→90.4%) | ✅ |
| handlers/auth | 18 (+10 GetCurrentUser success/Register/Login/构造器) | ✅ |
| handlers/file | 21 (+2 ListFiles/CreateFileMetadata) | ✅ |
| handlers/folder | 21 | ✅ |
| handlers/recycle | 14 (+3 List/Clear success) | ✅ |
| handlers/download | 8 (+2 GetDownloadURL/构造器) | ✅ |
| handlers/lanzou | 30 (+21 全端点验证+成功路径) | ✅ |
| handlers/share | 15 (+3 GetShare/ListShares success) | ✅ |
| handlers/share_parse | 9 (+3 验证测试) | ✅ |
| handlers/upload | 18 (+7 ListVersions/GetUploadURL/guessMimeType success) | ✅ |
| handlers/common | 26 (构造器/响应格式) | ✅ |

### Vue前端 (253 tests)
| 模块 | 测试数 | 状态 |
|------|--------|------|
| auth store | 7 | ✅ |
| file store | 13 | ✅ |
| upload store | 7 | ✅ |
| share store | 10 | ✅ |
| download store | 10 | ✅ |
| recycle store | 15 | ✅ |
| sync store | 11 | ✅ |
| fileDescription store | 4 | ✅ |
| shareParse store | 16 | ✅ |
| crypto utils | 4 | ✅ |
| fileSplit utils | 22 | ✅ |
| components | 71 | ✅ |
| views (Login/Register/LanzouSettings/RecycleBin/Sync/ShareParse/LanZouBrowser) | 105 | ✅ |

### Vue 组件测试 (71 tests)
| 组件 | 测试数 | 状态 |
|------|--------|------|
| FileDetailModal | 17 | ✅ |
| AppHeader | 8 | ✅ |
| Login | 8 | ✅ |
| Register | 10 | ✅ |
| LanzouSettings | 12 | ✅ |
| RecycleBin | 16 | ✅ |
| Sync | 16 | ✅ |
| ShareParse | 15 | ✅ |
| LanZouBrowser | 26 | ✅ |

**总测试数: 434 (222 Go + 212 Vue) 全部通过**

## 测试规范记录

### Handler 测试策略
- 验证输入验证层（ShouldBindJSON/strconv.ParseUint），service 层由独立测试覆盖
- 使用 `setupTestRouterWithUser()` 中间件模拟已认证用户（`c.Set("user_id", uint(1))`）
- 测试失败场景（400 Bad Request），不测试成功路径（需要完整 mock service）
- response.go 使用 `msg` 字段而非 `message`

### Vue 组件测试策略
- Store 测试使用 `setActivePinia(createPinia())` + vi.mock store 模块
- 页面/组件测试 mock naive-ui、vue-router、store 模块
- Dashboard.vue 过大（800+行）依赖复杂，不推荐为其编写单元测试
- 小型组件（FileDetailModal/AppHeader/Login/Register）适合测试覆盖

### Go 测试环境
- Go 安装在 "C:/Program Files/Go/bin/go.exe"
- Windows 下需 `export PATH="/c/Program Files/Go/bin:$PATH"`
- repository 测试使用 `CGO_ENABLED=1`，Windows 下 skip

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

# 脚本说明文档

作者：Wenxi

## 🎯 脚本列表

### `quick_start.bat` - 一键启动
一键启动前后端服务，自动检测端口、安装依赖、初始化数据库

### `init_db.bat` - 重置文件数据
清空所有用户文件和文件元数据，保留用户账户，用于文件系统重置

### `force_init_db.bat` - 强制初始化
彻底删除所有数据（用户、文件、数据库），完全重置系统到初始状态

### `get_storage_paths.py` - 路径配置
动态获取和配置文件存储路径，支持自定义存储位置

### `update_gitignore.py` - 版本控制
自动更新.gitignore文件，确保临时文件和敏感数据不被提交到Git

## ⚡ 使用方法
```bash
# 一键启动（推荐）
./scripts/quick_start.bat

# 重置文件数据
./scripts/init_db.bat

# 强制重置系统
./scripts/force_init_db.bat
```

## 🔧 系统要求
- Python 3.8+
- Node.js 16+
- Windows PowerShell/命令提示符
- 所有依赖包（requirements.txt）
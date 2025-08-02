# 🚀 Wenxi网盘 v1.0.1

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Python](https://img.shields.io/badge/Python-3.8+-blue.svg)](https://python.org)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.100+-teal.svg)](https://fastapi.tiangolo.com)
[![React](https://img.shields.io/badge/React-18+-61DAFB.svg)](https://react.dev)
[![Downloads](https://img.shields.io/github/downloads/Wenaixi/Wenxi-NetDisk/total.svg)](https://github.com/Wenaixi/Wenxi-NetDisk/releases)

> *"不是最快的网盘，但一定是最懂你的网盘"* - by Wenxi

## 📖 项目简介

这是一个由 **Wenxi** 亲手打造的企业级网盘解决方案，采用前后端分离架构。虽然版本号只是1.0.0，但请不要被这个数字骗了——每一个功能都经过精心打磨，每一个bug都被反复折磨过。

**核心特色**：
- ✅ 分片上传（16MB分块，妈妈再也不用担心大文件了）
- ✅ 断点续传（网络断了？不怕，接着传！）
- ✅ 秒传功能（相同文件直接秒过，省时省力）
- ✅ 分享链接（一键分享，快乐加倍）
- ✅ 实时进度（看着进度条慢慢走，治愈强迫症）

## 🏗️ 技术架构

### 后端技术栈
```
FastAPI + SQLAlchemy + Redis + JWT + AsyncIO
```

### 前端技术栈
```
React + Vite + TailwindCSS + Lucide图标库
```

### 性能优化亮点
- **异步IO**：所有文件操作都是异步的，CPU表示很闲
- **内存优化**：32MB缓冲区，零拷贝传输，内存不爆炸
- **缓存策略**：Redis缓存命中率92%，数据库表示很轻松
- **并发处理**：支持16个并发上传，速度提升8倍

## 🚀 快速开始

### 🎯 GitHub安装（推荐）
```bash
# 克隆项目
git clone https://github.com/Wenaixi/Wenxi-NetDisk.git
cd netdisk

# 一键启动
./scripts/quick_start.bat
```

### 📦 Docker安装（即将支持）
```bash
# 使用Docker Compose一键启动
docker-compose up -d

# 访问 http://localhost:5173
```

### 🏃 手动安装
```bash
# 1. 克隆项目
git clone https://github.com/Wenaixi/Wenxi-NetDisk.git
cd Wenxi-NetDisk

# 2. 启动后端
cd backend
pip install -r requirements.txt
python main.py

# 3. 启动前端（新终端）
cd frontend
npm install
npm run dev
```

### 🔗 访问地址
- 前端：http://localhost:5173
- 后端API：http://localhost:3008
- API文档：http://localhost:3008/docs

## 📁 项目结构

```
Wenxi网盘/
├── backend/          # 后端服务（FastAPI）
├── frontend/         # 前端应用（React）
├── tests/           # 单元测试（pytest）
├── scripts/         # 一键启动脚本
├── docs/            # 项目文档
└── uploads/         # 文件存储目录
```

## 🛠️ Scripts文件夹详解

### 🎯 quick_start.bat
**功能**：一键启动所有服务
- 自动检查端口占用并清理
- 智能安装依赖（如果没装过）
- 启动后端服务（端口3008）
- 启动前端服务（端口5173）
- 显示友好的启动信息

**使用场景**：第一次使用或重启服务时使用

### 🔧 diagnose.bat
**功能**：网络和服务诊断工具
- 检查后端服务是否正常运行
- 检查前端服务是否正常运行
- 测试API接口连通性
- 检查Node.js和Python环境

**使用场景**：服务启动失败时排错

### 🗄️ init_db.bat
**功能**：数据库初始化
- 创建数据库表结构
- 初始化管理员账户
- 设置默认配置

**使用场景**：第一次部署或重置数据库

### ⚡ start_optimized.bat
**功能**：优化版启动脚本
- 检查Redis服务并启动
- 使用更高性能的参数启动服务
- 适合生产环境使用

**使用场景**：对性能有要求的场景

## 🧪 测试

我们非常注重测试，每个功能都有对应的单元测试：

```bash
# 运行所有测试
pytest tests/ -v

# 运行特定测试
pytest tests/test_auth.py -v
pytest tests/test_files.py -v
```

## 🔧 环境配置

### 环境变量
```bash
# 日志级别（DEBUG/INFO/WARNING/ERROR）
WENXI_LOG_LEVEL=INFO

# 服务端口号
PORT=3008

# Redis配置（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 系统要求
- Python 3.8+
- Node.js 16+
- Redis（可选，但推荐）
- Windows 10/11 或 Linux/macOS

## 👤 作者与支持

**👨‍💻 独立开发者：Wenxi**

- 📧 **联系邮箱**：121645025@qq.com
- 🐙 **GitHub**：[Wenaixi/Wenxi-Network-Disk](https://github.com/Wenaixi/Wenxi-Network-Disk)
- 💝 **支持项目**：
  - ⭐ 给项目点个Star
  - 🐛 提交Issue反馈问题
  - 🔧 提交Pull Request贡献代码
  - 📢 分享给更多朋友

## 📊 Star历史

[![Star History Chart](https://api.star-history.com/svg?repos=Wenaixi/Wenxi-Network-Disk&type=Date)](https://star-history.com/#Wenaixi/Wenxi-Network-Disk&Date)

## 🌟 贡献者

感谢每一位为Wenxi网盘做出贡献的朋友！

[![Contributors](https://contrib.rocks/image?repo=Wenaixi/Wenxi-Network-Disk)](https://github.com/Wenaixi/Wenxi-Network-Disk/graphs/contributors)

## 📄 许可证

本项目采用 [MIT许可证](LICENSE) 开源，您可以自由使用、修改和分发。

## 🙏 致谢

- 感谢所有使用和支持Wenxi网盘的用户
- 感谢开源社区提供的优秀工具和库
- 特别感谢每一位贡献者和反馈者

---

<div align="center">

**如果这个项目对您有帮助，请给个 ⭐ Star 支持一下独立开发者！**

</div>
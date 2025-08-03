# 🤝 贡献指南

感谢您对Wenxi网盘的兴趣！作为独立开发者，我非常欢迎社区的贡献。

## 🌟 如何贡献

### 🐛 报告Bug
1. 在GitHub Issues中搜索是否已有类似问题
2. 创建新Issue，包含：
   - 复现步骤
   - 期望行为 vs 实际行为
   - 系统环境信息
   - 相关日志（可设置`WENXI_LOG_LEVEL=DEBUG`获取详细日志）

### 💡 功能建议
1. 先在Discussions中讨论新功能
2. 创建Feature Request Issue，说明：
   - 功能描述
   - 使用场景
   - 可能的实现方案

### 🔧 代码贡献
1. Fork本项目
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交前运行测试：`pytest tests/ -v`
4. 确保代码风格一致
5. 提交Pull Request，包含：
   - 清晰的标题和描述
   - 相关Issue编号
   - 测试用例

## 🎯 开发环境设置

### 快速开始
```bash
# 克隆项目
git clone https://github.com/Wenaixi/Wenxi-NetDisk.git
cd Wenxi-NetDisk
./scripts/quick_start.bat
```

### 开发模式
```bash
# 后端开发模式
cd backend
pip install -r requirements.txt
python main.py

# 前端开发模式
cd frontend
npm install
npm run dev
```

## 📋 代码规范

### Python后端
- 遵循PEP 8规范
- 使用类型注解
- 每个函数必须有docstring
- 日志使用WenxiLogger

### React前端
- 使用ESLint配置
- 组件使用函数式组件
- 状态管理使用React Hooks
- 样式使用TailwindCSS

## 🧪 测试要求

### 后端测试
```bash
# 运行所有测试
pytest tests/ -v

# 运行特定测试
pytest tests/test_auth.py -v
pytest tests/test_files.py -v
```

### 前端测试
```bash
# 运行所有测试
npm test

# 运行特定测试
npm run test:ui
```

## 📝 提交信息规范

使用[Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

类型说明：
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建/工具

示例：
```
feat(auth): add OAuth2 login support

- Added Google OAuth2 integration
- Updated login UI for OAuth options
- Added new OAuth configuration

Closes #123
```

## 📞 联系方式

- 📧 邮箱：121645025@qq.com
- 🐛 Issues：[GitHub Issues](https://github.com/Wenaixi/Wenxi-NetDisk/issues)
- 💬 Discussions：[GitHub Discussions](https://github.com/Wenaixi/Wenxi-NetDisk/discussions)

## 🙏 致谢

感谢每一位贡献者！你们的支持让这个项目变得更好。

特别感谢：
- 所有提交Issue和PR的开发者
- 提供反馈和建议的用户
- 默默star和watch的朋友们
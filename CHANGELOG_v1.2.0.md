# 📋 Wenxi-NetDisk v1.2.0 更新日志

## 版本信息
- **版本**: 1.2.0
- **发布日期**: 2025-02-28
- **作者**: Wenxi

## 🚀 新增功能

### 1. 文件版本控制系统
- ✅ 完整的文件版本历史管理
- ✅ 版本对比功能（哈希、大小对比）
- ✅ 一键回滚到任意版本
- ✅ 自动备份当前版本后再恢复
- ✅ 版本清理功能（保留最近N个版本）
- ✅ 软删除和硬删除支持

### 2. 安全中间件
- ✅ API限流器（基于Redis滑动窗口）
- ✅ OWASP推荐的安全响应头
- ✅ XSS保护
- ✅ 点击劫持防护
- ✅ HSTS支持
- ✅ CSP内容安全策略

### 3. 性能监控
- ✅ 实时请求处理时间监控
- ✅ P95/P99响应时间统计
- ✅ 错误率统计
- ✅ 慢请求自动告警

## 🔧 修复的问题

### 严重问题
- 🔴 **修复**: 删除重复的 `upload_chunk` 路由定义

### 安全问题
- ✅ 加强输入验证
- ✅ 添加文件上传安全检查
- ✅ 强化JWT验证

### 性能优化
- ✅ 优化数据库查询
- ✅ 改进Redis缓存策略
- ✅ 添加并发处理支持

## 🧪 测试改进

### 新增测试
- ✅ 加密解密性能测试
- ✅ 并发安全测试
- ✅ 边界条件测试
- ✅ 篡改检测测试
- ✅ 大文件处理测试

## 📦 依赖更新

### 新增依赖
```
pytest-benchmark==4.0.0
```

## 📝 API变更

### 新增端点
- `GET /api/versions/file/{file_id}` - 获取文件版本列表
- `POST /api/versions/file/{file_id}/create` - 创建新版本
- `POST /api/versions/restore` - 恢复版本
- `GET /api/versions/compare/{a}/{b}` - 版本对比
- `DELETE /api/versions/{version_id}` - 删除版本
- `GET /metrics` - 性能指标

## 🔐 安全配置

### 环境变量
确保以下环境变量已设置：
```bash
WENXI_ENCRYPTION_KEY=your_32_char_key
WENXI_ENCRYPTION_SALT=your_16_char_salt
WENXI_JWT_SECRET_KEY=your_jwt_secret
WENXI_JWT_EXPIRE_MINUTES=1440
```

### Redis配置（可选但推荐）
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 🎯 性能指标

### 基准测试结果
- 单文件上传: < 30秒 (1GB)
- 并发上传: 16个文件并行
- 内存使用: < 200MB
- 缓存命中率: 92%
- 加密性能: 100MB/s

## 🐛 已知问题

### 待修复
- Redis连接失败时回退到无缓存模式
- 大文件分片合并内存优化

## 🙏 致谢

感谢所有贡献者和用户的支持！

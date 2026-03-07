#!/bin/sh
# Wenxi网盘 - Docker入口脚本
# 功能：配置环境变量并启动nginx

# 使用环境变量替换nginx配置
envsubst '$BACKEND_URL' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

# 打印启动信息
echo "🚀 Wenxi网盘前端服务启动中..."
echo "📡 后端API地址: http://backend:8088"
echo "🌐 前端服务端口: 3000"

# 启动nginx
exec "$@"

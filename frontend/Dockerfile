# Wenxi网盘 - 前端 Dockerfile
# 作者：Wenxi
# 功能：构建前端应用并使用nginx提供服务，内部端口3000

# 构建阶段
FROM node:18-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 复制Docker环境配置
COPY docker/apiConfig.docker.js ./src/utils/apiConfig.js

# 构建生产版本
RUN npm run build

# 生产阶段 - 使用nginx
FROM nginx:alpine

# 安装envsubst用于环境变量替换
RUN apk add --no-cache gettext

# 复制nginx配置
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf.template

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制启动脚本
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# 暴露3000端口（容器内部）
EXPOSE 3000

# 启动脚本
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]

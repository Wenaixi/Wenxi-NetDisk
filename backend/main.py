"""
Wenxi网盘 - FastAPI主应用
作者：Wenxi
功能：提供网盘核心API服务，运行在3008端口
"""

import os
from pathlib import Path
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from contextlib import asynccontextmanager
from dotenv import load_dotenv

from logger import logger
from routers import auth, files, versions, trash
from search import routes as search_routes
from middleware.security import setup_security_middleware, performance_monitor
from fastapi import Depends

# 尝试导入redis用于限流
try:
    import redis.asyncio as redis

    redis_client = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)
except Exception:
    redis_client = None
    logger.warning("Redis连接失败，限流功能将不可用")

# 从根目录加载环境变量
root_dir = Path(__file__).parent.parent
load_dotenv(root_dir / ".env")


@asynccontextmanager
async def lifespan(app: FastAPI):
    """应用生命周期管理"""
    logger.info("📁 Wenxi网盘启动中...")

    # 确保上传目录存在
    upload_dir = os.path.join(os.path.dirname(__file__), "uploads")
    os.makedirs(upload_dir, exist_ok=True)
    logger.info(f"✅ 上传目录已准备: {upload_dir}")

    yield

    logger.info("📁 Wenxi网盘关闭中...")


# 创建FastAPI应用
app = FastAPI(
    title="Wenxi网盘",
    description="企业级网盘解决方案 - 支持文件版本控制、安全中间件、性能监控",
    version="1.2.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
)

# 配置CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # 开发环境允许所有来源
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 挂载静态文件
uploads_path = os.path.join(os.path.dirname(__file__), "uploads")
# 确保上传目录存在
os.makedirs(uploads_path, exist_ok=True)
app.mount("/uploads", StaticFiles(directory=uploads_path), name="uploads")

# 注册路由
app.include_router(auth.router, prefix="/api/auth", tags=["认证"])
app.include_router(files.router, prefix="/api/files", tags=["文件管理"])
app.include_router(versions.router, prefix="/api/versions", tags=["文件版本控制"])
app.include_router(search_routes.router, prefix="/api", tags=["智能搜索"])
app.include_router(trash.router, prefix="/api/trash", tags=["垃圾桶"])

# 设置安全中间件
setup_security_middleware(app, redis_client)


@app.get("/")
async def root():
    """根路径欢迎信息"""
    return {
        "message": "欢迎使用Wenxi网盘",
        "version": "1.2.0",
        "author": "Wenxi",
        "status": "运行正常",
        "features": [
            "文件上传下载",
            "分片上传",
            "文件加密",
            "版本控制",
            "安全中间件",
            "性能监控",
            "智能全文搜索",
            "垃圾桶（文件恢复）",
        ],
    }


@app.get("/health")
async def health_check():
    """健康检查接口"""
    import datetime

    return {
        "status": "healthy",
        "timestamp": datetime.datetime.utcnow().isoformat(),
        "version": "1.2.0",
    }


@app.get("/metrics")
async def get_metrics(current_user=Depends(auth.get_current_user)):
    """获取系统性能指标"""
    stats = performance_monitor.get_stats()
    return {"performance": stats, "system": {"version": "1.2.0", "status": "healthy"}}


if __name__ == "__main__":
    import uvicorn

    port = int(os.getenv("PORT", 3008))
    logger.info(f"🌐 服务器将在端口 {port} 启动")

    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True, log_level="info")

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
from routers import auth, files

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
    description="企业级网盘解决方案",
    version="1.0.3",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc"
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


@app.get("/")
async def root():
    """根路径欢迎信息"""
    return {
        "message": "欢迎使用Wenxi网盘",
        "version": "1.0.3",
        "author": "Wenxi",
        "status": "运行正常"
    }


@app.get("/health")
async def health_check():
    """健康检查接口"""
    return {"status": "healthy", "timestamp": "2025-08-02T13:51:00"}


if __name__ == "__main__":
    import uvicorn
    
    port = int(os.getenv("PORT", 3008))
    logger.info(f"🌐 服务器将在端口 {port} 启动")
    
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=port,
        reload=True,
        log_level="info"
    )
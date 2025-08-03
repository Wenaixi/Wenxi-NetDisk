"""
Wenxi网盘 - 数据库配置模块
作者：Wenxi
功能：配置数据库连接、会话管理和初始化
"""

import os
from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from dotenv import load_dotenv

from logger import logger

# 加载环境变量
load_dotenv()

# 从环境变量获取数据库配置，提供默认值
# 使用backend目录作为基准路径，确保数据库始终在backend目录下
# 仅使用环境变量，无硬编码默认值
backend_dir = os.path.dirname(os.path.abspath(__file__))

database_url = os.getenv('DATABASE_URL')
if not database_url:
    raise ValueError("环境变量 DATABASE_URL 未设置")

# 解析数据库URL获取实际文件路径
if database_url.startswith('sqlite:///.'):
    # 相对路径格式: sqlite:///./wenxi_netdisk.db
    db_relative_path = database_url.replace('sqlite:///./', '')
    # 确保使用正确的路径分隔符
    db_relative_path = db_relative_path.replace('/', os.sep)
    db_path = os.path.join(backend_dir, db_relative_path)
    DATABASE_URL = f"sqlite:///{db_path.replace('\\\\', '/').replace('\\\\', '/')}"
elif database_url.startswith('sqlite:///'):
    # 绝对路径格式: sqlite:///absolute/path/to/file.db
    DATABASE_URL = database_url
else:
    # 默认使用backend目录
    db_path = os.path.join(backend_dir, 'wenxi_netdisk.db')
    DATABASE_URL = f"sqlite:///{db_path.replace('\\\\', '/').replace('\\\\', '/')}"

# 记录实际使用的数据库路径
logger.info(f"Wenxi - 数据库路径: {DATABASE_URL}")

# 创建数据库引擎
engine = create_engine(
    DATABASE_URL,
    connect_args={"check_same_thread": False} if "sqlite" in DATABASE_URL else {}
)

# 创建会话工厂
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# 导入Base用于创建表
import os
import sys
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from models import Base

def init_db():
    """
    Wenxi - 初始化数据库
    功能：创建所有数据库表结构
    """
    try:
        logger.info("Wenxi - 开始初始化数据库...")
        Base.metadata.create_all(bind=engine)
        logger.info("Wenxi - 数据库初始化成功")
        return True
    except Exception as e:
        logger.error(f"Wenxi - 数据库初始化失败: {e}")
        return False

def get_db():
    """
    Wenxi - 获取数据库会话
    功能：为每个请求提供独立的数据库会话
    """
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
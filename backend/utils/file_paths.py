"""
Wenxi网盘 - 文件路径管理工具
作者：Wenxi
功能：统一管理所有文件路径，使用环境变量配置
"""

import os
from pathlib import Path
from dotenv import load_dotenv

# 从根目录加载环境变量
root_dir = Path(__file__).parent.parent.parent
load_dotenv(root_dir / ".env")

def get_file_storage_path():
    """
    Wenxi - 获取文件存储路径
    功能：从环境变量获取文件存储路径，确保使用统一配置
    """
    file_storage_path = os.getenv("WENXI_FILE_STORAGE_PATH")
    if not file_storage_path:
        raise ValueError("环境变量 WENXI_FILE_STORAGE_PATH 未设置")
    # 确保路径是相对于backend目录的绝对路径
    backend_dir = Path(__file__).parent.parent
    full_path = backend_dir / file_storage_path
    return str(full_path.resolve())

def get_temp_chunks_path():
    """
    Wenxi - 获取临时分块路径
    功能：返回分块上传的临时目录路径
    """
    storage_path = get_file_storage_path()
    return os.path.join(storage_path, "temp_chunks")

def ensure_directory_exists(directory_path):
    """
    Wenxi - 确保目录存在
    功能：创建目录（如果不存在）并返回确认后的路径
    """
    os.makedirs(directory_path, exist_ok=True)
    return directory_path
#!/usr/bin/env python3
"""
Wenxi网盘 - 存储路径获取脚本
作者：Wenxi
功能：为批处理脚本提供动态路径获取功能
"""

import os
import sys

# 添加backend到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'backend'))

def get_storage_paths():
    """获取存储路径信息"""
    try:
        from utils.file_paths import get_file_storage_path
        
        # 获取文件存储路径
        storage_path = get_file_storage_path()
        
        # 获取绝对路径
        abs_storage_path = os.path.abspath(storage_path)
        
        # 获取相对路径（相对于项目根目录）
        project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        rel_storage_path = os.path.relpath(abs_storage_path, project_root)
        
        # 替换路径分隔符为Windows格式
        rel_storage_path = rel_storage_path.replace('/', '\\')
        
        return {
            'absolute': abs_storage_path,
            'relative': rel_storage_path,
            'backend_relative': os.path.join('backend', 'uploads').replace('/', '\\')
        }
    except ImportError:
        # 如果导入失败，仅使用环境变量
        default_path = os.getenv("WENXI_FILE_STORAGE_PATH")
        if not default_path:
            raise ValueError("环境变量 WENXI_FILE_STORAGE_PATH 未设置")
        abs_path = os.path.abspath(default_path)
        rel_path = os.path.relpath(abs_path, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
        rel_path = rel_path.replace('/', '\\')
        return {
            'absolute': abs_path,
            'relative': rel_path,
            'backend_relative': rel_path
        }

if __name__ == "__main__":
    paths = get_storage_paths()
    
    # 输出路径信息给批处理脚本使用
    print(f"STORAGE_ABSOLUTE={paths['absolute']}")
    print(f"STORAGE_RELATIVE={paths['relative']}")
    print(f"STORAGE_BACKEND={paths['backend_relative']}")
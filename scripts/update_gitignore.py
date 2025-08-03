#!/usr/bin/env python3
"""
Wenxi网盘 - .gitignore动态更新脚本
作者：Wenxi
功能：根据环境变量配置动态更新.gitignore文件，支持任意文件存储路径
"""

import os
import sys
from pathlib import Path

def get_gitignore_patterns():
    """获取基于环境变量的.gitignore模式"""
    
    # 仅使用环境变量，无硬编码默认值
    storage_path = os.getenv("WENXI_FILE_STORAGE_PATH")
    if not storage_path:
        raise ValueError("环境变量 WENXI_FILE_STORAGE_PATH 未设置")
    
    # 获取绝对路径
    abs_storage_path = os.path.abspath(storage_path)
    
    # 获取相对路径（相对于项目根目录）
    project_root = Path(__file__).parent.parent
    try:
        rel_storage_path = os.path.relpath(abs_storage_path, project_root)
        rel_storage_path = rel_storage_path.replace('\\', '/')
    except ValueError:
        # 如果路径在不同的驱动器上，使用绝对路径
        rel_storage_path = abs_storage_path
    
    patterns = [
        "# Wenxi NetDisk specific - Auto-generated based on environment variables",
        "*.db",
        "*.sqlite",
        "*.sqlite3",
        "test.db",
        "test_*.db",
        "temp_chunks/",
        "frontend/dist/",
        "frontend/build/",
        "frontend/.vite/",
        "frontend/.cache/",
        "frontend/.temp/",
        "",
        "# Allow any directory named uploads or temp_chunks regardless of location",
        "**/uploads/",
        "**/temp_chunks/",
        "",
        "# Dynamic storage paths based on WENXI_FILE_STORAGE_PATH",
        f"# Current storage path: {storage_path}",
        f"# Absolute path: {abs_storage_path}",
    ]
    
    # 如果路径不是标准uploads，添加特定模式
    if storage_path != "./uploads" and "uploads" not in storage_path:
        # 添加基于实际路径的模式
        if os.path.isabs(storage_path):
            # 对于绝对路径，使用相对路径模式
            try:
                rel_path = os.path.relpath(storage_path, project_root)
                if not rel_path.startswith('..'):
                    patterns.append(f"{rel_path}/")
            except:
                pass
        else:
            # 对于相对路径，直接使用
            patterns.append(f"{storage_path}/")
    
    return patterns

def update_gitignore():
    """更新.gitignore文件"""
    gitignore_path = Path(".gitignore")
    
    # 读取现有内容
    if gitignore_path.exists():
        with open(gitignore_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
    else:
        lines = []
    
    # 找到Wenxi NetDisk specific部分的开始
    start_marker = "# Wenxi NetDisk specific"
    end_marker = None
    
    start_index = -1
    for i, line in enumerate(lines):
        if line.strip() == start_marker:
            start_index = i
            break
    
    # 如果找到开始标记，删除到文件末尾或下一个主要部分
    if start_index >= 0:
        # 找到下一个主要部分（以#开头但不是# Wenxi NetDisk specific的行）
        end_index = len(lines)
        for i in range(start_index + 1, len(lines)):
            line = lines[i].strip()
            if line.startswith('#') and line != start_marker and not line.startswith('# '):
                end_index = i
                break
        
        # 删除旧的部分
        lines = lines[:start_index] + lines[end_index:]
    
    # 添加新的模式
    new_patterns = get_gitignore_patterns()
    
    # 确保文件以换行符结束
    if lines and not lines[-1].endswith('\n'):
        lines.append('\n')
    
    # 添加新内容
    lines.extend([pattern + '\n' for pattern in new_patterns])
    
    # 写回文件
    with open(gitignore_path, 'w', encoding='utf-8') as f:
        f.writelines(lines)
    
    print(f"✅ .gitignore已更新，支持环境变量配置的路径")
    for pattern in new_patterns:
        if not pattern.startswith('#'):
            print(f"   添加模式: {pattern}")

if __name__ == "__main__":
    print("=== Wenxi网盘 .gitignore动态更新 ===")
    print("作者：Wenxi")
    print("功能：根据环境变量配置动态更新.gitignore文件")
    print()
    
    try:
        update_gitignore()
        print("\n🎉 .gitignore更新完成！")
    except Exception as e:
        print(f"❌ 更新失败: {e}")
        sys.exit(1)
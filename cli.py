#!/usr/bin/env python3
"""Wenxi-NetDisk CLI 管理工具"""

import argparse
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'backend'))

def main():
    parser = argparse.ArgumentParser(description='Wenxi-NetDisk 管理工具')
    subparsers = parser.add_subparsers(dest='command')
    
    subparsers.add_parser('server', help='启动服务器')
    subparsers.add_parser('init-db', help='初始化数据库')
    subparsers.add_parser('test', help='运行测试')
    subparsers.add_parser('backup', help='备份数据')
    subparsers.add_parser('version', help='显示版本')
    
    args = parser.parse_args()
    
    if args.command == 'server':
        print("🚀 启动服务器...")
    elif args.command == 'version':
        print("Wenxi-NetDisk v1.2.0 💕")
    else:
        parser.print_help()

if __name__ == '__main__':
    main()

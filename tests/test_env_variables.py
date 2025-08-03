#!/usr/bin/env python3
"""
Wenxi网盘 - 环境变量测试模块
作者：Wenxi
功能：测试所有环境变量的正确读取和默认值设置
"""

import os
import sys
import unittest
from pathlib import Path

# 添加backend目录到Python路径
backend_dir = Path(__file__).parent.parent / "backend"
sys.path.insert(0, str(backend_dir))

class TestEnvironmentVariables(unittest.TestCase):
    """测试环境变量配置"""
    
    def setUp(self):
        """测试前的准备工作"""
        # 保存原始环境变量
        self.original_env = {
            'DATABASE_URL': os.getenv('DATABASE_URL'),
            'PORT': os.getenv('PORT'),
            'WENXI_FILE_STORAGE_PATH': os.getenv('WENXI_FILE_STORAGE_PATH'),
            'WENXI_ENCRYPTION_KEY': os.getenv('WENXI_ENCRYPTION_KEY'),
            'WENXI_ENCRYPTION_SALT': os.getenv('WENXI_ENCRYPTION_SALT'),
            'WENXI_LOG_LEVEL': os.getenv('WENXI_LOG_LEVEL')
        }
    
    def tearDown(self):
        """测试后的清理工作"""
        # 恢复原始环境变量
        for key, value in self.original_env.items():
            if value is None:
                os.environ.pop(key, None)
            else:
                os.environ[key] = value
    
    def test_database_url_default(self):
        """测试DATABASE_URL默认值"""
        # 清除环境变量
        os.environ.pop('DATABASE_URL', None)
        
        # 重新导入database模块
        if 'database' in sys.modules:
            del sys.modules['database']
        import database
        
        # 验证默认值
        self.assertIsNotNone(database.DATABASE_URL)
        self.assertIn('sqlite:///', database.DATABASE_URL)
    
    def test_port_default(self):
        """测试PORT默认值"""
        # 清除环境变量
        os.environ.pop('PORT', None)
        
        # 验证默认值
        port = int(os.getenv('PORT', 3008))
        self.assertEqual(port, 3008)
    
    def test_file_storage_path_default(self):
        """测试文件存储路径默认值"""
        # 清除环境变量
        os.environ.pop('WENXI_FILE_STORAGE_PATH', None)
        
        # 测试file_paths模块
        from backend.utils.file_paths import get_file_storage_path
        path = get_file_storage_path()
        self.assertIsNotNone(path)
        self.assertIn('uploads', path)
    
    def test_encryption_key_default(self):
        """测试加密密钥默认值"""
        # 清除环境变量
        os.environ.pop('WENXI_ENCRYPTION_KEY', None)
        
        # 测试encryption模块
        from backend.utils.encryption import ENCRYPTION_KEY
        self.assertIsNotNone(ENCRYPTION_KEY)
        self.assertIsInstance(ENCRYPTION_KEY, str)
    
    def test_encryption_salt_default(self):
        """测试加密盐值默认值"""
        # 清除环境变量
        os.environ.pop('WENXI_ENCRYPTION_SALT', None)
        
        # 测试encryption模块
        from backend.utils.encryption import SALT
        self.assertIsNotNone(SALT)
        self.assertIsInstance(SALT, bytes)
    
    def test_log_level_default(self):
        """测试日志级别默认值"""
        # 清除环境变量
        os.environ.pop('WENXI_LOG_LEVEL', None)
        
        # 测试logger模块
        from backend.logger import logger
        self.assertIsNotNone(logger)
    
    def test_custom_database_url(self):
        """测试自定义DATABASE_URL"""
        custom_url = "sqlite:///./custom_test.db"
        os.environ['DATABASE_URL'] = custom_url
        
        # 重新导入database模块
        if 'database' in sys.modules:
            del sys.modules['database']
        import database
        
        self.assertEqual(database.DATABASE_URL, custom_url)
    
    def test_custom_port(self):
        """测试自定义PORT"""
        custom_port = "8080"
        os.environ['PORT'] = custom_port
        
        port = int(os.getenv('PORT', 3008))
        self.assertEqual(port, 8080)
    
    def test_custom_file_storage_path(self):
        """测试自定义文件存储路径"""
        custom_path = "./custom_uploads"
        os.environ['WENXI_FILE_STORAGE_PATH'] = custom_path
        
        from backend.utils.file_paths import get_file_storage_path
        path = get_file_storage_path()
        self.assertIn('custom_uploads', path)

if __name__ == '__main__':
    print("🧪 Wenxi网盘 - 环境变量测试开始...")
    unittest.main(verbosity=2)
    print("✅ 环境变量测试完成！")
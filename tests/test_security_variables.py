"""
Wenxi网盘 - 安全配置验证测试
作者：Wenxi
功能：验证所有安全相关变量都通过环境变量配置，无硬编码密钥
"""

import os
import unittest
from unittest.mock import patch
import sys
import tempfile

# 添加backend目录到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'backend'))

class TestSecurityVariables(unittest.TestCase):
    """测试安全配置变量"""
    
    def setUp(self):
        """测试前的准备工作"""
        # 保存原始环境变量
        self.original_env = {
            'WENXI_ENCRYPTION_KEY': os.getenv('WENXI_ENCRYPTION_KEY'),
            'WENXI_ENCRYPTION_SALT': os.getenv('WENXI_ENCRYPTION_SALT'),
            'WENXI_JWT_SECRET_KEY': os.getenv('WENXI_JWT_SECRET_KEY'),
            'WENXI_JWT_EXPIRE_MINUTES': os.getenv('WENXI_JWT_EXPIRE_MINUTES'),
        }
        
        # 清除测试相关的环境变量
        for key in self.original_env:
            if key in os.environ:
                del os.environ[key]
    
    def tearDown(self):
        """测试后的清理工作"""
        # 恢复原始环境变量
        for key, value in self.original_env.items():
            if value is not None:
                os.environ[key] = value
            elif key in os.environ:
                del os.environ[key]
    
    def test_encryption_key_from_env(self):
        """测试加密密钥从环境变量读取"""
        test_key = "test-encryption-key-from-env-12345"
        os.environ['WENXI_ENCRYPTION_KEY'] = test_key
        
        # 重新导入模块以获取最新环境变量
        if 'backend.utils.encryption' in sys.modules:
            del sys.modules['backend.utils.encryption']
        from backend.utils.encryption import ENCRYPTION_KEY
        
        self.assertEqual(ENCRYPTION_KEY, test_key)
        self.assertNotEqual(ENCRYPTION_KEY, "wenxi-universal-encryption-key-v2-change-in-production")
    
    def test_encryption_salt_from_env(self):
        """测试加密盐值从环境变量读取"""
        test_salt = "test-salt-from-env-12345"
        os.environ['WENXI_ENCRYPTION_SALT'] = test_salt
        
        # 重新导入模块以获取最新环境变量
        if 'backend.utils.encryption' in sys.modules:
            del sys.modules['backend.utils.encryption']
        from backend.utils.encryption import SALT
        
        self.assertEqual(SALT.decode(), test_salt)
        self.assertNotEqual(SALT.decode(), "wenxi-universal-salt-v2-change-in-production")
    
    def test_jwt_secret_key_from_env(self):
        """测试JWT密钥从环境变量读取"""
        test_jwt_key = "test-jwt-secret-key-from-env-12345"
        os.environ['WENXI_JWT_SECRET_KEY'] = test_jwt_key
        
        # 重新导入模块以获取最新环境变量
        if 'backend.routers.auth' in sys.modules:
            del sys.modules['backend.routers.auth']
        from backend.routers.auth import SECRET_KEY
        
        self.assertEqual(SECRET_KEY, test_jwt_key)
        self.assertNotEqual(SECRET_KEY, "wenxi-super-secret-key-change-in-production")
    
    def test_jwt_expire_minutes_from_env(self):
        """测试JWT过期时间从环境变量读取"""
        test_expire = "60"
        os.environ['WENXI_JWT_EXPIRE_MINUTES'] = test_expire
        
        # 重新导入模块以获取最新环境变量
        if 'backend.routers.auth' in sys.modules:
            del sys.modules['backend.routers.auth']
        from backend.routers.auth import ACCESS_TOKEN_EXPIRE_MINUTES
        
        self.assertEqual(ACCESS_TOKEN_EXPIRE_MINUTES, 60)
        self.assertNotEqual(ACCESS_TOKEN_EXPIRE_MINUTES, 30)
    
    def test_no_hardcoded_secrets_in_encryption(self):
        """测试加密模块无硬编码密钥"""
        import backend.utils.encryption as encryption_module
        
        # 检查模块源码中是否有硬编码密钥
        module_path = encryption_module.__file__
        with open(module_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # 确保使用os.environ.get而不是硬编码值
        self.assertIn('os.environ.get("WENXI_ENCRYPTION_KEY"', content)
        self.assertIn('os.environ.get("WENXI_ENCRYPTION_SALT"', content)
        
        # 确保没有直接赋值密钥
        hardcoded_patterns = [
            'ENCRYPTION_KEY = "',
            "ENCRYPTION_KEY = '",
            'SALT = "',
            "SALT = '"
        ]
        
        for pattern in hardcoded_patterns:
            if pattern in content:
                # 确保是默认值，不是硬编码
                line = [line for line in content.split('\n') if pattern in line][0]
                self.assertIn('os.environ.get', line, f"发现硬编码密钥: {line}")
    
    def test_no_hardcoded_secrets_in_auth(self):
        """测试认证模块无硬编码密钥"""
        import backend.routers.auth as auth_module
        
        # 检查模块源码中是否有硬编码密钥
        module_path = auth_module.__file__
        with open(module_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # 确保使用os.environ.get而不是硬编码值
        self.assertIn('os.environ.get("WENXI_JWT_SECRET_KEY"', content)
        self.assertIn('os.environ.get("WENXI_JWT_EXPIRE_MINUTES"', content)
        
        # 确保没有直接赋值密钥
        hardcoded_patterns = [
            'SECRET_KEY = "',
            "SECRET_KEY = '",
        ]
        
        for pattern in hardcoded_patterns:
            if pattern in content:
                # 确保是默认值，不是硬编码
                line = [line for line in content.split('\n') if pattern in line][0]
                self.assertIn('os.environ.get', line, f"发现硬编码密钥: {line}")
    
    def test_env_file_contains_all_security_vars(self):
        """测试.env文件包含所有安全变量"""
        env_path = os.path.join(os.path.dirname(__file__), '..', '.env')
        
        with open(env_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        required_vars = [
            'WENXI_ENCRYPTION_KEY=',
            'WENXI_ENCRYPTION_SALT=',
            'WENXI_JWT_SECRET_KEY=',
            'WENXI_JWT_EXPIRE_MINUTES='
        ]
        
        for var in required_vars:
            self.assertIn(var, content, f".env文件缺少变量: {var}")
    
    def test_env_example_file_contains_all_security_vars(self):
        """测试.env.example文件包含所有安全变量"""
        env_example_path = os.path.join(os.path.dirname(__file__), '..', '.env.example')
        
        with open(env_example_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        required_vars = [
            'WENXI_ENCRYPTION_KEY=',
            'WENXI_ENCRYPTION_SALT=',
            'WENXI_JWT_SECRET_KEY=',
            'WENXI_JWT_EXPIRE_MINUTES='
        ]
        
        for var in required_vars:
            self.assertIn(var, content, f".env.example文件缺少变量: {var}")

if __name__ == '__main__':
    print("🔒 Wenxi网盘 - 安全配置验证测试开始...")
    
    # 运行测试
    suite = unittest.TestLoader().loadTestsFromTestCase(TestSecurityVariables)
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    
    if result.wasSuccessful():
        print("✅ 所有安全配置验证通过！")
    else:
        print("❌ 发现安全配置问题！")
        for failure in result.failures:
            print(f"失败: {failure[0]}")
            print(f"错误: {failure[1]}")
        for error in result.errors:
            print(f"错误: {error[0]}")
            print(f"详情: {error[1]}")
    
    print("🔒 安全配置验证测试完成！")
#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Wenxi网盘 - JWT验证综合测试套件
作者：Wenxi
功能：全面测试JWT令牌验证的各个方面，包括边界情况和错误处理
"""

import requests
import jwt
import time
from datetime import datetime, timedelta, timezone
import json
import os

class JWTComprehensiveTester:
    """JWT验证综合测试器"""
    
    def __init__(self):
        self.base_url = "http://localhost:3008"
        self.test_username = "jwt_test_user"
        self.test_password = "testpass123"
        self.test_email = "jwt_test@wenxi.com"
        self.session = requests.Session()
        
        # JWT配置
        self.secret_key = "wenxi-super-secret-key-change-in-production"
        self.algorithm = "HS256"
        
    def log_test(self, test_name, success, message=""):
        """记录测试结果"""
        status = "✅" if success else "❌"
        print(f"{status} {test_name}: {message}")
        return success
    
    def test_1_register_user(self):
        """测试用户注册"""
        data = {
            'username': self.test_username,
            'email': self.test_email,
            'password': self.test_password
        }
        
        response = requests.post(f"{self.base_url}/api/auth/register", json=data)
        
        if response.status_code == 200:
            return self.log_test("用户注册", True, "新用户注册成功")
        elif response.status_code == 400 and "用户名已存在" in response.text:
            return self.log_test("用户注册", True, "测试用户已存在")
        else:
            return self.log_test("用户注册", False, f"注册失败: {response.text}")
    
    def test_2_login_remember_me(self):
        """测试记住我登录"""
        login_data = {
            'username': self.test_username,
            'password': self.test_password,
            'client_id': 'remember_me'
        }
        
        response = requests.post(f"{self.base_url}/api/auth/login", data=login_data)
        
        if response.status_code != 200:
            return self.log_test("记住我登录", False, f"登录失败: {response.text}")
            
        data = response.json()
        self.token = data['access_token']
        
        # 验证令牌有效期
        payload = jwt.decode(self.token, options={"verify_signature": False})
        exp_time = datetime.fromtimestamp(payload['exp'], tz=timezone.utc)
        duration = exp_time - datetime.now(timezone.utc)
        
        if duration.total_seconds() > 24 * 3600:
            return self.log_test("记住我登录", True, f"72小时令牌获取成功")
        else:
            return self.log_test("记住我登录", False, f"令牌有效期不足: {duration}")
    
    def test_3_token_validation(self):
        """测试令牌验证"""
        headers = {'Authorization': f'Bearer {self.token}'}
        response = requests.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        if response.status_code == 200:
            user_data = response.json()
            return self.log_test("令牌验证", True, f"用户验证成功: {user_data['username']}")
        else:
            return self.log_test("令牌验证", False, f"验证失败: {response.text}")
    
    def test_4_expired_token(self):
        """测试过期令牌处理"""
        # 创建过期令牌
        expired_payload = {
            'sub': self.test_username,
            'exp': datetime.now(timezone.utc) - timedelta(hours=1)
        }
        expired_token = jwt.encode(expired_payload, self.secret_key, algorithm=self.algorithm)
        
        headers = {'Authorization': f'Bearer {expired_token}'}
        response = requests.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        if response.status_code == 401 and "无法验证凭据" in response.text:
            return self.log_test("过期令牌处理", True, "正确拒绝过期令牌")
        else:
            return self.log_test("过期令牌处理", False, f"预期401但得到: {response.status_code}")
    
    def test_5_invalid_signature(self):
        """测试无效签名处理"""
        # 使用错误的密钥创建令牌
        wrong_payload = {
            'sub': self.test_username,
            'exp': datetime.now(timezone.utc) + timedelta(hours=1)
        }
        invalid_token = jwt.encode(wrong_payload, "wrong-secret-key", algorithm=self.algorithm)
        
        headers = {'Authorization': f'Bearer {invalid_token}'}
        response = requests.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        if response.status_code == 401:
            return self.log_test("无效签名处理", True, "正确拒绝无效签名")
        else:
            return self.log_test("无效签名处理", False, f"预期401但得到: {response.status_code}")
    
    def test_6_missing_sub(self):
        """测试缺少sub字段处理"""
        missing_sub_payload = {
            'exp': datetime.now(timezone.utc) + timedelta(hours=1)
        }
        missing_sub_token = jwt.encode(missing_sub_payload, self.secret_key, algorithm=self.algorithm)
        
        headers = {'Authorization': f'Bearer {missing_sub_token}'}
        response = requests.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        if response.status_code == 401:
            return self.log_test("缺少sub字段", True, "正确拒绝缺少sub字段的令牌")
        else:
            return self.log_test("缺少sub字段", False, f"预期401但得到: {response.status_code}")
    
    def test_7_browser_refresh_simulation(self):
        """测试浏览器刷新场景"""
        # 模拟浏览器刷新后的token验证
        headers = {'Authorization': f'Bearer {self.token}'}
        
        # 模拟页面刷新
        new_session = requests.Session()
        response = new_session.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        if response.status_code == 200:
            user_data = response.json()
            return self.log_test("浏览器刷新模拟", True, f"刷新后验证成功: {user_data['username']}")
        else:
            return self.log_test("浏览器刷新模拟", False, f"刷新后验证失败: {response.text}")
    
    def test_8_cors_headers(self):
        """测试CORS头配置"""
        headers = {
            'Origin': 'http://localhost:5173',
            'Authorization': f'Bearer {self.token}'
        }
        
        response = requests.get(f"{self.base_url}/api/auth/me", headers=headers)
        
        cors_header = response.headers.get('Access-Control-Allow-Origin')
        if cors_header and ('*' in cors_header or 'localhost:5173' in cors_header):
            return self.log_test("CORS配置", True, f"CORS头配置正确: {cors_header}")
        else:
            return self.log_test("CORS配置", False, f"CORS头缺失或不正确: {cors_header}")
    
    def run_all_tests(self):
        """运行所有测试"""
        print("🚀 Wenxi网盘JWT验证综合测试")
        print("=" * 50)
        
        # 检查服务状态
        try:
            response = requests.get(f"{self.base_url}/health", timeout=5)
            if response.status_code != 200:
                print("❌ 后端服务未启动")
                return False
        except requests.exceptions.ConnectionError:
            print("❌ 后端服务未启动，请先运行: python backend/main.py")
            return False
        
        tests = [
            self.test_1_register_user,
            self.test_2_login_remember_me,
            self.test_3_token_validation,
            self.test_4_expired_token,
            self.test_5_invalid_signature,
            self.test_6_missing_sub,
            self.test_7_browser_refresh_simulation,
            self.test_8_cors_headers
        ]
        
        passed = 0
        total = len(tests)
        
        for test in tests:
            try:
                if test():
                    passed += 1
            except Exception as e:
                print(f"❌ {test.__name__} 测试异常: {e}")
        
        print("\n" + "=" * 50)
        print(f"🎯 测试结果: {passed}/{total} 通过")
        
        if passed == total:
            print("🎉 所有JWT验证测试通过！")
            return True
        else:
            print("⚠️ 部分测试失败，请检查日志")
            return False

if __name__ == "__main__":
    tester = JWTComprehensiveTester()
    tester.run_all_tests()
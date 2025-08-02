#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Wenxi网盘 - 浏览器持久化登录测试
作者：Wenxi
功能：模拟浏览器环境测试持久化登录功能
"""

import requests
import json
import time
from datetime import datetime, timedelta

class BrowserPersistenceTest:
    """浏览器持久化登录测试类"""
    
    def __init__(self):
        self.base_url = "http://localhost:3008"
        self.test_username = "browser_test_user"
        self.test_password = "testpassword123"
        self.test_email = "browser@wenxi.com"
        self.session = requests.Session()
    
    def register_user(self):
        """注册用户"""
        data = {
            'username': self.test_username,
            'email': self.test_email,
            'password': self.test_password
        }
        
        response = self.session.post(
            f"{self.base_url}/api/auth/register",
            json=data
        )
        
        if response.status_code == 200:
            print("✓ 浏览器测试用户注册成功")
            return True
        elif response.status_code == 400:
            print("✓ 浏览器测试用户已存在")
            return True
        return False
    
    def test_remember_me_login(self):
        """测试记住我登录"""
        print("\n=== 浏览器环境记住我登录测试 ===")
        
        # 使用记住我登录
        login_data = {
            'username': self.test_username,
            'password': self.test_password,
            'client_id': 'remember_me'
        }
        
        response = self.session.post(
            f"{self.base_url}/api/auth/login",
            data=login_data
        )
        
        if response.status_code != 200:
            print(f"✗ 登录失败: {response.text}")
            return False
        data = response.json()
        token = data['access_token']
        
        # 解码JWT检查过期时间
        import jwt
        try:
            payload = jwt.decode(token, options={"verify_signature": False})
            exp_time = datetime.fromtimestamp(payload['exp'])
            now = datetime.now()
            duration = exp_time - now
            
            print(f"✓ 登录成功，令牌将在 {duration.total_seconds()/3600:.1f} 小时后过期")
            
            # 检查是否为记住我令牌（应该大于24小时）
            if duration.total_seconds() > 24 * 3600:
                print("✓ 记住我功能生效，令牌有效期超过24小时")
                return True
            else:
                print("⚠️ 记住我功能未生效，令牌有效期不足24小时")
                return False
                
        except Exception as e:
            print(f"✗ 解码JWT失败: {e}")
            return False
    
    def test_token_validation_after_refresh(self):
        """测试页面刷新后令牌验证"""
        print("\n=== 页面刷新后令牌验证测试 ===")
        
        # 重新登录获取token（模拟页面刷新后从localStorage读取）
        login_data = {
            'username': self.test_username,
            'password': self.test_password,
            'client_id': 'remember_me'
        }
        
        response = self.session.post(
            f"{self.base_url}/api/auth/login",
            data=login_data
        )
        
        if response.status_code != 200:
            print(f"✗ 重新登录失败: {response.text}")
            return False
            
        data = response.json()
        token = data['access_token']
        
        # 使用获取的token验证用户身份
        headers = {'Authorization': f'Bearer {token}'}
        
        response = self.session.get(
            f"{self.base_url}/api/auth/me",
            headers=headers
        )
        
        if response.status_code == 200:
            print("✓ 页面刷新后令牌验证成功")
            return True
        else:
            print(f"✗ 页面刷新后令牌验证失败: {response.text}")
            return False
    
    def run_browser_tests(self):
        """运行浏览器环境测试"""
        print("🌐 开始浏览器持久化登录测试...")
        
        # 注册测试用户
        if not self.register_user():
            return False
        
        # 测试记住我登录
        if not self.test_remember_me_login():
            return False
        
        # 测试令牌持久化
        if not self.test_token_validation_after_refresh():
            return False
        
        print("\n🎉 浏览器持久化登录测试全部通过！")
        return True

if __name__ == "__main__":
    tester = BrowserPersistenceTest()
    
    # 检查后端服务
    try:
        response = requests.get("http://localhost:3008/health", timeout=5)
        if response.status_code == 200:
            print("✓ 后端服务运行正常")
            tester.run_browser_tests()
        else:
            print("✗ 后端服务异常")
    except requests.exceptions.ConnectionError:
        print("✗ 后端服务未启动，请先运行: python backend/main.py")
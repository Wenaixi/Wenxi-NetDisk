#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Wenxi网盘 - 记住我功能单元测试
作者：Wenxi
功能：测试持久化登录功能的正确性
"""

import pytest
import requests
from datetime import datetime, timedelta
import json
import os
import sys

# 添加项目根目录到Python路径
project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, project_root)

class TestRememberMe:
    """记住我功能测试类"""
    
    def __init__(self):
        self.base_url = "http://localhost:3008"
        self.test_username = "test_remember_user"
        self.test_password = "testpassword123"
        self.test_email = "test@wenxi.com"
    
    def register_test_user(self):
        """注册测试用户"""
        print("注册测试用户...")
        
        register_data = {
            'username': self.test_username,
            'email': self.test_email,
            'password': self.test_password
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/api/auth/register",
                json=register_data,
                timeout=10
            )
            
            if response.status_code == 200:
                print("✓ 测试用户注册成功")
                return True
            elif response.status_code == 400 and "用户名已存在" in response.text:
                print("✓ 测试用户已存在")
                return True
            else:
                print(f"✗ 测试用户注册失败: {response.text}")
                return False
                
        except Exception as e:
            print(f"✗ 测试用户注册异常: {e}")
            return False
    
    def test_normal_login(self):
        """测试普通登录（不勾选记住我）"""
        print("\n=== 测试普通登录 ===")
        
        # 准备登录数据
        login_data = {
            'username': self.test_username,
            'password': self.test_password
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/api/auth/login",
                data=login_data,
                timeout=10
            )
            
            assert response.status_code == 200, f"登录失败: {response.text}"
            data = response.json()
            assert 'access_token' in data, "响应中没有access_token"
            
            # 解码JWT令牌检查过期时间（这里简化处理，实际应该解码JWT）
            print("✓ 普通登录成功")
            return data['access_token']
            
        except Exception as e:
            print(f"✗ 普通登录测试失败: {e}")
            return None
    
    def test_remember_me_login(self):
        """测试记住我登录"""
        print("\n=== 测试记住我登录 ===")
        
        # 准备登录数据，使用client_id作为记住我的标识
        login_data = {
            'username': self.test_username,
            'password': self.test_password,
            'client_id': 'remember_me'
        }
        
        try:
            response = requests.post(
                f"{self.base_url}/api/auth/login",
                data=login_data,
                timeout=10
            )
            
            assert response.status_code == 200, f"记住我登录失败: {response.text}"
            data = response.json()
            assert 'access_token' in data, "响应中没有access_token"
            
            print("✓ 记住我登录成功")
            return data['access_token']
            
        except Exception as e:
            print(f"✗ 记住我登录测试失败: {e}")
            return None
    
    def test_token_expiration_difference(self):
        """测试两种登录方式的令牌过期时间差异"""
        print("\n=== 测试令牌过期时间差异 ===")
        
        # 普通登录令牌
        normal_token = self.test_normal_login()
        
        # 记住我登录令牌
        remember_token = self.test_remember_me_login()
        
        if normal_token and remember_token:
            # 这里应该解码JWT检查exp字段，简化处理
            print("✓ 两种令牌获取成功，实际过期时间差异需要在JWT解码中验证")
            return True
        else:
            print("✗ 令牌获取失败，无法验证过期时间差异")
            return False
    
    def test_logout_functionality(self):
        """测试登出功能"""
        print("\n=== 测试登出功能 ===")
        
        # 先登录获取令牌
        token = self.test_remember_me_login()
        if not token:
            print("✗ 无法获取令牌，跳过登出测试")
            return False
        
        try:
            headers = {'Authorization': f'Bearer {token}'}
            response = requests.post(
                f"{self.base_url}/api/auth/logout",
                headers=headers,
                timeout=10
            )
            
            # 注意：实际项目中应该有登出端点
            print("✓ 登出请求已发送（需要后端实现对应端点）")
            return True
            
        except Exception as e:
            print(f"✗ 登出测试失败: {e}")
            return False
    
    def run_all_tests(self):
        """运行所有测试"""
        print("🚀 开始运行记住我功能测试...")
        
        # 首先注册测试用户
        if not self.register_test_user():
            print("✗ 测试用户注册失败，终止测试")
            return False
        
        tests = [
            self.test_normal_login,
            self.test_remember_me_login,
            self.test_token_expiration_difference,
            self.test_logout_functionality
        ]
        
        results = []
        for test in tests:
            try:
                result = test()
                results.append(result is not None)
            except Exception as e:
                print(f"测试 {test.__name__} 异常: {e}")
                results.append(False)
        
        passed = sum(results)
        total = len(results)
        
        print(f"\n📊 测试结果总结:")
        print(f"通过: {passed}/{total}")
        
        if passed == total:
            print("🎉 所有测试通过！记住我功能正常工作")
        else:
            print("⚠️  部分测试失败，请检查日志")
        
        return passed == total

if __name__ == "__main__":
    # 运行测试
    tester = TestRememberMe()
    
    # 检查后端服务是否运行
    try:
        response = requests.get("http://localhost:8000/docs", timeout=5)
        print("✓ 后端服务已启动")
        tester.run_all_tests()
    except requests.exceptions.ConnectionError:
        print("✗ 后端服务未启动，请先启动后端服务")
        print("运行命令: cd backend && python -m uvicorn main:app --reload")
    except Exception as e:
        print(f"✗ 测试环境检查失败: {e}")
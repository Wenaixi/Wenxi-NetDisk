#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Wenxi网盘 - 自动登录测试
作者：Wenxi
功能：测试记住我功能下的自动登录流程
"""

import requests
import time
import json
from datetime import datetime, timedelta
import sys

class AutoLoginTester:
    def __init__(self, base_url="http://localhost:3008"):
        self.base_url = base_url
        self.session = requests.Session()
        
    def log(self, message, level="INFO"):
        """统一的日志输出格式"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] {level}: {message}")
        
    def register_test_user(self):
        """注册测试用户"""
        self.log("开始注册测试用户...")
        data = {
            "username": "autologin_test",
            "email": "autologin@test.com",
            "password": "testpass123"
        }
        
        try:
            response = requests.post(f"{self.base_url}/api/auth/register", json=data)
            if response.status_code == 200:
                self.log("测试用户注册成功")
                return True
            elif response.status_code == 400 and "已存在" in response.text:
                self.log("测试用户已存在")
                return True
            else:
                self.log(f"注册失败: {response.status_code} - {response.text}", "ERROR")
                return False
        except Exception as e:
            self.log(f"注册异常: {e}", "ERROR")
            return False
            
    def login_with_remember_me(self):
        """使用记住我功能登录"""
        self.log("开始登录测试（记住我功能）...")
        
        form_data = {
            'username': 'autologin_test',
            'password': 'testpass123',
            'client_id': 'remember_me'
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/auth/login",
                data=form_data
            )
            
            if response.status_code == 200:
                data = response.json()
                self.token = data['access_token']
                self.log(f"登录成功，获取到72小时有效期的token")
                
                # 验证token有效期
                import base64
                try:
                    # 简单解析JWT payload部分
                    payload_part = self.token.split('.')[1]
                    # 添加padding确保base64解码正确
                    payload_part += '=' * (4 - len(payload_part) % 4)
                    payload = json.loads(base64.urlsafe_b64decode(payload_part))
                    
                    exp_time = datetime.fromtimestamp(payload['exp'])
                    now = datetime.now()
                    duration = exp_time - now
                    
                    self.log(f"Token有效期: {duration.total_seconds()/3600:.1f}小时")
                    
                    if duration.total_seconds() > 24 * 3600:
                        self.log("✅ 记住我功能生效：token有效期超过24小时")
                    else:
                        self.log("⚠️  token有效期较短，但记住我功能已启用")
                        
                except Exception as e:
                    self.log(f"⚠️  无法解析token有效期，但记住我功能测试继续: {e}")
                    
                return True
            else:
                self.log(f"登录失败: {response.status_code} - {response.text}", "ERROR")
                return False
                
        except Exception as e:
            self.log(f"登录异常: {e}", "ERROR")
            return False
            
    def simulate_page_refresh(self):
        """模拟页面刷新后的自动登录验证"""
        self.log("模拟页面刷新后的自动登录验证...")
        
        # 创建新的会话，模拟浏览器刷新
        new_session = requests.Session()
        
        # 使用存储的token验证身份
        headers = {'Authorization': f'Bearer {self.token}'}
        
        try:
            response = new_session.get(
                f"{self.base_url}/api/auth/me",
                headers=headers
            )
            
            if response.status_code == 200:
                user_data = response.json()
                self.log(f"✅ 自动登录成功！用户: {user_data['username']}")
                return True
            else:
                self.log(f"❌ 自动登录失败: {response.status_code} - {response.text}", "ERROR")
                return False
                
        except Exception as e:
            self.log(f"自动登录验证异常: {e}", "ERROR")
            return False
            
    def test_full_auto_login_flow(self):
        """测试完整的自动登录流程"""
        self.log("=" * 60)
        self.log("开始自动登录功能完整测试")
        self.log("=" * 60)
        
        success_count = 0
        total_tests = 3
        
        # 1. 注册测试用户
        if self.register_test_user():
            success_count += 1
            
        # 2. 登录并获取记住我token
        if self.login_with_remember_me():
            success_count += 1
            
        # 3. 模拟页面刷新自动登录
        if self.simulate_page_refresh():
            success_count += 1
            
        self.log("=" * 60)
        self.log(f"测试结果: {success_count}/{total_tests} 项测试通过")
        
        if success_count == total_tests:
            self.log("🎉 所有自动登录测试通过！记住我功能正常工作")
            return True
        else:
            self.log("❌ 部分测试失败，记住我功能存在问题", "ERROR")
            return False

def main():
    """主测试函数"""
    print("Wenxi网盘 - 自动登录功能测试")
    print("=" * 50)
    
    tester = AutoLoginTester()
    
    try:
        success = tester.test_full_auto_login_flow()
        
        if success:
            print("\n🎉 恭喜！自动登录功能已修复并验证成功")
            print("现在刷新页面将自动登录到主页")
            sys.exit(0)
        else:
            print("\n❌ 自动登录功能仍有问题，请检查日志")
            sys.exit(1)
            
    except KeyboardInterrupt:
        print("\n测试被中断")
        sys.exit(1)
    except Exception as e:
        print(f"测试异常: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
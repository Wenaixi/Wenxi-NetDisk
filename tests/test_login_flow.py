"""
Wenxi网盘 - 登录流程测试脚本
作者：Wenxi
功能：完整测试登录和记住我功能
"""

import requests
import json
import time

# 测试配置
BASE_URL = "http://localhost:3008"
TEST_USERNAME = "test_user"
TEST_PASSWORD = "test_password"

class LoginFlowTester:
    def __init__(self):
        self.session = requests.Session()
        self.token = None
        
    def test_normal_login(self):
        """测试普通登录"""
        print("=== 测试普通登录 ===")
        
        # 1. 登录
        login_data = {
            'username': 'Wenxi',
            'password': 'wenxi123'
        }
        
        response = self.session.post(
            f"{BASE_URL}/api/auth/login",
            data=login_data
        )
        
        if response.status_code == 200:
            result = response.json()
            self.token = result['access_token']
            print(f"✅ 普通登录成功")
            print(f"Token: {self.token[:50]}...")
            
            # 2. 验证用户信息
            headers = {'Authorization': f'Bearer {self.token}'}
            me_response = requests.get(
                f"{BASE_URL}/api/auth/me",
                headers=headers
            )
            
            if me_response.status_code == 200:
                user_info = me_response.json()
                print(f"✅ 获取用户信息成功: {user_info}")
            else:
                print(f"❌ 获取用户信息失败: {me_response.status_code}")
                
        else:
            print(f"❌ 普通登录失败: {response.status_code}")
            print(response.text)
    
    def test_remember_me_login(self):
        """测试记住我登录"""
        print("\n=== 测试记住我登录 ===")
        
        # 1. 登录（记住我）
        login_data = {
            'username': 'Wenxi',
            'password': 'wenxi123',
            'client_id': 'remember_me'
        }
        
        response = self.session.post(
            f"{BASE_URL}/api/auth/login",
            data=login_data
        )
        
        if response.status_code == 200:
            result = response.json()
            self.token = result['access_token']
            print(f"✅ 记住我登录成功")
            print(f"Token: {self.token[:50]}...")
            
            # 2. 解码JWT查看过期时间
            import base64
            try:
                # JWT解码（简化版）
                payload_part = self.token.split('.')[1]
                # 添加padding
                payload_part += '=' * (4 - len(payload_part) % 4)
                payload = json.loads(base64.urlsafe_b64decode(payload_part))
                exp_time = payload.get('exp', 0)
                current_time = int(time.time())
                expires_in_hours = (exp_time - current_time) / 3600
                
                print(f"✅ Token将在 {expires_in_hours:.1f} 小时后过期")
                
                if expires_in_hours > 24:
                    print("✅ 记住我功能生效（超过24小时）")
                else:
                    print("❌ 记住我功能未生效（少于24小时）")
                    
            except Exception as e:
                print(f"❌ 解码JWT失败: {e}")
                
        else:
            print(f"❌ 记住我登录失败: {response.status_code}")
            print(response.text)
    
    def test_token_validation(self):
        """测试token验证"""
        print("\n=== 测试Token验证 ===")
        
        if not self.token:
            print("❌ 没有token，先执行登录测试")
            return
            
        headers = {'Authorization': f'Bearer {self.token}'}
        
        # 测试验证接口
        response = requests.get(
            f"{BASE_URL}/api/auth/me",
            headers=headers
        )
        
        if response.status_code == 200:
            user_info = response.json()
            print(f"✅ Token验证成功: {user_info}")
        else:
            print(f"❌ Token验证失败: {response.status_code}")
            print(response.text)

if __name__ == "__main__":
    tester = LoginFlowTester()
    
    print("Wenxi网盘登录流程测试")
    print("=" * 50)
    
    # 测试普通登录
    tester.test_normal_login()
    
    # 测试记住我登录
    tester.test_remember_me_login()
    
    # 测试token验证
    tester.test_token_validation()
    
    print("\n" + "=" * 50)
    print("测试完成")
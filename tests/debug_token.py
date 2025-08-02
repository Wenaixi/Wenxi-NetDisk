#!/usr/bin/env python3
"""
Wenxi网盘 - JWT令牌调试工具
作者：Wenxi
功能：深度调试JWT令牌验证流程
"""

import requests
import jwt
import time
import json
from datetime import datetime, timezone

# 配置
BASE_URL = "http://localhost:3008"
USERNAME = "debug_user"
PASSWORD = "debugpass123"

# JWT配置
SECRET_KEY = "wenxi-super-secret-key-change-in-production"
ALGORITHM = "HS256"

def test_token_validation():
    """测试令牌验证流程"""
    print("🧪 Wenxi网盘JWT调试工具")
    print("=" * 50)
    
    # 1. 测试登录获取token
    print("\n1. 测试登录获取token...")
    login_data = {
        'username': USERNAME,
        'password': PASSWORD,
        'client_id': 'remember_me'
    }
    
    try:
        response = requests.post(f"{BASE_URL}/api/auth/login", data=login_data)
        response.raise_for_status()
        
        token_data = response.json()
        token = token_data['access_token']
        print(f"✅ 登录成功，获取token: {token[:50]}...")
        
        # 2. 解码并分析token
        print("\n2. 解码token内容...")
        try:
            decoded = jwt.decode(token, options={"verify_signature": False})
            print(f"📋 Token内容:")
            print(f"   sub: {decoded.get('sub')}")
            print(f"   exp: {decoded.get('exp')} ({datetime.fromtimestamp(decoded.get('exp'), tz=timezone.utc)})")
            print(f"   remember_me: {decoded.get('remember_me')}")
            
            # 3. 验证token签名
            print("\n3. 验证token签名...")
            try:
                verified = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
                print("✅ Token签名验证通过")
            except jwt.InvalidTokenError as e:
                print(f"❌ Token签名验证失败: {e}")
                
        except Exception as e:
            print(f"❌ Token解码失败: {e}")
            
        # 4. 测试/me接口
        print("\n4. 测试/me接口验证...")
        headers = {'Authorization': f'Bearer {token}'}
        
        try:
            me_response = requests.get(f"{BASE_URL}/api/auth/me", headers=headers)
            print(f"   状态码: {me_response.status_code}")
            
            if me_response.status_code == 200:
                user_data = me_response.json()
                print(f"✅ /me接口成功: {user_data}")
            else:
                print(f"❌ /me接口失败: {me_response.text}")
                
        except Exception as e:
            print(f"❌ /me接口请求失败: {e}")
            
        # 5. 检查服务器时间
        print("\n5. 检查服务器时间...")
        try:
            health_response = requests.get(f"{BASE_URL}/health")
            server_time = health_response.json().get('timestamp')
            print(f"   服务器时间: {server_time}")
            print(f"   本地时间: {datetime.now(timezone.utc)}")
        except Exception as e:
            print(f"❌ 获取服务器时间失败: {e}")
            
    except Exception as e:
        print(f"❌ 登录失败: {e}")

if __name__ == "__main__":
    test_token_validation()
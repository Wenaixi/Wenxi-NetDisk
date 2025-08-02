"""
Wenxi网盘 - 认证模块单元测试
作者：Wenxi
功能：测试用户注册、登录、JWT认证
"""

import pytest


class TestAuth:
    """认证功能测试"""

    def test_register_success(self, client):
        """测试用户注册成功"""
        import uuid
        unique_username = f"testuser_{uuid.uuid4().hex[:8]}"
        unique_email = f"test_{uuid.uuid4().hex[:8]}@example.com"
        response = client.post("/api/auth/register", json={
            "username": unique_username,
            "email": unique_email,
            "password": "testpassword123"
        })
        assert response.status_code == 200
        data = response.json()
        assert data["username"] == unique_username
        assert "id" in data

    def test_register_duplicate_username(self, client):
        """测试重复用户名注册失败"""
        import uuid
        duplicate_username = f"duplicate_{uuid.uuid4().hex[:8]}"
        
        # 先注册一个用户
        client.post("/api/auth/register", json={
            "username": duplicate_username,
            "email": f"duplicate_{uuid.uuid4().hex[:8]}@example.com",
            "password": "testpassword123"
        })
        
        # 尝试重复注册
        response = client.post("/api/auth/register", json={
            "username": duplicate_username,
            "email": f"new_{uuid.uuid4().hex[:8]}@example.com",
            "password": "testpassword123"
        })
        assert response.status_code == 400
        assert "用户名已存在" in response.json()["detail"]

    def test_login_success(self, client):
        """测试用户登录成功"""
        import uuid
        login_username = f"logintest_{uuid.uuid4().hex[:8]}"
        
        # 先注册用户
        client.post("/api/auth/register", json={
            "username": login_username,
            "email": f"login_{uuid.uuid4().hex[:8]}@example.com",
            "password": "testpassword123"
        })
        
        # 测试登录
        response = client.post("/api/auth/login", data={
            "username": login_username,
            "password": "testpassword123"
        })
        assert response.status_code == 200
        data = response.json()
        assert "access_token" in data
        assert data["token_type"] == "bearer"

    def test_login_invalid_credentials(self, client):
        """测试登录失败"""
        response = client.post("/api/auth/login", data={
            "username": "nonexistent",
            "password": "wrongpassword"
        })
        assert response.status_code == 401
        assert "用户名或密码错误" in response.json()["detail"]

    def test_get_current_user(self, client):
        """测试获取当前用户信息"""
        import uuid
        current_username = f"currentuser_{uuid.uuid4().hex[:8]}"
        
        # 注册用户并登录
        client.post("/api/auth/register", json={
            "username": current_username,
            "email": f"current_{uuid.uuid4().hex[:8]}@example.com",
            "password": "testpassword123"
        })
        
        login_response = client.post("/api/auth/login", data={
            "username": current_username,
            "password": "testpassword123"
        })
        token = login_response.json()["access_token"]
        
        # 获取当前用户信息
        response = client.get("/api/auth/me", headers={
            "Authorization": f"Bearer {token}"
        })
        assert response.status_code == 200
        data = response.json()
        assert data["username"] == current_username

if __name__ == "__main__":
    pytest.main([__file__])
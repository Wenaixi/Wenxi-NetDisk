"""
Wenxi网盘 - 导航跳转测试
作者：Wenxi
功能：测试前端路由导航功能
"""

import pytest
from fastapi.testclient import TestClient


class TestNavigation:
    """导航跳转测试"""

    def test_register_login_navigation(self, client):
        """测试注册后切换到登录模式"""
        import uuid
        unique_username = f"navtest_{uuid.uuid4().hex[:8]}"
        unique_email = f"nav_{uuid.uuid4().hex[:8]}@example.com"
        
        # 注册用户
        response = client.post("/api/auth/register", json={
            "username": unique_username,
            "email": unique_email,
            "password": "testpassword123"
        })
        assert response.status_code == 200
        
        # 登录用户
        response = client.post("/api/auth/login", data={
            "username": unique_username,
            "password": "testpassword123"
        })
        assert response.status_code == 200
        assert "access_token" in response.json()

    def test_login_redirect_flow(self, client):
        """测试登录流程完整性"""
        import uuid
        test_username = f"loginflow_{uuid.uuid4().hex[:8]}"
        
        # 先注册用户
        client.post("/api/auth/register", json={
            "username": test_username,
            "email": f"flow_{uuid.uuid4().hex[:8]}@example.com",
            "password": "testpassword123"
        })
        
        # 测试登录获取token
        response = client.post("/api/auth/login", data={
            "username": test_username,
            "password": "testpassword123"
        })
        
        token = response.json()["access_token"]
        
        # 使用token访问受保护路由
        response = client.get("/api/files/list", headers={
            "Authorization": f"Bearer {token}"
        })
        assert response.status_code == 200


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
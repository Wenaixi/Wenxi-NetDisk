"""
Wenxi网盘 - Dashboard认证测试
作者：Wenxi
功能：测试Dashboard页面的认证和权限验证
"""

import pytest

class TestDashboardAuth:
    """Dashboard认证测试类"""

    def test_dashboard_api_call_with_auth(self, client, auth_headers):
        """测试Dashboard在有认证状态下正确调用API"""
        # 使用已认证的客户端测试文件列表API
        response = client.get("/api/files/list", headers=auth_headers)
        assert response.status_code == 200

    def test_dashboard_401_without_auth(self, client):
        """测试Dashboard在未认证状态下访问API返回401"""
        response = client.get("/api/files/list")
        assert response.status_code == 401

if __name__ == "__main__":
    pytest.main([__file__, "-v"])
"""
Wenxi网盘 - 文件上传认证测试
作者：Wenxi
功能：测试文件上传功能的认证和权限验证
"""

import pytest
import tempfile
import os
from fastapi.testclient import TestClient


class TestFileUploadAuth:
    """文件上传认证测试类"""

    def test_upload_without_auth(self, client):
        """测试未认证状态下上传文件应返回401"""
        # 创建临时测试文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                response = client.post(
                    "/api/files/upload",
                    files={"file": ("test.txt", f, "text/plain")}
                    # 不添加认证头部
                )
            
            assert response.status_code == 401
            # 未认证场景返回的是FastAPI默认消息
            assert "Not authenticated" in response.json()["detail"]
        finally:
            os.unlink(tmp_path)

    def test_upload_with_valid_auth(self, client, auth_headers):
        """测试使用有效token上传文件成功"""
        # 创建临时测试文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                response = client.post(
                    "/api/files/upload",
                    files={"file": ("test.txt", f, "text/plain")},
                    headers=auth_headers
                )
            
            assert response.status_code == 200
            data = response.json()
            assert data["filename"] == "test.txt"
            assert data["file_size"] > 0
            assert "download_url" in data
        finally:
            os.unlink(tmp_path)

    def test_upload_with_invalid_auth(self, client):
        """测试使用无效token上传文件"""
        # 创建临时测试文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                response = client.post(
                    "/api/files/upload",
                    files={"file": ("test.txt", f, "text/plain")},
                    headers={"Authorization": "Bearer invalid_token"}
                )
            
            assert response.status_code == 401
            assert "无法验证凭据" in response.json()["detail"]
        finally:
            os.unlink(tmp_path)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
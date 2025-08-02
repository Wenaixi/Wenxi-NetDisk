"""
Wenxi网盘 - 文件管理单元测试
作者：Wenxi
功能：测试文件上传、下载、分享、删除等核心功能
"""

import os
import tempfile
import pytest

class TestFiles:
    """文件功能测试"""

    @pytest.fixture(autouse=True)
    def setup_method(self, client):
        """测试前准备"""
        import uuid
        self.client = client
        unique_username = f"fileuser_{uuid.uuid4().hex[:8]}"
        unique_email = f"file_{uuid.uuid4().hex[:8]}@example.com"
        
        # 注册用户并登录
        client.post("/api/auth/register", json={
            "username": unique_username,
            "email": unique_email,
            "password": "testpassword123"
        })
        
        login_response = client.post("/api/auth/login", data={
            "username": unique_username,
            "password": "testpassword123"
        })
        self.token = login_response.json()["access_token"]
        self.headers = {"Authorization": f"Bearer {self.token}"}

    def test_upload_file(self):
        """测试文件上传"""
        # 创建临时测试文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                response = self.client.post(
                    "/api/files/upload",
                    files={"file": ("test.txt", f, "text/plain")},
                    headers=self.headers
                )
            
            assert response.status_code == 200
            data = response.json()
            assert data["filename"] == "test.txt"
            assert data["file_size"] > 0
            assert "download_url" in data
        finally:
            os.unlink(tmp_path)

    def test_list_files(self):
        """测试获取文件列表"""
        # 先上传一个文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                self.client.post(
                    "/api/files/upload",
                    files={"file": ("test.txt", f, "text/plain")},
                    headers=self.headers
                )
            
            # 获取文件列表
            response = self.client.get("/api/files/list", headers=self.headers)
            assert response.status_code == 200
            data = response.json()
            assert len(data) >= 1
            assert any(f["original_filename"] == "test.txt" for f in data)
        finally:
            os.unlink(tmp_path)

    def test_delete_file(self):
        """测试删除文件"""
        # 上传文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False) as tmp:
            tmp.write("测试文件内容")
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                upload_response = self.client.post(
                    "/api/files/upload",
                    files={"file": ("delete_test.txt", f, "text/plain")},
                    headers=self.headers
                )
            
            file_id = upload_response.json()["id"]
            
            # 删除文件
            response = self.client.delete(f"/api/files/{file_id}", headers=self.headers)
            assert response.status_code == 200
            assert response.json()["message"] == "文件删除成功"
            
            # 验证文件已被删除
            list_response = self.client.get("/api/files/list", headers=self.headers)
            files = list_response.json()
            assert not any(f["id"] == file_id for f in files)
        finally:
            os.unlink(tmp_path)

    def test_download_file(self):
        """测试文件下载"""
        test_content = "下载测试内容"
        
        # 创建临时测试文件
        with tempfile.NamedTemporaryFile(mode='w', suffix='.txt', delete=False, encoding='utf-8') as tmp:
            tmp.write(test_content)
            tmp_path = tmp.name

        try:
            with open(tmp_path, 'rb') as f:
                upload_response = self.client.post(
                    "/api/files/upload",
                    files={"file": ("download_test.txt", f, "text/plain")},
                    headers=self.headers
                )
            
            file_id = upload_response.json()["id"]
            
            # 下载文件
            response = self.client.get(f"/api/files/download/{file_id}", headers=self.headers)
            assert response.status_code == 200
            
            # 解码响应内容进行比较
            downloaded_content = response.content.decode('utf-8')
            assert downloaded_content == test_content
        finally:
            # 确保文件句柄已关闭
            try:
                os.unlink(tmp_path)
            except PermissionError:
                # Windows下可能暂时无法删除，稍后再试
                pass

    def test_unauthorized_access(self):
        """测试未授权访问"""
        response = self.client.get("/api/files/list")
        assert response.status_code == 401

if __name__ == "__main__":
    pytest.main([__file__])
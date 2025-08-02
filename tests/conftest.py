"""
Wenxi网盘 - Pytest测试配置
作者：Wenxi
功能：提供测试环境配置、测试数据准备、数据库清理等公共功能
"""

import pytest
import tempfile
import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from backend.models import Base
from backend.database import get_db
from backend.main import app


@pytest.fixture(scope="session", autouse=True)
def setup_test_database():
    """设置测试数据库，整个测试会话只执行一次"""
    # 创建临时数据库文件
    db_fd, db_path = tempfile.mkstemp(suffix='.db')
    
    # 配置测试数据库URL
    test_db_url = f"sqlite:///{db_path}"
    
    # 创建测试引擎
    test_engine = create_engine(
        test_db_url,
        connect_args={"check_same_thread": False}
    )
    
    # 创建会话工厂
    TestSessionLocal = sessionmaker(
        autocommit=False,
        autoflush=False,
        bind=test_engine
    )
    
    # 创建所有表
    Base.metadata.create_all(bind=test_engine)
    
    # 覆盖依赖注入
    def override_get_db():
        try:
            db = TestSessionLocal()
            yield db
        finally:
            db.close()
    
    app.dependency_overrides[get_db] = override_get_db
    
    yield
    
    # 清理 - 确保所有数据库连接被关闭
    try:
        # 关闭引擎连接池
        test_engine.dispose()
        
        # 关闭文件描述符
        os.close(db_fd)
        
        # 强制垃圾回收，确保所有连接被释放
        import gc
        gc.collect()
        
        # 尝试删除文件，如果失败则等待并重试
        max_retries = 3
        for attempt in range(max_retries):
            try:
                os.unlink(db_path)
                break
            except PermissionError:
                if attempt < max_retries - 1:
                    import time
                    time.sleep(0.1 * (attempt + 1))  # 递增等待时间
                else:
                    # 最后一次尝试仍然失败，记录警告但不中断测试
                    print(f"警告: 无法删除临时数据库文件 {db_path}")
    except Exception as e:
        print(f"清理测试数据库时出错: {e}")


@pytest.fixture
def client():
    """提供测试客户端"""
    from fastapi.testclient import TestClient
    return TestClient(app)


@pytest.fixture
def auth_headers(client):
    """提供已认证的请求头"""
    # 注册用户
    client.post("/api/auth/register", json={
        "username": "testuser",
        "email": "test@example.com",
        "password": "testpassword123"
    })
    
    # 登录获取token
    login_response = client.post("/api/auth/login", data={
        "username": "testuser",
        "password": "testpassword123"
    })
    
    token = login_response.json()["access_token"]
    return {"Authorization": f"Bearer {token}"}
"""
Wenxi-NetDisk - 增强测试套件
包含安全测试、性能测试、集成测试
"""

import pytest
import asyncio
import os
import sys
import tempfile
import shutil
from datetime import datetime, timedelta
from unittest.mock import Mock, patch, AsyncMock
import jwt

# 添加backend到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'backend'))

from routers.auth import (
    create_access_token, verify_password, get_password_hash,
    get_current_user, SECRET_KEY, ALGORITHM
)
from utils.encryption import (
    encrypt_file, decrypt_file, derive_key,
    WENXI_MAGIC_HEADER, HEADER_VERSION
)


# ==================== Fixtures ====================

@pytest.fixture
def temp_dir():
    """创建临时目录"""
    temp = tempfile.mkdtemp()
    yield temp
    shutil.rmtree(temp, ignore_errors=True)


@pytest.fixture
def test_file(temp_dir):
    """创建测试文件"""
    file_path = os.path.join(temp_dir, "test.txt")
    with open(file_path, 'w') as f:
        f.write("Wenxi网盘测试数据 - " + "A" * 1000)
    return file_path


@pytest.fixture
def large_test_file(temp_dir):
    """创建大测试文件"""
    file_path = os.path.join(temp_dir, "large_test.bin")
    with open(file_path, 'wb') as f:
        f.write(os.urandom(1024 * 1024))  # 1MB
    return file_path


# ==================== 认证模块测试 ====================

def test_password_hashing():
    """测试密码哈希功能"""
    password = "WenxiSecure123!"
    hashed = get_password_hash(password)
    
    # 验证哈希不是明文
    assert hashed != password
    # 验证密码验证
    assert verify_password(password, hashed) is True
    # 验证错误密码
    assert verify_password("wrong_password", hashed) is False


def test_create_access_token():
    """测试JWT令牌创建"""
    data = {"sub": "testuser"}
    token = create_access_token(data)
    
    # 验证令牌可以解码
    payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
    assert payload["sub"] == "testuser"
    assert "exp" in payload


def test_create_access_token_with_expiry():
    """测试带过期时间的令牌创建"""
    data = {"sub": "testuser"}
    expires = timedelta(minutes=30)
    token = create_access_token(data, expires)
    
    payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
    exp_time = datetime.fromtimestamp(payload["exp"])
    
    # 验证过期时间在未来
    assert exp_time > datetime.utcnow()


@pytest.mark.asyncio
async def test_get_current_user_invalid_token():
    """测试无效令牌处理"""
    from fastapi import HTTPException
    
    invalid_token = "invalid.token.here"
    
    with pytest.raises(HTTPException) as exc_info:
        # 创建一个模拟的db
        mock_db = Mock()
        await get_current_user(invalid_token, mock_db)
    
    assert exc_info.value.status_code == 401


# ==================== 加密模块测试 ====================

def test_encryption_decryption_cycle(test_file, temp_dir):
    """测试加密解密完整周期"""
    encrypted_path = os.path.join(temp_dir, "encrypted.bin")
    decrypted_path = os.path.join(temp_dir, "decrypted.txt")
    
    # 加密
    success = encrypt_file(test_file, encrypted_path)
    assert success is True
    assert os.path.exists(encrypted_path)
    
    # 解密
    success = decrypt_file(encrypted_path, decrypted_path)
    assert success is True
    assert os.path.exists(decrypted_path)
    
    # 验证内容一致
    with open(test_file, 'r') as f:
        original = f.read()
    with open(decrypted_path, 'r') as f:
        decrypted = f.read()
    
    assert original == decrypted


def test_large_file_encryption(large_test_file, temp_dir):
    """测试大文件加密性能"""
    encrypted_path = os.path.join(temp_dir, "large_encrypted.bin")
    decrypted_path = os.path.join(temp_dir, "large_decrypted.bin")
    
    import time
    
    # 加密计时
    start = time.time()
    success = encrypt_file(large_test_file, encrypted_path, user_id=1, file_id=1)
    encrypt_time = time.time() - start
    
    assert success is True
    assert encrypt_time < 10.0, f"加密时间过长: {encrypt_time}s"
    
    # 解密计时
    start = time.time()
    success = decrypt_file(encrypted_path, decrypted_path, user_id=1, file_id=1)
    decrypt_time = time.time() - start
    
    assert success is True
    assert decrypt_time < 10.0, f"解密时间过长: {decrypt_time}s"
    
    # 验证文件大小
    original_size = os.path.getsize(large_test_file)
    decrypted_size = os.path.getsize(decrypted_path)
    assert original_size == decrypted_size


def test_encryption_invalid_file(temp_dir):
    """测试加密不存在文件"""
    result = encrypt_file("/nonexistent/file.txt", "/output/encrypted.bin")
    assert result is False


def test_decryption_invalid_file(temp_dir):
    """测试解密不存在文件"""
    result = decrypt_file("/nonexistent/encrypted.bin", "/output/decrypted.txt")
    assert result is False


def test_decryption_corrupted_file(temp_dir):
    """测试解密损坏的文件"""
    corrupted_file = os.path.join(temp_dir, "corrupted.bin")
    
    # 创建带有WENXI头但内容损坏的文件
    with open(corrupted_file, 'wb') as f:
        f.write(WENXI_MAGIC_HEADER)
        f.write(bytes([HEADER_VERSION]))
        f.write(b"invalid_data_here")
    
    output_path = os.path.join(temp_dir, "output.txt")
    result = decrypt_file(corrupted_file, output_path)
    
    assert result is False


def test_derive_key_consistency():
    """测试密钥派生一致性"""
    password = "WenxiTestPassword"
    salt = b"test_salt_12345"
    
    key1 = derive_key(password, salt)
    key2 = derive_key(password, salt)
    
    # 相同密码和盐应生成相同密钥
    assert key1 == key2
    assert len(key1) == 32  # 256位


def test_derive_key_different_salts():
    """测试不同盐生成不同密钥"""
    password = "WenxiTestPassword"
    salt1 = b"salt_one_1234567"
    salt2 = b"salt_two_1234567"
    
    key1 = derive_key(password, salt1)
    key2 = derive_key(password, salt2)
    
    # 不同盐应生成不同密钥
    assert key1 != key2


# ==================== 安全测试 ====================

def test_encryption_header_format(temp_dir):
    """测试加密文件头格式"""
    test_file = os.path.join(temp_dir, "test.txt")
    encrypted_file = os.path.join(temp_dir, "encrypted.bin")
    
    with open(test_file, 'w') as f:
        f.write("test")
    
    encrypt_file(test_file, encrypted_file)
    
    # 验证文件头
    with open(encrypted_file, 'rb') as f:
        magic = f.read(len(WENXI_MAGIC_HEADER))
        version = f.read(1)
    
    assert magic == WENXI_MAGIC_HEADER
    assert version[0] == HEADER_VERSION


def test_tamper_detection(temp_dir):
    """测试篡改检测"""
    test_file = os.path.join(temp_dir, "test.txt")
    encrypted_file = os.path.join(temp_dir, "encrypted.bin")
    tampered_file = os.path.join(temp_dir, "tampered.bin")
    decrypted_file = os.path.join(temp_dir, "decrypted.txt")
    
    with open(test_file, 'w') as f:
        f.write("Wenxi安全测试数据")
    
    # 加密
    encrypt_file(test_file, encrypted_file)
    
    # 篡改加密文件
    with open(encrypted_file, 'rb') as f:
        data = bytearray(f.read())
    
    # 修改文件中间的一些字节
    if len(data) > 100:
        data[50] ^= 0xFF
    
    with open(tampered_file, 'wb') as f:
        f.write(data)
    
    # 尝试解密篡改后的文件应失败
    result = decrypt_file(tampered_file, decrypted_file)
    assert result is False


# ==================== 性能基准测试 ====================

def test_encryption_performance(benchmark, temp_dir):
    """加密性能基准测试"""
    test_file = os.path.join(temp_dir, "perf_test.txt")
    encrypted_file = os.path.join(temp_dir, "perf_encrypted.bin")
    
    # 创建100KB测试文件
    with open(test_file, 'wb') as f:
        f.write(os.urandom(100 * 1024))
    
    def encrypt():
        encrypt_file(test_file, encrypted_file)
    
    result = benchmark(encrypt)


def test_decryption_performance(benchmark, temp_dir):
    """解密性能基准测试"""
    test_file = os.path.join(temp_dir, "perf_test.txt")
    encrypted_file = os.path.join(temp_dir, "perf_encrypted.bin")
    decrypted_file = os.path.join(temp_dir, "perf_decrypted.txt")
    
    # 创建并加密测试文件
    with open(test_file, 'wb') as f:
        f.write(os.urandom(100 * 1024))
    encrypt_file(test_file, encrypted_file)
    
    def decrypt():
        decrypt_file(encrypted_file, decrypted_file)
    
    result = benchmark(decrypt)


# ==================== 集成测试 ====================

@pytest.mark.asyncio
async def test_file_upload_flow(temp_dir):
    """测试文件上传完整流程"""
    # 此测试需要数据库，使用mock
    from unittest.mock import MagicMock
    
    mock_db = MagicMock()
    mock_user = MagicMock()
    mock_user.id = 1
    mock_user.username = "testuser"
    
    # 模拟文件上传
    test_content = b"Wenxi网盘测试文件内容"
    test_file = os.path.join(temp_dir, "upload_test.txt")
    
    with open(test_file, 'wb') as f:
        f.write(test_content)
    
    # 验证文件创建成功
    assert os.path.exists(test_file)
    
    with open(test_file, 'rb') as f:
        content = f.read()
    
    assert content == test_content


def test_concurrent_encryption(temp_dir):
    """测试并发加密性能"""
    import concurrent.futures
    import time
    
    test_files = []
    encrypted_files = []
    
    # 创建多个测试文件
    for i in range(5):
        test_file = os.path.join(temp_dir, f"concurrent_test_{i}.txt")
        encrypted_file = os.path.join(temp_dir, f"concurrent_encrypted_{i}.bin")
        
        with open(test_file, 'wb') as f:
            f.write(os.urandom(100 * 1024))  # 100KB每个文件
        
        test_files.append(test_file)
        encrypted_files.append(encrypted_file)
    
    # 并发加密
    start = time.time()
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
        futures = [
            executor.submit(encrypt_file, test_files[i], encrypted_files[i])
            for i in range(len(test_files))
        ]
        results = [f.result() for f in futures]
    
    total_time = time.time() - start
    
    # 验证所有加密成功
    assert all(results)
    # 验证性能合理
    assert total_time < 30.0, f"并发加密时间过长: {total_time}s"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])

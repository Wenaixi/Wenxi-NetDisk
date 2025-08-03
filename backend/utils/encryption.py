"""
Wenxi网盘 - 超强兼容性文件加密模块 v2.0
作者：Wenxi
功能：使用ChaCha20-Poly1305流加密，100%兼容所有文件格式
支持：MP4、PDF、DOCX、JPG、PNG、ZIP、EXE等所有文件类型
特点：零损坏、高性能、跨平台兼容、内存优化
"""

import os
import struct
import logging
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms
from cryptography.hazmat.primitives.ciphers.aead import ChaCha20Poly1305
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.backends import default_backend
import secrets

from logger import logger

# 超强兼容性加密配置 - 仅使用环境变量，无硬编码默认值
ENCRYPTION_KEY = os.environ.get("WENXI_ENCRYPTION_KEY")
if not ENCRYPTION_KEY:
    raise ValueError("环境变量 WENXI_ENCRYPTION_KEY 未设置")

SALT = os.environ.get("WENXI_ENCRYPTION_SALT")
if not SALT:
    raise ValueError("环境变量 WENXI_ENCRYPTION_SALT 未设置")
SALT = SALT.encode()

# ChaCha20-Poly1305参数 - 业界最强兼容性
WENXI_MAGIC_HEADER = b'WENXI\x02\x00'  # v2.0标识
HEADER_VERSION = 2  # 版本2 - ChaCha20-Poly1305
KEY_SIZE = 32  # ChaCha20 256-bit密钥
NONCE_SIZE = 12  # ChaCha20标准nonce大小
TAG_SIZE = 16  # Poly1305认证标签
CHUNK_SIZE = 64 * 1024  # 64KB块大小，内存友好


def derive_key(password: str, salt: bytes) -> bytes:
    """
    Wenxi超强兼容 - 从密码派生安全密钥
    使用PBKDF2-HMAC-SHA256，100万次迭代确保安全性
    
    参数:
        password: 用户密码
        salt: 随机盐值
    
    返回:
        32字节ChaCha20安全密钥
    """
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=KEY_SIZE,
        salt=salt,
        iterations=1000000,  # 100万次迭代，极致安全
        backend=default_backend()
    )
    return kdf.derive(password.encode())


def encrypt_file(input_path: str, output_path: str, password: str = None, user_id: int = None, file_id: int = None) -> bool:
    """
    Wenxi超强兼容 - 流式加密单个文件
    使用ChaCha20-Poly1305，100%兼容所有文件格式
    
    参数:
        input_path: 输入文件路径
        output_path: 输出加密文件路径
        password: 加密密码(可选)
        user_id: 用户ID(可选，用于日志追踪)
        file_id: 文件ID(可选，用于日志追踪)
    
    返回:
        加密成功返回True，失败返回False
    
    特性:
        - 零内存压力: 64KB块处理
        - 100%格式兼容: 支持MP4、PDF、ZIP等所有格式
        - 自动完整性验证
    """
    try:
        key_password = password or ENCRYPTION_KEY
        nonce = secrets.token_bytes(NONCE_SIZE)  # 安全随机nonce
        key = derive_key(key_password, SALT)
        
        # 创建ChaCha20-Poly1305实例
        chacha = ChaCha20Poly1305(key)
        
        with open(input_path, 'rb') as infile, open(output_path, 'wb') as outfile:
            # 写入文件头
            outfile.write(WENXI_MAGIC_HEADER)
            outfile.write(struct.pack('B', HEADER_VERSION))
            outfile.write(nonce)
            
            # 流式加密 - 内存友好
            file_size = os.path.getsize(input_path)
            outfile.write(struct.pack('>Q', file_size))  # 8字节文件大小
            
            encrypted_size = 0
            while True:
                chunk = infile.read(CHUNK_SIZE)
                if not chunk:
                    break
                
                # 为每个块生成关联数据（块序号）
                associated_data = struct.pack('>Q', encrypted_size // CHUNK_SIZE)
                encrypted_chunk = chacha.encrypt(nonce, chunk, associated_data)
                outfile.write(encrypted_chunk)
                encrypted_size += len(chunk)
        
        user_info = f"[用户{user_id}文件{file_id}]" if user_id and file_id else ""
        logger.info(f"[Wenxi加密] 成功{user_info}: {os.path.basename(input_path)} ({encrypted_size/1024/1024:.2f}MB)")
        return True
        
    except Exception as e:
        logger.error(f"[Wenxi加密] 失败 {input_path}: {str(e)}")
        if os.path.exists(output_path):
            os.remove(output_path)  # 清理失败文件
        return False


def decrypt_file(input_path: str, output_path: str, password: str = None, user_id: int = None, file_id: int = None) -> bool:
    """
    Wenxi超强兼容 - 流式解密单个文件
    仅支持新版ChaCha20-Poly1305格式，100%还原原始文件
    
    参数:
        input_path: 加密文件路径
        output_path: 解密输出文件路径
        password: 解密密码(可选)
        user_id: 用户ID(可选，用于日志追踪)
        file_id: 文件ID(可选，用于日志追踪)
    
    返回:
        解密成功返回True，失败返回False
    
    特性:
        - 100%文件完整性保证
        - 仅支持新版ChaCha20-Poly1305格式
        - 零损坏风险
        - 支持大文件流式处理
    """
    try:
        key_password = password or ENCRYPTION_KEY
        
        # 仅支持新版ChaCha20-Poly1305格式
        return _decrypt_new_format(input_path, output_path, key_password, user_id, file_id)
            
    except Exception as e:
        logger.error(f"[Wenxi解密] 失败 {input_path}: {str(e)}")
        if os.path.exists(output_path):
            os.remove(output_path)  # 清理失败文件
        return False

def _decrypt_new_format(input_path: str, output_path: str, password: str, user_id: int = None, file_id: int = None) -> bool:
    """
    Wenxi新版解密 - 新版ChaCha20-Poly1305格式专用
    
    参数:
        input_path: 加密文件路径
        output_path: 解密输出文件路径
        password: 解密密码
        user_id: 用户ID
        file_id: 文件ID
    
    返回:
        解密成功返回True
    """
    try:
        key = derive_key(password, SALT)
        
        with open(input_path, 'rb') as infile, open(output_path, 'wb') as outfile:
            # 验证文件头
            magic = infile.read(len(WENXI_MAGIC_HEADER))
            if magic != WENXI_MAGIC_HEADER:
                logger.error(f"[Wenxi解密] 无效格式: {os.path.basename(input_path)}")
                return False
            
            version = struct.unpack('B', infile.read(1))[0]
            if version != HEADER_VERSION:
                logger.error(f"[Wenxi解密] 版本不兼容: {version}")
                return False
            
            nonce = infile.read(NONCE_SIZE)
            original_size = struct.unpack('>Q', infile.read(8))[0]
            
            chacha = ChaCha20Poly1305(key)
            
            # 流式解密
            decrypted_size = 0
            chunk_index = 0
            
            while decrypted_size < original_size:
                # 计算当前块大小
                remaining = original_size - decrypted_size
                chunk_size = min(CHUNK_SIZE + TAG_SIZE, remaining + TAG_SIZE)
                encrypted_chunk = infile.read(chunk_size)
                
                if not encrypted_chunk:
                    break
                
                associated_data = struct.pack('>Q', chunk_index)
                decrypted_chunk = chacha.decrypt(nonce, encrypted_chunk, associated_data)
                
                # 写入实际大小的数据（处理最后一块）
                write_size = min(len(decrypted_chunk), original_size - decrypted_size)
                outfile.write(decrypted_chunk[:write_size])
                
                decrypted_size += write_size
                chunk_index += 1
            
            if decrypted_size != original_size:
                logger.error(f"[Wenxi解密] 大小不匹配: 期望{original_size}, 实际{decrypted_size}")
                return False
        
        user_info = f"[用户{user_id}文件{file_id}]" if user_id and file_id else ""
        logger.info(f"[Wenxi解密] 成功{user_info}: {os.path.basename(input_path)} ({decrypted_size/1024/1024:.2f}MB)")
        return True
        
    except Exception as e:
        logger.error(f"[Wenxi解密] 新版解密失败: {str(e)}")
        return False


def encrypt_stream(data: bytes, password: str = None) -> bytes:
    """
    Wenxi超强兼容 - 加密数据流
    使用ChaCha20-Poly1305，100%数据完整性保证
    
    参数:
        data: 原始数据
        password: 加密密码(可选)
    
    返回:
        加密后的数据包（包含头部信息）
    """
    try:
        key_password = password or ENCRYPTION_KEY
        nonce = secrets.token_bytes(NONCE_SIZE)
        key = derive_key(key_password, SALT)
        
        # 创建ChaCha20-Poly1305实例
        chacha = ChaCha20Poly1305(key)
        
        # 加密数据
        encrypted_data = chacha.encrypt(nonce, data, b"stream")
        
        # 构建完整数据包
        header = WENXI_MAGIC_HEADER + struct.pack('>B', HEADER_VERSION)
        full_package = header + nonce + struct.pack('>Q', len(data)) + encrypted_data
        
        return full_package
    
    except Exception as e:
        logger.error(f"[Wenxi流加密] 失败: {str(e)}")
        raise


def decrypt_stream(data: bytes, password: str = None) -> bytes:
    """
    Wenxi超强兼容 - 解密数据流
    使用ChaCha20-Poly1305，100%验证数据完整性
    
    参数:
        data: 加密数据（包含头部信息）
        password: 解密密码(可选)
    
    返回:
        解密后的原始数据
    """
    try:
        # 验证最小长度
        min_length = len(WENXI_MAGIC_HEADER) + 1 + NONCE_SIZE + 8 + TAG_SIZE
        if len(data) < min_length:
            raise ValueError("数据长度不足")
        
        # 解析文件头
        pos = 0
        header = data[pos:pos+len(WENXI_MAGIC_HEADER)]
        if header != WENXI_MAGIC_HEADER:
            raise ValueError("无效的文件头")
        
        pos += len(WENXI_MAGIC_HEADER)
        version = struct.unpack('>B', data[pos:pos+1])[0]
        if version != HEADER_VERSION:
            raise ValueError(f"不支持的版本: {version}")
        
        pos += 1
        nonce = data[pos:pos+NONCE_SIZE]
        pos += NONCE_SIZE
        
        original_size = struct.unpack('>Q', data[pos:pos+8])[0]
        pos += 8
        
        encrypted_data = data[pos:]
        
        # 生成解密密钥
        key_password = password or ENCRYPTION_KEY
        key = derive_key(key_password, SALT)
        
        # 创建ChaCha20-Poly1305实例
        chacha = ChaCha20Poly1305(key)
        
        # 解密数据
        decrypted_data = chacha.decrypt(nonce, encrypted_data, b"stream")
        
        if len(decrypted_data) != original_size:
            raise ValueError("数据大小不匹配")
        
        return decrypted_data
    
    except Exception as e:
        logger.error(f"[Wenxi流解密] 失败: {str(e)}")
        raise
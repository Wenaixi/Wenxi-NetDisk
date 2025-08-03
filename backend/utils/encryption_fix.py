"""
Wenxi网盘 - 加密格式兼容性修复工具
作者：Wenxi
功能：修复旧版加密文件的格式兼容性问题，支持向后兼容
"""

import os
import struct
import logging
from utils.encryption import encrypt_file, decrypt_file, WENXI_MAGIC_HEADER, HEADER_VERSION
from logger import logger

def check_file_format(file_path):
    """
    Wenxi兼容性检查 - 验证文件格式
    
    参数:
        file_path: 文件路径
    
    返回:
        文件格式信息字符串
    """
    try:
        with open(file_path, 'rb') as f:
            header = f.read(8)  # 读取8字节检查魔数
            
        if len(header) < 8:
            return "文件过小，可能已损坏"
            
        magic = header[:6]
        if magic == WENXI_MAGIC_HEADER:
            version = struct.unpack('B', header[6:7])[0] if len(header) >= 7 else 0
            return f"Wenxi格式 v{version}.0 - ChaCha20-Poly1305"
        elif header.startswith(b'\x00\x00\x00\x20'):
            return "旧版AES格式 - 需要升级"
        elif header.startswith(b'PK\x03\x04'):
            return "ZIP格式 - 可能未加密"
        else:
            # 检查是否是旧版AES加密文件（无魔数）
            return "未知格式 - 可能是旧版加密"
            
    except Exception as e:
        return f"检查失败: {e}"

def upgrade_old_format(old_file_path, new_file_path, password=None):
    """
    Wenxi格式升级 - 将旧版AES格式升级为ChaCha20-Poly1305格式
    
    参数:
        old_file_path: 旧格式文件路径
        new_file_path: 新格式文件路径
        password: 密码
    
    返回:
        升级成功返回True
    """
    try:
        # 创建临时解密文件
        temp_decrypt_path = old_file_path + ".temp_decrypt"
        
        # 尝试使用旧版AES解密（向后兼容）
        from utils.old_encryption import decrypt_file_old_format
        
        logger.info(f"[Wenxi升级] 开始升级文件: {os.path.basename(old_file_path)}")
        
        # 解密旧格式文件
        decrypt_success = decrypt_file_old_format(old_file_path, temp_decrypt_path, password)
        if not decrypt_success:
            logger.error("[Wenxi升级] 旧格式解密失败")
            return False
            
        # 使用新格式重新加密
        encrypt_success = encrypt_file(temp_decrypt_path, new_file_path, password)
        
        # 清理临时文件
        if os.path.exists(temp_decrypt_path):
            os.remove(temp_decrypt_path)
            
        if encrypt_success:
            logger.info("[Wenxi升级] 格式升级成功")
            return True
        else:
            logger.error("[Wenxi升级] 新格式加密失败")
            return False
            
    except Exception as e:
        logger.error(f"[Wenxi升级] 升级过程失败: {e}")
        return False

def batch_upgrade_directory(directory_path, password=None):
    """
    Wenxi批量升级 - 升级目录中的所有旧格式文件
    
    参数:
        directory_path: 目录路径
        password: 密码
    
    返回:
        升级统计信息
    """
    upgrade_stats = {
        "total_files": 0,
        "upgraded": 0,
        "failed": 0,
        "skipped": 0
    }
    
    if not os.path.exists(directory_path):
        logger.error(f"[Wenxi批量升级] 目录不存在: {directory_path}")
        return upgrade_stats
    
    for filename in os.listdir(directory_path):
        file_path = os.path.join(directory_path, filename)
        if not os.path.isfile(file_path):
            continue
            
        upgrade_stats["total_files"] += 1
        
        format_info = check_file_format(file_path)
        logger.info(f"[Wenxi批量升级] {filename}: {format_info}")
        
        if "旧版AES格式" in format_info or "未知格式" in format_info:
            new_path = file_path + ".upgraded"
            success = upgrade_old_format(file_path, new_path, password)
            
            if success:
                # 替换原文件
                os.remove(file_path)
                os.rename(new_path, file_path)
                upgrade_stats["upgraded"] += 1
                logger.info(f"[Wenxi批量升级] 成功升级: {filename}")
            else:
                upgrade_stats["failed"] += 1
                logger.error(f"[Wenxi批量升级] 升级失败: {filename}")
        else:
            upgrade_stats["skipped"] += 1
            logger.info(f"[Wenxi批量升级] 跳过: {filename}")
    
    logger.info(f"[Wenxi批量升级] 完成统计: {upgrade_stats}")
    return upgrade_stats

def validate_encryption_integrity():
    """
    Wenxi完整性验证 - 验证新加密格式的正确性
    
    返回:
        验证结果
    """
    import tempfile
    import os
    
    test_content = b"Wenxi Network Disk integrity test content - verify ChaCha20-Poly1305 encryption decryption" * 1000
    
    try:
        # 创建测试文件
        with tempfile.NamedTemporaryFile(delete=False) as temp_original:
            temp_original.write(test_content)
            original_path = temp_original.name
            
        with tempfile.NamedTemporaryFile(delete=False) as temp_encrypted:
            encrypted_path = temp_encrypted.name
            
        with tempfile.NamedTemporaryFile(delete=False) as temp_decrypted:
            decrypted_path = temp_decrypted.name
            
        try:
            # 测试加密
            encrypt_success = encrypt_file(original_path, encrypted_path)
            if not encrypt_success:
                return {"success": False, "error": "加密失败"}
            
            # 测试解密
            decrypt_success = decrypt_file(encrypted_path, decrypted_path)
            if not decrypt_success:
                return {"success": False, "error": "解密失败"}
            
            # 验证内容一致性
            with open(decrypted_path, 'rb') as f:
                decrypted_content = f.read()
            
            if decrypted_content != test_content:
                return {"success": False, "error": "内容不一致"}
            
            # 验证文件格式
            format_info = check_file_format(encrypted_path)
            if "ChaCha20-Poly1305" not in format_info:
                return {"success": False, "error": "格式不正确"}
                
            return {"success": True, "format": format_info}
            
        finally:
            # 清理测试文件
            for path in [original_path, encrypted_path, decrypted_path]:
                try:
                    if os.path.exists(path):
                        os.remove(path)
                except:
                    pass
                    
    except Exception as e:
        return {"success": False, "error": str(e)}

if __name__ == "__main__":
    # 运行完整性验证
    result = validate_encryption_integrity()
    print(f"Wenxi完整性验证: {result}")
    
    # 检查uploads目录
    uploads_dir = os.path.join(os.path.dirname(__file__), '..', 'uploads')
    if os.path.exists(uploads_dir):
        print("\nWenxi - 检查uploads目录格式...")
        stats = batch_upgrade_directory(uploads_dir)
        print(f"升级统计: {stats}")
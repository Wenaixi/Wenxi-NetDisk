"""
Wenxi网盘 - 文件版本控制系统
作者：Wenxi
功能：支持文件版本历史、回滚、版本对比
"""

import os
import uuid
import hashlib
from datetime import datetime
from typing import List, Optional, Dict
from sqlalchemy import Column, Integer, String, DateTime, ForeignKey, Text, BigInteger
from sqlalchemy.orm import relationship
from pydantic import BaseModel
from fastapi import APIRouter, Depends, HTTPException, UploadFile, File
from sqlalchemy.orm import Session

from logger import logger
from database import get_db, Base
from routers.auth import get_current_user
from models import User, File as FileModel


# ==================== 数据库模型 ====================

class FileVersion(Base):
    """文件版本模型"""
    __tablename__ = "file_versions"
    
    id = Column(Integer, primary_key=True, index=True)
    file_id = Column(Integer, ForeignKey("files.id"), nullable=False, index=True)
    version_number = Column(Integer, nullable=False)
    storage_path = Column(String(500), nullable=False)  # 实际存储路径
    file_hash = Column(String(64), nullable=False)  # SHA256哈希
    file_size = Column(BigInteger, nullable=False)
    created_by = Column(Integer, ForeignKey("users.id"), nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)
    comment = Column(Text, nullable=True)  # 版本备注
    is_deleted = Column(Integer, default=0)  # 软删除标记
    
    # 关系
    file = relationship("File", back_populates="versions")
    creator = relationship("User")


# 添加到File模型中
# versions = relationship("FileVersion", back_populates="file", order_by="FileVersion.version_number.desc()")


# ==================== Pydantic模型 ====================

class VersionInfo(BaseModel):
    """版本信息响应"""
    id: int
    version_number: int
    file_hash: str
    file_size: int
    created_at: datetime
    comment: Optional[str]
    created_by_username: str
    
    class Config:
        from_attributes = True


class VersionCompareResult(BaseModel):
    """版本对比结果"""
    version_a: int
    version_b: int
    is_identical: bool
    size_diff: int
    hash_match: bool


class VersionRestoreRequest(BaseModel):
    """版本恢复请求"""
    target_version_id: int
    comment: Optional[str] = None


class FileVersionCreate(BaseModel):
    """创建版本请求"""
    comment: Optional[str] = None


# ==================== 版本管理器 ====================

class VersionManager:
    """文件版本管理器"""
    
    def __init__(self, db: Session):
        self.db = db
        self.versions_dir = os.path.join(os.path.dirname(__file__), "..", "uploads", "versions")
        os.makedirs(self.versions_dir, exist_ok=True)
    
    def _get_version_path(self, file_id: int, version_num: int) -> str:
        """生成版本文件存储路径"""
        version_dir = os.path.join(self.versions_dir, str(file_id))
        os.makedirs(version_dir, exist_ok=True)
        return os.path.join(version_dir, f"v{version_num}")
    
    def _calculate_hash(self, file_path: str) -> str:
        """计算文件SHA256哈希"""
        sha256 = hashlib.sha256()
        with open(file_path, 'rb') as f:
            for chunk in iter(lambda: f.read(8192), b""):
                sha256.update(chunk)
        return sha256.hexdigest()
    
    def get_next_version_number(self, file_id: int) -> int:
        """获取下一个版本号"""
        latest = self.db.query(FileVersion).filter(
            FileVersion.file_id == file_id,
            FileVersion.is_deleted == 0
        ).order_by(FileVersion.version_number.desc()).first()
        
        return (latest.version_number + 1) if latest else 1
    
    def create_version(
        self,
        file_id: int,
        source_path: str,
        user_id: int,
        comment: Optional[str] = None
    ) -> FileVersion:
        """
        创建新版本
        
        Args:
            file_id: 文件ID
            source_path: 源文件路径
            user_id: 创建者ID
            comment: 版本备注
        
        Returns:
            创建的版本记录
        """
        # 获取下一个版本号
        version_number = self.get_next_version_number(file_id)
        
        # 生成存储路径
        version_path = self._get_version_path(file_id, version_number)
        
        # 复制文件到版本目录
        import shutil
        shutil.copy2(source_path, version_path)
        
        # 计算哈希
        file_hash = self._calculate_hash(version_path)
        file_size = os.path.getsize(version_path)
        
        # 创建版本记录
        version = FileVersion(
            file_id=file_id,
            version_number=version_number,
            storage_path=version_path,
            file_hash=file_hash,
            file_size=file_size,
            created_by=user_id,
            comment=comment
        )
        
        self.db.add(version)
        self.db.commit()
        self.db.refresh(version)
        
        logger.info(f"创建文件版本: file_id={file_id}, version={version_number}")
        return version
    
    def get_versions(self, file_id: int, include_deleted: bool = False) -> List[FileVersion]:
        """获取文件的所有版本"""
        query = self.db.query(FileVersion).filter(FileVersion.file_id == file_id)
        
        if not include_deleted:
            query = query.filter(FileVersion.is_deleted == 0)
        
        return query.order_by(FileVersion.version_number.desc()).all()
    
    def get_version(self, version_id: int) -> Optional[FileVersion]:
        """获取指定版本"""
        return self.db.query(FileVersion).filter(
            FileVersion.id == version_id,
            FileVersion.is_deleted == 0
        ).first()
    
    def restore_version(
        self,
        version_id: int,
        target_path: str,
        user_id: int
    ) -> bool:
        """
        恢复到指定版本
        
        Args:
            version_id: 版本ID
            target_path: 恢复目标路径
            user_id: 操作用户ID
        
        Returns:
            是否成功
        """
        version = self.get_version(version_id)
        if not version:
            logger.error(f"版本不存在: {version_id}")
            return False
        
        if not os.path.exists(version.storage_path):
            logger.error(f"版本文件不存在: {version.storage_path}")
            return False
        
        try:
            import shutil
            shutil.copy2(version.storage_path, target_path)
            logger.info(f"恢复版本: version_id={version_id}, to={target_path}")
            return True
        except Exception as e:
            logger.error(f"恢复版本失败: {e}")
            return False
    
    def compare_versions(self, version_a_id: int, version_b_id: int) -> VersionCompareResult:
        """
        对比两个版本
        
        Returns:
            对比结果
        """
        version_a = self.get_version(version_a_id)
        version_b = self.get_version(version_b_id)
        
        if not version_a or not version_b:
            raise HTTPException(status_code=404, detail="版本不存在")
        
        return VersionCompareResult(
            version_a=version_a.version_number,
            version_b=version_b.version_number,
            is_identical=version_a.file_hash == version_b.file_hash,
            size_diff=version_a.file_size - version_b.file_size,
            hash_match=version_a.file_hash == version_b.file_hash
        )
    
    def delete_version(self, version_id: int, soft_delete: bool = True) -> bool:
        """
        删除版本
        
        Args:
            version_id: 版本ID
            soft_delete: 是否软删除
        """
        version = self.get_version(version_id)
        if not version:
            return False
        
        if soft_delete:
            version.is_deleted = 1
            self.db.commit()
            logger.info(f"软删除版本: {version_id}")
        else:
            # 硬删除
            if os.path.exists(version.storage_path):
                try:
                    os.remove(version.storage_path)
                except Exception as e:
                    logger.error(f"删除版本文件失败: {e}")
                    return False
            
            self.db.delete(version)
            self.db.commit()
            logger.info(f"硬删除版本: {version_id}")
        
        return True
    
    def cleanup_old_versions(self, file_id: int, keep_count: int = 10) -> int:
        """
        清理旧版本，只保留最近N个
        
        Returns:
            删除的版本数
        """
        versions = self.get_versions(file_id, include_deleted=False)
        
        if len(versions) <= keep_count:
            return 0
        
        to_delete = versions[keep_count:]  # 保留最近的N个
        deleted_count = 0
        
        for version in to_delete:
            if self.delete_version(version.id, soft_delete=False):
                deleted_count += 1
        
        logger.info(f"清理旧版本: file_id={file_id}, deleted={deleted_count}")
        return deleted_count


# ==================== API路由 ====================

router = APIRouter(prefix="/versions", tags=["文件版本"])


@router.get("/file/{file_id}", response_model=List[VersionInfo])
async def list_versions(
    file_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """获取文件版本列表"""
    # 检查文件所有权
    file = db.query(FileModel).filter(
        FileModel.id == file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=404, detail="文件不存在")
    
    manager = VersionManager(db)
    versions = manager.get_versions(file_id)
    
    return [
        VersionInfo(
            id=v.id,
            version_number=v.version_number,
            file_hash=v.file_hash,
            file_size=v.file_size,
            created_at=v.created_at,
            comment=v.comment,
            created_by_username=v.creator.username if v.creator else "未知"
        )
        for v in versions
    ]


@router.post("/file/{file_id}/create")
async def create_new_version(
    file_id: int,
    comment: Optional[str] = None,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """手动创建新版本"""
    file = db.query(FileModel).filter(
        FileModel.id == file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=404, detail="文件不存在")
    
    # 检查文件是否存在
    file_path = os.path.join(os.path.dirname(__file__), "..", file.file_path)
    file_path = os.path.abspath(file_path)
    
    if not os.path.exists(file_path):
        raise HTTPException(status_code=404, detail="文件物理路径不存在")
    
    manager = VersionManager(db)
    version = manager.create_version(file_id, file_path, current_user.id, comment)
    
    return {
        "message": "版本创建成功",
        "version_id": version.id,
        "version_number": version.version_number
    }


@router.post("/restore")
async def restore_version(
    request: VersionRestoreRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """恢复到指定版本"""
    manager = VersionManager(db)
    
    version = manager.get_version(request.target_version_id)
    if not version:
        raise HTTPException(status_code=404, detail="版本不存在")
    
    # 检查文件所有权
    file = db.query(FileModel).filter(
        FileModel.id == version.file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=403, detail="无权操作此文件")
    
    # 先创建当前版本的备份
    file_path = os.path.join(os.path.dirname(__file__), "..", file.file_path)
    file_path = os.path.abspath(file_path)
    
    if os.path.exists(file_path):
        backup_comment = f"恢复前的自动备份 (目标版本: {version.version_number})"
        manager.create_version(file.id, file_path, current_user.id, backup_comment)
    
    # 执行恢复
    success = manager.restore_version(request.target_version_id, file_path, current_user.id)
    
    if not success:
        raise HTTPException(status_code=500, detail="恢复失败")
    
    return {
        "message": f"已成功恢复到版本 {version.version_number}",
        "restored_version": version.version_number
    }


@router.get("/compare/{version_a_id}/{version_b_id}", response_model=VersionCompareResult)
async def compare_versions(
    version_a_id: int,
    version_b_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """对比两个版本"""
    manager = VersionManager(db)
    
    # 获取版本信息
    version_a = manager.get_version(version_a_id)
    version_b = manager.get_version(version_b_id)
    
    if not version_a or not version_b:
        raise HTTPException(status_code=404, detail="版本不存在")
    
    # 检查权限（两个版本必须属于同一文件，且用户有权限）
    if version_a.file_id != version_b.file_id:
        raise HTTPException(status_code=400, detail="只能对比同一文件的不同版本")
    
    file = db.query(FileModel).filter(
        FileModel.id == version_a.file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=403, detail="无权访问此文件")
    
    return manager.compare_versions(version_a_id, version_b_id)


@router.delete("/{version_id}")
async def delete_version(
    version_id: int,
    hard_delete: bool = False,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """删除版本"""
    manager = VersionManager(db)
    
    version = manager.get_version(version_id)
    if not version:
        raise HTTPException(status_code=404, detail="版本不存在")
    
    # 检查权限
    file = db.query(FileModel).filter(
        FileModel.id == version.file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=403, detail="无权操作此文件")
    
    # 至少保留一个版本
    versions = manager.get_versions(file.id)
    if len(versions) <= 1:
        raise HTTPException(status_code=400, detail="至少需要保留一个版本")
    
    success = manager.delete_version(version_id, soft_delete=not hard_delete)
    
    if not success:
        raise HTTPException(status_code=500, detail="删除失败")
    
    return {
        "message": "版本已删除",
        "version_id": version_id,
        "hard_delete": hard_delete
    }


@router.post("/file/{file_id}/cleanup")
async def cleanup_old_versions(
    file_id: int,
    keep_count: int = 10,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db)
):
    """清理旧版本"""
    file = db.query(FileModel).filter(
        FileModel.id == file_id,
        FileModel.owner_id == current_user.id
    ).first()
    
    if not file:
        raise HTTPException(status_code=404, detail="文件不存在")
    
    manager = VersionManager(db)
    deleted_count = manager.cleanup_old_versions(file_id, keep_count)
    
    return {
        "message": f"已清理 {deleted_count} 个旧版本",
        "deleted_count": deleted_count,
        "kept_count": keep_count
    }

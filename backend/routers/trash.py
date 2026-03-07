"""
Wenxi网盘 - 垃圾桶模块
作者：Wenxi
功能：管理已删除文件，支持恢复、永久删除、自动清理
"""

import os
import asyncio
from datetime import datetime, timedelta
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, BackgroundTasks
from sqlalchemy.orm import Session
from pydantic import BaseModel

from logger import logger
from database import get_db
from models import File as FileModel, User
from routers.auth import get_current_user
from search.service import search_service


router = APIRouter(prefix="/trash", tags=["trash"])


class TrashItemResponse(BaseModel):
    """垃圾桶项目响应"""

    id: int
    filename: str
    original_filename: str
    file_size: int
    mime_type: Optional[str]
    created_at: datetime
    deleted_at: datetime
    days_until_deletion: int  # 距离永久删除的天数


class TrashRestoreResponse(BaseModel):
    """恢复文件响应"""

    message: str
    file_id: int
    restored_path: str


class TrashEmptyResponse(BaseModel):
    """清空垃圾桶响应"""

    message: str
    deleted_count: int


# 配置：文件在垃圾桶中保留的天数
TRASH_RETENTION_DAYS = 30


@router.get("/", response_model=List[TrashItemResponse])
async def get_trash_items(
    current_user: User = Depends(get_current_user), db: Session = Depends(get_db)
):
    """
    Wenxi - 获取垃圾桶列表
    功能：获取当前用户已删除的文件列表
    """
    try:
        files = (
            db.query(FileModel)
            .filter(FileModel.owner_id == current_user.id, FileModel.is_deleted == True)
            .order_by(FileModel.deleted_at.desc())
            .all()
        )

        result = []
        for file in files:
            # 计算距离永久删除的天数
            days_until_deletion = 0
            if file.deleted_at:
                deletion_date = file.deleted_at + timedelta(days=TRASH_RETENTION_DAYS)
                days_until_deletion = max(0, (deletion_date - datetime.now()).days)

            result.append(
                TrashItemResponse(
                    id=file.id,
                    filename=file.filename,
                    original_filename=file.original_filename,
                    file_size=file.file_size,
                    mime_type=file.mime_type,
                    created_at=file.created_at,
                    deleted_at=file.deleted_at,
                    days_until_deletion=days_until_deletion,
                )
            )

        logger.info(
            f"用户 {current_user.username} 获取垃圾桶列表: {len(result)} 个文件"
        )
        return result

    except Exception as e:
        logger.error(f"获取垃圾桶列表失败: {e}")
        raise HTTPException(status_code=500, detail="获取垃圾桶列表失败")


@router.post("/{file_id}/restore", response_model=TrashRestoreResponse)
async def restore_file(
    file_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 恢复文件
    功能：将文件从垃圾桶恢复到原始位置
    """
    try:
        file = (
            db.query(FileModel)
            .filter(
                FileModel.id == file_id,
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == True,
            )
            .first()
        )

        if not file:
            raise HTTPException(status_code=404, detail="文件不存在或不在垃圾桶中")

        # 恢复文件状态
        file.is_deleted = False
        file.deleted_at = None
        restored_path = file.original_path
        file.original_path = None

        db.commit()

        # 重新添加到搜索索引（异步执行）
        try:
            asyncio.create_task(search_service.index_file(file_id, db))
        except Exception as e:
            logger.warning(f"Wenxi - 异步索引恢复的文件失败: {e}")

        logger.info(f"用户 {current_user.username} 恢复文件: {file.original_filename}")

        return TrashRestoreResponse(
            message="文件恢复成功", file_id=file_id, restored_path=restored_path
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"恢复文件失败: {e}")
        raise HTTPException(status_code=500, detail="恢复文件失败")


@router.delete("/{file_id}")
async def permanent_delete_file(
    file_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 永久删除文件
    功能：从垃圾桶永久删除文件（不可恢复）
    """
    try:
        file = (
            db.query(FileModel)
            .filter(
                FileModel.id == file_id,
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == True,
            )
            .first()
        )

        if not file:
            raise HTTPException(status_code=404, detail="文件不存在或不在垃圾桶中")

        # 获取文件路径
        file_path = os.path.join(os.path.dirname(__file__), "..", file.file_path)
        file_path = os.path.abspath(file_path)

        # 删除物理文件
        if os.path.exists(file_path):
            try:
                os.remove(file_path)
                logger.info(f"已删除物理文件: {file_path}")
            except Exception as e:
                logger.warning(f"删除物理文件失败: {file_path}, 错误: {e}")

        # 从搜索索引中移除
        try:
            asyncio.create_task(search_service.remove_file_index(file_id))
        except Exception as e:
            logger.warning(f"Wenxi - 异步移除文件索引失败: {e}")

        # 删除数据库记录
        db.delete(file)
        db.commit()

        logger.info(
            f"用户 {current_user.username} 永久删除文件: {file.original_filename}"
        )

        return {
            "message": "文件已永久删除",
            "file_id": file_id,
            "filename": file.original_filename,
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"永久删除文件失败: {e}")
        raise HTTPException(status_code=500, detail="永久删除文件失败")


@router.delete("/", response_model=TrashEmptyResponse)
async def empty_trash(
    current_user: User = Depends(get_current_user), db: Session = Depends(get_db)
):
    """
    Wenxi - 清空垃圾桶
    功能：永久删除垃圾桶中的所有文件
    """
    try:
        # 获取所有已删除的文件
        files = (
            db.query(FileModel)
            .filter(FileModel.owner_id == current_user.id, FileModel.is_deleted == True)
            .all()
        )

        deleted_count = 0
        for file in files:
            # 获取文件路径
            file_path = os.path.join(os.path.dirname(__file__), "..", file.file_path)
            file_path = os.path.abspath(file_path)

            # 删除物理文件
            if os.path.exists(file_path):
                try:
                    os.remove(file_path)
                    logger.info(f"已删除物理文件: {file_path}")
                except Exception as e:
                    logger.warning(f"删除物理文件失败: {file_path}, 错误: {e}")

            # 从搜索索引中移除
            try:
                asyncio.create_task(search_service.remove_file_index(file.id))
            except Exception as e:
                logger.warning(f"Wenxi - 异步移除文件索引失败: {e}")

            # 删除数据库记录
            db.delete(file)
            deleted_count += 1

        db.commit()

        logger.info(
            f"用户 {current_user.username} 清空垃圾桶: 删除 {deleted_count} 个文件"
        )

        return TrashEmptyResponse(message="垃圾桶已清空", deleted_count=deleted_count)

    except Exception as e:
        logger.error(f"清空垃圾桶失败: {e}")
        raise HTTPException(status_code=500, detail="清空垃圾桶失败")


@router.post("/cleanup-expired", response_model=TrashEmptyResponse)
async def cleanup_expired_files(
    current_user: User = Depends(get_current_user), db: Session = Depends(get_db)
):
    """
    Wenxi - 清理过期文件
    功能：删除超过保留期的文件（通常由系统自动调用）
    """
    try:
        # 计算过期时间
        expiration_date = datetime.now() - timedelta(days=TRASH_RETENTION_DAYS)

        # 获取所有过期的已删除文件
        files = (
            db.query(FileModel)
            .filter(
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == True,
                FileModel.deleted_at <= expiration_date,
            )
            .all()
        )

        deleted_count = 0
        for file in files:
            # 获取文件路径
            file_path = os.path.join(os.path.dirname(__file__), "..", file.file_path)
            file_path = os.path.abspath(file_path)

            # 删除物理文件
            if os.path.exists(file_path):
                try:
                    os.remove(file_path)
                    logger.info(f"已删除过期物理文件: {file_path}")
                except Exception as e:
                    logger.warning(f"删除物理文件失败: {file_path}, 错误: {e}")

            # 从搜索索引中移除
            try:
                asyncio.create_task(search_service.remove_file_index(file.id))
            except Exception as e:
                logger.warning(f"Wenxi - 异步移除文件索引失败: {e}")

            # 删除数据库记录
            db.delete(file)
            deleted_count += 1

        db.commit()

        if deleted_count > 0:
            logger.info(
                f"用户 {current_user.username} 清理过期文件: 删除 {deleted_count} 个文件"
            )

        return TrashEmptyResponse(
            message=f"已清理 {deleted_count} 个过期文件", deleted_count=deleted_count
        )

    except Exception as e:
        logger.error(f"清理过期文件失败: {e}")
        raise HTTPException(status_code=500, detail="清理过期文件失败")


@router.get("/stats")
async def get_trash_stats(
    current_user: User = Depends(get_current_user), db: Session = Depends(get_db)
):
    """
    Wenxi - 获取垃圾桶统计信息
    功能：获取垃圾桶中的文件数量和总大小
    """
    try:
        files = (
            db.query(FileModel)
            .filter(FileModel.owner_id == current_user.id, FileModel.is_deleted == True)
            .all()
        )

        total_size = sum(file.file_size for file in files)
        file_count = len(files)

        # 计算即将过期（7天内）的文件数量
        expiration_threshold = datetime.now() - timedelta(days=TRASH_RETENTION_DAYS - 7)
        expiring_soon = sum(
            1
            for file in files
            if file.deleted_at and file.deleted_at <= expiration_threshold
        )

        return {
            "file_count": file_count,
            "total_size": total_size,
            "expiring_soon": expiring_soon,
            "retention_days": TRASH_RETENTION_DAYS,
        }

    except Exception as e:
        logger.error(f"获取垃圾桶统计失败: {e}")
        raise HTTPException(status_code=500, detail="获取垃圾桶统计失败")

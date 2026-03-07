"""
Wenxi网盘 - 文件夹管理模块
作者：Wenxi
功能：支持文件夹的创建、管理、嵌套结构和文件移动
"""

from datetime import datetime
from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from pydantic import BaseModel, Field

from logger import logger
from database import get_db
from models import Folder, File as FileModel, User
from routers.auth import get_current_user

router = APIRouter()


class FolderCreateRequest(BaseModel):
    """创建文件夹请求"""

    name: str = Field(..., min_length=1, max_length=255, description="文件夹名称")
    parent_id: Optional[int] = Field(None, description="父文件夹ID，为None表示根目录")


class FolderUpdateRequest(BaseModel):
    """更新文件夹请求"""

    name: str = Field(..., min_length=1, max_length=255, description="新文件夹名称")


class FolderMoveRequest(BaseModel):
    """移动文件夹请求"""

    target_parent_id: Optional[int] = Field(
        None, description="目标父文件夹ID，为None表示根目录"
    )


class FolderResponse(BaseModel):
    """文件夹响应"""

    id: int
    name: str
    parent_id: Optional[int]
    user_id: int
    created_at: datetime
    updated_at: datetime
    full_path: str
    file_count: int = 0
    children_count: int = 0

    class Config:
        from_attributes = True


class FolderTreeResponse(BaseModel):
    """文件夹树形结构响应"""

    id: int
    name: str
    parent_id: Optional[int]
    children: List["FolderTreeResponse"] = []


class FileMoveRequest(BaseModel):
    """移动文件到文件夹请求"""

    target_folder_id: Optional[int] = Field(
        None, description="目标文件夹ID，为None表示根目录"
    )


class FolderContentsResponse(BaseModel):
    """文件夹内容响应"""

    folder: Optional[FolderResponse]
    files: list
    subfolders: List[FolderResponse]
    breadcrumbs: List[dict]


# 递归引用处理
FolderTreeResponse.model_rebuild()


def get_folder_breadcrumbs(db: Session, folder: Optional[Folder]) -> List[dict]:
    """获取文件夹面包屑导航"""
    breadcrumbs = []
    if folder:
        current = folder
        while current:
            breadcrumbs.insert(0, {"id": current.id, "name": current.name})
            current = current.parent
    return breadcrumbs


def check_circular_reference(
    db: Session, folder_id: int, target_parent_id: int
) -> bool:
    """检查是否会形成循环引用"""
    if folder_id == target_parent_id:
        return True

    current = db.query(Folder).filter(Folder.id == target_parent_id).first()
    while current:
        if current.parent_id == folder_id:
            return True
        current = current.parent
    return False


@router.post("", response_model=FolderResponse)
async def create_folder(
    request: FolderCreateRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 创建文件夹
    功能：创建新文件夹，支持嵌套结构
    """
    try:
        # 检查父文件夹是否存在
        if request.parent_id:
            parent_folder = (
                db.query(Folder)
                .filter(
                    Folder.id == request.parent_id, Folder.user_id == current_user.id
                )
                .first()
            )
            if not parent_folder:
                raise HTTPException(status_code=404, detail="父文件夹不存在")

        # 检查同一父目录下是否已有同名文件夹
        existing = (
            db.query(Folder)
            .filter(
                Folder.name == request.name,
                Folder.parent_id == request.parent_id,
                Folder.user_id == current_user.id,
            )
            .first()
        )

        if existing:
            raise HTTPException(status_code=400, detail="该位置已存在同名文件夹")

        # 创建文件夹
        new_folder = Folder(
            name=request.name, parent_id=request.parent_id, user_id=current_user.id
        )

        db.add(new_folder)
        db.commit()
        db.refresh(new_folder)

        logger.info(f"用户 {current_user.username} 创建文件夹: {request.name}")

        # 计算统计信息
        file_count = (
            db.query(FileModel)
            .filter(FileModel.folder_id == new_folder.id, FileModel.is_deleted == False)
            .count()
        )

        children_count = (
            db.query(Folder).filter(Folder.parent_id == new_folder.id).count()
        )

        return FolderResponse(
            id=new_folder.id,
            name=new_folder.name,
            parent_id=new_folder.parent_id,
            user_id=new_folder.user_id,
            created_at=new_folder.created_at,
            updated_at=new_folder.updated_at,
            full_path=new_folder.get_full_path(),
            file_count=file_count,
            children_count=children_count,
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"创建文件夹失败: {e}")
        raise HTTPException(status_code=500, detail="创建文件夹失败")


@router.get("/tree", response_model=List[FolderTreeResponse])
async def get_folder_tree(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 获取文件夹树形结构
    功能：获取用户的完整文件夹树
    """
    try:

        def build_tree(parent_id: Optional[int]) -> List[FolderTreeResponse]:
            folders = (
                db.query(Folder)
                .filter(
                    Folder.user_id == current_user.id, Folder.parent_id == parent_id
                )
                .all()
            )

            return [
                FolderTreeResponse(
                    id=folder.id,
                    name=folder.name,
                    parent_id=folder.parent_id,
                    children=build_tree(folder.id),
                )
                for folder in folders
            ]

        return build_tree(None)

    except Exception as e:
        logger.error(f"获取文件夹树失败: {e}")
        raise HTTPException(status_code=500, detail="获取文件夹树失败")


@router.get("/list", response_model=List[FolderResponse])
async def list_folders(
    parent_id: Optional[int] = None,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 获取文件夹列表
    功能：获取指定父文件夹下的所有子文件夹
    """
    try:
        # 验证父文件夹权限
        if parent_id:
            parent = (
                db.query(Folder)
                .filter(Folder.id == parent_id, Folder.user_id == current_user.id)
                .first()
            )
            if not parent:
                raise HTTPException(status_code=404, detail="父文件夹不存在")

        folders = (
            db.query(Folder)
            .filter(Folder.user_id == current_user.id, Folder.parent_id == parent_id)
            .all()
        )

        result = []
        for folder in folders:
            file_count = (
                db.query(FileModel)
                .filter(FileModel.folder_id == folder.id, FileModel.is_deleted == False)
                .count()
            )

            children_count = (
                db.query(Folder).filter(Folder.parent_id == folder.id).count()
            )

            result.append(
                FolderResponse(
                    id=folder.id,
                    name=folder.name,
                    parent_id=folder.parent_id,
                    user_id=folder.user_id,
                    created_at=folder.created_at,
                    updated_at=folder.updated_at,
                    full_path=folder.get_full_path(),
                    file_count=file_count,
                    children_count=children_count,
                )
            )

        return result

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"获取文件夹列表失败: {e}")
        raise HTTPException(status_code=500, detail="获取文件夹列表失败")


@router.get("/{folder_id}/contents", response_model=FolderContentsResponse)
async def get_folder_contents(
    folder_id: int,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 获取文件夹内容
    功能：获取文件夹内的所有文件和子文件夹
    """
    try:
        folder = (
            db.query(Folder)
            .filter(Folder.id == folder_id, Folder.user_id == current_user.id)
            .first()
        )

        if not folder:
            raise HTTPException(status_code=404, detail="文件夹不存在")

        # 获取子文件夹
        subfolders = (
            db.query(Folder)
            .filter(Folder.parent_id == folder_id, Folder.user_id == current_user.id)
            .all()
        )

        # 获取文件夹内的文件
        files = (
            db.query(FileModel)
            .filter(
                FileModel.folder_id == folder_id,
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == False,
            )
            .all()
        )

        # 获取面包屑导航
        breadcrumbs = get_folder_breadcrumbs(db, folder)

        return FolderContentsResponse(
            folder=FolderResponse(
                id=folder.id,
                name=folder.name,
                parent_id=folder.parent_id,
                user_id=folder.user_id,
                created_at=folder.created_at,
                updated_at=folder.updated_at,
                full_path=folder.get_full_path(),
                file_count=len(files),
                children_count=len(subfolders),
            ),
            files=[
                {
                    "id": f.id,
                    "filename": f.filename,
                    "original_filename": f.original_filename,
                    "file_size": f.file_size,
                    "mime_type": f.mime_type,
                    "created_at": f.created_at,
                    "is_shared": f.is_shared,
                }
                for f in files
            ],
            subfolders=[
                FolderResponse(
                    id=f.id,
                    name=f.name,
                    parent_id=f.parent_id,
                    user_id=f.user_id,
                    created_at=f.created_at,
                    updated_at=f.updated_at,
                    full_path=f.get_full_path(),
                    file_count=db.query(FileModel)
                    .filter(FileModel.folder_id == f.id, FileModel.is_deleted == False)
                    .count(),
                    children_count=db.query(Folder)
                    .filter(Folder.parent_id == f.id)
                    .count(),
                )
                for f in subfolders
            ],
            breadcrumbs=breadcrumbs,
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"获取文件夹内容失败: {e}")
        raise HTTPException(status_code=500, detail="获取文件夹内容失败")


@router.get("/root/contents", response_model=FolderContentsResponse)
async def get_root_contents(
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 获取根目录内容
    功能：获取根目录下的所有文件和文件夹
    """
    try:
        # 获取根目录下的文件夹
        subfolders = (
            db.query(Folder)
            .filter(Folder.parent_id == None, Folder.user_id == current_user.id)
            .all()
        )

        # 获取根目录下的文件
        files = (
            db.query(FileModel)
            .filter(
                FileModel.folder_id == None,
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == False,
            )
            .all()
        )

        return FolderContentsResponse(
            folder=None,
            files=[
                {
                    "id": f.id,
                    "filename": f.filename,
                    "original_filename": f.original_filename,
                    "file_size": f.file_size,
                    "mime_type": f.mime_type,
                    "created_at": f.created_at,
                    "is_shared": f.is_shared,
                }
                for f in files
            ],
            subfolders=[
                FolderResponse(
                    id=f.id,
                    name=f.name,
                    parent_id=f.parent_id,
                    user_id=f.user_id,
                    created_at=f.created_at,
                    updated_at=f.updated_at,
                    full_path=f.get_full_path(),
                    file_count=db.query(FileModel)
                    .filter(FileModel.folder_id == f.id, FileModel.is_deleted == False)
                    .count(),
                    children_count=db.query(Folder)
                    .filter(Folder.parent_id == f.id)
                    .count(),
                )
                for f in subfolders
            ],
            breadcrumbs=[],
        )

    except Exception as e:
        logger.error(f"获取根目录内容失败: {e}")
        raise HTTPException(status_code=500, detail="获取根目录内容失败")


@router.put("/{folder_id}", response_model=FolderResponse)
async def rename_folder(
    folder_id: int,
    request: FolderUpdateRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 重命名文件夹
    功能：修改文件夹名称
    """
    try:
        folder = (
            db.query(Folder)
            .filter(Folder.id == folder_id, Folder.user_id == current_user.id)
            .first()
        )

        if not folder:
            raise HTTPException(status_code=404, detail="文件夹不存在")

        # 检查同一父目录下是否已有同名文件夹
        existing = (
            db.query(Folder)
            .filter(
                Folder.name == request.name,
                Folder.parent_id == folder.parent_id,
                Folder.user_id == current_user.id,
                Folder.id != folder_id,
            )
            .first()
        )

        if existing:
            raise HTTPException(status_code=400, detail="该位置已存在同名文件夹")

        folder.name = request.name
        db.commit()
        db.refresh(folder)

        logger.info(
            f"用户 {current_user.username} 重命名文件夹: {folder_id} -> {request.name}"
        )

        file_count = (
            db.query(FileModel)
            .filter(FileModel.folder_id == folder.id, FileModel.is_deleted == False)
            .count()
        )

        children_count = db.query(Folder).filter(Folder.parent_id == folder.id).count()

        return FolderResponse(
            id=folder.id,
            name=folder.name,
            parent_id=folder.parent_id,
            user_id=folder.user_id,
            created_at=folder.created_at,
            updated_at=folder.updated_at,
            full_path=folder.get_full_path(),
            file_count=file_count,
            children_count=children_count,
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"重命名文件夹失败: {e}")
        raise HTTPException(status_code=500, detail="重命名文件夹失败")


@router.put("/{folder_id}/move", response_model=FolderResponse)
async def move_folder(
    folder_id: int,
    request: FolderMoveRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 移动文件夹
    功能：将文件夹移动到新的父文件夹
    """
    try:
        folder = (
            db.query(Folder)
            .filter(Folder.id == folder_id, Folder.user_id == current_user.id)
            .first()
        )

        if not folder:
            raise HTTPException(status_code=404, detail="文件夹不存在")

        # 检查目标父文件夹
        if request.target_parent_id:
            target_folder = (
                db.query(Folder)
                .filter(
                    Folder.id == request.target_parent_id,
                    Folder.user_id == current_user.id,
                )
                .first()
            )

            if not target_folder:
                raise HTTPException(status_code=404, detail="目标文件夹不存在")

            # 检查循环引用
            if check_circular_reference(db, folder_id, request.target_parent_id):
                raise HTTPException(
                    status_code=400, detail="不能将文件夹移动到自身或其子文件夹中"
                )

        # 检查目标位置是否已有同名文件夹
        existing = (
            db.query(Folder)
            .filter(
                Folder.name == folder.name,
                Folder.parent_id == request.target_parent_id,
                Folder.user_id == current_user.id,
                Folder.id != folder_id,
            )
            .first()
        )

        if existing:
            raise HTTPException(status_code=400, detail="目标位置已存在同名文件夹")

        folder.parent_id = request.target_parent_id
        db.commit()
        db.refresh(folder)

        logger.info(
            f"用户 {current_user.username} 移动文件夹: {folder_id} -> parent {request.target_parent_id}"
        )

        file_count = (
            db.query(FileModel)
            .filter(FileModel.folder_id == folder.id, FileModel.is_deleted == False)
            .count()
        )

        children_count = db.query(Folder).filter(Folder.parent_id == folder.id).count()

        return FolderResponse(
            id=folder.id,
            name=folder.name,
            parent_id=folder.parent_id,
            user_id=folder.user_id,
            created_at=folder.created_at,
            updated_at=folder.updated_at,
            full_path=folder.get_full_path(),
            file_count=file_count,
            children_count=children_count,
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"移动文件夹失败: {e}")
        raise HTTPException(status_code=500, detail="移动文件夹失败")


@router.delete("/{folder_id}")
async def delete_folder(
    folder_id: int,
    move_to_root: bool = True,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 删除文件夹
    功能：删除文件夹，可选择将内部文件移至根目录或一起删除
    """
    try:
        folder = (
            db.query(Folder)
            .filter(Folder.id == folder_id, Folder.user_id == current_user.id)
            .first()
        )

        if not folder:
            raise HTTPException(status_code=404, detail="文件夹不存在")

        # 检查是否有子文件夹
        children_count = db.query(Folder).filter(Folder.parent_id == folder_id).count()

        if children_count > 0:
            raise HTTPException(status_code=400, detail="请先删除子文件夹")

        # 处理文件夹内的文件
        if move_to_root:
            # 将文件移至根目录
            db.query(FileModel).filter(FileModel.folder_id == folder_id).update(
                {"folder_id": None}
            )

        db.delete(folder)
        db.commit()

        logger.info(f"用户 {current_user.username} 删除文件夹: {folder_id}")

        return {
            "message": "文件夹已删除",
            "folder_id": folder_id,
            "files_moved_to_root": move_to_root,
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"删除文件夹失败: {e}")
        raise HTTPException(status_code=500, detail="删除文件夹失败")


@router.post("/files/{file_id}/move")
async def move_file_to_folder(
    file_id: int,
    request: FileMoveRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 移动文件到文件夹
    功能：将文件移动到指定文件夹或根目录
    """
    try:
        file = (
            db.query(FileModel)
            .filter(
                FileModel.id == file_id,
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == False,
            )
            .first()
        )

        if not file:
            raise HTTPException(status_code=404, detail="文件不存在")

        # 检查目标文件夹
        if request.target_folder_id:
            folder = (
                db.query(Folder)
                .filter(
                    Folder.id == request.target_folder_id,
                    Folder.user_id == current_user.id,
                )
                .first()
            )

            if not folder:
                raise HTTPException(status_code=404, detail="目标文件夹不存在")

        file.folder_id = request.target_folder_id
        db.commit()
        db.refresh(file)

        logger.info(
            f"用户 {current_user.username} 移动文件 {file_id} 到文件夹 {request.target_folder_id}"
        )

        return {
            "message": "文件移动成功",
            "file_id": file_id,
            "target_folder_id": request.target_folder_id,
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"移动文件失败: {e}")
        raise HTTPException(status_code=500, detail="移动文件失败")


@router.post("/files/batch-move")
async def batch_move_files(
    file_ids: List[int],
    request: FileMoveRequest,
    current_user: User = Depends(get_current_user),
    db: Session = Depends(get_db),
):
    """
    Wenxi - 批量移动文件
    功能：将多个文件批量移动到指定文件夹
    """
    try:
        # 检查目标文件夹
        if request.target_folder_id:
            folder = (
                db.query(Folder)
                .filter(
                    Folder.id == request.target_folder_id,
                    Folder.user_id == current_user.id,
                )
                .first()
            )

            if not folder:
                raise HTTPException(status_code=404, detail="目标文件夹不存在")

        # 验证所有文件所有权
        files = (
            db.query(FileModel)
            .filter(
                FileModel.id.in_(file_ids),
                FileModel.owner_id == current_user.id,
                FileModel.is_deleted == False,
            )
            .all()
        )

        if len(files) != len(file_ids):
            raise HTTPException(status_code=404, detail="部分文件不存在或无权访问")

        # 批量更新
        db.query(FileModel).filter(
            FileModel.id.in_(file_ids), FileModel.owner_id == current_user.id
        ).update({"folder_id": request.target_folder_id})

        db.commit()

        logger.info(
            f"用户 {current_user.username} 批量移动 {len(file_ids)} 个文件到文件夹 {request.target_folder_id}"
        )

        return {
            "message": f"成功移动 {len(file_ids)} 个文件",
            "file_ids": file_ids,
            "target_folder_id": request.target_folder_id,
        }

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"批量移动文件失败: {e}")
        raise HTTPException(status_code=500, detail="批量移动文件失败")

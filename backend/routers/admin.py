"""
Wenxi网盘 - 管理员API路由
作者：Wenxi
功能：提供用户管理、角色分配、权限管理等管理员功能

API列表：
- GET    /admin/users              获取用户列表
- GET    /admin/users/{id}         获取用户详情
- POST   /admin/users              创建用户
- PUT    /admin/users/{id}         更新用户信息
- DELETE /admin/users/{id}         删除用户
- PUT    /admin/users/{id}/role    修改用户角色
- PUT    /admin/users/{id}/permissions  修改用户权限
- GET    /admin/roles              获取角色列表
- GET    /admin/permissions        获取权限列表
- GET    /admin/stats              获取系统统计
"""

from typing import List, Optional
from datetime import datetime, timezone

from fastapi import APIRouter, Depends, HTTPException, status, Query
from pydantic import BaseModel, EmailStr
from sqlalchemy.orm import Session
from sqlalchemy import func

import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from database import get_db
from models import User, File
from routers.auth import get_current_user, get_password_hash
from permissions import (
    Permission,
    Role,
    PermissionChecker,
    get_all_permissions_info,
    get_permission_dict,
)
from decorators.permissions import (
    require_admin,
    require_super_admin,
    AdminRequired,
    SuperAdminRequired,
)
from logger import logger

router = APIRouter()


# ========== Pydantic模型定义 ==========


class UserCreate(BaseModel):
    """创建用户请求模型"""

    username: str
    email: EmailStr
    password: str
    role: int = Role.USER
    is_active: bool = True


class UserUpdate(BaseModel):
    """更新用户请求模型"""

    username: Optional[str] = None
    email: Optional[EmailStr] = None
    is_active: Optional[bool] = None


class UserRoleUpdate(BaseModel):
    """更新用户角色请求模型"""

    role: int


class UserPermissionsUpdate(BaseModel):
    """更新用户权限请求模型"""

    permissions: int
    mode: str = "set"  # set: 设置, add: 添加, remove: 移除


class UserResponse(BaseModel):
    """用户响应模型"""

    id: int
    username: str
    email: str
    role: int
    role_name: str
    permissions: int
    permission_list: List[str]
    is_active: bool
    created_at: datetime
    updated_at: datetime
    file_count: int

    class Config:
        from_attributes = True


class RoleInfo(BaseModel):
    """角色信息模型"""

    id: int
    name: str
    display_name: str
    description: str
    permissions: int
    permission_list: List[str]
    is_system: bool
    is_editable: bool


class PermissionInfo(BaseModel):
    """权限信息模型"""

    value: int
    name: str
    display_name: str
    icon: str
    category: str


class SystemStats(BaseModel):
    """系统统计模型"""

    total_users: int
    total_files: int
    total_storage: int  # 字节
    active_users: int
    inactive_users: int
    users_by_role: dict


# ========== 辅助函数 ==========


def get_user_response(user: User, db: Session) -> dict:
    """构建用户响应数据"""
    role_info = Role.get_role(user.role)
    file_count = (
        db.query(File)
        .filter(File.owner_id == user.id, File.is_deleted == False)
        .count()
    )

    return {
        "id": user.id,
        "username": user.username,
        "email": user.email,
        "role": user.role,
        "role_name": role_info["display_name"] if role_info else "未知角色",
        "permissions": user.permissions,
        "permission_list": PermissionChecker.get_permission_list(user.permissions),
        "is_active": user.is_active,
        "created_at": user.created_at,
        "updated_at": user.updated_at,
        "file_count": file_count,
    }


# ========== API路由 ==========


@router.get("/users", response_model=List[UserResponse])
async def list_users(
    skip: int = Query(0, ge=0, description="跳过的记录数"),
    limit: int = Query(20, ge=1, le=100, description="返回的记录数"),
    role: Optional[int] = Query(None, description="按角色筛选"),
    is_active: Optional[bool] = Query(None, description="按状态筛选"),
    search: Optional[str] = Query(None, description="搜索用户名或邮箱"),
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """
    获取用户列表（管理员）

    支持筛选和搜索：
    - role: 按角色ID筛选
    - is_active: 按激活状态筛选
    - search: 搜索用户名或邮箱
    """
    logger.info(f"管理员 {current_user.username} 查询用户列表")

    query = db.query(User)

    # 应用筛选
    if role is not None:
        query = query.filter(User.role == role)
    if is_active is not None:
        query = query.filter(User.is_active == is_active)
    if search:
        search_term = f"%{search}%"
        query = query.filter(
            (User.username.ilike(search_term)) | (User.email.ilike(search_term))
        )

    # 超级管理员可以看到所有用户，普通管理员不能看到超级管理员
    if not current_user.is_super_admin():
        query = query.filter(User.role > Role.SUPER_ADMIN)

    # 排序和分页
    users = query.order_by(User.created_at.desc()).offset(skip).limit(limit).all()

    return [get_user_response(user, db) for user in users]


@router.get("/users/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: int,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """获取用户详情（管理员）"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查：普通管理员不能查看超级管理员
    if user.role == Role.SUPER_ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权查看此用户信息")

    logger.info(f"管理员 {current_user.username} 查看用户 {user.username} 详情")
    return get_user_response(user, db)


@router.post("/users", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(
    user_data: UserCreate,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """
    创建新用户（管理员）

    可以指定用户的角色和初始权限
    """
    logger.info(f"管理员 {current_user.username} 尝试创建用户: {user_data.username}")

    # 检查权限：普通管理员不能创建超级管理员或管理员
    if user_data.role <= Role.ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权创建管理员账户")

    # 检查用户名是否已存在
    if db.query(User).filter(User.username == user_data.username).first():
        raise HTTPException(status_code=400, detail="用户名已存在")

    # 检查邮箱是否已存在
    if db.query(User).filter(User.email == user_data.email).first():
        raise HTTPException(status_code=400, detail="邮箱已存在")

    # 创建用户
    hashed_password = get_password_hash(user_data.password)
    default_permissions = Role.get_permissions(user_data.role)

    new_user = User(
        username=user_data.username,
        email=user_data.email,
        hashed_password=hashed_password,
        role=user_data.role,
        permissions=int(default_permissions),
        is_active=user_data.is_active,
    )

    db.add(new_user)
    db.commit()
    db.refresh(new_user)

    logger.info(
        f"管理员 {current_user.username} 成功创建用户: {new_user.username} (ID: {new_user.id})"
    )
    return get_user_response(new_user, db)


@router.put("/users/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: int,
    user_data: UserUpdate,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """更新用户信息（管理员）"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查：普通管理员不能修改超级管理员
    if user.role == Role.SUPER_ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权修改此用户")

    # 不能修改自己的信息通过此接口（使用普通用户接口）
    if user.id == current_user.id:
        raise HTTPException(status_code=400, detail="不能通过管理员接口修改自己的信息")

    # 更新字段
    if user_data.username is not None:
        # 检查新用户名是否已存在
        existing = (
            db.query(User)
            .filter(User.username == user_data.username, User.id != user_id)
            .first()
        )
        if existing:
            raise HTTPException(status_code=400, detail="用户名已存在")
        user.username = user_data.username

    if user_data.email is not None:
        # 检查新邮箱是否已存在
        existing = (
            db.query(User)
            .filter(User.email == user_data.email, User.id != user_id)
            .first()
        )
        if existing:
            raise HTTPException(status_code=400, detail="邮箱已存在")
        user.email = user_data.email

    if user_data.is_active is not None:
        user.is_active = user_data.is_active

    db.commit()
    db.refresh(user)

    logger.info(f"管理员 {current_user.username} 更新用户 {user.username} 信息")
    return get_user_response(user, db)


@router.delete("/users/{user_id}")
async def delete_user(
    user_id: int,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """删除用户（管理员）"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查
    if user.id == current_user.id:
        raise HTTPException(status_code=400, detail="不能删除自己的账户")

    if user.role == Role.SUPER_ADMIN:
        raise HTTPException(status_code=403, detail="不能删除超级管理员账户")

    if user.role == Role.ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="只有超级管理员可以删除管理员账户")

    # 获取用户的文件
    user_files = db.query(File).filter(File.owner_id == user.id).all()

    # 删除物理文件
    for file_record in user_files:
        try:
            file_path = os.path.join(
                os.path.dirname(__file__), "..", file_record.file_path
            )
            if os.path.exists(file_path):
                os.remove(file_path)
        except Exception as e:
            logger.error(f"删除文件失败: {file_record.file_path}, 错误: {e}")

    # 删除数据库记录
    db.query(File).filter(File.owner_id == user.id).delete()
    db.delete(user)
    db.commit()

    logger.info(
        f"管理员 {current_user.username} 删除用户 {user.username} (ID: {user_id})"
    )
    return {"message": f"用户 {user.username} 已成功删除"}


@router.put("/users/{user_id}/role", response_model=UserResponse)
async def update_user_role(
    user_id: int,
    role_data: UserRoleUpdate,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """
    修改用户角色（管理员）

    会自动更新用户的权限为对应角色的默认权限
    """
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查
    if user.id == current_user.id:
        raise HTTPException(status_code=400, detail="不能修改自己的角色")

    # 只有超级管理员可以分配超级管理员或管理员角色
    if role_data.role <= Role.ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权分配管理员角色")

    # 不能修改超级管理员的角色
    if user.role == Role.SUPER_ADMIN:
        raise HTTPException(status_code=403, detail="不能修改超级管理员的角色")

    # 获取角色的默认权限
    new_permissions = Role.get_permissions(role_data.role)

    # 更新用户
    user.role = role_data.role
    user.permissions = int(new_permissions)

    db.commit()
    db.refresh(user)

    logger.info(
        f"管理员 {current_user.username} 将用户 {user.username} 角色改为 {Role.get_role(role_data.role)['display_name']}"
    )
    return get_user_response(user, db)


@router.put("/users/{user_id}/permissions", response_model=UserResponse)
async def update_user_permissions(
    user_id: int,
    perm_data: UserPermissionsUpdate,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """
    修改用户权限（管理员）

    mode参数：
    - set: 设置权限（覆盖原有）
    - add: 添加权限
    - remove: 移除权限
    """
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查
    if user.role == Role.SUPER_ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权修改超级管理员的权限")

    # 只有超级管理员可以授予管理权限
    admin_perms = (
        Permission.MANAGE_USERS | Permission.MANAGE_ROLES | Permission.SYSTEM_SETTINGS
    )
    if (perm_data.permissions & admin_perms) and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权授予管理权限")

    # 应用权限变更
    if perm_data.mode == "set":
        user.permissions = perm_data.permissions
    elif perm_data.mode == "add":
        user.permissions |= perm_data.permissions
    elif perm_data.mode == "remove":
        user.permissions &= ~perm_data.permissions
    else:
        raise HTTPException(status_code=400, detail="无效的操作模式")

    db.commit()
    db.refresh(user)

    logger.info(
        f"管理员 {current_user.username} 更新用户 {user.username} 权限 (模式: {perm_data.mode})"
    )
    return get_user_response(user, db)


@router.get("/roles", response_model=List[RoleInfo])
async def list_roles(current_user: User = Depends(AdminRequired)):
    """获取角色列表"""
    roles = []
    for role_id, role in Role.ROLES.items():
        roles.append(
            {
                "id": role_id,
                "name": role["name"],
                "display_name": role["display_name"],
                "description": role["description"],
                "permissions": int(role["permissions"]),
                "permission_list": PermissionChecker.get_permission_list(
                    int(role["permissions"])
                ),
                "is_system": role["is_system"],
                "is_editable": role["is_editable"],
            }
        )
    return roles


@router.get("/permissions", response_model=List[PermissionInfo])
async def list_permissions(
    category: Optional[str] = Query(None, description="按类别筛选"),
    current_user: User = Depends(AdminRequired),
):
    """
    获取权限列表

    类别：file(文件), advanced(高级), admin(管理), system(系统)
    """
    permissions = get_all_permissions_info()

    if category:
        permissions = [p for p in permissions if p["category"] == category]

    return permissions


@router.get("/stats", response_model=SystemStats)
async def get_system_stats(
    current_user: User = Depends(AdminRequired), db: Session = Depends(get_db)
):
    """获取系统统计信息"""
    # 用户统计
    total_users = db.query(User).count()
    active_users = db.query(User).filter(User.is_active == True).count()
    inactive_users = total_users - active_users

    # 按角色统计
    users_by_role = {}
    for role_id, role in Role.ROLES.items():
        count = db.query(User).filter(User.role == role_id).count()
        users_by_role[role["name"]] = count

    # 文件统计
    total_files = db.query(File).filter(File.is_deleted == False).count()
    total_storage = (
        db.query(func.sum(File.file_size)).filter(File.is_deleted == False).scalar()
        or 0
    )

    return {
        "total_users": total_users,
        "total_files": total_files,
        "total_storage": int(total_storage),
        "active_users": active_users,
        "inactive_users": inactive_users,
        "users_by_role": users_by_role,
    }


@router.get("/users/{user_id}/permissions")
async def get_user_permissions_detail(
    user_id: int,
    current_user: User = Depends(AdminRequired),
    db: Session = Depends(get_db),
):
    """获取用户权限详情"""
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(status_code=404, detail="用户不存在")

    # 权限检查
    if user.role == Role.SUPER_ADMIN and not current_user.is_super_admin():
        raise HTTPException(status_code=403, detail="无权查看此用户的权限")

    return {
        "user_id": user.id,
        "username": user.username,
        "role": user.role,
        "role_name": Role.get_role(user.role)["display_name"]
        if Role.get_role(user.role)
        else "未知",
        "permissions": user.permissions,
        "permission_dict": get_permission_dict(user.permissions),
        "permission_list": PermissionChecker.get_permission_list(user.permissions),
    }


# 导入os用于删除文件
import os

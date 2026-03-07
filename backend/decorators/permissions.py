"""
Wenxi网盘 - 权限装饰器模块
作者：Wenxi
功能：提供FastAPI路由的权限控制装饰器

使用示例：
    @router.post("/upload")
    @require_permission(Permission.UPLOAD)
    async def upload_file(current_user: User = Depends(get_current_user)):
        pass

    @router.delete("/admin/users/{user_id}")
    @require_admin()
    async def delete_user(user_id: int, current_user: User = Depends(get_current_user)):
        pass
"""

from functools import wraps
from typing import Callable, List, Optional, Union

from fastapi import HTTPException, status, Depends
from sqlalchemy.orm import Session

import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from permissions import Permission, PermissionChecker, Role
from routers.auth import get_current_user
from models import User


class PermissionDenied(HTTPException):
    """权限拒绝异常"""

    def __init__(self, detail: str = "权限不足，无法访问此资源"):
        super().__init__(status_code=status.HTTP_403_FORBIDDEN, detail=detail)


def require_permission(permission: Permission):
    """
    要求特定权限的装饰器

    Args:
        permission: 要求的权限

    使用示例：
        @router.post("/upload")
        @require_permission(Permission.UPLOAD)
        async def upload_file(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            # 获取当前用户
            current_user = kwargs.get("current_user")
            if not current_user:
                # 尝试从参数中查找User类型的参数
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查权限
            if not current_user.has_permission(permission):
                raise PermissionDenied(f"需要权限: {permission.name}")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


def require_any_permission(*permissions: Permission):
    """
    要求任意一个权限的装饰器（满足其一即可）

    Args:
        *permissions: 权限列表

    使用示例：
        @router.get("/files")
        @require_any_permission(Permission.VIEW_OTHERS, Permission.DOWNLOAD)
        async def list_files(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            current_user = kwargs.get("current_user")
            if not current_user:
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查是否拥有任意权限
            if not PermissionChecker.check_any(
                current_user.permissions, list(permissions)
            ):
                perm_names = ", ".join([p.name for p in permissions])
                raise PermissionDenied(f"需要以下任意权限之一: {perm_names}")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


def require_all_permissions(*permissions: Permission):
    """
    要求同时拥有所有权限的装饰器

    Args:
        *permissions: 权限列表

    使用示例：
        @router.post("/admin/action")
        @require_all_permissions(Permission.MANAGE_USERS, Permission.MANAGE_ROLES)
        async def admin_action(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            current_user = kwargs.get("current_user")
            if not current_user:
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查是否同时拥有所有权限
            if not PermissionChecker.check_all(
                current_user.permissions, list(permissions)
            ):
                perm_names = ", ".join([p.name for p in permissions])
                raise PermissionDenied(f"需要同时拥有以下权限: {perm_names}")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


def require_admin():
    """
    要求管理员角色的装饰器

    使用示例：
        @router.get("/admin/dashboard")
        @require_admin()
        async def admin_dashboard(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            current_user = kwargs.get("current_user")
            if not current_user:
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查是否为管理员
            if not current_user.is_admin():
                raise PermissionDenied("需要管理员权限")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


def require_super_admin():
    """
    要求超级管理员角色的装饰器

    使用示例：
        @router.post("/admin/system-settings")
        @require_super_admin()
        async def update_system_settings(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            current_user = kwargs.get("current_user")
            if not current_user:
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查是否为超级管理员
            if not current_user.is_super_admin():
                raise PermissionDenied("需要超级管理员权限")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


def require_role(role_id: int):
    """
    要求特定角色的装饰器

    Args:
        role_id: 要求的角色ID

    使用示例：
        @router.get("/special")
        @require_role(Role.USER)
        async def special_endpoint(...):
            pass
    """

    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(*args, **kwargs):
            current_user = kwargs.get("current_user")
            if not current_user:
                for arg in args:
                    if isinstance(arg, User):
                        current_user = arg
                        break
                if not current_user:
                    for value in kwargs.values():
                        if isinstance(value, User):
                            current_user = value
                            break

            if not current_user:
                raise HTTPException(
                    status_code=status.HTTP_401_UNAUTHORIZED, detail="未提供用户信息"
                )

            # 检查角色
            if current_user.role != role_id:
                role_info = Role.get_role(role_id)
                role_name = role_info["display_name"] if role_info else f"角色{role_id}"
                raise PermissionDenied(f"需要{role_name}身份")

            return await func(*args, **kwargs)

        return wrapper

    return decorator


# 依赖注入函数（用于FastAPI的Depends）

check_permission = require_permission


def PermissionRequired(permission: Permission):
    """
    FastAPI依赖函数 - 检查单个权限

    使用示例：
        @router.post("/upload")
        async def upload_file(
            current_user: User = Depends(get_current_user),
            _: None = Depends(PermissionRequired(Permission.UPLOAD))
        ):
            pass
    """

    def checker(current_user: User = Depends(get_current_user)):
        if not current_user.has_permission(permission):
            raise PermissionDenied(f"需要权限: {permission.name}")
        return current_user

    return checker


def AdminRequired(current_user: User = Depends(get_current_user)):
    """
    FastAPI依赖函数 - 要求管理员

    使用示例：
        @router.get("/admin/dashboard")
        async def admin_dashboard(
            current_user: User = Depends(get_current_user),
            _: None = Depends(AdminRequired)
        ):
            pass
    """
    if not current_user.is_admin():
        raise PermissionDenied("需要管理员权限")
    return current_user


def SuperAdminRequired(current_user: User = Depends(get_current_user)):
    """
    FastAPI依赖函数 - 要求超级管理员

    使用示例：
        @router.post("/admin/system-settings")
        async def update_settings(
            current_user: User = Depends(get_current_user),
            _: None = Depends(SuperAdminRequired)
        ):
            pass
    """
    if not current_user.is_super_admin():
        raise PermissionDenied("需要超级管理员权限")
    return current_user

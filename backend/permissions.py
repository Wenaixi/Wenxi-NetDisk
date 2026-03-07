"""
Wenxi网盘 - 权限管理系统 v1.0
作者：Wenxi
功能：基于位掩码的细粒度权限控制

权限设计原理：
- 每个权限对应一个2的幂次方位
- 用户权限值 = 所有拥有权限的位或运算结果
- 检查权限：(user_permissions & permission) != 0

角色权限矩阵：
- 超级管理员(0xFF): 所有权限
- 管理员(0x7F): 除系统设置外的所有权限
- 普通用户(0x0F): 上传、下载、删除自己、分享
- 访客(0x02): 仅下载
"""

from enum import IntFlag, auto
from typing import List, Optional


class Permission(IntFlag):
    """
    权限位掩码定义

    使用IntFlag支持位运算：
    - user.has_permission(Permission.UPLOAD | Permission.DELETE_OWN)
    - user.check_permission(Permission.MANAGE_USERS)
    """

    # 文件操作权限
    UPLOAD = 1 << 0  # 1    - 上传文件
    DOWNLOAD = 1 << 1  # 2    - 下载文件
    DELETE_OWN = 1 << 2  # 4    - 删除自己的文件
    SHARE_FILES = 1 << 3  # 8    - 分享文件

    # 高级文件权限
    VIEW_OTHERS = 1 << 4  # 16   - 查看他人文件
    DELETE_OTHERS = 1 << 5  # 32   - 删除他人文件
    EDIT_OTHERS = 1 << 6  # 64   - 编辑他人文件

    # 管理权限
    MANAGE_USERS = 1 << 7  # 128  - 管理用户（创建、编辑、删除）
    MANAGE_ROLES = 1 << 8  # 256  - 管理角色权限
    VIEW_AUDIT = 1 << 9  # 512  - 查看审计日志

    # 系统权限
    SYSTEM_SETTINGS = 1 << 10  # 1024 - 系统设置
    BACKUP_RESTORE = 1 << 11  # 2048 - 备份和恢复

    # 组合权限（方便使用）
    FILE_BASIC = UPLOAD | DOWNLOAD | DELETE_OWN | SHARE_FILES  # 基本文件操作
    FILE_ALL = FILE_BASIC | VIEW_OTHERS | DELETE_OTHERS | EDIT_OTHERS  # 所有文件操作
    MANAGE_ALL = MANAGE_USERS | MANAGE_ROLES | VIEW_AUDIT  # 所有管理权限
    SYSTEM_ALL = SYSTEM_SETTINGS | BACKUP_RESTORE  # 所有系统权限


class Role:
    """
    预定义角色及其默认权限

    角色ID规则：
    - 0: 超级管理员（系统内置，不可删除）
    - 1: 管理员（系统内置，不可删除）
    - 2: 普通用户（系统内置，不可删除）
    - 3: 访客（系统内置，不可删除）
    - 100+: 自定义角色
    """

    # 角色常量
    SUPER_ADMIN = 0
    ADMIN = 1
    USER = 2
    GUEST = 3

    # 角色定义
    ROLES = {
        SUPER_ADMIN: {
            "name": "super_admin",
            "display_name": "超级管理员",
            "description": "系统最高权限，可进行所有操作",
            "permissions": Permission(0xFFFF),  # 所有权限
            "is_system": True,
            "is_editable": False,
        },
        ADMIN: {
            "name": "admin",
            "display_name": "管理员",
            "description": "系统管理员，除系统设置外的所有权限",
            "permissions": Permission.FILE_ALL
            | Permission.MANAGE_USERS
            | Permission.VIEW_AUDIT,
            "is_system": True,
            "is_editable": False,
        },
        USER: {
            "name": "user",
            "display_name": "普通用户",
            "description": "普通用户，可管理自己的文件",
            "permissions": Permission.FILE_BASIC,
            "is_system": True,
            "is_editable": False,
        },
        GUEST: {
            "name": "guest",
            "display_name": "访客",
            "description": "访客，仅能下载被分享的文件",
            "permissions": Permission.DOWNLOAD,
            "is_system": True,
            "is_editable": False,
        },
    }

    @classmethod
    def get_role(cls, role_id: int) -> Optional[dict]:
        """获取角色定义"""
        return cls.ROLES.get(role_id)

    @classmethod
    def get_role_by_name(cls, role_name: str) -> Optional[dict]:
        """通过名称获取角色"""
        for role_id, role in cls.ROLES.items():
            if role["name"] == role_name:
                return {**role, "id": role_id}
        return None

    @classmethod
    def get_all_roles(cls) -> List[dict]:
        """获取所有角色列表"""
        return [{**role, "id": role_id} for role_id, role in cls.ROLES.items()]

    @classmethod
    def get_permissions(cls, role_id: int) -> Permission:
        """获取角色的权限值"""
        role = cls.ROLES.get(role_id)
        if role:
            return role["permissions"]
        return Permission(0)

    @classmethod
    def has_permission(cls, role_id: int, permission: Permission) -> bool:
        """检查角色是否拥有指定权限"""
        role_permissions = cls.get_permissions(role_id)
        return bool(role_permissions & permission)


class PermissionChecker:
    """
    权限检查器

    使用方法：
    1. 直接检查权限值：
       PermissionChecker.check(user_permissions, Permission.UPLOAD)

    2. 检查多个权限（需同时满足）：
       PermissionChecker.check_all(user_permissions, [Permission.UPLOAD, Permission.DELETE_OWN])

    3. 检查任意权限（满足其一）：
       PermissionChecker.check_any(user_permissions, [Permission.UPLOAD, Permission.DELETE_OWN])
    """

    @staticmethod
    def check(user_permissions: int, permission: Permission) -> bool:
        """
        检查用户是否拥有指定权限

        Args:
            user_permissions: 用户的权限值（整数）
            permission: 要检查的权限

        Returns:
            bool: 是否拥有权限
        """
        return bool(user_permissions & permission)

    @staticmethod
    def check_all(user_permissions: int, permissions: List[Permission]) -> bool:
        """
        检查用户是否同时拥有所有指定权限

        Args:
            user_permissions: 用户的权限值（整数）
            permissions: 权限列表

        Returns:
            bool: 是否同时拥有所有权限
        """
        required = Permission(0)
        for perm in permissions:
            required |= perm
        return (user_permissions & required) == required

    @staticmethod
    def check_any(user_permissions: int, permissions: List[Permission]) -> bool:
        """
        检查用户是否拥有任意一个指定权限

        Args:
            user_permissions: 用户的权限值（整数）
            permissions: 权限列表

        Returns:
            bool: 是否拥有任意一个权限
        """
        for perm in permissions:
            if user_permissions & perm:
                return True
        return False

    @staticmethod
    def get_permission_list(permissions_value: int) -> List[str]:
        """
        将权限值转换为权限名称列表

        Args:
            permissions_value: 权限值（整数）

        Returns:
            List[str]: 权限名称列表
        """
        result = []
        for perm in Permission:
            if permissions_value & perm:
                result.append(perm.name)
        return result

    @staticmethod
    def get_permission_dict(permissions_value: int) -> dict:
        """
        将权限值转换为权限字典

        Args:
            permissions_value: 权限值（整数）

        Returns:
            dict: 权限字典 {permission_name: bool}
        """
        return {
            "upload": bool(permissions_value & Permission.UPLOAD),
            "download": bool(permissions_value & Permission.DOWNLOAD),
            "delete_own": bool(permissions_value & Permission.DELETE_OWN),
            "share_files": bool(permissions_value & Permission.SHARE_FILES),
            "view_others": bool(permissions_value & Permission.VIEW_OTHERS),
            "delete_others": bool(permissions_value & Permission.DELETE_OTHERS),
            "edit_others": bool(permissions_value & Permission.EDIT_OTHERS),
            "manage_users": bool(permissions_value & Permission.MANAGE_USERS),
            "manage_roles": bool(permissions_value & Permission.MANAGE_ROLES),
            "view_audit": bool(permissions_value & Permission.VIEW_AUDIT),
            "system_settings": bool(permissions_value & Permission.SYSTEM_SETTINGS),
            "backup_restore": bool(permissions_value & Permission.BACKUP_RESTORE),
        }


# 快捷函数
def has_permission(user_permissions: int, permission: Permission) -> bool:
    """快捷函数：检查单个权限"""
    return PermissionChecker.check(user_permissions, permission)


def has_any_permission(user_permissions: int, *permissions: Permission) -> bool:
    """快捷函数：检查任意权限"""
    return PermissionChecker.check_any(user_permissions, list(permissions))


def has_all_permissions(user_permissions: int, *permissions: Permission) -> bool:
    """快捷函数：检查所有权限"""
    return PermissionChecker.check_all(user_permissions, list(permissions))


# 权限描述映射（用于前端显示）
PERMISSION_DESCRIPTIONS = {
    Permission.UPLOAD: {"name": "上传文件", "icon": "upload", "category": "file"},
    Permission.DOWNLOAD: {"name": "下载文件", "icon": "download", "category": "file"},
    Permission.DELETE_OWN: {
        "name": "删除自己的文件",
        "icon": "delete",
        "category": "file",
    },
    Permission.SHARE_FILES: {"name": "分享文件", "icon": "share", "category": "file"},
    Permission.VIEW_OTHERS: {
        "name": "查看他人文件",
        "icon": "eye",
        "category": "advanced",
    },
    Permission.DELETE_OTHERS: {
        "name": "删除他人文件",
        "icon": "delete-user",
        "category": "advanced",
    },
    Permission.EDIT_OTHERS: {
        "name": "编辑他人文件",
        "icon": "edit",
        "category": "advanced",
    },
    Permission.MANAGE_USERS: {
        "name": "管理用户",
        "icon": "user-management",
        "category": "admin",
    },
    Permission.MANAGE_ROLES: {
        "name": "管理角色",
        "icon": "role-management",
        "category": "admin",
    },
    Permission.VIEW_AUDIT: {
        "name": "查看审计日志",
        "icon": "audit",
        "category": "admin",
    },
    Permission.SYSTEM_SETTINGS: {
        "name": "系统设置",
        "icon": "settings",
        "category": "system",
    },
    Permission.BACKUP_RESTORE: {
        "name": "备份和恢复",
        "icon": "backup",
        "category": "system",
    },
}


def get_permission_info(permission: Permission) -> dict:
    """获取权限的详细信息"""
    info = PERMISSION_DESCRIPTIONS.get(permission, {})
    return {
        "value": permission.value,
        "name": permission.name,
        "display_name": info.get("name", permission.name),
        "icon": info.get("icon", "default"),
        "category": info.get("category", "other"),
    }


def get_all_permissions_info() -> List[dict]:
    """获取所有权限的详细信息列表"""
    return [get_permission_info(perm) for perm in Permission]


def get_permissions_by_category(category: str) -> List[dict]:
    """按类别获取权限信息"""
    return [
        get_permission_info(perm)
        for perm in Permission
        if PERMISSION_DESCRIPTIONS.get(perm, {}).get("category") == category
    ]
